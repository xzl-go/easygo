package core

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	aborted    bool
}

func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Req = r
	c.Path = r.URL.Path
	c.Method = r.Method
	c.Params = make(map[string]string)
	c.StatusCode = 0
	c.aborted = false
}

func (c *Context) Abort() {
	c.aborted = true
}

func (c *Context) IsAborted() bool {
	return c.aborted
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Write([]byte(format))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) BindJSON(obj interface{}) error {
	body, err := io.ReadAll(c.Req.Body)
	if err != nil {
		return err
	}
	defer c.Req.Body.Close()

	return json.Unmarshal(body, obj)
}

func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) GetParam(key string) string {
	return c.Params[key]
}
