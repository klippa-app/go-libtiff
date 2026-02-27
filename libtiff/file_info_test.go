package libtiff_test

import (
	"context"
	"os"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TIFFIsBigEndian", func() {
	ctx := context.Background()

	It("returns a boolean for an existing TIFF file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		_, err = tiffFile.TIFFIsBigEndian(ctx)
		Expect(err).To(BeNil())
	})

	It("returns a consistent value for a written TIFF", func() {
		readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
		defer cleanup()

		bigEndian, err := readTiff.TIFFIsBigEndian(ctx)
		Expect(err).To(BeNil())
		// WASM is little-endian, so written files should be little-endian.
		Expect(bigEndian).To(BeFalse())
	})
})

var _ = Describe("TIFFIsBigTIFF", func() {
	ctx := context.Background()

	It("returns false for a regular TIFF file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		isBigTIFF, err := tiffFile.TIFFIsBigTIFF(ctx)
		Expect(err).To(BeNil())
		Expect(isBigTIFF).To(BeFalse())
	})
})

var _ = Describe("TIFFIsByteSwapped", func() {
	ctx := context.Background()

	It("returns a boolean for an existing TIFF file", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		_, err = tiffFile.TIFFIsByteSwapped(ctx)
		Expect(err).To(BeNil())
	})
})

var _ = Describe("TIFFIsCODECConfigured", func() {
	ctx := context.Background()

	It("returns true for COMPRESSION_NONE", func() {
		configured, err := instance.TIFFIsCODECConfigured(ctx, uint16(libtiff.COMPRESSION_NONE))
		Expect(err).To(BeNil())
		Expect(configured).To(BeTrue())
	})

	It("returns true for COMPRESSION_LZW", func() {
		configured, err := instance.TIFFIsCODECConfigured(ctx, uint16(libtiff.COMPRESSION_LZW))
		Expect(err).To(BeNil())
		Expect(configured).To(BeTrue())
	})

	It("returns true for COMPRESSION_JPEG", func() {
		configured, err := instance.TIFFIsCODECConfigured(ctx, uint16(libtiff.COMPRESSION_JPEG))
		Expect(err).To(BeNil())
		Expect(configured).To(BeTrue())
	})

	It("returns false for an invalid codec", func() {
		configured, err := instance.TIFFIsCODECConfigured(ctx, 9999)
		Expect(err).To(BeNil())
		Expect(configured).To(BeFalse())
	})
})

var _ = Describe("TIFFRGBAImageOK", func() {
	ctx := context.Background()

	It("returns true for a standard RGB TIFF", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		ok, msg, err := tiffFile.TIFFRGBAImageOK(ctx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
		Expect(msg).To(BeEmpty())
	})

	It("returns true for a written RGBA TIFF", func() {
		img := createTestRGBA(8, 8)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		ok, msg, err := tiffFile.TIFFRGBAImageOK(ctx)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
		Expect(msg).To(BeEmpty())
	})
})

var _ = Describe("TIFFFileName", func() {
	ctx := context.Background()

	It("returns the file name used when opening", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		name, err := tiffFile.TIFFFileName(ctx)
		Expect(err).To(BeNil())
		Expect(name).To(Equal("/testdata/multipage-sample.tif"))
	})

	It("returns the custom name used with TIFFOpenFileFromReader", func() {
		readTiff, cleanup := writeMinimalTiff(ctx, func(_ context.Context, _ *libtiff.File) {})
		defer cleanup()

		name, err := readTiff.TIFFFileName(ctx)
		Expect(err).To(BeNil())
		Expect(name).To(Equal("test.tif"))
	})
})

var _ = Describe("TIFFCreateDirectory", func() {
	ctx := context.Background()

	It("creates a new directory in a writable TIFF", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-createdir-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write a first directory.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		// Create a fresh directory.
		Expect(writeTiff.TIFFCreateDirectory(ctx)).To(Succeed())

		// Write tags in the new directory.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 2)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 2)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 2)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{64, 64, 64, 64})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

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
	})
})

var _ = Describe("TIFFUnlinkDirectory", func() {
	ctx := context.Background()

	It("removes a directory from a multi-directory TIFF", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-unlink-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write 3 directories.
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

		// Reopen for read-write and unlink directory 2 (1-based, so the second directory with width=20).
		rwFile, err := os.OpenFile(tmpPath, os.O_RDWR, 0)
		Expect(err).To(BeNil())
		defer rwFile.Close()

		rwMode := "r+"
		rwTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", rwFile, 0, &libtiff.OpenOptions{
			FileMode: &rwMode,
		})
		Expect(err).To(BeNil())

		Expect(rwTiff.TIFFUnlinkDirectory(ctx, 2)).To(Succeed())
		rwTiff.Close(ctx)
		rwFile.Close()

		// Verify only 2 directories remain.
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

		// First directory should still be width=10.
		width, err := readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(uint32(10)))

		// Second directory should now be width=30.
		Expect(readTiff.TIFFSetDirectory(ctx, 1)).To(Succeed())
		width, err = readTiff.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())
		Expect(width).To(Equal(uint32(30)))
	})
})

var _ = Describe("TIFFCreateEXIFDirectory", func() {
	ctx := context.Background()

	It("creates an EXIF sub-IFD and writes it to the file", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-exif-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write the main IFD first.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		// Create the EXIF sub-IFD after the main IFD is written.
		Expect(writeTiff.TIFFCreateEXIFDirectory(ctx)).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Verify the file is valid.
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

var _ = Describe("TIFFCreateGPSDirectory", func() {
	ctx := context.Background()

	It("creates a GPS sub-IFD without error", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-gps-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write the main IFD first.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		// Create the GPS sub-IFD after the main IFD is written.
		Expect(writeTiff.TIFFCreateGPSDirectory(ctx)).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Verify the file is valid.
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

var _ = Describe("WarnHandler", func() {
	ctx := context.Background()

	It("receives warning messages from libtiff", func() {
		var warnings []string
		warnAboutUnknown := true
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", &libtiff.OpenOptions{
			WarnHandler: func(module string, message string) {
				warnings = append(warnings, message)
			},
			WarnAboutUnknownTags: &warnAboutUnknown,
		})
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		// The file should open successfully regardless of warnings.
		width, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())
		Expect(width).To(BeNumerically(">", 0))
	})
})

var _ = Describe("TIFFWriteCustomDirectory EXIF roundtrip", func() {
	ctx := context.Background()

	It("writes an EXIF sub-IFD and reads it back via TIFFReadEXIFDirectory", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-exif-roundtrip-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write the main IFD first.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		// Create the EXIF sub-IFD and set at least one tag
		// (an empty directory causes malloc(0) → NULL in WASM).
		Expect(writeTiff.TIFFCreateEXIFDirectory(ctx)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.EXIFTAG_COLORSPACE, 1)).To(Succeed())

		// Write the EXIF sub-IFD and get its offset.
		exifOffset, err := writeTiff.TIFFWriteCustomDirectory(ctx)
		Expect(err).To(BeNil())
		Expect(exifOffset).To(BeNumerically(">", 0))

		// Go back to the main IFD and link the EXIF sub-IFD.
		Expect(writeTiff.TIFFSetDirectory(ctx, 0)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint64_t(ctx, libtiff.TIFFTAG_EXIFIFD, exifOffset)).To(Succeed())
		Expect(writeTiff.TIFFRewriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Reopen and verify we can read the EXIF directory.
		readFile, err := os.Open(tmpPath)
		Expect(err).To(BeNil())
		defer readFile.Close()

		stat, err := readFile.Stat()
		Expect(err).To(BeNil())

		readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
		Expect(err).To(BeNil())
		defer readTiff.Close(ctx)

		// Read the EXIF offset tag.
		readOffset, err := readTiff.TIFFGetFieldUint64_t(ctx, libtiff.TIFFTAG_EXIFIFD)
		Expect(err).To(BeNil())
		Expect(readOffset).To(Equal(exifOffset))

		// Read the EXIF directory.
		Expect(readTiff.TIFFReadEXIFDirectory(ctx, readOffset)).To(Succeed())
	})
})

var _ = Describe("TIFFWriteCustomDirectory GPS roundtrip", func() {
	ctx := context.Background()

	It("writes a GPS sub-IFD and reads it back via TIFFReadGPSDirectory", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-gps-roundtrip-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		// Write the main IFD first.
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP, 1)).To(Succeed())
		Expect(writeTiff.TIFFWriteEncodedStrip(ctx, 0, []byte{128})).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		// Create the GPS sub-IFD and set at least one tag
		// (an empty directory causes malloc(0) → NULL in WASM).
		Expect(writeTiff.TIFFCreateGPSDirectory(ctx)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldString(ctx, libtiff.GPSTAG_LATITUDEREF, "N")).To(Succeed())

		// Write the GPS sub-IFD and get its offset.
		gpsOffset, err := writeTiff.TIFFWriteCustomDirectory(ctx)
		Expect(err).To(BeNil())
		Expect(gpsOffset).To(BeNumerically(">", 0))

		// Go back to the main IFD and link the GPS sub-IFD.
		Expect(writeTiff.TIFFSetDirectory(ctx, 0)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint64_t(ctx, libtiff.TIFFTAG_GPSIFD, gpsOffset)).To(Succeed())
		Expect(writeTiff.TIFFRewriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Reopen and verify we can read the GPS directory.
		readFile, err := os.Open(tmpPath)
		Expect(err).To(BeNil())
		defer readFile.Close()

		stat, err := readFile.Stat()
		Expect(err).To(BeNil())

		readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
		Expect(err).To(BeNil())
		defer readTiff.Close(ctx)

		// Read the GPS offset tag.
		readOffset, err := readTiff.TIFFGetFieldUint64_t(ctx, libtiff.TIFFTAG_GPSIFD)
		Expect(err).To(BeNil())
		Expect(readOffset).To(Equal(gpsOffset))

		// Read the GPS directory.
		Expect(readTiff.TIFFReadGPSDirectory(ctx, readOffset)).To(Succeed())
	})
})

var _ = Describe("OpenOptions via TIFFOpenFileFromReader", func() {
	ctx := context.Background()

	It("returns an error when MaxSingleMemAlloc is reached via Reader", func() {
		filePath := "../testdata/multipage-sample.tif"
		f, err := os.Open(filePath)
		Expect(err).To(BeNil())
		defer f.Close()

		stat, err := f.Stat()
		Expect(err).To(BeNil())

		maxMemory := int32(30)
		tiffFile, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", f, uint64(stat.Size()), &libtiff.OpenOptions{
			MaxSingleMemAlloc: &maxMemory,
		})
		Expect(err).To(MatchError(ContainSubstring("is beyond the 30 byte limit defined in open options")))
		Expect(tiffFile).To(BeNil())
	})

	It("returns an error when MaxCumulatedMemAlloc is reached via Reader", func() {
		filePath := "../testdata/multipage-sample.tif"
		f, err := os.Open(filePath)
		Expect(err).To(BeNil())
		defer f.Close()

		stat, err := f.Stat()
		Expect(err).To(BeNil())

		maxMemory := int32(30)
		tiffFile, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", f, uint64(stat.Size()), &libtiff.OpenOptions{
			MaxCumulatedMemAlloc: &maxMemory,
		})
		Expect(err).To(MatchError(ContainSubstring("is beyond the 30 cumulated byte limit defined in open options")))
		Expect(tiffFile).To(BeNil())
	})
})

var _ = Describe("TIFFSetSubDirectory", func() {
	ctx := context.Background()

	It("returns an error for an invalid offset", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		// Offset 0 is not a valid sub-directory offset.
		err = tiffFile.TIFFSetSubDirectory(ctx, 0)
		Expect(err).To(HaveOccurred())
	})
})
