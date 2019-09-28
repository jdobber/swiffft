package handlers

import (
	"encoding/json"
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

		url := c.Request().URL.String()

		p, _ := NewIIIFQueryParser(c)

		// check cache for requested info.json
		if cmd.Args.CacheOptions.InfoJSON {
			body, err = cmd.Cache.Get(url)
			if err == nil {
				return c.Blob(http.StatusOK, "application/json", body)
			}
		}

		// not in cache -> go on
		identifier, _ := p.GetIIIFParameter("identifier")

		body, err = cmd.Cache.Get(identifier)
		if err != nil {

			body, err = sources.ReadFromSources(cmd.Sources, identifier)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

		profile, err := iiifprofile.NewProfile(cmd.Args.Endpoint, image, level)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		data, err := json.Marshal(profile)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		// set into cache
		if cmd.Args.CacheOptions.InfoJSON {

			err = cmd.Cache.Set(url, data)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "could not set item into cache")
			}
		}

		// return info.json
		return c.Blob(http.StatusOK, "application/json", data)
	}

	return fn
}
