GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/gin-gonic/gin gopkg.in/redis.v5 github.com/pr8kerl/redirektor

all: redirektor

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

redirektor: main.go config.go redirektor.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@

windows:
	  gox -os="windows"

.PHONY: $(DEPS) clean

clean:
	rm -f redirektor
