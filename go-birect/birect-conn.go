package birect

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
	"github.com/marcuswestin/birect/go-birect/internal/wire"
	"github.com/marcuswestin/go-errs"
	"github.com/marcuswestin/go-ws"
)

type reqId uint32
type resChan chan *wire.Response
type Conn struct {
	wsConn    *ws.Conn
	lastReqId reqId
	resChans  map[reqId]resChan
	JSONReqHandlerMap
}

func newConn(wsConn *ws.Conn, jsonHandlers JSONReqHandlerMap) *Conn {
	return &Conn{
		wsConn:            wsConn,
		resChans:          make(map[reqId]resChan, 1),
		JSONReqHandlerMap: jsonHandlers,
	}
}

var _ connInterface = &Conn{}

type connInterface interface {
	SendJSONReq(name string, paramsObj interface{}, resValPtr interface{}) (err error)
	// SendProtoReq(name ReqName, req ProtoReq, resValPtr proto.ProtoMessage) (err error)
	// Close()
}

type request interface {
	// Sending side
	encode() ([]byte, error)
	// Receiving side
	ParseParams(valPrt interface{})
}
type response interface {
	// Sending side
	dataType() wire.DataType
	encode() ([]byte, error)
	// Receiving side
}

// JSON Req/Res
///////////////

type JSONReqHandler func(req *JSONReq) (resValue interface{}, err error)
type JSONReqHandlerMap map[string]JSONReqHandler
type JSONReq struct {
	*Conn
	data []byte
}

func (j *JSONReq) ParseParams(valPtr interface{}) {
	err := json.Unmarshal(j.data, valPtr)
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to parse params"))
	}
}

type JSONRes struct{ resValPtr interface{} }

func (j *JSONRes) encode() ([]byte, error) {
	return json.Marshal(j.resValPtr)
}
func (j *JSONRes) dataType() wire.DataType {
	return wire.DataType_JSON
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

// Proto Req/Res
////////////////

// TODO

// Internal - Outgoing wrappers
///////////////////////////////

func (c *Conn) Log(argv ...interface{}) {
	log.Println(argv...)
}

func (c *Conn) sendRequestAndWaitForResponse(reqId reqId, wireReq *wire.Request, resValPtr interface{}) (err error) {
	c.resChans[reqId] = make(resChan)
	defer delete(c.resChans, reqId)

	c.Log("REQ", wireReq.Name, "ReqId:", reqId, "len:", len(wireReq.Data))
	err = c.sendWrapper(&wire.Wrapper{&wire.Wrapper_Request{wireReq}})
	if err != nil {
		return
	}

	// Receive response
	wireRes := <-c.resChans[reqId]
	c.Log("RCV", wireReq.Name, reqId, wireRes.Type, len(wireRes.Data))
	switch wireRes.Type {
	case wire.DataType_JSON:
		return json.Unmarshal(wireRes.Data, resValPtr)
	case wire.DataType_Proto:
		return proto.Unmarshal(wireRes.Data, resValPtr.(proto.Message))
	default:
		return errors.New("Bad response wire type: " + wireRes.Type.String())
	}
	return
}

func (c *Conn) sendResponse(wireReq *wire.Request, response response) {
	wireRes := &wire.Response{ReqId: wireReq.ReqId}
	data, err := response.encode()
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to encode response"))
	}
	wireRes.Type = response.dataType()
	wireRes.Data = data
	err = c.sendWrapper(&wire.Wrapper{&wire.Wrapper_Response{wireRes}})
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to send response"))
	}
}
func (c *Conn) sendErrorResponse(wireReq *wire.Request, err error) {
	wireRes := &wire.Response{
		ReqId:   wireReq.ReqId,
		IsError: true,
		// Type: ???,
		// Data: []byte(responseError.Error()) ???,
	}
	c.sendWrapper(&wire.Wrapper{&wire.Wrapper_Response{wireRes}})
}
func (c *Conn) nextReqId() reqId {
	rawReqId := atomic.AddUint32((*uint32)(&c.lastReqId), 1)
	return reqId(rawReqId)
}
func (c *Conn) sendWrapper(wrapper *wire.Wrapper) (err error) {
	wireData, err := proto.Marshal(wrapper)
	if err != nil {
		return
	}
	c.Log("SND Wrapper len:", len(wireData), wrapper)
	return c.wsConn.SendBinary(wireData)
}

// Internal - incoming wrappers
///////////////////////////////

func (c *Conn) readAndHandleWireWrapperReader(reader io.Reader) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to read wrapper"))
	}
	c.readAndHandleWireWrapper(data)
}
func (c *Conn) readAndHandleWireWrapper(data []byte) {
	if len(data) == 0 {
		panic(errs.New(nil, "Empty data"))
	}

	var wireWrapper wire.Wrapper
	err := proto.Unmarshal(data, &wireWrapper)
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to decode wire wrapper"))
	}

	c.Log("readAndHandleWireWrapper", wireWrapper.Content)
	switch content := wireWrapper.Content.(type) {
	case *wire.Wrapper_Message:
		c.handleMessage(content.Message)
	case *wire.Wrapper_Request:
		c.handleRequest(content.Request)
	case *wire.Wrapper_Response:
		c.handleResponse(content.Response)
	default:
		panic(errs.New(errs.Info{"Wrapper": wireWrapper}, "Unknown wire wrapper content type"))
	}
}

func (c *Conn) handleMessage(msg *wire.Message) {
	panic(errs.New(nil, "TODO: handleMessage"))
}
func (c *Conn) handleRequest(wireReq *wire.Request) {
	c.Log("HANDLE REQ", wireReq)
	switch wireReq.Type {
	case wire.DataType_JSON:
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
		c.sendResponse(wireReq, &JSONRes{resVal})

	default:
		panic(errs.New(errs.Info{"Type": wireReq.Type}, "Bad wireReq.Type"))
	}
}
func (c *Conn) handleResponse(wireRes *wire.Response) {
	c.Log("HANDLE RES", wireRes)
	c.resChans[reqId(wireRes.ReqId)] <- wireRes
}
