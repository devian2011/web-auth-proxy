all: test build

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY
