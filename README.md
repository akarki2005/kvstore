# kvstore

A lightweight distributed KV store built in Go, with TCP networking and concurrent client handling.

## Quick Start

### Start the Server

First, `cd` into the project directory.

Then, run the server with `$ go run main.go`.

### Start a Client

Open up a new terminal, then run `$ telnet 127.0.0.1 8080`.

### Commands

```
GET [key] # Get the value corresponding to a key

SET [key] [value] # Set a key to a corresponding value

DELETE [key] # Remove a key
```