package redmine

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Api struct {
	url        string
	apiKey     string
	Activities Activities
}

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func NewApi(apiUrl string, apiKey string) (*Api, error) {
	u, err := url.Parse(apiUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid url")
	}

	if apiUrl[len(apiUrl)-1:] != "/" {
		apiUrl = apiUrl + "/"
	}

	api := &Api{url: apiUrl, apiKey: apiKey, Activities: NewActivities()}
	go api.CollectActivities()

	return api, nil
}

func (api *Api) Create(method string, jsonBuffer []byte) ([]byte, error) {
	req, err := http.NewRequest(
		"POST",
		api.url+method,
		bytes.NewReader(jsonBuffer),
	)
	if err != nil {
		return nil, err
	}

	return api.request(req)
}

func (api *Api) request(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Redmine-API-Key", api.apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
