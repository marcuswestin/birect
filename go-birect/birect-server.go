package birect

import (
	"log"
	"sync"

	"github.com/marcuswestin/go-ws"
)

type Server struct {
	JSONReqHandlerMap
	ProtoReqHandlerMap
	connByWSConnMutex *sync.Mutex
	connByWSConn      map[*ws.Conn]*Conn
}

func UpgradeRequests(pattern string) (server *Server) {
	server = &Server{
		make(JSONReqHandlerMap),
		make(ProtoReqHandlerMap),
		&sync.Mutex{},
		make(map[*ws.Conn]*Conn, 10000),
	}
	ws.UpgradeRequests(pattern, func(event *ws.Event, wsConn *ws.Conn) {
		log.Println("Server:", event)
		switch event.Type {
		case ws.Connected:
			server.registerConn(wsConn)
		case ws.BinaryMessage:
			server.connByWSConn[wsConn].readAndHandleWireWrapperReader(event)
		case ws.Disconnected:
			server.deregisterConn(wsConn)
		default:
			panic("birect.Server unknown event: " + event.String())
		}
	})
	return server
}

func (s *Server) Log(args ...interface{}) {
	log.Println(args...)
}

// Internal
///////////

func (s *Server) registerConn(wsConn *ws.Conn) {
	s.connByWSConnMutex.Lock()
	defer s.connByWSConnMutex.Unlock()
	s.connByWSConn[wsConn] = newConn(wsConn, s.JSONReqHandlerMap, s.ProtoReqHandlerMap)
}
func (s *Server) deregisterConn(wsConn *ws.Conn) {
	s.connByWSConnMutex.Lock()
	defer s.connByWSConnMutex.Unlock()
	delete(s.connByWSConn, wsConn)
}
