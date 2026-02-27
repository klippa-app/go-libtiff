package libtiff_test

import (
	"context"
	"os"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TIFFStripSize", func() {
	ctx := context.Background()

	It("returns the strip size for a strip-based image", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		size, err := tiffFile.TIFFStripSize(ctx)
		Expect(err).To(BeNil())
		Expect(size).To(BeNumerically(">", 0))
	})
})

var _ = Describe("TIFFNumberOfStrips", func() {
	ctx := context.Background()

	It("returns the number of strips for a strip-based image", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		count, err := tiffFile.TIFFNumberOfStrips(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(BeNumerically(">", 0))
	})
})

var _ = Describe("TIFFTileSize", func() {
	ctx := context.Background()

	It("returns width * rowsperstrip * spp for a strip-based image", func() {
		// After TIFFReadDirectory, libtiff sets td_tilewidth=imagewidth and
		// td_tilelength=rowsperstrip for strip-based images. So TIFFTileSize
		// returns imagewidth * rowsperstrip * samplesPerPixel, which may be
		// larger than TIFFStripSize (which caps to the actual image height).
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		tiled, err := tiffFile.TIFFIsTiled(ctx)
		Expect(err).To(BeNil())
		Expect(tiled).To(BeFalse())

		tileSize, err := tiffFile.TIFFTileSize(ctx)
		Expect(err).To(BeNil())

		width, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
		Expect(err).To(BeNil())

		rps, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP)
		Expect(err).To(BeNil())

		spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
		Expect(err).To(BeNil())

		expectedTileSize := int64(width) * int64(rps) * int64(spp)
		Expect(tileSize).To(Equal(expectedTileSize))
	})

	It("returns the correct tile size for a tiled image", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-tilesize-tiled-test-*.tif")
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
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH, 128)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH, 128)).To(Succeed())

		// Write 4 tiles (2x2 grid of 128x128 tiles for a 256x256 image).
		tileData := make([]byte, 128*128)
		for i := uint32(0); i < 4; i++ {
			Expect(writeTiff.TIFFWriteEncodedTile(ctx, i, tileData)).To(Succeed())
		}
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

		tiled, err := readTiff.TIFFIsTiled(ctx)
		Expect(err).To(BeNil())
		Expect(tiled).To(BeTrue())

		// 128 * 128 * 1 byte per sample = 16384 bytes per tile.
		size, err := readTiff.TIFFTileSize(ctx)
		Expect(err).To(BeNil())
		Expect(size).To(Equal(int64(16384)))
	})
})

var _ = Describe("TIFFNumberOfTiles", func() {
	ctx := context.Background()

	It("returns the number of strips for a strip-based image", func() {
		// After TIFFReadDirectory, libtiff unifies strip/tile handling internally.
		// For strip-based images, TIFFNumberOfTiles equals TIFFNumberOfStrips.
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		tiled, err := tiffFile.TIFFIsTiled(ctx)
		Expect(err).To(BeNil())
		Expect(tiled).To(BeFalse())

		numTiles, err := tiffFile.TIFFNumberOfTiles(ctx)
		Expect(err).To(BeNil())

		numStrips, err := tiffFile.TIFFNumberOfStrips(ctx)
		Expect(err).To(BeNil())

		Expect(numTiles).To(Equal(numStrips))
	})

	It("returns the correct tile count for a tiled image", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-numtiles-test-*.tif")
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
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH, 128)).To(Succeed())
		Expect(writeTiff.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH, 128)).To(Succeed())

		tileData := make([]byte, 128*128)
		for i := uint32(0); i < 4; i++ {
			Expect(writeTiff.TIFFWriteEncodedTile(ctx, i, tileData)).To(Succeed())
		}
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

		tiled, err := readTiff.TIFFIsTiled(ctx)
		Expect(err).To(BeNil())
		Expect(tiled).To(BeTrue())

		// 256x256 image with 128x128 tiles = 4 tiles.
		count, err := readTiff.TIFFNumberOfTiles(ctx)
		Expect(err).To(BeNil())
		Expect(count).To(Equal(uint32(4)))
	})
})

var _ = Describe("TIFFComputeStrip", func() {
	ctx := context.Background()

	It("returns strip 0 for row 0", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		strip, err := tiffFile.TIFFComputeStrip(ctx, 0, 0)
		Expect(err).To(BeNil())
		Expect(strip).To(Equal(uint32(0)))
	})

	It("returns a higher strip number for rows beyond the first strip", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
			RowsPerStrip: 4,
		})
		defer cleanup()

		// Row 0 should be in strip 0.
		strip0, err := tiffFile.TIFFComputeStrip(ctx, 0, 0)
		Expect(err).To(BeNil())
		Expect(strip0).To(Equal(uint32(0)))

		// Row 4 should be in strip 1.
		strip1, err := tiffFile.TIFFComputeStrip(ctx, 4, 0)
		Expect(err).To(BeNil())
		Expect(strip1).To(Equal(uint32(1)))

		// Row 12 should be in strip 3.
		strip3, err := tiffFile.TIFFComputeStrip(ctx, 12, 0)
		Expect(err).To(BeNil())
		Expect(strip3).To(Equal(uint32(3)))
	})
})

var _ = Describe("TIFFComputeTile", func() {
	ctx := context.Background()

	It("returns 0 for a non-tiled image", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		tile, err := tiffFile.TIFFComputeTile(ctx, 0, 0, 0, 0)
		Expect(err).To(BeNil())
		Expect(tile).To(Equal(uint32(0)))
	})

	It("returns different tile numbers for different coordinates in a tiled image", func() {
		img := createTestRGBA(64, 64)
		tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
			TileWidth:  32,
			TileHeight: 32,
		})
		defer cleanup()

		// Top-left tile.
		tile0, err := tiffFile.TIFFComputeTile(ctx, 0, 0, 0, 0)
		Expect(err).To(BeNil())
		Expect(tile0).To(Equal(uint32(0)))

		// Top-right tile.
		tile1, err := tiffFile.TIFFComputeTile(ctx, 32, 0, 0, 0)
		Expect(err).To(BeNil())
		Expect(tile1).To(Equal(uint32(1)))

		// Bottom-left tile.
		tile2, err := tiffFile.TIFFComputeTile(ctx, 0, 32, 0, 0)
		Expect(err).To(BeNil())
		Expect(tile2).To(Equal(uint32(2)))

		// Bottom-right tile.
		tile3, err := tiffFile.TIFFComputeTile(ctx, 32, 32, 0, 0)
		Expect(err).To(BeNil())
		Expect(tile3).To(Equal(uint32(3)))
	})
})

var _ = Describe("TIFFIsTiled", func() {
	ctx := context.Background()

	It("returns false for a strip-based image", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		tiled, err := tiffFile.TIFFIsTiled(ctx)
		Expect(err).To(BeNil())
		Expect(tiled).To(BeFalse())
	})
})

var _ = Describe("TIFFScanlineSize", func() {
	ctx := context.Background()

	It("returns the scanline size for a strip-based image", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		size, err := tiffFile.TIFFScanlineSize(ctx)
		Expect(err).To(BeNil())
		Expect(size).To(BeNumerically(">", 0))
	})

	It("returns width * spp for a written RGBA image", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		size, err := tiffFile.TIFFScanlineSize(ctx)
		Expect(err).To(BeNil())
		// 16 pixels * 4 bytes/pixel = 64.
		Expect(size).To(Equal(int64(64)))
	})
})

var _ = Describe("TIFFVStripSize", func() {
	ctx := context.Background()

	It("returns the size of a strip with a given number of rows", func() {
		tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		size, err := tiffFile.TIFFVStripSize(ctx, 10)
		Expect(err).To(BeNil())
		Expect(size).To(BeNumerically(">", 0))
	})

	It("equals nrows * scanline size for a written image", func() {
		img := createTestRGBA(16, 16)
		tiffFile, cleanup := writeAndReopen(ctx, img, nil)
		defer cleanup()

		scanlineSize, err := tiffFile.TIFFScanlineSize(ctx)
		Expect(err).To(BeNil())

		vstripSize, err := tiffFile.TIFFVStripSize(ctx, 4)
		Expect(err).To(BeNil())
		Expect(vstripSize).To(Equal(scanlineSize * 4))
	})
})

var _ = Describe("TIFFDefaultTileSize", func() {
	ctx := context.Background()

	It("returns default tile dimensions", func() {
		tmpFile, err := os.CreateTemp("", "libtiff-tilesize-test-*.tif")
		Expect(err).To(BeNil())
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		fileMode := "w"
		tiffFile, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
			FileMode: &fileMode,
		})
		Expect(err).To(BeNil())
		defer tiffFile.Close(ctx)

		Expect(tiffFile.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH, 256)).To(Succeed())
		Expect(tiffFile.TIFFSetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH, 256)).To(Succeed())
		Expect(tiffFile.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE, 8)).To(Succeed())
		Expect(tiffFile.TIFFSetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL, 3)).To(Succeed())

		tw, th, err := tiffFile.TIFFDefaultTileSize(ctx)
		Expect(err).To(BeNil())
		Expect(tw).To(BeNumerically(">", 0))
		Expect(th).To(BeNumerically(">", 0))
	})
})
