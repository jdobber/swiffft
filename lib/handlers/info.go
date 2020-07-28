package handlers

import (
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"

	"github.com/jdobber/swiffft/lib/cmd"
	"github.com/jdobber/swiffft/lib/sources"
	"github.com/jdobber/swiffft/lib/support"

	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
	iiifprofile "github.com/jdobber/go-iiif-mod/lib/profile"
)

// InfoHandler ...
func InfoHandler() echo.HandlerFunc {

	fn := func(c echo.Context) error {

		var body []byte
		var err error
		var sourceInfo *sources.SourceInfo
		var newSourceInfo *sources.SourceInfo

		url := c.Request().URL.String()

		p, _ := NewIIIFQueryParser(c)

		// check cache for requested info.json
		if cmd.Args.CacheOptions.InfoJSON {
			body, err = cmd.Cache.Get(url)
			if err == nil {
				sourceInfo, err = sources.Unpack(body)
				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
				support.SetETagHeader(c, sourceInfo.ETag, sourceInfo.LastModified)
				if support.CheckETagHeader(c) {
					return c.NoContent(http.StatusNotModified)
				}
				return c.JSONBlob(http.StatusOK, sourceInfo.Payload)
			}
		}

		identifier, _ := p.GetIIIFParameter("identifier")
		sourceInfo, err = sources.ReadFromSources(cmd.Sources, identifier)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		image, err := iiifimage.NewNativeImage(identifier, sourceInfo.Payload)
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

		newSourceInfo = &sources.SourceInfo{
			Payload:      data,
			LastModified: sourceInfo.LastModified,
			ETag:         support.GetETag(&data),
		}

		// set into cache
		if cmd.Args.CacheOptions.InfoJSON {
			packed, _ := sources.Pack(newSourceInfo)
			err = cmd.Cache.Set(url, packed)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "could not set item into cache")
			}
		}

		// return NOT MODIFIED if etag matches
		support.SetETagHeader(c, newSourceInfo.ETag, newSourceInfo.LastModified)
		if support.CheckETagHeader(c) {
			return c.NoContent(http.StatusNotModified)
		}

		// return info.json
		return c.JSONBlob(http.StatusOK, data)
	}

	return fn
}
