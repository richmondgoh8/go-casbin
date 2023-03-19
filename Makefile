.PHONY : run gen mock badge test

run:
	go run cmd/server.go dev
gen:
	go generate ./...
mock:
	mockgen -source=./internal/core/ports/ports.go -destination=./internal/mocks/core/ports/ports.go
badge:
	@which gobadge >/dev/null 2>&1 || go install github.com/AlexBeauchemin/gobadge
	bash test.sh
test:
	go test ./...
lint:
	golangci-lint run
postgres:
	docker-compose -f docker-compose-postgres.yml up -d -V
destroy:
	docker-compose --log-level ERROR -f docker-compose-postgres.yml down --remove-orphans