
all: build

build:
	@GO111MODULE=on go build -o bin/trans ./cmd/trans

manifest: build
	@./bin/trans -manifest trans -location ./plugin/trans.vim

