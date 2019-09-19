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

		p, _ := NewIIIFQueryParser(c)
		identifier, _ := p.GetIIIFParameter("identifier")

		if cmd.Cache.Has(identifier) {

			body, err = cmd.Cache.Get(identifier)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "could not get item from cache")
			}
		} else {

			body, err = sources.ReadFromSources(cmd.Sources, identifier)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			err = cmd.Cache.Set(identifier, body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "could not set item into cache")
			}
		}

		image, err := iiifimage.NewNativeImage(identifier, body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "could not load image")
		}

		level, err := iiiflevel.NewLevelFromConfig(cmd.Config, &cmd.Args.Endpoint)
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

		format, _ := p.GetIIIFParameter("format")
		opts := iiifimage.EncodingOptions{
			Format:  format,
			Quality: 70,
		}

		data, err := image.Encode(&opts)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.Blob(http.StatusOK, "image/"+format, data)

	}

	return fn
}
