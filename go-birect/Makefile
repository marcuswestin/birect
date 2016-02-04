test: ${GO_PROTOS} lint vet
	go test -v .
test-race: test
	go test --race -v .
lint:
	golint .
	test -z "$$(golint .)"
vet:
	go vet .

GO_PROTOS := go-birect/internal/wire/*.go
PROTOEASY := ${GOPATH}/bin/protoeasy
${GO_PROTOS}: proto/*.proto ${PROTOEASY}
	${protoeasy} --go ./proto --out ./go-birect/internal/wire
${PROTOEASY}:
	go get go.pedge.io/protoeasy/cmd/protoeasy
