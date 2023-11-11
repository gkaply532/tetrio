package tetrio

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

type Response struct {
	Success bool            `json:"success"`
	Error   json.RawMessage `json:"error"`
	Data    json.RawMessage `json:"data"`
	Cache   CacheInfo       `json:"cache"`
}

type CacheInfo struct {
	At     CacheTime `json:"cached_at"`
	Until  CacheTime `json:"cached_until"`
	Status string    `json:"status"`
}

func (c CacheInfo) Duration() time.Duration {
	return c.Until.Sub(c.At.Time)
}

type CacheTime struct {
	time.Time
}

func (ct *CacheTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	t, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	ct.Time = time.UnixMilli(t)
	return nil
}

func (ct CacheTime) MarshalJSON() ([]byte, error) {
	t := ct.UnixMilli()
	str := strconv.FormatInt(t, 10)
	return []byte(str), nil
}

type StatusError struct {
	Code    int
	Content []byte
}

func (e StatusError) Error() string {
	return fmt.Sprintf("tetrio: non 200 status, %d: %q...", e.Code, e.Content[:32])
}

type Error string // The API returns string errors so... not much can be done.

func (e Error) Error() string {
	return fmt.Sprintf("tetrio: API error: %#q", string(e))
}

// The returned responses "data" field is *usually* an object with a single key
// that contains the actual data, this function removes that unnecessary
// wrapping.
func unwrapResponse(data json.RawMessage) json.RawMessage {
	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err == nil && len(m) == 1 {
		// Looks weird but this works and is idiomatic.
		for _, v := range m {
			return v
		}
	}
	// Error or not a single keyed wrapper object. If there was an error then it
	// probably wasn't even an object in that case.
	return data
}

func parseResponse(resp io.Reader) (*Response, error) {
	var result Response
	err := json.NewDecoder(resp).Decode(&result)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, Error(result.Error)
	}
	return &result, nil
}

func parseDataFromJSON[T any](jsonReader io.Reader) (T, error) {
	var z T
	resp, err := parseResponse(jsonReader)
	if err != nil {
		return z, err
	}
	unwrapped := unwrapResponse(resp.Data)
	var result T
	err = json.Unmarshal(unwrapped, &result)
	if err != nil {
		return z, err
	}
	return result, nil
}
