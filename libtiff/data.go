package libtiff

import (
	"context"
	"errors"
)

type cString struct {
	Pointer uint64
	Free    func(ctx context.Context) error
}

func (i *Instance) newCString(ctx context.Context, input string) (*cString, error) {
	inputLength := uint64(len(input)) + 1

	pointer, err := i.malloc(ctx, inputLength)
	if err != nil {
		return nil, err
	}

	// Write string + null terminator.
	if !i.internalInstance.Module.Memory().Write(uint32(pointer), append([]byte(input), byte(0))) {
		return nil, errors.New("could not write cString data")
	}

	return &cString{
		Pointer: pointer,
		Free: func(ctx context.Context) error {
			return i.free(ctx, pointer)
		},
	}, nil
}

func (i *Instance) readCString(pointer uint32) string {
	cStringData := []byte{}
	for {
		data, success := i.internalInstance.Module.Memory().Read(pointer, 1)
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

func (i *Instance) malloc(ctx context.Context, size uint64) (uint64, error) {
	results, err := i.internalInstance.CallExportedFunction(ctx, "malloc", size)
	if err != nil {
		return 0, err
	}

	pointer := results[0]
	ok := i.internalInstance.Module.Memory().Write(uint32(pointer), make([]byte, size))
	if !ok {
		return 0, errors.New("could not write nulls to memory")
	}

	return pointer, nil
}

func (i *Instance) free(ctx context.Context, pointer uint64) error {
	_, err := i.internalInstance.CallExportedFunction(ctx, "free", pointer)
	if err != nil {
		return err
	}
	return nil
}
