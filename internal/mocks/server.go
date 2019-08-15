package mocks

import (
	"net/http"
	"net/http/httptest"
)

func MockServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}
