package db

import (
	"category-crud/config"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

func Configure(config config.Template) (*sql.DB, *goqu.Database, error) {
	const dialect = "postgres"
	db, err := sql.Open(dialect, config.DB.ConnectionString)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	if config.App.Debug {
		log.Default().Println("Connected to database")
	}

	db.SetMaxOpenConns(config.DB.MaxOpenConns)
	db.SetMaxIdleConns(config.DB.MaxIdleConns)

	return db, goqu.New(dialect, db), nil

}
