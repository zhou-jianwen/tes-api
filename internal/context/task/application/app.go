package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"code.byted.org/epscp/vetes-api/internal/apiserver/options"
	"code.byted.org/epscp/vetes-api/internal/context/task/application/command"
	"code.byted.org/epscp/vetes-api/internal/context/task/application/query"
	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
	"code.byted.org/epscp/vetes-api/internal/context/task/infra/normalize"
	"code.byted.org/epscp/vetes-api/internal/context/task/infra/persistence/sql"
	"code.byted.org/epscp/vetes-api/pkg/consts"
)

// TaskService ...
type TaskService struct {
	TaskCommands *command.Commands
	TaskQueries  *query.Queries
}

// NewTaskService ...
func NewTaskService(ctx context.Context, opts *options.Options) (*TaskService, error) {
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

	normalizer, err := normalize.NewNormalizer(opts.Normalize)
	if err != nil {
		return nil, err
	}
	svc := domain.NewService(repo, normalizer)
	taskCommands := command.NewCommands(svc)
	taskQueries := query.NewQueries(readModel)

	return &TaskService{
		TaskCommands: taskCommands,
		TaskQueries:  taskQueries,
	}, nil
}
