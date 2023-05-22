.PHONY: all migrate-up migrate-down run build

all: build run migrate-up

build:
	docker-compose build

run:
	docker-compose up

migrate-up:
	migrate -path=./shared/db/posgtresql/migrations -database 'postgres://user:password@localhost:5432/newsletter_app?sslmode=disable' up

migrate-down:
	migrate -path=./shared/db/posgtresql/migrations -database 'postgres://user:password@localhost:5432/newsletter_app?sslmode=disable' down