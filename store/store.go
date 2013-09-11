package store

import (
	cfg "../config"
)

type Store interface {
	Put(name string, key string, data string) error
	Get(name string, key string) (string, error)
	Exists(name string, key string) bool
	Delete(name, key string) error
	DeleteName(name string) error
}

func EnsureStore() Store {

	storeType := cfg.GetOrElse("metadata_store_type", "memory")
	var store Store
	switch storeType {
	case "memory":
		store = EnsureKVMemoryStore()
	case "file":
		store = EnsureKVStore()
	case "riak":
		store = EnsureRiakStore()
	default:
		store = EnsureKVMemoryStore()
	}
	return store
}
