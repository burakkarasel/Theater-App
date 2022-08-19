sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run cmd/web/main.go
PHONY: sqlc test down server