package query

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/pkg/utils"
)

func TestList(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().List(gomock.Any(), &ListFilter{}).
		Return([]*Cluster{{
			ID:                 "cluster-id",
			HeartbeatTimestamp: time.Now().UTC().Truncate(time.Second),
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
		}}, nil)

	handler := NewListHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO(), &ListQuery{Filter: &ListFilter{}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.HaveLen(1))
}
