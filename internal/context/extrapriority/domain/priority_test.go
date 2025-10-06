package domain

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestNewExtraPriorityID(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name         string
		accountID    string
		userID       string
		submissionID string
		runID        string
		expID        string
		expErr       bool
	}{
		{
			name:   "all empty",
			expErr: true,
		},
		{
			name:      "normal: account",
			accountID: "ac1",
			expID:     "ac1///",
			expErr:    false,
		},
		{
			name:      "normal: user",
			accountID: "ac1",
			userID:    "u1",
			expID:     "ac1/u1//",
			expErr:    false,
		},
		{
			name:   "empty accountID with non-empty userID",
			userID: "u1",
			expErr: true,
		},
		{
			name:         "normal: submission",
			submissionID: "sb1",
			expID:        "//sb1/",
			expErr:       false,
		},
		{
			name:   "normal: run",
			runID:  "r1",
			expID:  "///r1",
			expErr: false,
		},
		{
			name:         "account with submission",
			accountID:    "ac1",
			submissionID: "sb1",
			expErr:       true,
		},
		{
			name:      "account with run",
			accountID: "ac1",
			runID:     "r1",
			expErr:    true,
		},
		{
			name:         "submission with run",
			submissionID: "sb1",
			runID:        "r1",
			expErr:       true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewExtraPriorityID(test.accountID, test.userID, test.submissionID, test.runID)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
			g.Expect(id).To(gomega.Equal(test.expID))
		})
	}
}
