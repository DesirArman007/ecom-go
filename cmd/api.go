package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/DesirArman007/ecom/internal/adapters/postgresql/sqlc"
	"github.com/DesirArman007/ecom/internal/orders"
	"github.com/DesirArman007/ecom/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type application struct {
	config config
	db     *pgx.Conn
}

type config struct {
	addr string //:8080
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID) // imp for Ratelimiting,
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover from panics and returns a 500 error code

	// sets a timeout value on the request context (ctx), that will signal through
	// ctx.Done() that the work should be canceled if the request takes longer than
	// 60 seconds.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Working just fine"))
	})

	productService := products.NewService(repo.New(app.db)) // create a new instance of the product service
	productsHandler := products.NewHandler(productService)  // create a new instance of the product handler
	r.Get("/products", productsHandler.ListProducts)
	r.Get("/products/{id}", productsHandler.FindProductByID)
	r.Post("/products", productsHandler.CreateProduct)

	orderService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderService)
	r.Post("/orders", ordersHandler.PlaceOrder)
	r.Get("/orders/{id}", ordersHandler.GetOrderByID)

	return r
}

// run
func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}

	log.Printf("Server has atarted at addr %s", app.config.addr)

	return server.ListenAndServe()
}
