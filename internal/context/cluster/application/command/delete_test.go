package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/cluster/domain"
)

func TestDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Delete(gomock.Any(), "cluster-id").
		Return(nil)

	handler := NewDeleteHandler(fakeService)
	err := handler.Handle(context.TODO(), &DeleteCommand{ID: "cluster-id"})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
