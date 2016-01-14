package birect_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/marcuswestin/birect/go-birect"
)

func TestSetupServer(t *testing.T) {
	var addr = "localhost:8087"
	var path = "/ws/echo"

	type EchoParams struct{ Text string }
	type EchoResponse struct{ Text string }

	// Server side
	server := birect.UpgradeRequests(path)
	server.HandleJSONReq("Echo", func(req *birect.JSONReq) (res interface{}, err error) {
		var par EchoParams
		req.ParseParams(&par)

		res = EchoResponse{par.Text}
		return
	})
	go http.ListenAndServe(addr, nil)

	// Client side
	conn, err := birect.Connect("http://" + addr + path)
	assert(t, err == nil)
	var par = EchoParams{"Hi!"}
	var res EchoResponse
	err = conn.SendJSONReq("Echo", par, &res)
	assert(t, err == nil)
	assert(t, res.Text == par.Text)
}

// Misc
///////

func assert(t *testing.T, ok bool, msg ...interface{}) {
	if !ok {
		// t.Fatal("assert failed", msg)
		log.Panic(msg...)
	}
}
