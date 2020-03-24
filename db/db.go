package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (dbs *sql.DB) {
	dbs, err := sql.Open("mysql", "root@tcp(127.0.0.1:3307)/deployments")
	if err != nil {
		panic(err.Error())
	}

	return dbs
}

func Exec(db *sql.DB, query string) []Services {
	var service Services
	services_response := []Services{}

	selDB, err2 := db.Query(query)

	if err2 != nil {
		panic(err2.Error())
	}

	for selDB.Next() {
		_ = selDB.Scan(&service.SERVICE_TYPE_ID, &service.SERVICE_NAME)
		services_response = append(services_response, service)
	}

	return services_response
}
