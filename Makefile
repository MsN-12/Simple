createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen --build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/MsN-12/simpleBank/db/sqlc Store
.PHONY: createdb migrateup migratedown sqlc test server mock