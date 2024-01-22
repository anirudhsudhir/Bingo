# Bingo

A WIP pastebin written in Go.

Built while following the "Let's Go" book.

## Usage

Clone the repository and build the application

```bash
cd bingo
go build ./cmd/web
```

Run the pastebin and open localhost:4000 in the browser

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
