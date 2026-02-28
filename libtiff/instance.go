package libtiff

import (
	"context"
	_ "embed"
	"errors"

	"github.com/klippa-app/go-libtiff/internal/instance"
	"github.com/tetratelabs/wazero/api"
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
		Stdout:           config.Stdout,
		Stderr:           config.Stderr,
		RandSource:       config.RandSource,
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

// TIFFGetConfiguredCODECs returns the list of all configured (available) compression codecs.
// The caller can use this to discover which compression schemes are supported.
func (i *Instance) TIFFGetConfiguredCODECs(ctx context.Context) ([]Codec, error) {
	results, err := i.internalInstance.CallExportedFunction(ctx, "TIFFGetConfiguredCODECs")
	if err != nil {
		return nil, err
	}

	if results[0] == 0 {
		return nil, errors.New("could not get configured codecs")
	}

	arrayPointer := uint32(results[0])

	// The TIFFCodec struct in WASM32 is:
	//   char *name;          // 4 bytes (offset 0)
	//   uint16_t scheme;     // 2 bytes (offset 4)
	//   // 2 bytes padding   (offset 6)
	//   TIFFInitMethod init; // 4 bytes (offset 8)
	// Total: 12 bytes per entry
	const codecStructSize = 12

	var codecs []Codec

	i.internalInstance.CallLock.Lock()
	defer i.internalInstance.CallLock.Unlock()

	for offset := uint32(0); ; offset += codecStructSize {
		namePtr, ok := i.internalInstance.Module.Memory().ReadUint32Le(arrayPointer + offset)
		if !ok {
			break
		}

		scheme, ok := i.internalInstance.Module.Memory().ReadUint16Le(arrayPointer + offset + 4)
		if !ok {
			break
		}

		// End of array: name pointer is NULL and scheme is 0.
		if namePtr == 0 && scheme == 0 {
			break
		}

		// Read the codec name - need to temporarily unlock since readCString locks.
		i.internalInstance.CallLock.Unlock()
		name := i.readCString(namePtr)
		i.internalInstance.CallLock.Lock()

		codecs = append(codecs, Codec{
			Name:   name,
			Scheme: scheme,
		})
	}

	// Free the array allocated by TIFFGetConfiguredCODECs.
	i.internalInstance.CallLock.Unlock()
	err = i.free(ctx, uint64(arrayPointer))
	i.internalInstance.CallLock.Lock()
	if err != nil {
		return codecs, err
	}

	return codecs, nil
}

// TIFFDataWidth returns the size in bytes of the given TIFF data type.
func (i *Instance) TIFFDataWidth(ctx context.Context, dataType TIFFDataType) (int, error) {
	results, err := i.internalInstance.CallExportedFunction(ctx, "TIFFDataWidth", api.EncodeU32(uint32(dataType)))
	if err != nil {
		return 0, err
	}

	return int(api.DecodeI32(results[0])), nil
}
