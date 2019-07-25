package internal

import (
	"redmine-automatization-bot/internal/redmine"
	"sync"
)

type RedmineApis struct {
	mx   sync.Mutex
	apis map[int]*redmine.Api
}

func NewRedmineApis(userStorage *UserStorage) *RedmineApis {
	redmineApis = &RedmineApis{apis: make(map[int]*redmine.Api)}
	for _, user := range userStorage.users {
		if user.RedmineApiKey == "" {
			continue
		}

		api, err := redmine.NewApi(user.RedmineUrl, user.RedmineApiKey)
		if err != nil {
			panic(err)
		}

		redmineApis.apis[user.Id] = api
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
