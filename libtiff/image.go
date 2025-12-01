package libtiff

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/klippa-app/go-libtiff/internal/image/image_jpeg"
	"github.com/tetratelabs/wazero/api"
)

// ToGoImage convert the current directory in the open TIFF file to RGBA, the
// caller is responsible for closing since the returned cleanup function will
// free the allocated memory.
func (f *File) ToGoImage(ctx context.Context) (image.Image, func(context.Context) error, error) {
	width, height, err := f.GetDimensions(ctx)
	if err != nil {
		return nil, nil, err
	}

	img := &image.RGBA{}
	img.Rect = image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: width, Y: height}}
	img.Stride = img.Rect.Max.X * 4
	nBytes := img.Rect.Max.X * img.Rect.Max.Y * 4
	imagePointer, err := f.instance.Malloc(ctx, uint64(nBytes))
	if err != nil {
		return nil, nil, err
	}

	cleanupFunc := func(ctx context.Context) error {
		return f.instance.Free(ctx, imagePointer)
	}

	results, err := f.instance.internalInstance.Module.ExportedFunction("TIFFReadRGBAImageOriented").Call(ctx, f.pointer, api.EncodeU32(uint32(img.Rect.Max.X)), api.EncodeU32(uint32(img.Rect.Max.Y)), imagePointer, api.EncodeU32(uint32(ORIENTATION_TOPLEFT)), 0)
	if err != nil {
		cleanupErr := cleanupFunc(ctx)
		return nil, nil, errors.Join(err, cleanupErr)
	}

	if results[0] != 1 {
		cleanupErr := cleanupFunc(ctx)
		return nil, nil, errors.Join(errors.New("error while converting tiff to RGBA"), cleanupErr)
	}

	// We directly open a view on the image data in the Wazero memory so that
	// we don't have to do any image copying.
	memoryView, ok := f.instance.internalInstance.Module.Memory().Read(uint32(imagePointer), uint32(nBytes))
	if !ok {
		cleanupErr := cleanupFunc(ctx)
		return nil, nil, errors.Join(errors.New("memory view not found"), cleanupErr)
	}

	img.Pix = memoryView

	return img, cleanupFunc, nil
}

type ImageOptionsOutputFormat string // The file format to render output as.

const (
	ImageOptionsOutputFormatJPEG ImageOptionsOutputFormat = "jpg" // Render the file as a JPEG file.
	ImageOptionsOutputFormatPNG  ImageOptionsOutputFormat = "png" // Render the file as a PNG file.
)

type ImageOptionsOutputTarget string // The file target output.

const (
	ImageOptionsOutputTargetBytes ImageOptionsOutputTarget = "bytes" // Returns the file as a byte array in the response.
	ImageOptionsOutputTargetFile  ImageOptionsOutputTarget = "file"  // Writes away the file to a given path or a generated tmp file.
)

type ImageOptions struct {
	OutputFormat   ImageOptionsOutputFormat // The format to output the image as
	OutputTarget   ImageOptionsOutputTarget // Where to output the image
	OutputQuality  int                      // Only used when OutputFormat RenderToFileOutputFormatJPG. Ranges from 1 to 100 inclusive, higher is better. The default is 95.
	Progressive    bool                     // Only used when OutputFormat RenderToFileOutputFormatJPG and with build tag libtiff_use_turbojpeg. Will render a progressive jpeg.
	MaxFileSize    int64                    // The maximum file size, when OutputFormat RenderToFileOutputFormatJPG, it will try to lower the quality it until it fits.
	TargetFilePath string                   // When OutputTarget is file, the path to write it to.
}

// ToImage convert the current directory in the open TIFF file to an image file.
func (f *File) ToImage(ctx context.Context, options *ImageOptions) ([]byte, error) {
	if options == nil {
		return nil, errors.New("options cannot be nil")
	}

	renderedImage, cleanup, err := f.ToGoImage(ctx)
	if err != nil {
		return nil, err
	}
	defer cleanup(ctx)

	var imgBuf bytes.Buffer

	if options.OutputFormat == ImageOptionsOutputFormatJPEG {
		opt := image_jpeg.Options{
			Options: &jpeg.Options{
				Quality: 95,
			},
			Progressive: options.Progressive,
		}

		if options.OutputQuality > 0 {
			opt.Options.Quality = options.OutputQuality
		}

		for {
			err := image_jpeg.Encode(&imgBuf, renderedImage.(*image.RGBA), opt)
			if err != nil {
				return nil, err
			}

			if options.MaxFileSize == 0 || int64(imgBuf.Len()) < options.MaxFileSize {
				break
			}

			opt.Quality -= 10

			if opt.Quality <= 45 {
				return nil, errors.New("TIFF image would exceed maximum filesize")
			}

			imgBuf.Reset()
		}
	} else if options.OutputFormat == ImageOptionsOutputFormatPNG {
		err := png.Encode(&imgBuf, renderedImage)
		if err != nil {
			return nil, err
		}

		if options.MaxFileSize != 0 && int64(imgBuf.Len()) > options.MaxFileSize {
			return nil, errors.New("TIFF image would exceed maximum filesize")
		}
	} else {
		return nil, errors.New("invalid output format given")
	}

	if options.OutputTarget == ImageOptionsOutputTargetBytes {
		imageBytes := imgBuf.Bytes()
		return imageBytes, nil
	} else if options.OutputTarget == ImageOptionsOutputTargetFile {
		var targetFile *os.File
		if options.TargetFilePath == "" {
			return nil, errors.New("target file path can't be empty")
		}

		if options.TargetFilePath != "" {
			existingFile, err := os.Create(options.TargetFilePath)
			if err != nil {
				return nil, err
			}
			targetFile = existingFile
		}

		_, err := targetFile.Write(imgBuf.Bytes())
		if err != nil {
			return nil, err
		}

		err = targetFile.Close()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid output target given")
	}

	return nil, nil
}
