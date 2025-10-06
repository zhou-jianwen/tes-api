package options

import (
	"github.com/spf13/pflag"

	"github.com/GBA-BI/tes-api/pkg/log"

	"github.com/GBA-BI/tes-api/internal/context/task/infra/normalize"
	"github.com/GBA-BI/tes-api/pkg/db"
	"github.com/GBA-BI/tes-api/pkg/server"
)

// Options ...
type Options struct {
	Log       *log.Options       `mapstructure:"log"`
	Server    *server.Options    `mapstructure:"server"`
	DB        *db.Options        `mapstructure:"db"`
	Normalize *normalize.Options `mapstructure:"normalize"`
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		Log:       log.NewOptions(),
		Server:    server.NewOptions(),
		DB:        db.NewOptions(),
		Normalize: normalize.NewOptions(),
	}
}

// Validate ...
func (o *Options) Validate() error {
	if err := o.Log.Validate(); err != nil {
		return err
	}
	if err := o.Server.Validate(); err != nil {
		return err
	}
	if err := o.DB.Validate(); err != nil {
		return err
	}
	if err := o.Normalize.Validate(); err != nil {
		return err
	}
	return nil
}

// AddFlags ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.Log.AddFlags(fs)
	o.Server.AddFlags(fs)
	o.DB.AddFlags(fs)
	o.Normalize.AddFlags(fs)
}
