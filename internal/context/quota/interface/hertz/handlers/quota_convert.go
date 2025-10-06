package handlers

import (
	"code.byted.org/epscp/vetes-api/internal/context/quota/application/command"
	"code.byted.org/epscp/vetes-api/internal/context/quota/application/query"
)

func (r *GetQuotaRequest) toDTO() *query.GetQuery {
	return &query.GetQuery{
		Global:    r.Global,
		AccountID: r.AccountID,
		UserID:    r.UserID,
	}
}

func (r *PutQuotaRequest) toDTO() *command.PutCommand {
	res := &command.PutCommand{
		Global:    r.Global,
		AccountID: r.AccountID,
		UserID:    r.UserID,
	}
	if r.ResourceQuota != nil {
		res.ResourceQuota = &command.ResourceQuota{
			Count:    r.ResourceQuota.Count,
			CPUCores: r.ResourceQuota.CPUCores,
			RamGB:    r.ResourceQuota.RamGB,
			DiskGB:   r.ResourceQuota.DiskGB,
		}
		if r.ResourceQuota.GPUQuota != nil {
			res.ResourceQuota.GPUQuota = &command.GPUQuota{
				GPU: r.ResourceQuota.GPUQuota.GPU,
			}
		}
	}
	return res
}

func (r *DeleteQuotaRequest) toDTO() *command.DeleteCommand {
	return &command.DeleteCommand{
		Global:    r.Global,
		AccountID: r.AccountID,
		UserID:    r.UserID,
	}
}

func quotaDTOToVO(quota *query.Quota) *GetQuotaResponse {
	res := &GetQuotaResponse{
		Global:    quota.Global,
		Default:   quota.Default,
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
