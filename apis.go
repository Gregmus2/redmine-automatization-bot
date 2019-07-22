package main

import (
	"redmine-automatization-bot/redmine"
	"sync"
)

type RedmineApis struct {
	mx sync.Mutex
	apis map[int]*redmine.Api
}

var redmineApis *RedmineApis

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