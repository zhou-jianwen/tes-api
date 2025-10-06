package sql

import "time"

// Task ...
type Task struct {
	TaskBasic
	Inputs  []*Input  `gorm:"column:inputs;type:LONGTEXT;serializer:json"`
	Outputs []*Output `gorm:"column:outputs;type:LONGTEXT;serializer:json"`
}

// TaskBasic ...
type TaskBasic struct {
	TaskStatus
	Name          string            `gorm:"column:name;type:VARCHAR(512);not null;default:'';index:name"`
	Description   string            `gorm:"column:description;type:LONGTEXT"`
	Resources     *Resources        `gorm:"embedded"`
	Executors     []*Executor       `gorm:"column:executors;type:LONGTEXT;serializer:json"`
	Volumes       []string          `gorm:"column:volumes;type:LONGTEXT;serializer:json"`
	Tags          map[string]string `gorm:"column:tags;type:LONGTEXT;serializer:json"`
	BioosInfo     *BioosInfo        `gorm:"embedded"`
	PriorityValue int               `gorm:"column:priority_value;type:BIGINT;not null;default:0"`
}

// TaskStatus ...
type TaskStatus struct {
	TaskState
	Logs         []*TaskLog `gorm:"column:logs;type:LONGTEXT;serializer:json"`
	CreationTime time.Time  `gorm:"column:creation_time;type:DATETIME;not null"`
	// ClusterID may be updated to empty string, we have to mark it as pointer because
	// gorm do not update default value
	ClusterID *string `gorm:"column:cluster_id;type:VARCHAR(32);not null;default:'';index:state_cluster,priority:2"`

	StatusResourceVersion int `gorm:"column:status_resource_version;type:BIGINT;not null;default:0"`
}

// TaskState ...
type TaskState struct {
	ID    string `gorm:"column:id;type:VARCHAR(16);not null;primaryKey"`
	State string `gorm:"column:state;type:VARCHAR(16);not null;index:state_cluster,priority:1;index:state_account_user,priority:1"`
}

// Input ...
type Input struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	URL         string `json:"url,omitempty"`
	Content     string `json:"content,omitempty"`
}

// Output ...
type Output struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	URL         string `json:"url"`
}

// Resources ...
type Resources struct {
	CPUCores   int      `gorm:"column:cpu_cores;type:SMALLINT;not null"`
	RamGB      float64  `gorm:"column:ram_gb;type:DOUBLE;not null"` // nolint
	DiskGB     float64  `gorm:"column:disk_gb;type:DOUBLE;not null"`
	BootDiskGB *int     `gorm:"column:boot_disk_gb;type:SMALLINT"`
	GPUCount   *float64 `gorm:"column:gpu_count;type:DOUBLE"`
	GPUType    *string  `gorm:"column:gpu_type;type:VARCHAR(32)"`
}

// Executor ...
type Executor struct {
	Image   string            `json:"image"`
	Command []string          `json:"command"`
	Workdir string            `json:"workdir,omitempty"`
	Stdin   string            `json:"stdin,omitempty"`
	Stdout  string            `json:"stdout,omitempty"`
	Stderr  string            `json:"stderr,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// TaskLog ...
type TaskLog struct {
	ClusterID  string           `json:"cluster_id"`
	Logs       [][]*ExecutorLog `json:"logs,omitempty"`
	StartTime  *time.Time       `json:"start_time,omitempty"`
	EndTime    *time.Time       `json:"end_time,omitempty"`
	SystemLogs []string         `json:"system_logs,omitempty"`
}

// ExecutorLog ...
type ExecutorLog struct {
	ExecutorID string     `json:"executor_id"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
}

// BioosInfo ...
type BioosInfo struct {
	AccountID    string         `gorm:"column:account_id;type:VARCHAR(32);not null;default:'';index:state_account_user,priority:2"`
	UserID       string         `gorm:"column:user_id;type:VARCHAR(32);not null;default:'';index:state_account_user,priority:3"`
	SubmissionID string         `gorm:"column:submission_id;type:VARCHAR(32);not null;default:''"`
	RunID        string         `gorm:"column:run_id;type:VARCHAR(32);not null;default:''"`
	Meta         *BioosInfoMeta `gorm:"column:meta;type:longtext;serializer:json"`
}

// BioosInfoMeta ...
type BioosInfoMeta struct {
	AAIPassport     *string          `json:"aai_passport,omitempty"`
	MountTOS        *bool            `json:"mount_tos,omitempty"`
	BucketsAuthInfo *BucketsAuthInfo `json:"buckets_auth_info,omitempty"`
}

// BucketsAuthInfo ...
type BucketsAuthInfo struct {
	ReadOnly  []string                  `json:"read_only,omitempty"`
	ReadWrite []string                  `json:"read_write,omitempty"`
	External  []*ExternalBucketAuthInfo `json:"external,omitempty"`
}

// ExternalBucketAuthInfo ...
type ExternalBucketAuthInfo struct {
	Bucket string `json:"bucket,omitempty"`
	AK     string `json:"ak,omitempty"`
	SK     string `json:"sk,omitempty"`
}

// TableName ...
func (t *Task) TableName() string {
	return "task"
}
