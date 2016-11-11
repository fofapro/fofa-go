.PHONY: build install fmt vet test doc

GOPATH := ${GOPATH}
export GOPATH

default: build

build: vet
	govendor build +local
install: vet
	govendor install +local
fmt:
	go fmt ./...

test:
	govendor test  ./fofa

vet: 
	go vet ./...

doc:
	godoc -http=:6060 -index
