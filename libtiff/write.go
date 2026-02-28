package libtiff

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/tetratelabs/wazero/api"
)

// TIFFWriteEncodedStrip writes a strip of data to the TIFF file.
func (f *File) TIFFWriteEncodedStrip(ctx context.Context, strip uint32, data []byte) error {
	size := uint64(len(data))
	dataPointer, err := f.instance.malloc(ctx, size)
	if err != nil {
		return err
	}
	defer f.instance.free(ctx, dataPointer)

	// Write data into WASM memory.
	f.instance.internalInstance.CallLock.Lock()
	ok := f.instance.internalInstance.Module.Memory().Write(uint32(dataPointer), data)
	f.instance.internalInstance.CallLock.Unlock()
	if !ok {
		return errors.New("could not write strip data to WASM memory")
	}

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFWriteEncodedStrip", f.pointer, api.EncodeU32(strip), dataPointer, size)
	if err != nil {
		return err
	}

	err = f.GetError()
	if err != nil {
		return err
	}

	// TIFFWriteEncodedStrip returns -1 on error.
	if api.DecodeI32(results[0]) == -1 {
		return errors.New("error writing encoded strip")
	}

	return nil
}

// TIFFWriteDirectory writes the current directory to the TIFF file.
func (f *File) TIFFWriteDirectory(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFWriteDirectory", f.pointer)
	if err != nil {
		return err
	}

	err = f.GetError()
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not write directory")
	}

	return nil
}

// TIFFDefaultStripSize returns a sensible default strip size for the given request.
func (f *File) TIFFDefaultStripSize(ctx context.Context, request uint32) (uint32, error) {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFDefaultStripSize", f.pointer, api.EncodeU32(request))
	if err != nil {
		return 0, err
	}

	return api.DecodeU32(results[0]), nil
}

// TIFFSetFieldExtraSamples sets the EXTRASAMPLES tag with the given sample types.
func (f *File) TIFFSetFieldExtraSamples(ctx context.Context, sampleTypes []uint16) error {
	count := uint16(len(sampleTypes))
	size := uint64(count) * 2

	arrayPointer, err := f.instance.malloc(ctx, size)
	if err != nil {
		return err
	}
	defer f.instance.free(ctx, arrayPointer)

	// Write the uint16 array into WASM memory.
	f.instance.internalInstance.CallLock.Lock()
	buf := make([]byte, size)
	for i, v := range sampleTypes {
		binary.LittleEndian.PutUint16(buf[i*2:], v)
	}
	ok := f.instance.internalInstance.Module.Memory().Write(uint32(arrayPointer), buf)
	f.instance.internalInstance.CallLock.Unlock()
	if !ok {
		return errors.New("could not write extra samples array to WASM memory")
	}

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFSetFieldExtraSamples", f.pointer, api.EncodeU32(uint32(count)), arrayPointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not set extra samples")
	}

	return nil
}

// TIFFWriteEncodedTile writes a tile of data to the TIFF file.
func (f *File) TIFFWriteEncodedTile(ctx context.Context, tile uint32, data []byte) error {
	size := uint64(len(data))
	dataPointer, err := f.instance.malloc(ctx, size)
	if err != nil {
		return err
	}
	defer f.instance.free(ctx, dataPointer)

	f.instance.internalInstance.CallLock.Lock()
	ok := f.instance.internalInstance.Module.Memory().Write(uint32(dataPointer), data)
	f.instance.internalInstance.CallLock.Unlock()
	if !ok {
		return errors.New("could not write tile data to WASM memory")
	}

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFWriteEncodedTile", f.pointer, api.EncodeU32(tile), dataPointer, size)
	if err != nil {
		return err
	}

	err = f.GetError()
	if err != nil {
		return err
	}

	if api.DecodeI32(results[0]) == -1 {
		return errors.New("error writing encoded tile")
	}

	return nil
}

// TIFFWriteScanline writes a scanline of data at the given row.
// For planar images, sample specifies the plane (0-based); use 0 for contiguous images.
func (f *File) TIFFWriteScanline(ctx context.Context, data []byte, row uint32, sample uint16) error {
	size := uint64(len(data))
	dataPointer, err := f.instance.malloc(ctx, size)
	if err != nil {
		return err
	}
	defer f.instance.free(ctx, dataPointer)

	f.instance.internalInstance.CallLock.Lock()
	ok := f.instance.internalInstance.Module.Memory().Write(uint32(dataPointer), data)
	f.instance.internalInstance.CallLock.Unlock()
	if !ok {
		return errors.New("could not write scanline data to WASM memory")
	}

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFWriteScanline", f.pointer, dataPointer, api.EncodeU32(row), api.EncodeU32(uint32(sample)))
	if err != nil {
		return err
	}

	err = f.GetError()
	if err != nil {
		return err
	}

	if api.DecodeI32(results[0]) == -1 {
		return errors.New("error writing scanline")
	}

	return nil
}

// TIFFWriteRawStrip writes raw (pre-compressed) data for a strip.
func (f *File) TIFFWriteRawStrip(ctx context.Context, strip uint32, data []byte) error {
	size := uint64(len(data))
	dataPointer, err := f.instance.malloc(ctx, size)
	if err != nil {
		return err
	}
	defer f.instance.free(ctx, dataPointer)

	f.instance.internalInstance.CallLock.Lock()
	ok := f.instance.internalInstance.Module.Memory().Write(uint32(dataPointer), data)
	f.instance.internalInstance.CallLock.Unlock()
	if !ok {
		return errors.New("could not write raw strip data to WASM memory")
	}

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFWriteRawStrip", f.pointer, api.EncodeU32(strip), dataPointer, size)
	if err != nil {
		return err
	}

	err = f.GetError()
	if err != nil {
		return err
	}

	if api.DecodeI32(results[0]) == -1 {
		return errors.New("error writing raw strip")
	}

	return nil
}

// TIFFWriteRawTile writes raw (pre-compressed) data for a tile.
func (f *File) TIFFWriteRawTile(ctx context.Context, tile uint32, data []byte) error {
	size := uint64(len(data))
	dataPointer, err := f.instance.malloc(ctx, size)
	if err != nil {
		return err
	}
	defer f.instance.free(ctx, dataPointer)

	f.instance.internalInstance.CallLock.Lock()
	ok := f.instance.internalInstance.Module.Memory().Write(uint32(dataPointer), data)
	f.instance.internalInstance.CallLock.Unlock()
	if !ok {
		return errors.New("could not write raw tile data to WASM memory")
	}

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFWriteRawTile", f.pointer, api.EncodeU32(tile), dataPointer, size)
	if err != nil {
		return err
	}

	err = f.GetError()
	if err != nil {
		return err
	}

	if api.DecodeI32(results[0]) == -1 {
		return errors.New("error writing raw tile")
	}

	return nil
}

// TIFFDeferStrileArrayWriting defers writing of strip/tile offset and byte count
// arrays. This can improve write performance for large files by batching metadata updates.
func (f *File) TIFFDeferStrileArrayWriting(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFDeferStrileArrayWriting", f.pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not defer strile array writing")
	}

	return nil
}

// TIFFForceStrileArrayWriting forces writing of previously deferred strip/tile
// offset and byte count arrays.
func (f *File) TIFFForceStrileArrayWriting(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFForceStrileArrayWriting", f.pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("could not force strile array writing")
	}

	return nil
}

// TIFFPrintDirectory returns a human-readable string representation of all tags
// in the current directory. The flags parameter controls verbosity.
func (f *File) TIFFPrintDirectory(ctx context.Context, flags PrintDirectoryFlag) (string, error) {
	const bufSize = 65536

	bufPointer, err := f.instance.malloc(ctx, bufSize)
	if err != nil {
		return "", err
	}
	defer f.instance.free(ctx, bufPointer)

	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFPrintDirectoryToBuffer", f.pointer, bufPointer, api.EncodeI32(int32(bufSize)), uint64(flags))
	if err != nil {
		return "", err
	}

	bytesWritten := api.DecodeI32(results[0])
	if bytesWritten == -1 {
		return "", errors.New("could not print directory")
	}

	f.instance.internalInstance.CallLock.Lock()
	buf, ok := f.instance.internalInstance.Module.Memory().Read(uint32(bufPointer), uint32(bytesWritten))
	if !ok {
		f.instance.internalInstance.CallLock.Unlock()
		return "", errors.New("could not read directory output from WASM memory")
	}
	result := make([]byte, len(buf))
	copy(result, buf)
	f.instance.internalInstance.CallLock.Unlock()

	return string(result), nil
}

// TIFFFlush flushes pending writes to the file, including the directory.
func (f *File) TIFFFlush(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFFlush", f.pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("error flushing TIFF file")
	}

	return nil
}

// TIFFFlushData flushes pending data writes to the file without writing the directory.
func (f *File) TIFFFlushData(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFFlushData", f.pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("error flushing TIFF data")
	}

	return nil
}

// TIFFCheckpointDirectory writes the current directory to the file
// without closing it, allowing further writes to the same directory.
func (f *File) TIFFCheckpointDirectory(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFCheckpointDirectory", f.pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("error checkpointing directory")
	}

	return nil
}

// TIFFRewriteDirectory rewrites the current directory in place.
// This can be used to update a directory after it has been written.
func (f *File) TIFFRewriteDirectory(ctx context.Context) error {
	results, err := f.instance.internalInstance.CallExportedFunction(ctx, "TIFFRewriteDirectory", f.pointer)
	if err != nil {
		return err
	}

	if results[0] == 0 {
		return errors.New("error rewriting directory")
	}

	return nil
}
