      __ _   ___    __ _  _ __  _   _
     / _` | / _ \  / _` || '__|| | | |
    | (_| || (_) || (_| || |   | |_| |
     \__, | \___/  \__,_||_|    \__, |
      __/ |                      __/ |
     |___/                      |___/ 

# goary - sample web.go application

## Requirements

- import "github.com/sschober/web.go"

## Features

Very simple twitter-like application which lets you create small
tweets via a restfull interface. Supports json as on-the-wire
representation.

## Usage

Just build and run the binary. The server listens on localhost:9999.
A list of entries is reachable at http://localhost:9999/roars (GET).
That uri also understands POST-Requests (Content-Type:
application/json has to be set). Single entries are addressed by
appending their id to the formentioned uri.

## Author

Sven Schober
