// Package repotest gives tests an isolated in-memory database wired into repository.DB.
package repotest

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"github.com/randspace0/ketikode/backend/repository"
	"github.com/randspace0/ketikode/backend/repository/migrate"
)

func Open(t *testing.T) {
	t.Helper()
	// A unique file-less database per test; the pool is capped at one connection
	// because each in-memory connection would otherwise see its own empty database.
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory"), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { sqlDB.Close() })

	if err := migrate.Run(db); err != nil {
		t.Fatal(err)
	}
	repository.DB = db
}
