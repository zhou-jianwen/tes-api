package sql

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/quota/domain"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
	"code.byted.org/epscp/vetes-api/pkg/testutil"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

var id = "ac1/u1"

var quotaPO = &Quota{
	ID:        id,
	AccountID: "ac1",
	UserID:    "u1",
	ResourceQuota: &ResourceQuota{
		Count:    utils.Point(10),
		CPUCores: utils.Point(100),
		RamGB:    utils.Point[float64](200),
		DiskGB:   utils.Point[float64](1000),
		GPUQuota: &GPUQuota{
			GPU: map[string]float64{
				"gpu-01": 10,
				"gpu-02": 15,
			},
		},
	},
}

var quotaDO = &domain.Quota{
	ID:        id,
	AccountID: "ac1",
	UserID:    "u1",
	ResourceQuota: &domain.ResourceQuota{
		Count:    utils.Point(10),
		CPUCores: utils.Point(100),
		RamGB:    utils.Point[float64](200),
		DiskGB:   utils.Point[float64](1000),
		GPUQuota: &domain.GPUQuota{
			GPU: map[string]float64{
				"gpu-01": 10,
				"gpu-02": 15,
			},
		},
	},
}

var rows = []string{"id", "account_id", "user_id", "resource_quota"}

func TestGet(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `quota` WHERE `id` = ? ORDER BY `quota`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(rows).AddRow(quotaPO.ID, quotaPO.AccountID, quotaPO.UserID, testutil.MustJSONMarshal(quotaPO.ResourceQuota)))
	resp, err := r.Get(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(quotaDO))
}

func TestGetNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `quota` WHERE `id` = ? ORDER BY `quota`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(rows))
	_, err := r.Get(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestSave(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("INSERT INTO `quota` %s %s", testutil.GenInsertSql(rows), testutil.GenDuplicateKeySql(rows[1:] /*without id*/))).
		WithArgs(quotaPO.ID, quotaPO.AccountID, quotaPO.UserID, testutil.MustJSONMarshal(quotaPO.ResourceQuota)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Save(context.TODO(), quotaDO)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDelete(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `quota` WHERE `id` = ?").WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Delete(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
