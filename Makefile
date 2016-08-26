GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/boltdb/bolt



all: rewriter importer

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

rewriter: rewriter.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@

importer: importer.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@


.PHONY: $(DEPS) clean

clean:
	rm -f rewriter importer

