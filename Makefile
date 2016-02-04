setup: setup-go setup-js
test: test-go test-js

# go
####
#GO_PROTOS := go-birect/internal/wire/*.go
#PROTOEASY := ${GOPATH}/bin/protoeasy
setup-go:
	cd go-birect && glide install
test-go:# ${GO_PROTOS}
	cd go-birect && make test
#${GO_PROTOS}: proto/*.proto ${PROTOEASY}
#	${protoeasy} --go ./proto --out ./go-birect/internal/wire
#${PROTOEASY}:
#	go get go.pedge.io/protoeasy/cmd/protoeasy

# js
####
setup-js:
	cd js-birect && npm install
test-js:
	cd js-birect && make test

# dev
#####
setup-dev:
	bash scripts/install-protoc.sh
