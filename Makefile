setup: setup-go setup-js setup-dev
test: test-go test-js
clean-test: test-clean-go test-clean-js

# go
####
setup-go:
	cd go-birect && glide install
test-go:
	cd go-birect && make test

# js
####
setup-js:
	cd js-birect && npm install
test-js:
	cd js-birect && make test

# clean tests (circleci)
########################
test-clean-go:
	rm -rf ${GOPATH}/src/github.com/marcuswestin/birect
	go get github.com/marcuswestin/birect/go-birect
	go test --race -v github.com/marcuswestin/birect/go-birect
test-clean-js:
	rm -rf ./js-birect/node_modules
	cd ./js-birect && npm install
	cd ./js-birect && node birect_tests.js

# dev
#####
setup-dev:
	bash scripts/install-protoc.sh
