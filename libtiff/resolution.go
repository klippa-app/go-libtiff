package libtiff

import (
	"context"
)

// GetResolution returns the DPI of the current directory.
func (f *File) GetResolution(ctx context.Context) (float32, float32, error) {
	x, err := f.TIFFGetFieldFloat(ctx, TIFFTAG_XRESOLUTION)
	if err != nil {
		return 0, 0, err
	}

	y, err := f.TIFFGetFieldFloat(ctx, TIFFTAG_YRESOLUTION)
	if err != nil {
		return 0, 0, err
	}

	return x, y, nil
}
