package birect

import (
	"encoding/json"

	"github.com/marcuswestin/birect/go-birect/internal/wire"
	"github.com/marcuswestin/go-errs"
)

// JSONReqHandler functions get called on every json request
type JSONReqHandler func(req *JSONReq) (resValue interface{}, err error)

// SendJSONReq sends a request for the JSONReqHandler with the given `name`, along with the
// given paramsObj. When the server responds, SendJSONReq will parse the response into resValPtr.
func (c *Conn) SendJSONReq(name string, paramsObj interface{}, resValPtr interface{}) (err error) {
	data, err := json.Marshal(paramsObj)
	if err != nil {
		return
	}
	reqID := c.nextReqID()
	wireReq := &wire.Request{Type: wire.DataType_JSON, Name: name, ReqId: uint32(reqID), Data: data}
	return c.sendRequestAndWaitForResponse(reqID, wireReq, resValPtr)
}

// JSONReq wraps a request sent via SendJSONReq. Use ParseParams to access the JSON values.
type JSONReq struct {
	*Conn
	data []byte
}

// ParseParams parses the JSONReq values into the given valuePtr.
// valuePtr should be a pointer to a struct that can be JSON-parsed, e.g
//
// 	type params struct { Foo string }
// 	var p params
// 	jsonReq.ParseParams(&p)
func (j *JSONReq) ParseParams(valuePtr interface{}) {
	err := json.Unmarshal(j.data, valuePtr)
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to parse params"))
	}
}

// Internal
///////////

type jsonReqHandlerMap map[string]JSONReqHandler

func (m jsonReqHandlerMap) HandleJSONReq(reqName string, handler JSONReqHandler) {
	m[reqName] = handler
}

func (c *Conn) handleJSONWireReq(wireReq *wire.Request) {
	// Find handler
	handler, exists := c.jsonReqHandlerMap[wireReq.Name]
	if !exists {
		c.sendErrorResponse(wireReq, errs.New(nil, "Missing request handler"))
		return
	}
	// Execute handler
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errs.New(errs.Info{"Recovery": r})
			}
			c.Log("Error while handling request", wireReq.Name, err)
			c.sendErrorResponse(wireReq, errs.Wrap(err, errs.Info{"Name": wireReq.Name, "Data": wireReq.Data}))
		}
	}()
	resVal, err := handler(&JSONReq{c, wireReq.Data})
	if err != nil {
		c.sendErrorResponse(wireReq, errs.Wrap(err, errs.Info{"HandlerName": wireReq.Name}))
		return
	}
	// Send response
	c.sendResponse(wireReq, &jsonRes{resVal})
}

type jsonRes struct {
	resValPtr interface{}
}

func (j *jsonRes) encode() ([]byte, error) {
	return json.Marshal(j.resValPtr)
}
func (j *jsonRes) dataType() wire.DataType {
	return wire.DataType_JSON
}
