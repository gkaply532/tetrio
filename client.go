package tetrio

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gkaply532/tetrio/v2/types"
	"golang.org/x/time/rate"
)

var Upstream = "https://ch.tetr.io/api"

var limiter = rate.NewLimiter(rate.Every(time.Second), 1)

type Session struct {
	ID    string
	cache map[string]types.Response
}

func New() Session {
	n := rand.Uint64()
	bs := make([]byte, 4)
	binary.BigEndian.PutUint64(bs, n)
	h := hex.EncodeToString(bs)
	return Session{ID: "SESS-" + h}
}

type StatusError struct {
	Code    int
	Content []byte
}

func (e StatusError) Error() string {
	return fmt.Sprintf("tetrio: non 200 status, %d: %q...", e.Code, e.Content[:32])
}

// Send sends a GET request to the specified path rooted from Upstream, it sets
// the X-Session-ID header and checks the status code.
func (s Session) Send(ctx context.Context, path string) (*types.Response, error) {
	if cached, ok := s.cache[path]; ok && time.Now().Before(cached.Cache.Until.Time) {
		return &cached, nil
	}

	err := limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("GET: %q\n", path)

	req, err := http.NewRequestWithContext(ctx, "GET", Upstream+path, nil)
	if err != nil {
		return nil, err
	}

	if s.ID != "" {
		req.Header.Set("X-Session-ID", s.ID)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// The body is most likely an error page.
		bs, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}

		return nil, StatusError{
			Code:    resp.StatusCode,
			Content: bs,
		}
	}

	var result types.Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	s.cache[path] = result
	return &result, nil
}

type Error string // The API returns string errors so... not much can be done.

func (e Error) Error() string {
	return fmt.Sprintf("tetrio: API error: %#q", string(e))
}

func UnmarshalResponse(resp *types.Response, v any) error {
	if resp == nil {
		return errors.New("tetrio: nil response")
	}

	if !resp.Success {
		return Error(resp.Error)
	}

	data := resp.Data

	var m map[string]json.RawMessage
	err := json.Unmarshal(data, &m)
	if err == nil && len(m) == 1 {
		for _, u := range m {
			data = u
		}
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func do[T any](ctx context.Context, s Session, path string) (T, error) {
	var z T

	resp, err := s.Send(ctx, path)
	if err != nil {
		return z, err
	}

	var result T
	err = UnmarshalResponse(resp, &result)
	if err != nil {
		return z, err
	}

	return result, nil
}
