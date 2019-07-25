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
	url    string
	apiKey string
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

	return &Api{url: apiUrl, apiKey: apiKey}, nil
}

func (api *Api) Create(method string, jsonBuffer []byte) (string, error) {
	req, err := http.NewRequest(
		"POST",
		api.url+method,
		bytes.NewReader(jsonBuffer),
	)
	if err != nil {
		return "", err
	}

	return api.request(req)
}

func (api *Api) request(req *http.Request) (string, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Redmine-API-Key", api.apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
