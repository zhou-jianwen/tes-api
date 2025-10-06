package domain

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

var cluster = &Cluster{
	ID:                 "cluster-x",
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
}

func TestPut(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Save(gomock.Any(), gomock.Any()).
		Return(nil)

	svc := NewService(fakeRepo)
	err := svc.Put(context.TODO(), cluster)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), cluster.ID).
		Return(cluster, nil)
	fakeRepo.EXPECT().Delete(gomock.Any(), cluster.ID).
		Return(nil)

	svc := NewService(fakeRepo)
	err := svc.Delete(context.TODO(), cluster.ID)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDeleteNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), cluster.ID).
		Return(nil, apperrors.NewNotFoundError("cluster", cluster.ID))

	svc := NewService(fakeRepo)
	err := svc.Delete(context.TODO(), cluster.ID)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}
