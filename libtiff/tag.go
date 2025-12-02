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
	return fmt.Sprintf("Tag %d was not found in the current directory", e.Tag)
}

func (f *File) TIFFGetFieldUint16_t(ctx context.Context, tag TIFFTAG) (uint16, error) {
	valuePointer, err := f.instance.malloc(ctx, 2)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldUint16_t").Call(ctx, f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	readValue, success := f.instance.internalInstance.Module.Memory().ReadUint16Le(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return readValue, nil
}

func (f *File) TIFFGetFieldUint32_t(ctx context.Context, tag TIFFTAG) (uint32, error) {
	valuePointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldUint32_t").Call(ctx, f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
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

func (f *File) TIFFGetFieldInt(ctx context.Context, tag TIFFTAG) (int, error) {
	valuePointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldInt").Call(ctx, f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	readValue, success := f.instance.internalInstance.Module.Memory().ReadByte(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return int(readValue), nil
}

func (f *File) TIFFGetFieldFloat(ctx context.Context, tag TIFFTAG) (float32, error) {
	valuePointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldFloat").Call(ctx, f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
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

func (f *File) TIFFGetFieldDouble(ctx context.Context, tag TIFFTAG) (float64, error) {
	valuePointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldDouble").Call(ctx, f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	readValue, success := f.instance.internalInstance.Module.Memory().ReadFloat64Le(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return readValue, nil
}

func (f *File) TIFFGetFieldConstChar(ctx context.Context, tag TIFFTAG) (string, error) {
	valuePointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return "", err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFGetFieldConstChar").Call(ctx, f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return "", err
	}

	if results[0] == 0 {
		return "", &TagNotDefinedError{
			Tag: tag,
		}
	}

	readPointer, success := f.instance.internalInstance.Module.Memory().ReadUint32Le(uint32(valuePointer))
	if !success {
		return "", errors.New("could not read tag value")
	}

	readValue := f.instance.readCString(readPointer)

	return readValue, nil
}
