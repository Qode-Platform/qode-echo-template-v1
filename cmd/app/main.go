// A minimal Echo service shaped for the fleet: one group carries the whole app
// so the ingress prefix is applied in a single place.
package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// basePath returns the fleet's ingress prefix, normalised to "" or
// "/leading/no-trailing-slash". The fleet injects BASE_PATH as
// /direct/<agent>:<port> and nginx forwards it UNCHANGED, so every route must
// live under it. Empty means standalone: serve at the host root.
func basePath() string {
	raw := strings.Trim(strings.TrimSpace(os.Getenv("BASE_PATH")), "/")
	if raw == "" {
		return ""
	}
	return "/" + raw
}

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

	g := e.Group(basePath())
	g.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	g.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"service": "echo-template", "base_path": basePath()})
	})
	return e
}

func main() {
	addr := ":" + port()
	log.Printf("echo-template listening on %s (base_path=%q)", addr, basePath())
	if err := newServer().Start(addr); err != nil {
		log.Fatal(err)
	}
}
