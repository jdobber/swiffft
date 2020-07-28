package handlers

import (
	"net/http"

	echo "github.com/labstack/echo/v4"

	"github.com/jdobber/swiffft/lib/cmd"
	"github.com/jdobber/swiffft/lib/sources"
	"github.com/jdobber/swiffft/lib/support"

	iiifimage "github.com/jdobber/go-iiif-mod/lib/image"
	iiiflevel "github.com/jdobber/go-iiif-mod/lib/level"
)

// ImageHandler ...
func ImageHandler() echo.HandlerFunc {

	fn := func(c echo.Context) error {

		var body []byte
		var err error
		var sourceInfo *sources.SourceInfo
		var newSourceInfo *sources.SourceInfo

		url := c.Request().URL.String()

		p, _ := NewIIIFQueryParser(c)
		format, _ := p.GetIIIFParameter("format")

		// check cache for requested tile
		if cmd.Args.CacheOptions.Tiles {
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
				return c.Blob(http.StatusOK, "image/"+format, sourceInfo.Payload)
			}
		}

		identifier, _ := p.GetIIIFParameter("identifier")
		sourceInfo, err = sources.ReadFromSources(cmd.Sources, identifier)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		image, err := iiifimage.NewNativeImage(identifier, sourceInfo.Payload)
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

		newSourceInfo = &sources.SourceInfo{
			Payload:      data,
			LastModified: sourceInfo.LastModified,
			ETag:         support.GetETag(&data),
		}

		// set into cache
		if cmd.Args.CacheOptions.Tiles {
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

		// return image
		return c.Blob(http.StatusOK, "image/"+format, data)

	}

	return fn
}
