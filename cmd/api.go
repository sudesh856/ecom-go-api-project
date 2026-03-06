package main

import (
	
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"
	"log"

	

)


func (app *application) mount() http.Handler {

	 r := chi.NewRouter()

  // A good base middleware stack
  r.Use(middleware.RequestID)  //rate-limiting
  r.Use(middleware.RealIP)  //rate-limiting && analytics and tracing
  r.Use(middleware.Logger)  
  r.Use(middleware.Recoverer)  //recover from crashes

  // Set a timeout value on the request context (ctx), that will signal
  // through ctx.Done() that the request has timed out and further
  // processing should be stopped.
  r.Use(middleware.Timeout(60 * time.Second))

  r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Everything normal."))
  })
	// http.ListenAndServe(":3333", r)
	return r
}

func (app *application) run(h http.Handler) error{

	srv := &http.Server{
		Addr : ":" + app.config.addr,
		Handler: h,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 30,
		IdleTimeout: time.Minute,
	}

	log.Println("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()


}


type application struct {
	config config
}



type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}