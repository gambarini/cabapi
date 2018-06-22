package svc

import (
	"log"
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/gambarini/cabapi/api/internal/data"
	"errors"
	"strings"
)

var (
	ErrNotPickupDate     = errors.New("no pickup date found")
	ErrInvalidPickupDate = errors.New("invalid pickup date. Expected date=YYYY-MM-DD")
	Db *data.Db
)

func HandleTrips(r *mux.Router, db *data.Db) {

	Db = db

	r.HandleFunc("/trip/{date}", getTrips)

}

func getTrips(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	strDate := vars["date"]

	date, err := resolvePickupDate(strDate)

	if err == ErrInvalidPickupDate {
		log.Println(err)
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	medallions, err := resolveMedallions(request)

	if err == ErrInvalidPickupDate {
		log.Println(err)
		writer.WriteHeader(400)
		writer.Write([]byte(err.Error()))
		return
	}

	trips, err := Db.FindTrips(date, medallions, resolveUseCache(request))

	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	payloadJSON, err := json.Marshal(trips)

	if err != nil {
		log.Println(err)
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(payloadJSON)
}

func resolvePickupDate(dateStr string) (date time.Time, err error) {

	log.Printf("Date Param %s", dateStr)

	if dateStr == "" {
		return date, ErrNotPickupDate
	}

	date, err = time.Parse("2006-01-02", dateStr)

	if err != nil {
		return date, ErrInvalidPickupDate
	}

	log.Printf("Date: %s", date)

	return date, nil
}

func resolveMedallions(request *http.Request) (medallions []string, err error) {

	medlStr := request.URL.Query().Get("medallion")

	log.Printf("Medallion Param %s", medlStr)

	if medlStr == "" {
		return medallions, nil
	}

	medallions = strings.Split(medlStr, ",")

	log.Printf("medallions: %v", medallions)

	return medallions, nil
}

func resolveUseCache(request *http.Request) (useCache bool) {

	cache := request.URL.Query().Get("cache")

	if cache != "false" {
		return true
	}

	return false
}
