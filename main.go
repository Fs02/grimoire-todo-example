package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/mysql"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func dsn() string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))
}

func main() {
	port := os.Getenv("PORT")

	// initialize mysql adapter.
	adapter, err := mysql.Open(dsn())
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize grimoire's repo.
	repo := grimoire.New(adapter)
	resource := Resource{Repo: repo}

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	// List all todos
	r.Get("/", resource.Index)
	r.With(resource.BodyParser).Post("/", resource.Create)
	r.With(resource.Load).Get("/{ID}", resource.Show)
	r.With(resource.BodyParser, resource.Load).Patch("/{ID}", resource.Update)
	r.With(resource.Load).Delete("/{ID}", resource.Delete)
	r.Delete("/", resource.Clear)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
