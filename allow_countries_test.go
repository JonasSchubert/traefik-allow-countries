package AllowCountries_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	AllowCountries "github.com/jonasschubert/traefik-allow-countries"
)

const (
	xForwardedFor = "X-Forwarded-For"
	DE            = "99.220.109.148"
	UK            = "82.220.110.18"
	PrivateRange  = "192.168.1.1"
	Invalid       = "192.168.1.X"
)

func TestEmptyCidrFileFolder(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.CidrFileFolder = ""
	cfg.Countries = append(cfg.Countries, "DE")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	_, err := AllowCountries.New(ctx, next, cfg, "AllowCountries")

	// expect error
	if err == nil {
		t.Fatal("Empty CIDR file folder is not allowed")
	}
}

func TestEmptyAllowedCountryList(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.CidrFileFolder = "/usr/share/traefik/cidr"
	cfg.Countries = make([]string, 0)

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	_, err := AllowCountries.New(ctx, next, cfg, "AllowCountries")

	// expect error
	if err == nil {
		t.Fatal("Empty list of allowed countries is not allowed")
	}
}

func TestAllowedCountry(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.LogLocalRequests = false
	cfg.Countries = append(cfg.Countries, "DE")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := AllowCountries.New(ctx, next, cfg, "AllowCountries")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, DE)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func TestMultipleAllowedCountries(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.LogLocalRequests = false
	cfg.Countries = append(cfg.Countries, "DE", "UK")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := AllowCountries.New(ctx, next, cfg, "AllowCountries")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, UK)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func TestAllowLocalIP(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = true
	cfg.LogLocalRequests = false
	cfg.Countries = append(cfg.Countries, "DE")

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := AllowCountries.New(ctx, next, cfg, "AllowCountries")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add(xForwardedFor, PrivateRange)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func assertStatusCode(t *testing.T, req *http.Response, expected int) {
	t.Helper()

	received := req.StatusCode

	if received != expected {
		t.Errorf("invalid status code: %d <> %d", expected, received)
	}
}
