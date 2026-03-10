package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/sudesh856/ecom-go-api-project/internal/adaptors/postgresql/sqlc"
	"github.com/sudesh856/ecom-go-api-project/internal/orders"
	"github.com/sudesh856/ecom-go-api-project/internal/products"
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


  productService := products.NewService(repo.New(app.db))
  productHandler := products.NewHandler(productService)
  r.Get("/products", productHandler.ListProducts)
  r.Get("/products/{id}", productHandler.FindProduct)
  r.Post("/products", productHandler.CreateProduct)
	// http.ListenAndServe(":3333", r)


	orderService := orders.NewService(repo.New(app.db), app.db)
  	ordersHandler := orders.NewHandler(orderService)
	r.Post("/orders", ordersHandler.PlaceOrder)
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

	log.Printf("Server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()


}


type application struct {
	config config
	db *pgx.Conn
}



type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}