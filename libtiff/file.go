package libtiff

import (
	"context"
	"errors"
	"io"
	"iter"

	"github.com/tetratelabs/wazero/api"
)

type File struct {
	Pointer   uint64
	instance  *Instance
	closeFunc func(context.Context) error
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

			res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFReadDirectory").Call(ctx, f.Pointer)
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
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFCurrentDirectory").Call(ctx, f.Pointer)
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(res[0]), nil
}

// TIFFLastDirectory returns whether the current directory is the last directory.
func (f *File) TIFFLastDirectory(ctx context.Context) (bool, error) {
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFLastDirectory").Call(ctx, f.Pointer)
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
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFReadDirectory").Call(ctx, f.Pointer)
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
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFSetDirectory").Call(ctx, f.Pointer, api.EncodeU32(n))
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
	res, err := f.instance.internalInstance.Module.ExportedFunction("TIFFNumberOfDirectories").Call(ctx, f.Pointer)
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
		Pointer:  res[0],
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
func (i *Instance) TIFFClientOpen(ctx context.Context, filename string, reader io.ReadSeeker) (*File, error) {
	cStringFileName, err := i.NewCString(ctx, filename)
	if err != nil {
		return nil, err
	}
	defer cStringFileName.Free(ctx)

	cStringFileMode, err := i.NewCString(ctx, "r")
	if err != nil {
		return nil, err
	}
	defer cStringFileMode.Free(ctx)

	// Result is a pointer to struct_tiff
	res, err := i.internalInstance.Module.ExportedFunction("TIFFClientOpenGo").Call(ctx, cStringFileName.Pointer, cStringFileMode.Pointer, cStringFileMode.Pointer)
	if err != nil {
		return nil, err
	}

	if res[0] == 0 {
		return nil, errors.New("error while opening tiff file")
	}

	return &File{
		Pointer:  res[0],
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
