sqlc:
	sqlc generate
test:
	./test.sh
server:
	./start.sh
mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/burakkarasel/Theatre-API/internal/db/sqlc Store
	
PHONY: sqlc test down server mock