package main

type Storage interface {
	Close()
	GetAll(collection string) (map[string]string, error)
	Put(collection string, key string, value string) error
	CreateCollectionIfNotExist(collection string)
}
