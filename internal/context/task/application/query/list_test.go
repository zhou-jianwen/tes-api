package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/pkg/consts"
)

func TestListMinimal(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().ListMinimal(gomock.Any(), defaultPageSize, nil, &ListFilter{WithoutCluster: true}).
		Return([]*TaskMinimal{{ID: "task-1111", State: consts.TaskQueued}}, nil, nil)

	handler := NewListHandler(fakeReadModel)
	resp, nextPageToken, err := handler.Handle(context.TODO(), &ListQuery{
		View:      "", // default minimal
		PageSize:  0,  // default 256
		PageToken: nil,
		Filter:    &ListFilter{WithoutCluster: true},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.HaveLen(1))
	g.Expect(nextPageToken).To(gomega.BeNil())
}

func TestListBasic(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().ListBasic(gomock.Any(), 1024, nil, nil).
		Return([]*TaskBasic{{TaskMinimal: TaskMinimal{ID: "task-1111", State: consts.TaskQueued}}}, nil, nil)

	handler := NewListHandler(fakeReadModel)
	resp, nextPageToken, err := handler.Handle(context.TODO(), &ListQuery{
		View:     consts.BasicView,
		PageSize: 1024,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.HaveLen(1))
	g.Expect(nextPageToken).To(gomega.BeNil())
}

func TestListFull(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().ListFull(gomock.Any(), defaultPageSize, nil, nil).
		Return([]*Task{{TaskBasic: TaskBasic{TaskMinimal: TaskMinimal{ID: "task-1111", State: consts.TaskQueued}}}}, nil, nil)

	handler := NewListHandler(fakeReadModel)
	resp, nextPageToken, err := handler.Handle(context.TODO(), &ListQuery{
		View:     consts.FullView,
		PageSize: 0, // default 256
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.HaveLen(1))
	g.Expect(nextPageToken).To(gomega.BeNil())
}
