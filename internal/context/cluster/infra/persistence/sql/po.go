package sql

import "time"

// Cluster ...
type Cluster struct {
	ID                 string    `gorm:"column:id;type:VARCHAR(32);not null;primaryKey"`
	HeartbeatTimestamp time.Time `gorm:"column:heartbeat_timestamp;type:DATETIME;not null"`
	Capacity           *Capacity `gorm:"column:capacity;type:LONGTEXT;serializer:json"`
	Limits             *Limits   `gorm:"column:limits;type:LONGTEXT;serializer:json"`
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

// TableName ...
func (c *Cluster) TableName() string {
	return "cluster"
}
