.PHONY: build build-backend build-frontend test test-backend test-frontend run-backend run-frontend lint lint-backend lint-frontend

build: build-backend build-frontend

build-backend:
	cd backend && go build ./...

build-frontend:
	cd frontend && npm run build

test: test-backend test-frontend

test-backend:
	cd backend && go test ./... -cover

test-frontend:
	cd frontend && npm run test

run-backend:
	cd backend && go run ./cmd/server

run-frontend:
	cd frontend && npm run dev

lint: lint-backend lint-frontend

lint-backend:
	cd backend && go vet ./...

lint-frontend:
	cd frontend && npm run lint
