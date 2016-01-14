package birect_test


// Create package birect-hub, which has a redis adaptor and an in-memory adaptor

func forwardMessagesFromHub(server *birect.Server) {
	hub := redis.Connect()
	go func() {
		for hubMsg := range hub {
			hubMsg := HubMessage.Unmarshal(hubMsg.Bytes)
			if hubMsg.Type == UserEvent {
				server.Dispatch(hubMsg.UserId, hubMsg.Bytes)
			} else if hubMsg.Type == Internal {
				// E.g load new config, etc.
			}
		}
	}()
}

func dispatchUserEvent(userId UserId, event Event) {
	hubMsg := HubMessage{userId, event}
	hub.Publish(userId)
}

func ExampleConnect_client() {
	client, err := birect.Connect("/birect/ws")
	go readEvents(client)
}

// type AuthParams struct {
// 	Username string
// 	Password string
// }
// type EchoParams struct {
// 	Text string
// }
// type EchoResponse struct {
// 	Text string
// }

func ExampleListen_server() {
	birect.Listen("/birect/ws", func(req birect.Request) (res *birect.Response) {
		switch (req.Name) {
			case "Auth":
				var params AuthParams
				req.MustUnmarshal(&params)
				if params.Password != "open sesame" {
					return errs.New("Wrong password")
				}
				req.Conn.Subscribe(params.Username)
			
			case "Echo":
				return EchoResponse{req.Data()}
				
			case "SendTo":
				var params struct{
					To string
					Text string
				}
				req.UnmarshalJSON(&params)
				birect.SendTo(params.To, params.Text)
				return nil
			
			case "Upload":
				info, reader, err := req.NextAttachment()
				birect.Broadcast([]byte("New upload"))
				allData, err := iotuils.ReadAll(reader)
				
			case "Broadcast":
				birect.Broadcast(req.Data())
				return nil
				
			default:
				return errs.New("Unkown request")
		}
	})
	
	
	conn, err := birect.Connect("/birect/ws")
	if err != nil {
		panic(err)
	}
	
	conn.RequestJSON("Echo", map[string]string{"Text":"Hi!"})
	
	server := birect.Listen("/birect/ws")
	// ...
	var userId = 3716823
	server.Subscribe(conn, userId)
	// ..
	// Send messages
	server.Send(userId, []byte("Hi"))
	server.SendMulti()
	
	serverSend := func() {	
		const msg = "Hi!"
		data, err := conn.Request([]byte(msg))
		fmt.Println("Request:", data, err)
		
		var jsonRes EchoResponse
		err = conn.RequestJSON(req, EchoRequest{msg}, &jsonRes)
		fmt.Println("RequestJSON:", jsonRes, err)

		var protoRes EchoResponse
		err = conn.RequestJSON(req, EchoRequest{msg}, &protoRes)
		fmt.Println("RequestProto:", protoRes, err)
	}
	
	// Output:
	// Request: Hi! nil
	// RequestJSON: {"Text":"Hi!"} nil
	// RequestProto: {"Text":"Hi!"} nil
}

// Implements proto.Message
type ProtoRequest struct {
	Text string
}
// Implements proto.Message
type EchoResponse struct {
	Text string
}

func readEvents(client *birect.Client) {
	for event := range client.EventChan {
		
	}
}

func assert(ok, msg ...string) {
	if !ok {
		panic(fmt.Sprint(msg...))
	}
}

func (c *Conn) Request(req proto.Message, res proto.Message)