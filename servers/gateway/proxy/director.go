package proxy

import (
	"net/http"
	"net/url"
)

type Director func(r *http.Request)

func CustomDirector(targets []*url.URL) Director {
	counter := 0

	return func(r *http.Request) {
		targ := targets[counter%len(targets)]
		counter++
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Host = targ.Host
		r.URL.Host = targ.Host
		r.URL.Scheme = targ.Scheme
	}
}