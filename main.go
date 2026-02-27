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
	availableBinaries = append(availableBinaries, "tiff2img", "img2tiff")
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

	if os.Args[1] == "img2tiff" {
		err := img2tiff()
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

func img2tiff() error {
	var (
		compression    string
		quality        int
		append         bool
		software       string
		dateTime       string
		artist         string
		predictor      string
		xResolution    float32
		yResolution    float32
		resolutionUnit string
		description    string
		copyright      string
		documentName   string
		pageName       string
		hostComputer   string
		make_          string
		model          string
		rowsPerStrip   uint32
		orientation    string
		tileWidth      uint32
		tileHeight     uint32
		pageNumber     uint16
		totalPages     uint16
	)

	rootCmd := &cobra.Command{
		Use:   "img2tiff [input...] [output]",
		Short: "A CLI tool to convert JPEG/PNG images to TIFF",
		Args: func(cmd *cobra.Command, args []string) error {
			// args[0] is the command name, need at least 1 input + 1 output.
			return cobra.MinimumNArgs(3)(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			inputs := args[1 : len(args)-1]
			output := args[len(args)-1]

			// Register decoders.
			_ = jpeg.Decode
			_ = png.Decode

			// Map compression string to TIFFTAG.
			var comp libtiff.TIFFTAG
			switch strings.ToLower(compression) {
			case "none":
				comp = libtiff.COMPRESSION_NONE
			case "lzw":
				comp = libtiff.COMPRESSION_LZW
			case "deflate":
				comp = libtiff.COMPRESSION_ADOBE_DEFLATE
			case "jpeg":
				comp = libtiff.COMPRESSION_JPEG
			case "packbits":
				comp = libtiff.COMPRESSION_PACKBITS
			case "ccitt3":
				comp = libtiff.COMPRESSION_CCITTFAX3
			case "ccitt4":
				comp = libtiff.COMPRESSION_CCITTFAX4
			default:
				log.Fatal(fmt.Errorf("unsupported compression: %s (use none, lzw, deflate, jpeg, packbits, ccitt3, or ccitt4)", compression))
			}

			// Map predictor string.
			var pred libtiff.TIFFTAG
			switch strings.ToLower(predictor) {
			case "none", "":
				// Don't set predictor.
			case "horizontal":
				pred = libtiff.PREDICTOR_HORIZONTAL
			case "floatingpoint":
				pred = libtiff.PREDICTOR_FLOATINGPOINT
			default:
				log.Fatal(fmt.Errorf("unsupported predictor: %s (use none, horizontal, or floatingpoint)", predictor))
			}

			// Map resolution unit string.
			var resUnit libtiff.TIFFTAG
			switch strings.ToLower(resolutionUnit) {
			case "", "default":
				// Don't set resolution unit.
			case "none":
				resUnit = libtiff.RESUNIT_NONE
			case "inch":
				resUnit = libtiff.RESUNIT_INCH
			case "centimeter":
				resUnit = libtiff.RESUNIT_CENTIMETER
			default:
				log.Fatal(fmt.Errorf("unsupported resolution unit: %s (use none, inch, or centimeter)", resolutionUnit))
			}

			// Map orientation string.
			var orient libtiff.TIFFTAG
			switch strings.ToLower(orientation) {
			case "", "default":
				// Don't set orientation (use default TOPLEFT).
			case "topleft":
				orient = libtiff.ORIENTATION_TOPLEFT
			case "topright":
				orient = libtiff.ORIENTATION_TOPRIGHT
			case "botright":
				orient = libtiff.ORIENTATION_BOTRIGHT
			case "botleft":
				orient = libtiff.ORIENTATION_BOTLEFT
			case "lefttop":
				orient = libtiff.ORIENTATION_LEFTTOP
			case "righttop":
				orient = libtiff.ORIENTATION_RIGHTTOP
			case "rightbot":
				orient = libtiff.ORIENTATION_RIGHTBOT
			case "leftbot":
				orient = libtiff.ORIENTATION_LEFTBOT
			default:
				log.Fatal(fmt.Errorf("unsupported orientation: %s (use topleft, topright, botright, botleft, lefttop, righttop, rightbot, or leftbot)", orientation))
			}

			instance, err := libtiff.GetInstance(ctx, &libtiff.Config{
				CompilationCache: compilationCache,
			})
			if err != nil {
				log.Fatal(err)
			}
			defer instance.Close(ctx)

			// Open or create the output file.
			var outputFile *os.File
			var fileMode string
			if append {
				outputFile, err = os.OpenFile(output, os.O_RDWR, 0)
				if err != nil {
					log.Fatal(fmt.Errorf("could not open existing tiff file for appending: %w", err))
				}
				fileMode = "a"
			} else {
				outputFile, err = os.Create(output)
				if err != nil {
					log.Fatal(err)
				}
				fileMode = "w"
			}
			defer outputFile.Close()
			tiffFile, err := instance.TIFFOpenFileFromReadWriteSeeker(ctx, path.Base(output), outputFile, 0, &libtiff.OpenOptions{
				FileMode: &fileMode,
			})
			if err != nil {
				log.Fatal(fmt.Errorf("could not open tiff file for writing: %w", err))
			}
			defer tiffFile.Close(ctx)

			for i, input := range inputs {
				// Open and decode the input image.
				inputFile, err := os.Open(input)
				if err != nil {
					log.Fatal(err)
				}

				img, _, err := image.Decode(inputFile)
				inputFile.Close()
				if err != nil {
					log.Fatal(fmt.Errorf("could not decode input image %s: %w", input, err))
				}

				err = tiffFile.FromGoImage(ctx, img, &libtiff.FromGoImageOptions{
					Compression:    comp,
					Quality:        quality,
					Software:       software,
					DateTime:       dateTime,
					Artist:         artist,
					Predictor:      pred,
					XResolution:    xResolution,
					YResolution:    yResolution,
					ResolutionUnit: resUnit,
					Description:    description,
					Copyright:      copyright,
					DocumentName:   documentName,
					PageName:       pageName,
					HostComputer:   hostComputer,
					Make:           make_,
					Model:          model,
					RowsPerStrip:   rowsPerStrip,
					Orientation:    orient,
					TileWidth:      tileWidth,
					TileHeight:     tileHeight,
					PageNumber:     pageNumber,
					TotalPages:     totalPages,
				})
				if err != nil {
					log.Fatal(fmt.Errorf("could not write image %s to tiff: %w", input, err))
				}

				log.Printf("Written image %d/%d: %s", i+1, len(inputs), input)
			}

			if append {
				log.Printf("Appended %d image(s) to TIFF file %s", len(inputs), output)
			} else {
				log.Printf("Created TIFF file %s with %d image(s)", output, len(inputs))
			}
		},
	}

	rootCmd.Flags().StringVarP(&compression, "compression", "", "deflate", "Compression type: none, lzw, deflate, jpeg, packbits, ccitt3, or ccitt4")
	rootCmd.Flags().IntVarP(&quality, "quality", "", 75, "JPEG compression quality (1-100), only used with --compression jpeg")
	rootCmd.Flags().BoolVarP(&append, "append", "", false, "Append to an existing TIFF file instead of creating a new one")
	rootCmd.Flags().StringVarP(&software, "software", "", "", "TIFFTAG_SOFTWARE value (default: go-libtiff/libtiff-{version})")
	rootCmd.Flags().StringVarP(&dateTime, "datetime", "", "", "TIFFTAG_DATETIME value in YYYY:MM:DD HH:MM:SS format (default: current time)")
	rootCmd.Flags().StringVarP(&artist, "artist", "", "", "TIFFTAG_ARTIST value (omitted if empty)")
	rootCmd.Flags().StringVarP(&predictor, "predictor", "", "none", "Predictor: none, horizontal, or floatingpoint")
	rootCmd.Flags().Float32VarP(&xResolution, "xresolution", "", 0, "X resolution (DPI)")
	rootCmd.Flags().Float32VarP(&yResolution, "yresolution", "", 0, "Y resolution (DPI)")
	rootCmd.Flags().StringVarP(&resolutionUnit, "resolution-unit", "", "", "Resolution unit: none, inch, or centimeter")
	rootCmd.Flags().StringVarP(&description, "description", "", "", "TIFFTAG_IMAGEDESCRIPTION value (omitted if empty)")
	rootCmd.Flags().StringVarP(&copyright, "copyright", "", "", "TIFFTAG_COPYRIGHT value (omitted if empty)")
	rootCmd.Flags().StringVarP(&documentName, "document-name", "", "", "TIFFTAG_DOCUMENTNAME value (omitted if empty)")
	rootCmd.Flags().StringVarP(&pageName, "page-name", "", "", "TIFFTAG_PAGENAME value (omitted if empty)")
	rootCmd.Flags().StringVarP(&hostComputer, "host-computer", "", "", "TIFFTAG_HOSTCOMPUTER value (omitted if empty)")
	rootCmd.Flags().StringVarP(&make_, "make", "", "", "TIFFTAG_MAKE value (omitted if empty)")
	rootCmd.Flags().StringVarP(&model, "model", "", "", "TIFFTAG_MODEL value (omitted if empty)")
	rootCmd.Flags().Uint32VarP(&rowsPerStrip, "rows-per-strip", "", 0, "Rows per strip (0 = auto)")
	rootCmd.Flags().StringVarP(&orientation, "orientation", "", "", "Orientation: topleft, topright, botright, botleft, lefttop, righttop, rightbot, or leftbot")
	rootCmd.Flags().Uint32VarP(&tileWidth, "tile-width", "", 0, "Tile width (0 = strip-based)")
	rootCmd.Flags().Uint32VarP(&tileHeight, "tile-height", "", 0, "Tile height (0 = strip-based)")
	rootCmd.Flags().Uint16VarP(&pageNumber, "page-number", "", 0, "Page number (0-based) for TIFFTAG_PAGENUMBER")
	rootCmd.Flags().Uint16VarP(&totalPages, "total-pages", "", 0, "Total pages for TIFFTAG_PAGENUMBER (tag is omitted if 0)")

	rootCmd.SetOut(os.Stdout)
	return rootCmd.Execute()
}
