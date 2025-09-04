package repository

import (
	"errors"
	"fmt"
	"github.com/Asylann/OrderService/server/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

var DB *sqlx.DB

func InitDBConn(cfg config.Config) error {
	var err error
	DB, err = sqlx.Open("postgres", cfg.DataBaseSource)
	if err != nil {
		return errors.New(fmt.Sprintf("Error during Open DB: %s", err))
	}
	if err = DB.Ping(); err != nil {
		return errors.New(fmt.Sprintf("Error during Open DB: %s", err))
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(4 * time.Minute)

	fmt.Println("Postgres DB is connected!!")
	return nil
}
