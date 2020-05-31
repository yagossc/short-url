# Solution Description

The solution to the challenge can be split in three subsections:
- Shortened url generation **strategy**
- Code structure/divisions
- Database structure

## URL Generation Strategy
- Generate a random int64 number
- Do a base convertion (decimal to 62) mapping to chars in a - z, A -
  Z, 0 - 9
- Before inserting a new one, check for collisions and keep trying
- Insert the new URL friendly base converted string

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
- A simple table called url_map containing columns: url_id, url_short
  and url_long
- Indexing of the column url_short to simplify queries and ensure
  performance
- No nullable fields required
