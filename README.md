# Solution Description

The comprehensive solution's explanation will be split in three sections:
- Shortened url generation **strategy**
- Code structure/divisions
- Database structure
- In depth implementation key points

## URL Generation Strategy
### Takeaways
- Generate a random int64 number
- Do a base conversion (decimal to 62) mapping to chars in [a-z A-Z 0-9]
- Before inserting a new one, check for collisions and keep trying
- Insert the new URL friendly base converted string

### Reasoning
As a first instinct when reading the word "unique", one might
immediately go for "hashes". So, are hashes a good choice?

No, they're not a good fit for this specific scenario, because:
- Length : Most normal hashing algorithms produce long strings, which
goes against the point of a URL shortener.
- Non-Printable Characters: There might be the need to do some parsing
  of the final hash results to generate a valid URL.
- Collisions: One might say, but the final solutions actually takes
  this into account.

Given the nature of an URL, we're looking at using alpha numeric values
and even though hyphen(-) and underscore(_) are allowed in a URL,
they'll be avoided in order to not build "bad" looking URLs like
`http://domain.co/c-____`.

What now? Taking into consideration that [a-z A-Z 0-9] is a 62 value
range, we'll simplify the problem to a base conversion.

What about the non-sequencial constraint and predictability of a valid
URL? A random generated value will be passed to the base conversion
and avoid these problems.

How about uniqueness? Collisions will be checked during insertion.

Performance? To ensure performance, measures were taken in the data
storage level, since the retrieval time was the focus. For more info
see [this section](#database-structure).

## Code structure
### Summary
- package api:            the server/routes actions and routines
- package store(storage): the database interactions coming from the api
- package query:          helper package for building and running the queries
- package app:            the application data structures
- package shortener:      where the shortened path of the URL comes from

### Tree
```bash
.
├── README.md
├── api/                       # the server/routes actions and routines
│   ├── routes.go              # this file can be read as a list of routes
│   ├── server.go              # server actions i.e.: start listening, add a route
│   └── url.go                 # the request handlers for each url related endpoint
├── app/                       # the application data structures
│   └── url.go                 # the url related data definitions
├── build/                     # files related to deployment, docker, etc.
│   ├── cloudflare.compose
│   ├── initialization/
│   │   └── url.sql
│   └── postgresql.dockerfile
├── config.go                 # applications's config data definitions and a convenient load function
├── go.mod
├── go.sum
├── main.go
├── query/                    # an SQL builder and executor helper package
│   ├── builder.go
│   ├── executor.go
│   └── nocopy.go
├── shortener/                # where the shortened path of the URL comes from
│   └── converter.go          # exports the GetShortURL function
└── store/                    # the database interactions coming from the api
    └── url.go                # the url related database interactions

7 directories, 18 files
```

## Database structure
### Takeaways
- Simple model/table called url_map containing columns: url_id, url_short
  and url_long
- Indexing of the column url_short to simplify queries and ensure
  performance
- No nullable fields required

### Reasoning
Since the importance of performance was put on the retrieval side of
the operations: the redirect of the short URL. This means the lookup
is the critical piece. Therefore, the indexing of a good and fast DBM
like PostgreSQL should do the trick, keeping in mind this is a short
time simple solution.

### Implemetation details
```bash
# /short-url/build/dep/postgresql/init/url.sql

CREATE TABLE url_map (
  url_id SERIAL NOT NULL,
  url_short text NOT NULL,
  url_long text NOT NULL,

  CONSTRAINT pk_url PRIMARY KEY (url_id),
  CONSTRAINT uq_url_short UNIQUE (url_short)
);

# Relevant line
CREATE INDEX ix_url_short ON url_map (url_short);
```

# How to build

This project was built using go v1.14 (and go modules), if you have that installed, a
simple `go build` at project's root should do the trick.

> Note: Single binary, therefore, no Makefile needed.

# Running

After building the binary, simply execute it: `./short-url`

# References

- [Project layout](https://github.com/golang-standards/project-layout)
- [Code organization](https://blog.golang.org/organizing-go-code)
- [Naming conventions](https://blog.golang.org/package-names)
- [Golangci-lint](https://github.com/golangci/golangci-lint)
- [sqlx](https://github.com/jmoiron/sqlx)
- [echo](https://github.com/labstack/echo)
