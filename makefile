include .env
export

run:
	go run ./cmd/api/main.go

test:
	go clean -testcache
	go test -p=1 ./internal/...

swag:
	swag init -g ./internal/transport/http/http.go