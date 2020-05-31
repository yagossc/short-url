package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	Port     uint64
	DBDriver string
	DBURL    string
}

func loadConfig() config {
	// must be separate calls, to avoid the 'early exit' when one of the files
	// does not exist
	_ = godotenv.Load(".env.local")
	_ = godotenv.Load(".env")

	var cfg config

	if s, ok := os.LookupEnv("TAPI_PORT"); ok {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			fmt.Printf("%v\n", err)
			v = 8080
		}

		cfg.Port = v
	}

	if s, ok := os.LookupEnv("TAPI_DB_DRIVER"); ok {
		if s == "" {
			s = "pgx"
		}

		cfg.DBDriver = s
	}

	if s, ok := os.LookupEnv("TAPI_DB_URL"); ok {
		cfg.DBURL = s
	}

	return cfg

}
