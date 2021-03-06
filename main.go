package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yagossc/short-url/api"
	"github.com/yagossc/short-url/query"
)

func main() {

	// load configuration
	cfg := loadConfig()

	// database connection
	db, err := openDBConnection(cfg)
	if err != nil {
		os.Exit(1)
	}
	defer db.Close()

	// 'generic' db layer (prototype version)
	executor := query.NewExecutor(db)

	// server configuration
	e := echo.New()
	e.Use(middleware.Recover())

	// Create server
	s := api.NewServer(executor, e, cfg.BaseURL+":"+strconv.FormatUint(cfg.Port, 10))

	// API routes
	s.Routes()

	log.Fatal(s.Start(":" + strconv.FormatUint(cfg.Port, 10)))
}

func openDBConnection(cfg config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < 5; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}

		fmt.Printf("Connecting to database (tries=%d)... ", i+1)
		db, err = sqlx.Open(cfg.DBDriver, cfg.DBURL)
		if err != nil {
			fmt.Printf("ERROR!\n%v\n\n", err)
			continue
		}

		err = db.Ping()
		if err != nil {
			fmt.Printf("ERROR!\ndatabase error: %v\n", err)
		} else {
			break
		}
	}

	return db, err
}
