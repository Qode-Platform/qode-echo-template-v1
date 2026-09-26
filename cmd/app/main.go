// A minimal Echo service shaped for the fleet: it serves at the root of its own
// hostname, so routes mount directly on the Echo instance.
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func port() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return p
	}
	return "8080"
}

func newServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"service": "echo-template"})
	})
	return e
}

func main() {
	addr := ":" + port()
	log.Printf("echo-template listening on %s", addr)
	if err := newServer().Start(addr); err != nil {
		log.Fatal(err)
	}
}
