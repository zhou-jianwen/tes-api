package consts

const (
	// XRequestIDKey is request id key in log
	XRequestIDKey = "X-Request-ID"
)

// api prefix
const (
	Ga4ghAPIPrefix = "/api/ga4gh/tes/v1"
	OtherAPIPrefix = "/api/v1"
)

// MySQLType is mysql db type
const MySQLType = "mysql"

// task view types
const (
	MinimalView = "MINIMAL"
	BasicView   = "BASIC"
	FullView    = "FULL"
)

// task state
const (
	TaskQueued        = "QUEUED"
	TaskInitializing  = "INITIALIZING"
	TaskRunning       = "RUNNING"
	TaskComplete      = "COMPLETE"
	TaskSystemError   = "SYSTEM_ERROR"
	TaskExecutorError = "EXECUTOR_ERROR"
	TaskCanceling     = "CANCELING"
	TaskCanceled      = "CANCELED"
)

// GlobalQuotaID is id of global quota
const GlobalQuotaID = "global"

// DefaultQuotaAccountID is AccountID of default quota
const DefaultQuotaAccountID = "0"
