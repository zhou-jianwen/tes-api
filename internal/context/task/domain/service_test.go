package domain

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/pkg/consts"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

var id = "task-1111"
var now = time.Now().UTC().Truncate(time.Second)

func TestCreate(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeNormalizer := NewFakeNormalizer(ctrl)
	fakeNormalizer.EXPECT().Normalize(gomock.Any()).
		Return(nil)
	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().CheckIDExist(gomock.Any(), gomock.Any()).
		Return(false, nil)
	fakeRepo.EXPECT().Create(gomock.Any(), gomock.Any()).
		Return(nil)

	svc := NewService(fakeRepo, fakeNormalizer)
	_, err := svc.Create(context.TODO(), &Task{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestCancel(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().GetStatus(gomock.Any(), id).
		Return(&TaskStatus{
			ID:           id,
			State:        consts.TaskRunning,
			ClusterID:    "cluster-01",
			CreationTime: now,
		}, nil)
	fakeRepo.EXPECT().UpdateStatus(gomock.Any(), &TaskStatus{
		ID:           id,
		State:        consts.TaskCanceling,
		ClusterID:    "cluster-01",
		CreationTime: now,
	}).Return(true, nil)

	svc := NewService(fakeRepo, nil)
	err := svc.Cancel(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeRepo := NewFakeRepo(ctrl)
	fakeRepo.EXPECT().GetStatus(gomock.Any(), id).
		Return(&TaskStatus{
			ID:           id,
			State:        consts.TaskQueued,
			ClusterID:    "",
			CreationTime: now,
		}, nil)
	fakeRepo.EXPECT().UpdateStatus(gomock.Any(), &TaskStatus{
		ID:           id,
		State:        consts.TaskQueued,
		ClusterID:    "cluster-01",
		CreationTime: now,
		Logs: []*TaskLog{{
			StartTime: &now,
		}},
	}).Return(true, nil)

	svc := NewService(fakeRepo, nil)
	err := svc.Update(context.TODO(), id, utils.Point(consts.TaskQueued), utils.Point("cluster-01"), []*TaskLog{{StartTime: &now}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
