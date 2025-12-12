package imports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	tiffErrors "github.com/klippa-app/go-libtiff/errors"

	"github.com/tetratelabs/wazero/api"
)

type File struct {
	ParamPointer uint64
	FileSize     uint64
	Reader       io.ReadSeeker
	Error        error
	WarnHandler  func(module string, message string)
}

func (f *File) GetError() error {
	if f.Error == nil {
		return nil
	}

	// Clear out error.
	var copyError = f.Error
	f.Error = nil

	return copyError
}

var FileReaders = struct {
	Refs    map[uint32]*File
	Counter uint32
	Mutex   *sync.RWMutex
}{
	Refs:    map[uint32]*File{},
	Counter: 0,
	Mutex:   &sync.RWMutex{},
}

type TIFFReadWriteProcGoCB struct {
}

func (cb TIFFReadWriteProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	paramPointer := uint32(stack[0])
	pBufPointer := uint32(stack[1])
	size := uint32(stack[2])

	mem := mod.Memory()
	param, ok := mem.ReadUint32Le(paramPointer)
	if !ok {
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	// Check if we have the file referenced in param.
	FileReaders.Mutex.RLock()
	openFile, ok := FileReaders.Refs[param]
	FileReaders.Mutex.RUnlock()
	if !ok {
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	// Read the requested data into a buffer.
	readBuffer := make([]byte, size)
	n, err := openFile.Reader.Read(readBuffer)

	// Clear out the error if we have EOF but read the requested size.
	// This is to handle some edge case clients that return EOF as err when
	// reading the exact amount of bytes requested until the end of the file.
	if err != nil && errors.Is(err, io.EOF) && n == int(size) {
		err = nil
	}

	if n == 0 || err != nil {
		if err != nil && openFile.WarnHandler != nil {
			openFile.WarnHandler("TIFFReadWriteProcGoCB", fmt.Sprintf("Read %d (requested %d) bytes with err: %v", n, size, err))
		}
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	ok = mem.Write(pBufPointer, readBuffer)
	if !ok {
		if openFile.WarnHandler != nil {
			openFile.WarnHandler("TIFFReadWriteProcGoCB", fmt.Sprintf("Could not write memory at %d", pBufPointer))
		}
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	stack[0] = uint64(n)
	return
}

type TIFFSeekProcGoCB struct {
}

func (cb TIFFSeekProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	paramPointer := uint32(stack[0])
	offset := uint32(stack[1])
	whence := api.DecodeI32(stack[2])

	mem := mod.Memory()
	param, ok := mem.ReadUint32Le(paramPointer)
	if !ok {
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	// Check if we have the file referenced in param.
	FileReaders.Mutex.RLock()
	openFile, ok := FileReaders.Refs[param]
	FileReaders.Mutex.RUnlock()
	if !ok {
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	newOffset, err := openFile.Reader.Seek(int64(offset), int(whence))
	if err != nil {
		if openFile.WarnHandler != nil {
			openFile.WarnHandler("TIFFSeekProcGoCB", fmt.Sprintf("Could not seek to %d with whence %d: %v", offset, whence, err))
		}
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	stack[0] = uint64(newOffset)

	return
}

type TIFFCloseProcGoCB struct {
}

func (cb TIFFCloseProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	// We don't really have to do anything on close.
	// File reader ownership is handled by the user.
	return
}

type TIFFSizeProcGoCB struct {
}

func (cb TIFFSizeProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	paramPointer := uint32(stack[0])

	mem := mod.Memory()
	param, ok := mem.ReadUint32Le(paramPointer)
	if !ok {
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	// Check if we have the file referenced in param.
	FileReaders.Mutex.RLock()
	openFile, ok := FileReaders.Refs[param]
	FileReaders.Mutex.RUnlock()
	if !ok {
		stack[0] = uint64(0) // Should we return -1 like libtiff here? How does that work with uint?
		return
	}

	stack[0] = openFile.FileSize

	return
}

type TIFFMapFileProcGoCB struct {
}

func (cb TIFFMapFileProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	// We don't support this.
	stack[0] = uint64(0)
	return
}

type TIFFUnmapFileProcGoCB struct{}

func (cb TIFFUnmapFileProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	// We don't support this.
	return
}

func readCString(mod api.Module, pointer uint32) string {
	cStringData := []byte{}
	for {
		data, success := mod.Memory().Read(pointer, 1)
		if !success {
			return string(cStringData)
		}

		if data[0] == 0x00 {
			break
		}

		cStringData = append(cStringData, data[0])
		pointer++
	}

	return string(cStringData)
}

func alloc(ctx context.Context, mod api.Module, size uint64) uint64 {
	results, err := mod.ExportedFunction("malloc").Call(ctx, size)
	if err != nil {
		return 0
	}

	pointer := results[0]
	ok := mod.Memory().Write(uint32(results[0]), make([]byte, size))
	if !ok {
		return 0
	}

	return pointer
}

func formatError(ctx context.Context, mod api.Module, fmtPointer uint32, vaListPointer uint32) error {
	// 1024 must be enough right?
	errorPointer := alloc(ctx, mod, 1024)
	defer func() {
		// Cleanup error string
		mod.ExportedFunction("free").Call(ctx, errorPointer)
	}()
	results, err := mod.ExportedFunction("vsprintf").Call(ctx, errorPointer, api.EncodeU32(fmtPointer), api.EncodeU32(vaListPointer))
	if err != nil {
		return err
	}
	if results[0] == 0 {
		return errors.New("could not read error from memory")
	}
	errorText := readCString(mod, uint32(errorPointer))
	return errors.New(errorText)
}

type TIFFOpenOptionsSetErrorHandlerExtRGoCB struct {
}

func (cb TIFFOpenOptionsSetErrorHandlerExtRGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	// Returning 0 will cause libtiff to use the default error handler.
	_ = stack[0] // Pointer to tiff file.
	userDataPointer := uint32(stack[1])
	modulePointer := uint32(stack[2])
	fmtPointer := uint32(stack[3])
	vaListPointer := uint32(stack[4])

	mem := mod.Memory()
	param, ok := mem.ReadUint32Le(userDataPointer)
	if !ok {
		stack[0] = uint64(0)
		return
	}

	// Check if we have the file referenced in param.
	FileReaders.Mutex.RLock()
	openFile, ok := FileReaders.Refs[param]
	FileReaders.Mutex.RUnlock()
	if !ok {
		stack[0] = uint64(0)
		return
	}

	openFile.Error = &tiffErrors.TiffError{
		Module:    readCString(mod, modulePointer),
		TiffError: formatError(ctx, mod, fmtPointer, vaListPointer),
	}

	stack[0] = uint64(1)
	return
}

type TIFFOpenOptionsSetWarningHandlerExtRGoCB struct {
}

func (cb TIFFOpenOptionsSetWarningHandlerExtRGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	// Returning 0 will cause libtiff to use the default error handler.
	_ = stack[0] // Pointer to tiff file.
	userDataPointer := uint32(stack[1])
	modulePointer := uint32(stack[2])
	fmtPointer := uint32(stack[3])
	vaListPointer := uint32(stack[4])

	mem := mod.Memory()
	param, ok := mem.ReadUint32Le(userDataPointer)
	if !ok {
		stack[0] = uint64(0)
		return
	}

	// Check if we have the file referenced in param.
	FileReaders.Mutex.RLock()
	openFile, ok := FileReaders.Refs[param]
	FileReaders.Mutex.RUnlock()
	if !ok {
		stack[0] = uint64(0)
		return
	}

	if openFile.WarnHandler != nil {
		openFile.WarnHandler(readCString(mod, modulePointer), formatError(ctx, mod, fmtPointer, vaListPointer).Error())
	}

	stack[0] = uint64(1)
	return
}
