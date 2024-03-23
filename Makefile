.PHONY: full-run

full-run: run

run:
	docker compose up --build -d

stop:
	docker compose down