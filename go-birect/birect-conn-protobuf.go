package birect

import (
	"github.com/golang/protobuf/proto"
	"github.com/marcuswestin/birect/go-birect/internal/wire"
	"github.com/marcuswestin/go-errs"
)

type ProtoReqHandler func(req *ProtoReq) (resValue proto.Message, err error)
type ProtoReqHandlerMap map[string]ProtoReqHandler
type ProtoReq struct {
	*Conn
	data []byte
}

func (m ProtoReqHandlerMap) HandleProtoReq(reqName string, handler ProtoReqHandler) {
	m[reqName] = handler
}

func (c *Conn) SendProtoReq(name string, paramsObj proto.Message, resValPtr proto.Message) (err error) {
	data, err := proto.Marshal(paramsObj)
	if err != nil {
		return
	}
	reqId := c.nextReqId()
	wireReq := &wire.Request{Type: wire.DataType_Proto, Name: name, ReqId: uint32(reqId), Data: data}
	return c.sendRequestAndWaitForResponse(reqId, wireReq, resValPtr)
}

func (p *ProtoReq) ParseParams(valuePtr proto.Message) {
	err := proto.Unmarshal(p.data, valuePtr)
	if err != nil {
		panic(errs.Wrap(err, nil, "Unable to parse params"))
	}
}

// Internal
///////////

func (c *Conn) handleProtoWireReq(wireReq *wire.Request) {
	// Find handler
	handler, exists := c.ProtoReqHandlerMap[wireReq.Name]
	if !exists {
		c.sendErrorResponse(wireReq, errs.New(nil, "Missing request handler"))
		return
	}
	// Execute handler
	resVal, err := handler(&ProtoReq{c, wireReq.Data})
	if err != nil {
		c.sendErrorResponse(wireReq, err)
		return
	}
	// Send response
	c.sendResponse(wireReq, &protoRes{resVal})
}

type protoRes struct {
	resValPtr proto.Message
}

func (j *protoRes) encode() ([]byte, error) {
	return proto.Marshal(j.resValPtr)
}
func (j *protoRes) dataType() wire.DataType {
	return wire.DataType_Proto
}
