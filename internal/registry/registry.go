package registry

import (
	"context"
	"sort"
	"sync"

	"github.com/klippa-app/go-libtiff/internal/instance"
)

var programRegistry = make(map[string][]byte)
var lock sync.RWMutex

func Register(name string, program []byte) {
	lock.Lock()
	defer lock.Unlock()
	programRegistry[name] = program
}

func Get(name string) []byte {
	lock.RLock()
	defer lock.RUnlock()
	return programRegistry[name]
}

func List() []string {
	lock.RLock()
	defer lock.RUnlock()

	keys := make([]string, len(programRegistry))

	i := 0
	for k := range programRegistry {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return keys
}

func Run(ctx context.Context, name string, args ...string) error {
	wazeroInstance, err := instance.GetInstance(ctx, &instance.Config{
		WASMData:     programRegistry[name],
		IsProgramRun: true,
	})
	if err != nil {
		return err
	}
	defer wazeroInstance.Close(ctx)

	if err := wazeroInstance.RunProgram(ctx, args...); err != nil {
		return err
	}

	return nil
}
