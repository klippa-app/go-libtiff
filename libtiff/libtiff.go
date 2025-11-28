package libtiff

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/klippa-app/go-libtiff/internal/instance"
	"github.com/tetratelabs/wazero/api"
)

//go:embed libtiff.wasm
var wasmBinary []byte

func Convert(ctx context.Context, filePath string) error {
	wazeroInstance, err := instance.GetInstance(ctx, &instance.Config{
		WASMData: wasmBinary,
	})
	if err != nil {
		return err
	}

	defer wazeroInstance.Close(ctx)

	cStringFilePath, err := NewCString(ctx, wazeroInstance.Module, filePath)
	if err != nil {
		return err
	}
	defer cStringFilePath.Free(ctx)

	cStringFileMode, err := NewCString(ctx, wazeroInstance.Module, "r")
	if err != nil {
		return err
	}
	defer cStringFileMode.Free(ctx)

	tiffStart := time.Now()

	// Result is a pointer to struct_tiff
	res, err := wazeroInstance.Module.ExportedFunction("TIFFOpen").Call(ctx, cStringFilePath.Pointer, cStringFileMode.Pointer)
	if err != nil {
		return err
	}
	if res[0] == 0 {
		return errors.New("Error while opening tiff file")
	}

	tiffFilePointer := res[0]
	defer func() {
		_, _ = wazeroInstance.Module.ExportedFunction("TIFFClose").Call(ctx, tiffFilePointer)
	}()

	width, height, err := getWH(ctx, wazeroInstance.Module, tiffFilePointer)
	if err != nil {
		return err
	}
	log.Printf("Size=Width: %d, Height: %d", width, height)

	/*
		x, y, err := getRes(ctx, mod, tiffFilePointer)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("DPI=X: %f, Y: %f", x, y)

	*/

	log.Printf("Time to read tiff: %dms", time.Since(tiffStart).Milliseconds())
	tiffStart = time.Now()

	img := &image.RGBA{}
	img.Rect = image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{int(width), int(height)}}
	img.Stride = img.Rect.Max.X * 4
	nBytes := img.Rect.Max.X * img.Rect.Max.Y * 4
	imagePointer, err := Malloc(ctx, wazeroInstance.Module, uint64(nBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer Free(ctx, wazeroInstance.Module, imagePointer)

	results, err := wazeroInstance.Module.ExportedFunction("TIFFReadRGBAImageOriented").Call(ctx, tiffFilePointer, api.EncodeU32(uint32(img.Rect.Max.X)), api.EncodeU32(uint32(img.Rect.Max.Y)), imagePointer, api.EncodeU32(uint32(ORIENTATION_TOPLEFT)), 0)
	if err != nil {
		log.Fatal(err)
	}

	if results[0] != 1 {
		log.Fatal("Error while reading tiff file")
	}

	memoryView, ok := wazeroInstance.Module.Memory().Read(uint32(imagePointer), uint32(nBytes))
	if !ok {
		log.Fatal(fmt.Errorf("memory view not found"))
	}
	img.Pix = memoryView

	log.Printf("Time to convert to RGBA: %dms", time.Since(tiffStart).Milliseconds())

	pngTime := time.Now()

	f, err := os.Create(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Time to convert RGBA to PNG: %dms", time.Since(pngTime).Milliseconds())

	return nil
}

type CString struct {
	Pointer uint64
	Free    func(ctx context.Context)
}

func NewCString(ctx context.Context, mod api.Module, input string) (*CString, error) {
	inputLength := uint64(len(input)) + 1

	pointer, err := Malloc(ctx, mod, inputLength)
	if err != nil {
		return nil, err
	}

	// Write string + null terminator.
	if !mod.Memory().Write(uint32(pointer), append([]byte(input), byte(0))) {
		return nil, errors.New("could not write CString data")
	}

	return &CString{
		Pointer: pointer,
		Free: func(ctx context.Context) {
			Free(ctx, mod, pointer)
		},
	}, nil
}

type TIFFTAG uint32

// https://gitlab.com/libtiff/libtiff/-/blob/master/libtiff/tiff.h
var (
	ORIENTATION_TOPLEFT   = TIFFTAG(1)
	TIFFTAG_IMAGEWIDTH    = TIFFTAG(256)
	TIFFTAG_IMAGELENGTH   = TIFFTAG(257)
	TIFFTAG_BITSPERSAMPLE = TIFFTAG(258)
	TIFFTAG_COMPRESSION   = TIFFTAG(259)
	TIFFTAG_XRESOLUTION   = TIFFTAG(282)
	TIFFTAG_YRESOLUTION   = TIFFTAG(283)
)

func getWH(ctx context.Context, mod api.Module, filePointer uint64) (int, int, error) {
	widthPointer, err := Malloc(ctx, mod, 4)
	if err != nil {
		return 0, 0, err
	}
	defer Free(ctx, mod, widthPointer)

	lengthPointer, err := Malloc(ctx, mod, 4)
	if err != nil {
		return 0, 0, err
	}
	defer Free(ctx, mod, lengthPointer)

	if err := TIFFGetFieldUint32_t(ctx, mod, filePointer, TIFFTAG_IMAGEWIDTH, widthPointer); err != nil {
		return 0, 0, err
	}

	if err := TIFFGetFieldUint32_t(ctx, mod, filePointer, TIFFTAG_IMAGELENGTH, lengthPointer); err != nil {
		return 0, 0, err
	}

	// Expected data type of TIFFTAG_IMAGEWIDTH and TIFFTAG_IMAGELENGTH is uint32_t.

	readWidth, success := mod.Memory().ReadUint32Le(uint32(widthPointer))
	if !success {
		return 0, 0, errors.New("Could not read width")
	}

	readLength, success := mod.Memory().ReadUint32Le(uint32(lengthPointer))
	if !success {
		return 0, 0, errors.New("Could not read length")
	}

	return int(readWidth), int(readLength), nil
}

func getRes(ctx context.Context, mod api.Module, filePointer uint64) (float32, float32, error) {
	xPointer, err := Malloc(ctx, mod, 4)
	if err != nil {
		return 0, 0, err
	}
	defer Free(ctx, mod, xPointer)

	yPointer, err := Malloc(ctx, mod, 4)
	if err != nil {
		return 0, 0, err
	}
	defer Free(ctx, mod, yPointer)

	if err := TIFFGetFieldFloat(ctx, mod, filePointer, TIFFTAG_XRESOLUTION, xPointer); err != nil {
		return 0, 0, err
	}

	if err := TIFFGetFieldFloat(ctx, mod, filePointer, TIFFTAG_YRESOLUTION, yPointer); err != nil {
		return 0, 0, err
	}

	// Expected type of TIFFTAG_XRESOLUTION / TIFFTAG_YRESOLUTION is float.

	readX, success := mod.Memory().ReadFloat32Le(uint32(xPointer))
	if !success {
		return 0, 0, errors.New("Could not read x")
	}

	readY, success := mod.Memory().ReadFloat32Le(uint32(yPointer))
	if !success {
		return 0, 0, errors.New("Could not read y")
	}

	return readX, readY, nil
}

func TIFFGetFieldUint32_t(ctx context.Context, mod api.Module, filePointer uint64, tag TIFFTAG, valuePointer uint64) error {
	results, err := mod.ExportedFunction("TIFFGetFieldUint32_t").Call(ctx, filePointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return fmt.Errorf("could not get field data for tag %d", tag)
	}

	return nil
}

func TIFFGetFieldFloat(ctx context.Context, mod api.Module, filePointer uint64, tag TIFFTAG, valuePointer uint64) error {
	results, err := mod.ExportedFunction("TIFFGetFieldFloat").Call(ctx, filePointer, api.EncodeU32(uint32(tag)), valuePointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return fmt.Errorf("could not get field data for tag %d", tag)
	}

	return nil
}

func Malloc(ctx context.Context, mod api.Module, size uint64) (uint64, error) {
	results, err := mod.ExportedFunction("malloc").Call(ctx, size)
	if err != nil {
		return 0, err
	}

	pointer := results[0]
	ok := mod.Memory().Write(uint32(results[0]), make([]byte, size))
	if !ok {
		return 0, errors.New("could not write nulls to memory")
	}

	return pointer, nil
}

func Free(ctx context.Context, mod api.Module, pointer uint64) error {
	_, err := mod.ExportedFunction("free").Call(ctx, pointer)
	if err != nil {
		return err
	}
	return nil
}
