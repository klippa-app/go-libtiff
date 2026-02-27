package libtiff

import (
	"context"
	"fmt"
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
	// Predictor sets TIFFTAG_PREDICTOR. Only meaningful for LZW and Deflate.
	// If 0, the tag is not set.
	Predictor TIFFTAG
	// XResolution sets TIFFTAG_XRESOLUTION. If <= 0, the tag is not set.
	XResolution float32
	// YResolution sets TIFFTAG_YRESOLUTION. If <= 0, the tag is not set.
	YResolution float32
	// ResolutionUnit sets TIFFTAG_RESOLUTIONUNIT. If 0, the tag is not set.
	ResolutionUnit TIFFTAG
	// Description sets TIFFTAG_IMAGEDESCRIPTION. If empty, the tag is not written.
	Description string
	// Copyright sets TIFFTAG_COPYRIGHT. If empty, the tag is not written.
	Copyright string
	// DocumentName sets TIFFTAG_DOCUMENTNAME. If empty, the tag is not written.
	DocumentName string
	// PageName sets TIFFTAG_PAGENAME. If empty, the tag is not written.
	PageName string
	// HostComputer sets TIFFTAG_HOSTCOMPUTER. If empty, the tag is not written.
	HostComputer string
	// Make sets TIFFTAG_MAKE. If empty, the tag is not written.
	Make string
	// Model sets TIFFTAG_MODEL. If empty, the tag is not written.
	Model string
	// RowsPerStrip overrides the default strip size. If 0, uses TIFFDefaultStripSize.
	// Ignored for JPEG compression and tile-based output.
	RowsPerStrip uint32
	// Orientation sets TIFFTAG_ORIENTATION. If 0, defaults to ORIENTATION_TOPLEFT.
	Orientation TIFFTAG
	// TileWidth sets the tile width for tile-based output. Both TileWidth and
	// TileHeight must be set to enable tiled output. If 0, strip-based output is used.
	TileWidth uint32
	// TileHeight sets the tile height for tile-based output. Both TileWidth and
	// TileHeight must be set to enable tiled output. If 0, strip-based output is used.
	TileHeight uint32
	// BilevelThreshold sets the luminance threshold for CCITT bilevel conversion.
	// Pixels with luminance >= threshold become white, below become black.
	// If 0, defaults to 128.
	BilevelThreshold uint8
	// PageNumber sets TIFFTAG_PAGENUMBER. The first value is the page number
	// (0-based), the second is the total number of pages. Both must be set
	// (non-zero TotalPages) for the tag to be written.
	PageNumber uint16
	// TotalPages is the total number of pages for TIFFTAG_PAGENUMBER.
	TotalPages uint16
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

	// Validate tile options.
	useTiles := false
	if options != nil {
		if (options.TileWidth > 0) != (options.TileHeight > 0) {
			return fmt.Errorf("both TileWidth and TileHeight must be set for tile-based output")
		}
		useTiles = options.TileWidth > 0 && options.TileHeight > 0
	}

	isCCITT := compression == COMPRESSION_CCITTFAX3 || compression == COMPRESSION_CCITTFAX4

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
	if isCCITT {
		samplesPerPixel = 1
	}

	// Set TIFF tags.
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_IMAGEWIDTH, width); err != nil {
		return err
	}
	if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_IMAGELENGTH, height); err != nil {
		return err
	}
	bitsPerSample := uint16(8)
	if isCCITT {
		bitsPerSample = 1
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_BITSPERSAMPLE, bitsPerSample); err != nil {
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
	} else if isCCITT {
		if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PHOTOMETRIC, uint16(PHOTOMETRIC_MINISWHITE)); err != nil {
			return err
		}
		if compression == COMPRESSION_CCITTFAX3 {
			if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_GROUP3OPTIONS, uint32(GROUP3OPT_FILLBITS)); err != nil {
				return err
			}
		}
	} else {
		if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PHOTOMETRIC, uint16(PHOTOMETRIC_RGB)); err != nil {
			return err
		}
	}

	// Set predictor if specified.
	if options != nil && options.Predictor != 0 {
		if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_PREDICTOR, uint16(options.Predictor)); err != nil {
			return err
		}
	}

	// Set orientation.
	orientation := ORIENTATION_TOPLEFT
	if options != nil && options.Orientation != 0 {
		orientation = options.Orientation
	}
	if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_ORIENTATION, uint16(orientation)); err != nil {
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

	// Write optional metadata tags.
	if options != nil {
		optionalStringTags := []struct {
			value string
			tag   TIFFTAG
		}{
			{options.Description, TIFFTAG_IMAGEDESCRIPTION},
			{options.Copyright, TIFFTAG_COPYRIGHT},
			{options.DocumentName, TIFFTAG_DOCUMENTNAME},
			{options.PageName, TIFFTAG_PAGENAME},
			{options.HostComputer, TIFFTAG_HOSTCOMPUTER},
			{options.Make, TIFFTAG_MAKE},
			{options.Model, TIFFTAG_MODEL},
		}
		for _, t := range optionalStringTags {
			if t.value != "" {
				if err := f.TIFFSetFieldString(ctx, t.tag, t.value); err != nil {
					return err
				}
			}
		}
	}

	// Set page number tag.
	if options != nil && options.TotalPages > 0 {
		if err := f.TIFFSetFieldTwoUint16(ctx, TIFFTAG_PAGENUMBER, options.PageNumber, options.TotalPages); err != nil {
			return err
		}
	}

	// Set resolution tags.
	if options != nil {
		if options.XResolution > 0 {
			if err := f.TIFFSetFieldFloat(ctx, TIFFTAG_XRESOLUTION, options.XResolution); err != nil {
				return err
			}
		}
		if options.YResolution > 0 {
			if err := f.TIFFSetFieldFloat(ctx, TIFFTAG_YRESOLUTION, options.YResolution); err != nil {
				return err
			}
		}
		if options.ResolutionUnit != 0 {
			if err := f.TIFFSetFieldUint16_t(ctx, TIFFTAG_RESOLUTIONUNIT, uint16(options.ResolutionUnit)); err != nil {
				return err
			}
		}
	}

	if !isJPEG && !isCCITT {
		extraSample := EXTRASAMPLE_ASSOCALPHA
		if alphaMode == AlphaUnassociated {
			extraSample = EXTRASAMPLE_UNASSALPHA
		}
		if err := f.TIFFSetFieldExtraSamples(ctx, []uint16{uint16(extraSample)}); err != nil {
			return err
		}
	}

	// Set up tile or strip layout.
	if useTiles {
		tileWidth := options.TileWidth
		tileHeight := options.TileHeight
		if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_TILEWIDTH, tileWidth); err != nil {
			return err
		}
		if err := f.TIFFSetFieldUint32_t(ctx, TIFFTAG_TILELENGTH, tileHeight); err != nil {
			return err
		}
	} else {
		// Get a sensible strip size.
		var rowsPerStrip uint32
		if isJPEG {
			// JPEG requires writing the entire image as a single strip to avoid
			// MCU boundary alignment issues with partial last strips.
			rowsPerStrip = height
		} else if options != nil && options.RowsPerStrip > 0 {
			rowsPerStrip = options.RowsPerStrip
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
	}

	// CCITT bilevel output.
	if isCCITT {
		threshold := uint8(128)
		if options != nil && options.BilevelThreshold != 0 {
			threshold = options.BilevelThreshold
		}
		bytesPerRow := (int(width) + 7) / 8

		if useTiles {
			return f.writeTiles(ctx, bounds, options.TileWidth, options.TileHeight, 0, func(tileData []byte, tileX, tileY, tw, th int) {
				tileBytesPerRow := (tw + 7) / 8
				for row := 0; row < th; row++ {
					imgY := tileY + row
					for col := 0; col < tw; col++ {
						imgX := tileX + col
						var lum uint8
						if imgX < bounds.Max.X && imgY < bounds.Max.Y {
							r, g, b, _ := img.At(imgX, imgY).RGBA()
							lum = uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
						}
						byteIdx := row*tileBytesPerRow + col/8
						bitIdx := uint(7 - col%8)
						if lum < threshold {
							tileData[byteIdx] |= 1 << bitIdx // black = 1 for MINISWHITE
						}
					}
				}
			})
		}

		strip := uint32(0)
		rowsPerStrip := uint32(0)
		if options != nil && options.RowsPerStrip > 0 {
			rowsPerStrip = options.RowsPerStrip
		}
		if rowsPerStrip == 0 {
			var err error
			rowsPerStrip, err = f.TIFFDefaultStripSize(ctx, 0)
			if err != nil {
				return err
			}
		}
		// For CCITT we already set ROWSPERSTRIP above (in the !useTiles branch),
		// but we need the value to iterate.
		return f.writeStrips(ctx, bounds, rowsPerStrip, bytesPerRow, &strip, func(stripData []byte, y, rows int) {
			for row := 0; row < rows; row++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := img.At(x, y+row).RGBA()
					lum := uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)
					col := x - bounds.Min.X
					byteIdx := row*bytesPerRow + col/8
					bitIdx := uint(7 - col%8)
					if lum < threshold {
						stripData[byteIdx] |= 1 << bitIdx // black = 1 for MINISWHITE
					}
				}
			}
		})
	}

	// Write image data.
	bytesPerPixel := int(samplesPerPixel)
	bytesPerRow := int(width) * bytesPerPixel

	if useTiles {
		tileWidth := options.TileWidth
		tileHeight := options.TileHeight

		// Fast path for tiles.
		if rgbaImg, ok := img.(*image.RGBA); ok && !isJPEG && alphaMode == AlphaAssociated {
			return f.writeTiles(ctx, bounds, tileWidth, tileHeight, bytesPerPixel, func(tileData []byte, tileX, tileY, tw, th int) {
				tileBytesPerRow := tw * bytesPerPixel
				for row := 0; row < th; row++ {
					imgY := tileY + row
					if imgY >= bounds.Max.Y {
						break
					}
					srcY := imgY - bounds.Min.Y
					dstStart := row * tileBytesPerRow
					cols := tw
					if tileX+cols > bounds.Max.X {
						cols = bounds.Max.X - tileX
					}
					srcStart := srcY*rgbaImg.Stride + (tileX-bounds.Min.X)*4
					copy(tileData[dstStart:dstStart+cols*4], rgbaImg.Pix[srcStart:srcStart+cols*4])
				}
			})
		}
		if nrgbaImg, ok := img.(*image.NRGBA); ok && !isJPEG && alphaMode == AlphaUnassociated {
			return f.writeTiles(ctx, bounds, tileWidth, tileHeight, bytesPerPixel, func(tileData []byte, tileX, tileY, tw, th int) {
				tileBytesPerRow := tw * bytesPerPixel
				for row := 0; row < th; row++ {
					imgY := tileY + row
					if imgY >= bounds.Max.Y {
						break
					}
					srcY := imgY - bounds.Min.Y
					dstStart := row * tileBytesPerRow
					cols := tw
					if tileX+cols > bounds.Max.X {
						cols = bounds.Max.X - tileX
					}
					srcStart := srcY*nrgbaImg.Stride + (tileX-bounds.Min.X)*4
					copy(tileData[dstStart:dstStart+cols*4], nrgbaImg.Pix[srcStart:srcStart+cols*4])
				}
			})
		}

		// Generic tile path.
		if alphaMode == AlphaUnassociated {
			return f.writeTiles(ctx, bounds, tileWidth, tileHeight, bytesPerPixel, func(tileData []byte, tileX, tileY, tw, th int) {
				tileBytesPerRow := tw * bytesPerPixel
				for row := 0; row < th; row++ {
					imgY := tileY + row
					if imgY >= bounds.Max.Y {
						break
					}
					for col := 0; col < tw; col++ {
						imgX := tileX + col
						if imgX >= bounds.Max.X {
							break
						}
						c := color.NRGBAModel.Convert(img.At(imgX, imgY)).(color.NRGBA)
						offset := row*tileBytesPerRow + col*bytesPerPixel
						tileData[offset] = c.R
						tileData[offset+1] = c.G
						tileData[offset+2] = c.B
						if !isJPEG {
							tileData[offset+3] = c.A
						}
					}
				}
			})
		}

		return f.writeTiles(ctx, bounds, tileWidth, tileHeight, bytesPerPixel, func(tileData []byte, tileX, tileY, tw, th int) {
			tileBytesPerRow := tw * bytesPerPixel
			for row := 0; row < th; row++ {
				imgY := tileY + row
				if imgY >= bounds.Max.Y {
					break
				}
				for col := 0; col < tw; col++ {
					imgX := tileX + col
					if imgX >= bounds.Max.X {
						break
					}
					r, g, b, a := img.At(imgX, imgY).RGBA()
					offset := row*tileBytesPerRow + col*bytesPerPixel
					tileData[offset] = uint8(r >> 8)
					tileData[offset+1] = uint8(g >> 8)
					tileData[offset+2] = uint8(b >> 8)
					if !isJPEG {
						tileData[offset+3] = uint8(a >> 8)
					}
				}
			}
		})
	}

	// Strip-based output.
	strip := uint32(0)
	rowsPerStrip := uint32(0)
	if isJPEG {
		rowsPerStrip = height
	} else if options != nil && options.RowsPerStrip > 0 {
		rowsPerStrip = options.RowsPerStrip
	} else {
		var err error
		rowsPerStrip, err = f.TIFFDefaultStripSize(ctx, 0)
		if err != nil {
			return err
		}
	}

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

// writeTiles writes image data tile by tile, calling fillTile to populate each tile's pixel data,
// then writes the TIFF directory. bytesPerPixel is 0 for bilevel (bit-packed) data.
func (f *File) writeTiles(ctx context.Context, bounds image.Rectangle,
	tileWidth, tileHeight uint32, bytesPerPixel int,
	fillTile func(tileData []byte, tileX, tileY, tw, th int)) error {

	tw := int(tileWidth)
	th := int(tileHeight)
	var tileSize int
	if bytesPerPixel == 0 {
		// Bilevel: bit-packed.
		tileSize = ((tw + 7) / 8) * th
	} else {
		tileSize = tw * th * bytesPerPixel
	}

	tile := uint32(0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += th {
		for x := bounds.Min.X; x < bounds.Max.X; x += tw {
			tileData := make([]byte, tileSize)
			fillTile(tileData, x, y, tw, th)

			if err := f.TIFFWriteEncodedTile(ctx, tile, tileData); err != nil {
				return err
			}
			tile++
		}
	}

	return f.TIFFWriteDirectory(ctx)
}
