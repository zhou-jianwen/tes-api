package query

import (
	"context"

	"code.byted.org/epscp/vetes-api/pkg/utils"
)

// ReadModel ...
type ReadModel interface {
	ListMinimal(ctx context.Context, pageSize int, pageToken *utils.PageToken, filter *ListFilter) ([]*TaskMinimal, *utils.PageToken, error)
	ListBasic(ctx context.Context, pageSize int, pageToken *utils.PageToken, filter *ListFilter) ([]*TaskBasic, *utils.PageToken, error)
	ListFull(ctx context.Context, pageSize int, pageToken *utils.PageToken, filter *ListFilter) ([]*Task, *utils.PageToken, error)
	GetMinimal(ctx context.Context, id string) (*TaskMinimal, error)
	GetBasic(ctx context.Context, id string) (*TaskBasic, error)
	GetFull(ctx context.Context, id string) (*Task, error)
	GatherResources(ctx context.Context, filter *GatherFilter) (*TasksResources, error)
	ListAccounts(ctx context.Context) ([]*AccountInfo, error)
}
