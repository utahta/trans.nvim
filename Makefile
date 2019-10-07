export GO111MODULE=on

all: build

build:
	go build -o bin/trans ./src/cmd/trans

manifest: build
	./bin/trans -manifest trans -location ./plugin/trans.vim

