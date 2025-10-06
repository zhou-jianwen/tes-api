package sql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/cluster/domain"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/testutil"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

var id = "cluster-x"
var now = time.Now().UTC().Truncate(time.Second)

var clusterPO = &Cluster{
	ID:                 id,
	HeartbeatTimestamp: now,
	Capacity: &Capacity{
		Count:    utils.Point(100),
		CPUCores: utils.Point(500),
		RamGB:    utils.Point[float64](2000),
		DiskGB:   utils.Point[float64](10000),
		GPUCapacity: &GPUCapacity{
			GPU: map[string]float64{
				"gpu-01": 8,
				"gpu-02": 10,
			},
		},
	},
	Limits: &Limits{
		CPUCores: utils.Point(50),
		RamGB:    utils.Point[float64](200),
		GPULimit: &GPULimit{
			GPU: map[string]float64{
				"gpu-01": 1,
				"gpu-02": 1,
			},
		},
	},
}

var clusterDO = &domain.Cluster{
	ID:                 id,
	HeartbeatTimestamp: now,
	Capacity: &domain.Capacity{
		Count:    utils.Point(100),
		CPUCores: utils.Point(500),
		RamGB:    utils.Point[float64](2000),
		DiskGB:   utils.Point[float64](10000),
		GPUCapacity: &domain.GPUCapacity{
			GPU: map[string]float64{
				"gpu-01": 8,
				"gpu-02": 10,
			},
		},
	},
	Limits: &domain.Limits{
		CPUCores: utils.Point(50),
		RamGB:    utils.Point[float64](200),
		GPULimit: &domain.GPULimit{
			GPU: map[string]float64{
				"gpu-01": 1,
				"gpu-02": 1,
			},
		},
	},
}

var rows = []string{"id", "heartbeat_timestamp", "capacity", "limits"}

func TestGet(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `cluster` WHERE `id` = ? ORDER BY `cluster`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(rows).AddRow(clusterPO.ID, clusterPO.HeartbeatTimestamp,
			testutil.MustJSONMarshal(clusterPO.Capacity), testutil.MustJSONMarshal(clusterPO.Limits)))
	resp, err := r.Get(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(clusterDO))
}

func TestGetNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `cluster` WHERE `id` = ? ORDER BY `cluster`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(rows))
	_, err := r.Get(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestSave(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("INSERT INTO `cluster` %s %s", testutil.GenInsertSql(rows), testutil.GenDuplicateKeySql(rows[1:] /*without id*/))).
		WithArgs(clusterPO.ID, clusterPO.HeartbeatTimestamp,
			testutil.MustJSONMarshal(clusterPO.Capacity), testutil.MustJSONMarshal(clusterPO.Limits)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Save(context.TODO(), clusterDO)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDelete(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `cluster` WHERE `id` = ?").WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Delete(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
