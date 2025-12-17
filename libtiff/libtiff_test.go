package libtiff_test

import (
	"context"
	"image"
	"image/color"
	"log"
	"os"
	"path"
	"sync"

	"github.com/klippa-app/go-libtiff/internal/imports"
	"github.com/klippa-app/go-libtiff/libtiff"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	"github.com/tetratelabs/wazero"
)

var instance *libtiff.Instance

var _ = BeforeSuite(func() {
	// Set ENV to ensure resulting values.
	err := os.Setenv("TZ", "UTC")
	Expect(err).To(BeNil())

	instance, err = libtiff.GetInstance(context.Background(), &libtiff.Config{
		FSConfig: wazero.NewFSConfig().WithDirMount("../testdata", "/testdata"),
	})
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	err := instance.Close(context.Background())
	Expect(err).To(BeNil())

	Eventually(Goroutines).ShouldNot(HaveLeaked())

	// Check if all files are closed.
	Expect(imports.FileReaders.Refs).To(HaveLen(0))
})

var _ = Describe("files", func() {
	Context("a normal tiff file", func() {
		When("is opened with a file path", func() {
			It("opens without errors", func() {
				ctx := context.Background()
				tiffFile, err := instance.TIFFOpenFileFromPath(ctx, "/testdata/lena512color.jpeg.tiff", nil)
				Expect(err).To(BeNil())
				Expect(tiffFile).To(Not(BeNil()))
				if tiffFile != nil {
					defer tiffFile.Close(ctx)
				}
			})
		})

		When("is opened with a file reader", func() {
			It("opens without errors", func() {
				ctx := context.Background()

				filePath := "../testdata/lena512color.jpeg.tiff"
				f, err := os.Open(filePath)
				Expect(err).To(BeNil())
				Expect(f).To(Not(BeNil()))
				defer f.Close()

				stat, err := f.Stat()
				Expect(err).To(BeNil())
				Expect(stat).To(Not(BeNil()))

				tiffFile, err := instance.TIFFOpenFileFromReader(ctx, path.Base(filePath), f, uint64(stat.Size()), nil)
				Expect(err).To(BeNil())
				Expect(tiffFile).To(Not(BeNil()))
				defer tiffFile.Close(ctx)
			})
		})
	})

	Context("a multipage tiff file", func() {
		var file *libtiff.File
		var realFile *os.File

		BeforeEach(func() {
			filePath := "../testdata/multipage-sample.tif"
			f, err := os.Open(filePath)
			Expect(err).To(BeNil())
			Expect(f).To(Not(BeNil()))
			realFile = f

			stat, err := f.Stat()
			Expect(err).To(BeNil())
			Expect(stat).To(Not(BeNil()))

			tiffFile, err := instance.TIFFOpenFileFromReader(context.Background(), path.Base(filePath), f, uint64(stat.Size()), nil)
			Expect(err).To(BeNil())
			Expect(tiffFile).To(Not(BeNil()))
			file = tiffFile
		})

		AfterEach(func() {
			file.Close(context.Background())
			realFile.Close()
		})

		It("allows traversing the directories", func() {
			files := []int{}
			for i, err := range file.Directories(context.Background()) {
				files = append(files, i)
				Expect(err).To(BeNil())
			}
			Expect(files).To(Equal([]int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			}))
		})
	})
})

var _ = Describe("tags", func() {
	Context("a normal tiff file", func() {
		var file *libtiff.File

		BeforeEach(func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/lena512color.jpeg.tiff", nil)
			Expect(err).To(BeNil())
			Expect(tiffFile).To(Not(BeNil()))
			file = tiffFile
		})

		AfterEach(func() {
			file.Close(context.Background())
		})

		It("returns the correct dimensions", func() {
			width, height, err := file.GetDimensions(context.Background())
			Expect(err).To(BeNil())
			Expect(width).To(Equal(512))
			Expect(height).To(Equal(512))
		})

		It("returns the correct Uint16_t tag values", func() {
			val, err := file.TIFFGetFieldUint16_t(context.Background(), libtiff.TIFFTAG_BITSPERSAMPLE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint16(8)))
		})

		It("returns the correct Uint32_t tag values", func() {
			val, err := file.TIFFGetFieldUint32_t(context.Background(), libtiff.TIFFTAG_IMAGEWIDTH)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(512)))

			val, err = file.TIFFGetFieldUint32_t(context.Background(), libtiff.TIFFTAG_IMAGELENGTH)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(512)))
		})

		It("returns the correct int tag values", func() {
			val, err := file.TIFFGetFieldInt(context.Background(), libtiff.TIFFTAG_JPEGQUALITY)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(int(75)))
		})

		/*
			// We don't have a sample with doubles.
				It("returns the correct double tag values", func() {
					val, err := file.TIFFGetFieldInt(context.Background(), libtiff.TIFFTAG_STONITS)
					Expect(err).To(BeNil())
					Expect(val).To(Equal(int(75)))
				})
		*/

		It("returns the correct const char values", func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/multipage-sample.tif", nil)
			Expect(err).To(BeNil())
			Expect(tiffFile).To(Not(BeNil()))
			defer tiffFile.Close(context.Background())

			val, err := tiffFile.TIFFGetFieldConstChar(context.Background(), libtiff.TIFFTAG_SOFTWARE)
			Expect(err).To(BeNil())
			Expect(val).To(Equal("IrfanView"))
		})

		It("returns a typed error when a tag is of a different type", func() {
			val, err := file.TIFFGetFieldUint32_t(context.Background(), libtiff.TIFFTAG_XRESOLUTION)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_XRESOLUTION,
			}))
			Expect(val).To(Equal(uint32(0)))
		})

		It("returns a typed error when a tag can't be found", func() {
			val, err := file.TIFFGetFieldFloat(context.Background(), libtiff.TIFFTAG_XRESOLUTION)
			Expect(err).To(Equal(&libtiff.TagNotDefinedError{
				Tag: libtiff.TIFFTAG_XRESOLUTION,
			}))
			Expect(val).To(Equal(float32(0)))
		})

		Context("with no resolution tags", func() {
			It("returns the right resolution", func() {
				x, y, err := file.GetResolution(context.Background())
				Expect(err).To(Equal(&libtiff.TagNotDefinedError{
					Tag: libtiff.TIFFTAG_XRESOLUTION,
				}))
				Expect(x).To(Equal(float32(0)))
				Expect(y).To(Equal(float32(0)))
			})
		})
	})

	Context("a tiff file with the resolution tag", func() {
		var file *libtiff.File

		BeforeEach(func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/multipage-sample.tif", nil)
			Expect(err).To(BeNil())
			Expect(tiffFile).To(Not(BeNil()))
			file = tiffFile
		})

		AfterEach(func() {
			file.Close(context.Background())
		})

		It("returns the right resolution", func() {
			x, y, err := file.GetResolution(context.Background())
			Expect(err).To(BeNil())
			Expect(x).To(Equal(float32(96)))
			Expect(y).To(Equal(float32(96)))
		})

		It("returns the correct float tag values", func() {
			val, err := file.TIFFGetFieldFloat(context.Background(), libtiff.TIFFTAG_XRESOLUTION)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(float32(96)))

			val, err = file.TIFFGetFieldFloat(context.Background(), libtiff.TIFFTAG_YRESOLUTION)
			Expect(err).To(BeNil())
			Expect(val).To(Equal(float32(96)))
		})
	})
})

var _ = Describe("limits", func() {
	Context("a normal tiff file", func() {
		var maxMemory = int32(30)

		It("returns an error when the MaxSingleMemAlloc is reached", func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/multipage-sample.tif", &libtiff.OpenOptions{
				MaxSingleMemAlloc: &maxMemory,
			})
			Expect(err).To(MatchError(ContainSubstring("Memory allocation of 983 bytes is beyond the 30 byte limit defined in open options")))
			Expect(tiffFile).To(BeNil())
		})

		It("returns an error when the MaxSingleMemAlloc is reached", func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/multipage-sample.tif", &libtiff.OpenOptions{
				MaxCumulatedMemAlloc: &maxMemory,
			})
			Expect(err).To(MatchError(ContainSubstring("Memory allocation of 983 bytes is beyond the 30 cumulated byte limit defined in open options")))
			Expect(tiffFile).To(BeNil())
		})
	})
})

var _ = Describe("directory", func() {
	Context("a normal tiff file", func() {
		var file *libtiff.File

		BeforeEach(func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/multipage-sample.tif", nil)
			Expect(err).To(BeNil())
			Expect(tiffFile).To(Not(BeNil()))
			file = tiffFile
		})

		AfterEach(func() {
			file.Close(context.Background())
		})

		It("correctly returns the number of directories", func() {
			val, err := file.TIFFNumberOfDirectories(context.Background())
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(10)))
		})

		It("correctly returns whether we are at the last directory", func() {
			val, err := file.TIFFLastDirectory(context.Background())
			Expect(err).To(BeNil())
			Expect(val).To(Equal(false))

			// Move to last directory.
			err = file.TIFFSetDirectory(context.Background(), 9)
			Expect(err).To(BeNil())

			val, err = file.TIFFLastDirectory(context.Background())
			Expect(err).To(BeNil())
			Expect(val).To(Equal(true))
		})

		It("returns an error when moving past the directory list", func() {
			err := file.TIFFSetDirectory(context.Background(), 10)
			Expect(err).To(Not(BeNil()))
		})

		It("allows changing the directory", func() {
			val, err := file.TIFFCurrentDirectory(context.Background())
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(0)))

			// Move to last directory.
			err = file.TIFFSetDirectory(context.Background(), 5)
			Expect(err).To(BeNil())

			val, err = file.TIFFCurrentDirectory(context.Background())
			Expect(err).To(BeNil())
			Expect(val).To(Equal(uint32(5)))
		})

		It("allows traversing the directories", func() {
			files := []int{}
			for i, err := range file.Directories(context.Background()) {
				files = append(files, i)
				Expect(err).To(BeNil())
			}
			Expect(files).To(Equal([]int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			}))
		})
	})
})

var _ = Describe("image", func() {
	Context("a normal tiff file", func() {
		var file *libtiff.File

		BeforeEach(func() {
			tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/lena512color.jpeg.tiff", nil)
			Expect(err).To(BeNil())
			Expect(tiffFile).To(Not(BeNil()))
			file = tiffFile
		})

		AfterEach(func() {
			file.Close(context.Background())
		})

		It("allows rendering the image to a Go image", func() {
			goImage, cleanup, err := file.ToGoImage(context.Background())
			Expect(err).To(BeNil())
			Expect(goImage).To(Not(BeNil()))
			defer cleanup(context.Background())
			Expect(goImage.Bounds()).To(Equal(image.Rectangle{
				Min: image.Point{
					X: 0, Y: 0,
				},
				Max: image.Point{
					X: 512, Y: 512,
				},
			}))
			Expect(goImage.ColorModel()).To(Equal(color.RGBAModel))
		})

		It("allows rendering the image to a JPEG image", func() {
			image, err := file.ToImage(context.Background(), &libtiff.ImageOptions{
				OutputFormat:  libtiff.ImageOptionsOutputFormatJPEG,
				OutputTarget:  libtiff.ImageOptionsOutputTargetBytes,
				OutputQuality: 95,
			})
			Expect(err).To(BeNil())
			Expect(image).To(Not(BeNil()))
			Expect(image).To(Or(HaveLen(64150), HaveLen(65413))) // First is Go, second is libjpegturbo.
		})

		It("allows rendering the image to a PNG image", func() {
			image, err := file.ToImage(context.Background(), &libtiff.ImageOptions{
				OutputFormat:  libtiff.ImageOptionsOutputFormatPNG,
				OutputTarget:  libtiff.ImageOptionsOutputTargetBytes,
				OutputQuality: 95,
			})
			Expect(err).To(BeNil())
			Expect(image).To(Not(BeNil()))
			Expect(image).To(HaveLen(321744))
		})
	})
})

var _ = Describe("multithreading", func() {
	Context("a normal tiff file", func() {
		It("allows multiple tiff files to be processed at the same time", func() {
			// 4 concurrent tiff files.
			sem := make(chan struct{}, 4)
			var wg sync.WaitGroup

			for i := 0; i < 100; i++ {
				wg.Add(1)

				go func(id int) {
					defer wg.Done()
					defer GinkgoRecover()

					sem <- struct{}{}
					defer func() { <-sem }()

					tiffFile, err := instance.TIFFOpenFileFromPath(context.Background(), "/testdata/multipage-sample.tif", nil)
					Expect(err).To(BeNil())
					Expect(tiffFile).To(Not(BeNil()))
					defer func() {
						err = tiffFile.Close(context.Background())
						if err != nil {
							log.Printf("Close failed: %v", err)
						}
					}()

					val, err := tiffFile.TIFFCurrentDirectory(context.Background())
					Expect(err).To(BeNil())
					Expect(val).To(Equal(uint32(0)))

					image, cleanup, err := tiffFile.ToGoImage(context.Background())
					Expect(err).To(BeNil())
					Expect(image).To(Not(BeNil()))
					err = cleanup(context.Background())
					Expect(err).To(BeNil())

					// Move to last directory.
					err = tiffFile.TIFFSetDirectory(context.Background(), 5)
					Expect(err).To(BeNil())

					val, err = tiffFile.TIFFCurrentDirectory(context.Background())
					Expect(err).To(BeNil())
					Expect(val).To(Equal(uint32(5)))

					image, cleanup, err = tiffFile.ToGoImage(context.Background())
					Expect(err).To(BeNil())
					Expect(image).To(Not(BeNil()))
					err = cleanup(context.Background())
					Expect(err).To(BeNil())
				}(i)
			}

			wg.Wait()
		})
	})
})
