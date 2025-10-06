package server

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/spf13/pflag"
)

// HTTPOptions ...
type HTTPOptions struct {
	Port               uint16 `mapstructure:"port"`
	MetricsPort        uint16 `mapstructure:"metricsPort"`
	MaxRequestBodySize int    `mapstructure:"maxRequestBodySize"`
}

// NewHTTPOptions ...
func NewHTTPOptions() *HTTPOptions {
	return &HTTPOptions{
		Port:               8080,
		MetricsPort:        9090,
		MaxRequestBodySize: 128 * 1024 * 1024, // 128MiB
	}
}

// Validate ...
func (o *HTTPOptions) Validate() error {
	if o.Port == o.MetricsPort {
		return fmt.Errorf("port and metrics port cannot be the same")
	}
	return nil
}

// AddFlags ...
func (o *HTTPOptions) AddFlags(fs *pflag.FlagSet) {
	fs.Uint16Var(&o.Port, "http-port", o.Port, "http listen port")
	fs.Uint16Var(&o.MetricsPort, "http-metrics-port", o.MetricsPort, "http listen port for metrics")
	fs.IntVar(&o.MaxRequestBodySize, "http-max-request-body-size", o.MaxRequestBodySize, "http max request body size in bytes")
}

// RouteRegister ...
type RouteRegister interface {
	AddRoute(route.IRouter)
}
