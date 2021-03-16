package router

import (
	"io"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	middlecalled := false
	router := NewRouter()
	router.Use(func(c *Context) {
		middlecalled = true
		c.Next()
	})
	router.For("/hello/:place", func(c *Context) {
		c.Out.Write([]byte(c.Vars["place"]))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/hello/world", nil)
	router.ServeHTTP(recorder, request)
	body, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		t.Error(err)
	}

	if !middlecalled {
		t.Error("middleware not called")
	}
	if recorder.Result().StatusCode != 200 {
		t.Error("expected 200")
	}
	if string(body) != "world" {
		t.Errorf("invalid response body, got %#v", body)
	}
}

func TestRouterPaths(t *testing.T) {
	router := NewRouter()
	router.For("/path/to/foo", func(c *Context) {
		c.Out.Write([]byte("foo"))
	})
	router.For("/path/to/bar", func(c *Context) {
		c.Out.Write([]byte("bar"))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/path/to/bar", nil)
	router.ServeHTTP(recorder, request)

	body, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "bar" {
		t.Error("expected 2nd route to be called")
	}
}

func TestRouterMiddlewareIntercept(t *testing.T) {
	router := NewRouter()
	router.Use(func(c *Context) {
		c.Out.WriteHeader(404)
	})
	router.For("/hello/:place", func(c *Context) {
		c.Out.WriteHeader(200)
		c.Out.Write([]byte(c.Vars["place"]))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/hello/world", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 404 {
		t.Error("expected 404")
	}
	body, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		t.Error(err)
	}
	if len(body) != 0 {
		t.Errorf("expected empty body, got %v", body)
	}
}

func TestRouterMiddlewareOrder(t *testing.T) {
	router := NewRouter()

	router.Use(func(c *Context) {
		c.Out.Write([]byte("foo"))
		c.Next()
	})
	router.Use(func(c *Context) {
		c.Out.Write([]byte("bar"))
		c.Next()
	})
	router.For("/hello/:place", func(c *Context) {
		c.Out.Write([]byte("!!!"))
	})

	router.Use(func(c *Context) {
		c.Out.Write([]byte("baz"))
		c.Next()
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/hello/world", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Result().StatusCode != 200 {
		t.Error("expected 200")
	}
	body, err := io.ReadAll(recorder.Result().Body)
	if err != nil {
		t.Error(err)
	}
	if string(body) != "foobar!!!" {
		t.Errorf("invalid body, got %#v", string(body))
	}
}
