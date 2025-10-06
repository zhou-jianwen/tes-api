package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"code.byted.org/epscp/vetes-api/internal/apiserver/options"
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application/command"
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application/query"
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/domain"
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/infra/persistence/sql"
	"code.byted.org/epscp/vetes-api/pkg/consts"
)

// ExtraPriorityService ...
type ExtraPriorityService struct {
	ExtraPriorityCommands *command.Commands
	ExtraPriorityQueries  *query.Queries
}

// NewExtraPriorityService ...
func NewExtraPriorityService(ctx context.Context, opts *options.Options) (*ExtraPriorityService, error) {
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
	extraPriorityCommands := command.NewCommands(svc)
	extraPriorityQueries := query.NewQueries(readModel)

	return &ExtraPriorityService{
		ExtraPriorityCommands: extraPriorityCommands,
		ExtraPriorityQueries:  extraPriorityQueries,
	}, nil
}
