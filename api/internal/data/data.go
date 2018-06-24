package data

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/gambarini/cabapi/api/internal/cache"
	_ "github.com/go-sql-driver/mysql"
)

type (
	Db struct {
		db    *sql.DB
		cache *cache.Cache
	}
)

func NewDb() (data *Db, err error) {

	var db *sql.DB
	wait := make(chan int)

	go func() {

		defer close(wait)
		retry := 1

		for {
			db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/cabapi")

			if err == nil || retry == 5 {
				return
			}

			retry++
			time.Sleep(time.Second * 5)
		}
	}()

	<-wait

	if err != nil {
		return data, fmt.Errorf("failed to connect to mysql, %s", err)
	}

	c, err := cache.NewCache()
	if err != nil {
		return data, fmt.Errorf("failed to get cache connection, %s", err)
	}

	return &Db{
		db:    db,
		cache: c,
	}, nil
}

func (db *Db) Close() {
	db.db.Close()
	db.cache.Close()
}
