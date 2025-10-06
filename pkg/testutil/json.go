package testutil

import "encoding/json"

// MustJSONMarshal ...
func MustJSONMarshal(c interface{}) string {
	res, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(res)
}

// MustJSONUnmarshal ...
func MustJSONUnmarshal(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}
