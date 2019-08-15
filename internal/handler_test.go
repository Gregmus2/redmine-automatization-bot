package internal

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jarcoal/httpmock"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"redmine-automatization-bot/internal/global"
	_ "redmine-automatization-bot/internal/handlers"
	"redmine-automatization-bot/internal/mocks"
	"redmine-automatization-bot/internal/redmine/redmine_responses"
	"testing"
)

var _ = func() interface{} {
	_testing = true
	return nil
}()

func mockServer() *httptest.Server {
	return mocks.MockServer(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func mockClient() {
	httpmock.Activate()

	res, _ := json.Marshal(tgbotapi.APIResponse{
		Ok:          true,
		Result:      nil,
		ErrorCode:   0,
		Description: "",
		Parameters:  nil,
	})

	httpmock.RegisterResponder("POST", `=~^https://api.telegram.org/bottoken/\w+\z`,
		httpmock.NewStringResponder(200, string(res)))
}

func TestMain(m *testing.M) {
	mockClient()
	defer httpmock.DeactivateAndReset()
	Bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		panic(err)
	}
	Bot.Debug = false

	server := mockServer()
	defer server.Close()

	storage := mocks.NewMockStorage()
	err = storage.Put(global.RedmineUrlsCollection, "1", []byte(server.URL))
	if err != nil {
		panic("error on put in mock storage")
	}
	err = storage.Put(global.ApiKeysCollection, "1", []byte("SomeTestKey"))
	if err != nil {
		panic("error on put in mock storage")
	}

	us, err := global.NewUserStorage(storage)
	if err != nil {
		panic("NewUserStorage return error")
	}

	global.RA = global.NewRedmineApis(us)

	os.Exit(m.Run())
}

func TestExistsAuthorize(t *testing.T) {
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

//func TestNewAuthorize(t *testing.T) {
//	t.Fatal(Bot)
//	message := tgbotapi.Message{
//		From: &tgbotapi.User{
//			ID: 2,
//		},
//		Chat: &tgbotapi.Chat{
//			ID: 1,
//		},
//	}
//
//	api, status := authorize(&message)
//	if status == true {
//		t.Fatal("expected not authorize, but authorize")
//	}
//
//	if api != nil {
//		t.Fatal("expected nil api var")
//	}
//}
