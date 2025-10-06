package handlers

// PutExtraPriorityRequest ...
type PutExtraPriorityRequest struct {
	AccountID          string `query:"account_id" json:"-"`
	UserID             string `query:"user_id" json:"-"`
	SubmissionID       string `query:"submission_id" json:"-"`
	RunID              string `query:"run_id" json:"-"`
	ExtraPriorityValue int    `json:"extra_priority_value"`
}

// PutExtraPriorityResponse ...
type PutExtraPriorityResponse struct{}

// ListExtraPriorityRequest ...
type ListExtraPriorityRequest struct {
	AccountID    string `query:"account_id"`
	SubmissionID string `query:"submission_id"`
	RunID        string `query:"run_id"`
}

// ListExtraPriorityResponse ...
type ListExtraPriorityResponse []*ExtraPriority

// DeleteExtraPriorityRequest ...
type DeleteExtraPriorityRequest struct {
	AccountID    string `query:"account_id"`
	UserID       string `query:"user_id"`
	SubmissionID string `query:"submission_id"`
	RunID        string `query:"run_id"`
}

// DeleteExtraPriorityResponse ...
type DeleteExtraPriorityResponse struct{}

// ExtraPriority ...
type ExtraPriority struct {
	AccountID          string `json:"account_id,omitempty"`
	UserID             string `json:"user_id,omitempty"`
	SubmissionID       string `json:"submission_id,omitempty"`
	RunID              string `json:"run_id,omitempty"`
	ExtraPriorityValue int    `json:"extra_priority_value"`
}
