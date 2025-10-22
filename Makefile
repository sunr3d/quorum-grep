all: up build-client

up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

rebuild: clean up

clean:
	docker compose down -v
	rm -f ./mygrep

build-client:
	go build -o mygrep ./cmd/client/main.go

protogen:
	protoc --go_out=proto --go-grpc_out=proto api/grep_service/grep.proto

logs:
	docker compose logs -f

test:
	go test -v ./...

test-comparison: all
	cd test_files && ./test_comparison.sh

fmt:
	go fmt ./...

lint:
	golangci-lint run