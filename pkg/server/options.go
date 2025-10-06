package server

import "github.com/spf13/pflag"

// Options ...
type Options struct {
	HTTP *HTTPOptions `mapstructure:"http"`
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		HTTP: NewHTTPOptions(),
	}
}

// Validate ...
func (o *Options) Validate() error {
	if err := o.HTTP.Validate(); err != nil {
		return err
	}
	return nil
}

// AddFlags  ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.HTTP.AddFlags(fs)
}
