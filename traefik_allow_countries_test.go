package traefik_allow_countries_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	AllowCountries "github.com/JonasSchubert/traefik-allow-countries"
)

const (
	xForwardedFor = "X-Forwarded-For"
	AT            = "2.16.16.0"
	DE            = "2.56.20.0"
	GB            = "1.186.0.0"
	US            = "2.56.8.0"
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

	cfg.CidrFileFolder = ".test-data"
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
	cfg.CidrFileFolder = ".test-data"
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

func TestAllowedCountry_WithCustomFileExtension(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.LogLocalRequests = false
	cfg.CidrFileFolder = ".test-data"
	cfg.FileExtension = "netset"
	cfg.Countries = append(cfg.Countries, "AT")

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

	req.Header.Add(xForwardedFor, AT)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func TestMultipleAllowedCountries(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.LogLocalRequests = false
	cfg.CidrFileFolder = ".test-data"
	cfg.Countries = append(cfg.Countries, "DE", "GB")

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

	req.Header.Add(xForwardedFor, GB)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusOK)
}

func TestAllowLocalIP(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = true
	cfg.LogLocalRequests = false
	cfg.CidrFileFolder = ".test-data"
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

func TestBlockedCountry(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.LogLocalRequests = false
	cfg.CidrFileFolder = ".test-data"
	cfg.Countries = append(cfg.Countries, "DE", "GB")

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

	req.Header.Add(xForwardedFor, US)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusForbidden)
}

func TestBlockedCountry_WithCustomFileExtension(t *testing.T) {
	cfg := AllowCountries.CreateConfig()

	cfg.AllowLocalRequests = false
	cfg.LogLocalRequests = false
	cfg.CidrFileFolder = ".test-data"
	cfg.FileExtension = "netset"
	cfg.Countries = append(cfg.Countries, "AT")

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

	req.Header.Add(xForwardedFor, US)

	handler.ServeHTTP(recorder, req)

	assertStatusCode(t, recorder.Result(), http.StatusForbidden)
}
