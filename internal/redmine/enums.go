package redmine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"redmine-automatization-bot/internal/redmine/redmine_responses"
	"sort"
	"sync"
)

type Activities struct {
	mx  sync.Mutex
	m   map[uint]string
	ids []uint
}

func NewActivities() Activities {
	return Activities{m: make(map[uint]string)}
}

func (a *Activities) FindAll() map[uint]string {
	return a.m
}

func (a *Activities) Add(id uint, name string) {
	a.mx.Lock()
	defer a.mx.Unlock()
	_, exists := a.m[id]
	if exists {
		return
	}

	a.m[id] = name
	a.ids = append(a.ids, id)
	sort.Slice(a.ids, func(i, j int) bool { return a.ids[i] < a.ids[j] })
}

func (a *Activities) ToText() string {
	b := new(bytes.Buffer)
	for _, value := range a.ids {
		_, _ = fmt.Fprintf(b, "%d - %s\n", value, a.m[value])
	}
	return b.String()
}

func (api *Api) CollectActivities() {
	req, err := http.NewRequest("GET", api.url+"time_entries.json", nil)
	if err != nil {
		log.Panic(err)
		return
	}

	q := req.URL.Query()
	q.Add("limit", "100")
	req.URL.RawQuery = q.Encode()

	resp, err := api.request(req)
	if err != nil {
		log.Panic(err)
		return
	}

	timeEntries := &redmine_responses.TimeEntriesResponse{}
	err = json.Unmarshal(resp, timeEntries)
	if err != nil {
		log.Panic(err)
		return
	}

	for _, timeEntry := range timeEntries.TimeEntries {
		api.Activities.Add(timeEntry.Activity.ID, timeEntry.Activity.Name)
	}
}
