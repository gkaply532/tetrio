package tetrio

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

// TODO(gkm): Create sessions on demand instead of using a single, global one.
var sessionID = "SESS-" + must(randomHexString(5))
var limiter = rate.NewLimiter(rate.Every(2*time.Second), 2)

const apiBase = "https://ch.tetr.io/api"

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

func randomHexString(n int) (string, error) {
	buf := make([]byte, n)
	got, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	if got != n {
		return "", errors.New("randomHexString: couldn't read random bytes")
	}
	return hex.EncodeToString(buf[:got]), nil
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// makeRequest sends a get request to the specified path rooted from apiBase, it
// also sets the session header and checks the status code.
func makeRequest(path string) (io.ReadCloser, error) {
	limiter.Wait(context.TODO())
	log.Printf("GET: %q\n", path)
	req, err := http.NewRequest("GET", apiBase+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Session-Id", sessionID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// The body is most likely an error page.
		buf := &bytes.Buffer{}
		err := resp.Write(buf)
		if err != nil {
			return nil, err
		}
		return nil, StatusError{
			Code:    resp.StatusCode,
			Content: buf.Bytes(),
		}
	}
	return resp.Body, nil
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

func send[T any](path string) (T, error) {
	var z T
	resp, err := makeRequest(path)
	if err != nil {
		return z, err
	}
	defer resp.Close()
	return parseDataFromJSON[T](resp)
}
