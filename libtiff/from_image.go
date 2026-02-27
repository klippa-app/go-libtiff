package libtiff

import (
	"context"
	"image"
	"image/color"
	"strings"
	"time"
)

// AlphaMode controls how the alpha channel is stored in the TIFF file.
type AlphaMode int

const (
	// AlphaAuto auto-detects: *image.NRGBA uses unassociated, everything else uses associated.
	AlphaAuto AlphaMode = iota
	// AlphaAssociated stores premultiplied (associated) alpha.
	AlphaAssociated
	// AlphaUnassociated stores non-premultiplied (unassociated) alpha.
	AlphaUnassociated
)

type FromGoImageOptions struct {
	Compression TIFFTAG
	// Quality sets the compression quality level (1-100). Only used for JPEG
	// compression. If 0, the default quality (75) is used.
	Quality int
	// AlphaMode controls premultiplied vs non-premultiplied alpha storage.
	// Zero value (AlphaAuto) auto-detects based on image type.
	AlphaMode AlphaMode
	// Software sets the TIFFTAG_SOFTWARE tag. If empty, defaults to
	// "go-libtiff/libtiff-{version}" where version is the linked libtiff version.
	Software string
	// DateTime sets the TIFFTAG_DATETIME tag in "YYYY:MM:DD HH:MM:SS" format.
	// If empty, defaults to the current time.
	DateTime string
	// Artist sets the TIFFTAG_ARTIST tag. If empty, the tag is not written.
	Artist string
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

	// Determine alpha mode.
	alphaMode := AlphaAuto
	if options != nil {
		alphaMode = options.AlphaMode
	}
	if alphaMode == AlphaAuto {
		if _, ok := img.(*image.NRGBA); ok {
			alphaMode = AlphaUnassociated
		} else {
			alphaMode = AlphaAssociated
		}
	}

	isJPEG := compression == COMPRESSION_JPEG
	samplesPerPixel := uint16(4)
	if isJPEG {
		samplesPerPixel = 3 // JPEG does not support alpha channel.
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
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_SAMPLESPERPIXEL, samplesPerPixel); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_COMPRESSION, uint16(compression)); err != nil {
		return err
	}
	if isJPEG {
		quality := 75
		if options != nil && options.Quality > 0 {
			quality = options.Quality
		}
		if err := f.TIFFSetFieldInt(ctx, TIFFTAG_JPEGQUALITY, quality); err != nil {
			return err
		}
		// JPEG in TIFF requires YCBCR photometric.
		if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PHOTOMETRIC, uint16(PHOTOMETRIC_YCBCR)); err != nil {
			return err
		}
		// Tell libtiff we provide RGB data and it should convert to YCbCr.
		// Without this, libtiff expects raw subsampled YCbCr data and
		// miscalculates the scanline size.
		if err := f.TIFFSetFieldInt(ctx, TIFFTAG_JPEGCOLORMODE, int(JPEGCOLORMODE_RGB)); err != nil {
			return err
		}
	} else {
		if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PHOTOMETRIC, uint16(PHOTOMETRIC_RGB)); err != nil {
			return err
		}
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_ORIENTATION, uint16(ORIENTATION_TOPLEFT)); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PLANARCONFIG, uint16(PLANARCONFIG_CONTIG)); err != nil {
		return err
	}

	// Set metadata tags.
	software := ""
	if options != nil {
		software = options.Software
	}
	if software == "" {
		software = "go-libtiff"
		if ver, err := f.instance.TIFFGetVersion(ctx); err == nil {
			// TIFFGetVersion returns "LIBTIFF, Version X.Y.Z\n...".
			if i := strings.Index(ver, "Version "); i != -1 {
				v := ver[i+len("Version "):]
				if j := strings.IndexAny(v, "\n\r"); j != -1 {
					v = v[:j]
				}
				software = "go-libtiff/libtiff-" + strings.TrimSpace(v)
			}
		}
	}
	if err := f.TIFFSetFieldString(ctx, TIFFTAG_SOFTWARE, software); err != nil {
		return err
	}

	dateTime := ""
	if options != nil {
		dateTime = options.DateTime
	}
	if dateTime == "" {
		dateTime = time.Now().Format("2006:01:02 15:04:05")
	}
	if err := f.TIFFSetFieldString(ctx, TIFFTAG_DATETIME, dateTime); err != nil {
		return err
	}

	artist := ""
	if options != nil {
		artist = options.Artist
	}
	if artist != "" {
		if err := f.TIFFSetFieldString(ctx, TIFFTAG_ARTIST, artist); err != nil {
			return err
		}
	}

	if !isJPEG {
		extraSample := EXTRASAMPLE_ASSOCALPHA
		if alphaMode == AlphaUnassociated {
			extraSample = EXTRASAMPLE_UNASSALPHA
		}
		if err := f.TIFFSetFieldExtraSamples(ctx, []uint16{uint16(extraSample)}); err != nil {
			return err
		}
	}

	// Get a sensible strip size.
	var rowsPerStrip uint32
	if isJPEG {
		// JPEG requires writing the entire image as a single strip to avoid
		// MCU boundary alignment issues with partial last strips.
		rowsPerStrip = height
	} else {
		var err error
		rowsPerStrip, err = f.TIFFDefaultStripSize(ctx, 0)
		if err != nil {
			return err
		}
	}
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_ROWSPERSTRIP, rowsPerStrip); err != nil {
		return err
	}

	// Write image data strip by strip.
	bytesPerPixel := int(samplesPerPixel)
	bytesPerRow := int(width) * bytesPerPixel
	strip := uint32(0)

	// Fast path: direct pixel access for matching image types (non-JPEG only).
	if rgbaImg, ok := img.(*image.RGBA); ok && !isJPEG && alphaMode == AlphaAssociated {
		return f.writeStrips(ctx, bounds, rowsPerStrip, bytesPerRow, &strip, func(stripData []byte, y, rows int) {
			for row := 0; row < rows; row++ {
				srcY := y + row - bounds.Min.Y
				srcStart := srcY*rgbaImg.Stride + (bounds.Min.X * 4)
				copy(stripData[row*bytesPerRow:], rgbaImg.Pix[srcStart:srcStart+bytesPerRow])
			}
		})
	}
	if nrgbaImg, ok := img.(*image.NRGBA); ok && !isJPEG && alphaMode == AlphaUnassociated {
		return f.writeStrips(ctx, bounds, rowsPerStrip, bytesPerRow, &strip, func(stripData []byte, y, rows int) {
			for row := 0; row < rows; row++ {
				srcY := y + row - bounds.Min.Y
				srcStart := srcY*nrgbaImg.Stride + (bounds.Min.X * 4)
				copy(stripData[row*bytesPerRow:], nrgbaImg.Pix[srcStart:srcStart+bytesPerRow])
			}
		})
	}

	// Generic path.
	if alphaMode == AlphaUnassociated {
		return f.writeStrips(ctx, bounds, rowsPerStrip, bytesPerRow, &strip, func(stripData []byte, y, rows int) {
			for row := 0; row < rows; row++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					c := color.NRGBAModel.Convert(img.At(x, y+row)).(color.NRGBA)
					offset := row*bytesPerRow + (x-bounds.Min.X)*bytesPerPixel
					stripData[offset] = c.R
					stripData[offset+1] = c.G
					stripData[offset+2] = c.B
					if !isJPEG {
						stripData[offset+3] = c.A
					}
				}
			}
		})
	}

	return f.writeStrips(ctx, bounds, rowsPerStrip, bytesPerRow, &strip, func(stripData []byte, y, rows int) {
		for row := 0; row < rows; row++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y+row).RGBA()
				offset := row*bytesPerRow + (x-bounds.Min.X)*bytesPerPixel
				stripData[offset] = uint8(r >> 8)
				stripData[offset+1] = uint8(g >> 8)
				stripData[offset+2] = uint8(b >> 8)
				if !isJPEG {
					stripData[offset+3] = uint8(a >> 8)
				}
			}
		}
	})
}

// writeStrips writes image data strip by strip, calling fillStrip to populate each strip's pixel data,
// then writes the TIFF directory.
func (f *File) writeStrips(ctx context.Context, bounds image.Rectangle, rowsPerStrip uint32, bytesPerRow int, strip *uint32, fillStrip func(stripData []byte, y, rows int)) error {
	for y := bounds.Min.Y; y < bounds.Max.Y; y += int(rowsPerStrip) {
		rows := int(rowsPerStrip)
		if y+rows > bounds.Max.Y {
			rows = bounds.Max.Y - y
		}

		stripData := make([]byte, rows*bytesPerRow)
		fillStrip(stripData, y, rows)

		if err := f.TIFFWriteEncodedStrip(ctx, *strip, stripData); err != nil {
			return err
		}
		*strip++
	}

	return f.TIFFWriteDirectory(ctx)
}
