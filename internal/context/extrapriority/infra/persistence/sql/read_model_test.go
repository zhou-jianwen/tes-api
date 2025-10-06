package sql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/application/query"
	"code.byted.org/epscp/vetes-api/pkg/testutil"
)

var priorityDTO = &query.ExtraPriority{
	AccountID:          "ac1",
	UserID:             "r1",
	SubmissionID:       "sb1",
	RunID:              "r1",
	ExtraPriorityValue: 100,
}

func TestList(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `extra_priority` WHERE `account_id` = ? AND `submission_id` = ? AND `run_id` = ?").
		WithArgs("ac1", "sb1", "r1").
		WillReturnRows(sqlmock.NewRows(rows).AddRow(priorityPO.ID, priorityPO.AccountID, priorityPO.UserID,
			priorityPO.SubmissionID, priorityPO.RunID, priorityPO.ExtraPriorityValue))
	resp, err := r.List(context.TODO(), &query.ListFilter{
		AccountID:    "ac1",
		SubmissionID: "sb1",
		RunID:        "r1",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.ExtraPriority{priorityDTO}))
}
