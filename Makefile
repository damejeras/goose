migration:
	@mkdir -p db/migrations
	@read -p "Enter migration name: " name; \
	docker run --rm --user $(shell id -u):$(shell id -g) -v $(PWD)/db/migrations:/migrations --network host migrate/migrate create -ext sql -dir /migrations $$name

queries:
	docker run --rm --user $(shell id -u):$(shell id -g) -v $(PWD):/src -w /src sqlc/sqlc generate

html:
	docker run --rm --user $(shell id -u):$(shell id -g) -v $(PWD):/src -w /src ghcr.io/a-h/templ:latest generate

lint:
	docker run --rm -t -v /tmp/golangci:/root/.cache:rw -v $(PWD):/src -w /src golangci/golangci-lint:v2.2.2 golangci-lint run

run:
	@CLIENT_ID=$$(op read op://Projects/rent-a-flat/google_oauth_client_id) && \
	CLIENT_SECRET=$$(op read op://Projects/rent-a-flat/google_oauth_secret) && \
	wgo -file=.go -file=.templ -xfile=_templ.go templ generate :: go run cmd/main.go --google.client-id=$$CLIENT_ID --google.client-secret=$$CLIENT_SECRET


.PHONY: migration queries html lint run
