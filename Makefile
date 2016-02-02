setup:
	bash scripts/install-protoc.sh
	cd go-birect && glide install
	cd js-birect && npm install

test: test-go test-js

GO_PROTOS = go-birect/internal/wire/*.go
PROTOEASY = ${GOPATH}/bin/protoeasy

test-go: ${GO_PROTOS}
	cd go-birect && make test

test-js:
	cd js-birect && make test

${GO_PROTOS}: proto/*.proto ${PROTOEASY}
	${protoeasy} --go ./proto --out ./go-birect/internal/wire

${PROTOEASY}:
	go get go.pedge.io/protoeasy/cmd/protoeasy
