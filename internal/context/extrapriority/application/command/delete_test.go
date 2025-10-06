package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/extrapriority/domain"
)

func TestDelete(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Delete(gomock.Any(), "", "", "sb1", "").
		Return(nil)

	handler := NewDeleteHandler(fakeService)
	err := handler.Handle(context.TODO(), &DeleteCommand{SubmissionID: "sb1"})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
