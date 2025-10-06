package validator

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestValidateABSPath(t *testing.T) {
	g := gomega.NewWithT(t)

	type Obj struct {
		Path string `validate:"omitempty,abspath"`
	}

	tests := []struct {
		name     string
		path     string
		expValid bool
	}{
		{
			name:     "empty",
			path:     "",
			expValid: true,
		},
		{
			name:     "normal root",
			path:     "/",
			expValid: true,
		},
		{
			name:     "normal abs path",
			path:     "/abc/dev",
			expValid: true,
		},
		{
			name:     "invalid",
			path:     "ttt/ddd",
			expValid: false,
		},
	}

	RegisterValidators()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(&Obj{Path: test.path})
			g.Expect(test.expValid).To(gomega.Equal(err == nil))
		})
	}
}
