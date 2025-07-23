all: run

fastrun:
	go run src/main.go --config config.yaml

run:
	docker compose up --build -d

stop:
	docker compose down -v

