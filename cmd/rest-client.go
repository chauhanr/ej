package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	AUTH_HEADER     = "Authorization"
	AUTH_BASIC_TYPE = "Basic "
	CONTENT_TYPE    = "Content-Type"
)

type HttpClient struct {
	Client *http.Client
}

func (h *HttpClient) HEAD(URL string, c EJConfig) (int, error) {
	r, err := http.NewRequest(http.MethodHead, URL, nil)
	r.Header.Set(AUTH_HEADER, AUTH_BASIC_TYPE+c.B64)

	res, err := h.Client.Do(r)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return res.StatusCode, nil
}

/*GET will send params like jql etc.*/
func (h *HttpClient) GET(URL string, c EJConfig) (string, int) {
	//fmt.Printf("URL: %s\n", URL)
	r, err := http.NewRequest(http.MethodGet, URL, nil)
	r.Header.Set(AUTH_HEADER, AUTH_BASIC_TYPE+c.B64)
	r.Header.Set(CONTENT_TYPE, "application/json")

	res, err := h.Client.Do(r)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", http.StatusInternalServerError
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return "", http.StatusInternalServerError
		}
		body := string(data)
		return body, http.StatusOK
	}

	switch res.StatusCode {
	case http.StatusForbidden:
		fmt.Printf("Use do not have the right access to this resource\n")
		return "", http.StatusForbidden
	case http.StatusUnauthorized:
		fmt.Printf("User is not authenticated or has wrong credentials\n")
		return "", http.StatusUnauthorized
	default:
		fmt.Printf("Status Code %d, Internal Server Error occured\n", res.StatusCode)
		return "", http.StatusInternalServerError
	}
}

func (h *HttpClient) POST(URL string, c EJConfig, payload interface{}) (string, int) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(payload)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", http.StatusInternalServerError
	}

	r, err := http.NewRequest(http.MethodPost, URL, &buf)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", http.StatusInternalServerError
	}

	r.Header.Set(AUTH_HEADER, AUTH_BASIC_TYPE+c.B64)
	r.Header.Set(CONTENT_TYPE, "application/json")
	res, err := h.Client.Do(r)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", http.StatusInternalServerError
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return "", http.StatusInternalServerError
		}
		body := string(data)
		return body, http.StatusOK
	}

	switch res.StatusCode {
	case http.StatusForbidden:
		fmt.Printf("Use do not have the right access to this resource")
		return "", http.StatusForbidden
	case http.StatusUnauthorized:
		fmt.Printf("User is not authenticated or has wrong credentials")
		return "", http.StatusUnauthorized
	default:
		fmt.Printf("Internal Server Error occured")
		return "", http.StatusInternalServerError
	}
}
