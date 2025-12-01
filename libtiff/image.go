package libtiff

import (
	"context"
	"errors"
	"image"

	"github.com/tetratelabs/wazero/api"
)

// ToImage convert the current directory in the open TIFF file to RGBA, the
// caller is responsible for closing since the returned cleanup function will
// free the allocated memory.
func (f *File) ToImage(ctx context.Context) (image.Image, func(context.Context) error, error) {
	width, height, err := f.GetDimensions(ctx)
	if err != nil {
		return nil, nil, err
	}

	img := &image.RGBA{}
	img.Rect = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: width, Y: height}}
	img.Stride = img.Rect.Max.X * 4
	nBytes := img.Rect.Max.X * img.Rect.Max.Y * 4
	imagePointer, err := f.instance.Malloc(ctx, uint64(nBytes))
	if err != nil {
		return nil, nil, err
	}

	cleanupFunc := func(ctx context.Context) error {
		return f.instance.Free(ctx, imagePointer)
	}

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFReadRGBAImageOriented").Call(ctx, f.pointer, api.EncodeU32(uint32(img.Rect.Max.X)), api.EncodeU32(uint32(img.Rect.Max.Y)), imagePointer, api.EncodeU32(uint32(ORIENTATION_TOPLEFT)), 0)
	if err != nil {
		cleanupErr := cleanupFunc(ctx)
		return nil, nil, errors.Join(err, cleanupErr)
	}

	if results[0] != 1 {
		cleanupErr := cleanupFunc(ctx)
		return nil, nil, errors.Join(errors.New("error while converting tiff to RGBA"), cleanupErr)
	}

	// We directly open a view on the image data in the Wazero memory so that
	// we don't have to do any image copying.
	memoryView, ok := f.instance.internalInstance.Module.Memory().Read(uint32(imagePointer), uint32(nBytes))
	if !ok {
		cleanupErr := cleanupFunc(ctx)
		return nil, nil, errors.Join(errors.New("memory view not found"), cleanupErr)
	}

	img.Pix = memoryView

	return img, cleanupFunc, nil
}
