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
