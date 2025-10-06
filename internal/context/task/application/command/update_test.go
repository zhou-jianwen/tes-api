package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

func TestUpdate(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now().UTC().Truncate(time.Second)

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Update(gomock.Any(), "task-1111", utils.Point(consts.TaskQueued), utils.Point(""), gomock.Any()).
		Return(nil)

	handler := NewUpdateHandler(fakeService)
	err := handler.Handle(context.TODO(), &UpdateCommand{
		ID:        "task-1111",
		ClusterID: utils.Point(""),
		State:     utils.Point(consts.TaskQueued),
		Logs: []*TaskLog{{
			ClusterID: "cluster-01",
			Logs: [][]*ExecutorLog{{{
				ExecutorID: "executor-01",
				StartTime:  &now,
				EndTime:    &now,
			}}},
			StartTime:  &now,
			EndTime:    &now,
			SystemLogs: []string{"abcd efg"},
		}},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
