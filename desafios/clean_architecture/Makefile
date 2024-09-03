createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up

migratedown:
	migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose down

sqlcgenerate:
	sqlc generate

db:
	docker-compose exec mysql mysql -uroot -proot orders

gqlgen:
	go run github.com/99designs/gqlgen generate

protoc:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

evans:
	evans -r repl

run:
	cd cmd/ordersystem && go run main.go wire_gen.go


.PHONY: migrate migratedown createmigration sqlcgenerate db gqlgen protoc evans run