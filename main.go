package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	_ "github.com/klippa-app/go-libtiff/fax2ps"
	_ "github.com/klippa-app/go-libtiff/fax2tiff"
	"github.com/klippa-app/go-libtiff/internal/image/image_jpeg"
	"github.com/klippa-app/go-libtiff/internal/registry"
	"github.com/klippa-app/go-libtiff/libtiff"
	_ "github.com/klippa-app/go-libtiff/mkg3states"
	_ "github.com/klippa-app/go-libtiff/pal2rgb"
	_ "github.com/klippa-app/go-libtiff/ppm2tiff"
	_ "github.com/klippa-app/go-libtiff/raw2tiff"
	_ "github.com/klippa-app/go-libtiff/rgb2ycbcr"
	_ "github.com/klippa-app/go-libtiff/thumbnail"
	_ "github.com/klippa-app/go-libtiff/tiff2bw"
	_ "github.com/klippa-app/go-libtiff/tiff2pdf"
	_ "github.com/klippa-app/go-libtiff/tiff2ps"
	_ "github.com/klippa-app/go-libtiff/tiff2rgba"
	_ "github.com/klippa-app/go-libtiff/tiffcmp"
	_ "github.com/klippa-app/go-libtiff/tiffcp"
	_ "github.com/klippa-app/go-libtiff/tiffcrop"
	_ "github.com/klippa-app/go-libtiff/tiffdither"
	_ "github.com/klippa-app/go-libtiff/tiffdump"
	_ "github.com/klippa-app/go-libtiff/tiffinfo"
	_ "github.com/klippa-app/go-libtiff/tiffmedian"
	_ "github.com/klippa-app/go-libtiff/tiffset"
	_ "github.com/klippa-app/go-libtiff/tiffsplit"
	"github.com/tetratelabs/wazero/sys"

	"github.com/spf13/cobra"
	"github.com/tetratelabs/wazero"
)

var compilationCache wazero.CompilationCache

func init() {
	compilationCacheOption := os.Getenv("LIBTIFF_COMPILATION_CACHE_DIR")
	if compilationCacheOption == "memory" {
		compilationCache = wazero.NewCompilationCache()
	} else if compilationCacheOption != "" {
		var err error
		compilationCache, err = wazero.NewCompilationCacheWithDir(compilationCacheOption)
		if err != nil {
			log.Fatal(fmt.Errorf("could not create compilation cache directory"))
		}
	}
}

func main() {
	ctx := context.Background()

	availableBinaries := registry.List()
	availableBinaries = append(availableBinaries, "tiff2img")
	incorrectStartArgument := func() {
		log.Fatalf("You should minimally start the program with one of the following arguments: %s", strings.Join(availableBinaries, ", "))
	}
	if len(os.Args) < 2 {
		incorrectStartArgument()
	}
	if !slices.Contains(availableBinaries, os.Args[1]) {
		incorrectStartArgument()
	}

	if os.Args[1] == "tiff2img" {
		err := tiff2img()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := registry.Run(ctx, os.Args[1], os.Args[2:]...)
	if err != nil {
		if exitErr, ok := err.(*sys.ExitError); ok {
			os.Exit(int(exitErr.ExitCode()))
		} else {
			log.Fatal(err)
		}
	}
}

func tiff2img() error {
	var (
		// Used for flags.
		fileType    string
		quality     int
		progressive bool
	)

	rootCmd := &cobra.Command{
		Use:   "tiff2img [input] [output]",
		Short: "A CLI tool to convert tiff to images, if the file has multiple images, the output path has to contain %d",
		Args: func(cmd *cobra.Command, args []string) error {
			return cobra.ExactArgs(3)(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			input := args[1]

			openFile, err := os.Open(input)
			if err != nil {
				log.Fatal(err)
			}
			defer openFile.Close()

			stat, err := openFile.Stat()
			if err != nil {
				log.Fatal(err)
			}

			output := args[2]
			instance, err := libtiff.GetInstance(ctx, &libtiff.Config{
				CompilationCache: compilationCache,
			})
			if err != nil {
				log.Fatal(err)
			}
			defer instance.Close(ctx)

			file, err := instance.TIFFOpenFileFromReader(ctx, path.Base(input), openFile, uint64(stat.Size()), nil)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close(ctx)

			imageCount, err := file.TIFFNumberOfDirectories(ctx)
			if err != nil {
				log.Fatal(err)
			}

			if imageCount == 0 {
				log.Fatal(errors.New("tiff has no images"))
			}

			if imageCount > 1 && !strings.Contains(output, "%d") {
				log.Fatal(errors.New("tiff has multiple images and output path does not contain %d"))
			}

			for i := range file.Directories(ctx) {
				func() {
					renderedImage, cleanup, err := file.ToGoImage(ctx)
					if err != nil {
						log.Fatal(fmt.Errorf("could not convert tiff image %d to go image: %w", i, err))
					}
					defer cleanup(ctx)

					outputPath := output
					outputPath = strings.Replace(outputPath, "%d", fmt.Sprintf("%d", i), 1)
					outFile, err := os.Create(outputPath)
					if err != nil {
						log.Fatal(fmt.Errorf("could not create output path %s for tiff image %d: %w", outputPath, i, err))
					}

					defer outFile.Close()
					if fileType == "jpeg" {
						err = image_jpeg.Encode(outFile, renderedImage.(*image.RGBA), image_jpeg.Options{
							Options: &jpeg.Options{
								Quality: quality,
							},
							Progressive: progressive,
						})
						if err != nil {
							log.Fatal(fmt.Errorf("could not create output jpeg %s for tiff image %d: %w", outputPath, i, err))
						}
					} else if fileType == "png" {
						err = png.Encode(outFile, renderedImage)
						if err != nil {
							log.Fatal(fmt.Errorf("could not create output png %s for tiff image %d: %w", outputPath, i, err))
						}
					} else {
						log.Fatal(fmt.Errorf("invalid output filetype: %s", fileType))
					}

					log.Printf("Created file %s for TIFF image %d", outputPath, i)
				}()
			}
		},
	}

	rootCmd.Flags().IntVarP(&quality, "quality", "", 95, "The quality to render the image in, only used for jpeg.")
	rootCmd.Flags().StringVarP(&fileType, "file-type", "", "jpeg", "The file type to render in, jpeg or png")
	rootCmd.Flags().BoolVarP(&progressive, "progressive", "", false, "Create progressive images, only used for jpeg.")

	rootCmd.SetOut(os.Stdout)
	return rootCmd.Execute()
}
