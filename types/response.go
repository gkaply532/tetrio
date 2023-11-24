package types

import (
	"encoding/json"
)

type Response struct {
	Success bool            `json:"success"`
	Error   json.RawMessage `json:"error"`
	Data    json.RawMessage `json:"data"`
	Cache   CacheInfo       `json:"cache"`
}
