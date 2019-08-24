export GO111MODULE = on

.PHONY: test test-cover

# for dev
dev:
	fresh

# for test
test: export GO_ENV=test
test:
	go test -race -cover -tags="brotli" ./...

test-cover: export GO_ENV=test
test-cover:
	go test -race -tags="brotli" -coverprofile=test.out ./... && go tool cover --html=test.out

build-linux:
	packr2
	GOOS=linux go build -tags="brotli" -o tiny-linux

clean:
	packr2 clean
