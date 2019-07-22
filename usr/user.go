package usr

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"sync"
)

type User struct {
	Id int
	RedmineApiKey string
}

type UserStorage struct {
	mx sync.Mutex
	users map[int]User
}

var users *UserStorage

func init() {
	// todo чтение из базы
	users = &UserStorage{users: make(map[int]User)}
}

func Find(id int) *User {
	users.mx.Lock()
	defer users.mx.Unlock()

	user, exists := users.users[id]
	if !exists {
		return nil
	}

	return &user
}

func Register(id int, apiKey string) (*User, error) {
	users.mx.Lock()
	defer users.mx.Unlock()

	user, exists := users.users[id]
	if exists {
		return &user, errors.New("user already exists")
	}

	user = User{Id:id, RedmineApiKey:apiKey}
	users.users[id] = user
	// todo сохранить в базу

	return &user, nil
}