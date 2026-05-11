.PHONY: all

all: test forum hashgen mockgen

test:
	export CGO_ENABLED=1; go test -v ./...

forum:
	export CGO_ENABLED=1; go build -o bin/forum cmd/forum/main.go

hashgen:
	export CGO_ENABLED=1; go build -o bin/hashgen cmd/hashgen/main.go

mockgen:
	export CGO_ENABLED=1; go build -o bin/mockgen cmd/mockgen/main.go

clean:
	rm -rf bin
