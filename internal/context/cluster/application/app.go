package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"code.byted.org/epscp/vetes-api/internal/apiserver/options"
	"code.byted.org/epscp/vetes-api/internal/context/cluster/application/command"
	"code.byted.org/epscp/vetes-api/internal/context/cluster/application/query"
	"code.byted.org/epscp/vetes-api/internal/context/cluster/domain"
	"code.byted.org/epscp/vetes-api/internal/context/cluster/infra/persistence/sql"
	"code.byted.org/epscp/vetes-api/pkg/consts"
)

// ClusterService ...
type ClusterService struct {
	ClusterCommands *command.Commands
	ClusterQueries  *query.Queries
}

// NewClusterService ...
func NewClusterService(ctx context.Context, opts *options.Options) (*ClusterService, error) {
	var (
		err       error
		repo      domain.Repo
		readModel query.ReadModel
	)

	switch opts.DB.Type {
	case consts.MySQLType:
		var db *gorm.DB
		if db, err = opts.DB.MySQL.GetGORMInstance(); err != nil {
			return nil, err
		}
		if repo, err = sql.NewRepo(ctx, db); err != nil {
			return nil, err
		}
		if readModel, err = sql.NewReadModel(ctx, db); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported db type")
	}

	svc := domain.NewService(repo)
	clusterCommands := command.NewCommands(svc)
	clusterQueries := query.NewQueries(readModel)

	return &ClusterService{
		ClusterCommands: clusterCommands,
		ClusterQueries:  clusterQueries,
	}, nil
}
