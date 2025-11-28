package libtiff

import (
	"context"
	"errors"
)

type CString struct {
	Pointer uint64
	Free    func(ctx context.Context) error
}

func (i *Instance) NewCString(ctx context.Context, input string) (*CString, error) {
	inputLength := uint64(len(input)) + 1

	pointer, err := i.Malloc(ctx, inputLength)
	if err != nil {
		return nil, err
	}

	// Write string + null terminator.
	if !i.internalInstance.Module.Memory().Write(uint32(pointer), append([]byte(input), byte(0))) {
		return nil, errors.New("could not write CString data")
	}

	return &CString{
		Pointer: pointer,
		Free: func(ctx context.Context) error {
			return i.Free(ctx, pointer)
		},
	}, nil
}

func (i *Instance) Malloc(ctx context.Context, size uint64) (uint64, error) {
	results, err := i.internalInstance.Module.ExportedFunction("malloc").Call(ctx, size)
	if err != nil {
		return 0, err
	}

	pointer := results[0]
	ok := i.internalInstance.Module.Memory().Write(uint32(results[0]), make([]byte, size))
	if !ok {
		return 0, errors.New("could not write nulls to memory")
	}

	return pointer, nil
}

func (i *Instance) Free(ctx context.Context, pointer uint64) error {
	_, err := i.internalInstance.Module.ExportedFunction("free").Call(ctx, pointer)
	if err != nil {
		return err
	}
	return nil
}
