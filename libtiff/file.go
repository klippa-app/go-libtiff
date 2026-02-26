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

func (f *File) GetError() error {
	if f.readerFile == nil {
		return nil
	}

	if f.readerFile.Error == nil {
		return nil
	}

	return f.readerFile.GetError()
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

			res, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadDirectory", f.pointer)
			if err != nil {
				yieldRes = yield(n, err)
				if !yieldRes {
					break
				}
			}

			err = f.GetError()
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
	res, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFCurrentDirectory", f.pointer)
	if err != nil {
		return 0, err
	}

	err = f.GetError()
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(res[0]), nil
}

// TIFFLastDirectory returns whether the current directory is the last directory.
func (f *File) TIFFLastDirectory(ctx context.Context) (bool, error) {
	res, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFLastDirectory", f.pointer)
	if err != nil {
		return false, err
	}

	err = f.GetError()
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
	res, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadDirectory", f.pointer)
	if err != nil {
		return err
	}

	err = f.GetError()
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
	res, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetDirectory", f.pointer, api.EncodeU32(n))
	if err != nil {
		return err
	}

	err = f.GetError()
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
	res, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFNumberOfDirectories", f.pointer)
	if err != nil {
		return 0, err
	}

	err = f.GetError()
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(res[0]), nil
}

type OpenOptions struct {
	MaxSingleMemAlloc    *int32
	MaxCumulatedMemAlloc *int32
	WarnHandler          func(module string, message string)
	WarnAboutUnknownTags *bool
}

// TIFFOpenFileFromPath opens a file from a path. Be aware that this is limited to the
// virtual filesystem given to the instance.
func (i *Instance) TIFFOpenFileFromPath(ctx context.Context, filePath string, options *OpenOptions) (*File, error) {
	imports.FileReaders.Mutex.Lock()
	fileReaderIndex := imports.FileReaders.Counter
	imports.FileReaders.Counter++
	imports.FileReaders.Mutex.Unlock()

	paramPointer, err := i.malloc(ctx, 4)
	if err != nil {
		return nil, err
	}

	cleanupFileReader := func(ctx context.Context) error {
		err := i.free(ctx, paramPointer)
		if err != nil {
			return err
		}
		return nil
	}

	// Prevent concurrent memory usage.
	i.internalInstance.CallLock.Lock()
	ok := i.internalInstance.Module.Memory().WriteUint32Le(uint32(paramPointer), fileReaderIndex)
	if !ok {
		i.internalInstance.CallLock.Unlock()
		cleanupFileReader(ctx)
		return nil, errors.New("could not write file reader param to memory")
	}
	i.internalInstance.CallLock.Unlock()

	newFileReader := &imports.File{
		ParamPointer: paramPointer,
	}

	imports.FileReaders.Mutex.Lock()
	imports.FileReaders.Refs[fileReaderIndex] = newFileReader
	imports.FileReaders.Mutex.Unlock()

	var oldCleanup = cleanupFileReader
	newCleanup := func(ctx context.Context) error {
		err := oldCleanup(ctx)
		if err != nil {
			return err
		}
		imports.FileReaders.Mutex.Lock()
		delete(imports.FileReaders.Refs, fileReaderIndex)
		imports.FileReaders.Mutex.Unlock()
		return nil
	}
	cleanupFileReader = newCleanup

	cStringFilePath, err := i.newCString(ctx, filePath)
	if err != nil {
		return nil, err
	}
	defer cStringFilePath.Free(ctx)

	cStringFileMode, err := i.newCString(ctx, "r")
	if err != nil {
		return nil, err
	}
	defer cStringFileMode.Free(ctx)

	TIFFOpenOptionsAlloc, err := i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsAlloc")
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	if TIFFOpenOptionsAlloc[0] == 0 {
		cleanupFileReader(ctx)
		return nil, errors.New("error while allocating tiff file options")
	}
	TIFFOpenOptionsAllocPointer := TIFFOpenOptionsAlloc[0]

	var oldCleanup2 = cleanupFileReader
	newCleanup = func(ctx context.Context) error {
		err := oldCleanup2(ctx)
		if err != nil {
			return err
		}

		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsFree", TIFFOpenOptionsAllocPointer)
		return err
	}
	cleanupFileReader = newCleanup

	if options != nil && options.MaxSingleMemAlloc != nil {
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetMaxSingleMemAlloc", TIFFOpenOptionsAllocPointer, api.EncodeI32(*options.MaxSingleMemAlloc))
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	if options != nil && options.MaxCumulatedMemAlloc != nil {
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetMaxCumulatedMemAlloc", TIFFOpenOptionsAllocPointer, api.EncodeI32(*options.MaxCumulatedMemAlloc))
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	if options != nil && options.WarnAboutUnknownTags != nil {
		value := int32(0)
		if *options.WarnAboutUnknownTags {
			value = int32(1)
		}
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetWarnAboutUnknownTags", TIFFOpenOptionsAllocPointer, api.EncodeI32(value))
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	if options != nil && options.WarnHandler != nil {
		newFileReader.WarnHandler = options.WarnHandler
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetWarningHandlerExtRGo", TIFFOpenOptionsAllocPointer, paramPointer)
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetErrorHandlerExtRGo", TIFFOpenOptionsAllocPointer, paramPointer)
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	// Result is a pointer to struct_tiff
	res, err := i.internalInstance.CallExportedFunction(ctx, "TIFFOpenExtGo", cStringFilePath.Pointer, cStringFileMode.Pointer, TIFFOpenOptionsAllocPointer)
	if err != nil {
		return nil, err
	}

	err = newFileReader.GetError()
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	if res[0] == 0 {
		return nil, errors.New("error while opening tiff file")
	}
	filePointer := res[0]

	return &File{
		pointer:  filePointer,
		instance: i,
		closeFunc: func(ctx context.Context) error {
			_, err := i.internalInstance.CallExportedFunction(ctx, "TIFFClose", filePointer)
			if err != nil {
				return err
			}

			return cleanupFileReader(ctx)
		},
	}, nil
}

type fakeReadWriteSeeker struct {
	io.ReadSeeker
}

func (*fakeReadWriteSeeker) Write(p []byte) (n int, err error) {
	return 0, errors.New("the given reader can't be written to")
}

// TIFFOpenFileFromReader can open a TIFF file from a ReadSeeker.
// The filename property is for meaningful errors/warning and the TIFFFileName
// method, it's not required to enter the actual filename.
// The fileSize is for some validation checks and memory allocation limits, the
// fileSize is not absolutely required, but the file might not always be opened
// correctly if the fileSize is not given.
func (i *Instance) TIFFOpenFileFromReader(ctx context.Context, filename string, reader io.ReadSeeker, fileSize uint64, options *OpenOptions) (*File, error) {
	return i.TIFFOpenFileFromReadWriteSeeker(ctx, filename, &fakeReadWriteSeeker{
		ReadSeeker: reader,
	}, fileSize, options)
}

// TIFFOpenFileFromReadWriteSeeker can open a TIFF file from a ReadWriteSeeker.
// The filename property is for meaningful errors/warning and the TIFFFileName
// method, it's not required to enter the actual filename.
// The fileSize is for some validation checks and memory allocation limits, the
// fileSize is not absolutely required, but the file might not always be opened
// correctly if the fileSize is not given.
func (i *Instance) TIFFOpenFileFromReadWriteSeeker(ctx context.Context, filename string, readWriteSeeker io.ReadWriteSeeker, fileSize uint64, options *OpenOptions) (*File, error) {
	imports.FileReaders.Mutex.Lock()
	fileReaderIndex := imports.FileReaders.Counter
	imports.FileReaders.Counter++
	imports.FileReaders.Mutex.Unlock()

	paramPointer, err := i.malloc(ctx, 4)
	if err != nil {
		return nil, err
	}

	cleanupFileReader := func(ctx context.Context) error {
		err := i.free(ctx, paramPointer)
		if err != nil {
			return err
		}
		return nil
	}

	// Prevent concurrent memory usage.
	i.internalInstance.CallLock.Lock()
	ok := i.internalInstance.Module.Memory().WriteUint32Le(uint32(paramPointer), fileReaderIndex)
	if !ok {
		i.internalInstance.CallLock.Unlock()
		cleanupFileReader(ctx)
		return nil, errors.New("could not write file reader param to memory")
	}
	i.internalInstance.CallLock.Unlock()

	newFileReader := &imports.File{
		ParamPointer:    paramPointer,
		FileSize:        fileSize,
		ReadWriteSeeker: readWriteSeeker,
	}

	imports.FileReaders.Mutex.Lock()
	imports.FileReaders.Refs[fileReaderIndex] = newFileReader
	imports.FileReaders.Mutex.Unlock()

	var oldCleanup = cleanupFileReader
	newCleanup := func(ctx context.Context) error {
		err := oldCleanup(ctx)
		if err != nil {
			return err
		}
		imports.FileReaders.Mutex.Lock()
		delete(imports.FileReaders.Refs, fileReaderIndex)
		imports.FileReaders.Mutex.Unlock()
		return nil
	}
	cleanupFileReader = newCleanup

	cStringFileName, err := i.newCString(ctx, filename)
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}
	defer cStringFileName.Free(ctx)

	cStringFileMode, err := i.newCString(ctx, "r")
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}
	defer cStringFileMode.Free(ctx)

	TIFFOpenOptionsAlloc, err := i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsAlloc")
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	if TIFFOpenOptionsAlloc[0] == 0 {
		cleanupFileReader(ctx)
		return nil, errors.New("error while allocating tiff file options")
	}
	TIFFOpenOptionsAllocPointer := TIFFOpenOptionsAlloc[0]

	var oldCleanup2 = cleanupFileReader
	newCleanup = func(ctx context.Context) error {
		err := oldCleanup2(ctx)
		if err != nil {
			return err
		}

		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsFree", TIFFOpenOptionsAllocPointer)
		return err
	}
	cleanupFileReader = newCleanup

	if options != nil && options.MaxSingleMemAlloc != nil {
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetMaxSingleMemAlloc", TIFFOpenOptionsAllocPointer, api.EncodeI32(*options.MaxSingleMemAlloc))
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	if options != nil && options.MaxCumulatedMemAlloc != nil {
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetMaxCumulatedMemAlloc", TIFFOpenOptionsAllocPointer, api.EncodeI32(*options.MaxCumulatedMemAlloc))
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	if options != nil && options.WarnAboutUnknownTags != nil {
		value := int32(0)
		if *options.WarnAboutUnknownTags {
			value = int32(1)
		}
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetWarnAboutUnknownTags", TIFFOpenOptionsAllocPointer, api.EncodeI32(value))
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	if options != nil && options.WarnHandler != nil {
		newFileReader.WarnHandler = options.WarnHandler
		_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetWarningHandlerExtRGo", TIFFOpenOptionsAllocPointer, paramPointer)
		if err != nil {
			cleanupFileReader(ctx)
			return nil, err
		}
	}

	_, err = i.internalInstance.CallExportedFunction(ctx, "TIFFOpenOptionsSetErrorHandlerExtRGo", TIFFOpenOptionsAllocPointer, paramPointer)
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	// Result is a pointer to struct_tiff.
	res, err := i.internalInstance.CallExportedFunction(ctx, "TIFFClientOpenExtGo", cStringFileName.Pointer, cStringFileMode.Pointer, paramPointer, TIFFOpenOptionsAllocPointer)
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	err = newFileReader.GetError()
	if err != nil {
		cleanupFileReader(ctx)
		return nil, err
	}

	if res[0] == 0 {
		cleanupFileReader(ctx)
		return nil, errors.New("error while opening tiff file")
	}

	filePointer := res[0]

	newFile := &File{
		pointer:    filePointer,
		readerFile: newFileReader,
		instance:   i,
		closeFunc: func(ctx context.Context) error {
			_, err := i.internalInstance.CallExportedFunction(ctx, "TIFFClose", filePointer)
			if err != nil {
				return err
			}

			return cleanupFileReader(ctx)
		},
	}

	return newFile, nil
}
