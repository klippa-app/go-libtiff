package fax2ps

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed fax2ps.wasm
var wasmBinary []byte

func init() {
	registry.Register("fax2ps", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "fax2ps", args...)
}
