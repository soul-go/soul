package context

import (
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/proto"
)

const (
	typeJson  = "application/json"
	typePlain = "text/plain"
	typeProto = "application/protobuf"
)

type Handler func(*Context)

type C interface {
	Publish(*R)
	Request() *R
}

type Context struct {
	C       C
	M       M
	TraceID string
	SpanID  string
	Request *R
	Headers M
}

type R struct {
	Body    []byte
	Headers M
	Type    string
}

type M map[string]interface{}

func NewContext(c C) *Context {
	context := &Context{C: c}
	context.setRequest()
	return context
}

func (c *Context) Bind(v interface{}) error {
	r := c.Request
	switch r.Type {
	case typeJson:
		return json.Unmarshal(r.Body, v)
	case typeProto:
		return proto.Unmarshal(r.Body, v.(proto.Message))
	default:
		return errors.New("")
	}
}

func (c *Context) JSON(v interface{}) {
	body, _ := json.Marshal(v)
	c.publish(body, typeJson)
}

func (c *Context) Proto(pb proto.Message) {
	body, _ := proto.Marshal(pb)
	c.publish(body, typeProto)
}

func (c *Context) String(text string) {
	body := []byte(text)
	c.publish(body, typePlain)
}

func (c *Context) publish(body []byte, t string) {
	r := &R{
		Body:    body,
		Headers: c.Headers,
		Type:    t,
	}
	c.C.Publish(r)
}

func (c *Context) setRequest() {
	r := c.C.Request()
	c.Request = r
}