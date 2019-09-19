package handlers

import (
	"net/http"

	echo "github.com/labstack/echo/v4"

	"github.com/jdobber/swiffft/lib/cmd"
	"github.com/jdobber/swiffft/lib/sources"

	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifprofile "github.com/jdobber/go-iiif-mod/lib/profile"
)

// InfoHandler ...
func InfoHandler() echo.HandlerFunc {

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

		profile, err := iiifprofile.NewProfile(&cmd.Args.Endpoint, image, level)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, profile)
	}

	return fn
}
