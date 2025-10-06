package sql

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/task/application/query"
	"code.byted.org/epscp/vetes-api/pkg/consts"
	apperrors "code.byted.org/epscp/vetes-api/pkg/errors"
	"code.byted.org/epscp/vetes-api/pkg/testutil"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

var taskDTO = &query.Task{
	TaskBasic: query.TaskBasic{
		TaskMinimal: query.TaskMinimal{
			ID:    id,
			State: consts.TaskRunning,
		},
		Name:        "name",
		Description: "description",
		Resources: &query.Resources{
			CPUCores: 1,
			RamGB:    2,
			DiskGB:   10,
			GPU: &query.GPUResource{
				Count: 2,
				Type:  "gpu-01",
			},
		},
		Executors: []*query.Executor{{
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
		Logs: []*query.TaskLog{{
			ClusterID: "cluster-01",
			Logs: [][]*query.ExecutorLog{{{
				ExecutorID: "executor-01",
				StartTime:  &now,
				EndTime:    &now,
			}}},
			StartTime:  &now,
			SystemLogs: []string{"system log"},
		}},
		CreationTime: now,
		BioosInfo: &query.BioosInfo{
			AccountID:    "account-01",
			UserID:       "user-01",
			SubmissionID: "submission-01",
			RunID:        "run-01",
			Meta: &query.BioosInfoMeta{
				AAIPassport: utils.Point("aai-passport"),
				MountTOS:    utils.Point(false),
				BucketsAuthInfo: &query.BucketsAuthInfo{
					ReadOnly:  []string{"ro"},
					ReadWrite: []string{"rw"},
					External: []*query.ExternalBucketAuthInfo{{
						Bucket: "bucket",
						AK:     "ak",
						SK:     "sk",
					}},
				},
			},
		},
		PriorityValue: 100,
		ClusterID:     "cluster-01",
	},
	Inputs: []*query.Input{{
		Name:        "filein",
		Description: "filein description",
		Path:        "/base/filein.txt",
		Type:        "FILE",
		URL:         "s3://abc/filein.txt",
		Content:     "cccc",
	}},
	Outputs: []*query.Output{{
		Name:        "fileout",
		Description: "fileout description",
		Path:        "/base/fileout.txt",
		Type:        "FILE",
		URL:         "s3://abc/fileout.txt",
	}},
}

func TestListMinimal(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` ORDER BY `id` LIMIT 10",
		testutil.GenSelectFieldsSql("task", taskStateRows))).
		WillReturnRows(sqlmock.NewRows(taskStateRows).AddRow(taskPO.ID, taskPO.State))
	resp, nextPageToken, err := r.ListMinimal(context.TODO(), 10, nil, nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(nextPageToken).To(gomega.BeNil())
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.TaskMinimal{&taskDTO.TaskMinimal}))
}

func TestListMinimalWithPageToken(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` > ? ORDER BY `id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStateRows))).
		WithArgs("task-1111").
		WillReturnRows(sqlmock.NewRows(taskStateRows).AddRow(taskPO.ID, taskPO.State))
	resp, nextPageToken, err := r.ListMinimal(context.TODO(), 1, &utils.PageToken{LastID: "task-1111"}, nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(nextPageToken).To(gomega.BeEquivalentTo(&utils.PageToken{LastID: id}))
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.TaskMinimal{&taskDTO.TaskMinimal}))
}

func TestListMinimalWithFilter(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `name` LIKE ? AND `state` IN (?,?) AND `cluster_id` = ? AND `id` > ? ORDER BY `id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStateRows))).
		WithArgs("task\\%1\\_1%", consts.TaskRunning, consts.TaskQueued, "cluster-01", "task-1111").
		WillReturnRows(sqlmock.NewRows(taskStateRows).AddRow(taskPO.ID, taskPO.State))
	resp, nextPageToken, err := r.ListMinimal(context.TODO(), 1, &utils.PageToken{LastID: "task-1111"}, &query.ListFilter{
		NamePrefix: "task%1_1",
		State:      []string{consts.TaskRunning, consts.TaskQueued},
		ClusterID:  "cluster-01",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(nextPageToken).To(gomega.BeEquivalentTo(&utils.PageToken{LastID: id}))
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.TaskMinimal{&taskDTO.TaskMinimal}))
}

func TestListMinimalWithFilterWithoutCluster(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `cluster_id` = '' AND `id` > ? ORDER BY `id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStateRows))).
		WithArgs("task-1111").
		WillReturnRows(sqlmock.NewRows(taskStateRows))
	resp, nextPageToken, err := r.ListMinimal(context.TODO(), 1, &utils.PageToken{LastID: "task-1111"}, &query.ListFilter{
		WithoutCluster: true,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(nextPageToken).To(gomega.BeNil())
	g.Expect(resp).To(gomega.BeEmpty())
}

func TestListBasic(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` ORDER BY `id` LIMIT 10",
		testutil.GenSelectFieldsSql("task", taskBasicRow))).
		WillReturnRows(sqlmock.NewRows(taskBasicRow).AddRow(taskPO.ID, taskPO.State,
			testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion, taskPO.Name, taskPO.Description,
			taskPO.Resources.CPUCores, taskPO.Resources.RamGB, taskPO.Resources.DiskGB, taskPO.Resources.BootDiskGB,
			taskPO.Resources.GPUCount, taskPO.Resources.GPUType,
			testutil.MustJSONMarshal(taskPO.Executors), testutil.MustJSONMarshal(taskPO.Volumes),
			testutil.MustJSONMarshal(taskPO.Tags),
			taskPO.BioosInfo.AccountID, taskPO.BioosInfo.UserID, taskPO.BioosInfo.SubmissionID,
			taskPO.BioosInfo.RunID,
			testutil.MustJSONMarshal(taskPO.BioosInfo.Meta), taskPO.PriorityValue))
	resp, nextPageToken, err := r.ListBasic(context.TODO(), 10, nil, nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(nextPageToken).To(gomega.BeNil())
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.TaskBasic{&taskDTO.TaskBasic}))
}

func TestListFull(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `task` ORDER BY `id` LIMIT 10").
		WillReturnRows(sqlmock.NewRows(taskRows).AddRow(taskPO.ID, taskPO.State,
			testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion, taskPO.Name, taskPO.Description,
			taskPO.Resources.CPUCores, taskPO.Resources.RamGB, taskPO.Resources.DiskGB, taskPO.Resources.BootDiskGB,
			taskPO.Resources.GPUCount, taskPO.Resources.GPUType,
			testutil.MustJSONMarshal(taskPO.Executors), testutil.MustJSONMarshal(taskPO.Volumes),
			testutil.MustJSONMarshal(taskPO.Tags),
			taskPO.BioosInfo.AccountID, taskPO.BioosInfo.UserID, taskPO.BioosInfo.SubmissionID,
			taskPO.BioosInfo.RunID,
			testutil.MustJSONMarshal(taskPO.BioosInfo.Meta), taskPO.PriorityValue,
			testutil.MustJSONMarshal(taskPO.Inputs), testutil.MustJSONMarshal(taskPO.Outputs)))
	resp, nextPageToken, err := r.ListFull(context.TODO(), 10, nil, nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(nextPageToken).To(gomega.BeNil())
	g.Expect(resp).To(gomega.BeEquivalentTo([]*query.Task{taskDTO}))
}

func TestGetMinimal(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStateRows))).WithArgs(id).
		WillReturnRows(sqlmock.NewRows(taskStateRows).AddRow(taskPO.ID, taskPO.State))
	resp, err := r.GetMinimal(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&taskDTO.TaskMinimal))
}

func TestGetMinimalNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskStateRows))).WithArgs(id).
		WillReturnRows(sqlmock.NewRows(taskStateRows))
	_, err := r.GetMinimal(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestGetBasic(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskBasicRow))).WithArgs(id).
		WillReturnRows(sqlmock.NewRows(taskBasicRow).AddRow(taskPO.ID, taskPO.State,
			testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion, taskPO.Name, taskPO.Description,
			taskPO.Resources.CPUCores, taskPO.Resources.RamGB, taskPO.Resources.DiskGB, taskPO.Resources.BootDiskGB,
			taskPO.Resources.GPUCount, taskPO.Resources.GPUType,
			testutil.MustJSONMarshal(taskPO.Executors), testutil.MustJSONMarshal(taskPO.Volumes),
			testutil.MustJSONMarshal(taskPO.Tags),
			taskPO.BioosInfo.AccountID, taskPO.BioosInfo.UserID, taskPO.BioosInfo.SubmissionID,
			taskPO.BioosInfo.RunID,
			testutil.MustJSONMarshal(taskPO.BioosInfo.Meta), taskPO.PriorityValue))
	resp, err := r.GetBasic(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&taskDTO.TaskBasic))
}

func TestGetBasicNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery(fmt.Sprintf("SELECT %s FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1",
		testutil.GenSelectFieldsSql("task", taskBasicRow))).WithArgs(id).
		WillReturnRows(sqlmock.NewRows(taskBasicRow))
	_, err := r.GetBasic(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestGetFull(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(taskRows).AddRow(taskPO.ID, taskPO.State,
			testutil.MustJSONMarshal(taskPO.Logs), taskPO.CreationTime, taskPO.ClusterID,
			taskPO.StatusResourceVersion, taskPO.Name, taskPO.Description,
			taskPO.Resources.CPUCores, taskPO.Resources.RamGB, taskPO.Resources.DiskGB, taskPO.Resources.BootDiskGB,
			taskPO.Resources.GPUCount, taskPO.Resources.GPUType,
			testutil.MustJSONMarshal(taskPO.Executors), testutil.MustJSONMarshal(taskPO.Volumes),
			testutil.MustJSONMarshal(taskPO.Tags),
			taskPO.BioosInfo.AccountID, taskPO.BioosInfo.UserID, taskPO.BioosInfo.SubmissionID,
			taskPO.BioosInfo.RunID,
			testutil.MustJSONMarshal(taskPO.BioosInfo.Meta), taskPO.PriorityValue,
			testutil.MustJSONMarshal(taskPO.Inputs), testutil.MustJSONMarshal(taskPO.Outputs)))
	resp, err := r.GetFull(context.TODO(), id)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(taskDTO))
}

func TestGetFullNotFound(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT * FROM `task` WHERE `id` = ? ORDER BY `task`.`id` LIMIT 1").WithArgs(id).
		WillReturnRows(sqlmock.NewRows(taskRows))
	_, err := r.GetFull(context.TODO(), id)
	g.Expect(apperrors.IsCode(err, apperrors.NotFoundCode)).To(gomega.BeTrue())
}

func TestGatherResources(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT COUNT(*) AS `count`, SUM(`cpu_cores`) AS `cpu_cores`, SUM(`ram_gb`) AS `ram_gb`, SUM(`disk_gb`) AS `disk_gb` FROM `task`").
		WillReturnRows(sqlmock.NewRows([]string{"count", "cpu_cores", "ram_gb", "disk_gb"}).AddRow(5, 10, 20, 30))
	mock.ExpectQuery("SELECT `gpu_type`, SUM(`gpu_count`) AS `gpu_count` FROM `task` WHERE `gpu_type` IS NOT NULL AND `gpu_count` IS NOT NULL GROUP BY `gpu_type`").
		WillReturnRows(sqlmock.NewRows([]string{"gpu_type", "gpu_count"}).AddRow("gpu-01", 2).AddRow("gpu-02", 3))
	resp, err := r.GatherResources(context.TODO(), nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&query.TasksResources{
		Count:    5,
		CPUCores: 10,
		RamGB:    20,
		DiskGB:   30,
		GPU: map[string]float64{
			"gpu-01": 2,
			"gpu-02": 3,
		},
	}))
}

func TestGatherResourcesEmptyGPU(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT COUNT(*) AS `count`, SUM(`cpu_cores`) AS `cpu_cores`, SUM(`ram_gb`) AS `ram_gb`, SUM(`disk_gb`) AS `disk_gb` FROM `task`").
		WillReturnRows(sqlmock.NewRows([]string{"count", "cpu_cores", "ram_gb", "disk_gb"}).AddRow(5, 10, 20, 30))
	mock.ExpectQuery("SELECT `gpu_type`, SUM(`gpu_count`) AS `gpu_count` FROM `task` WHERE `gpu_type` IS NOT NULL AND `gpu_count` IS NOT NULL GROUP BY `gpu_type`").
		WillReturnRows(sqlmock.NewRows([]string{"gpu_type", "gpu_count"}))
	resp, err := r.GatherResources(context.TODO(), nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&query.TasksResources{
		Count:    5,
		CPUCores: 10,
		RamGB:    20,
		DiskGB:   30,
		GPU:      nil,
	}))
}

func TestGatherResourcesEmpty(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT COUNT(*) AS `count`, SUM(`cpu_cores`) AS `cpu_cores`, SUM(`ram_gb`) AS `ram_gb`, SUM(`disk_gb`) AS `disk_gb` FROM `task`").
		WillReturnRows(sqlmock.NewRows([]string{"count", "cpu_cores", "ram_gb", "disk_gb"}))
	mock.ExpectQuery("SELECT `gpu_type`, SUM(`gpu_count`) AS `gpu_count` FROM `task` WHERE `gpu_type` IS NOT NULL AND `gpu_count` IS NOT NULL GROUP BY `gpu_type`").
		WillReturnRows(sqlmock.NewRows([]string{"gpu_type", "gpu_count"}))
	resp, err := r.GatherResources(context.TODO(), nil)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&query.TasksResources{
		Count:    0,
		CPUCores: 0,
		RamGB:    0,
		DiskGB:   0,
		GPU:      nil,
	}))
}

func TestGatherResourcesWithFilterOfCluster(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT COUNT(*) AS `count`, SUM(`cpu_cores`) AS `cpu_cores`, SUM(`ram_gb`) AS `ram_gb`, SUM(`disk_gb`) AS `disk_gb` FROM `task` "+
		"WHERE `state` IN (?,?,?) AND `cluster_id` = ?").
		WithArgs(consts.TaskQueued, consts.TaskRunning, consts.TaskCanceling, "cluster-01").
		WillReturnRows(sqlmock.NewRows([]string{"count", "cpu_cores", "ram_gb", "disk_gb"}).AddRow(5, 10, 20, 30))
	mock.ExpectQuery("SELECT `gpu_type`, SUM(`gpu_count`) AS `gpu_count` FROM `task` "+
		"WHERE `state` IN (?,?,?) AND `cluster_id` = ? "+
		"AND (`gpu_type` IS NOT NULL AND `gpu_count` IS NOT NULL) "+
		"GROUP BY `gpu_type`").
		WithArgs(consts.TaskQueued, consts.TaskRunning, consts.TaskCanceling, "cluster-01").
		WillReturnRows(sqlmock.NewRows([]string{"gpu_type", "gpu_count"}).AddRow("gpu-01", 2).AddRow("gpu-02", 3))
	resp, err := r.GatherResources(context.TODO(), &query.GatherFilter{
		State:     []string{consts.TaskQueued, consts.TaskRunning, consts.TaskCanceling},
		ClusterID: "cluster-01",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&query.TasksResources{
		Count:    5,
		CPUCores: 10,
		RamGB:    20,
		DiskGB:   30,
		GPU: map[string]float64{
			"gpu-01": 2,
			"gpu-02": 3,
		},
	}))
}

func TestGatherResourcesWithFilterOfAccount(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT COUNT(*) AS `count`, SUM(`cpu_cores`) AS `cpu_cores`, SUM(`ram_gb`) AS `ram_gb`, SUM(`disk_gb`) AS `disk_gb` FROM `task` "+
		"WHERE `state` IN (?,?,?) AND `cluster_id` <> '' AND `account_id` = ? AND `user_id` = ?").
		WithArgs(consts.TaskQueued, consts.TaskRunning, consts.TaskCanceling, "account-01", "user-01").
		WillReturnRows(sqlmock.NewRows([]string{"count", "cpu_cores", "ram_gb", "disk_gb"}).AddRow(5, 10, 20, 30))
	mock.ExpectQuery("SELECT `gpu_type`, SUM(`gpu_count`) AS `gpu_count` FROM `task` "+
		"WHERE `state` IN (?,?,?) AND `cluster_id` <> '' AND `account_id` = ? AND `user_id` = ? "+
		"AND (`gpu_type` IS NOT NULL AND `gpu_count` IS NOT NULL) "+
		"GROUP BY `gpu_type`").
		WithArgs(consts.TaskQueued, consts.TaskRunning, consts.TaskCanceling, "account-01", "user-01").
		WillReturnRows(sqlmock.NewRows([]string{"gpu_type", "gpu_count"}).AddRow("gpu-01", 2).AddRow("gpu-02", 3))
	resp, err := r.GatherResources(context.TODO(), &query.GatherFilter{
		State:       []string{consts.TaskQueued, consts.TaskRunning, consts.TaskCanceling},
		WithCluster: true,
		AccountID:   "account-01",
		UserID:      "user-01",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.BeEquivalentTo(&query.TasksResources{
		Count:    5,
		CPUCores: 10,
		RamGB:    20,
		DiskGB:   30,
		GPU: map[string]float64{
			"gpu-01": 2,
			"gpu-02": 3,
		},
	}))
}

func TestListAccounts(t *testing.T) {
	g := gomega.NewWithT(t)
	mock, gormDB := testutil.NewSqlMock()
	r := &readModel{db: gormDB}
	mock.ExpectQuery("SELECT DISTINCT `account_id`,`user_id` FROM `task`").
		WillReturnRows(sqlmock.NewRows([]string{"account_id", "user_id"}).
			AddRow("account-01", "").
			AddRow("account-01", "user-01").
			AddRow("account-02", "user-02").
			AddRow("", ""))
	resp, err := r.ListAccounts(context.TODO())
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.ConsistOf([]*query.AccountInfo{
		{AccountID: "account-01", UserIDs: []string{"", "user-01"}},
		{AccountID: "account-02", UserIDs: []string{"user-02"}},
	}))
}
