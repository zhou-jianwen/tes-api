package domain

import (
	"testing"
	"time"

	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/pkg/consts"
)

func TestUpdateState(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name     string
		oldState string
		newState string
		expErr   bool
	}{
		{
			name:     "finished -> self",
			oldState: consts.TaskSystemError,
			newState: consts.TaskSystemError,
			expErr:   false,
		},
		{
			name:     "finished -> other: invalid",
			oldState: consts.TaskCanceled,
			newState: consts.TaskExecutorError,
			expErr:   true,
		},
		{
			name:     "canceling -> self",
			oldState: consts.TaskCanceling,
			newState: consts.TaskCanceling,
			expErr:   false,
		},
		{
			name:     "canceling -> executing: invalid",
			oldState: consts.TaskCanceling,
			newState: consts.TaskRunning,
			expErr:   true,
		},
		{
			name:     "canceling -> canceled",
			oldState: consts.TaskCanceling,
			newState: consts.TaskCanceled,
			expErr:   false,
		},
		{
			name:     "other -> canceled: invalid",
			oldState: consts.TaskRunning,
			newState: consts.TaskCanceled,
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := &TaskStatus{State: test.oldState}
			err := task.UpdateState(test.newState)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
			if err == nil {
				g.Expect(task.State).To(gomega.Equal(test.newState))
			}
		})
	}
}

func TestUpdateClusterID(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name         string
		oldState     string
		oldClusterID string
		newClusterID string
		expErr       bool
	}{
		{
			name:         "same cluster",
			oldState:     consts.TaskRunning,
			oldClusterID: "cluster-01",
			newClusterID: "cluster-01",
			expErr:       false,
		},
		{
			name:         "schedule",
			oldState:     consts.TaskQueued,
			oldClusterID: "",
			newClusterID: "cluster-01",
			expErr:       false,
		},
		{
			name:         "change",
			oldState:     consts.TaskQueued,
			oldClusterID: "cluster-01",
			newClusterID: "cluster-02",
			expErr:       false,
		},
		{
			name:         "non queue: invalid",
			oldState:     consts.TaskInitializing,
			oldClusterID: "",
			newClusterID: "cluster-01",
			expErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := &TaskStatus{State: test.oldState, ClusterID: test.oldClusterID}
			err := task.UpdateClusterID(test.newClusterID)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
			if err == nil {
				g.Expect(task.ClusterID).To(gomega.Equal(test.newClusterID))
			}
		})
	}
}

func TestUpdateLogs(t *testing.T) {
	g := gomega.NewWithT(t)

	now := time.Now().UTC().Truncate(time.Second)
	t1 := now
	t2 := t1.Add(time.Second)
	t3 := t2.Add(time.Second)
	t4 := t3.Add(time.Second)
	t5 := t4.Add(time.Second)

	tests := []struct {
		name         string
		creationTime time.Time
		oldLogs      []*TaskLog
		newLogs      []*TaskLog
		expLogs      []*TaskLog
		expErr       bool
	}{
		{
			name:         "normal1: init TaskLog",
			creationTime: now,
			oldLogs:      nil,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
			}},
			expLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
			}},
		},
		{
			name:         "normal2: init ExecutorLog",
			creationTime: now,
			oldLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
			}},
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
				}}},
			}},
			expLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
				}}},
			}},
		},
		{
			name:         "normal3: end ExecutorLog, append another",
			creationTime: now,
			oldLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
				}}},
			}},
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
				}}},
			}},
			expLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
				}}},
			}},
		},
		{
			name:         "normal4: end ExecutorLogs, append another",
			creationTime: now,
			oldLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
				}}},
			}},
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-02",
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}},
			expLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}},
		},
		{
			name:         "normal5: end TaskLog, append another",
			creationTime: now,
			oldLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &now,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}},
			newLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				EndTime:    &now,
				SystemLogs: []string{"abc"},
			}, {
				ClusterID:  "cluster-02",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ddd"},
			}},
			expLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"abc"},
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}, {
				ClusterID:  "cluster-02",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ddd"},
			}},
		},
		{
			name:         "normal6: append systemLogs",
			creationTime: now,
			oldLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"abc"},
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}, {
				ClusterID:  "cluster-02",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ddd"},
			}},
			newLogs: []*TaskLog{{
				ClusterID:  "cluster-02",
				SystemLogs: []string{"", "ttt"},
			}},
			expLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"abc"},
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}, {
				ClusterID:  "cluster-02",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ddd", "ttt"},
			}},
		},
		{
			name:         "normal7, update systemLogs",
			creationTime: now,
			oldLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"abc"},
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}, {
				ClusterID:  "cluster-02",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ddd", "ttt"},
			}},
			newLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				SystemLogs: []string{"ccc"},
			}},
			expLogs: []*TaskLog{{
				ClusterID:  "cluster-01",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ccc"},
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "ex-01-01",
					StartTime:  &now,
					EndTime:    &now,
				}, {
					ExecutorID: "ex-01-02",
					StartTime:  &now,
					EndTime:    &now,
				}}, {{
					ExecutorID: "ex-02-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
			}, {
				ClusterID:  "cluster-02",
				StartTime:  &now,
				EndTime:    &now,
				SystemLogs: []string{"ddd", "ttt"},
			}},
		},
		{
			name:         "normal time",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t2,
				EndTime:   &t5,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t3,
					EndTime:    &t4,
				}}},
			}},
			expLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t2,
				EndTime:   &t5,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t3,
					EndTime:    &t4,
				}}},
			}},
		},
		{
			name:         "invalid executor StartTime after EndTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t2,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t4,
					EndTime:    &t3,
				}}},
			}},
			expErr: true,
		},
		{
			name:         "invalid empty executor StartTime with EndTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t2,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					EndTime:    &t4,
				}}},
			}},
			expErr: true,
		},
		{
			name:         "invalid executor time with empty task StartTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t3,
				}}},
			}},
			expErr: true,
		},
		{
			name:         "invalid executor StartTime before task StartTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t3,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t2,
				}}},
			}},
			expErr: true,
		},
		{
			name:         "invalid empty executor EndTime with taskEndTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t2,
				EndTime:   &t5,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t3,
				}}},
			}},
			expErr: true,
		},
		{
			name:         "invalid executor EndTime after taskEndTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t2,
				EndTime:   &t4,
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-id",
					StartTime:  &t3,
					EndTime:    &t5,
				}}},
			}},
			expErr: true,
		},
		{
			name:         "invalid task StartTime after EndTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t3,
				EndTime:   &t2,
			}},
			expErr: true,
		},
		{
			name:         "invalid empty task StartTime with EndTime",
			creationTime: t1,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				EndTime:   &t3,
			}},
			expErr: true,
		},
		{
			name:         "invalid task StartTime after CreationTime",
			creationTime: t2,
			newLogs: []*TaskLog{{
				ClusterID: "cluster-01",
				StartTime: &t1,
			}},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := &TaskStatus{CreationTime: test.creationTime, Logs: test.oldLogs}
			err := task.UpdateLogs(test.newLogs)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
			if err == nil {
				g.Expect(task.Logs).To(gomega.BeEquivalentTo(test.expLogs))
			}
		})
	}
}
