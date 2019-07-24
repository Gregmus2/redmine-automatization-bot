package main

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"strconv"
	"sync"
)

type User struct {
	Id            int
	RedmineApiKey string
	RedmineUrl    string
}

type UserStorage struct {
	mx      sync.Mutex
	users   map[int]User
	storage Storage
}

const ApiKeysCollection string = "api_keys"
const RedmineUrlsCollection string = "redmine_urls"

func NewUserStorage(storage Storage) (*UserStorage, error) {
	storage.CreateCollectionIfNotExist(RedmineUrlsCollection)
	storage.CreateCollectionIfNotExist(ApiKeysCollection)

	us := &UserStorage{
		users:   make(map[int]User),
		storage: storage,
	}

	urls, err := storage.GetAll(RedmineUrlsCollection)
	if err != nil {
		return nil, err
	}

	keys, err := storage.GetAll(ApiKeysCollection)
	if err != nil {
		return nil, err
	}

	for id, v := range urls {
		intId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}

		us.users[intId] = User{
			Id:            intId,
			RedmineUrl:    v,
			RedmineApiKey: keys[id],
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

func (us *UserStorage) Register(userId int, redmineUrl string) (*User, error) {
	us.mx.Lock()
	defer us.mx.Unlock()

	user, exists := us.users[userId]
	if exists {
		return &user, errors.New("user already exists")
	}

	err := us.storage.Put(RedmineUrlsCollection, strconv.Itoa(userId), redmineUrl)
	if err != nil {
		return nil, err
	}

	user = User{Id: userId, RedmineUrl: redmineUrl}
	us.users[userId] = user

	return &user, nil
}

func (us *UserStorage) AddApiKey(userId int, apiKey string) error {
	err := us.storage.Put(ApiKeysCollection, strconv.Itoa(userId), apiKey)
	if err != nil {
		return err
	}

	user := us.Find(userId)
	user.RedmineApiKey = apiKey

	return nil
}
