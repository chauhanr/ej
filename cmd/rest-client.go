package cmd

import (
	"net/http"
)

type HttpClient struct {
	Client *http.Client
}

func (h *HttpClient) HEAD(url string, c EJConfig) (int, error) {
	r, err := http.NewRequest(http.MethodHead, url, nil)
	r.Header.Set("Authorization", "Basic "+c.B64)

	res, err := h.Client.Do(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return res.StatusCode, nil
}
