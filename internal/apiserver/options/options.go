package options

import (
	"github.com/spf13/pflag"

	"code.byted.org/epscp/go-common/log"

	"code.byted.org/epscp/vetes-api/internal/context/task/infra/normalize"
	"code.byted.org/epscp/vetes-api/pkg/db"
	"code.byted.org/epscp/vetes-api/pkg/server"
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
