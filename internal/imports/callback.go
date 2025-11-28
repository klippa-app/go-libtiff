package imports

import (
	"context"
	"log"

	"github.com/tetratelabs/wazero/api"
)

type TIFFReadWriteProcGoCB struct {
}

func (cb TIFFReadWriteProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	log.Println("TIFFReadWriteProcGoCB")
	return
}

type TIFFSeekProcGoCB struct {
}

func (cb TIFFSeekProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	log.Println("TIFFSeekProcGoCB")
	return
}

type TIFFCloseProcGoCB struct {
}

func (cb TIFFCloseProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	log.Println("TIFFCloseProcGoCB")
	return
}

type TIFFSizeProcGoCB struct {
}

func (cb TIFFSizeProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	log.Println("TIFFSizeProcGoCB")
	return
}

type TIFFMapFileProcGoCB struct {
}

func (cb TIFFMapFileProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	log.Println("TIFFMapFileProcGoCB")
	return
}

type TIFFUnmapFileProcGoCB struct {
}

func (cb TIFFUnmapFileProcGoCB) Call(ctx context.Context, mod api.Module, stack []uint64) {
	log.Println("TIFFUnmapFileProcGoCB")
	return
}
