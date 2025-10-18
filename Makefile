up:
	docker compose up -d --build

down:
	docker compose down

restart: down up

rebuild: clean up

clean:
	docker compose down -v

protogen:
	protoc --go_out=proto --go-grpc_out=proto api/grep_service/grep.proto

logs:
	docker compose logs -f app

test: mocks
	go test -v ./...

mocks:
	go generate ./...

fmt:
	go fmt ./...