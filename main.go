package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Cat-Empire/cat-backend/ent"
	"github.com/Cat-Empire/cat-backend/graph"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	host     = getEnv("DATABASE_HOST", "localhost")
	port     = getEnv("DATABASE_PORT", "5432")
	user     = getEnv("DATABASE_USERNAME", "catempire")
	password = getEnv("DATABASE_PASSWORD", "catempire123")
	dbname   = getEnv("DATABASE_NAME", "cats")
)

func main() {
	e := echo.New()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	client, err := ent.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error: posgres: %v\n", err)
	}
	defer client.Close()

	// Run the migration here
	if err := client.Schema.Create(context.Background()); !errors.Is(err, nil) {
		log.Fatalf("Error: failed creating schema resources %v\n", err)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	srv := handler.NewDefaultServer(graph.NewSchema(client))
	{
		e.POST("/users", func(c echo.Context) error {
			// Create an article entity
			a, err := client.Post.Create().
				SetTitle("title 1").
				SetDescription("description 1").
				Save(c.Request().Context())
			if !errors.Is(err, nil) {
				log.Fatalf("Error: failed creating article %v\n", err)
			}

			u, err := client.User.
				Create().
				SetName("Bob").
				SetAge(21).
				AddPosts(a). // Add article to the user
				Save(c.Request().Context())

			if !errors.Is(err, nil) {
				log.Fatalf("Error: failed creating user %v\n", err)
			}

			return c.JSON(http.StatusCreated, u)
		})
		e.POST("/query", func(c echo.Context) error {
			srv.ServeHTTP(c.Response(), c.Request())
			return nil
		})

		e.GET("/playground", func(c echo.Context) error {
			playground.Handler("GraphQL", "/query").ServeHTTP(c.Response(), c.Request())
			return nil
		})

	}
	e.Logger.Fatal(e.Start(":8080"))
}
