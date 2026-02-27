package libtiff

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

// TIFFReadEncodedStrip reads and decompresses a strip of data.
// Returns the decompressed strip bytes.
func (f *File) TIFFReadEncodedStrip(ctx context.Context, strip uint32) ([]byte, error) {
	stripSize, err := f.TIFFStripSize(ctx)
	if err != nil {
		return nil, err
	}

	bufPointer, err := f.instance.malloc(ctx, uint64(stripSize))
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadEncodedStrip", f.pointer, api.EncodeU32(strip), bufPointer, api.EncodeI32(int32(-1)))
	if err != nil {
		return nil, err
	}

	bytesRead := api.DecodeI32(results[0])
	if bytesRead == -1 {
		return nil, errors.New("error reading encoded strip")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bytesRead))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read strip data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}

// TIFFReadEncodedTile reads and decompresses a tile of data.
// Returns the decompressed tile bytes.
func (f *File) TIFFReadEncodedTile(ctx context.Context, tile uint32) ([]byte, error) {
	tileSize, err := f.TIFFTileSize(ctx)
	if err != nil {
		return nil, err
	}

	bufPointer, err := f.instance.malloc(ctx, uint64(tileSize))
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadEncodedTile", f.pointer, api.EncodeU32(tile), bufPointer, api.EncodeI32(int32(-1)))
	if err != nil {
		return nil, err
	}

	bytesRead := api.DecodeI32(results[0])
	if bytesRead == -1 {
		return nil, errors.New("error reading encoded tile")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bytesRead))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read tile data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}

// TIFFReadScanline reads a single scanline of data at the given row.
// For planar images, sample specifies the plane (0-based); use 0 for contiguous images.
func (f *File) TIFFReadScanline(ctx context.Context, row uint32, sample uint16) ([]byte, error) {
	scanlineSize, err := f.TIFFScanlineSize(ctx)
	if err != nil {
		return nil, err
	}

	bufPointer, err := f.instance.malloc(ctx, uint64(scanlineSize))
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadScanline", f.pointer, bufPointer, api.EncodeU32(row), api.EncodeU32(uint32(sample)))
	if err != nil {
		return nil, err
	}

	if api.DecodeI32(results[0]) == -1 {
		return nil, errors.New("error reading scanline")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(scanlineSize))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read scanline data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}

// TIFFReadRGBAStrip reads a strip of data and converts it to RGBA format.
// row is the first row in the strip. Returns RGBA pixels (4 bytes per pixel: R, G, B, A).
func (f *File) TIFFReadRGBAStrip(ctx context.Context, row uint32) ([]byte, error) {
	width, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_IMAGEWIDTH)
	if err != nil {
		return nil, err
	}

	rowsPerStrip, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_ROWSPERSTRIP)
	if err != nil {
		return nil, err
	}

	height, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_IMAGELENGTH)
	if err != nil {
		return nil, err
	}

	// The actual number of rows in this strip may be less than rowsPerStrip
	// for the last strip.
	rows := rowsPerStrip
	if row+rows > height {
		rows = height - row
	}

	bufSize := uint64(width) * uint64(rows) * 4
	bufPointer, err := f.instance.malloc(ctx, bufSize)
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadRGBAStrip", f.pointer, api.EncodeU32(row), bufPointer)
	if err != nil {
		return nil, err
	}

	if results[0] == 0 {
		return nil, errors.New("error reading RGBA strip")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bufSize))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read RGBA strip data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}

// TIFFReadRGBATile reads a tile of data and converts it to RGBA format.
// col and row are the origin of the tile. Returns RGBA pixels (4 bytes per pixel: R, G, B, A).
func (f *File) TIFFReadRGBATile(ctx context.Context, col, row uint32) ([]byte, error) {
	tileWidth, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_TILEWIDTH)
	if err != nil {
		return nil, err
	}

	tileLength, err := f.TIFFGetFieldUint32_t(ctx, TIFFTAG_TILELENGTH)
	if err != nil {
		return nil, err
	}

	bufSize := uint64(tileWidth) * uint64(tileLength) * 4
	bufPointer, err := f.instance.malloc(ctx, bufSize)
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadRGBATile", f.pointer, api.EncodeU32(col), api.EncodeU32(row), bufPointer)
	if err != nil {
		return nil, err
	}

	if results[0] == 0 {
		return nil, errors.New("error reading RGBA tile")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bufSize))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read RGBA tile data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}

// TIFFReadRawStrip reads the raw (compressed) data for a strip.
func (f *File) TIFFReadRawStrip(ctx context.Context, strip uint32) ([]byte, error) {
	// Use decompressed strip size as a buffer upper bound.
	stripSize, err := f.TIFFStripSize(ctx)
	if err != nil {
		return nil, err
	}

	bufPointer, err := f.instance.malloc(ctx, uint64(stripSize))
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadRawStrip", f.pointer, api.EncodeU32(strip), bufPointer, api.EncodeI32(int32(stripSize)))
	if err != nil {
		return nil, err
	}

	bytesRead := api.DecodeI32(results[0])
	if bytesRead == -1 {
		return nil, errors.New("error reading raw strip")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bytesRead))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read raw strip data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}

// TIFFReadRawTile reads the raw (compressed) data for a tile.
func (f *File) TIFFReadRawTile(ctx context.Context, tile uint32) ([]byte, error) {
	// Use decompressed tile size as a buffer upper bound.
	tileSize, err := f.TIFFTileSize(ctx)
	if err != nil {
		return nil, err
	}

	bufPointer, err := f.instance.malloc(ctx, uint64(tileSize))
	if err != nil {
		return nil, err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFReadRawTile", f.pointer, api.EncodeU32(tile), bufPointer, api.EncodeI32(int32(tileSize)))
	if err != nil {
		return nil, err
	}

	bytesRead := api.DecodeI32(results[0])
	if bytesRead == -1 {
		return nil, errors.New("error reading raw tile")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bytesRead))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return nil, errors.New("could not read raw tile data from WASM memory")
	}
	data := make([]byte, len(buf))
	copy(data, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return data, nil
}
