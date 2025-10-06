package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/quota/domain"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

func TestGet(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().GetOrDefault(gomock.Any(), false, "ac1", "").
		Return(&domain.Quota{
			AccountID: consts.DefaultQuotaAccountID,
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
		}, nil)

	handler := NewGetHandler(fakeService)
	resp, err := handler.Handle(context.TODO(), &GetQuery{
		Global:    false,
		AccountID: "ac1",
		UserID:    "",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Quota{
		Global:    false,
		Default:   true,
		AccountID: "ac1",
		UserID:    "",
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
	}))
}
