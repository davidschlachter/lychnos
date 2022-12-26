backend:
	cd src/backend && go build .

backend.watch:
	cd src/backend && go run github.com/cespare/reflex@latest -s -r '\.go$$' go run github.com/joho/godotenv/cmd/godotenv@latest go run .

frontend: frontend.install
	npm run build

frontend.install:
	cd src/frontend && npm install

frontend.watch: frontend.install
	cd src/frontend && BROWSER=none npm run start

dev: backend.watch frontend.watch

clean:
	rm src/backend/backend
	rm src/backend/database.sqlite3
	rm src/frontend/node_modules

test:
	cd src/backend && go test ./...

.PHONY: backend backend.watch frontend frontend.install frontend.watch dev clean test
