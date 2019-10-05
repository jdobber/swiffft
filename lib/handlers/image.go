package handlers

import (
	"net/http"

	echo "github.com/labstack/echo/v4"

	"github.com/jdobber/swiffft/lib/cmd"
	"github.com/jdobber/swiffft/lib/sources"

	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
)

// ImageHandler ...
func ImageHandler() echo.HandlerFunc {

	fn := func(c echo.Context) error {

		var body []byte
		var err error

		url := c.Request().URL.String()

		p, _ := NewIIIFQueryParser(c)
		format, _ := p.GetIIIFParameter("format")

		// check cache for requested tile
		if cmd.Args.CacheOptions.Tiles {
			body, err = cmd.Cache.Get(url)
			if err == nil {
				return c.Blob(http.StatusOK, "image/"+format, body)
			}
		}

		// not in cache -> go on
		identifier, _ := p.GetIIIFParameter("identifier")

		// check cache for image
		body, err = cmd.Cache.Get(identifier)
		if err != nil {

			body, err = sources.ReadFromSources(cmd.Sources, identifier)
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}

			if cmd.Args.CacheOptions.Images {
				err = cmd.Cache.Set(identifier, body)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, "could not set item into cache")
				}
			}
		}

		image, err := iiifimage.NewNativeImage(identifier, body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not load image")
		}

		level, err := iiiflevel.NewLevelFromConfig(cmd.Config, cmd.Args.Endpoint)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		iiifparams, err := p.GetIIIFParameters()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		transformation, err := iiifimage.NewTransformation(level,
			iiifparams.Region,
			iiifparams.Size,
			iiifparams.Rotation,
			iiifparams.Quality,
			iiifparams.Format)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if transformation.HasTransformation() {

			_, err := image.Transform(transformation)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

		}

		opts := iiifimage.EncodingOptions{
			Format:  format,
			Quality: 70,
		}

		data, err := image.Encode(&opts)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// set into cache
		if cmd.Args.CacheOptions.Tiles {
			err = cmd.Cache.Set(url, data)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "could not set item into cache")
			}
		}

		// return image
		return c.Blob(http.StatusOK, "image/"+format, data)

	}

	return fn
}
