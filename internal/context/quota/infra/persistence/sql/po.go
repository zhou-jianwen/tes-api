package sql

// Quota ...
type Quota struct {
	ID            string         `gorm:"column:id;type:VARCHAR(128);not null;primaryKey"`
	AccountID     string         `gorm:"column:account_id;type:VARCHAR(32);not null;default:''"`
	UserID        string         `gorm:"column:user_id;type:VARCHAR(32);not null;default:''"`
	ResourceQuota *ResourceQuota `gorm:"column:resource_quota;type:LONGTEXT;serializer:json"`
}

// ResourceQuota ...
type ResourceQuota struct {
	Count    *int      `json:"count,omitempty"`
	CPUCores *int      `json:"cpu_cores,omitempty"`
	RamGB    *float64  `json:"ram_gb,omitempty"` // nolint
	DiskGB   *float64  `json:"disk_gb,omitempty"`
	GPUQuota *GPUQuota `json:"gpu_quota,omitempty"`
}

// GPUQuota ...
type GPUQuota struct {
	GPU map[string]float64 `json:"gpu,omitempty"`
}

// TableName ...
func (q *Quota) TableName() string {
	return "quota"
}
