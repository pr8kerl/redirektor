GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/boltdb/bolt google.golang.org/grpc github.com/mitchellh/cli gopkg.in/redis.v4 gopkg.in/gcfg.v1 github.com/cfdrake/go-gdbm github.com/jsimonetti/berkeleydb



all: redirektor 

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

redirektor: cmd/redirektor-redis/main.go cmd/redirektor-redis/config.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch $@

redirektor-boltdb: cmd/redirektor-boltdb/main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch bin/$@

csv2boltdb: cmd/csv2boltdb/main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch bin/$@

csv2gdbm: cmd/csv2gdbm/main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch bin/$@

csv2dbm: cmd/csv2dbm/main.go cmd/csv2dbm/import.go cmd/csv2dbm/export.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o bin/$@ -v $^
		touch bin/$@


.PHONY: $(DEPS) clean

clean:
	rm -f bin/*

