package database

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

//Init ...
func Init() {
	log.Printf("Init database")
	dbUser := "epgasiygcgfrzf"
	dbPassword := "2b25ff8a36774ced0e133a57c01e0d38b9b55f98d01b7cfbb98be3c55f5978af8"
	dbName := "fs"
	dbHost := "localhost"
	dbPort := "5432"

	port, e := strconv.Atoi(dbPort)
	if e != nil {
		log.Fatal(e)
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, port, dbUser, dbPassword, dbName)

	var err error
	if db != nil {
		db.Close()
	}

	db, err = ConnectDB(dbInfo)
	if err != nil {
		log.Fatal(err)
	}
}

//InitForTest inits connections for specific test schema. Set `search_path` only works  for single connection.
// So we need init connection  for specific schema  to prevent unexpected result. See more at https://stackoverflow.com/a/56368340
func InitForTest(testSchema string) {
	dbUser := "epgasiygcgfrzf"
	dbPassword := "b25ff8a36774ced0e133a57c01e0d38b9b55f98d01b7cfbb98be3c55f5978af8"
	dbName := "fs"
	dbHost := "localhost"
	dbPort := "5432"

	port, e := strconv.Atoi(dbPort)
	if e != nil {
		log.Fatal(e)
	}
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		dbHost, port, dbUser, dbPassword, dbName, testSchema)

	var err error

	if db != nil {
		db.Close()
	}

	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}
}

//ConnectDB ...
func ConnectDB(dbinfo string) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(3)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(500)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Hour)

	env := "prod"

	// enable logmode to see sql in debug mode
	if strings.ToLower(env) == "debug" {
		db.LogMode(true)
	}

	return db, nil
}

//GetDB ...
func GetDB() *gorm.DB {
	return db
}

//CloseDB before close application
func CloseDB() {
	db.Close()
	db = nil
}

//SetDB set shared db
func SetDB(database *gorm.DB) {
	db = database
}
