package libtiff_test

import (
	"context"
	"os"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// writeMinimalTiff creates a minimal 1x1 grayscale TIFF, calls setup to set
// additional tags, writes a single strip and directory, then reopens for reading.
func writeMinimalTiff(ctx context.Context, setup func(ctx context.Context, f *libtiff.File)) (*libtiff.File, func()) {
	tmpFile, err := os.CreateTemp("", "libtiff-write-test-*.tif")
	Expect(err).To(BeNil())
	tmpPath := tmpFile.Name()

	fileMode := "w"
	writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
		FileMode: &fileMode,
	})
	Expect(err).To(BeNil())

	// Set minimum required tags for a valid 1x1 grayscale TIFF.
	Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
	Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
	Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
	Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
	Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
	Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())

	setup(ctx, writeTiff)

	Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
	Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

	writeTiff.Close(ctx)
	tmpFile.Close()

	// Reopen for reading.
	readFile, err := os.Open(tmpPath)
	Expect(err).To(BeNil())

	stat, err := readFile.Stat()
	Expect(err).To(BeNil())

	readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
	Expect(err).To(BeNil())

	cleanup := func() {
		readTiff.Close(ctx)
		readFile.Close()
		os.Remove(tmpPath)
	}

	return readTiff, cleanup
}

var _ = Describe("TIFFOpenFileFromReadWriteSeeker", func() {
	ctx := context.Background()

	It("opens a new file for writing", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-open-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		tiffFile, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		Expect(tiffFile).ToNot(BeNil())
		tiffFile.Close(ctx)
	})

	It("opens an existing file for appending", func() {
		// First create a valid TIFF.
		tmpFile, err := os.CreateTemp("", "libtiff-append-test-*.tif")
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
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Reopen for appending.
		appendFile, err := os.OpenFile(tmpPath, os.O_RDWR, 0)
		Expect(err).To(BeNil())
		defer appendFile.Close()

		appendMode := "a"
		appendTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", appendFile, 0, &libtiff.OpenOptions{
			FileMode: &appendMode,
		})
		Expect(err).To(BeNil())
		Expect(appendTiff).ToNot(BeNil())

		// Write a second directory.
		Expect(appendTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 2)).To(Succeed())
		Expect(appendTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 2)).To(Succeed())
		Expect(appendTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(appendTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(appendTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(appendTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 2)).To(Succeed())
		Expect(appendTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{64, 64})).To(Succeed())
		Expect(appendTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		appendTiff.Close(ctx)
		appendFile.Close()

		// Verify both directories exist.
		readFile, err := os.Open(tmpPath)
		Expect(err).To(BeNil())
		defer readFile.Close()

		stat, err := readFile.Stat()
		Expect(err).To(BeNil())

		readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
		Expect(err).To(BeNil())
		defer readTiff.Close(ctx)

		numDirs, err := readTiff.TIFFNumberOfDirectories(ctx)
		Expect(err).To(BeNil())
		Expect(numDirs).To(Equal(uint32(2)))

		// First directory: 1x1.
		width, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(uint32(1)))

		// Second directory: 2x2.
		Expect(readTiff.TIFFSetDirectory(ctx, 1)).To(Succeed())
		width, err = readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(uint32(2)))
	})
})

var _ = Describe("TIFFSetField", func() {
	ctx := context.Background()

	Context("TIFFSetFieldUint16_t", func() {
		It("writes and reads back a uint16 tag value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_ORIENTATION, uint16(libtiff.ORIENTATION_BOTRIGHT))).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_ORIENTATION)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint16(libtiff.ORIENTATION_BOTRIGHT)))
		})
	})

	Context("TIFFSetFieldUint32_t", func() {
		It("writes and reads back a uint32 tag value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				// IMAGEWIDTH is already set to 1 by the helper; override to verify.
				Expect(f.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 42)).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(42)))
		})
	})

	Context("TIFFSetFieldInt", func() {
		It("writes an int tag value without error", func() {
			// JPEGQUALITY is a pseudo-tag that doesn't persist in the file,
			// but setting it should not error on a JPEG-compressed TIFF.
			tmpFile, err := os.CreateTemp("", "libtiff-int-test-*.tif")
			Expect(err).To(BeNil())
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			fileMode := "w"
			writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
				FileMode: &fileMode,
			})
			Expect(err).To(BeNil())

			Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 8)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 8)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 3)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION, uint16(libtiff.COMPRESSION_JPEG))).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_YCBCR))).To(Succeed())
			Expect(writeTiff.TIFFSetFieldInt(ctx, libtiff.TIFFTAG_JPEGCOLORMODE, int(libtiff.JPEGCOLORMODE_RGB))).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 8)).To(Succeed())

			Expect(writeTiff.TIFFSetFieldInt(ctx, libtiff.TIFFTAG_JPEGQUALITY, 50)).To(Succeed())

			// Write a minimal 8x8 RGB strip (JPEG requires MCU-aligned dimensions).
			strip := make([]byte, 8*8*3)
			Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, strip)).To(Succeed())
			Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())
			writeTiff.Close(ctx)
		})
	})

	Context("TIFFSetFieldFloat", func() {
		It("writes and reads back a float tag value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldFloat(ctx, libtiff.TIFFTAG_XRESOLUTION, 300.0)).To(Succeed())
				Expect(f.TIFFSetFieldFloat(ctx, libtiff.TIFFTAG_YRESOLUTION, 150.0)).To(Succeed())
			})
			defer cleanup()

			xres, err := readTiff.TIFFGetFieldFloat(ctx, libtiff.TIFFTAG_XRESOLUTION)
			Expect(err).To(BeNil())
			Expect(xres).To(Equal(float32(300.0)))

			yres, err := readTiff.TIFFGetFieldFloat(ctx, libtiff.TIFFTAG_YRESOLUTION)
			Expect(err).To(BeNil())
			Expect(yres).To(Equal(float32(150.0)))
		})
	})

	Context("TIFFSetFieldDouble", func() {
		It("writes and reads back a double tag value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldDouble(ctx, libtiff.TIFFTAG_STONITS, 0.005)).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldDouble(ctx, libtiff.TIFFTAG_STONITS)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(0.005))
		})
	})

	Context("TIFFSetFieldString", func() {
		It("writes and reads back a string tag value", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_SOFTWARE, "test-software-1.0")).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("test-software-1.0"))
		})

		It("writes and reads back the DATETIME tag", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_DATETIME, "2024:01:15 10:30:00")).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_DATETIME)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("2024:01:15 10:30:00"))
		})

		It("writes and reads back the ARTIST tag", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_ARTIST, "Jane Doe")).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_ARTIST)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("Jane Doe"))
		})

		It("handles an empty string", func() {
			readTiff, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldString(ctx, libtiff.TIFFTAG_SOFTWARE, "")).To(Succeed())
			})
			defer cleanup()

			val, err := readTiff.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(""))
		})
	})

	Context("TIFFSetFieldExtraSamples", func() {
		It("sets associated alpha without error", func() {
			// Extra samples requires spp > photometric channels, so use RGBA.
			_, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 4)).To(Succeed())
				Expect(f.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_RGB))).To(Succeed())
				Expect(f.TIFFSetFieldExtraSamples(ctx, []uint16{uint16(libtiff.EXTRASAMPLE_ASSOCALPHA)})).To(Succeed())
			})
			defer cleanup()
		})

		It("sets unassociated alpha without error", func() {
			_, cleanup := writeMinimalTiff(ctx, func(ctx context.Context, f *libtiff.File) {
				Expect(f.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 4)).To(Succeed())
				Expect(f.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_RGB))).To(Succeed())
				Expect(f.TIFFSetFieldExtraSamples(ctx, []uint16{uint16(libtiff.EXTRASAMPLE_UNASSALPHA)})).To(Succeed())
			})
			defer cleanup()
		})
	})
})

var _ = Describe("TIFFWriteEncodedStrip", func() {
	ctx := context.Background()

	It("writes strip data that survives a roundtrip", func() {
		readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {
			// No extra tags needed; the helper writes pixel value 128.
		})
		defer cleanup()

		goImage, imgCleanup, err := readTiff.ToGoImage(ctx)
		Expect(err).To(BeNil())
		defer imgCleanup(ctx)

		// The 1x1 grayscale pixel should be readable.
		Expect(goImage.Bounds().Dx()).To(Equal(1))
		Expect(goImage.Bounds().Dy()).To(Equal(1))
	})

	It("writes multiple strips correctly", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-strips-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 4)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 4)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		// 2 rows per strip → 2 strips for a 4-row image.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 2)).To(Succeed())

		// Strip 0: rows 0-1 (4 bytes each row, 8 bytes total).
		strip0 := []byte{10, 20, 30, 40, 50, 60, 70, 80}
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, strip0)).To(Succeed())

		// Strip 1: rows 2-3.
		strip1 := []byte{90, 100, 110, 120, 130, 140, 150, 160}
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 1, strip1)).To(Succeed())

		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())
		writeTiff.Close(ctx)
		tmpFile.Close()

		// Reopen and verify dimensions.
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
		Expect(width).To(Equal(uint32(4)))

		height, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH)
		Expect(err).To(BeNil())
		Expect(height).To(Equal(uint32(4)))
	})
})

var _ = Describe("TIFFWriteDirectory", func() {
	ctx := context.Background()

	It("writes multiple directories to a single file", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-multidir-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write 3 directories with different widths.
		for _, w := range []uint32{10, 20, 30} {
			Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, w)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
			Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())

			strip := make([]byte, w)
			Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, strip)).To(Succeed())
			Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())
		}

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Reopen and verify.
		readFile, err := os.Open(tmpPath)
		Expect(err).To(BeNil())
		defer readFile.Close()

		stat, err := readFile.Stat()
		Expect(err).To(BeNil())

		readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
		Expect(err).To(BeNil())
		defer readTiff.Close(ctx)

		numDirs, err := readTiff.TIFFNumberOfDirectories(ctx)
		Expect(err).To(BeNil())
		Expect(numDirs).To(Equal(uint32(3)))

		for i, expectedWidth := range []uint32{10, 20, 30} {
			Expect(readTiff.TIFFSetDirectory(ctx, uint32(i))).To(Succeed())
			width, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
			Expect(err).To(BeNil())
			Expect(width).To(Equal(expectedWidth))
		}
	})
})

var _ = Describe("TIFFDefaultStripSize", func() {
	ctx := context.Background()

	It("returns a non-zero strip size", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-stripsize-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		defer writeTiff.Close(ctx)

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 100)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 200)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 3)).To(Succeed())

		stripSize, err := writeTiff.TIFFDefaultStripSize(ctx, 0)
		Expect(err).To(BeNil())
		Expect(stripSize).To(BeNumerically(">", 0))
		Expect(stripSize).To(BeNumerically("<=", 200))
	})
})

var _ = Describe("TIFFWriteScanline", func() {
	ctx := context.Background()

	It("writes an image scanline by scanline", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-scanline-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 4)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 3)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 3)).To(Succeed())

		for row := uint32(0); row < 3; row++ {
			scanline := []byte{byte(row * 10), byte(row*10 + 1), byte(row*10 + 2), byte(row*10 + 3)}
			Expect(writeTiff.TIFFWriteScanline(ctx, scanline, row, 0)).To(Succeed())
		}

		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())
		writeTiff.Close(ctx)
		tmpFile.Close()

		// Verify.
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
		Expect(width).To(Equal(uint32(4)))

		height, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH)
		Expect(err).To(BeNil())
		Expect(height).To(Equal(uint32(3)))
	})
})

var _ = Describe("TIFFWriteRawStrip", func() {
	ctx := context.Background()

	It("writes raw uncompressed strip data", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-rawstrip-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 4)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 2)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 2)).To(Succeed())

		rawData := []byte{10, 20, 30, 40, 50, 60, 70, 80}
		Expect(writeTiff.TIFFWriteRawStrip(ctx, 0, rawData)).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())
		writeTiff.Close(ctx)
		tmpFile.Close()

		// Verify it's a valid TIFF.
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
		Expect(width).To(Equal(uint32(4)))
	})
})

var _ = Describe("TIFFWriteRawTile", func() {
	ctx := context.Background()

	It("writes raw tile data to a tiled TIFF", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-rawtile-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 256)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 256)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH, 256)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH, 256)).To(Succeed())

		tileData := make([]byte, 256*256)
		Expect(writeTiff.TIFFWriteRawTile(ctx, 0, tileData)).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())
		writeTiff.Close(ctx)
		tmpFile.Close()

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
		Expect(width).To(Equal(uint32(256)))
	})
})

var _ = Describe("TIFFFlush", func() {
	ctx := context.Background()

	It("flushes pending writes without error", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-flush-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		defer writeTiff.Close(ctx)

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())

		Expect(writeTiff.TIFFFlush(ctx)).To(Succeed())
	})
})

var _ = Describe("TIFFFlushData", func() {
	ctx := context.Background()

	It("flushes pending data without error", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-flushdata-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		defer writeTiff.Close(ctx)

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())

		Expect(writeTiff.TIFFFlushData(ctx)).To(Succeed())
	})
})

var _ = Describe("TIFFCheckpointDirectory", func() {
	ctx := context.Background()

	It("checkpoints the directory without closing it", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-checkpoint-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		defer writeTiff.Close(ctx)

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())

		Expect(writeTiff.TIFFCheckpointDirectory(ctx)).To(Succeed())
	})
})

var _ = Describe("TIFFRewriteDirectory", func() {
	ctx := context.Background()

	It("rewrites the current directory in place", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-rewrite-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		defer writeTiff.Close(ctx)

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())

		// First write directory, then rewrite it.
		Expect(writeTiff.TIFFCheckpointDirectory(ctx)).To(Succeed())
		Expect(writeTiff.TIFFRewriteDirectory(ctx)).To(Succeed())
	})
})
