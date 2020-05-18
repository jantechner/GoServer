package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func New(mux *mux.Router, addr string) *http.Server {
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv
}