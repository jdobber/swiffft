package metrics

// see also: https://github.com/brancz/prometheus-example-app/blob/master/main.go

import (
	"time"
	"log"
	"fmt"

	echo "github.com/labstack/echo/v4"

	"github.com/prometheus/client_golang/prometheus"
	//"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpReqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	registry *prometheus.Registry
)

func RecordMetrics() {
	go func() {
		for {
			//RequestsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)	
	log.Printf("%s took %s ", name, elapsed)
}

func init() {
	registry = prometheus.NewRegistry()
	registry.MustRegister(httpReqs)
	//RecordMetrics()
}

// Handler ...
func Handler() echo.HandlerFunc {

	fn := func(c echo.Context) error {
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}

	return fn
}

// HandlerWrapper ...
func HandlerWrapper(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if err := next(c); err != nil {
			c.Error(err)
		}

		// exit if metrics route
		if c.Request().URL.Path == "/metrics" {
			return nil
		}

		//defer timeTrack(time.Now(), "xxx")		
		httpReqs.WithLabelValues(fmt.Sprintf("%d", c.Response().Status), c.Request().Method).Inc()
		
		return nil
	}
}
