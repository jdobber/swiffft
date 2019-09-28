package handlers

import (
	echo "github.com/labstack/echo/v4"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricHandler ...
func MetricsHandler() echo.HandlerFunc {

	fn := func(c echo.Context) error {
		h := promhttp.Handler()
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}

	return fn
}
