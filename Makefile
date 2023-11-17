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
proto:
	rm pb/*.go ;rm doc/swagger/*.swagger.json;protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank proto/*.proto; statik -src./doc/swagger -dest=./doc
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
.PHONY: createdb migrateup migratedown sqlc test server mock proto redis new_migration db_schema