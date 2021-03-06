package main

import (
	"log"
	"net/http"
	"os"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"

	"github.com/labstack/echo-contrib/prometheus"
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

	/*
		INIT CACHE
	*/
	if cmd.Args.CacheOptions.UseCache {
		cmd.Cache, err = caches.NewFastCache(cmd.Args.CacheOptions.Size)
		log.Printf("Use FastCache with size %d MB.\n", cmd.Args.CacheOptions.Size)
	} else {
		cmd.Cache, err = caches.NewNullCache()
	}
	cmd.Check(err)

	/*
		INIT SERVER AND ROUTES
	*/
	e := echo.New()

	if cmd.Args.MetricsOptions.EnableMetrics {
		log.Printf("Enable metrics with namespace=%s.\n", cmd.Args.MetricsOptions.Namespace)
		p := prometheus.NewPrometheus(cmd.Args.MetricsOptions.Namespace, nil)
		p.Use(e)
	}

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

	// IIIF Routes
	g := e.Group("/iiif")
	g.GET("/:identifier/:region/:size/:rotation/:quality", handlers.ImageHandler())
	g.GET("/:identifier/info.json", handlers.InfoHandler())

	e.Logger.Fatal(e.Start(cmd.Args.Host))
}
