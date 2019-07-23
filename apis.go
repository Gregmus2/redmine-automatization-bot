package main

import (
	"redmine-automatization-bot/redmine"
	"sync"
)

type RedmineApis struct {
	mx   sync.Mutex
	apis map[int]*redmine.Api
}

func NewRedmineApis(userStorage *UserStorage) *RedmineApis {
	redmineApis = &RedmineApis{apis: make(map[int]*redmine.Api)}
	for _, user := range userStorage.users {
		redmineApis.apis[user.Id] = redmine.NewApi(user.RedmineApiKey)
	}

	return redmineApis
}

func (apis *RedmineApis) Find(key int) (*redmine.Api, bool) {
	apis.mx.Lock()
	defer apis.mx.Unlock()
	val, ok := apis.apis[key]

	return val, ok
}

func (apis *RedmineApis) Save(key int, value *redmine.Api) {
	apis.mx.Lock()
	defer apis.mx.Unlock()
	apis.apis[key] = value
}
