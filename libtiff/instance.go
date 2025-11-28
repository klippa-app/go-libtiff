package libtiff

import (
	"context"
	_ "embed"
	"errors"

	"github.com/klippa-app/go-libtiff/internal/instance"
)

//go:embed libtiff.wasm
var wasmBinary []byte

type Instance struct {
	internalInstance *instance.Instance
}

func GetInstance(ctx context.Context, config *Config) (*Instance, error) {
	if config == nil {
		return nil, errors.New("config must be given")
	}
	internalInstance, err := instance.GetInstance(ctx, &instance.Config{
		WASMData:         wasmBinary,
		CompilationCache: config.CompilationCache,
		FSConfig:         config.FSConfig,
		Debug:            config.Debug,
	})
	if err != nil {
		return nil, err
	}
	return &Instance{
		internalInstance: internalInstance,
	}, nil
}

func (i *Instance) Close(ctx context.Context) error {
	return i.internalInstance.Close(ctx)
}
