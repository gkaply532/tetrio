package tetrio

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	apiBase    string
	session    Session // The default session.
	httpClient *http.Client
	limiter    *rate.Limiter
}

func NewClient(h *http.Client) *Client {
	client := &Client{
		apiBase:    "https://ch.tetr.io/api",
		limiter:    rate.NewLimiter(rate.Every(time.Second), 1),
		httpClient: h,
	}
	client.session = newSession(client)
	return client
}

func (c *Client) Session() Session {
	return newSession(c)
}
