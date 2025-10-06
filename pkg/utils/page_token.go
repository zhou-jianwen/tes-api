package utils

import (
	"encoding/base64"
	"encoding/json"

	applog "github.com/GBA-BI/tes-api/pkg/log"

	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

// PageToken ...
type PageToken struct {
	LastID string `json:"last_id"`
}

// GenPageToken ...
func GenPageToken(token *PageToken) string {
	if token == nil {
		return ""
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil { // never happen
		applog.Errorw("marshal pageToken", "err", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(tokenBytes)
}

// ParsePageToken ...
func ParsePageToken(tokenStr string) (token *PageToken, err error) {
	if tokenStr == "" {
		return nil, nil
	}

	defer func() {
		if err != nil {
			applog.Errorw("parse pageToken", "err", err)
			err = apperrors.NewInvalidError("pageToken")
		}
	}()

	tokenBytes, err := base64.StdEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, err
	}

	token = new(PageToken)
	if err = json.Unmarshal(tokenBytes, token); err != nil {
		return nil, err
	}
	return token, nil
}
