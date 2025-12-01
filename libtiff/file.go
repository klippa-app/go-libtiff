package libtiff

import (
	"context"
	"errors"
	"io"
	"iter"

	"github.com/klippa-app/go-libtiff/internal/imports"

	"github.com/tetratelabs/wazero/api"
)

type File struct {
	pointer    uint64
	readerFile *imports.File
	instance   *Instance
	closeFunc  func(context.Context) error
}

func (f *File) Close(ctx context.Context) error {
	return f.closeFunc(ctx)
}

func (f *File) Directories(ctx context.Context) iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		n := 0
		for {
			yieldRes := yield(n, nil)
			if !yieldRes {
				break
			}

			n++

			res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFReadDirectory").Call(ctx, f.pointer)
			if err != nil {
				yieldRes = yield(n, err)
				if !yieldRes {
					break
				}
			}

			// No more images in the file.
			if res[0] == 0 {
				break
			}
		}
	}
}

// TIFFCurrentDirectory returns the index of the current directory.
func (f *File) TIFFCurrentDirectory(ctx context.Context) (uint32, error) {
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFCurrentDirectory").Call(ctx, f.pointer)
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(res[0]), nil
}

// TIFFLastDirectory returns whether the current directory is the last directory.
func (f *File) TIFFLastDirectory(ctx context.Context) (bool, error) {
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFLastDirectory").Call(ctx, f.pointer)
	if err != nil {
		return false, err
	}

	// TIFFLastDirectory() returns a non-zero value if the current directory
	// is the last directory in the file; otherwise zero is returned.
	if res[0] == 0 {
		return false, nil
	}

	return true, nil
}

// TIFFReadDirectory will read the next directory.
func (f *File) TIFFReadDirectory(ctx context.Context) error {
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFReadDirectory").Call(ctx, f.pointer)
	if err != nil {
		return err
	}

	if res[0] == 0 {
		return errors.New("could not read directory")
	}

	return nil
}

// TIFFSetDirectory will set and read the given directory.
func (f *File) TIFFSetDirectory(ctx context.Context, n uint32) error {
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFSetDirectory").Call(ctx, f.pointer, api.EncodeU32(n))
	if err != nil {
		return err
	}

	if res[0] == 0 {
		return errors.New("could not set directory")
	}
	return nil
}

// TIFFNumberOfDirectories will return the amount of directories.
func (f *File) TIFFNumberOfDirectories(ctx context.Context) (uint32, error) {
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFNumberOfDirectories").Call(ctx, f.pointer)
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(res[0]), nil
}

// TIFFOpenFile opens a file from a path. Be aware that this is limited to the
// virtual filesystem given to the instance.
func (i *Instance) TIFFOpenFile(ctx context.Context, filePath string) (*File, error) {
	cStringFilePath, err := i.NewCString(ctx, filePath)
	if err != nil {
		return nil, err
	}
	defer cStringFilePath.Free(ctx)

	cStringFileMode, err := i.NewCString(ctx, "r")
	if err != nil {
		return nil, err
	}
	defer cStringFileMode.Free(ctx)

	// Result is a pointer to struct_tiff
	res, err := i.internalInstance.Module.ExportedFunction("TIFFOpen").Call(ctx, cStringFilePath.Pointer, cStringFileMode.Pointer)
	if err != nil {
		return nil, err
	}

	if res[0] == 0 {
		return nil, errors.New("error while opening tiff file")
	}

	return &File{
		pointer:  res[0],
		instance: i,
		closeFunc: func(ctx context.Context) error {
			_, err := i.internalInstance.Module.ExportedFunction("TIFFClose").Call(ctx, res[0])
			if err != nil {
				return err
			}
			return nil
		},
	}, nil
}

// TIFFClientOpen can open a TIFF file from a reader.
func (i *Instance) TIFFClientOpen(ctx context.Context, filename string, reader io.ReadSeeker, fileSize uint64) (*File, error) {
	imports.FileReaders.Mutex.Lock()
	fileReaderIndex := imports.FileReaders.Counter
	imports.FileReaders.Counter++
	imports.FileReaders.Mutex.Unlock()

	paramPointer, err := i.Malloc(ctx, 4)
	if err != nil {
		return nil, err
	}

	cleanupFileReader := func(ctx context.Context) error {
		err = i.Free(ctx, paramPointer)
		if err != nil {
			return err
		}
		return nil
	}

	ok := i.internalInstance.Module.Memory().WriteUint32Le(uint32(paramPointer), fileReaderIndex)
	if !ok {
		cleanupFileReader(ctx)
		return nil, errors.New("could not write file reader param to memory")
	}

	newFileReader := &imports.File{
		ParamPointer: paramPointer,
		FileSize:     fileSize,
		Reader:       reader,
	}

	imports.FileReaders.Mutex.Lock()
	imports.FileReaders.Refs[fileReaderIndex] = newFileReader
	imports.FileReaders.Mutex.Unlock()

	var oldCleanup = cleanupFileReader
	newCleanup := func(ctx context.Context) error {
		err = oldCleanup(ctx)
		if err != nil {
			return err
		}
		imports.FileReaders.Mutex.Lock()
		delete(imports.FileReaders.Refs, fileReaderIndex)
		imports.FileReaders.Mutex.Unlock()
		return nil
	}
	cleanupFileReader = newCleanup

	cStringFileName, err := i.NewCString(ctx, filename)
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}
	defer cStringFileName.Free(ctx)

	cStringFileMode, err := i.NewCString(ctx, "r")
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}
	defer cStringFileMode.Free(ctx)

	// Result is a pointer to struct_tiff.
	res, err := i.internalInstance.Module.ExportedFunction("TIFFClientOpenGo").Call(ctx, cStringFileName.Pointer, cStringFileMode.Pointer, paramPointer)
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	if res[0] == 0 {
		cleanupFileReader(ctx)
		return nil, errors.New("error while opening tiff file")
	}

	newFile := &File{
		pointer:    res[0],
		readerFile: newFileReader,
		instance:   i,
		closeFunc: func(ctx context.Context) error {
			_, err := i.internalInstance.Module.ExportedFunction("TIFFClose").Call(ctx, res[0])
			if err != nil {
				return err
			}

			return cleanupFileReader(ctx)
		},
	}

	return newFile, nil
}
