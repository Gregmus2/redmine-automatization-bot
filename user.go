package main

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"strconv"
	"sync"
)

type User struct {
	Id            int
	RedmineApiKey string
}

type UserStorage struct {
	mx      sync.Mutex
	users   map[int]User
	storage Storage
}

const COLLECTION string = "users"

func NewUserStorage(storage Storage) (*UserStorage, error) {
	storage.CreateCollectionIfNotExist(COLLECTION)

	us := &UserStorage{
		users:   make(map[int]User),
		storage: storage,
	}

	data, err := storage.GetAll(COLLECTION)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		id, err := strconv.Atoi(k)
		if err != nil {
			return nil, err
		}

		us.users[id] = User{
			Id:            id,
			RedmineApiKey: v,
		}
	}

	return us, nil
}

func (us *UserStorage) Find(id int) *User {
	us.mx.Lock()
	defer us.mx.Unlock()

	user, exists := us.users[id]
	if !exists {
		return nil
	}

	return &user
}

func (us *UserStorage) Register(id int, apiKey string) (*User, error) {
	us.mx.Lock()
	defer us.mx.Unlock()

	user, exists := us.users[id]
	if exists {
		return &user, errors.New("user already exists")
	}

	err := us.storage.Put(COLLECTION, strconv.Itoa(id), apiKey)
	if err != nil {
		return nil, err
	}

	user = User{Id: id, RedmineApiKey: apiKey}
	us.users[id] = user

	return &user, nil
}
