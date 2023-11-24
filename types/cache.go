package types

import (
	"strconv"
	"time"
)

type CacheInfo struct {
	At     CacheTime `json:"cached_at"`
	Until  CacheTime `json:"cached_until"`
	Status string    `json:"status"`
}

func (c CacheInfo) TotalDuration() time.Duration {
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
