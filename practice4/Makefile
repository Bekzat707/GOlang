.PHONY: run stop clean

run:
	docker compose --env-file .env up --build

stop:
	docker compose down

clean:
	docker compose down -v
