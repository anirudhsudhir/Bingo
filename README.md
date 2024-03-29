# Bingo

A server-side rendered pastebin written in Go.

It persists data to a MySQL database and features middleware, logging, dependency injection and support for sessions.

## Usage

1. Clone the repository and build the application

```bash
cd bingo
go build ./cmd/web
```

2. Create and set up a new MySQL database using the contents of 'init_db.sql'

3. Add the MySQL database user and password to a new secrets.env file

```text
Sample secrets.env

DBuser = "user"
DBpass = "pass"
```

4. Run the pastebin and open localhost:4000 in the browser

```bash
./web
```

## Command-line flags

Note: All flags can be viewed by invoking Go's built-in help flag

```bash
./web -help
```

1. addr : Sets the HTTP network address (defaults to localhost:4000)

```bash
./web -addr=":5000"
```

2. dsn: Sets the data source name of the MySQL database (defaults to the user and password set in secrets.env)

```bash
./web -dsn="[user]:[password]@/bingo?parseTime=true"
```

## Note

A session key will be generated and stored in secrets.env when the application is run for the first time. This will be used to sign cookies created by the pastebin.
To create a new session key, delete the `SessionKey = "randomsessionkey"` entry in secrets.env and clear the "session" cookie stored in the browser.

This project was built while following the "Let's Go" book.
