package normalize

import (
	"math"
	"strings"

	applog "github.com/GBA-BI/tes-api/pkg/log"

	"github.com/GBA-BI/tes-api/internal/context/task/domain"
	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// normalizer ...
type normalizer struct {
	opts *Options
}

// NewNormalizer ...
func NewNormalizer(opts *Options) (domain.Normalizer, error) {
	res := &normalizer{opts: opts}
	return res, nil
}

var _ domain.Normalizer = (*normalizer)(nil)

// Normalize ...
func (n *normalizer) Normalize(task *domain.Task) error {
	if err := checkPath(task, n.opts.ExecutorBasePath); err != nil {
		return err
	}

	setDefaultResources(task)

	normalizeDiskGB(task, n.opts.DiskGB)
	normalizeBootDiskGB(task, n.opts.BootDiskGB)
	normalizeGPU(task, n.opts.GPU)

	return nil
}

func checkPath(task *domain.Task, executorBasePath string) error {
	for _, input := range task.Inputs {
		if input == nil {
			continue
		}
		if !strings.HasPrefix(input.Path, executorBasePath) {
			return apperrors.NewInvalidError("input.path should under " + executorBasePath)
		}
	}
	for _, output := range task.Outputs {
		if output == nil {
			continue
		}
		if !strings.HasPrefix(output.Path, executorBasePath) {
			return apperrors.NewInvalidError("output.path should under " + executorBasePath)
		}
	}
	return nil
}

func setDefaultResources(task *domain.Task) {
	if task.Resources == nil {
		task.Resources = &domain.Resources{}
	}
	if task.Resources.CPUCores == 0 {
		task.Resources.CPUCores = 1
		applog.Warnw("set cpu default", "task", task.ID, "cpu", task.Resources.CPUCores)
	}
	if task.Resources.RamGB == 0 {
		task.Resources.RamGB = 1
		applog.Warnw("set ram default", "task", task.ID, "ram", task.Resources.RamGB)
	}
	if task.Resources.DiskGB == 0 {
		task.Resources.DiskGB = 1
		applog.Warnw("set disk default", "task", task.ID, "disk", task.Resources.DiskGB)
	}
	if task.Resources.GPU != nil && task.Resources.GPU.Count == 0 {
		task.Resources.GPU.Count = 1
		applog.Warnw("set gpu count default", "task", task.ID, "gpuCount", task.Resources.GPU.Count)
	}
}

func normalizeDiskGB(task *domain.Task, options DiskGBOptions) {
	if !options.Enable {
		return
	}
	newDiskGB := task.Resources.DiskGB
	if task.Resources.DiskGB < options.Min {
		newDiskGB = options.Min
	} else if task.Resources.DiskGB > options.Max {
		newDiskGB = options.Max
	}
	if options.IsInteger {
		newDiskGB = math.Ceil(newDiskGB)
	}
	if task.Resources.DiskGB != newDiskGB {
		applog.Warnw("normalize diskGB", "task", task.ID, "origin", task.Resources.DiskGB, "new", newDiskGB)
	}
	task.Resources.DiskGB = newDiskGB
	return
}

func normalizeBootDiskGB(task *domain.Task, options BootDiskGBOptions) {
	if !options.Enable {
		return
	}
	if task.Resources.BootDiskGB == nil {
		return
	}
	newBootDiskGB := *task.Resources.BootDiskGB
	if *task.Resources.BootDiskGB < options.Min {
		newBootDiskGB = options.Min
	} else if *task.Resources.BootDiskGB > options.Max {
		newBootDiskGB = options.Max
	}
	if *task.Resources.BootDiskGB != newBootDiskGB {
		applog.Warnw("normalize bootDiskGB", "task", task.ID, "origin", *task.Resources.BootDiskGB, "new", newBootDiskGB)
	}
	task.Resources.BootDiskGB = &newBootDiskGB
	return
}

func normalizeGPU(task *domain.Task, options GPUOptions) {
	if !options.Enable {
		return
	}
	if task.Resources.GPU == nil || task.Resources.GPU.Count == 0 {
		return
	}
	newGPUCount := task.Resources.GPU.Count
	if options.IsInteger {
		newGPUCount = math.Ceil(newGPUCount)
	}
	if task.Resources.GPU.Count != newGPUCount {
		applog.Warnw("normalize gpu count", "task", task.ID, "origin", task.Resources.GPU.Count, "new", newGPUCount)
	}
	task.Resources.GPU.Count = newGPUCount
}
