package tetrio

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net/http"
)

type Session struct {
	c  *Client
	ID string
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

func newSession(c *Client) Session {
	return Session{
		c:  c,
		ID: "SESS-" + must(randomHexString(5)),
	}
}

// makeRequest sends a get request to the specified path rooted from apiBase, it
// also sets the session header and checks the status code.
func (s Session) makeRequest(ctx context.Context, path string) (io.ReadCloser, error) {
	s.c.limiter.Wait(ctx)
	log.Printf("GET: %q\n", path)
	req, err := http.NewRequestWithContext(ctx, "GET", s.c.apiBase+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Session-Id", s.ID)
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

func send[T any](ctx context.Context, s Session, path string) (T, error) {
	var z T
	resp, err := s.makeRequest(ctx, path)
	if err != nil {
		return z, err
	}
	defer resp.Close()
	return parseDataFromJSON[T](resp)
}
