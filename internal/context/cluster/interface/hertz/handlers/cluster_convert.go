package handlers

import (
	"time"

	"github.com/GBA-BI/tes-api/internal/context/cluster/application/command"
	"github.com/GBA-BI/tes-api/internal/context/cluster/application/query"
)

func (r *PutClusterRequest) toDTO() *command.PutCommand {
	return &command.PutCommand{
		ID:       r.ID,
		Capacity: r.Capacity.toDTO(),
		Limits:   r.Limits.toDTO(),
	}
}

func (c *Capacity) toDTO() *command.Capacity {
	if c == nil {
		return nil
	}
	return &command.Capacity{
		Count:       c.Count,
		CPUCores:    c.CPUCores,
		RamGB:       c.RamGB,
		DiskGB:      c.DiskGB,
		GPUCapacity: c.GPUCapacity.toDTO(),
	}
}

func (c *GPUCapacity) toDTO() *command.GPUCapacity {
	if c == nil {
		return nil
	}
	return &command.GPUCapacity{GPU: c.GPU}
}

func (l *Limits) toDTO() *command.Limits {
	if l == nil {
		return nil
	}
	return &command.Limits{
		CPUCores: l.CPUCores,
		RamGB:    l.RamGB,
		GPULimit: l.GPULimit.toDTO(),
	}
}

func (l *GPULimit) toDTO() *command.GPULimit {
	if l == nil {
		return nil
	}
	return &command.GPULimit{GPU: l.GPU}
}

func (r *ListClustersRequest) toDTO() *query.ListQuery {
	return &query.ListQuery{Filter: &query.ListFilter{}}
}

func (r *DeleteClusterRequest) toDTO() *command.DeleteCommand {
	return &command.DeleteCommand{
		ID: r.ID,
	}
}

func clusterDTOToVO(cluster *query.Cluster) *Cluster {
	if cluster == nil {
		return nil
	}
	res := &Cluster{
		ID:       cluster.ID,
		Capacity: nil,
		Limits:   nil,
	}
	if !cluster.HeartbeatTimestamp.IsZero() {
		res.HeartbeatTimestamp = cluster.HeartbeatTimestamp.Format(time.RFC3339)
	}
	if cluster.Capacity != nil {
		res.Capacity = &Capacity{
			Count:    cluster.Capacity.Count,
			CPUCores: cluster.Capacity.CPUCores,
			RamGB:    cluster.Capacity.RamGB,
			DiskGB:   cluster.Capacity.DiskGB,
		}
		if cluster.Capacity.GPUCapacity != nil {
			res.Capacity.GPUCapacity = &GPUCapacity{GPU: cluster.Capacity.GPUCapacity.GPU}
		}
	}
	if cluster.Limits != nil {
		res.Limits = &Limits{
			CPUCores: cluster.Limits.CPUCores,
			RamGB:    cluster.Limits.RamGB,
		}
		if cluster.Limits.GPULimit != nil {
			res.Limits.GPULimit = &GPULimit{GPU: cluster.Limits.GPULimit.GPU}
		}
	}
	return res
}
