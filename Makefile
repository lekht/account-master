all: run

run:
	docker compose up --build -d

stop:
	docker compose down -v
