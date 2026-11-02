package db

import (
	"context"
	"finchd/lib"
	"finchd/lib/api"
	"fmt"
	"iter"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DEFAULT_CHUNK_SIZE = 50

func InitDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to %s: %w", path, err)
	}

	return db, db.AutoMigrate(&api.JobEntry{})
}

func writeChunk[T any](db *gorm.DB, data []T) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(data, len(data)).Error
	})
}
func WriteData(output string, jobs iter.Seq2[*api.JobEntry, error], count *int, ctx *context.Context) (err error) {
	db, err := InitDB(output)
	if err != nil {
		return
	}
	_sql, err := db.DB()
	if err != nil {
		return
	}
	defer func() {
		if cerr := _sql.Close(); cerr != nil {
			log.Println("close db:", cerr)
		}
	}()
	counted := false
	defer func() {
		if counted {
			fmt.Println()
		}
	}()
	rows := make([]api.JobEntry, 0, DEFAULT_CHUNK_SIZE) // def chunk size i suppose
	for data, err := range jobs {
		if err != nil {
			break
		}
		if ctx != nil && (*ctx).Err() != nil {
			return (*ctx).Err()
		}

		rows = append(rows, *data)
		if len(rows) < DEFAULT_CHUNK_SIZE {
			continue
		}
		if err = writeChunk(db, rows); err != nil {
			break
		}
		*count += len(rows)
		fmt.Printf("\r\033[2K%s recv %d\t(if this is stuck more than 30s, Ctrl-C to stop early)", lib.InteractiveLoading(), *count)
		counted = true
		rows = rows[:0]
	}
	if err != nil {
		return
	}
	if len(rows) == 0 {
		return
	}
	if err = writeChunk(db, rows); err != nil {
		return
	}
	*count += len(rows)
	return
}
