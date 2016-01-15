package birect

import (
	"encoding/json"

	"github.com/marcuswestin/birect/go-birect/internal/wire"
	"github.com/marcuswestin/go-errs"
)

type JSONReqHandler func(req *JSONReq) (resValue interface{}, err error)
type JSONReqHandlerMap map[string]JSONReqHandler
type JSONReq struct {
	*Conn
	data []byte
}

func (m JSONReqHandlerMap) HandleJSONReq(reqName string, handler JSONReqHandler) {
	m[reqName] = handler
}

func (c *Conn) SendJSONReq(name string, paramsObj interface{}, resValPtr interface{}) (err error) {
	data, err := json.Marshal(paramsObj)
	if err != nil {
		return
	}
	reqId := c.nextReqId()
	wireReq := &wire.Request{Type: wire.DataType_JSON, Name: name, ReqId: uint32(reqId), Data: data}
	return c.sendRequestAndWaitForResponse(reqId, wireReq, resValPtr)
}

func (j *JSONReq) ParseParams(valuePtr interface{}) {
	err := json.Unmarshal(j.data, valuePtr)
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to parse params"))
	}
}

// Internal
///////////

func (c *Conn) handleJSONWireReq(wireReq *wire.Request) {
	// Find handler
	handler, exists := c.JSONReqHandlerMap[wireReq.Name]
	if !exists {
		c.sendErrorResponse(wireReq, errs.New(nil, "Missing request handler"))
		return
	}
	// Execute handler
	resVal, err := handler(&JSONReq{c, wireReq.Data})
	if err != nil {
		c.sendErrorResponse(wireReq, err)
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
