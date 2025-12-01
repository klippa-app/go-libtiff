package imports

import (
	"context"
	"io"
	"sync"

	"github.com/tetratelabs/wazero/api"
)

type File struct {
	ParamPointer uint64
	FileSize     uint64
	Reader       io.ReadSeeker
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

	// Read the requested data into a buffer.
	readBuffer := make([]byte, size)
	n, err := openFile.Reader.Read(readBuffer)
	if n == 0 || err != nil {
		stack[0] = uint64(0)
		return
	}

	ok = mem.Write(pBufPointer, readBuffer)
	if !ok {
		stack[0] = uint64(0)
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
	whence := int(stack[2])

	mem := mod.Memory()
	param, ok := mem.ReadUint32Le(paramPointer)
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

	newOffset, err := openFile.Reader.Seek(int64(offset), whence)
	if err != nil {
		stack[0] = uint64(0)
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
