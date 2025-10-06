package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

func TestList(t *testing.T) {
	g := gomega.NewWithT(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeReadModel := NewFakeReadModel(ctrl)
	fakeReadModel.EXPECT().List(gomock.Any(), &ListFilter{AccountID: "ac1"}).
		Return([]*ExtraPriority{{
			AccountID:          "ac1",
			UserID:             "u1",
			SubmissionID:       "",
			RunID:              "",
			ExtraPriorityValue: 1000,
		}}, nil)

	handler := NewListHandler(fakeReadModel)
	resp, err := handler.Handle(context.TODO(), &ListQuery{Filter: &ListFilter{AccountID: "ac1"}})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(resp).To(gomega.HaveLen(1))
}

func TestListValidate(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name         string
		accountID    string
		submissionID string
		runID        string
		expErr       bool
	}{
		{
			name:   "all empty",
			expErr: false,
		},
		{
			name:      "normal: account",
			accountID: "ac1",
			expErr:    false,
		},
		{
			name:         "normal: submission",
			submissionID: "sb1",
			expErr:       false,
		},
		{
			name:   "normal: run",
			runID:  "r1",
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
			q := &ListQuery{Filter: &ListFilter{
				AccountID:    test.accountID,
				SubmissionID: test.submissionID,
				RunID:        test.runID,
			}}
			err := q.validate()
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
		})
	}
}
