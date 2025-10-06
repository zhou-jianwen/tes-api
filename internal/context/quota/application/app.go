package application

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/GBA-BI/tes-api/internal/apiserver/options"
	"github.com/GBA-BI/tes-api/internal/context/quota/application/command"
	"github.com/GBA-BI/tes-api/internal/context/quota/application/query"
	"github.com/GBA-BI/tes-api/internal/context/quota/domain"
	"github.com/GBA-BI/tes-api/internal/context/quota/infra/persistence/sql"
	"github.com/GBA-BI/tes-api/pkg/consts"
)

// QuotaService ...
type QuotaService struct {
	QuotaQueries  *query.Queries
	QuotaCommands *command.Commands
}

// NewQuotaService ...
func NewQuotaService(ctx context.Context, opts *options.Options) (*QuotaService, error) {
	var (
		err  error
		repo domain.Repo
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
	default:
		return nil, fmt.Errorf("unsupported db type")
	}

	svc := domain.NewService(repo)
	quotaQueries := query.NewQueries(svc)
	quotaCommands := command.NewCommands(svc)

	return &QuotaService{
		QuotaQueries:  quotaQueries,
		QuotaCommands: quotaCommands,
	}, nil
}
