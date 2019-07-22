package redmine

import (
	"encoding/json"
	"time"
)

type TimeEntry struct {
	TimeEntry TimeEntryBody `json:"time_entry"`
}

type TimeEntryBody struct {
	IssueId uint `json:"issue_id"`
	SpentOn string `json:"spent_on"`
	Hours float32 `json:"hours"`
	ActivityId uint8 `json:"activity_id"`
	Comments string `json:"comments"`
}

const TimeEntriesUri string = "time_entries.json"

func (api *Api) CreateTimeEntry(issue uint, hours float32, activity uint8, comments string) (string, error) {
	date := time.Now().Format("2006-01-02")
	timeEntry := TimeEntry{TimeEntry: TimeEntryBody{
		IssueId: issue,
		SpentOn: date,
		Hours: hours,
		ActivityId: activity,
		Comments: comments,
	}}

	jsonBuffer, err := json.Marshal(timeEntry)
	if err != nil {
		return "", err
	}

	return api.Create(TimeEntriesUri, jsonBuffer)
}