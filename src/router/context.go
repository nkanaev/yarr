package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

func (c *Context) HTML(status int, tmpl *template.Template, data interface{}) {
	c.Out.WriteHeader(status)
	c.Out.Header().Set("Content-Type", "text/html")
	tmpl.Execute(c.Out, data)
}

func (c *Context) VarInt64(key string) (int64, error) {
	if val, ok := c.Vars[key]; ok {
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, fmt.Errorf("no such var: %s", key)
}

func (c *Context) QueryInt64(key string) (int64, error) {
	query := c.Req.URL.Query()
	return strconv.ParseInt(query.Get(key), 10, 64)
}

func (c *Context) Redirect(url string) {
	if url == "" {
		url = "/"
	}
	http.Redirect(c.Out, c.Req, url, http.StatusFound)
}
