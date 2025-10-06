package sql

import (
	"github.com/GBA-BI/tes-api/internal/context/quota/domain"
)

func (q *Quota) toDO() *domain.Quota {
	res := &domain.Quota{
		ID:        q.ID,
		AccountID: q.AccountID,
		UserID:    q.UserID,
	}
	if q.ResourceQuota != nil {
		res.ResourceQuota = &domain.ResourceQuota{
			Count:    q.ResourceQuota.Count,
			CPUCores: q.ResourceQuota.CPUCores,
			RamGB:    q.ResourceQuota.RamGB,
			DiskGB:   q.ResourceQuota.DiskGB,
		}
		if q.ResourceQuota.GPUQuota != nil {
			res.ResourceQuota.GPUQuota = &domain.GPUQuota{
				GPU: q.ResourceQuota.GPUQuota.GPU,
			}
		}
	}
	return res
}

func quotaDOToPO(quota *domain.Quota) *Quota {
	if quota == nil {
		return nil
	}
	res := &Quota{
		ID:        quota.ID,
		AccountID: quota.AccountID,
		UserID:    quota.UserID,
	}
	if quota.ResourceQuota != nil {
		res.ResourceQuota = &ResourceQuota{
			Count:    quota.ResourceQuota.Count,
			CPUCores: quota.ResourceQuota.CPUCores,
			RamGB:    quota.ResourceQuota.RamGB,
			DiskGB:   quota.ResourceQuota.DiskGB,
		}
		if quota.ResourceQuota.GPUQuota != nil {
			res.ResourceQuota.GPUQuota = &GPUQuota{
				GPU: quota.ResourceQuota.GPUQuota.GPU,
			}
		}
	}
	return res
}
