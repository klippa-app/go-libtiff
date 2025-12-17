package instance

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/klippa-app/go-libtiff/internal/imports"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/experimental/logging"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Config struct {
	CompilationCache wazero.CompilationCache
	WASMData         []byte
	IsProgramRun     bool
	FSConfig         wazero.FSConfig
	Debug            bool
	Stdout           io.Writer
	Stderr           io.Writer
	RandSource       io.Reader
}

type Instance struct {
	runtime        wazero.Runtime
	Module         api.Module
	config         wazero.ModuleConfig
	compiledModule wazero.CompiledModule
	callLock       sync.Mutex
}

// This lock is needed because of a bug in Wazero where sometimes it would
// return an error when you do multiple CompileModule calls at the same time
// on the same WASM binary when using a shared compilation cache.
// See: https://github.com/wazero/wazero/issues/2459
var compilationLock = sync.Mutex{}

func GetInstance(ctx context.Context, config *Config) (*Instance, error) {
	if config == nil {
		return nil, errors.New("config must be given")
	}
	if config.Debug {
		ctx = experimental.WithFunctionListenerFactory(ctx, logging.NewHostLoggingListenerFactory(os.Stdout, logging.LogScopeFilesystem))
	}
	runtimeConfig := wazero.NewRuntimeConfig()
	if config.CompilationCache != nil {
		runtimeConfig = runtimeConfig.WithCompilationCache(config.CompilationCache)
	}

	wazeroRuntime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, wazeroRuntime); err != nil {
		return nil, fmt.Errorf("could not instantiate wasi_snapshot_preview1: %w", err)
	}

	if config.CompilationCache != nil {
		compilationLock.Lock()
	}

	compiledModule, err := wazeroRuntime.CompileModule(ctx, config.WASMData)
	if err != nil {
		if config.CompilationCache != nil {
			compilationLock.Unlock()
		}
		return nil, err
	}

	if config.CompilationCache != nil {
		compilationLock.Unlock()
	}

	if _, err := imports.Instantiate(ctx, wazeroRuntime, compiledModule); err != nil {
		return nil, fmt.Errorf("could not instantiate imports: %w", err)
	}

	fsConfig := config.FSConfig
	if fsConfig == nil {
		fsConfig = wazero.NewFSConfig()

		// On Windows we mount the volume of the current working directory as
		// root. On Linux we mount / as root.
		if runtime.GOOS == "windows" {
			cwdDir, err := os.Getwd()
			if err != nil {
				return nil, err
			}

			volumeName := filepath.VolumeName(cwdDir)
			if volumeName != "" {
				fsConfig = fsConfig.WithDirMount(fmt.Sprintf("%s\\", volumeName), "/")
			}
		} else {
			fsConfig = fsConfig.WithDirMount("/", "/")
		}
	}

	moduleConfig := wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithRandSource(rand.Reader).
		WithFSConfig(fsConfig).
		WithSysNanotime().
		WithSysNanosleep().
		WithName("")

	if config.Stderr != nil {
		moduleConfig = moduleConfig.WithStderr(config.Stderr)
	}

	if config.Stdout != nil {
		moduleConfig = moduleConfig.WithStdout(config.Stdout)
	}

	if config.RandSource != nil {
		moduleConfig = moduleConfig.WithRandSource(config.RandSource)
	}

	if config.IsProgramRun {
		return &Instance{
			runtime:        wazeroRuntime,
			config:         moduleConfig,
			compiledModule: compiledModule,
		}, nil
	}

	moduleConfig = moduleConfig.WithStartFunctions("_initialize")
	mod, err := wazeroRuntime.InstantiateModule(ctx, compiledModule, moduleConfig)
	if err != nil {
		return nil, err
	}

	return &Instance{
		runtime:        wazeroRuntime,
		Module:         mod,
		compiledModule: compiledModule,
	}, nil
}

func (i *Instance) RunProgram(ctx context.Context, args ...string) error {
	i.config = i.config.WithStartFunctions("_start")
	i.config = i.config.WithArgs(args...)
	mod, err := i.runtime.InstantiateModule(ctx, i.compiledModule, i.config)
	i.Module = mod
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) Close(ctx context.Context) error {
	if i.Module != nil {
		if err := i.Module.Close(ctx); err != nil {
			return fmt.Errorf("could not close module: %w", err)
		}
	}
	if i.runtime != nil {
		if err := i.runtime.Close(ctx); err != nil {
			return fmt.Errorf("could not close runtime: %w", err)
		}
	}
	if i.compiledModule != nil {
		if err := i.compiledModule.Close(ctx); err != nil {
			return fmt.Errorf("could not close compiled module: %w", err)
		}
	}
	return nil
}

func (i *Instance) CallExportedFunction(ctx context.Context, name string, args ...uint64) ([]uint64, error) {
	i.callLock.Lock()
	defer i.callLock.Unlock()
	return i.Module.ExportedFunction(name).Call(ctx, args...)
}
