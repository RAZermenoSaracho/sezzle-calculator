package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, contentType, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		wantStatus  int
		wantBody    string
	}{
		{"success", "application/json", `{"expression":"1 + 2 * 3"}`, 200, `{"result":7}`},
		{"success with charset", "application/json; charset=utf-8", `{"expression":"sqrt(16) + 50%"}`, 200, `{"result":4.5}`},
		{"decimal result", "application/json", `{"expression":"1 / 4"}`, 200, `{"result":0.25}`},
		{"empty expression", "application/json", `{"expression":"  "}`, 400, `{"error":"empty expression"}`},
		{"syntax error", "application/json", `{"expression":"1 +"}`, 400, `{"error":"invalid expression"}`},
		{"division by zero", "application/json", `{"expression":"1 / 0"}`, 400, `{"error":"division by zero"}`},
		{"invalid sqrt", "application/json", `{"expression":"sqrt(-4)"}`, 400, `{"error":"invalid square root"}`},
		{"invalid number", "application/json", `{"expression":"10 ^ 1000"}`, 400, `{"error":"invalid number"}`},
		{"too complex", "application/json", `{"expression":"` + strings.Repeat("(", 200) + `"}`, 400, `{"error":"expression too complex"}`},
		{"malformed json", "application/json", `{"expression":`, 400, `{"error":"invalid request"}`},
		{"empty body", "application/json", ``, 400, `{"error":"invalid request"}`},
		{"missing field", "application/json", `{}`, 400, `{"error":"invalid request"}`},
		{"null field", "application/json", `{"expression":null}`, 400, `{"error":"invalid request"}`},
		{"non-string field", "application/json", `{"expression":42}`, 400, `{"error":"invalid request"}`},
		{"unknown field", "application/json", `{"expression":"1","extra":true}`, 400, `{"error":"invalid request"}`},
		{"trailing data", "application/json", `{"expression":"1"} {"expression":"2"}`, 400, `{"error":"invalid request"}`},
		{"not an object", "application/json", `"1 + 2"`, 400, `{"error":"invalid request"}`},
		{"oversized body", "application/json", `{"expression":"` + strings.Repeat("1", 5000) + `"}`, 413, `{"error":"request too large"}`},
		{"wrong content type", "text/plain", `{"expression":"1"}`, 415, `{"error":"unsupported media type"}`},
		{"missing content type", "", `{"expression":"1"}`, 415, `{"error":"unsupported media type"}`},
	}
	h := NewHandler()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, h, http.MethodPost, "/api/calculate", tc.contentType, tc.body)
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d (body %s)", rec.Code, tc.wantStatus, rec.Body)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != tc.wantBody {
				t.Errorf("body = %s, want %s", got, tc.wantBody)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	tests := []struct {
		method, path, allow string
	}{
		{http.MethodGet, "/api/calculate", "POST"},
		{http.MethodPut, "/api/calculate", "POST"},
		{http.MethodPost, "/api/health", "GET"},
	}
	h := NewHandler()
	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			rec := do(t, h, tc.method, tc.path, "application/json", "")
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want 405", rec.Code)
			}
			if got := rec.Header().Get("Allow"); got != tc.allow {
				t.Errorf("Allow = %q, want %q", got, tc.allow)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != `{"error":"method not allowed"}` {
				t.Errorf("body = %s", got)
			}
		})
	}
}

func TestHealth(t *testing.T) {
	rec := do(t, NewHandler(), http.MethodGet, "/api/health", "", "")
	if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != `{"status":"ok"}` {
		t.Errorf("got %d %s", rec.Code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
}

func TestUnexpectedFailuresDoNotLeak(t *testing.T) {
	const secret = "secret internal detail"
	tests := []struct {
		name     string
		evaluate func(string) (float64, error)
	}{
		{"unknown error", func(string) (float64, error) { return 0, errors.New(secret) }},
		{"panic", func(string) (float64, error) { panic(secret) }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, newHandler(tc.evaluate), http.MethodPost, "/api/calculate", "application/json", `{"expression":"1"}`)
			if rec.Code != http.StatusInternalServerError {
				t.Errorf("status = %d, want 500", rec.Code)
			}
			if got := strings.TrimSpace(rec.Body.String()); got != `{"error":"internal error"}` {
				t.Errorf("body = %s", got)
			}
			if strings.Contains(rec.Body.String(), secret) {
				t.Error("response leaks internal detail")
			}
		})
	}
}
