
all: build

build:
	@GO111MODULE=on go build -o bin/trans ./src/cmd/trans

manifest: build
	@./bin/trans -manifest trans -location ./plugin/trans.vim

