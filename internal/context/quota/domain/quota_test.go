package domain

import (
	"testing"

	"github.com/onsi/gomega"

	"github.com/GBA-BI/tes-api/pkg/consts"
)

func TestNewQuotaID(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name      string
		global    bool
		accountID string
		userID    string
		expID     string
		expErr    bool
	}{
		{
			name:   "normal: global",
			global: true,
			expID:  consts.GlobalQuotaID,
			expErr: false,
		},
		{
			name:      "global with non-empty accountID",
			global:    true,
			accountID: "aaa",
			expErr:    true,
		},
		{
			name:   "global with non-empty userID",
			global: true,
			userID: "bbb",
			expErr: true,
		},
		{
			name:      "global with non-empty accountID and userID",
			global:    true,
			accountID: "aaa",
			userID:    "bbb",
			expErr:    true,
		},
		{
			name:   "empty accountID and userID",
			expErr: true,
		},
		{
			name:   "empty accountID and non-empty userID",
			userID: "bbb",
			expErr: true,
		},
		{
			name:      "normal: default account",
			accountID: consts.DefaultQuotaAccountID,
			expID:     consts.DefaultQuotaAccountID,
			expErr:    false,
		},
		{
			name:      "default account with non-empty userID",
			accountID: consts.DefaultQuotaAccountID,
			userID:    "bbb",
			expErr:    true,
		},
		{
			name:      "normal: accountID with empty userID",
			accountID: "aaa",
			expID:     "aaa/",
			expErr:    false,
		},
		{
			name:      "normal: accountID with non-empty userID",
			accountID: "aaa",
			userID:    "bbb",
			expID:     "aaa/bbb",
			expErr:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id, err := NewQuotaID(test.global, test.accountID, test.userID)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
			g.Expect(id).To(gomega.Equal(test.expID))
		})
	}
}
