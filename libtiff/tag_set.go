package libtiff

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

func (f *File) TIFFSetFieldUint16_t(ctx context.Context, tag TIFFTAG, val uint16) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldUint16_t", f.pointer, api.EncodeU32(uint32(tag)), api.EncodeU32(uint32(val)))
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}

func (f *File) TIFFSetFieldUint32_t(ctx context.Context, tag TIFFTAG, val uint32) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldUint32_t", f.pointer, api.EncodeU32(uint32(tag)), api.EncodeU32(val))
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}

func (f *File) TIFFSetFieldInt(ctx context.Context, tag TIFFTAG, val int) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldInt", f.pointer, api.EncodeU32(uint32(tag)), api.EncodeI32(int32(val)))
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}

func (f *File) TIFFSetFieldFloat(ctx context.Context, tag TIFFTAG, val float32) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldFloat", f.pointer, api.EncodeU32(uint32(tag)), api.EncodeF32(val))
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}

func (f *File) TIFFSetFieldDouble(ctx context.Context, tag TIFFTAG, val float64) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldDouble", f.pointer, api.EncodeU32(uint32(tag)), api.EncodeF64(val))
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}

func (f *File) TIFFSetFieldTwoUint16(ctx context.Context, tag TIFFTAG, val1, val2 uint16) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldTwoUint16", f.pointer, api.EncodeU32(uint32(tag)), api.EncodeU32(uint32(val1)), api.EncodeU32(uint32(val2)))
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}

func (f *File) TIFFSetFieldString(ctx context.Context, tag TIFFTAG, val string) error {
	cString, err := f.instance.newCString(ctx, val)
	if err != nil {
		return err
	}
	defer cString.Free(ctx)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldString", f.pointer, api.EncodeU32(uint32(tag)), cString.Pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set tag value")
	}

	return nil
}
