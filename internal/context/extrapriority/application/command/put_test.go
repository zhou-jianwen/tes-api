package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/extrapriority/domain"
)

func TestPut(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Put(gomock.Any(), "ac1", "", "", "", 10).
		Return(nil)

	handler := NewPutHandler(fakeService)
	err := handler.Handle(context.TODO(), &PutCommand{AccountID: "ac1", ExtraPriorityValue: 10})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
