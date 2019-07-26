package redmine_responses

type TimeEntriesResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
}

type TimeEntry struct {
	ID        uint       `json:"id"`
	Project   IDNamePair `json:"project"`
	Issue     IDObject   `json:"issue"`
	Activity  IDNamePair `json:"activity"`
	Hours     float32    `json:"hours"`
	Comments  string     `json:"comments"`
	SpentOn   string     `json:"spent_on"`
	CreatedOn string     `json:"created_on"`
	UpdatedOn string     `json:"updated_on"`
}

type IDNamePair struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type IDObject struct {
	ID uint `json:"id"`
}
