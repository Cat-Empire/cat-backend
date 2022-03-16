package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
  "os"

	"github.com/Cat-Empire/cat-backend/ent"
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

  // RUn the migration here
  if err := client.Schema.Create(context.Background()); !errors.Is(err, nil) {
    log.Fatalf("Error: failed creating schema resources %v\n", err)
  }
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
