include config.env
export $(shell sed 's/=.*//' config.env)

build:
	docker-compose build rest-wallets

run:
	docker-compose up rest-wallets

migration-create:
	migrate create -ext sql -dir ./schema -seq $(name)

migrate:
	migrate -path ./schema -database 'postgres://postgres:$(POSTGRES_PASSWORD)@0.0.0.0:5436/postgres?sslmode=disable' up

test:
	go test -v ./...