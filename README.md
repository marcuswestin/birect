Birect
======

Efficient, scalable and language-agnostic bidirectional realtime messaging with request/response support.


Setup
-----

OS X:

	# protoc:
	./scripts/install-protoc.sh


Usage: Go
---------

#### Server

See http://godoc.org/pkg/github.com/marcuswestin/birect/go-birect/#example_UpgradeRequests_server

#### Client

See http://godoc.org/pkg/github.com/marcuswestin/birect/go-birect/#example_Connect_client


Usage: JS
---------

Coming.


TODO
----

- [ ] Implement protobuf-based Conn
- [ ] De-duplicate json/protobuf Conn code
- [ ] Consider implementing text-based Conn
- [ ] Tests for protobuf code
- [ ] Tests for error handling
- [ ] Pluggable Client/Server logging
- [ ] Regular messages (wire.Message)
- [ ] Tests for regular messaging
- [ ] Error encoding, decoding and receiving
- [ ] Specific panic handling (params encode/decode, etc)
- [ ] JS conn + Connect
- [ ] JS Attachment sending
- [ ] Go Attachment recieving
- [ ] Consider all names
	- [ ] Req/Res vs Request/Response
	- [ ] Go, JS, Objc, Java
- [X] Fix `make test-race`
