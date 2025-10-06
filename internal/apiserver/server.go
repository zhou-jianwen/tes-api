package apiserver

import (
	"context"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	_ "github.com/GBA-BI/tes-api/docs" // for swagger
	"github.com/GBA-BI/tes-api/internal/apiserver/options"
	clusterapp "github.com/GBA-BI/tes-api/internal/context/cluster/application"
	clusterhertz "github.com/GBA-BI/tes-api/internal/context/cluster/interface/hertz"
	extrapriorityapp "github.com/GBA-BI/tes-api/internal/context/extrapriority/application"
	extrapriorityhertz "github.com/GBA-BI/tes-api/internal/context/extrapriority/interface/hertz"
	quotaapp "github.com/GBA-BI/tes-api/internal/context/quota/application"
	quotahertz "github.com/GBA-BI/tes-api/internal/context/quota/interface/hertz"
	taskapp "github.com/GBA-BI/tes-api/internal/context/task/application"
	taskhertz "github.com/GBA-BI/tes-api/internal/context/task/interface/hertz"
	"github.com/GBA-BI/tes-api/pkg/version"
	"github.com/GBA-BI/tes-api/pkg/viper"
)

const (
	component = "veTES-api"
)

func newServerCommand(ctx context.Context, opts *options.Options) *cobra.Command {
	return &cobra.Command{
		Use:          component,
		Short:        "veTES api server",
		Long:         "veTES api server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			version.PrintVersionOrContinue()

			if err := opts.Validate(); err != nil {
				return err
			}

			applog.RegisterLogger(opts.Log)
			defer applog.Sync()

			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				applog.Infow("FLAG", flag.Name, flag.Value)
			})

			return run(ctx, opts)
		},
	}
}

func run(ctx context.Context, opts *options.Options) (err error) {
	applog.Infow("run veTES api server")

	taskService, err := taskapp.NewTaskService(ctx, opts)
	if err != nil {
		return err
	}
	clusterService, err := clusterapp.NewClusterService(ctx, opts)
	if err != nil {
		return err
	}
	quotaService, err := quotaapp.NewQuotaService(ctx, opts)
	if err != nil {
		return err
	}
	extraPriorityService, err := extrapriorityapp.NewExtraPriorityService(ctx, opts)
	if err != nil {
		return err
	}

	httpServer := setupHTTPServer(opts.Server.HTTP,
		taskhertz.NewRouterRegister(taskService),
		clusterhertz.NewRouterRegister(clusterService),
		quotahertz.NewRouterRegister(quotaService),
		extrapriorityhertz.NewRouterRegister(extraPriorityService),
	)

	httpServer.Spin()
	return nil
}

// NewServerCommand create a veTES api server command.
func NewServerCommand(ctx context.Context) (*cobra.Command, error) {
	opts := options.NewOptions()
	cmd := newServerCommand(ctx, opts)

	opts.AddFlags(cmd.Flags())
	version.AddFlags(cmd.Flags())
	cmd.Flags().AddFlag(pflag.Lookup(viper.ConfigFlagName))
	if err := viper.LoadConfig(opts); err != nil {
		return nil, err
	}
	return cmd, nil
}
