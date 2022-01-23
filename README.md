# Budget annually, report monthly

Use [Firefly III](https://github.com/firefly-iii/firefly-iii) to store transactions, but implement my own budgeting system. Goal: have a single app for recording and reporting, instead of using Firefly and an Excel spreadsheet.

## Screenshots

![Category summary](docs/category-summary.png) ![Transactions list](docs/transactions-list.png) ![New transaction](docs/new-transaction.png)

## Installation

You'll need to create a database and user for `lychnos`. I'm using MySQL:

```mysql
CREATE USER 'lychnos_user'@'localhost' IDENTIFIED BY 'password';
CREATE DATABASE lychnos_db;
GRANT ALL PRIVILEGES ON lychnos_db.* TO 'lychnos_user'@'localhost';
FLUSH PRIVILEGES;
```

Then, copy `.env.sample` in `src/backend` to `.env` and update the database connection string. Be sure to load the environment variables from `.env` before starting `lychnos`. The easiest way to start the backend is to `go run .` in `src/backend`. Since the application does not provide authentication, a reverse proxy should provide access control.

## Deployment

I put the backend behind an nginx reverse proxy for the `/api` path, and serve the React frontend under `/app` (statically with nginx) on the same domain.

Relevant excerpt from my nginx config:

```
location /app/ {
	# React build served statically
	proxy_pass http://192.168.72.221:80/;
}
location /api/ {
	# Go service, listening on port 8080
	# Important: don't keep a trailing path slash on proxy path (prevent /api/api)
	proxy_pass http://192.168.72.221:8080;
}
```
