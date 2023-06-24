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
