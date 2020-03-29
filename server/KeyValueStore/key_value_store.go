package KeyValueStore

import "sync"

var (
	keyValueStore = make(map[string]string)
	rwm sync.RWMutex
)

func PUT(key, value string) {
	rwm.Lock()
	defer rwm.Unlock()
	keyValueStore[key] = value
}

func GET(key string) string {
	return keyValueStore[key]
}

func DELETE(key string) string {
	value := GET(key)

	rwm.Lock()
	defer rwm.Unlock()
	delete(keyValueStore, key)
	return value
}

type KeyValueStore interface {
	PUT(key, value string) string
	GET(key string) string
	DELETE(key string) string
}
