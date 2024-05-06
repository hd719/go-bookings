// This file contains how we are connecting our app to a specific DB

package dbrepo

import (
	"database/sql"

	"github.com/hd719/go-bookings/internal/config"
	"github.com/hd719/go-bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// type mongoDBRepo struct {
// 	App *config.AppConfig
// 	DB  *nosql.DB
// }

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

// func NewMongoRepo(conn *nosql.DB, a *config.AppConfig) repository.DatabaseRepo {
// 	return &postgresDBRepo{
// 		App: a,
// 		DB:  conn,
// 	}
// }
