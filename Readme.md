      __ _   ___    __ _  _ __  _   _
     / _` | / _ \  / _` || '__|| | | |
    | (_| || (_) || (_| || |   | |_| |
     \__, | \___/  \__,_||_|    \__, |
      __/ |                      __/ |
     |___/                      |___/ 

# goary - sample web.go application

## Requirements

    import "github.com/sschober/web.go"

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

## Resources/API

- /roars (Collection-style)

  - GET: List of all entries.
    Default Content-Type: text/plain
    Optional Content-Type: application/json

  - POST: Create new entry.
    Content-Type: application/json

- /roars/{id}

  - GET: Returns specified entry with {id}
    Default Content-Type: text/plain
    Optional Content-Type: application/json

  - DELETE: Deletes specified entry with {id}

## Examples (Client)

List all entries as strings:

    $ curl http://localhost:9999/roars

List all entries as json:

    $ curl -H "Content-Type: application/json" \
    http://localhost:9999/roars

List a single entry as string:

    $ curl http://localhost:9999/roars/0

Add an entry:

    $ curl -i -X POST -d \
    '{"Author":"Sven","Text":"Hello!!!!","CreationDate":"Sat, 10 Jul 2010 19:46:48 CEST"}' \
    -H "Content-Type: application/json" \
    http://localhost:9999/roars

Delete an entry:

    $ curl -i -X DELETE http://localhost:9999/roars/2

## Author

Sven Schober
