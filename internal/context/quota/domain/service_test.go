package domain

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/pkg/consts"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

var resourceQuota = &ResourceQuota{
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
}

func TestGetAccountCustom(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/").
		Return(&Quota{
			ID:            "ac1/",
			AccountID:     "ac1",
			UserID:        "",
			ResourceQuota: resourceQuota,
		}, nil)

	svc := NewService(fakeRepo)
	resp, err := svc.GetOrDefault(context.TODO(), false, "ac1", "")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Quota{
		ID:            "ac1/",
		AccountID:     "ac1",
		UserID:        "",
		ResourceQuota: resourceQuota,
	}))
}

func TestGetAccountNotFoundDefault(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/").
		Return(nil, apperrors.NewNotFoundError("quota", "ac1/")) // custom not found
	fakeRepo.EXPECT().Get(gomock.Any(), consts.DefaultQuotaAccountID).
		Return(&Quota{
			ID:            consts.DefaultQuotaAccountID,
			AccountID:     consts.DefaultQuotaAccountID,
			UserID:        "",
			ResourceQuota: resourceQuota,
		}, nil)

	svc := NewService(fakeRepo)
	resp, err := svc.GetOrDefault(context.TODO(), false, "ac1", "")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Quota{
		ID:            consts.DefaultQuotaAccountID,
		AccountID:     consts.DefaultQuotaAccountID,
		UserID:        "",
		ResourceQuota: resourceQuota,
	}))
}

func TestGetAccountNOtFoundDefaultNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/").
		Return(nil, apperrors.NewNotFoundError("quota", "ac1/")) // custom not found
	fakeRepo.EXPECT().Get(gomock.Any(), consts.DefaultQuotaAccountID).
		Return(nil, apperrors.NewNotFoundError("quota", "0")) // default not found

	svc := NewService(fakeRepo)
	_, err := svc.GetOrDefault(context.TODO(), false, "ac1", "")
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestGetUser(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/u1").
		Return(&Quota{
			ID:            "ac1/",
			AccountID:     "ac1",
			UserID:        "u1",
			ResourceQuota: resourceQuota,
		}, nil)

	svc := NewService(fakeRepo)
	resp, err := svc.GetOrDefault(context.TODO(), false, "ac1", "u1")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Quota{
		ID:            "ac1/",
		AccountID:     "ac1",
		UserID:        "u1",
		ResourceQuota: resourceQuota,
	}))
}

func TestGetUserNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/u1").
		Return(nil, apperrors.NewNotFoundError("quota", "ac1/u1"))

	svc := NewService(fakeRepo)
	_, err := svc.GetOrDefault(context.TODO(), false, "ac1", "u1")
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestGetGlobal(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), consts.GlobalQuotaID).
		Return(&Quota{
			ID:            consts.GlobalQuotaID,
			AccountID:     "",
			UserID:        "",
			ResourceQuota: resourceQuota,
		}, nil)

	svc := NewService(fakeRepo)
	resp, err := svc.GetOrDefault(context.TODO(), true, "", "")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Quota{
		ID:            consts.GlobalQuotaID,
		AccountID:     "",
		UserID:        "",
		ResourceQuota: resourceQuota,
	}))
}

func TestGetGlobalNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), consts.GlobalQuotaID).
		Return(nil, apperrors.NewNotFoundError("quota", "global"))

	svc := NewService(fakeRepo)
	_, err := svc.GetOrDefault(context.TODO(), true, "", "")
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestGetDefault(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), consts.DefaultQuotaAccountID).
		Return(&Quota{
			ID:            consts.DefaultQuotaAccountID,
			AccountID:     consts.DefaultQuotaAccountID,
			UserID:        "",
			ResourceQuota: resourceQuota,
		}, nil)

	svc := NewService(fakeRepo)
	resp, err := svc.GetOrDefault(context.TODO(), false, consts.DefaultQuotaAccountID, "")
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Quota{
		ID:            consts.DefaultQuotaAccountID,
		AccountID:     consts.DefaultQuotaAccountID,
		UserID:        "",
		ResourceQuota: resourceQuota,
	}))
}

func TestGetDefaultNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), consts.DefaultQuotaAccountID).
		Return(nil, apperrors.NewNotFoundError("quota", "0"))

	svc := NewService(fakeRepo)
	_, err := svc.GetOrDefault(context.TODO(), false, consts.DefaultQuotaAccountID, "")
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestPut(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Save(gomock.Any(), &Quota{
		ID:            "ac1/",
		AccountID:     "ac1",
		UserID:        "",
		ResourceQuota: resourceQuota,
	}).Return(nil)

	svc := NewService(fakeRepo)
	err := svc.Put(context.TODO(), false, "ac1", "", resourceQuota)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/u1").
		Return(&Quota{
			ID:            "ac1/u1",
			AccountID:     "ac1",
			UserID:        "u1",
			ResourceQuota: resourceQuota,
		}, nil)
	fakeRepo.EXPECT().Delete(gomock.Any(), "ac1/u1").
		Return(nil)

	svc := NewService(fakeRepo)
	err := svc.Delete(context.TODO(), false, "ac1", "u1")
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDeleteNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1/u1").
		Return(nil, apperrors.NewNotFoundError("quota", "ac1/u1"))

	svc := NewService(fakeRepo)
	err := svc.Delete(context.TODO(), false, "ac1", "u1")
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}
