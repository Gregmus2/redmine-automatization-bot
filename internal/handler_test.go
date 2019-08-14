package internal

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"net/http/httptest"
	"redmine-automatization-bot/internal/global"
	_ "redmine-automatization-bot/internal/handlers"
	"redmine-automatization-bot/internal/mocks"
	"redmine-automatization-bot/internal/redmine/redmine_responses"
	"testing"
)

func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/xml")
		timeEntry := redmine_responses.TimeEntriesResponse{
			TimeEntries: []redmine_responses.TimeEntry{
				{
					ID: 1,
					Project: redmine_responses.IDNamePair{
						ID:   1,
						Name: "test",
					},
					Issue: redmine_responses.IDObject{ID: 1},
					Activity: redmine_responses.IDNamePair{
						ID:   1,
						Name: "test",
					},
					Hours:     0.5,
					Comments:  "",
					SpentOn:   "2000-01-01",
					CreatedOn: "2000-01-01",
					UpdatedOn: "2000-01-01",
				},
			},
		}
		bytes, err := json.Marshal(timeEntry)
		if err != nil {
			log.Panic(err)
		}

		_, err = fmt.Fprintln(w, string(bytes))
		if err != nil {
			log.Panic(err)
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

var _ = func() interface{} {
	_testing = true
	return nil
}()

func TestAuthorize(t *testing.T) {
	server := mockServer()
	defer server.Close()

	storage := mocks.NewMockStorage()
	err := storage.Put(global.RedmineUrlsCollection, "1", []byte(server.URL))
	if err != nil {
		t.Fatal("error on put in mock storage")
	}
	err = storage.Put(global.ApiKeysCollection, "1", []byte("SomeTestKey"))
	if err != nil {
		t.Fatal("error on put in mock storage")
	}

	us, err := global.NewUserStorage(storage)
	if err != nil {
		t.Fatal("NewUserStorage return error")
	}

	global.RA = global.NewRedmineApis(us)
	message := tgbotapi.Message{
		From: &tgbotapi.User{
			ID: 1,
		},
	}

	api, status := authorize(&message)
	if status == false {
		t.Fatal("authorize failed")
	}

	if api == nil {
		t.Fatal("nil api returned")
	}
}
