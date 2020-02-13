package cmd

import (
	"net/http"
	"testing"

	"github.com/chauhanr/ej/cmd/mock"
)

const (
	URL403 = "http://almsmart.hclets.com/403"
	URL401 = "http://almsmart.hclets.com/401"
	URL200 = "http://almsmart.hclets.com"
)

var testConf200 = EJConfig{
	Url:      URL200,
	Username: "Ritesh",
	Password: "password",
	B64:      EncodeCreds("Ritesh", "password"),
}

var testConf403 = EJConfig{
	Url:      URL403,
	Username: "Ritesh",
	Password: "password403",
	B64:      EncodeCreds("Ritesh", "password403"),
}

var testConf401 = EJConfig{
	Url:      URL401,
	Username: "Ritesh",
	Password: "password401",
	B64:      EncodeCreds("Ritesh", "password401"),
}

func TestLoginAuthCheck200(t *testing.T) {
	hc, teardown := mock.TestingHttpClient(http.HandlerFunc(AuthHandler))
	defer teardown()
	cli := HttpClient{Client: hc}
	code := isUserAuthCorrect(testConf200.Url, testConf200, cli)
	if code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d for case: %v", http.StatusOK, code, testConf200)
	}
}

func TestLoginAuthCheck403(t *testing.T) {
	hc, teardown := mock.TestingHttpClient(http.HandlerFunc(AuthHandler))
	defer teardown()
	cli := HttpClient{Client: hc}
	code := isUserAuthCorrect(testConf403.Url, testConf403, cli)
	if code != http.StatusForbidden {
		t.Errorf("Expected status code %d but got %d for case: %v", http.StatusForbidden, code, testConf403)
	}
}

func TestLoginAuthCheck401(t *testing.T) {
	hc, teardown := mock.TestingHttpClient(http.HandlerFunc(AuthHandler))
	defer teardown()
	cli := HttpClient{Client: hc}
	code := isUserAuthCorrect(testConf401.Url, testConf401, cli)
	if code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d for case: %v", http.StatusUnauthorized, code, testConf401)
	}
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	b := "Basic "
	if auth == b+testConf200.B64 {
		w.WriteHeader(http.StatusOK)
	} else if auth == b+testConf403.B64 {
		w.WriteHeader(http.StatusForbidden)
	} else if auth == b+testConf401.B64 {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getAuthHeader(u, p string) string {
	return "Basic " + EncodeCreds(u, p)
}
