package db

import (
"database/sql"
"fmt"

_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "postgres"
)

var connectionString = fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
var DBCon *sql.DB

func init(){
	dbCon, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	DBCon = dbCon
}