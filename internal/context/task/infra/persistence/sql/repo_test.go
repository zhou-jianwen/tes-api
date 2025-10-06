package sql

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/task/domain"
	"github.com/GBA-BI/tes-api/pkg/consts"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
	"github.com/GBA-BI/tes-api/pkg/testutil"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

var id = "task-1234"
var now = time.Now().UTC().Truncate(time.Second)

var taskPO = &Task{
	TaskBasic: TaskBasic{
		TaskStatus: TaskStatus{
			TaskState: TaskState{
				ID:    id,
				State: consts.TaskRunning,
			},
			Logs: []*TaskLog{{
				ClusterID: "cluster-01",
				Logs: [][]*ExecutorLog{{{
					ExecutorID: "executor-01",
					StartTime:  &now,
					EndTime:    &now,
				}}},
				StartTime:  &now,
				SystemLogs: []string{"system log"},
			}},
			CreationTime:          now,
			ClusterID:             utils.Point("cluster-01"),
			StatusResourceVersion: 0,
		},
		Name:        "name",
		Description: "description",
		Resources: &Resources{
			CPUCores: 1,
			RamGB:    2,
			DiskGB:   10,
			GPUType:  utils.Point("gpu-01"),
			GPUCount: utils.Point[float64](2),
		},
		Executors: []*Executor{{
			Image:   "image:tag",
			Command: []string{"command"},
			Workdir: "/base/workdir",
			Stdin:   "/base/stdin",
			Stdout:  "/base/stdout",
			Stderr:  "/base/stderr",
			Env:     map[string]string{"abc": "def"},
		}},
		Volumes: []string{"volume"},
		Tags:    map[string]string{"kkk": "vvv"},
		BioosInfo: &BioosInfo{
			AccountID:    "account-01",
			UserID:       "user-01",
			SubmissionID: "submission-01",
			RunID:        "run-01",
			Meta: &BioosInfoMeta{
				AAIPassport: utils.Point("aai-passport"),
				MountTOS:    utils.Point(false),
				BucketsAuthInfo: &BucketsAuthInfo{
					ReadOnly:  []string{"ro"},
					ReadWrite: []string{"rw"},
					External: []*ExternalBucketAuthInfo{{
						Bucket: "bucket",
						AK:     "ak",
						SK:     "sk",
					}},
				},
			},
		},
		PriorityValue: 100,
	},
	Inputs: []*Input{{
		Name:        "filein",
		Description: "filein description",
		Path:        "/base/filein.txt",
		Type:        "FILE",
		URL:         "s3://abc/filein.txt",
		Content:     "cccc",
	}},
	Outputs: []*Output{{
		Name:        "fileout",
		Description: "fileout description",
		Path:        "/base/fileout.txt",
		Type:        "FILE",
		URL:         "s3://abc/fileout.txt",
	}},
}

var taskDO = &domain.Task{
	TaskStatus: domain.TaskStatus{
		ID:    id,
		State: consts.TaskRunning,
		Logs: []*domain.TaskLog{{
			ClusterID: "cluster-01",
			Logs: [][]*domain.ExecutorLog{{{
				ExecutorID: "executor-01",
				StartTime:  &now,
				EndTime:    &now,
			}}},
			StartTime:  &now,
			SystemLogs: []string{"system log"},
		}},
		CreationTime:          now,
		ClusterID:             "cluster-01",
		StatusResourceVersion: 0,
	},
	Name:        "name",
	Description: "description",
	Inputs: []*domain.Input{{
		Name:        "filein",
		Description: "filein description",
		Path:        "/base/filein.txt",
		Type:        "FILE",
		URL:         "s3://abc/filein.txt",
		Content:     "cccc",
	}},
	Outputs: []*domain.Output{{
		Name:        "fileout",
		Description: "fileout description",
		Path:        "/base/fileout.txt",
		Type:        "FILE",
		URL:         "s3://abc/fileout.txt",
	}},
	Resources: &domain.Resources{
		CPUCores: 1,
		RamGB:    2,
		DiskGB:   10,
		GPU: &domain.GPUResource{
			Count: 2,
			Type:  "gpu-01",
		},
	},
	Executors: []*domain.Executor{{
		Image:   "image:tag",
		Command: []string{"command"},
		Workdir: "/base/workdir",
		Stdin:   "/base/stdin",
		Stdout:  "/base/stdout",
		Stderr:  "/base/stderr",
		Env:     map[string]string{"abc": "def"},
	}},
	Volumes: []string{"volume"},
	Tags:    map[string]string{"kkk": "vvv"},
	BioosInfo: &domain.BioosInfo{
		AccountID:    "account-01",
		UserID:       "user-01",
		SubmissionID: "submission-01",
		RunID:        "run-01",
		Meta: &domain.BioosInfoMeta{
			AAIPassport: utils.Point("aai-passport"),
			MountTOS:    utils.Point(false),
			BucketsAuthInfo: &domain.BucketsAuthInfo{
				ReadOnly:  []string{"ro"},
				ReadWrite: []string{"rw"},
				External: []*domain.ExternalBucketAuthInfo{{
					Bucket: "bucket",
					AK:     "ak",
					SK:     "sk",
				}},
			},
		},
	},
	PriorityValue: 100,
}

var taskStateRows = []string{"id", "state"}
var taskStatusRows = append(taskStateRows, []string{"logs", "creation_time", "cluster_id", "status_resource_version"}...)
var taskBasicRow = append(taskStatusRows, []string{"name", "description",
	"cpu_cores", "ram_gb", "disk_gb", "boot_disk_gb", "gpu_count", "gpu_type", "executors", "volumes", "tags",
	"account_id", "user_id", "submission_id", "run_id", `meta`, "priority_value"}...)
var taskRows = append(taskBasicRow, []string{"inputs", "outputs"}...)

func TestCreate(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("INSERT INTO `task` %s", testutil.GenInsertSql(taskRows))).
		WithArgs(taskPO.ID, taskPO.State,
			testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion, taskPO.Name, taskPO.Description,
			taskPO.Resources.CPUCores, taskPO.Resources.RamGB, taskPO.Resources.DiskGB, taskPO.Resources.BootDiskGB,
			taskPO.Resources.GPUCount, taskPO.Resources.GPUType,
			testutil.MustJSONMarshal(taskPO.Executors), testutil.MustJSONMarshal(taskPO.Volumes),
			testutil.MustJSONMarshal(taskPO.Tags),
			taskPO.BioosInfo.AccountID, taskPO.BioosInfo.UserID, taskPO.BioosInfo.SubmissionID,
			taskPO.BioosInfo.RunID,
			testutil.MustJSONMarshal(taskPO.BioosInfo.Meta), taskPO.PriorityValue,
			testutil.MustJSONMarshal(taskPO.Inputs), testutil.MustJSONMarshal(taskPO.Outputs)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err := r.Create(context.TODO(), taskDO)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}

func TestGetStatus(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStatusRows))).
		WithArgs(id).WillReturnRows(sqlmock.NewRows(taskStatusRows).AddRow(taskPO.ID, taskPO.State,
		testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID, taskPO.StatusResourceVersion))
	resp, err := r.GetStatus(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&taskDO.TaskStatus))
}

func TestGetStatusNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStatusRows))).
		WithArgs(id).WillReturnRows(sqlmock.NewRows(taskStatusRows))
	_, err := r.GetStatus(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestUpdateStatus(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("UPDATE `task` SET %s WHERE `id` = ? AND `status_resource_version` = ?", testutil.GenUpdateSql(taskStatusRows))).
		WithArgs(taskPO.ID, taskPO.State, testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion+1, id, taskPO.StatusResourceVersion).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	updated, err := r.UpdateStatus(context.TODO(), &taskDO.TaskStatus)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updated).To(gomega.BeTrue())
}

func TestUpdateStatusConflict(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectBegin()
	mock.ExpectExec(fmt.Sprintf("UPDATE `task` SET %s WHERE `id` = ? AND `status_resource_version` = ?", testutil.GenUpdateSql(taskStatusRows))).
		WithArgs(taskPO.ID, taskPO.State, testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion+1, id, taskPO.StatusResourceVersion).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	updated, err := r.UpdateStatus(context.TODO(), &taskDO.TaskStatus)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(updated).To(gomega.BeFalse())
}

func TestCheckIDExist(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &repo{db: gormDB}
	mock.ExpectQuery("SELECT count(*) FROM `task` WHERE `id` = ?").WithArgs(id).
		WillReturnRows(testutil.NewCountRows(0))
	exist, err := r.CheckIDExist(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(exist).To(gomega.BeFalse())
}
