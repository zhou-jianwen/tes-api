package sql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/cluster/application/query"
	"code.byted.org/epscp/vetes-api/pkg/testutil"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

var clusterDTO = &query.Cluster{
	ID:                 id,
	HeartbeatTimestamp: now,
	Capacity: &query.Capacity{
		Count:    utils.Point(100),
		CPUCores: utils.Point(500),
		RamGB:    utils.Point[float64](2000),
		DiskGB:   utils.Point[float64](10000),
		GPUCapacity: &query.GPUCapacity{
			GPU: map[string]float64{
				"gpu-01": 8,
				"gpu-02": 10,
			},
		},
	},
	Limits: &query.Limits{
		CPUCores: utils.Point(50),
		RamGB:    utils.Point[float64](200),
		GPULimit: &query.GPULimit{
			GPU: map[string]float64{
				"gpu-01": 1,
				"gpu-02": 1,
			},
		},
	},
}

func TestList(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `cluster`").
		WillReturnRows(sqlmock.NewRows(rows).AddRow(clusterPO.ID, clusterPO.HeartbeatTimestamp,
			testutil.MustJSONMarshal(clusterPO.Capacity), testutil.MustJSONMarshal(clusterPO.Limits)))
	resp, err := r.List(context.TODO(), nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.Cluster{clusterDTO}))
}
