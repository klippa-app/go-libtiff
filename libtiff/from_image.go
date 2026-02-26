package libtiff

import (
	"context"
	"image"
	"image/color"
)

type FromGoImageOptions struct {
	Compression TIFFTAG
}

// FromGoImage writes a Go image to the open TIFF file.
// The file must have been opened for writing via TIFFOpenFileFromReadWriteSeeker.
func (f *File) FromGoImage(ctx context.Context, img image.Image, options *FromGoImageOptions) error {
	bounds := img.Bounds()
	width := uint32(bounds.Dx())
	height := uint32(bounds.Dy())

	compression := COMPRESSION_NONE
	if options != nil && options.Compression != 0 {
		compression = options.Compression
	}

	// Set TIFF tags.
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_IMAGEWIDTH, width); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_IMAGELENGTH, height); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_BITSPERSAMPLE, 8); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_SAMPLESPERPIXEL, 4); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_COMPRESSION, uint16(compression)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PHOTOMETRIC, uint16(PHOTOMETRIC_RGB)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_ORIENTATION, uint16(ORIENTATION_TOPLEFT)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PLANARCONFIG, uint16(PLANARCONFIG_CONTIG)); err != nil {
		return err
	}

	// Set EXTRASAMPLES to indicate the alpha channel is associated alpha.
	if err := f.TIFFSetFieldExtraSamples(ctx, []uint16{uint16(EXTRASAMPLE_ASSOCALPHA)}); err != nil {
		return err
	}

	// Get a sensible strip size.
	rowsPerStrip, err := f.TIFFDefaultStripSize(ctx, 0)
	if err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_ROWSPERSTRIP, rowsPerStrip); err != nil {
		return err
	}

	// Write image data strip by strip.
	bytesPerRow := int(width) * 4
	strip := uint32(0)

	// Optimization: use direct pixel access for *image.RGBA.
	if rgbaImg, ok := img.(*image.RGBA); ok {
		for y := bounds.Min.Y; y < bounds.Max.Y; y += int(rowsPerStrip) {
			rows := int(rowsPerStrip)
			if y+rows > bounds.Max.Y {
				rows = bounds.Max.Y - y
			}

			stripSize := rows * bytesPerRow
			stripData := make([]byte, stripSize)

			for row := 0; row < rows; row++ {
				srcY := y + row - bounds.Min.Y
				srcStart := srcY*rgbaImg.Stride + (bounds.Min.X * 4)
				copy(stripData[row*bytesPerRow:], rgbaImg.Pix[srcStart:srcStart+bytesPerRow])
			}

			if err := f.TIFFWriteEncodedStrip(ctx, strip, stripData); err != nil {
				return err
			}
			strip++
		}
	} else {
		// Generic path: use At() for any image type.
		for y := bounds.Min.Y; y < bounds.Max.Y; y += int(rowsPerStrip) {
			rows := int(rowsPerStrip)
			if y+rows > bounds.Max.Y {
				rows = bounds.Max.Y - y
			}

			stripSize := rows * bytesPerRow
			stripData := make([]byte, stripSize)

			for row := 0; row < rows; row++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, a := img.At(x, y+row).RGBA()
					offset := row*bytesPerRow + (x-bounds.Min.X)*4
					stripData[offset] = uint8(r >> 8)
					stripData[offset+1] = uint8(g >> 8)
					stripData[offset+2] = uint8(b >> 8)
					stripData[offset+3] = uint8(a >> 8)
				}
			}

			if err := f.TIFFWriteEncodedStrip(ctx, strip, stripData); err != nil {
				return err
			}
			strip++
		}
	}

	// Write the directory.
	if err := f.TIFFWriteDirectory(ctx); err != nil {
		return err
	}

	return nil
}

// FromGoImageNRGBA writes a Go image to the open TIFF file using unassociated alpha.
// Use this when the source image has non-premultiplied alpha (e.g. *image.NRGBA).
func (f *File) FromGoImageNRGBA(ctx context.Context, img image.Image, options *FromGoImageOptions) error {
	bounds := img.Bounds()
	width := uint32(bounds.Dx())
	height := uint32(bounds.Dy())

	compression := COMPRESSION_NONE
	if options != nil && options.Compression != 0 {
		compression = options.Compression
	}

	// Set TIFF tags.
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_IMAGEWIDTH, width); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_IMAGELENGTH, height); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_BITSPERSAMPLE, 8); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_SAMPLESPERPIXEL, 4); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_COMPRESSION, uint16(compression)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PHOTOMETRIC, uint16(PHOTOMETRIC_RGB)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_ORIENTATION, uint16(ORIENTATION_TOPLEFT)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PLANARCONFIG, uint16(PLANARCONFIG_CONTIG)); err != nil {
		return err
	}

	// Set EXTRASAMPLES to indicate the alpha channel is unassociated alpha.
	if err := f.TIFFSetFieldExtraSamples(ctx, []uint16{uint16(EXTRASAMPLE_UNASSALPHA)}); err != nil {
		return err
	}

	// Get a sensible strip size.
	rowsPerStrip, err := f.TIFFDefaultStripSize(ctx, 0)
	if err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_ROWSPERSTRIP, rowsPerStrip); err != nil {
		return err
	}

	// Write image data strip by strip using non-premultiplied alpha.
	bytesPerRow := int(width) * 4
	strip := uint32(0)

	// Optimization: use direct pixel access for *image.NRGBA.
	if nrgbaImg, ok := img.(*image.NRGBA); ok {
		for y := bounds.Min.Y; y < bounds.Max.Y; y += int(rowsPerStrip) {
			rows := int(rowsPerStrip)
			if y+rows > bounds.Max.Y {
				rows = bounds.Max.Y - y
			}

			stripSize := rows * bytesPerRow
			stripData := make([]byte, stripSize)

			for row := 0; row < rows; row++ {
				srcY := y + row - bounds.Min.Y
				srcStart := srcY*nrgbaImg.Stride + (bounds.Min.X * 4)
				copy(stripData[row*bytesPerRow:], nrgbaImg.Pix[srcStart:srcStart+bytesPerRow])
			}

			if err := f.TIFFWriteEncodedStrip(ctx, strip, stripData); err != nil {
				return err
			}
			strip++
		}
	} else {
		// Generic path: convert to NRGBA using color model.
		for y := bounds.Min.Y; y < bounds.Max.Y; y += int(rowsPerStrip) {
			rows := int(rowsPerStrip)
			if y+rows > bounds.Max.Y {
				rows = bounds.Max.Y - y
			}

			stripSize := rows * bytesPerRow
			stripData := make([]byte, stripSize)

			for row := 0; row < rows; row++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := color.NRGBAModel.Convert(img.At(x, y+row)).(color.NRGBA)
					offset := row*bytesPerRow + (x-bounds.Min.X)*4
					stripData[offset] = c.R
					stripData[offset+1] = c.G
					stripData[offset+2] = c.B
					stripData[offset+3] = c.A
				}
			}

			if err := f.TIFFWriteEncodedStrip(ctx, strip, stripData); err != nil {
				return err
			}
			strip++
		}
	}

	// Write the directory.
	if err := f.TIFFWriteDirectory(ctx); err != nil {
		return err
	}

	return nil
}
