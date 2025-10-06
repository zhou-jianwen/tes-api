package sql

import (
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/application/query"
	"github.com/GBA-BI/tes-api/internal/context/extrapriority/domain"
)

func (e *ExtraPriority) toDTO() *query.ExtraPriority {
	return &query.ExtraPriority{
		AccountID:          e.AccountID,
		UserID:             e.UserID,
		SubmissionID:       e.SubmissionID,
		RunID:              e.RunID,
		ExtraPriorityValue: e.ExtraPriorityValue,
	}
}

func (e *ExtraPriority) toDO() *domain.ExtraPriority {
	return &domain.ExtraPriority{
		ID:                 e.ID,
		AccountID:          e.AccountID,
		UserID:             e.UserID,
		SubmissionID:       e.SubmissionID,
		RunID:              e.RunID,
		ExtraPriorityValue: e.ExtraPriorityValue,
	}
}

func extraPriorityDOToPO(e *domain.ExtraPriority) *ExtraPriority {
	return &ExtraPriority{
		ID:                 e.ID,
		AccountID:          e.AccountID,
		UserID:             e.UserID,
		SubmissionID:       e.SubmissionID,
		RunID:              e.RunID,
		ExtraPriorityValue: e.ExtraPriorityValue,
	}
}
