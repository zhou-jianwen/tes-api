package domain

import (
	"fmt"

	"code.byted.org/epscp/vetes-api/pkg/consts"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
)

// NewQuotaID ...
func NewQuotaID(global bool, accountID, userID string) (string, error) {
	if global {
		if accountID != "" || userID != "" {
			return "", apperrors.NewInvalidError("global with non-empty account_id or user_id")
		}
		return consts.GlobalQuotaID, nil
	}

	if accountID == "" {
		return "", apperrors.NewInvalidError("empty account_id")
	}

	if accountID == consts.DefaultQuotaAccountID {
		if userID != "" {
			return "", apperrors.NewInvalidError("default account_id with non-empty user_id")
		}
		return consts.DefaultQuotaAccountID, nil
	}

	return fmt.Sprintf("%s/%s", accountID, userID), nil
}

// Quota ...
type Quota struct {
	ID            string
	AccountID     string
	UserID        string
	ResourceQuota *ResourceQuota
}

// ResourceQuota ...
type ResourceQuota struct {
	Count    *int
	CPUCores *int
	RamGB    *float64 // nolint
	DiskGB   *float64
	GPUQuota *GPUQuota
}

// GPUQuota ...
type GPUQuota struct {
	GPU map[string]float64
}

// IsGlobal ...
func (q *Quota) IsGlobal() bool {
	return q.ID == consts.GlobalQuotaID
}

// IsDefault ...
func (q *Quota) IsDefault() bool {
	return q.AccountID == consts.DefaultQuotaAccountID && q.UserID == ""
}
