package redmine

import (
	"encoding/json"
	"time"
)

type Issue struct {
	Issue IssueBody `json:"issue"`
}

type IssueBody struct {
	ProjectId      string  `json:"project_id"`
	TrackerId      uint    `json:"tracker_id"`
	StatusId       uint    `json:"status_id"`
	PriorityId     uint    `json:"priority_id"`
	Subject        string  `json:"subject"`
	Description    string  `json:"description"`
	AssignedToId   uint    `json:"assigned_to_id"`
	ParentIssueId  uint    `json:"parent_issue_id"`
	StartDate      string  `json:"start_date"`
	EstimatedHours float32 `json:"estimated_hours"`
}

func (api *Api) CreateIssue(issueBody IssueBody) ([]byte, error) {
	issueBody.StartDate = time.Now().Format("2006-01-02")
	timeEntry := Issue{Issue: issueBody}

	jsonBuffer, err := json.Marshal(timeEntry)
	if err != nil {
		return nil, err
	}

	return api.Create(TimeEntriesUri, jsonBuffer)
}
