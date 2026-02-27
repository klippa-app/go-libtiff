package libtiff_test

import (
	"context"
	"image"
	"os"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TIFFReadEncodedStrip", func() {
	ctx := context.Background()

	It("reads and decompresses a strip from an existing TIFF", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		data, err := tiffFile.TIFFReadEncodedStrip(ctx, 0)
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeEmpty())
	})

	It("roundtrips through write and read", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
			Compression: libtiff.COMPRESSION_LZW,
		})
		defer cleanup()

		data, err := tiffFile.TIFFReadEncodedStrip(ctx, 0)
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeEmpty())
	})
})

var _ = Describe("TIFFReadEncodedTile", func() {
	ctx := context.Background()

	It("returns an error for non-tiled images", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		// TIFFTileSize returns 0 for non-tiled images, so we expect an error
		// or empty result.
		_, err = tiffFile.TIFFReadEncodedTile(ctx, 0)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("TIFFReadScanline", func() {
	ctx := context.Background()

	It("reads a scanline from a written uncompressed TIFF", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		data, err := tiffFile.TIFFReadScanline(ctx, 0, 0)
		Expect(err).To(BeNil())
		// 16 pixels * 4 bytes/pixel = 64 bytes per scanline.
		Expect(data).To(HaveLen(64))
	})

	It("reads different data for different rows", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		row0, err := tiffFile.TIFFReadScanline(ctx, 0, 0)
		Expect(err).To(BeNil())

		row1, err := tiffFile.TIFFReadScanline(ctx, 1, 0)
		Expect(err).To(BeNil())

		Expect(row0).ToNot(Equal(row1))
	})
})

var _ = Describe("TIFFReadRGBAStrip", func() {
	ctx := context.Background()

	It("reads a strip as RGBA data", func() {
		img := createTestRGBA(8, 8)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		data, err := tiffFile.TIFFReadRGBAStrip(ctx, 0)
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeEmpty())
		// Each pixel is 4 bytes (RGBA).
		Expect(len(data) % 4).To(Equal(0))
	})
})

var _ = Describe("TIFFReadRGBATile", func() {
	ctx := context.Background()

	It("returns an error for non-tiled images", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		// Non-tiled images don't have TILEWIDTH/TILELENGTH tags.
		_, err = tiffFile.TIFFReadRGBATile(ctx, 0, 0)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("TIFFReadRawStrip", func() {
	ctx := context.Background()

	It("reads raw compressed data for a strip", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
			Compression: libtiff.COMPRESSION_LZW,
		})
		defer cleanup()

		data, err := tiffFile.TIFFReadRawStrip(ctx, 0)
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeEmpty())
	})

	It("returns raw data that is smaller than decoded data for compressed images", func() {
		// Create a highly compressible image (solid color).
		img := image.NewRGBA(image.Rect(0, 0, 64, 64))

		tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
			Compression: libtiff.COMPRESSION_LZW,
		})
		defer cleanup()

		raw, err := tiffFile.TIFFReadRawStrip(ctx, 0)
		Expect(err).To(BeNil())

		decoded, err := tiffFile.TIFFReadEncodedStrip(ctx, 0)
		Expect(err).To(BeNil())

		Expect(len(raw)).To(BeNumerically("<", len(decoded)))
	})
})

var _ = Describe("TIFFReadRawTile", func() {
	ctx := context.Background()

	It("returns an error for non-tiled images", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/multipage-sample.tif", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		_, err = tiffFile.TIFFReadRawTile(ctx, 0)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("TIFFWriteEncodedTile", func() {
	ctx := context.Background()

	It("writes a tile to a tiled TIFF", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-tile-write-test-*.tif")
		Expect(err).To(BeNil())
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		fileMode := "w"
		writeTiff, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())

		tileW := uint32(256)
		tileH := uint32(256)

		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 256)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 256)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 1)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC, uint16(libtiff.PHOTOMETRIC_MINISBLACK))).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH, tileW)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH, tileH)).To(Succeed())

		tileData := make([]byte, tileW*tileH)
		Expect(writeTiff.TIFFWriteEncodedTile(ctx, 0, tileData)).To(Succeed())
		Expect(writeTiff.TIFFWriteDirectory(ctx)).To(Succeed())

		writeTiff.Close(ctx)
		tmpFile.Close()

		// Verify we can read it back.
		readFile, err := os.Open(tmpPath)
		Expect(err).To(BeNil())
		defer readFile.Close()

		stat, err := readFile.Stat()
		Expect(err).To(BeNil())

		readTiff, err := instance.TIFFOpenFileFromReader(ctx, "test.tif", readFile, uint64(stat.Size()), nil)
		Expect(err).To(BeNil())
		defer readTiff.Close(ctx)

		tiled, err := readTiff.TIFFIsTiled(ctx)
		Expect(err).To(BeNil())
		Expect(tiled).To(BeTrue())
	})
})
