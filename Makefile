test: test-go test-js

GO_PROTOS = go-birect/internal/wire/*.go

test-go: ${GO_PROTOS}
	cd go-birect && make test

test-js:
	cd js-birect && make test

${GO_PROTOS}: proto/*.proto
	protoeasy --go ./proto --out ./go-birect/internal/wire
	
