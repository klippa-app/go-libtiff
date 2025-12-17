package libtiff

import (
	"context"
	"errors"
)

func (i *Instance) TIFFGetVersion(ctx context.Context) (string, error) {
	results, err := i.internalInstance.CallExportedFunction(ctx, "TIFFGetVersion")
	if err != nil {
		return "", err
	}

	if results[0] == 0 {
		return "", errors.New("could not call TIFFGetVersion")
	}

	stringPointer := results[0]
	readValue := i.readCString(uint32(stringPointer))

	return readValue, nil
}
