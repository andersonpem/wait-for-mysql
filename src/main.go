package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := flag.String("user", "", "Database user")
	dbPassword := flag.String("password", "", "Database password")
	dbHost := flag.String("host", "", "Database host")
	dbPort := flag.String("port", "", "Database port")
	dbName := flag.String("name", "", "Database name")

	// Parse the flags
	flag.Parse()

	// Validate required flags
	if *dbUser == "" || *dbHost == "" || *dbPort == "" {
		log.Fatal("Error: Required arguments -user, -host, or -port not found")
	}
	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, *dbName)

	for {
		db, err := sql.Open("mysql", dsn)

		if err != nil {
			log.Printf("(Open) Waiting for MySQL: %s", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			if strings.Contains(err.Error(), "Access denied for user") {
				// If the server responds with 'Access denied', print a message and exit with code 0
				fmt.Println("(Ping) Server is up! But access was denied. \nMake sure to use proper credentials when actually connecting for real. Exiting!")
				db.Close()
				break
			} else {
				// If an error is encountered while connecting, wait for 1 second and try again
				log.Printf("(Ping) Waiting for MySQL: %s", err.Error())
				time.Sleep(1 * time.Second)
			}
		} else {
			// If the database responds, print a success message and exit with code 0
			fmt.Println("MySQL is up")
			db.Close()
			break
		}
		db.Close() // close the database connection to avoid connection leaks
	}
}
