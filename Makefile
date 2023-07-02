ENV ?="development"
PORT ?= 8080 

run.dev:
	go run ./cmd/api/ -d ${DB} -e ${ENV} -p ${PORT}

build.bin:
	rm -rf bin
	mkdir bin && go build -o bin ./cmd/api/

run.help:
	go run ./cmd/api/ --help

db.shell:
	docker exec -it postgresql_db psql --host=localhost --dbname=greenlight --username=postgres


db.up:
	docker compose up -d && docker start postgresql_db && export DB=postgres://postgres:123123@localhost/greenlight?sslmode=disable && export GREENLIGHT_DB_CONNECTION=postgres://postgres:123123@localhost/greenlight?sslmode=disable

db.migrate.up:
	migrate -path ./migrations -database "${GREENLIGHT_DB_CONNECTION}" up

db.migrate.down:
	migrate -path ./migrations -database "${GREENLIGHT_DB_CONNECTION}" down

db.migrate.to:
	migrate -path ./migrations -database "${GREENLIGHT_DB_CONNECTION}" goto ${VERSION}

db.migration.new:
	migrate create -seq -ext .sql -dir ./migrations ${NAME}