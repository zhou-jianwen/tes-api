package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"code.byted.org/epscp/vetes-api/internal/context/task/domain"
	"code.byted.org/epscp/vetes-api/pkg/utils"
)

func TestCreate(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := domain.NewFakeService(ctrl)
	fakeService.EXPECT().Create(gomock.Any(), gomock.Any()).
		Return("task-1234", nil)

	handler := NewCreateHandler(fakeService)
	id, err := handler.Handle(context.TODO(), &CreateCommand{
		Name:        "name",
		Description: "description",
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
		Resources: &Resources{
			CPUCores: 1,
			RamGB:    2,
			DiskGB:   10,
			GPU: &GPUResource{
				Count: 2,
				Type:  "gpu-01",
			},
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
				MountTOS:    utils.Point(true),
				BucketsAuthInfo: &BucketsAuthInfo{
					ReadWrite: []string{"rw"},
					ReadOnly:  []string{"ro"},
					External: []*ExternalBucketAuthInfo{{
						Bucket: "buk",
						AK:     "ak",
						SK:     "sk",
					}},
				},
			},
		},
		PriorityValue: 100,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(id).To(gomega.Equal("task-1234"))
}
