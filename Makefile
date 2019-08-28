export GO111MODULE = on

.PHONY: test test-cover

# for dev
dev:
	fresh

# build protoc
protoc:
	protoc -I pb/ pb/optim.proto --go_out=plugins=grpc:pb

# for test
test: export GO_ENV=test
test:
	go test -race -cover ./...

test-cover: export GO_ENV=test
test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

build-linux:
	GOOS=linux go build -o tiny-linux
