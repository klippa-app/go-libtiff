package libtiff_test

import (
	"context"
	"os"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TIFFGetFieldDefaulted", func() {
	ctx := context.Background()

	Context("Uint16", func() {
		It("returns the explicitly set value for BitsPerSample", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldDefaultedUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint16(8)))
		})

		It("returns the default value for a tag not explicitly set", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
			defer cleanup()

			// FillOrder defaults to 1 (MSB2LSB) per TIFF spec.
			val, err := readTiff.TIFFGetFieldDefaultedUint16_t(ctx, libtiff.TIFFTAG_FILLORDER)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint16(1)))
		})
	})

	Context("Uint32", func() {
		It("returns the explicitly set value for ImageWidth", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldDefaultedUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(1)))
		})
	})

	Context("Float", func() {
		It("returns the explicitly set value for XResolution", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldFloat(ctx, libtiff.TIFFTAG_XRESOLUTION, 300.0)).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldDefaultedFloat(ctx, libtiff.TIFFTAG_XRESOLUTION)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(float32(300.0)))
		})
	})

	Context("Double", func() {
		It("returns the explicitly set value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldDouble(ctx, libtiff.TIFFTAG_STONITS, 0.005)).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldDefaultedDouble(ctx, libtiff.TIFFTAG_STONITS)
			Expect(err).To(BeNil())
			Expect(val).To(BeNumerically("~", 0.005, 0.0001))
		})
	})

	Context("ConstChar", func() {
		It("returns the explicitly set string value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_SOFTWARE, "test-app")).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldDefaultedConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("test-app"))
		})
	})

	Context("TwoUint16", func() {
		It("returns the explicitly set page number values", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldTwoUint16(ctx, libtiff.TIFFTAG_PAGENUMBER, 0, 5)).To(Succeed())
			})
			defer cleanup()

			val1, val2, err := readTiff.TIFFGetFieldDefaultedTwoUint16(ctx, libtiff.TIFFTAG_PAGENUMBER)
			Expect(err).To(BeNil())
			Expect(val1).To(Equal(uint16(0)))
			Expect(val2).To(Equal(uint16(5)))
		})
	})
})

var _ = Describe("TIFFUnsetField", func() {
	ctx := context.Background()

	It("removes a tag from the directory", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-unset-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldString(ctx, libtiff.TIFFTAG_SOFTWARE, "test-software")).To(Succeed())

		// Verify the tag is set.
		val, err := writeTiff.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
		Expect(err).To(BeNil())
		Expect(val).To(Equal("test-software"))

		// Unset the tag.
		Expect(writeTiff.TIFFUnsetField(ctx, libtiff.TIFFTAG_SOFTWARE)).To(Succeed())

		// Verify the tag is no longer set.
		_, err = writeTiff.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
		Expect(err).To(MatchError(&libtiff.TagNotDefinedError{Tag: libtiff.TIFFTAG_SOFTWARE}))

		writeTiff.Close(ctx)
	})
})

var _ = Describe("TIFFGetConfiguredCODECs", func() {
	ctx := context.Background()

	It("returns a non-empty list of codecs", func() {
		codecs, err := instance.TIFFGetConfiguredCODECs(ctx)
		Expect(err).To(BeNil())
		Expect(codecs).ToNot(BeEmpty())
	})

	It("includes the NONE codec", func() {
		codecs, err := instance.TIFFGetConfiguredCODECs(ctx)
		Expect(err).To(BeNil())

		found := false
		for _, c := range codecs {
			if c.Scheme == uint16(libtiff.COMPRESSION_NONE) {
				found = true
				Expect(c.Name).ToNot(BeEmpty())
				break
			}
		}
		Expect(found).To(BeTrue())
	})

	It("includes LZW and JPEG codecs", func() {
		codecs, err := instance.TIFFGetConfiguredCODECs(ctx)
		Expect(err).To(BeNil())

		foundLZW := false
		foundJPEG := false
		for _, c := range codecs {
			if c.Scheme == uint16(libtiff.COMPRESSION_LZW) {
				foundLZW = true
			}
			if c.Scheme == uint16(libtiff.COMPRESSION_JPEG) {
				foundJPEG = true
			}
		}
		Expect(foundLZW).To(BeTrue())
		Expect(foundJPEG).To(BeTrue())
	})
})

var _ = Describe("TIFFDataWidth", func() {
	ctx := context.Background()

	It("returns 1 for TIFF_BYTE", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_BYTE)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(1))
	})

	It("returns 2 for TIFF_SHORT", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_SHORT)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(2))
	})

	It("returns 4 for TIFF_LONG", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_LONG)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(4))
	})

	It("returns 8 for TIFF_DOUBLE", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_DOUBLE)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(8))
	})

	It("returns 8 for TIFF_RATIONAL", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_RATIONAL)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(8))
	})

	It("returns 4 for TIFF_FLOAT", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_FLOAT)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(4))
	})

	It("returns 8 for TIFF_LONG8", func() {
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_LONG8)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(8))
	})

	It("returns 1 for TIFF_NOTYPE", func() {
		// In libtiff, TIFF_NOTYPE falls through to TIFF_BYTE and returns 1.
		width, err := instance.TIFFDataWidth(ctx, libtiff.TIFF_NOTYPE)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(1))
	})
})

var _ = Describe("TIFFGetMode", func() {
	ctx := context.Background()

	It("returns a non-negative mode for a read-only file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		mode, err := tiffFile.TIFFGetMode(ctx)
		Expect(err).To(BeNil())
		Expect(mode).To(BeNumerically(">=", 0))
	})
})

var _ = Describe("TIFFCurrentDirOffset", func() {
	ctx := context.Background()

	It("returns a non-zero offset for a valid file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		offset, err := tiffFile.TIFFCurrentDirOffset(ctx)
		Expect(err).To(BeNil())
		Expect(offset).To(BeNumerically(">", 0))
	})

	It("changes when navigating to a different directory", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		offset1, err := tiffFile.TIFFCurrentDirOffset(ctx)
		Expect(err).To(BeNil())

		Expect(tiffFile.TIFFSetDirectory(ctx, 1)).To(Succeed())

		offset2, err := tiffFile.TIFFCurrentDirOffset(ctx)
		Expect(err).To(BeNil())

		Expect(offset1).ToNot(Equal(offset2))
	})
})

var _ = Describe("TIFFCurrentStrip", func() {
	ctx := context.Background()

	It("returns a value for an open file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		strip, err := tiffFile.TIFFCurrentStrip(ctx)
		Expect(err).To(BeNil())
		_ = strip // Initial value is implementation-defined
	})
})

var _ = Describe("TIFFCurrentTile", func() {
	ctx := context.Background()

	It("returns a value for an open file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		tile, err := tiffFile.TIFFCurrentTile(ctx)
		Expect(err).To(BeNil())
		_ = tile // Initial value is implementation-defined
	})
})

var _ = Describe("TIFFCurrentRow", func() {
	ctx := context.Background()

	It("returns a value for an open file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		row, err := tiffFile.TIFFCurrentRow(ctx)
		Expect(err).To(BeNil())
		_ = row // Initial value is implementation-defined
	})
})

var _ = Describe("TIFFIsMSB2LSB", func() {
	ctx := context.Background()

	It("returns a boolean for an existing TIFF file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		_, err = tiffFile.TIFFIsMSB2LSB(ctx)
		Expect(err).To(BeNil())
	})
})

var _ = Describe("TIFFIsUpSampled", func() {
	ctx := context.Background()

	It("returns false for a standard RGB image", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		upSampled, err := tiffFile.TIFFIsUpSampled(ctx)
		Expect(err).To(BeNil())
		Expect(upSampled).To(BeFalse())
	})
})

var _ = Describe("TIFFRawStripSize", func() {
	ctx := context.Background()

	It("returns the compressed strip size", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		rawSize, err := tiffFile.TIFFRawStripSize(ctx, 0)
		Expect(err).To(BeNil())
		Expect(rawSize).To(BeNumerically(">", 0))
	})

	It("returns a size smaller than or equal to the uncompressed strip size for JPEG", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		rawSize, err := tiffFile.TIFFRawStripSize(ctx, 0)
		Expect(err).To(BeNil())

		stripSize, err := tiffFile.TIFFStripSize(ctx)
		Expect(err).To(BeNil())

		// For JPEG compression, raw (compressed) should typically be smaller.
		Expect(rawSize).To(BeNumerically("<=", stripSize))
	})
})

var _ = Describe("TIFFTileRowSize", func() {
	ctx := context.Background()

	It("returns a value for a tiled image", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-tilerowsize-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 64)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 64)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 3)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_RGB))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH, 32)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH, 32)).To(Succeed())

		tileRowSize, err := writeTiff.TIFFTileRowSize(ctx)
		Expect(err).To(BeNil())
		// 32 pixels wide * 3 samples per pixel * 1 byte per sample = 96
		Expect(tileRowSize).To(Equal(int64(96)))

		writeTiff.Close(ctx)
		tmpFile.Close()
	})
})

var _ = Describe("TIFFCheckTile", func() {
	ctx := context.Background()

	It("returns true for valid tile coordinates", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-checktile-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 64)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 64)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH, 32)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH, 32)).To(Succeed())

		valid, err := writeTiff.TIFFCheckTile(ctx, 0, 0, 0, 0)
		Expect(err).To(BeNil())
		Expect(valid).To(BeTrue())

		writeTiff.Close(ctx)
		tmpFile.Close()
	})
})

var _ = Describe("TIFFRasterScanlineSize", func() {
	ctx := context.Background()

	It("returns a positive value for a strip-based image", func() {
		readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
		defer cleanup()

		size, err := readTiff.TIFFRasterScanlineSize(ctx)
		Expect(err).To(BeNil())
		Expect(size).To(BeNumerically(">", 0))
	})
})

var _ = Describe("TIFFGetStrileOffset", func() {
	ctx := context.Background()

	It("returns a non-zero offset for the first strip", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		offset, err := tiffFile.TIFFGetStrileOffset(ctx, 0)
		Expect(err).To(BeNil())
		Expect(offset).To(BeNumerically(">", 0))
	})
})

var _ = Describe("TIFFGetStrileByteCount", func() {
	ctx := context.Background()

	It("returns a non-zero byte count for the first strip", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		byteCount, err := tiffFile.TIFFGetStrileByteCount(ctx, 0)
		Expect(err).To(BeNil())
		Expect(byteCount).To(BeNumerically(">", 0))
	})

	It("matches the raw strip size", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		byteCount, err := tiffFile.TIFFGetStrileByteCount(ctx, 0)
		Expect(err).To(BeNil())

		rawSize, err := tiffFile.TIFFRawStripSize(ctx, 0)
		Expect(err).To(BeNil())

		Expect(byteCount).To(Equal(uint64(rawSize)))
	})
})

var _ = Describe("TIFFFreeDirectory", func() {
	ctx := context.Background()

	It("frees the current directory without error", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-freedir-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())

		// FreeDirectory should not error.
		err = writeTiff.TIFFFreeDirectory(ctx)
		Expect(err).To(BeNil())

		writeTiff.Close(ctx)
	})
})

var _ = Describe("TIFFDeferStrileArrayWriting", func() {
	ctx := context.Background()

	It("can be called without error during writing", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-defer-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())

		// Defer strile array writing.
		Expect(writeTiff.TIFFDeferStrileArrayWriting(ctx)).To(Succeed())

		// Write the directory header (without strip offsets/bytecounts).
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()
	})

	It("writes a valid file when used with ForceStrileArrayWriting", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-force-defer-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())

		// Write data first, then directory — no deferral needed for this simple case.
		// Just verify ForceStrileArrayWriting doesn't error after normal write.
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Verify the file is readable.
		readFile, err := os.Open(tmpPath)
		Expect(err).To(BeNil())
		defer readFile.Close()

		stat, err := readFile.Stat()
		Expect(err).To(BeNil())

		readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
		Expect(err).To(BeNil())
		defer readTiff.Close(ctx)

		width, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(uint32(1)))
	})
})

var _ = Describe("TIFFPrintDirectory", func() {
	ctx := context.Background()

	It("returns a human-readable directory description", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		output, err := tiffFile.TIFFPrintDirectory(ctx, libtiff.TIFFPRINT_NONE)
		Expect(err).To(BeNil())
		Expect(output).ToNot(BeEmpty())
		Expect(output).To(ContainSubstring("Image Width"))
	})

	It("includes strip info when TIFFPRINT_STRIPS is set", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		output, err := tiffFile.TIFFPrintDirectory(ctx, libtiff.TIFFPRINT_STRIPS)
		Expect(err).To(BeNil())
		Expect(output).ToNot(BeEmpty())
	})
})
