package query

import "github.com/GBA-BI/tes-api/internal/context/quota/domain"

// Quota ...
type Quota struct {
	Global        bool
	Default       bool
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

func resourceQuotaDOToDTO(r *domain.ResourceQuota) *ResourceQuota {
	if r == nil {
		return nil
	}
	res := &ResourceQuota{
		Count:    r.Count,
		CPUCores: r.CPUCores,
		RamGB:    r.RamGB,
		DiskGB:   r.DiskGB,
	}
	if r.GPUQuota != nil {
		res.GPUQuota = &GPUQuota{
			GPU: r.GPUQuota.GPU,
		}
	}
	return res
}
