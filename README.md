# Budget annually, report monthly

Use [Firefly III](https://github.com/firefly-iii/firefly-iii) to store transactions, but implement my own budgeting system. Goal: have a single app for recording and reporting, instead of using Firefly and an Excel spreadsheet. I budget a yearly amount for each category, and want a simple app to record transactions and track my budget, while still being able to use advanced features in `firefly-iii` if needed. For more details, see `DESIGN.md`.

## Screenshots

<img alt="Category summary" src="docs/category-summary.png" width="375"/> <img alt="Transactions list" src="docs/transactions-list.png" width="375" /> <img alt="New transaction" src="docs/new-transaction.png" width="375" />

## Local development

Prerequisites: working installations of Go and NodeJS.

Make sure you have a working installation of `firefly-iii`. In `src/backend`, copy `.env.sample` to `.env` and update the `firefly-iii` API key (in Firefly: Options > Profile > OAuth > Personal Access Tokens) and URL. The database connection string is optional (if unspecified, a SQLite file will be created in `src/backend`).

Then, simply `make -j2 dev`, which will start the backend and frontend, and restart either if their source files are changed.

## Deployments

Since this application does not provide any authentication, a proxy should provide access control (e.g. using client certificate or basic authentication).

You can build the backend with `make backend` (`backend` binary in `src/backend`) and the frontend with `make frontend` (`build` folder in `src/frontend`).

I put the backend behind an nginx reverse proxy for the `/api` path, and serve the React frontend under `/app` (statically with nginx) on the same domain. You'll want both to be on the same host so that you don't have trouble with CORS.

Relevant excerpt from my `nginx` config:

```
# App icons and a redirect to the app (helps iOS get the right icon)
root /usr/local/www/lychnos-splash;

# React build served statically
location /app/ {
		alias /usr/local/www/lychnos/src/frontend/build/;
		index  index.html index.htm;
		try_files $uri /app/index.html; # allows React paths to work without backing files
}

# Go service, listening on port 8080
location /api/ {
		# Important: don't keep a trailing path slash on proxy path (prevent /api/api)
		proxy_pass http://localhost:8080;
}
```