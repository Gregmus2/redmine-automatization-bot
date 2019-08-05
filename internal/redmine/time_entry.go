package redmine

import (
	"encoding/json"
	"time"
)

type TimeEntry struct {
	TimeEntry TimeEntryBody `json:"time_entry"`
}

type TimeEntryBody struct {
	IssueId    uint    `json:"issue_id"`
	SpentOn    string  `json:"spent_on"`
	Hours      float32 `json:"hours"`
	ActivityId uint8   `json:"activity_id"`
	Comments   string  `json:"comments"`
}

const TimeEntriesUri string = "time_entries.json"

func (api *Api) CreateTimeEntry(timeEntryBody TimeEntryBody) ([]byte, error) {
	timeEntryBody.SpentOn = time.Now().Format("2006-01-02")
	timeEntry := TimeEntry{TimeEntry: timeEntryBody}

	jsonBuffer, err := json.Marshal(timeEntry)
	if err != nil {
		return nil, err
	}

	return api.Create(TimeEntriesUri, jsonBuffer)
}
