package application

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/pbergman/logger"
)

func NewClient(token string, logger *logger.Logger) *http.Client {
	return &http.Client{
		Transport: apiTransport{
			r: http.DefaultTransport,
			t: token,
			l: logger.WithName("client"),
		},
	}
}

type apiTransport struct {
	r http.RoundTripper
	t string
	l *logger.Logger
}

func (a apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Plex-Token", a.t)

	a.printRequest(req)

	resp, err := a.r.RoundTrip(req)

	if err != nil {
		return nil, err
	}

	a.printResponse(resp)

	return resp, nil
}

func (a apiTransport) printRequest(request *http.Request) {
	o, _ := httputil.DumpRequest(request, true)
	s := strings.Split(string(o), "\n")
	for i, c := 0, len(s); i < c; i++ {
		a.l.Debug(">> " + s[i])
	}
}

func (a apiTransport) printResponse(response *http.Response) {
	o, _ := httputil.DumpResponse(response, true)
	s := strings.Split(string(o), "\n")
	for i, c := 0, len(s); i < c; i++ {
		a.l.Debug("<< " + s[i])
	}
}
