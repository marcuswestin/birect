package birect_test

import (
	"fmt"
	"net/http"

	"github.com/marcuswestin/birect/go-birect"
)

func ExampleUpgradeRequests_server() {
	go http.ListenAndServe("localhost:8087", nil)
	server := birect.UpgradeRequests("/birect/upgrade")

	type EchoParams struct{ Text string }
	type EchoResponse struct{ Text string }
	server.HandleJSONReq("Echo", func(req *birect.JSONReq) (res interface{}, err error) {
		var par EchoParams
		req.ParseParams(&par)
		return EchoResponse{par.Text}, nil
	})
	// Output:
	//
}

func ExampleConnect_client() {
	conn, _ := birect.Connect("http://localhost:8087/birect/upgrade")

	type EchoParams struct{ Text string }
	type EchoResponse struct{ Text string }
	var par = EchoParams{"Hi!"}
	var res EchoResponse
	fmt.Println("Send:", par.Text)
	conn.SendJSONReq("Echo", par, &res)
	fmt.Println("Received:", res.Text)

	// Output:
	// Send: Hi!
	// Received: Hi!
}
