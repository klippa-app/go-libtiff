package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/klippa-app/go-libtiff/libtiff"

	"github.com/tetratelabs/wazero"
)

func main() {
	// Please be aware that the input path needs to be absolute, and that the
	// path needs to be available inside the Wazero runtime.
	input := "/testdata/multipage-sample.tif"

	// The output path is handled in Go, so that needs to be handled as such,
	// so not a Wazero runtime path, and can be relative.
	output := "./testdata/multipage-sample.tif-%d.jpeg"

	ctx := context.Background()
	instance, err := libtiff.GetInstance(ctx, &libtiff.Config{
		FSConfig: wazero.NewFSConfig().
			WithReadOnlyDirMount("./testdata", "/testdata"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer instance.Close(ctx)

	file, err := instance.TIFFOpenFileFromPath(ctx, input, nil)
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
		_, err := file.ToImage(ctx, &libtiff.ImageOptions{
			OutputFormat:   libtiff.ImageOptionsOutputFormatJPEG,
			OutputTarget:   libtiff.ImageOptionsOutputTargetFile,
			OutputQuality:  95,
			Progressive:    false,
			TargetFilePath: strings.Replace(output, "%d", fmt.Sprintf("%d", i), 1),
		})
		if err != nil {
			log.Fatal(fmt.Errorf("could not convert tiff image %d to image: %w", i, err))
		}
	}
}
