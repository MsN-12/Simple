DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down
sqlc:
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen --build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/MsN-12/simpleBank/db/sqlc Store
	mockgen --build_flags=--mod=mod -package mockwk -destination worker/mock/distributor.go github.com/MsN-12/simpleBank/worker TaskDistributor
proto:
	rm pb/*.go ;rm doc/swagger/*.swagger.json;protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank proto/*.proto;statik -src./doc/swagger -dest=./doc
redis:
	docker run --name redis --network bank-network -p 6379:6379 -d redis:7-alpine
new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
image:
	docker run --name simplebank -p 8080:8080 -e DB_SOURCE"postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" -e REDIS_ADDRESS="redis:6379" simplebank:latest
.PHONY: createdb migrateup migratedown sqlc test server mock proto redis new_migration db_schema postgres image