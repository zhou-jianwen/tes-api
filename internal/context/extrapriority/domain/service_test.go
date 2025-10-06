package domain

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

func TestPut(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Save(gomock.Any(), &ExtraPriority{
		ID:                 "ac1/u1//",
		AccountID:          "ac1",
		UserID:             "u1",
		SubmissionID:       "",
		RunID:              "",
		ExtraPriorityValue: 20,
	}).Return(nil)

	svc := NewService(fakeRepo)
	err := svc.Put(context.TODO(), "ac1", "u1", "", "", 20)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "ac1///").
		Return(&ExtraPriority{
			ID:                 "ac1///",
			AccountID:          "ac1",
			ExtraPriorityValue: 100,
		}, nil)
	fakeRepo.EXPECT().Delete(gomock.Any(), "ac1///").
		Return(nil)

	svc := NewService(fakeRepo)
	err := svc.Delete(context.TODO(), "ac1", "", "", "")
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestDeleteNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().Get(gomock.Any(), "///r1").
		Return(nil, apperrors.NewNotFoundError("extra priority", "///r1"))

	svc := NewService(fakeRepo)
	err := svc.Delete(context.TODO(), "", "", "", "r1")
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}
