package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/gkaply532/tetrio/v2/types"
	"golang.org/x/time/rate"
)

var upstreamURL = must(url.Parse("https://ch.tetr.io/api"))
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 1)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func main() {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(upstreamURL)
			log.Println(pr.Out.Method, pr.Out.URL)
		},
		ModifyResponse: func(r *http.Response) error {
			buf := &bytes.Buffer{}

			_, err := io.Copy(buf, r.Body)
			if err != nil {
				return err
			}

			err = r.Body.Close()
			if err != nil {
				return err
			}

			r.Body = io.NopCloser(buf)

			var resp types.Response
			err = json.Unmarshal(buf.Bytes(), &resp)
			if err != nil {
				// We don't really care that we couldn't to parse the response.
				// Just return the response without more modifications.
				return nil
			}

			maxage := int64(math.Ceil(resp.Cache.TotalDuration().Seconds()))
			age := int64(time.Since(resp.Cache.At.Time).Seconds())
			if age < 0 {
				age = 0
			}
			r.Header.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxage))
			r.Header.Set("Age", strconv.FormatInt(age, 10))

			return nil
		},
	}

	handler := http.StripPrefix("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := limiter.Wait(r.Context())
		if err != nil {
			return
		}
		proxy.ServeHTTP(w, r)
	}))

	http.Handle("/api", handler)
	http.Handle("/api/", handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this is the tetrio-gateway!\nthe endpoints are available in `/api/**`.")
	})

	log.Println("starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
