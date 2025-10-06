package query

import (
	"context"

	"github.com/GBA-BI/tes-api/internal/context/quota/domain"
	"github.com/GBA-BI/tes-api/pkg/validator"
)

// GetQuery ...
type GetQuery struct {
	Global    bool
	AccountID string
	UserID    string
}

func (q *GetQuery) setDefault() {}

func (q *GetQuery) validate() error {
	return validator.Validate(q)
}

// GetHandler ...
type GetHandler interface {
	Handle(ctx context.Context, query *GetQuery) (*Quota, error)
}

type getHandler struct {
	svc domain.Service
}

var _ GetHandler = (*getHandler)(nil)

// NewGetHandler ...
func NewGetHandler(svc domain.Service) GetHandler {
	return &getHandler{svc: svc}
}

// Handle ...
func (h *getHandler) Handle(ctx context.Context, query *GetQuery) (*Quota, error) {
	query.setDefault()
	if err := query.validate(); err != nil {
		return nil, err
	}
	quota, err := h.svc.GetOrDefault(ctx, query.Global, query.AccountID, query.UserID)
	if err != nil {
		return nil, err
	}
	return &Quota{
		Global:        quota.IsGlobal(),
		Default:       quota.IsDefault(),
		AccountID:     query.AccountID,
		UserID:        query.UserID,
		ResourceQuota: resourceQuotaDOToDTO(quota.ResourceQuota),
	}, nil
}
