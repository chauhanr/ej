package mock

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
)

func TestingHttpClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}
	return cli, s.Close
}
