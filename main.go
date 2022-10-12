package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" //Use of blank identifier _ when importing the driver
)

const (
	username = "root"
	password = "Password"
	hostname = "127.0.0.1:3306"
	dbname   = "golang_mysql"
)

func dsn(dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbname)
}

func main() {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelfunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)

	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	no, err := res.RowsAffected()

	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}

	log.Printf("rows affected %d\n", no)

	db.Close()

	db, err = sql.Open("mysql", dsn(dbname))

	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}

	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	err = db.PingContext(ctx)

	if err != nil {
		log.Printf("Error %s pingingdb", err)
		return
	}

	log.Printf("Connected to DB %s successfully /n", dbname)

}
