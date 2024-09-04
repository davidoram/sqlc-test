.PHONY: all

test:
	go test ./... -v

clean:
	mkdir -p bin
	rm -rf bin

sqlc-generate:
	docker run --network demo_network -e POSTGRES_PASSWORD=postgres --rm -v $(shell pwd):/src -w /src sqlc/sqlc:1.26.0 generate

build: sqlc-generate
	mkdir -p bin
	go build -o bin/sqlc-test main.go 


startup:
	docker network create demo_network
	docker run --network demo_network --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:16-alpine
	./scripts/wait_for_postgres.sh
	PGPASSWORD=postgres psql -h localhost -U postgres -c 'create database sqlc_test;'
	sql-migrate up -env="development"

shutdown:
	docker stop postgres || true
	docker rm postgres || true
	docker network rm demo_network || true


bootstrap: shutdown startup build
