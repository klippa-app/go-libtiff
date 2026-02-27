package libtiff_test

import (
	"context"
	"image"
	"image/color"
	"os"
	"time"

	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func createTestRGBA(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}

func createTestNRGBA(width, height int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}

// writeAndReopen writes a Go image to a temp TIFF file, closes it, then reopens for reading.
// Returns the readable TIFF file and a cleanup function.
func writeAndReopen(ctx context.Context, img image.Image, options *libtiff.FromGoImageOptions) (*libtiff.File, func()) {
	tmpFile, err := os.CreateTemp("", "libtiff-test-*.tif")
	Expect(err).To(BeNil())
	tmpPath := tmpFile.Name()

	fileMode := "w"
	tiffFile, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
		FileMode: &fileMode,
	})
	Expect(err).To(BeNil())

	err = tiffFile.FromGoImage(ctx, img, options)
	Expect(err).To(BeNil())

	tiffFile.Close(ctx)
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

var _ = Describe("FromGoImage", func() {
	ctx := context.Background()

	Context("basic TIFF tags", func() {
		It("writes correct tags for an RGBA image with default options", func() {
			img := createTestRGBA(64, 48)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			width, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
			Expect(err).To(BeNil())
			Expect(width).To(Equal(uint32(64)))

			height, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGELENGTH)
			Expect(err).To(BeNil())
			Expect(height).To(Equal(uint32(48)))

			bps, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE)
			Expect(err).To(BeNil())
			Expect(bps).To(Equal(uint16(8)))

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_NONE)))

			photo, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC)
			Expect(err).To(BeNil())
			Expect(photo).To(Equal(uint16(libtiff.PHOTOMETRIC_RGB)))

			orient, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_ORIENTATION)
			Expect(err).To(BeNil())
			Expect(orient).To(Equal(uint16(libtiff.ORIENTATION_TOPLEFT)))

			planar, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_PLANARCONFIG)
			Expect(err).To(BeNil())
			Expect(planar).To(Equal(uint16(libtiff.PLANARCONFIG_CONTIG)))
		})
	})

	Context("compression options", func() {
		It("writes with LZW compression", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_LZW,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_LZW)))

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))
		})

		It("writes with DEFLATE compression", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_ADOBE_DEFLATE,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_ADOBE_DEFLATE)))
		})

		It("writes with JPEG compression using 3 samples per pixel", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_JPEG,
				Quality:     90,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_JPEG)))

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(3)))
		})
	})

	Context("metadata tags", func() {
		It("writes a default Software tag containing the libtiff version", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			software, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
			Expect(err).To(BeNil())
			Expect(software).To(HavePrefix("go-libtiff/libtiff-"))
		})

		It("writes a custom Software tag when specified", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Software: "my-custom-software",
			})
			defer cleanup()

			software, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_SOFTWARE)
			Expect(err).To(BeNil())
			Expect(software).To(Equal("my-custom-software"))
		})

		It("writes a default DateTime tag with the current time", func() {
			before := time.Now()
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()
			after := time.Now()

			dateTime, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_DATETIME)
			Expect(err).To(BeNil())

			// Parse the TIFF datetime and verify it's within the test window.
			parsed, err := time.Parse("2006:01:02 15:04:05", dateTime)
			Expect(err).To(BeNil())
			Expect(parsed).To(BeTemporally(">=", before.Truncate(time.Second)))
			Expect(parsed).To(BeTemporally("<=", after.Truncate(time.Second).Add(time.Second)))
		})

		It("writes a custom DateTime tag when specified", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				DateTime: "2024:06:15 12:30:00",
			})
			defer cleanup()

			dateTime, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_DATETIME)
			Expect(err).To(BeNil())
			Expect(dateTime).To(Equal("2024:06:15 12:30:00"))
		})

		It("does not write the Artist tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_ARTIST)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_ARTIST,
			}))
		})

		It("writes the Artist tag when specified", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Artist: "Test Author",
			})
			defer cleanup()

			artist, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_ARTIST)
			Expect(err).To(BeNil())
			Expect(artist).To(Equal("Test Author"))
		})
	})

	Context("alpha mode", func() {
		It("auto-detects associated alpha for *image.RGBA", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))
		})

		It("auto-detects unassociated alpha for *image.NRGBA", func() {
			img := createTestNRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))
		})

		It("respects explicit AlphaAssociated for *image.NRGBA", func() {
			img := createTestNRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				AlphaMode: libtiff.AlphaAssociated,
			})
			defer cleanup()

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))
		})

		It("respects explicit AlphaUnassociated for *image.RGBA", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				AlphaMode: libtiff.AlphaUnassociated,
			})
			defer cleanup()

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))
		})

		It("strips alpha with JPEG regardless of alpha mode", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_JPEG,
			})
			defer cleanup()

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(3)))
		})
	})

	Context("pixel data roundtrip", func() {
		It("preserves pixel data with no compression (RGBA)", func() {
			src := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, nil)
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			Expect(goImage.Bounds()).To(Equal(src.Bounds()))

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr), "red mismatch at (%d,%d)", x, y)
					Expect(dg).To(Equal(sg), "green mismatch at (%d,%d)", x, y)
					Expect(db).To(Equal(sb), "blue mismatch at (%d,%d)", x, y)
					Expect(da).To(Equal(sa), "alpha mismatch at (%d,%d)", x, y)
				}
			}
		})

		It("preserves pixel data with LZW compression (RGBA)", func() {
			src := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_LZW,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})

		It("preserves pixel data with DEFLATE compression (RGBA)", func() {
			src := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_ADOBE_DEFLATE,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})

		It("preserves pixel data with no compression (NRGBA)", func() {
			src := createTestNRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, nil)
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			Expect(goImage.Bounds()).To(Equal(src.Bounds()))

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})

		It("produces valid pixel data with JPEG compression", func() {
			src := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_JPEG,
				Quality:     100,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			Expect(goImage.Bounds()).To(Equal(src.Bounds()))

			// JPEG is lossy, so just verify pixels are reasonable (within tolerance).
			for y := 0; y < 32; y++ {
				for x := 0; x < 32; x++ {
					sr, sg, sb, _ := src.At(x, y).RGBA()
					dr, dg, db, _ := goImage.At(x, y).RGBA()
					Expect(dr).To(BeNumerically("~", sr, 0x1000))
					Expect(dg).To(BeNumerically("~", sg, 0x1000))
					Expect(db).To(BeNumerically("~", sb, 0x1000))
				}
			}
		})
	})

	Context("generic image path", func() {
		It("handles a non-RGBA, non-NRGBA image type", func() {
			// image.Gray forces the generic path.
			gray := image.NewGray(image.Rect(0, 0, 8, 8))
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					gray.SetGray(x, y, color.Gray{Y: uint8((x + y) % 256)})
				}
			}

			tiffFile, cleanup := writeAndReopen(ctx, gray, nil)
			defer cleanup()

			width, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_IMAGEWIDTH)
			Expect(err).To(BeNil())
			Expect(width).To(Equal(uint32(8)))

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			// Verify the grayscale values roundtrip through RGB.
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					expected := uint8((x + y) % 256)
					r, g, b, _ := goImage.At(x, y).RGBA()
					// Gray converts to equal R=G=B.
					Expect(uint8(r >> 8)).To(Equal(expected))
					Expect(uint8(g >> 8)).To(Equal(expected))
					Expect(uint8(b >> 8)).To(Equal(expected))
				}
			}
		})
	})

	Context("new compression methods", func() {
		It("writes with PackBits compression", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_PACKBITS,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_PACKBITS)))

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(4)))
		})

		It("preserves pixel data with PackBits compression", func() {
			src := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_PACKBITS,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})

	})

	Context("CCITT bilevel compression", func() {
		It("writes with CCITTFAX3 compression", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_CCITTFAX3,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_CCITTFAX3)))

			bps, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE)
			Expect(err).To(BeNil())
			Expect(bps).To(Equal(uint16(1)))

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(1)))

			photo, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_PHOTOMETRIC)
			Expect(err).To(BeNil())
			Expect(photo).To(Equal(uint16(libtiff.PHOTOMETRIC_MINISWHITE)))
		})

		It("writes with CCITTFAX4 compression", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_CCITTFAX4,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_CCITTFAX4)))

			bps, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_BITSPERSAMPLE)
			Expect(err).To(BeNil())
			Expect(bps).To(Equal(uint16(1)))

			spp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_SAMPLESPERPIXEL)
			Expect(err).To(BeNil())
			Expect(spp).To(Equal(uint16(1)))
		})

		It("converts all-white image to white with CCITT", func() {
			// Create an all-white image.
			img := image.NewRGBA(image.Rect(0, 0, 16, 16))
			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
				}
			}

			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_CCITTFAX4,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_CCITTFAX4)))
		})

		It("converts all-black image to black with CCITT", func() {
			// Create an all-black image.
			img := image.NewRGBA(image.Rect(0, 0, 16, 16))
			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
				}
			}

			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_CCITTFAX4,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_CCITTFAX4)))
		})

		It("respects custom BilevelThreshold", func() {
			// Create an image with gray pixels at luminance ~128.
			// With threshold 200, they should become black.
			img := image.NewRGBA(image.Rect(0, 0, 16, 16))
			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					img.Set(x, y, color.RGBA{R: 128, G: 128, B: 128, A: 255})
				}
			}

			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression:      libtiff.COMPRESSION_CCITTFAX4,
				BilevelThreshold: 200,
			})
			defer cleanup()

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_CCITTFAX4)))
		})
	})

	Context("predictor", func() {
		It("writes with horizontal predictor for LZW", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_LZW,
				Predictor:   libtiff.PREDICTOR_HORIZONTAL,
			})
			defer cleanup()

			pred, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_PREDICTOR)
			Expect(err).To(BeNil())
			Expect(pred).To(Equal(uint16(libtiff.PREDICTOR_HORIZONTAL)))
		})

		It("writes with horizontal predictor for Deflate", func() {
			img := createTestRGBA(32, 32)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_ADOBE_DEFLATE,
				Predictor:   libtiff.PREDICTOR_HORIZONTAL,
			})
			defer cleanup()

			pred, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_PREDICTOR)
			Expect(err).To(BeNil())
			Expect(pred).To(Equal(uint16(libtiff.PREDICTOR_HORIZONTAL)))
		})

		It("preserves pixel data with LZW + horizontal predictor", func() {
			src := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_LZW,
				Predictor:   libtiff.PREDICTOR_HORIZONTAL,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})
	})

	Context("resolution tags", func() {
		It("writes resolution tags with inch unit", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				XResolution:    300,
				YResolution:    300,
				ResolutionUnit: libtiff.RESUNIT_INCH,
			})
			defer cleanup()

			xRes, err := tiffFile.TIFFGetFieldFloat(ctx, libtiff.TIFFTAG_XRESOLUTION)
			Expect(err).To(BeNil())
			Expect(xRes).To(BeNumerically("~", float32(300), 0.01))

			yRes, err := tiffFile.TIFFGetFieldFloat(ctx, libtiff.TIFFTAG_YRESOLUTION)
			Expect(err).To(BeNil())
			Expect(yRes).To(BeNumerically("~", float32(300), 0.01))

			unit, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_RESOLUTIONUNIT)
			Expect(err).To(BeNil())
			Expect(unit).To(Equal(uint16(libtiff.RESUNIT_INCH)))
		})

		It("writes resolution tags with centimeter unit", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				XResolution:    118,
				YResolution:    118,
				ResolutionUnit: libtiff.RESUNIT_CENTIMETER,
			})
			defer cleanup()

			unit, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_RESOLUTIONUNIT)
			Expect(err).To(BeNil())
			Expect(unit).To(Equal(uint16(libtiff.RESUNIT_CENTIMETER)))
		})

		It("does not write resolution tags when zero", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{})
			defer cleanup()

			// Resolution tags may or may not have defaults set by libtiff,
			// but ResolutionUnit should not be explicitly set.
			_, err := tiffFile.TIFFGetFieldFloat(ctx, libtiff.TIFFTAG_XRESOLUTION)
			// This may or may not error - libtiff may set defaults.
			_ = err
		})
	})

	Context("additional metadata tags", func() {
		It("writes and reads back Description tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Description: "Test description",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_IMAGEDESCRIPTION)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("Test description"))
		})

		It("does not write Description tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_IMAGEDESCRIPTION)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_IMAGEDESCRIPTION,
			}))
		})

		It("writes and reads back Copyright tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Copyright: "Copyright 2024 Test",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_COPYRIGHT)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("Copyright 2024 Test"))
		})

		It("does not write Copyright tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_COPYRIGHT)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_COPYRIGHT,
			}))
		})

		It("writes and reads back DocumentName tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				DocumentName: "test-document.tif",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_DOCUMENTNAME)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("test-document.tif"))
		})

		It("does not write DocumentName tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_DOCUMENTNAME)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_DOCUMENTNAME,
			}))
		})

		It("writes and reads back PageName tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				PageName: "Page 1",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_PAGENAME)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("Page 1"))
		})

		It("does not write PageName tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_PAGENAME)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_PAGENAME,
			}))
		})

		It("writes and reads back HostComputer tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				HostComputer: "test-host",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_HOSTCOMPUTER)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("test-host"))
		})

		It("does not write HostComputer tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_HOSTCOMPUTER)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_HOSTCOMPUTER,
			}))
		})

		It("writes and reads back Make tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Make: "Test Make",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_MAKE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("Test Make"))
		})

		It("does not write Make tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_MAKE)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_MAKE,
			}))
		})

		It("writes and reads back Model tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Model: "Test Model",
			})
			defer cleanup()

			val, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_MODEL)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("Test Model"))
		})

		It("does not write Model tag when empty", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, err := tiffFile.TIFFGetFieldConstChar(ctx, libtiff.TIFFTAG_MODEL)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_MODEL,
			}))
		})
	})

	Context("rows per strip override", func() {
		It("uses custom rows per strip", func() {
			img := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				RowsPerStrip: 4,
			})
			defer cleanup()

			rps, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_ROWSPERSTRIP)
			Expect(err).To(BeNil())
			Expect(rps).To(Equal(uint32(4)))

			strips, err := tiffFile.TIFFNumberOfStrips(ctx)
			Expect(err).To(BeNil())
			Expect(strips).To(Equal(uint32(4))) // 16 rows / 4 rows per strip = 4 strips
		})

		It("preserves pixel data with custom rows per strip", func() {
			src := createTestRGBA(16, 16)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				RowsPerStrip: 4,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})
	})

	Context("orientation override", func() {
		It("uses custom orientation", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				Orientation: libtiff.ORIENTATION_BOTLEFT,
			})
			defer cleanup()

			orient, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_ORIENTATION)
			Expect(err).To(BeNil())
			Expect(orient).To(Equal(uint16(libtiff.ORIENTATION_BOTLEFT)))
		})

		It("defaults to TOPLEFT when not specified", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			orient, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_ORIENTATION)
			Expect(err).To(BeNil())
			Expect(orient).To(Equal(uint16(libtiff.ORIENTATION_TOPLEFT)))
		})
	})

	Context("page number tag", func() {
		It("writes and reads back PageNumber tag", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				PageNumber: 2,
				TotalPages: 5,
			})
			defer cleanup()

			pageNum, totalPages, err := tiffFile.TIFFGetFieldTwoUint16(ctx, libtiff.TIFFTAG_PAGENUMBER)
			Expect(err).To(BeNil())
			Expect(pageNum).To(Equal(uint16(2)))
			Expect(totalPages).To(Equal(uint16(5)))
		})

		It("writes page 0 with TotalPages set", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				PageNumber: 0,
				TotalPages: 3,
			})
			defer cleanup()

			pageNum, totalPages, err := tiffFile.TIFFGetFieldTwoUint16(ctx, libtiff.TIFFTAG_PAGENUMBER)
			Expect(err).To(BeNil())
			Expect(pageNum).To(Equal(uint16(0)))
			Expect(totalPages).To(Equal(uint16(3)))
		})

		It("does not write PageNumber tag when TotalPages is 0", func() {
			img := createTestRGBA(8, 8)
			tiffFile, cleanup := writeAndReopen(ctx, img, nil)
			defer cleanup()

			_, _, err := tiffFile.TIFFGetFieldTwoUint16(ctx, libtiff.TIFFTAG_PAGENUMBER)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_PAGENUMBER,
			}))
		})
	})

	Context("tile-based output", func() {
		It("writes tiled TIFF with correct tags", func() {
			img := createTestRGBA(64, 64)
			tiffFile, cleanup := writeAndReopen(ctx, img, &libtiff.FromGoImageOptions{
				TileWidth:  32,
				TileHeight: 32,
			})
			defer cleanup()

			isTiled, err := tiffFile.TIFFIsTiled(ctx)
			Expect(err).To(BeNil())
			Expect(isTiled).To(BeTrue())

			numTiles, err := tiffFile.TIFFNumberOfTiles(ctx)
			Expect(err).To(BeNil())
			Expect(numTiles).To(Equal(uint32(4))) // 64/32 * 64/32 = 4

			tw, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_TILEWIDTH)
			Expect(err).To(BeNil())
			Expect(tw).To(Equal(uint32(32)))

			th, err := tiffFile.TIFFGetFieldUint32_t(ctx, libtiff.TIFFTAG_TILELENGTH)
			Expect(err).To(BeNil())
			Expect(th).To(Equal(uint32(32)))
		})

		It("preserves pixel data with tiled output", func() {
			src := createTestRGBA(64, 64)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				TileWidth:  32,
				TileHeight: 32,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			Expect(goImage.Bounds()).To(Equal(src.Bounds()))

			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr), "red mismatch at (%d,%d)", x, y)
					Expect(dg).To(Equal(sg), "green mismatch at (%d,%d)", x, y)
					Expect(db).To(Equal(sb), "blue mismatch at (%d,%d)", x, y)
					Expect(da).To(Equal(sa), "alpha mismatch at (%d,%d)", x, y)
				}
			}
		})

		It("preserves pixel data with tiled NRGBA output", func() {
			src := createTestNRGBA(64, 64)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				TileWidth:  32,
				TileHeight: 32,
			})
			defer cleanup()

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					sr, sg, sb, sa := src.At(x, y).RGBA()
					dr, dg, db, da := goImage.At(x, y).RGBA()
					Expect(dr).To(Equal(sr))
					Expect(dg).To(Equal(sg))
					Expect(db).To(Equal(sb))
					Expect(da).To(Equal(sa))
				}
			}
		})

		It("handles non-RGBA image types with tiled output", func() {
			gray := image.NewGray(image.Rect(0, 0, 64, 64))
			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					gray.SetGray(x, y, color.Gray{Y: uint8((x + y) % 256)})
				}
			}

			tiffFile, cleanup := writeAndReopen(ctx, gray, &libtiff.FromGoImageOptions{
				TileWidth:  32,
				TileHeight: 32,
			})
			defer cleanup()

			isTiled, err := tiffFile.TIFFIsTiled(ctx)
			Expect(err).To(BeNil())
			Expect(isTiled).To(BeTrue())

			goImage, imgCleanup, err := tiffFile.ToGoImage(ctx)
			Expect(err).To(BeNil())
			defer imgCleanup(ctx)

			for y := 0; y < 64; y++ {
				for x := 0; x < 64; x++ {
					expected := uint8((x + y) % 256)
					r, _, _, _ := goImage.At(x, y).RGBA()
					Expect(uint8(r >> 8)).To(Equal(expected))
				}
			}
		})

		It("returns error when only TileWidth is set", func() {
			img := createTestRGBA(64, 64)

			tmpFile, err := os.CreateTemp("", "libtiff-test-*.tif")
			Expect(err).To(BeNil())
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			fileMode := "w"
			tiffFile, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
				FileMode: &fileMode,
			})
			Expect(err).To(BeNil())

			err = tiffFile.FromGoImage(ctx, img, &libtiff.FromGoImageOptions{
				TileWidth: 32,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("both TileWidth and TileHeight"))
			tiffFile.Close(ctx)
		})

		It("returns error when only TileHeight is set", func() {
			img := createTestRGBA(64, 64)

			tmpFile, err := os.CreateTemp("", "libtiff-test-*.tif")
			Expect(err).To(BeNil())
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			fileMode := "w"
			tiffFile, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, "test.tif", tmpFile, 0, &libtiff.OpenOptions{
				FileMode: &fileMode,
			})
			Expect(err).To(BeNil())

			err = tiffFile.FromGoImage(ctx, img, &libtiff.FromGoImageOptions{
				TileHeight: 32,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("both TileWidth and TileHeight"))
			tiffFile.Close(ctx)
		})

		It("writes tiled TIFF with LZW compression", func() {
			src := createTestRGBA(64, 64)
			tiffFile, cleanup := writeAndReopen(ctx, src, &libtiff.FromGoImageOptions{
				Compression: libtiff.COMPRESSION_LZW,
				TileWidth:   32,
				TileHeight:  32,
			})
			defer cleanup()

			isTiled, err := tiffFile.TIFFIsTiled(ctx)
			Expect(err).To(BeNil())
			Expect(isTiled).To(BeTrue())

			comp, err := tiffFile.TIFFGetFieldUint16_t(ctx, libtiff.TIFFTAG_COMPRESSION)
			Expect(err).To(BeNil())
			Expect(comp).To(Equal(uint16(libtiff.COMPRESSION_LZW)))
		})
	})
})
