package normalize

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/internal/context/task/domain"
	"github.com/GBA-BI/tes-api/pkg/utils"
)

func TestNormalizeCheckPath(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name   string
		task   *domain.Task
		expErr bool
	}{
		{
			name: "invalid base: inputs",
			task: &domain.Task{
				Inputs:  []*domain.Input{{Path: "/basetxt.txt"}},
				Outputs: []*domain.Output{{Path: "/base/txt.txt"}},
			},
			expErr: true,
		},
		{
			name: "invalid base: outputs",
			task: &domain.Task{
				Inputs:  []*domain.Input{{Path: "/base/abc/txt.txt"}},
				Outputs: []*domain.Output{{Path: "base/txt.txt"}},
			},
			expErr: true,
		},
		{
			name: "normal inputs and outputs",
			task: &domain.Task{
				Inputs:  []*domain.Input{{Path: "/base/abc/txt.txt"}},
				Outputs: []*domain.Output{{Path: "/base/txt.txt"}},
				Resources: &domain.Resources{
					CPUCores: 4,
					RamGB:    8,
					DiskGB:   20,
				},
			},
			expErr: false,
		},
	}

	n, err := NewNormalizer(&Options{ExecutorBasePath: "/base/"})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = n.Normalize(test.task)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
		})
	}
}

func TestNormalizeSetDefault(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name    string
		task    *domain.Task
		taskExp *domain.Task
	}{
		{
			name: "normal cpu: set default",
			task: &domain.Task{},
			taskExp: &domain.Task{Resources: &domain.Resources{
				CPUCores: 1,
				RamGB:    1,
				DiskGB:   1,
			}},
		},
		{
			name: "normal gpu: set default",
			task: &domain.Task{Resources: &domain.Resources{GPU: &domain.GPUResource{Type: "Type2"}}},
			taskExp: &domain.Task{Resources: &domain.Resources{
				CPUCores: 1,
				RamGB:    1,
				DiskGB:   1,
				GPU:      &domain.GPUResource{Type: "Type2", Count: 1},
			}},
		},
		{
			name: "no set default",
			task: &domain.Task{Resources: &domain.Resources{
				CPUCores: 3,
				RamGB:    4,
				DiskGB:   10,
				GPU:      &domain.GPUResource{Type: "Type1", Count: 3},
			}},
			taskExp: &domain.Task{Resources: &domain.Resources{
				CPUCores: 3,
				RamGB:    4,
				DiskGB:   10,
				GPU:      &domain.GPUResource{Type: "Type1", Count: 3},
			}},
		},
	}
	n, err := NewNormalizer(&Options{ExecutorBasePath: "/base/"})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = n.Normalize(test.task)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(test.task).To(gomega.BeEquivalentTo(test.taskExp))
		})
	}
}

func TestNormalizeDiskGB(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name    string
		task    *domain.Task
		taskExp *domain.Task
	}{
		{
			name:    "normal",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30}},
		},
		{
			name:    "less than min",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 10}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 20}},
		},
		{
			name:    "more than max",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 10000}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 8192}},
		},
		{
			name:    "integer",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30.4}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 31}},
		},
	}

	n, err := NewNormalizer(&Options{
		ExecutorBasePath: "/base/",
		DiskGB: DiskGBOptions{
			Enable:    true,
			Min:       20,
			Max:       8192,
			IsInteger: true,
		},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = n.Normalize(test.task)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(test.task).To(gomega.BeEquivalentTo(test.taskExp))
		})
	}
}

func TestNormalizeBootDiskGB(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name    string
		task    *domain.Task
		taskExp *domain.Task
	}{
		{
			name:    "normal, nil",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30, BootDiskGB: nil}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30, BootDiskGB: nil}},
		},
		{
			name:    "normal",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30, BootDiskGB: utils.Point(50)}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30, BootDiskGB: utils.Point(50)}},
		},
		{
			name:    "less than min",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 20, BootDiskGB: utils.Point(10)}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 20, BootDiskGB: utils.Point(40)}},
		},
		{
			name:    "more than max",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 20, BootDiskGB: utils.Point(200)}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 20, BootDiskGB: utils.Point(100)}},
		},
	}

	n, err := NewNormalizer(&Options{
		ExecutorBasePath: "/base/",
		BootDiskGB: BootDiskGBOptions{
			Enable: true,
			Min:    40,
			Max:    100,
		},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = n.Normalize(test.task)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(test.task).To(gomega.BeEquivalentTo(test.taskExp))
		})
	}
}

func TestNormalizeGPU(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name    string
		task    *domain.Task
		taskExp *domain.Task
	}{
		{
			name:    "no gpu",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30}},
		},
		{
			name:    "integer",
			task:    &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30, GPU: &domain.GPUResource{Count: 0.2}}},
			taskExp: &domain.Task{Resources: &domain.Resources{CPUCores: 1, RamGB: 2, DiskGB: 30, GPU: &domain.GPUResource{Count: 1}}},
		},
	}

	n, err := NewNormalizer(&Options{
		ExecutorBasePath: "/base/",
		GPU: GPUOptions{
			Enable:    true,
			IsInteger: true,
		},
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = n.Normalize(test.task)
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(test.task).To(gomega.BeEquivalentTo(test.taskExp))
		})
	}
}
