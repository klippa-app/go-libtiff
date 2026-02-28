package libtiff

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

// TIFFField represents a TIFF tag field descriptor from the C library.
// It provides metadata about a tag such as its name, data type, and read/write counts.
type TIFFField struct {
	pointer  uint64
	instance *Instance
}

// TIFFGetTagListCount returns the number of tags set in the current directory.
func (f *File) TIFFGetTagListCount(ctx context.Context) (int, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetTagListCount", f.pointer)
	if err != nil {
		return 0, err
	}

	return int(api.DecodeI32(results[0])), nil
}

// TIFFGetTagListEntry returns the tag number at the given index in the current directory's tag list.
func (f *File) TIFFGetTagListEntry(ctx context.Context, tagIndex int) (uint32, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFGetTagListEntry", f.pointer, api.EncodeI32(int32(tagIndex)))
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// TIFFFieldWithTag returns the field descriptor for the given tag number.
// Returns nil if the tag is not known.
func (f *File) TIFFFieldWithTag(ctx context.Context, tag uint32) (*TIFFField, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldWithTag", f.pointer, api.EncodeU32(tag))
	if err != nil {
		return nil, err
	}

	if results[0] == 0 {
		return nil, nil
	}

	return &TIFFField{
		pointer:  results[0],
		instance: f.instance,
	}, nil
}

// TIFFFieldWithName returns the field descriptor for the given tag name.
// Returns nil if the tag name is not known.
func (f *File) TIFFFieldWithName(ctx context.Context, name string) (*TIFFField, error) {
	cStr, err := f.instance.newCString(ctx, name)
	if err != nil {
		return nil, err
	}
	defer cStr.Free(ctx)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldWithName", f.pointer, cStr.Pointer)
	if err != nil {
		return nil, err
	}

	if results[0] == 0 {
		return nil, nil
	}

	return &TIFFField{
		pointer:  results[0],
		instance: f.instance,
	}, nil
}

// Name returns the human-readable name of this field.
func (tf *TIFFField) Name(ctx context.Context) (string, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldName", tf.pointer)
	if err != nil {
		return "", err
	}

	if results[0] == 0 {
		return "", errors.New("could not get field name")
	}

	return tf.instance.readCString(uint32(results[0])), nil
}

// DataType returns the data type of this field.
func (tf *TIFFField) DataType(ctx context.Context) (TIFFDataType, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldDataType", tf.pointer)
	if err != nil {
		return 0, err
	}

	return TIFFDataType(api.DecodeU32(results[0])), nil
}

// Tag returns the tag number of this field.
func (tf *TIFFField) Tag(ctx context.Context) (uint32, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldTag", tf.pointer)
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// IsAnonymous returns true if this field is an unknown/anonymous tag.
func (tf *TIFFField) IsAnonymous(ctx context.Context) (bool, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldIsAnonymous", tf.pointer)
	if err != nil {
		return false, err
	}

	return results[0] != 0, nil
}

// PassCount returns true if the field requires a count parameter when reading/writing.
func (tf *TIFFField) PassCount(ctx context.Context) (bool, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldPassCount", tf.pointer)
	if err != nil {
		return false, err
	}

	return results[0] != 0, nil
}

// ReadCount returns the number of values to read for this field.
// Special values: TIFF_VARIABLE (-1), TIFF_VARIABLE2 (-3), TIFF_SPP (-2).
func (tf *TIFFField) ReadCount(ctx context.Context) (int, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldReadCount", tf.pointer)
	if err != nil {
		return 0, err
	}

	return int(api.DecodeI32(results[0])), nil
}

// WriteCount returns the number of values to write for this field.
// Special values: TIFF_VARIABLE (-1), TIFF_VARIABLE2 (-3), TIFF_SPP (-2).
func (tf *TIFFField) WriteCount(ctx context.Context) (int, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldWriteCount", tf.pointer)
	if err != nil {
		return 0, err
	}

	return int(api.DecodeI32(results[0])), nil
}

// SetGetSize returns the size in bytes of the internal set/get value for this field.
func (tf *TIFFField) SetGetSize(ctx context.Context) (int, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldSetGetSize", tf.pointer)
	if err != nil {
		return 0, err
	}

	return int(api.DecodeI32(results[0])), nil
}

// SetGetCountSize returns the size in bytes of the count parameter for variable-count fields.
func (tf *TIFFField) SetGetCountSize(ctx context.Context) (int, error) {
	results, err := tf.instance.internalInstance.CallExportedFunction(ctx, "TIFFFieldSetGetCountSize", tf.pointer)
	if err != nil {
		return 0, err
	}

	return int(api.DecodeI32(results[0])), nil
}
