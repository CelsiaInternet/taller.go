package module

import (
	"fmt"
	"net/http"

	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/envar"
	"github.com/celsiainternet/elvis/middleware"
	"github.com/celsiainternet/elvis/response"
	"github.com/celsiainternet/elvis/strs"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	v1 "github.com/celsiainternet/taller/internal/service/taller/v1"
)

type Server struct {
	http *http.Server
}

func New() (*Server, error) {
	server := Server{}

	port := envar.EnvarInt(3300, "PORT")
	if port == 0 {
		return nil, fmt.Errorf("variable PORT es requerida")
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	latest := v1.New()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.HTTPError(w, r, http.StatusNotFound, "404 Not Found")
	})

	r.Mount("/", latest)
	r.Mount("/v1", latest)

	handler := cors.AllowAll().Handler(r)
	addr := strs.Format(":%d", port)
	serv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	server.http = serv

	return &server, nil
}

func (serv *Server) Close() {
	v1.Close()

	console.LogK("Http", "Shutting down server...")
}

func (serv *Server) StartHttpServer() {
	if serv.http == nil {
		return
	}

	svr := serv.http
	console.LogKF("Http", "Running on http://localhost%s", svr.Addr)
	console.Fatal(serv.http.ListenAndServe())
}

func (serv *Server) Start() {
	go serv.StartHttpServer()

	v1.Banner()

	<-make(chan struct{})
}
