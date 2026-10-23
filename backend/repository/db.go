package repository

import (
	"fmt"
	"os"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

func DBConnect() *xorm.Engine {
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	dbEngine, err := xorm.NewEngine("postgres", dbURL)

	if err != nil {
		panic("Database engine creation error")
	}

	if err := dbEngine.Ping(); err != nil {
		panic("Database ping error")
	}

	return dbEngine
}

var DB = DBConnect()
