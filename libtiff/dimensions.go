package libtiff

import (
	"context"
)

// GetDimensions returns the width and height of the current directory in pixels.
func (f *File) GetDimensions(ctx context.Context) (int, int, error) {
	width, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_IMAGEWIDTH)
	if err != nil {
		return 0, 0, err
	}

	height, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_IMAGELENGTH)
	if err != nil {
		return 0, 0, err
	}

	return int(width), int(height), nil
}
