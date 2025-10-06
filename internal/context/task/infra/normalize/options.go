package normalize

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

// Options ...
type Options struct {
	ExecutorBasePath string            `mapstructure:"executorBasePath"`
	DiskGB           DiskGBOptions     `mapstructure:"diskGB"`
	BootDiskGB       BootDiskGBOptions `mapstructure:"bootDiskGB"`
	GPU              GPUOptions        `mapstructure:"gpu"`
}

// DiskGBOptions ...
type DiskGBOptions struct {
	Enable    bool    `mapstrucutre:"enable"`
	Min       float64 `mapstrucutre:"min"`
	Max       float64 `mapstrucutre:"max"`
	IsInteger bool    `mapstrucutre:"isInteger"`
}

// BootDiskGBOptions ...
type BootDiskGBOptions struct {
	Enable bool `mapstrucutre:"enable"`
	Min    int  `mapstrucutre:"min"`
	Max    int  `mapstrucutre:"max"`
}

// GPUOptions ...
type GPUOptions struct {
	Enable    bool `mapstrucutre:"enable"`
	IsInteger bool `mapstrucutre:"isInteger"`
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		ExecutorBasePath: "/cromwell-executions/",
		DiskGB: DiskGBOptions{
			Enable:    true,
			Min:       20,
			Max:       8192,
			IsInteger: true,
		},
		BootDiskGB: BootDiskGBOptions{
			Enable: true,
			Min:    40,
			Max:    100,
		},
		GPU: GPUOptions{
			Enable:    true,
			IsInteger: true,
		},
	}
}

// Validate ...
func (o *Options) Validate() error {
	if !filepath.IsAbs(o.ExecutorBasePath) {
		return fmt.Errorf("executorBasePath %s should be absolute path", o.ExecutorBasePath)
	}
	if !strings.HasSuffix(o.ExecutorBasePath, "/") {
		return fmt.Errorf("executorBasePath %s should ends with slash", o.ExecutorBasePath)
	}
	if o.DiskGB.Enable {
		if o.DiskGB.Min > o.DiskGB.Max {
			return fmt.Errorf("normalize diskGB min shoould be less than max")
		}
		if o.DiskGB.Min < 0 {
			return fmt.Errorf("normalize diskGB min should be positive")
		}
		if o.DiskGB.Max < 0 {
			return fmt.Errorf("normalize diskGB max should be positive")
		}
	}
	if o.BootDiskGB.Enable {
		if o.BootDiskGB.Min > o.BootDiskGB.Max {
			return fmt.Errorf("normalize bootDiskGB min shoould be less than max")
		}
		if o.BootDiskGB.Min < 0 {
			return fmt.Errorf("normalize bootDiskGB min should be positive")
		}
		if o.BootDiskGB.Max < 0 {
			return fmt.Errorf("normalize bootDiskGB max should be positive")
		}
	}
	return nil
}

// AddFlags ...
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ExecutorBasePath, "normalize-executor-base-path", o.ExecutorBasePath, "check executor base path")
	fs.BoolVar(&o.DiskGB.Enable, "normalize-diskgb-enable", o.DiskGB.Enable, "enable normalize disk in gb")
	fs.Float64Var(&o.DiskGB.Min, "normalize-diskgb-min", o.DiskGB.Min, "normalize disk min in gb")
	fs.Float64Var(&o.DiskGB.Max, "normalize-diskgb-max", o.DiskGB.Max, "normalize disk max in gb")
	fs.BoolVar(&o.DiskGB.IsInteger, "normalize-diskgb-integer", o.DiskGB.IsInteger, "normalize disk in gb as integer")
	fs.BoolVar(&o.BootDiskGB.Enable, "normalize-bootdiskgb-enable", o.BootDiskGB.Enable, "enable normalize bootDisk in gb")
	fs.IntVar(&o.BootDiskGB.Min, "normalize-bootdiskgb-min", o.BootDiskGB.Min, "normalize bootDisk min in gb")
	fs.IntVar(&o.BootDiskGB.Max, "normalize-bootdiskgb-max", o.BootDiskGB.Max, "normalize bootDisk max in gb")
	fs.BoolVar(&o.GPU.Enable, "normalize-gpu-enable", o.GPU.Enable, "enable normalize gpu")
	fs.BoolVar(&o.GPU.IsInteger, "normalize-gpu-integer", o.GPU.IsInteger, "normalize gpu count as integer")
}
