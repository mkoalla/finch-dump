package db

import (
	"context"
	"finchd/lib/api"
	"iter"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("couldn't connect to " + path)
	}

	return db, db.AutoMigrate(&api.JobEntry{})
}

func WriteData(output string, jobs iter.Seq2[*api.JobEntry, error], count *int, ctx *context.Context) (err error) {
	return
}
