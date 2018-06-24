package main

import (
	"log"
	"github.com/gambarini/cabapi/api/internal/data"
	"github.com/gambarini/cabapi/api/internal/server"
)

func main() {

	db, err := data.NewDb()

	if err != nil {
		log.Fatalf("Failed to connect db, %s", err)
	}

	defer db.Close()

	srv := server.NewServer(":8000", db)

	log.Println("Listening on port 8000. Ctrl+C to stop")
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatalf("Failed to start server, %s", err)
	}

}
