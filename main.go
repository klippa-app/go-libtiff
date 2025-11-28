package main

import "C"
import (
	"context"
	_ "embed"
	"log"
	"os"
	"slices"
	"strings"

	_ "github.com/klippa-app/go-libtiff/fax2ps"
	_ "github.com/klippa-app/go-libtiff/fax2tiff"
	"github.com/klippa-app/go-libtiff/internal/registry"
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
)

func main() {
	ctx := context.Background()

	availableBinaries := registry.List()
	incorrectStartArgument := func() {
		log.Fatalf("You should minimally start the program with one of the following arguments: %s", strings.Join(availableBinaries, ", "))
	}
	if len(os.Args) < 2 {
		incorrectStartArgument()
	}
	if !slices.Contains(availableBinaries, os.Args[1]) {
		incorrectStartArgument()
	}

	err := registry.Run(ctx, os.Args[1], os.Args[1:]...)
	if err != nil {
		log.Fatal(err)
	}
}
