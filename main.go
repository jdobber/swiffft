package main

import (
	"log"
	"net/http"
	"os"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"

	caches "github.com/jdobber/swiffft/lib/caches"
	"github.com/jdobber/swiffft/lib/cmd"
	"github.com/jdobber/swiffft/lib/handlers"
	sources "github.com/jdobber/swiffft/lib/sources"
)

func main() {
	var err error

	cmd.Args = cmd.Init()

	if cmd.Args.Config == "" {
		log.Fatal("Missing config file")
	}

	cmd.Config, err = iiifconfig.NewConfigFromFlag(cmd.Args.Config)
	cmd.Check(err)

	/*
		INIT SOURCES
	*/
	for _, source := range cmd.Args.Sources {
		switch source {
		case "minio":
			s, err := sources.NewMinioSource(cmd.Args.MinioOptions)
			cmd.Check(err)
			cmd.Sources = append(cmd.Sources, s)
		case "file":
			s, err := sources.NewFileSource(cmd.Args.FileSourceOptions.Prefix)
			cmd.Check(err)
			cmd.Sources = append(cmd.Sources, s)
		}
	}

	cmd.Cache, err = caches.NewFastCache()
	cmd.Check(err)

	/*
		INIT SERVER AND ROUTES
	*/
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${method}] ${status} ${uri} ${latency_human} ${bytes_out} \n",
		Output: os.Stdout,
	}))
	e.HTTPErrorHandler = handlers.ErrorHandler
	e.Logger.SetLevel(echolog.DEBUG)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "pong")
	})

	//e.GET("/debug/vars", expvar_handler)

	g := e.Group("/iiif")
	g.GET("/:identifier/:region/:size/:rotation/:quality", handlers.ImageHandler())
	g.GET("/:identifier/info.json", handlers.InfoHandler())

	e.Logger.Fatal(e.Start(cmd.Args.Host))
}
