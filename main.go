package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	iiifconfig "github.com/jdobber/go-iiif-mod/lib/config"
	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifparser "github.com/jdobber/go-iiif-mod/lib/parser"
	iiifprofile "github.com/jdobber/go-iiif-mod/lib/profile"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
	"github.com/whosonfirst/go-sanitize"

	caches "github.com/jdobber/swiffft/lib/caches"
	"github.com/jdobber/swiffft/lib/cmd"
	sources "github.com/jdobber/swiffft/lib/sources"
)

var (
	Config *iiifconfig.Config
	Source sources.Source
	Cache  caches.Cache
	Args   cmd.CommandOptions
)

func check(e error) {
	if e != nil {
		log.Fatalln(e)
		//panic(e)
	}
}

func main() {
	var err error

	Args = cmd.Init()

	if Args.Config == "" {
		log.Fatal("Missing config file")
	}

	Config, err = iiifconfig.NewConfigFromFlag(Args.Config)
	check(err)

	// Source, err = sources.NewFileSource("/home/jens/Bilder")
	Source, err = sources.NewMinioSource(Args.MinioOptions)
	check(err)

	Cache, err = caches.NewFastCache()
	check(err)

	/*
		INIT SERVER AND ROUTES
	*/
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${method}] ${status} ${uri} ${latency_human} ${bytes_out} \n",
		Output: os.Stdout,
	}))
	e.Logger.SetLevel(echolog.INFO)

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "pong")
	})

	//e.GET("/debug/vars", expvar_handler)

	g := e.Group("/iiif")
	g.GET("/:identifier/:region/:size/:rotation/:quality", ImageHandler())
	g.GET("/:identifier/info.json", InfoHandler())

	e.Logger.Fatal(e.Start(Args.Host))
}

func NewIIIFQueryParser(c echo.Context) (*iiifparser.IIIFQueryParser, error) {

	opts := sanitize.DefaultOptions()
	vars := make(map[string]string)

	for _, name := range c.ParamNames() {
		if name == "quality" {
			a := strings.Split(c.Param(name), ".")
			vars["quality"] = a[0]
			vars["format"] = a[1]
		} else {
			vars[name] = c.Param(name)
		}

	}

	p := iiifparser.IIIFQueryParser{
		Opts: opts,
		Vars: vars,
	}

	return &p, nil
}

func InfoHandler() echo.HandlerFunc {

	fn := func(c echo.Context) error {

		var body []byte
		var err error

		p, _ := NewIIIFQueryParser(c)
		identifier, _ := p.GetIIIFParameter("identifier")

		if Cache.Has(identifier) {

			body, err = Cache.Get(identifier)
			check(err)

		} else {

			body, err = Source.Read(identifier)
			check(err)

			err = Cache.Set(identifier, body)
			check(err)
		}

		image, err := iiifimage.NewNativeImage(identifier, body)
		check(err)

		level, err := iiiflevel.NewLevelFromConfig(Config, &Args.Endpoint)
		check(err)

		profile, err := iiifprofile.NewProfile(&Args.Endpoint, image, level)
		check(err)

		return c.JSON(http.StatusOK, profile)
	}

	return fn
}

func ImageHandler() echo.HandlerFunc {

	fn := func(c echo.Context) error {

		var body []byte
		var err error

		p, _ := NewIIIFQueryParser(c)
		identifier, _ := p.GetIIIFParameter("identifier")
		format, _ := p.GetIIIFParameter("format")

		if Cache.Has(identifier) {
			body, err = Cache.Get(identifier)
			check(err)
		} else {
			body, err = Source.Read(identifier)
			check(err)

			err = Cache.Set(identifier, body)
			check(err)
		}

		image, err := iiifimage.NewNativeImage(identifier, body)
		check(err)

		level, err := iiiflevel.NewLevelFromConfig(Config, &Args.Endpoint)
		check(err)

		iiifparams, err := p.GetIIIFParameters()
		check(err)

		transformation, err := iiifimage.NewTransformation(level,
			iiifparams.Region,
			iiifparams.Size,
			iiifparams.Rotation,
			iiifparams.Quality,
			iiifparams.Format)
		check(err)

		if transformation.HasTransformation() {

			_, err := image.Transform(transformation)
			check(err)

		}

		opts := iiifimage.EncodingOptions{
			Format:  format,
			Quality: 70,
		}

		data, err := image.Encode(&opts)
		check(err)

		return c.Blob(http.StatusOK, "image/"+format, data)

	}

	return fn
}
