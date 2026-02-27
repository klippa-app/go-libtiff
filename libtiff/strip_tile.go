package libtiff

import (
	"context"

	"github.com/tetratelabs/wazero/api"
)

// TIFFStripSize returns the size in bytes of a decompressed strip.
func (f *File) TIFFStripSize(ctx context.Context) (int64, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFStripSize", f.pointer)
	if err != nil {
		return 0, err
	}

	return int64(api.DecodeI32(results[0])), nil
}

// TIFFNumberOfStrips returns the total number of strips in the image.
func (f *File) TIFFNumberOfStrips(ctx context.Context) (uint32, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFNumberOfStrips", f.pointer)
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// TIFFTileSize returns the size in bytes of a decompressed tile.
func (f *File) TIFFTileSize(ctx context.Context) (int64, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFTileSize", f.pointer)
	if err != nil {
		return 0, err
	}

	return int64(api.DecodeI32(results[0])), nil
}

// TIFFNumberOfTiles returns the total number of tiles in the image.
func (f *File) TIFFNumberOfTiles(ctx context.Context) (uint32, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFNumberOfTiles", f.pointer)
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// TIFFComputeStrip returns the strip number containing the given row and sample.
func (f *File) TIFFComputeStrip(ctx context.Context, row uint32, sample uint16) (uint32, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFComputeStrip", f.pointer, api.EncodeU32(row), api.EncodeU32(uint32(sample)))
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// TIFFComputeTile returns the tile number containing the given coordinates and sample.
func (f *File) TIFFComputeTile(ctx context.Context, x, y, z uint32, sample uint16) (uint32, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFComputeTile", f.pointer, api.EncodeU32(x), api.EncodeU32(y), api.EncodeU32(z), api.EncodeU32(uint32(sample)))
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// TIFFIsTiled returns true if the image data is organized in tiles rather than strips.
func (f *File) TIFFIsTiled(ctx context.Context) (bool, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFIsTiled", f.pointer)
	if err != nil {
		return false, err
	}

	return results[0] != 0, nil
}

// TIFFScanlineSize returns the size in bytes of a decompressed scanline.
func (f *File) TIFFScanlineSize(ctx context.Context) (int64, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFScanlineSize", f.pointer)
	if err != nil {
		return 0, err
	}

	return int64(api.DecodeI32(results[0])), nil
}

// TIFFVStripSize returns the size in bytes of a strip with the given number of rows.
func (f *File) TIFFVStripSize(ctx context.Context, nrows uint32) (int64, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFVStripSize", f.pointer, api.EncodeU32(nrows))
	if err != nil {
		return 0, err
	}

	return int64(api.DecodeI32(results[0])), nil
}

// TIFFDefaultTileSize returns the default tile dimensions for the file.
func (f *File) TIFFDefaultTileSize(ctx context.Context) (tileWidth uint32, tileHeight uint32, err error) {
	twPointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, 0, err
	}
	defer f.instance.free(ctx, twPointer)

	thPointer, err := f.instance.malloc(ctx, 4)
	if err != nil {
		return 0, 0, err
	}
	defer f.instance.free(ctx, thPointer)

	// Write initial values of 0.
	f.instance.internalInstance.CallLock.Lock()
	f.instance.internalInstance.Module.Memory().WriteUint32Le(uint32(twPointer), 0)
	f.instance.internalInstance.Module.Memory().WriteUint32Le(uint32(thPointer), 0)
	f.instance.internalInstance.CallLock.Unlock()

	_, err = f.instance.internalInstance.CallExportedFunction(ctx, "TIFFDefaultTileSize", f.pointer, twPointer, thPointer)
	if err != nil {
		return 0, 0, err
	}

	f.instance.internalInstance.CallLock.Lock()
	tw, _ := f.instance.internalInstance.Module.Memory().ReadUint32Le(uint32(twPointer))
	th, _ := f.instance.internalInstance.Module.Memory().ReadUint32Le(uint32(thPointer))
	f.instance.internalInstance.CallLock.Unlock()

	return tw, th, nil
}
