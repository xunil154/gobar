
all: gobar gobar_core

clean:
	rm -rf bin/*

gobar: clean goget
	go build -o bin/gobar client.go

gobar_core: clean goget
	go build -o bin/core core.go

goget:
	cd client && go get || cd -
	cd core && go get || cd -

test: clean goget gobar
	go test github.com/xunil154/gobar/ui

run: gobar
	./bin/gobar

run_core: gobar_core
	./bin/core
