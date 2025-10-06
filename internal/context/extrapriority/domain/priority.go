package domain

import (
	"fmt"

	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// ExtraPriority ...
type ExtraPriority struct {
	ID                 string
	AccountID          string
	UserID             string
	SubmissionID       string
	RunID              string
	ExtraPriorityValue int
}

// NewExtraPriorityID ...
func NewExtraPriorityID(accountID, userID, submissionID, runID string) (string, error) {
	nonEmptyCnt := 0
	if accountID != "" {
		nonEmptyCnt++
	}
	if submissionID != "" {
		nonEmptyCnt++
	}
	if runID != "" {
		nonEmptyCnt++
	}
	if nonEmptyCnt != 1 {
		return "", apperrors.NewInvalidError("only one field must be non-empty in account_id|submission_id|run_id")
	}
	if accountID == "" && userID != "" {
		return "", apperrors.NewInvalidError("empty account_id with non-empty user_id")
	}
	return fmt.Sprintf("%s/%s/%s/%s", accountID, userID, submissionID, runID), nil
}
