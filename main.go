// Package classification Todos API
//
// Documentation for Todos API
//
// Schemes : http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	handlers "github.com/Firasso/DemoGoApi/handlers"
	custom_middleware "github.com/Firasso/DemoGoApi/middleware"
	todosTypes "github.com/Firasso/DemoGoApi/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"database/sql"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	var err error

	initDB := os.Getenv("INIT_DB") == "true"
	dbLink := os.Getenv("DB_URL")

	DB, err = sql.Open("postgres", dbLink)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("DB connected")
	}

	if initDB {
		content, err := ioutil.ReadFile("./sql/todos.sql")
		if err != nil {
			panic(err)
		}

		todosTable := string(content)

		_, err = DB.Exec(todosTable)
		if err != nil {
			panic(err)
		}

		var todos = todosTypes.Todos{}

		_, err = todos.AddTodo(todosTypes.NewTodo{Text: "first todo", Completed: true}, DB)

		if err != nil {
			panic(err)
		}
		_, err = todos.AddTodo(todosTypes.NewTodo{Text: "second todo", Completed: true}, DB)
		if err != nil {
			panic(err)
		}

		fmt.Println("DB initialized")
	}
}

var Router = chi.NewRouter()

func SetupRoutes() {
	Router.Route("/docs", func(r chi.Router) {
		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			swaggerURL, _ := url.Parse("http://swagger:8080")
			proxy := httputil.NewSingleHostReverseProxy(swaggerURL)
			proxy.ServeHTTP(w, req)
		})
	})

	Router.Handle("/metrics", promhttp.Handler())

	Router.Route("/todos", func(r chi.Router) {
		r.Get("/", handlers.GetTodos)
		r.Post("/", handlers.AddTodo)

		r.Route("/{todoID}", func(r chi.Router) {
			r.Get("/", handlers.GetSingleTodo)
			r.Put("/", handlers.UpdateTodo)
			r.Delete("/", handlers.DeleteTodo)
		})
	})
}

func SetupMiddleware(db *sql.DB) {
	Router.Use(middleware.Logger)
	Router.Use(custom_middleware.PassDBContext(db))
}

func main() {
	defer DB.Close()
	SetupMiddleware(DB)
	SetupRoutes()

	if err := http.ListenAndServe(":3000", Router); err != http.ErrServerClosed {
		panic(err)
	}
}
