postgres:
	docker run --name financetrackerdb -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:12-alpine

createdb:
	docker exec -it financetrackerdb createdb --username=root --owner=root finance_tracker
	
dropdb:
	docker exec -it financetrackerdb dropdb finance_tracker

migrateup:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/finance_tracker?sslmode=disable" -verbose up

migrateup1:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/finance_tracker?sslmode=disable" -verbose up 1

migratedown:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/finance_tracker?sslmode=disable" -verbose down

migratedown1:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/finance_tracker?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

templ:
	templ generate

test:
	go test -v -cover ./...

air:
	air

server:
	go run ./cmd/server/main.go

faking:
	go run ./cmd/faking/main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/moth13/finance_tracker/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown server mock migratedown1 migrateup1 air templ faking