package query

import "context"

// ReadModel ...
type ReadModel interface {
	List(ctx context.Context, filter *ListFilter) ([]*ExtraPriority, error)
}
