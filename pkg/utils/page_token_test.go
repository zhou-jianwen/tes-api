package utils

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestGenPageToken(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name  string
		token *PageToken
		exp   string
	}{
		{
			name:  "nil",
			token: nil,
			exp:   "",
		},
		{
			name:  "normal",
			token: &PageToken{LastID: "task-abcd"},
			exp:   "eyJsYXN0X2lkIjoidGFzay1hYmNkIn0=",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g.Expect(GenPageToken(test.token)).To(gomega.Equal(test.exp))
		})
	}
}

func TestParsePageToken(t *testing.T) {
	g := gomega.NewWithT(t)

	tests := []struct {
		name     string
		tokenStr string
		expToken *PageToken
		expErr   bool
	}{
		{
			name:     "empty",
			tokenStr: "",
			expToken: nil,
			expErr:   false,
		},
		{
			name:     "normal",
			tokenStr: "eyJsYXN0X2lkIjoidGFzay1hYmNkIn0=",
			expToken: &PageToken{LastID: "task-abcd"},
			expErr:   false,
		},
		{
			name:     "invalid base64",
			tokenStr: "eyJsYXN0X2lkIjoidGFzay1hYmNkIn0",
			expToken: nil,
			expErr:   true,
		},
		{
			name:     "invalid json",
			tokenStr: "ay1h",
			expToken: nil,
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ParsePageToken(test.tokenStr)
			g.Expect(err != nil).To(gomega.Equal(test.expErr))
			g.Expect(res).To(gomega.BeEquivalentTo(test.expToken))
		})
	}
}
