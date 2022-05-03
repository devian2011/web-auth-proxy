all: test build

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: make_certs
make_certs:
	openssl req -x509 -nodes -newkey rsa:2048 -keyout $(PWD)/config/certs/proxy.rsa.key -out $(PWD)/config/certs/proxy.rsa.crt -days 3650
	openssl req -x509 -nodes -newkey rsa:2048 -keyout $(PWD)/config/certs/admin.rsa.key -out $(PWD)/config/certs/admin.rsa.crt -days 3650
