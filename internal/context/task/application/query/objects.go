package query

import "time"

// TaskMinimal ...
type TaskMinimal struct {
	ID    string
	State string
}

// TaskBasic ...
type TaskBasic struct {
	TaskMinimal
	Name          string
	Description   string
	Resources     *Resources
	Executors     []*Executor
	Volumes       []string
	Tags          map[string]string
	Logs          []*TaskLog
	CreationTime  time.Time
	BioosInfo     *BioosInfo
	PriorityValue int
	ClusterID     string
}

// Task ...
type Task struct {
	TaskBasic
	Inputs  []*Input
	Outputs []*Output
}

// Input ...
type Input struct {
	Name        string
	Description string
	Path        string
	Type        string
	URL         string
	Content     string
}

// Output ...
type Output struct {
	Name        string
	Description string
	Path        string
	Type        string
	URL         string
}

// Resources ...
type Resources struct {
	CPUCores   int
	RamGB      float64 // nolint
	DiskGB     float64
	BootDiskGB *int
	GPU        *GPUResource
}

// GPUResource ...
type GPUResource struct {
	Count float64
	Type  string
}

// Executor ...
type Executor struct {
	Image   string
	Command []string
	Workdir string
	Stdin   string
	Stdout  string
	Stderr  string
	Env     map[string]string
}

// TaskLog ...
type TaskLog struct {
	ClusterID  string
	Logs       [][]*ExecutorLog
	StartTime  *time.Time
	EndTime    *time.Time
	SystemLogs []string
}

// ExecutorLog ...
type ExecutorLog struct {
	ExecutorID string
	StartTime  *time.Time
	EndTime    *time.Time
}

// BioosInfo ...
type BioosInfo struct {
	AccountID    string
	UserID       string
	SubmissionID string
	RunID        string
	Meta         *BioosInfoMeta
}

// BioosInfoMeta ...
type BioosInfoMeta struct {
	AAIPassport     *string
	MountTOS        *bool
	BucketsAuthInfo *BucketsAuthInfo
}

// BucketsAuthInfo ...
type BucketsAuthInfo struct {
	ReadOnly  []string
	ReadWrite []string
	External  []*ExternalBucketAuthInfo
}

// ExternalBucketAuthInfo ...
type ExternalBucketAuthInfo struct {
	Bucket string
	AK     string
	SK     string
}

// TasksResources ...
type TasksResources struct {
	Count    int
	CPUCores int
	RamGB    float64 // nolint
	DiskGB   float64
	GPU      map[string]float64
}

// AccountInfo ...
type AccountInfo struct {
	AccountID string
	UserIDs   []string
}
