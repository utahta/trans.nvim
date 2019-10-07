export GO111MODULE=on

all: build ## Runs a build task.

build: ## Build a trans binary.
	go build -o bin/trans ./src/cmd/trans

manifest: build ## Update a manifest.
	./bin/trans -manifest trans -location ./plugin/trans.vim

