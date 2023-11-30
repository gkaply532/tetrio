// Package tlhist is an API client for @p1nkl0bst3r's https://api.p1nkl0bst3r.xyz/tlhist/ TETRA LEAGUE history API.
package tlhist

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
)

const Endpoint = "https://api.p1nkl0bst3r.xyz/tlhist/"

func RawData(ctx context.Context, uid string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", Endpoint+url.PathEscape(uid), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
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
