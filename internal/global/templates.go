package global

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Templates struct {
	TextCommand map[string]string `json:"text_command"`
}

type TemplateStorage struct {
	templates map[int]Templates
	storage   Storage
}

const TemplatesCollection string = "templates"

func NewTemplateStorage(storage Storage) (*TemplateStorage, error) {
	storage.CreateCollectionIfNotExist(TemplatesCollection)

	ts := &TemplateStorage{
		templates: make(map[int]Templates),
		storage:   storage,
	}

	templatesJson, err := storage.GetAllRaw(TemplatesCollection)
	if err != nil {
		return nil, err
	}

	for userId, template := range templatesJson {
		templates := Templates{}
		err := json.Unmarshal(template, &templates)
		if err != nil {
			return nil, err
		}

		intUserId, err := strconv.Atoi(userId)
		if err != nil {
			return nil, err
		}

		ts.templates[intUserId] = templates
	}

	return ts, nil
}

func (ts *TemplateStorage) AddTemplate(userId int, name string, command string) error {
	templates, exists := ts.templates[userId]
	if !exists {
		templates = Templates{TextCommand: make(map[string]string)}
		ts.templates[userId] = templates
	}

	templates.TextCommand[name] = command

	bytes, err := json.Marshal(templates)
	if err != nil {
		return err
	}

	err = ts.storage.Put(TemplatesCollection, strconv.Itoa(userId), bytes)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TemplateStorage) RemoveTemplate(userId int, name string) error {
	templates, exists := ts.templates[userId]
	if !exists {
		return errors.New("you haven't any templates")
	}

	_, exists = templates.TextCommand[name]
	if !exists {
		return errors.New("template doesn't exist")
	}

	delete(templates.TextCommand, name)

	return nil
}

func (ts *TemplateStorage) GetTemplateNames(userId int) []string {
	templates, exists := ts.templates[userId]
	if !exists {
		return []string{}
	}

	commandsCount := len(templates.TextCommand)
	names := make([]string, commandsCount)
	for name, _ := range templates.TextCommand {
		names = append(names, name)
	}

	return names
}
