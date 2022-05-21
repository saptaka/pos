package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/saptaka/pos/api"
	"github.com/saptaka/pos/config"
	"github.com/saptaka/pos/repository"
)

type ApiServer interface {
	Listen(port int)
}

type server struct {
	mux *mux.Router
}

func (s *server) Listen(port int) {

	loggingRouter := handlers.LoggingHandler(os.Stdout, s.mux)
	log.Printf("Starting POS server on PORT : %d", port)
	srv := &http.Server{
		ReadTimeout:       180 * time.Second,
		WriteTimeout:      180 * time.Second,
		ReadHeaderTimeout: 180 * time.Second,
		Handler:           loggingRouter,
		Addr:              fmt.Sprintf(":%d", port),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewServer(cfg *config.Config, repo repository.Repo) ApiServer {

	muxRouter := mux.NewRouter()
	apiHandler := api.NewAPI(context.Background(), muxRouter, repo)
	apiHandler.Route()
	return &server{muxRouter}
}
