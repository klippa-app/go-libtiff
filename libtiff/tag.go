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

func (e *TagNotDefinedError) Is(err error) bool {
	if _, ok := err.(*TagNotDefinedError); ok {
		return true
	}
	return false
}

func (f *File) TIFFGetFieldUint16_t(ctx context.Context, tag TIFFTAG) (uint16, error) {
	valuePointer, err := f.instance.malloc(ctx, 2)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldUint16_t", f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()
	defer f.instance.internalInstance.CallLock.Unlock()

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

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldUint32_t", f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()
	defer f.instance.internalInstance.CallLock.Unlock()

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

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldInt", f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()
	defer f.instance.internalInstance.CallLock.Unlock()

	readValue, success := f.instance.internalInstance.Module.Memory().ReadByte(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return int(readValue), nil
}

func (f *File) TIFFGetFieldTwoUint16(ctx context.Context, tag TIFFTAG) (uint16, uint16, error) {
	valuePointer1, err := f.instance.malloc(ctx, 2)
	if err != nil {
		return 0, 0, err
	}
	defer f.instance.free(ctx, valuePointer1)

	valuePointer2, err := f.instance.malloc(ctx, 2)
	if err != nil {
		return 0, 0, err
	}
	defer f.instance.free(ctx, valuePointer2)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldTwoUint16", f.pointer, api.EncodeU32(uint32(tag)), valuePointer1, valuePointer2)
	if err != nil {
		return 0, 0, err
	}

	if results[0] == 0 {
		return 0, 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()
	defer f.instance.internalInstance.CallLock.Unlock()

	readValue1, success := f.instance.internalInstance.Module.Memory().ReadUint16Le(uint32(valuePointer1))
	if !success {
		return 0, 0, errors.New("could not read tag value")
	}

	readValue2, success := f.instance.internalInstance.Module.Memory().ReadUint16Le(uint32(valuePointer2))
	if !success {
		return 0, 0, errors.New("could not read tag value")
	}

	return readValue1, readValue2, nil
}

func (f *File) TIFFGetFieldFloat(ctx context.Context, tag TIFFTAG) (float32, error) {
	valuePointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldFloat", f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()
	defer f.instance.internalInstance.CallLock.Unlock()

	readValue, success := f.instance.internalInstance.Module.Memory().ReadFloat32Le(uint32(valuePointer))
	if !success {
		return 0, errors.New("could not read tag value")
	}

	return readValue, nil
}

func (f *File) TIFFGetFieldDouble(ctx context.Context, tag TIFFTAG) (float64, error) {
	valuePointer, err := f.instance.malloc(ctx, 8)
	if err != nil {
		return 0, err
	}
	defer f.instance.free(ctx, valuePointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldDouble", f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return 0, err
	}

	if results[0] == 0 {
		return 0, &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()
	defer f.instance.internalInstance.CallLock.Unlock()

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

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetFieldConstChar", f.pointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return "", err
	}

	if results[0] == 0 {
		return "", &TagNotDefinedError{
			Tag: tag,
		}
	}

	// Prevent concurrent memory usage.
	f.instance.internalInstance.CallLock.Lock()

	readPointer, success := f.instance.internalInstance.Module.Memory().ReadUint32Le(uint32(valuePointer))
	if !success {
		f.instance.internalInstance.CallLock.Unlock()
		return "", errors.New("could not read tag value")
	}

	f.instance.internalInstance.CallLock.Unlock()

	readValue := f.instance.readCString(readPointer)

	return readValue, nil
}
