run_local:
	docker-compose up --build --remove-orphans -d

test:
	go test -race ./...
