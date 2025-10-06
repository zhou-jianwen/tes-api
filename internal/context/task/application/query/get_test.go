package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/pkg/consts"
)

func TestGetMinimal(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().GetMinimal(gomock.Any(), "task-1234").
		Return(&TaskMinimal{ID: "task-1234", State: consts.TaskComplete}, nil)

	handler := NewGetHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO(), &GetQuery{
		ID:   "task-1234",
		View: "", // default minimal
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Task{
		TaskBasic: TaskBasic{
			TaskMinimal: TaskMinimal{
				ID:    "task-1234",
				State: consts.TaskComplete,
			},
		},
	}))
}

func TestGetBasic(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().GetBasic(gomock.Any(), "task-1234").
		Return(&TaskBasic{
			TaskMinimal: TaskMinimal{ID: "task-1234", State: consts.TaskComplete},
			Logs:        []*TaskLog{{SystemLogs: []string{"abcd"}}},
		}, nil)

	handler := NewGetHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO(), &GetQuery{
		ID:   "task-1234",
		View: consts.BasicView,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Task{
		TaskBasic: TaskBasic{
			TaskMinimal: TaskMinimal{
				ID:    "task-1234",
				State: consts.TaskComplete,
			},
			Logs: []*TaskLog{{}}, // no systemLogs
		},
	}))
}

func TestGetFull(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().GetFull(gomock.Any(), "task-1234").
		Return(&Task{
			TaskBasic: TaskBasic{
				TaskMinimal: TaskMinimal{ID: "task-1234", State: consts.TaskComplete},
				Logs:        []*TaskLog{{SystemLogs: []string{"abcd"}}},
			},
		}, nil)

	handler := NewGetHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO(), &GetQuery{
		ID:   "task-1234",
		View: consts.FullView,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&Task{
		TaskBasic: TaskBasic{
			TaskMinimal: TaskMinimal{
				ID:    "task-1234",
				State: consts.TaskComplete,
			},
			Logs: []*TaskLog{{SystemLogs: []string{"abcd"}}},
		},
	}))
}
