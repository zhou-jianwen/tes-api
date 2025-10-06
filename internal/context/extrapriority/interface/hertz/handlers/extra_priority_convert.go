package handlers

import (
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application/command"
	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application/query"
)

func (r *PutExtraPriorityRequest) toDTO() *command.PutCommand {
	return &command.PutCommand{
		AccountID:          r.AccountID,
		UserID:             r.UserID,
		SubmissionID:       r.SubmissionID,
		RunID:              r.RunID,
		ExtraPriorityValue: r.ExtraPriorityValue,
	}
}

func (r *ListExtraPriorityRequest) toDTO() *query.ListQuery {
	return &query.ListQuery{
		Filter: &query.ListFilter{
			AccountID:    r.AccountID,
			SubmissionID: r.SubmissionID,
			RunID:        r.RunID,
		},
	}
}

func (r *DeleteExtraPriorityRequest) toDTO() *command.DeleteCommand {
	return &command.DeleteCommand{
		AccountID:    r.AccountID,
		UserID:       r.UserID,
		SubmissionID: r.SubmissionID,
		RunID:        r.RunID,
	}
}

func extraPriorityDTOToVO(extraPriority *query.ExtraPriority) *ExtraPriority {
	return &ExtraPriority{
		AccountID:          extraPriority.AccountID,
		UserID:             extraPriority.UserID,
		SubmissionID:       extraPriority.SubmissionID,
		RunID:              extraPriority.RunID,
		ExtraPriorityValue: extraPriority.ExtraPriorityValue,
	}
}
