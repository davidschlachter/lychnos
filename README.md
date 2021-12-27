# Budget annually, report monthly

Use [Firefly III](https://github.com/firefly-iii/firefly-iii) to store transactions, but implement my own budgeting system. Goal: have a single app for recording and reporting, instead of using Firefly and an Excel spreadsheet.

## Installation

You'll need to create a database and user for `lychnos`. I'm using MySQL:

```mysql
CREATE USER 'lychnos_user'@'localhost' IDENTIFIED BY 'password';
CREATE DATABASE lychnos_db;
GRANT ALL PRIVILEGES ON lychnos_db.* TO 'lychnos_user'@'localhost';
FLUSH PRIVILEGES;
```

Then, copy `.env.sample` in `src/backend` to `.env` and update the database connection string. Be sure to load the environment variables from `.env` before starting `lychnos`. The easiest way to start the backend is to `go run .` in `src/backend`.