package modules

import "sync"

var (
	store map[string]any
	lock  sync.Mutex
)

func init() {
	store = map[string]any{}
	lock = sync.Mutex{}
}

type Module = interface{}

var (
	mux = sync.Map{}
)

func Register[T any](moduleName string, moduleCreate func(string) (T, error)) (T, error) {
	locker, _ := mux.LoadOrStore(moduleName, &sync.Mutex{})
	locker.(*sync.Mutex).Lock()
	defer locker.(*sync.Mutex).Unlock()

	var (
		module T
		ok     bool
	)

	module, ok = store[moduleName].(T)
	if ok {
		return module, nil
	}
	module, err := moduleCreate(moduleName)
	if err != nil {
		return module, err
	}
	store[moduleName] = module
	return module, nil
}

func Unregister(moduleName string) {
	lock.Lock()
	defer lock.Unlock()
	delete(store, moduleName)
}
