package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/chauhanr/ej/cmd/mock"
)

var tcase200 = struct {
	url  string
	Case mockJsonResponse
}{
	url:  "http://almsmart.hclets.com",
	Case: mockJsonResponse{Issue: "1234", Project: "projecta"},
}

func TestGETMethodClient200(t *testing.T) {

	hc, teardown := mock.TestingHttpClient(http.HandlerFunc(GetMethodHandler200))
	defer teardown()
	cli := HttpClient{Client: hc}
	ej := EJConfig{}
	res, code := cli.GET(tcase200.url, ej)
	if code != http.StatusOK {
		t.Errorf("Expected Response %v, but got status %d\n", tcase200.Case, code)
	}
	rj := mockJsonResponse{}
	err := json.Unmarshal([]byte(res), &rj)
	if err != nil {
		t.Errorf("Marshalling error %s\n", err)
	}
	if rj.Issue != tcase200.Case.Issue || rj.Project != tcase200.Case.Project {
		t.Errorf("Expected response %v, but got %v \n", tcase200.Case, rj)
	}
}

func TestGETMethodClient403(t *testing.T) {
	hc, teardown := mock.TestingHttpClient(http.HandlerFunc(GetMethodHandler403))
	defer teardown()
	cli := HttpClient{Client: hc}
	ej := EJConfig{}
	_, code := cli.GET(tcase200.url, ej)
	if code != http.StatusForbidden {
		t.Errorf("Expected status code %d but got code %d", http.StatusForbidden, code)
	}

}

func TestGETMethodClient401(t *testing.T) {
	hc, teardown := mock.TestingHttpClient(http.HandlerFunc(GetMethodHandler401))
	defer teardown()
	cli := HttpClient{Client: hc}
	ej := EJConfig{}
	_, code := cli.GET(tcase200.url, ej)
	if code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got code %d", http.StatusUnauthorized, code)
	}

}

func GetMethodHandler403(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
}

func GetMethodHandler401(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func GetMethodHandler200(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(tcase200.Case)
	_, err := w.Write(buf.Bytes())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

type mockJsonResponse struct {
	Issue   string `json:"issue"`
	Project string `json:"project"`
}
