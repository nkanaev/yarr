package router

import (
	"encoding/json"
	"log"
	"net/http"
)

type Context struct {
	Req   *http.Request
	Out   http.ResponseWriter

	Vars  map[string]string

	chain []Handler
	index int
}

func (c *Context) Next() {
	c.index++
	c.chain[c.index](c)
}

func (c *Context) JSON(status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	c.Out.WriteHeader(status)
	c.Out.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.Out.Write(body)
	c.Out.Write([]byte("\n"))
}
