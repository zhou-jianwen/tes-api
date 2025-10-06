package handlers

// PutClusterRequest ...
type PutClusterRequest struct {
	ID       string    `path:"id" json:"-"`
	Capacity *Capacity `json:"capacity,omitempty"`
	Limits   *Limits   `json:"limits,omitempty"`
}

// PutClusterResponse ...
type PutClusterResponse struct{}

// ListClustersRequest ...
type ListClustersRequest struct {
}

// ListClustersResponse ...
type ListClustersResponse []*Cluster

// DeleteClusterRequest ...
type DeleteClusterRequest struct {
	ID string `path:"id"`
}

// DeleteClusterResponse ...
type DeleteClusterResponse struct{}

// Cluster ...
type Cluster struct {
	ID                 string    `json:"id"`
	HeartbeatTimestamp string    `json:"heartbeat_timestamp"`
	Capacity           *Capacity `json:"capacity,omitempty"`
	Limits             *Limits   `json:"limits,omitempty"`
}

// Capacity ...
type Capacity struct {
	Count       *int         `json:"count,omitempty"`
	CPUCores    *int         `json:"cpu_cores,omitempty"`
	RamGB       *float64     `json:"ram_gb,omitempty"` // nolint
	DiskGB      *float64     `json:"disk_gb,omitempty"`
	GPUCapacity *GPUCapacity `json:"gpu_capacity,omitempty"`
}

// GPUCapacity ...
type GPUCapacity struct {
	GPU map[string]float64 `json:"gpu,omitempty"`
}

// Limits ...
type Limits struct {
	CPUCores *int      `json:"cpu_cores,omitempty"`
	RamGB    *float64  `json:"ram_gb,omitempty"` // nolint
	GPULimit *GPULimit `json:"gpu_limit,omitempty"`
}

// GPULimit ...
type GPULimit struct {
	GPU map[string]float64 `json:"gpu,omitempty"`
}
