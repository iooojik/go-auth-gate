run_local:
	docker-compose up --build --remove-orphans -d

test:
	go test -race ./...

mockery:
	mockery --all

lint:
	golangci-lint run -c .golangci.yml