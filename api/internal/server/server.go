package server

import (
	"net/http"
	"github.com/gambarini/cabapi/api/internal/data"
	"github.com/gorilla/mux"
)

type (
	Server struct {
		*http.Server
		Db *data.Db
	}
)

func NewServer(addr string, db *data.Db) *Server {

	server := &Server{
		Server: &http.Server{
			Addr: addr,
		},
		Db: db,
	}

	r := mux.NewRouter()

	http.Handle("/", r)

	HandleTrips(r, server)

	return server
}
