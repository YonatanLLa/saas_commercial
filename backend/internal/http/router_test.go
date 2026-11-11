package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"saas/backend/internal/http/api"
	"saas/backend/internal/platform/clientip"
)

func TestRouter_Healthz(t *testing.T) {
	t.Parallel()

	router := NewRouter(Deps{})

	t.Run("returns ok with security headers and X-Request-Id", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			t.Fatalf("Content-Type = %q, want application/json", ct)
		}
		if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
		}
		if got := rec.Header().Get("X-Frame-Options"); got != "DENY" {
			t.Fatalf("X-Frame-Options = %q, want DENY", got)
		}
		if got := rec.Header().Get("Referrer-Policy"); got != "strict-origin-when-cross-origin" {
			t.Fatalf("Referrer-Policy = %q, want strict-origin-when-cross-origin", got)
		}
		if got := rec.Header().Get("Content-Security-Policy"); got != "default-src 'none'; frame-ancestors 'none'" {
			t.Fatalf("Content-Security-Policy = %q", got)
		}
		if got := rec.Header().Get("X-Request-Id"); got == "" {
			t.Fatal("expected non-empty X-Request-Id response header")
		}
		if got := rec.Header().Get("X-API-Version"); got != "" {
			t.Fatalf("unexpected X-API-Version on /healthz: %q", got)
		}

		var body api.HealthStatus
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Status != api.HealthStatusStatusOk {
			t.Fatalf("body.Status = %q, want ok", body.Status)
		}
	})

	t.Run("preserves incoming X-Request-Id", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		req.Header.Set("X-Request-Id", "req-trace-001")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if got := rec.Header().Get("X-Request-Id"); got != "req-trace-001" {
			t.Fatalf("X-Request-Id = %q, want req-trace-001", got)
		}
	})
}

func TestRouter_Readyz(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 when all checks pass", func(t *testing.T) {
		t.Parallel()
		var observedDeadline bool
		router := NewRouter(Deps{
			ReadyChecks: map[string]Check{
				"database": func(ctx context.Context) error {
					_, observedDeadline = ctx.Deadline()
					return nil
				},
				"redis": func(context.Context) error { return nil },
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if !observedDeadline {
			t.Fatal("expected readiness check context to have a deadline")
		}

		var body api.HealthStatus
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Status != api.HealthStatusStatusOk {
			t.Fatalf("status = %q, want ok", body.Status)
		}
		if body.Checks == nil || (*body.Checks)["database"] != "ok" || (*body.Checks)["redis"] != "ok" {
			t.Fatalf("unexpected checks: %#v", body.Checks)
		}
	})

	t.Run("returns 503 degraded when a check fails", func(t *testing.T) {
		t.Parallel()
		router := NewRouter(Deps{
			ReadyChecks: map[string]Check{
				"database": func(context.Context) error { return nil },
				"redis":    func(context.Context) error { return errors.New("redis down") },
			},
		})

		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want 503", rec.Code)
		}

		var body api.HealthStatus
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Status != api.HealthStatusStatusDegraded {
			t.Fatalf("status = %q, want degraded", body.Status)
		}
		if body.Checks == nil || (*body.Checks)["database"] != "ok" || (*body.Checks)["redis"] != "error" {
			t.Fatalf("unexpected checks: %#v", body.Checks)
		}
	})
}

func TestRouter_SystemInfoV1(t *testing.T) {
	t.Parallel()

	t.Run("uses configured AppVersion", func(t *testing.T) {
		t.Parallel()
		router := NewRouter(Deps{
			AppVersion: "1.4.2",
		})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", ct)
		}
		if ver := rec.Header().Get("X-API-Version"); ver != "v1" {
			t.Fatalf("X-API-Version = %q, want v1", ver)
		}

		var got map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		want := map[string]string{
			"apiVersion": "v1",
			"service":    "saas-api",
			"status":     "ok",
			"version":    "1.4.2",
		}
		for k, wv := range want {
			if got[k] != wv {
				t.Fatalf("field %q = %v, want %q", k, got[k], wv)
			}
		}
	})

	t.Run("defaults AppVersion to 0.1.0 when empty", func(t *testing.T) {
		t.Parallel()
		router := NewRouter(Deps{})

		req := httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}

		var body api.ApiInfo
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Version != "0.1.0" || body.ApiVersion != "v1" || body.Service != "saas-api" || body.Status != api.ApiInfoStatusOk {
			t.Fatalf("unexpected ApiInfo: %+v", body)
		}
	})
}

func TestRouter_NotFoundAndMethodNotAllowed(t *testing.T) {
	t.Parallel()

	router := NewRouter(Deps{})

	t.Run("404 on unknown /api/v1 path returns standardized error and X-API-Version", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/unknown", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ver := rec.Header().Get("X-API-Version"); ver != "v1" {
			t.Fatalf("X-API-Version = %q, want v1", ver)
		}

		var body api.Error
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Error.Code != "not_found" {
			t.Fatalf("error.code = %q, want not_found", body.Error.Code)
		}
		if body.Error.RequestId == nil || *body.Error.RequestId == "" {
			t.Fatal("expected non-empty error.requestId")
		}
		if gotHdr := rec.Header().Get("X-Request-Id"); gotHdr != *body.Error.RequestId {
			t.Fatalf("X-Request-Id header = %q, body.error.requestId = %q", gotHdr, *body.Error.RequestId)
		}
	})

	t.Run("404 on /api/v10 prefix does not set X-API-Version v1", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/api/v10/unknown", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want 404", rec.Code)
		}
		if ver := rec.Header().Get("X-API-Version"); ver != "" {
			t.Fatalf("unexpected X-API-Version on /api/v10: %q", ver)
		}
	})

	t.Run("405 on wrong method on /api/v1/system/info", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/system/info", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want 405", rec.Code)
		}
		if ver := rec.Header().Get("X-API-Version"); ver != "v1" {
			t.Fatalf("X-API-Version = %q, want v1", ver)
		}

		var body api.Error
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Error.Code != "method_not_allowed" {
			t.Fatalf("error.code = %q, want method_not_allowed", body.Error.Code)
		}
		if body.Error.RequestId == nil || *body.Error.RequestId == "" {
			t.Fatal("expected non-empty error.requestId")
		}
	})
}

func TestRouter_IgnoresSpoofedForwardedHeaders(t *testing.T) {
	t.Parallel()

	t.Run("default nil ClientIP ignores X-Forwarded-For and preserves RemoteAddr", func(t *testing.T) {
		t.Parallel()
		var logBuf bytes.Buffer
		log := slog.New(slog.NewJSONHandler(&logBuf, nil))

		var seenRemoteAddr string
		router := NewRouter(Deps{
			Logger: log,
			MountWebhooks: func(r chi.Router) {
				r.Post("/test", func(w http.ResponseWriter, r *http.Request) {
					seenRemoteAddr = r.RemoteAddr
					w.WriteHeader(http.StatusOK)
				})
			},
		})

		req := httptest.NewRequest(http.MethodPost, "/webhooks/test", nil)
		req.RemoteAddr = "198.51.100.77:54321"
		req.Header.Set("X-Forwarded-For", "1.2.3.4")
		req.Header.Set("X-Real-IP", "5.6.7.8")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if seenRemoteAddr != "198.51.100.77:54321" {
			t.Fatalf("RemoteAddr was mutated to %q", seenRemoteAddr)
		}

		var entry map[string]any
		if err := json.Unmarshal(logBuf.Bytes(), &entry); err != nil {
			t.Fatalf("unmarshal log line: %v", err)
		}
		if entry["client_ip"] != "198.51.100.77" {
			t.Fatalf("logged client_ip = %v, want 198.51.100.77", entry["client_ip"])
		}
	})

	t.Run("configured trusted proxy resolves forwarded client_ip", func(t *testing.T) {
		t.Parallel()
		var logBuf bytes.Buffer
		log := slog.New(slog.NewJSONHandler(&logBuf, nil))
		resolver, err := clientip.NewResolver([]string{"10.0.0.0/8"}, "")
		if err != nil {
			t.Fatalf("NewResolver: %v", err)
		}

		router := NewRouter(Deps{
			Logger:   log,
			ClientIP: resolver,
		})

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		req.RemoteAddr = "10.1.2.3:4567"
		req.Header.Set("X-Forwarded-For", "9.9.9.9, 203.0.113.44")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}

		var entry map[string]any
		if err := json.Unmarshal(logBuf.Bytes(), &entry); err != nil {
			t.Fatalf("unmarshal log line: %v", err)
		}
		if entry["client_ip"] != "203.0.113.44" {
			t.Fatalf("logged client_ip = %v, want 203.0.113.44", entry["client_ip"])
		}
	})
}

func TestRouter_RealtimeAndWebhooksIsolation(t *testing.T) {
	t.Parallel()

	var rtHasDeadline, whHasDeadline bool
	router := NewRouter(Deps{
		MountRealtime: func(r chi.Router) {
			r.Get("/events", func(w http.ResponseWriter, r *http.Request) {
				_, rtHasDeadline = r.Context().Deadline()
				w.WriteHeader(http.StatusOK)
			})
		},
		MountWebhooks: func(r chi.Router) {
			r.Post("/whatsapp", func(w http.ResponseWriter, r *http.Request) {
				_, whHasDeadline = r.Context().Deadline()
				w.WriteHeader(http.StatusAccepted)
			})
		},
	})

	t.Run("/rt works without /api/v1 timeout or X-API-Version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/rt/events", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if rtHasDeadline {
			t.Fatal("/rt route should not inherit /api/v1 30s timeout")
		}
		if ver := rec.Header().Get("X-API-Version"); ver != "" {
			t.Fatalf("unexpected X-API-Version on /rt: %q", ver)
		}
	})

	t.Run("/webhooks works without /api/v1 timeout or X-API-Version", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhooks/whatsapp", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Fatalf("status = %d, want 202", rec.Code)
		}
		if whHasDeadline {
			t.Fatal("/webhooks route should not inherit /api/v1 30s timeout")
		}
		if ver := rec.Header().Get("X-API-Version"); ver != "" {
			t.Fatalf("unexpected X-API-Version on /webhooks: %q", ver)
		}
	})

	t.Run("onlyPathPrefix applies middleware only to /api/v1 routes", func(t *testing.T) {
		t.Parallel()
		var v1HasDeadline, rootHasDeadline bool
		wrapped := onlyPathPrefix(apiV1Prefix, middleware.Timeout(30*time.Second))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hasPathPrefix(r.URL.Path, apiV1Prefix) {
				_, v1HasDeadline = r.Context().Deadline()
			} else {
				_, rootHasDeadline = r.Context().Deadline()
			}
			w.WriteHeader(http.StatusOK)
		}))

		wrapped.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil))
		wrapped.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/healthz", nil))

		if !v1HasDeadline {
			t.Fatal("expected /api/v1 route to have timeout deadline")
		}
		if rootHasDeadline {
			t.Fatal("expected non-/api/v1 route not to have 30s timeout deadline")
		}
	})
}

func TestRouter_CORS(t *testing.T) {
	t.Parallel()

	router := NewRouter(Deps{
		AllowedOrigins: []string{"https://app.example.com"},
	})

	t.Run("allows configured origin preflight with Idempotency-Key", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/system/info", nil)
		req.Header.Set("Origin", "https://app.example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Authorization, Idempotency-Key")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent && rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200 or 204", rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
			t.Fatalf("Access-Control-Allow-Origin = %q, want https://app.example.com", got)
		}
		if got := rec.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
			t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
		}
		allowHdrs := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
		if !strings.Contains(allowHdrs, "idempotency-key") {
			t.Fatalf("Access-Control-Allow-Headers = %q, expected Idempotency-Key", allowHdrs)
		}
	})

	t.Run("rejects disallowed origin", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		req.Header.Set("Origin", "https://evil.example.com")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("unexpected Access-Control-Allow-Origin = %q for disallowed origin", got)
		}
	})
}

func TestRespond_ErrorHandlers(t *testing.T) {
	t.Parallel()

	t.Run("requestErrorHandler formats param validation error as JSON", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "req-400")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		requestErrorHandler(rec, req, &api.RequiredParamError{ParamName: "tenantId"})

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rec.Code)
		}
		var body api.Error
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Error.Code != "validation_error" {
			t.Fatalf("code = %q, want validation_error", body.Error.Code)
		}
		if body.Error.RequestId == nil || *body.Error.RequestId != "req-400" {
			t.Fatalf("requestId = %v, want req-400", body.Error.RequestId)
		}
		if body.Error.Details == nil || len(*body.Error.Details) != 1 || (*body.Error.Details)[0].Field != "tenantId" {
			t.Fatalf("unexpected details: %#v", body.Error.Details)
		}
	})

	t.Run("responseErrorHandler hides internal error message and returns 500 JSON", func(t *testing.T) {
		t.Parallel()
		var logBuf bytes.Buffer
		log := slog.New(slog.NewJSONHandler(&logBuf, nil))
		handler := newResponseErrorHandler(log)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/system/info", nil)
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "req-500")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()

		handler(rec, req, errors.New("secret sql detail"))

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", rec.Code)
		}
		if strings.Contains(rec.Body.String(), "secret sql detail") {
			t.Fatal("response leaked internal error string")
		}
		if !strings.Contains(logBuf.String(), "secret sql detail") {
			t.Fatal("expected internal error detail to be logged")
		}
		var body api.Error
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if body.Error.Code != "internal_error" {
			t.Fatalf("code = %q, want internal_error", body.Error.Code)
		}
		if body.Error.RequestId == nil || *body.Error.RequestId != "req-500" {
			t.Fatalf("requestId = %v, want req-500", body.Error.RequestId)
		}
	})
}
