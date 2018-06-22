package main

import (
	"net/http"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gambarini/cabapi/api/internal/svc"
	"github.com/gambarini/cabapi/api/internal/data"
)



func main() {

	db, err := data.NewDb()

	if err != nil {
		log.Fatalf("Failed to connect db, %s", err)
	}

	defer db.Close()

	r := mux.NewRouter()

	svc.HandleTrips(r, db)

	http.Handle("/", r)

	log.Println("Listening on port 8000. Ctrl+C to stop")
	err = http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatalf("Failed to start server, %s", err)
	}

}





