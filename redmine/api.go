package redmine

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

type Api struct {
	apiKey string
}

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func NewApi(apiKey string) *Api {
	return &Api{apiKey: apiKey}
}

func (api *Api) Create(method string, jsonBuffer []byte) (string, error) {
	req, err := http.NewRequest(
		"POST",
		"https://redmine.netpeak.net/" + method,
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