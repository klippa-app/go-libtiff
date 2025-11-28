package libtiff

import (
	"context"
	"errors"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

type TagNotDefinedError struct {
	Tag TIFFTAG
}

func (e TagNotDefinedError) Error() string {
	return fmt.Sprintf("Tag %d was not found in the curret directory", e.Tag)
}

func (f *File) TIFFGetFieldUint32_t(ctx context.Context, tag TIFFTAG) (uint32, error) {
	valuePointer, err := f.instance.Malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.Free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldUint32_t").Call(ctx, f.Pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	readValue, success := f.instance.internalInstance.Module.Memory().ReadUint32Le(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return readValue, nil
}

func (f *File) TIFFGetFieldFloat(ctx context.Context, tag TIFFTAG) (float32, error) {
	valuePointer, err := f.instance.Malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.Free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldFloat").Call(ctx, f.Pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	readValue, success := f.instance.internalInstance.Module.Memory().ReadFloat32Le(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return readValue, nil
}
