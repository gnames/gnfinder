package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gnames/gnfinder"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const withLogs = true

//go:embed static
var static embed.FS

// Run starts GNfinder service for its webiste and RESTful API.
func Run(gnf gnfinder.GNfinder, port int) {
	log.Printf("Starting the HTTP API server on port %d.", port)
	e := echo.New()

	var err error
	e.Renderer, err = NewTemplate()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	if withLogs {
		e.Use(middleware.Logger())
	}

	e.GET("/", home(gnf))
	e.GET("/apidoc", apidoc(gnf))
	e.POST("/find", find(gnf))

	e.GET("/api/v1/ping", pingAPI())
	e.GET("/api/v1/version", verAPI(gnf))
	e.POST("/api/v1/find", findAPI(gnf))

	fs := http.FileServer(http.FS(static))
	e.GET("/static/*", echo.WrapHandler(fs))

	addr := fmt.Sprintf(":%d", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
