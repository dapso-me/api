include .env
export

up:
	docker compose -f ./deployment/docker-compose-dev.yml build
	docker compose -f ./deployment/docker-compose-dev.yml up -d

down:
	docker compose -f ./deployment/docker-compose-dev.yml down
	

test:
	go clean -testcache
	go test -p=1 ./internal/...

swag:
	swag init -g ./internal/transport/http/http.go