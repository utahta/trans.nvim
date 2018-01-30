
all: build

build:
	@GOPATH=$(CURDIR) go build -o bin/trans ./cmd/trans

manifest: build
	@./bin/trans -manifest trans -location ./plugin/trans.vim

