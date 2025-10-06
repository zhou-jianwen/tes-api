package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/pkg/consts"
)

func TestGather(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().GatherResources(gomock.Any(), &GatherFilter{
		State:       []string{consts.TaskQueued},
		ClusterID:   "",
		WithCluster: true,
		AccountID:   "ac1",
		UserID:      "",
	}).Return(&TasksResources{
		Count:    1,
		CPUCores: 2,
		RamGB:    5,
		DiskGB:   10,
		GPU:      map[string]float64{"gpu-01": 1},
	}, nil)

	handler := NewGatherHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO(), &GatherQuery{Filter: &GatherFilter{
		State:       []string{consts.TaskQueued},
		ClusterID:   "",
		WithCluster: true,
		AccountID:   "ac1",
		UserID:      "",
	}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&TasksResources{
		Count:    1,
		CPUCores: 2,
		RamGB:    5,
		DiskGB:   10,
		GPU:      map[string]float64{"gpu-01": 1},
	}))
}
