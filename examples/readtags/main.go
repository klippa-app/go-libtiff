package main

import (
	"context"
	"errors"
	"log"

	"github.com/klippa-app/go-libtiff/libtiff"

	"github.com/tetratelabs/wazero"
)

func main() {
	// Please be aware that the input path needs to be absolute, and that the
	// path needs to be available inside the Wazero runtime.
	input := "/testdata/multipage-sample.tif"

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

	getTag := func(tag libtiff.TIFFTAG) (string, error) {
		value, err := file.TIFFGetFieldConstChar(ctx, tag)
		if err != nil {
			// This error will be returned if the tag is not defined or of the wrong type.
			if !errors.Is(err, &libtiff.TagNotDefinedError{}) {
				return "", err
			}
			return "", nil
		}
		return value, nil
	}

	software, err := getTag(libtiff.TIFFTAG_SOFTWARE)
	if err != nil {
		log.Fatal(err)
	}

	dateTime, err := getTag(libtiff.TIFFTAG_DATETIME)
	if err != nil {
		log.Fatal(err)
	}

	artist, err := getTag(libtiff.TIFFTAG_ARTIST)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Software: %s\nDatetime: %s\nArtist: %s\n", software, dateTime, artist)
}
