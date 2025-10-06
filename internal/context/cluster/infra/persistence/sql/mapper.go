package sql

import (
	"code.byted.org/epscp/vetes-api/internal/context/cluster/application/query"
	"code.byted.org/epscp/vetes-api/internal/context/cluster/domain"
)

func (c *Cluster) toDTO() *query.Cluster {
	if c == nil {
		return nil
	}
	res := &query.Cluster{
		ID:                 c.ID,
		HeartbeatTimestamp: c.HeartbeatTimestamp,
	}
	if c.Capacity != nil {
		res.Capacity = &query.Capacity{
			Count:    c.Capacity.Count,
			CPUCores: c.Capacity.CPUCores,
			RamGB:    c.Capacity.RamGB,
			DiskGB:   c.Capacity.DiskGB,
		}
		if c.Capacity.GPUCapacity != nil {
			res.Capacity.GPUCapacity = &query.GPUCapacity{
				GPU: c.Capacity.GPUCapacity.GPU,
			}
		}
	}
	if c.Limits != nil {
		res.Limits = &query.Limits{
			CPUCores: c.Limits.CPUCores,
			RamGB:    c.Limits.RamGB,
		}
		if c.Limits.GPULimit != nil {
			res.Limits.GPULimit = &query.GPULimit{
				GPU: c.Limits.GPULimit.GPU,
			}
		}
	}
	return res
}

func (c *Cluster) toDO() *domain.Cluster {
	if c == nil {
		return nil
	}
	res := &domain.Cluster{
		ID:                 c.ID,
		HeartbeatTimestamp: c.HeartbeatTimestamp,
	}
	if c.Capacity != nil {
		res.Capacity = &domain.Capacity{
			Count:    c.Capacity.Count,
			CPUCores: c.Capacity.CPUCores,
			RamGB:    c.Capacity.RamGB,
			DiskGB:   c.Capacity.DiskGB,
		}
		if c.Capacity.GPUCapacity != nil {
			res.Capacity.GPUCapacity = &domain.GPUCapacity{
				GPU: c.Capacity.GPUCapacity.GPU,
			}
		}
	}
	if c.Limits != nil {
		res.Limits = &domain.Limits{
			CPUCores: c.Limits.CPUCores,
			RamGB:    c.Limits.RamGB,
		}
		if c.Limits.GPULimit != nil {
			res.Limits.GPULimit = &domain.GPULimit{
				GPU: c.Limits.GPULimit.GPU,
			}
		}
	}
	return res
}

func clusterDOToPO(cluster *domain.Cluster) *Cluster {
	if cluster == nil {
		return nil
	}
	res := &Cluster{
		ID:                 cluster.ID,
		HeartbeatTimestamp: cluster.HeartbeatTimestamp,
	}
	if cluster.Capacity != nil {
		res.Capacity = &Capacity{
			Count:    cluster.Capacity.Count,
			CPUCores: cluster.Capacity.CPUCores,
			RamGB:    cluster.Capacity.RamGB,
			DiskGB:   cluster.Capacity.DiskGB,
		}
		if cluster.Capacity.GPUCapacity != nil {
			res.Capacity.GPUCapacity = &GPUCapacity{
				GPU: cluster.Capacity.GPUCapacity.GPU,
			}
		}
	}
	if cluster.Limits != nil {
		res.Limits = &Limits{
			CPUCores: cluster.Limits.CPUCores,
			RamGB:    cluster.Limits.RamGB,
		}
		if cluster.Limits.GPULimit != nil {
			res.Limits.GPULimit = &GPULimit{
				GPU: cluster.Limits.GPULimit.GPU,
			}
		}
	}
	return res
}
