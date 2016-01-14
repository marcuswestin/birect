package birect

// import "github.com/marcuswestin/go-ws"



// // pool maintains the set of active connections and can
// // broadcasts messages to the connections.
// type pool struct {
// 	// Registered connections.
// 	connections map[*ws.Conn]bool

// 	// Inbound messages from the connections.
// 	broadcastChan chan []byte

// 	// Register requests from the connections.
// 	registerChan chan *ws.Conn

// 	// Unregister requests from connections.
// 	unregisterChan chan *ws.Conn
// }

// func newPool() *pool {
// 	p := &pool{
// 		registerChan:   make(chan *ws.Conn),
// 		unregisterChan: make(chan *ws.Conn),
// 		connections:    make(map[*ws.Conn]bool),
// 	}
// 	go p._eventLoop()
// 	return p
// }

// // Broadcast will send data to all current connections.
// // A Broadcast will be sent to a snapshot of all connections
// // that are connected when Broadcast is called.
// func (p *pool) Broadcast(data []byte) {
// 	p.broadcastChan <- data
// }

// // Internal
// ///////////

// func (p *pool) _eventLoop() {
// 	for {
// 		select {
// 		case c := <-p.registerChan:
// 			p.connections[c] = true

// 		case c := <-p.unregisterChan:
// 			if _, ok := p.connections[c]; ok {
// 				delete(p.connections, c)
// 			}

// 		case m := <-p.broadcastChan:
// 			p._doBroadcast(m)
// 		}
// 	}
// }

// func (p *pool) _doBroadcast(message interface{}) {
// 	for c := range p.connections {
// 		select {
// 		case c.SendChan() <- message:
// 			// BinaryMessage sent
// 		default:
// 			// Unable to send to conn
// 			delete(p.connections, c)
// 			// close(c.send)
// 		}
// 	}
// }
