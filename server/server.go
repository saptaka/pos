package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fasthttp/router"
	"github.com/saptaka/pos/api"
	"github.com/saptaka/pos/config"
	"github.com/saptaka/pos/repository"
	"github.com/valyala/fasthttp"
)

type ApiServer interface {
	Listen(port int)
}

type server struct {
	mux *router.Router
}

func (s *server) Listen(port int) {

	// loggingRouter := handlers.LoggingHandler(os.Stdout, s.mux)
	log.Printf("Starting POS server on PORT : %d", port)
	srv := &fasthttp.Server{
		ReadTimeout:  180 * time.Second,
		WriteTimeout: 180 * time.Second,
		Handler:      s.mux.Handler,
	}

	err := srv.ListenAndServe(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewServer(cfg *config.Config, repo repository.Repo) ApiServer {

	muxRouter := router.New()
	apiHandler := api.NewAPI(context.Background(), muxRouter, repo)
	apiHandler.Route()
	return &server{muxRouter}
}
