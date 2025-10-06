package sql

// ExtraPriority ...
type ExtraPriority struct {
	ID                 string `gorm:"column:id;type:VARCHAR(128);not null;primaryKey"`
	AccountID          string `gorm:"column:account_id;type:VARCHAR(32);not null;default:''"`
	UserID             string `gorm:"column:user_id;type:VARCHAR(32);not null;default:''"`
	SubmissionID       string `gorm:"column:submission_id;type:VARCHAR(32);not null;default:''"`
	RunID              string `gorm:"column:run_id;type:VARCHAR(32);not null;default:''"`
	ExtraPriorityValue int    `gorm:"column:extra_priority_value;type:BIGINT;not null;default:0"`
}

// TableName ...
func (e *ExtraPriority) TableName() string {
	return "extra_priority"
}
