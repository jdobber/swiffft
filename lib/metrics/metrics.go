package metrics

// MetricsOptions ...
type MetricsOptions struct {
	EnableMetrics bool   `name:"metrics.enable" default:"false" desc:"Enable metrics."`
	Namespace     string `name:"metrics.namespace" default:"swiffft" desc:"Define a namespace for the metrics."`
}
