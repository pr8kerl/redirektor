GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/gin-gonic/gin gopkg.in/gcfg.v1 gopkg.in/redis.v5 github.com/asaskevich/govalidator

all: redirektor

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^
	npm install

redirektor: main.go config.go redirektor.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@

frontend:
		npm run build-dev

.PHONY: $(DEPS) clean

clean:
	rm -f redirektor
