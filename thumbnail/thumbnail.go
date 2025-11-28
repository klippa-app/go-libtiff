package thumbnail

import (
	"context"
	_ "embed"

	"github.com/klippa-app/go-libtiff/internal/registry"
)

//go:embed thumbnail.wasm
var wasmBinary []byte

func init() {
	registry.Register("thumbnail", wasmBinary)
}

func Run(ctx context.Context, args []string) error {
	return registry.Run(ctx, "thumbnail", args...)
}
