package support

import (
	"fmt"
	"hash/fnv"

	"github.com/labstack/echo/v4"
)

// GetETag  ...
func GetETag(payload *[]byte) string {
	// build ETag
	h := fnv.New128a()
	h.Write(*payload)
	return fmt.Sprintf("\"%x\"", h.Sum(nil))
}

// SetETagHeader ...
func SetETagHeader(c echo.Context, etag string, lastModified string) {
	c.Response().Header().Set("ETag", etag)
	c.Response().Header().Set(echo.HeaderLastModified, lastModified)
}

// CheckETagHeader ...
func CheckETagHeader(c echo.Context) bool {
	return c.Request().Header.Get("If-None-Match") == c.Response().Header().Get("ETag")
}
