// Package http arma el router, los middlewares comunes y los handlers.
package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"saas/backend/internal/http/api"
	"saas/backend/internal/platform/clientip"
)

const apiV1Prefix = "/api/v1"

// Check verifica una dependencia (BD, Redis…) para /readyz.
type Check func(ctx context.Context) error

// Deps agrupa las dependencias necesarias para construir el router HTTP.
type Deps struct {
	Logger         *slog.Logger
	AllowedOrigins []string
	ReadyChecks    map[string]Check
	AppVersion     string
	ClientIP       *clientip.Resolver
	MountRealtime  func(chi.Router)
	MountWebhooks  func(chi.Router)
}

// NewRouter configura los middlewares globales, los endpoints OpenAPI (/healthz,
// /readyz, /api/v1/*) y los grupos aislados para tiempo real (/rt) y webhooks (/webhooks).
func NewRouter(d Deps) http.Handler {
	log := d.Logger
	if log == nil {
		log = slog.Default()
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(requestIDHeader)
	r.Use(requestLogger(log, d.ClientIP))
	r.Use(middleware.Recoverer)
	r.Use(securityHeaders)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   d.AllowedOrigins, // orígenes explícitos (BR-137)
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-Id", "Idempotency-Key"},
		ExposedHeaders:   []string{"X-Request-Id", "X-API-Version"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	strictHandler := api.NewStrictHandlerWithOptions(NewServer(d), nil, api.StrictHTTPServerOptions{
		RequestErrorHandlerFunc:  requestErrorHandler,
		ResponseErrorHandlerFunc: newResponseErrorHandler(log),
	})
	api.HandlerWithOptions(strictHandler, api.ChiServerOptions{
		BaseRouter: r,
		Middlewares: []api.MiddlewareFunc{
			onlyPathPrefix(apiV1Prefix, middleware.Timeout(30*time.Second)),
			onlyPathPrefix(apiV1Prefix, apiVersionHeader),
		},
		ErrorHandlerFunc: requestErrorHandler,
	})

	// /rt queda fuera del timeout de 30s de /api/v1 para soportar SSE/WebSocket (sección 5.6).
	r.Route("/rt", func(r chi.Router) {
		if d.MountRealtime != nil {
			d.MountRealtime(r)
		}
	})

	// /webhooks queda aislado de los middlewares específicos de /api/v1 (sección 5.2).
	r.Route("/webhooks", func(r chi.Router) {
		if d.MountWebhooks != nil {
			d.MountWebhooks(r)
		}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		setAPIVersionHeaderIfV1(w, r)
		writeError(w, r, http.StatusNotFound, "not_found", "Recurso no encontrado.")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		setAPIVersionHeaderIfV1(w, r)
		writeError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Método no permitido.")
	})

	return r
}

func onlyPathPrefix(prefix string, mw func(http.Handler) http.Handler) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		wrapped := mw(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if hasPathPrefix(r.URL.Path, prefix) {
				wrapped.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func hasPathPrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func apiVersionHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", apiVersionV1)
		next.ServeHTTP(w, r)
	})
}

func setAPIVersionHeaderIfV1(w http.ResponseWriter, r *http.Request) {
	if r != nil && hasPathPrefix(r.URL.Path, apiV1Prefix) {
		w.Header().Set("X-API-Version", apiVersionV1)
	}
}

func requestIDHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if reqID := middleware.GetReqID(r.Context()); reqID != "" {
			w.Header().Set("X-Request-Id", reqID)
		}
		next.ServeHTTP(w, r)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

func requestLogger(log *slog.Logger, ipResolver *clientip.Resolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("http",
				"method", r.Method,
				"path", r.URL.Path, // sin query string: puede contener datos sensibles
				"status", ww.Status(),
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", middleware.GetReqID(r.Context()),
				"client_ip", ipResolver.Resolve(r),
			)
		})
	}
}
