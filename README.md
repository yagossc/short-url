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
Analyzing the problems's text brings up the following questions, that,
when answered, show path to the solution's motivations.

___
**Question: Are hashes a good choice here?**\
**Answer:** No, they're not a good fit for this specific scenario for
the followig reasons:
- **Length:** Common hashing algorithms produce long strings, which
goes against the point of a URL shortener.
- **Non-Printable Characters:** Final hash results might not be a valid URL.
- **Collisions:** One might say, but the final solution actually dabbles
  in this.

Given the nature of a URL, we're looking at using alpha numeric values
and even though hyphen(-) and underscore(_) are allowed in a URL,
they'll be avoided in order to not build "bad" looking URLs like
`http://domain.co/c-____`.

___
**Question: What will be the short URL identifiers and how to get them?**\
**Answer:** Taking into consideration what was said above about URLs
and that [a-z A-Z 0-9] is a 62 value range, we'll simplify the problem
to a base conversion.

___
**Question: Can this base represent enough values?**\
**Answer:** Some quick math.

With 62 chars and a unique string of, let's say, 7 characters long we
can represent:
- **62⁷ = 3,521,614,606,208 URLs.**

That's ~3.5 trillion URLs. What about 8 characters long?
- **62⁸ = 218,340,105,584,896 URLs.**

That's ~218 trillion URLs. How about 10 characters long?
- **62¹⁰ = 839,299,365,868,340,224 URLs.**

Well, that's a really large number.

___
**Question: What about the non-sequential constraint and
predictability of a valid URL?**\
**Answer:** A random generated value will be passed to the base conversion function
and avoid these problems.

___
**Question: How about uniqueness?**\
**Answer:** Collisions will be checked during insertion.

___
**Question: Performance?**\
**Answer:** To ensure performance, measures were taken in the data
storage level. For more info see [this section](#database-structure).

## Code structure
### Summary
Here's a quick summary of the code's relevant packages and their responsibilities:
- **package api:**            the server/routes actions and routines.
- **package app:**            the application data structures.
- **package history:**        utility service for checking dates.
- **package query:**          helper service for building and running SQL queries.
- **package shortener:**      service for generating the shortened URL identifier.
- **package store**(storage): service handler for database interactions coming from the API.

### Tree -F
A visual representation of the code's file structure:
```bash
.
├── README.md
├── api/                                  # the server/routes actions and routines
│   ├── aux_test.go                       # testing auxiliary routines
│   ├── reqHistory.go                     # the handlers for each history related endpoint
│   ├── reqHistory_test.go
│   ├── routes.go                         # this file can be read as a list of routes
│   ├── server.go                         # server actions i.e.: start listening, add a route
│   ├── url.go                            # the handlers for each url related endpoint
│   └── url_test.go
├── app/                                  # the application data structures
│   ├── reqHistory.go                     # the history related data definitions
│   └── url.go                            # the url related data definitions
├── build/                                # files related to deployment, docker, etc.
│   ├── cloudflare.compose
│   └── dep/
│       └── postgresql/
│           ├── init/
│           │   ├── 000001-url.sql
│           │   └── 000002-requests.sql
│           ├── postgresql.conf
│           └── postgresql.dockerfile
├── config.go                             # applications's config data definitions
├── go.mod
├── go.sum
├── history/
│   └── history.go                        # utility service for checking dates
├── internal/
│   └── dbtest/                           # useful db handler for testing
│       └── dbtest.go
├── main.go
├── query/                                # an SQL builder and executor helper service
│   ├── builder.go
│   ├── executor.go
│   └── nocopy.go
├── shortener/                            # service for generating the shortened URL identifier.
│   └── converter.go
└── store/                                # the database interactions coming from the api
    ├── reqHistory.go                     # the history related database interactions
    └── url.go                            # the url related database interactions

```

## Database structure

### Summary
- Table `url_map` contains columns `url_short` and `url_long`.
- Indexing of the column `url_short` avoids double way base
  conversions and tries to ensure performance.
- The requests history, i.e., table `req_history`, holds a
  timestamp(`req_time`) and a foreign key (`url_short`) from table `url_map`.
- Again indexing of the column `url_short` to improve/ensure performance.

### Reasoning
The importance of performance was put on the retrieval side of
operations: the redirect behind the short URL. This means lookup is
the critical piece.
Therefore, the DBM's (PostgreSQL) indexing should do the trick and,
although some other tweaks like the use of caches or some sort of
rotation between RAM and Disk usage could be beneficial, this will be
the only adopted strategy, for simplicity's sake.

### Data definition details
#### URL Mapping
```sql
# /short-url/build/dep/postgresql/init/000001-url.sql

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

#### Requests History
```sql
# /short-url/build/dep/postgresql/init/000002-requests.sql

CREATE TABLE req_history (
  req_id SERIAL NOT NULL,
  url_short TEXT NOT NULL,
  req_time BIGINT NOT NULL,

  CONSTRAINT pk_req PRIMARY KEY (req_id),
  CONSTRAINT fk_url FOREIGN KEY (url_short) REFERENCES url_map (url_short)
);

CREATE INDEX ix_req_url ON req_history (url_short)
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

### Request ocurrences counter
```golang
// GetEntriesInInvertval counts the ocurrences of a req given a time
// interval in hours. For example, to get the entries for the past day:
// ocurrences := GetEntriesInInvertval(entries, 24)
// To Get the entries for the past week:
// ocurrences := GetEntriesInInvertval(entries, 24*7)
func GetEntriesInInvertval(entries []app.ReqHistory, interval int) int {
    ocurrences := 0
    currTime := time.Now()

    for _, entry := range entries {
        t := time.Unix(entry.ReqTime, 0)
        diff := int64(currTime.Sub(t).Hours()) / int64(interval)

        if diff < 1 {
            ocurrences++
        }
    }

    return ocurrences
}
```

# Available Routes
### Static routes
**Shortener:**
<pre>
<b>METHOD:</b> POST
<b>PATH:</b> "/"
<b>BODY:</b> JSON {"url":"http://some.url"}
<b>STATUS:</b> 201
<b>RESPONSE:</b> JSON {"url":"http://shortened.url/id"}
</pre>

**History:**
<pre>
<b>METHOD:</b> GET
<b>PATH:</b> "/history"
<b>BODY:</b> JSON {"url":"http://some.url"}
<b>STATUS:</b> 200
<b>RESPONSE:</b> JSON {"count": number}
</pre>

<pre>
<b>METHOD:</b> GET
<b>PATH:</b> "/history/week"
<b>BODY:</b> JSON {"url":"http://some.url"}
<b>STATUS:</b> 200
<b>RESPONSE:</b> JSON {"count": number}
</pre>

<pre>
<b>METHOD:</b> GET
<b>PATH:</b> "/history/day"
<b>BODY:</b> JSON {"url":"http://some.url"}
<b>STATUS:</b> 200
<b>RESPONSE:</b> JSON {"count": number}
</pre>

### Dynamic routes
Each time a short URL is generated, a new route is added to the API in
the form of `APIBase`+`/ShortURLID`:
<pre>
<b>METHOD:</b> GET
<b>PATH:</b> "/ShortURLID"
<b>STATUS:</b> 301
<b>RESPONSE:</b> Redirect
</pre>

# How to test
## Go test tool
There are some tests for the API endpoints that can be run with one of
the commands bellow:
```
go test ./...

# for verbosity
go test -v ./...

# for coverage
go test -cover ./...

# for detailed coverage
go test ./... --coverprofile=coverage.out && \
go tool cover -html=coverage.out
```

## Testing with cURL
This assumes the application's running on `localhost:8080`.
### Static routes
**Shortener:**
```bash
curl -X POST -i -H "Content-Type: application/json" \
    -H "Accept: application/json"                   \
    http://localhost:8080/                          \
    -d '{"URL":"http://www.google.com"}'
```

**History:**\
- **Full**
```bash
curl -X GET -i -H "Content-Type: application/json" \
    -H "Accept: application/json"                  \
    http://localhost:8080/history                  \
    -d '{"URL":"http://localhost/ShortURLIdentifier"}'
```
- **Week**
```bash
curl -X -i GET -H "Content-Type: application/json" \
p    -H "Accept: application/json"                 \
    http://localhost:8080/history/week             \
    -d '{"URL":"http://localhost/ShortURLIdentifier"}'
```
- **Day**
```bash
curl -X GET -i -H "Content-Type: application/json" \
    -H "Accept: application/json"                  \
    http://localhost:8080/history/day              \
    -d '{"URL":"http://localhost/ShortURLIdentifier"}'
```

### Dynamic routes
The new shortened URL is the response from the `Shortener` endpoint.
```bash
curl -X GET -i -L http://localhost:8080/ShortURLIdentifier
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
