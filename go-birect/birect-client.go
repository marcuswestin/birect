package birect

import "github.com/marcuswestin/go-ws"

type Client struct {
	JSONReqHandlerMap
	*Conn
}

func Connect(url string) (client *Client, err error) {
	client = &Client{make(JSONReqHandlerMap), nil}
	wsConnChan := make(chan *ws.Conn)
	ws.Connect(url, func(event *ws.Event, conn *ws.Conn) {
		client.Log("Client:", event)
		switch event.Type {
		case ws.Connected:
			wsConnChan <- conn
		case ws.BinaryMessage:
			client.Conn.readAndHandleWireWrapperReader(event)

		default:
			panic("TODO Handle event: " + event.String())
		}
	})
	client.Conn = newConn(<-wsConnChan, client.JSONReqHandlerMap)
	return
}
