sqlc:
	sqlc generate
test:
	./test.sh
server:
	./start.sh
mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/burakkarasel/Theatre-API/internal/db/sqlc Store
down:
	./down.sh
	
PHONY: sqlc test down server mock