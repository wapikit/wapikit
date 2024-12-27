package database

import (
	"database/sql"

	_ "ariga.io/atlas-go-sdk/atlasexec"
	_ "ariga.io/atlas-provider-gorm/gormschema"
	_ "github.com/go-jet/jet/v2"
	_ "github.com/lib/pq"
)

var DatabaseConnection *sql.DB

func GetDbInstance(dbUrl string) *sql.DB {
	if DatabaseConnection != nil {
		return DatabaseConnection
	}
	// ! TODO: use env variables here
	dsn := dbUrl
	var err error
	DatabaseConnection, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return DatabaseConnection
}
