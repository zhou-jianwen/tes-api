package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/quota/domain"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

func TestPut(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Put(gomock.Any(), true, "", "", gomock.Any()).
		Return(nil)

	handler := NewPutHandler(fakeService)
	err := handler.Handle(context.TODO(), &PutCommand{
		Global: true,
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
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
