package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/GBA-BI/tes-api/internal/apiserver/options"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/application/command"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/application/query"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/domain"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/infra/persistence/sql"
	"github.com/GBA-BI/tes-api/pkg/consts"
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
