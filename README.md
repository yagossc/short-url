# Solution Description

The comprehensive solution's explanation is splitted in four sections:
- Shortened URL generation [strategy](#url-generation-strategy).
- [Code structure](#code-structure).
- [Database structure](#database-structure).
- In depth implementation [key points](#in-depth-implementation).

## URL Generation Strategy

### Summary
- Generate a random int64 number.
- Do a base conversion (decimal to base62) mapping to chars in [a-z A-Z 0-9].
- Before insertion, loop checking for collisions.
- Insert the new URL friendly base converted string.

### Reasoning
As a first instinct when reading the problems's text, one might
immediately go for "hashes". So, are hashes a good choice?

**No**, they're not a good fit for this specific scenario, because:
- **Length:** Common hashing algorithms produce long strings, which
goes against the point of a URL shortener.
- **Non-Printable Characters:** Final hash results might not be a valid URL.
- **Collisions:** One might say, but the final solution actually dabbles
  in this.

Given the nature of a URL, we're looking at using alpha numeric values
and even though hyphen(-) and underscore(_) are allowed in a URL,
they'll be avoided in order to not build "bad" looking URLs like
`http://domain.co/c-____`.

**What now?** Taking into consideration that [a-z A-Z 0-9] is a 62 value
range, we'll simplify the problem to a base conversion.

**Is this range enough?** Some quick math:

With 62 chars and a unique string of, let's say, 7 characters long we
can represent:
- **62⁷ = 3,521,614,606,208 URLs.**

That's ~3.5 trillion URLs. What about 8 characters long?
- **62⁸ = 218,340,105,584,896 URLs.**

That's ~218 trillion URLs. How about 10 characters long?
- **62¹⁰ = 839,299,365,868,340,224 URLs.**

Well, that's a really large number.


**What about the non-sequential constraint and predictability of a valid
URL?** A random generated value will be passed to the base conversion function
and avoid these problems.

**How about uniqueness?** Collisions will be checked during insertion.

**Performance?** To ensure performance, measures were taken in the data
storage level. For more info see [this section](#database-structure).

## Code structure
### Summary
Here's a quick summary of the code's packages and their responsibilities:
- **package api:**            the server/routes actions and routines.
- **package app:**            the application data structures.
- **package query:**          helper package for building and running the queries.
- **package shortener:**      where the shortened path of the URL comes
  from.
- **package store**(storage): the database interactions coming from the api.

### Tree
A visual representation of the code's file structure:
```bash
.
├── README.md
├── api/                               # the server/routes actions and routines
│   ├── routes.go                      # this file can be read as a list of routes
│   ├── server.go                      # server actions i.e.: start listening, add a route
│   └── url.go                         # the request handlers for each url related endpoint
├── app/                               # the application data structures
│   └── url.go                         # the url related data definitions
├── build/                             # files related to deployment, docker, etc.
│   ├── cloudflare.compose
│   └── dep/
│       └── postgresql/
│           ├── init/
│           │   └── url.sql
│           ├── postgresql.conf
│           └── postgresql.dockerfile
├── config.go                          # applications's config data definitions
├── go.mod
├── go.sum
├── main.go
├── query/                             # an SQL builder and executor helper package
│   ├── builder.go
│   ├── executor.go
│   └── nocopy.go
├── shortener/                         # where the shortened path of the URL comes from
│   └── converter.go                   # exports the GetShortURL function
└── store/                             # the database interactions coming from the api
    └── url.go                         # the url related database interactions

```

## Database structure

### Summary
- The data definition is quite simple, containing only the Short and Long URLs.
- Simple model/table named `url_map` containing columns: `url_id`, `url_short`
  and `url_long`.
- Indexing of the column `url_short` to avoid double way base conversions and ensure
  performance.

### Reasoning
The importance of performance was put on the retrieval side of
operations: the redirect behind short URL. This means the lookup is
the critical piece.
Therefore, the DBM's (PostgreSQL) indexing should do the trick and,
although some other tweaks like the use of caches or some sort of
rotation between RAM and Disk usage could be beneficial, this will be
the only adopted strategy, for simplicity's sake.

### Data definition details
```sql
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

## In depth implementation
Shown bellow are some key points/portions of code.
### Base conversion algorithm
```golang
// FILE: short-url/shortener/converter.go

var valueMapping = []string{
    "a", "b", "c", "d", "e", "f", "g", "h", "i",
    "j", "k", "l", "m", "n", "o", "p", "q", "r",
    "s", "t", "u", "v", "w", "x", "y", "z", "A",
    "B", "C", "D", "E", "F", "G", "H", "I", "J",
    "K", "L", "M", "N", "O", "P", "Q", "R", "S",
    "T", "U", "V", "W", "X", "Y", "Z", "0", "1",
    "2", "3", "4", "5", "6", "7", "8", "9"}

// GetShortURL converts a given decimal
// number to our base62 string
func GetShortURL(num int64) string {
    var converted strings.Builder

    for num > 0 {
        converted.WriteString(valueMapping[num%62])
        num = num / 62
    }

    return converted.String()
}
```

### Server route injection
```golang
// FILE: short-url/api/server.go

// AddRoute does the dynamic route injection and is what
// gives the API the expected url shortener behavior.
func (s *Server) AddRoute(id string) {
    newRoute := "/" + id
    s.e.GET(newRoute, s.redirect)
}
```

### Collision check during insertion
```golang
// FILE: short-url/api/url.go[48-62]
    ...
    short := shortener.GetShortURL(rand.Int63())
    existent, err := store.FindURLByShort(s.db, short)
    if err != nil {
        return err
    }

    // Check for a collision and try again if needed
    for existent != nil || short == "" {
        var err error
        short = shortener.GetShortURL(rand.Int63())
        existent, err = store.FindURLByShort(s.db, short)
        if err != nil {
            return err
        }
    }

    ...
```

# How to build

This project was built using go v1.14 (and go modules), if you have that installed, a
simple `go build` at the project's root should do the trick.

> Note: Single binary, therefore, no Makefile needed.

# Running

Make sure to start the database container:
```bash
> docker-compose -f build/cloudflare.compose up -d
```

Copy the .env file to .env.local if any customization is needed:
```bash
> cp .env .env.local
```
Now edit the file based on your needs. For example:
- Change the variable `TAPI_URL` to the desired base for the shortened
URLs.
- Change the applications listening port on `TAPI_PORT`.
- Change the "host" part of the variable `TAPI_DB_URL` to the correct database
location. In case the container is in the same machine, give it the
container's name.

And, after building the binary, simply execute it: `./short-url`

# References

- [Project layout](https://github.com/golang-standards/project-layout)
- [Code organization](https://blog.golang.org/organizing-go-code)
- [Naming conventions](https://blog.golang.org/package-names)
- [Golangci-lint](https://github.com/golangci/golangci-lint)
- [sqlx](https://github.com/jmoiron/sqlx)
- [echo](https://github.com/labstack/echo)
- [PostgreSQL Indexes](https://www.postgresql.org/docs/9.1/sql-createindex.html)
