test: test-go

GO_PROTOS = go-birect/internal/wire/*.go

test-go: ${GO_PROTOS}
	# go test -v github.com/marcuswestin/go-ws
	go test -v github.com/marcuswestin/birect/go-birect

${GO_PROTOS}: proto/*.proto
	protoeasy --go ./proto --out ./go-birect/internal/wire
