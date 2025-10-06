package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
)

func TestCancel(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Cancel(gomock.Any(), "task-1111").
		Return(nil)

	handler := NewCancelHandler(fakeService)
	err := handler.Handle(context.TODO(), &CancelCommand{ID: "task-1111"})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
