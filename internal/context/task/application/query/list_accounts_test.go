package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

func TestListAccounts(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().ListAccounts(gomock.Any()).
		Return([]*AccountInfo{{AccountID: "account-01", UserIDs: []string{"", "user-01"}}}, nil)

	handler := NewListAccountsHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO())
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.HaveLen(1))
}
