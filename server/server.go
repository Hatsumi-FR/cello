package server

import (
	"cello/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
)

type Server struct {
}

func (s Server) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/ws", handlers.WsHandler)

	log.Println("server started on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
