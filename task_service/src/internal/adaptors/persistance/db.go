//database configuration

package persistance

import (
	"database/sql"
	"task_service/src/internal/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB //pointer to database
}

func NewDatabase() (*Database, error) {
	config, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME, config.DB_SSLMODE)

	//printing database connection string
	fmt.Println("Database URL:", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db: db}, nil
} //It is giving a new object of Database struct

func (d *Database) Close() {
	d.db.Close()
}
func (d *Database) GetDB() *sql.DB {
	return d.db
} //It is giving the field of that Database Object

//! when NewDatabase() is returning the database, then what is the use of getDB() method?
