package sql

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/domain"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
	"code.byted.org/epscp/vetes-api/pkg/testutil"
)

var id = "ac1/u1/sb1/r1"

var priorityPO = &ExtraPriority{
	ID:                 id,
	AccountID:          "ac1",
	UserID:             "r1",
	SubmissionID:       "sb1",
	RunID:              "r1",
	ExtraPriorityValue: 100,
}

var priorityDO = &domain.ExtraPriority{
	ID:                 id,
	AccountID:          "ac1",
	UserID:             "r1",
	SubmissionID:       "sb1",
	RunID:              "r1",
	ExtraPriorityValue: 100,
}

var rows = []string{"id", "account_id", "user_id", "submission_id", "run_id", "extra_priority_value"}

func TestGet(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `extra_priority` WHERE `id` = ? ORDER BY `extra_priority`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(rows).AddRow(priorityPO.ID, priorityPO.AccountID, priorityPO.UserID,
			priorityPO.SubmissionID, priorityPO.RunID, priorityPO.ExtraPriorityValue))
	resp, err := r.Get(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(priorityDO))
}

func TestGetNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `extra_priority` WHERE `id` = ? ORDER BY `extra_priority`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(rows))
	_, err := r.Get(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestSave(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("INSERT INTO `extra_priority` %s %s", testutil.GenInsertSql(rows), testutil.GenDuplicateKeySql(rows[1:] /*without id*/))).
		WithArgs(priorityPO.ID, priorityPO.AccountID, priorityPO.UserID, priorityPO.SubmissionID, priorityPO.RunID, priorityPO.ExtraPriorityValue).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Save(context.TODO(), priorityDO)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDelete(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `extra_priority` WHERE `id` = ?").WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Delete(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
