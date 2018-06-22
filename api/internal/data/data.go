package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/jmoiron/sqlx"
	"github.com/gambarini/cabapi/api/internal/cache"
	"github.com/gambarini/cabapi/api/internal/model"
)

const (
	queryByMedallion = "SELECT medallion, DAY(pickup_datetime) day, MONTH(pickup_datetime) month, YEAR(pickup_datetime) year, count(*) trips " +
		"FROM cab_trip_data " +
		"GROUP BY medallion, DAY(pickup_datetime), MONTH(pickup_datetime), YEAR(pickup_datetime) " +
		"HAVING medallion IN (?) AND " +
		"day = ? AND " +
		"month = ? AND " +
		"year = ?"
	queryAll = "SELECT medallion, DAY(pickup_datetime) day, MONTH(pickup_datetime) month, YEAR(pickup_datetime) year, count(*) trips " +
		"FROM cab_trip_data " +
		"GROUP BY medallion, DAY(pickup_datetime), MONTH(pickup_datetime), YEAR(pickup_datetime) " +
		"HAVING day = ? AND " +
		"month = ? AND " +
		"year = ?"
)

type (
	Db struct {
		db       *sql.DB
		cache    *cache.Cache
	}
)

func NewDb() (data *Db, err error) {

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/cabapi")
	if err != nil {
		return data, fmt.Errorf("failed to connect to mysql, %s", err)
	}

	c, err := cache.NewCache()
	if err != nil {
		return data, fmt.Errorf("Failed to get cache connection, %s", err)
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

func (db *Db) FindTrips(date time.Time, medallions []string, useCache bool) (trips []model.Trips, err error) {

	trips = make([]model.Trips, 0)

	var rows *sql.Rows

	if len(medallions) == 0 {
		log.Printf("Query for date %s", date)

		rows, err = db.db.Query(queryAll, date.Day(), date.Month(), date.Year())

		if err != nil {
			return trips, err
		}

	} else {

		if useCache {
			trips, medallions = db.fetchCache(date, medallions)
		}

		if len(medallions) > 0 {

			log.Printf("Query for medallion %s, date %s", medallions, date)

			q, args, err := sqlx.In(queryByMedallion, medallions, date.Day(), int(date.Month()), date.Year())

			log.Printf("Query %s, %v", q, args)

			if err != nil {
				return trips, err
			}

			rows, err = db.db.Query(q, args...)

			if err != nil {
				return trips, err
			}
		}

	}

	if rows == nil {
		return trips, nil
	}

	tripsIn := make(chan model.Trips, 5)
	go db.updateCache(tripsIn)

	for rows.Next() {
		var medl string
		var d, m, y, count int

		err := rows.Scan(&medl, &d, &m, &y, &count)

		if err != nil {
			return trips, err
		}

		trip := model.Trips{
			Medallion: medl,
			Total:     count,
			Date:      fmt.Sprintf("%d-%d-%d", y, m, d),
		}

		trips = append(trips, trip)

		tripsIn <- trip
	}

	close(tripsIn)

	return trips, nil

}

func (db *Db) fetchCache(date time.Time, medallions []string) (trips []model.Trips, toDb []string) {

	for _, medl := range medallions {

		trip, err := db.cache.Get(medl, date)

		if err != nil {
			log.Printf("Failed to fetch cache, %s", err)
			toDb = append(toDb, medl)
			continue
		}

		log.Printf("From cache, %v", trip)
		trips = append(trips, trip)

	}

	return trips, toDb
}

func (db *Db) updateCache(tripsIn chan model.Trips) {

	for trip := range tripsIn {
		err := db.cache.Set(trip)

		if err != nil {
			log.Printf("Failed to update cache, %s", err)
			continue
		}

	}
}
