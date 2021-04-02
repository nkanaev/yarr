package server

import (
	"net/http/httptest"
	"testing"
)

func TestStatic(t *testing.T) {
	handler := NewServer(nil, "127.0.0.1:8000").handler()
	url := "/static/javascripts/app.js"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 200 {
		t.FailNow()
	}
}

func TestStaticWithBase(t *testing.T) {
	server := NewServer(nil, "127.0.0.1:8000")
	server.BasePath = "/sub"
	
	handler := server.handler()
	url := "/sub/static/javascripts/app.js"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", url, nil)
	handler.ServeHTTP(recorder, request)
	if recorder.Result().StatusCode != 200 {
		t.FailNow()
	}
}
