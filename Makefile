GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/boltdb/bolt google.golang.org/grpc github.com/mitchellh/cli gopkg.in/redis.v4



all: redirektor csvimporter

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

redirektor: cmd/redirektor/main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch $@

csvimporter: cmd/csvimporter/main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch $@


.PHONY: $(DEPS) clean

clean:
	rm -f bin/*

