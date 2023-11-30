// Package tlhist is an API client for @p1nkl0bst3r's https://api.p1nkl0bst3r.xyz/tlhist/ TETRA LEAGUE history API.
package tlhist

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

const Endpoint = "https://api.p1nkl0bst3r.xyz/tlhist/"

var ErrInvalidUID = errors.New("tlhist: invalid uid")
var ErrNotFound = errors.New("tlhist: not found")

var idRegexp = regexp.MustCompile("^[a-fA-F0-9]{24}$")

func isValidUID(id string) bool {
	return idRegexp.MatchString(id)
}

func RawData(ctx context.Context, uid string) (r io.ReadCloser, err error) {
	if !isValidUID(uid) {
		return nil, ErrInvalidUID
	}

	req, err := http.NewRequestWithContext(ctx, "GET", Endpoint+url.PathEscape(uid), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return resp.Body, nil
}

func AllRecords(ctx context.Context, uid string) ([]Record, error) {
	stream, err := RawData(ctx, uid)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	buffered := bufio.NewReader(stream)

	return NewReader(buffered).All(true)
}
