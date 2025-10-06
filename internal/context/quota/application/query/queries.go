package query

import "code.byted.org/epscp/vetes-api/internal/context/quota/domain"

// Queries ...
type Queries struct {
	Get GetHandler
}

// NewQueries ...
func NewQueries(svc domain.Service) *Queries {
	return &Queries{Get: NewGetHandler(svc)}
}
