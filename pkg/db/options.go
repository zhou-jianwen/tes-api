package db

import (
	"errors"

	"github.com/spf13/pflag"

	"code.byted.org/epscp/vetes-api/pkg/consts"
)

// Options ...
type Options struct {
	Type  string        `mapstrucure:"type"`
	MySQL *MySQLOptions `mapstructure:"mysql"`
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		Type:  consts.MySQLType,
		MySQL: NewMySQLOptions(),
	}
}

// Validate ...
func (o *Options) Validate() error {
	switch o.Type {
	case consts.MySQLType:
		if err := o.MySQL.Validate(); err != nil {
			return err
		}
	default:
		return errors.New("invalid db type")
	}

	return nil
}

// AddFlags ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Type, "db-type", o.Type, "db type")
	o.MySQL.AddFlags(fs)
}
