package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/cluster/domain"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

func TestPut(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Put(gomock.Any(), gomock.Any()).
		Return(nil)

	handler := NewPutHandler(fakeService)
	err := handler.Handle(context.TODO(), &PutCommand{
		ID: "cluster-01",
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
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
