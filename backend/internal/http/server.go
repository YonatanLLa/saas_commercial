package http

import (
	"context"
	"strings"
	"time"

	"saas/backend/internal/http/api"
)

const (
	defaultServiceName = "saas-api"
	defaultAppVersion  = "0.1.0"
	apiVersionV1       = "v1"
)

// Server compone los handlers de cada módulo para implementar api.StrictServerInterface.
// Nuevos módulos de dominio pueden incrustarse aquí sin reescribir el enrutador base.
type Server struct {
	SystemHandler
}

var _ api.StrictServerInterface = (*Server)(nil)

// NewServer construye el servidor HTTP compuesto a partir de las dependencias del router.
func NewServer(d Deps) *Server {
	version := strings.TrimSpace(d.AppVersion)
	if version == "" {
		version = defaultAppVersion
	}
	return &Server{
		SystemHandler: SystemHandler{
			service:     defaultServiceName,
			appVersion:  version,
			readyChecks: d.ReadyChecks,
		},
	}
}

// SystemHandler implementa los endpoints transversales de salud y versión del servicio.
type SystemHandler struct {
	service     string
	appVersion  string
	readyChecks map[string]Check
}

// GetHealth responde si el proceso HTTP está vivo (GET /healthz).
func (h SystemHandler) GetHealth(_ context.Context, _ api.GetHealthRequestObject) (api.GetHealthResponseObject, error) {
	return api.GetHealth200JSONResponse{
		Status: api.HealthStatusStatusOk,
	}, nil
}

// GetReadiness ejecuta los checks de dependencias con timeout de 2 segundos (GET /readyz).
func (h SystemHandler) GetReadiness(ctx context.Context, _ api.GetReadinessRequestObject) (api.GetReadinessResponseObject, error) {
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	status := api.HealthStatusStatusOk
	results := make(map[string]string, len(h.readyChecks))
	for name, check := range h.readyChecks {
		if check == nil {
			continue
		}
		if err := check(checkCtx); err != nil {
			results[name] = "error"
			status = api.HealthStatusStatusDegraded
			continue
		}
		results[name] = "ok"
	}

	body := api.HealthStatus{
		Status: status,
		Checks: &results,
	}
	if status == api.HealthStatusStatusDegraded {
		return api.GetReadiness503JSONResponse(body), nil
	}
	return api.GetReadiness200JSONResponse(body), nil
}

// GetSystemInfo devuelve la información básica y versión de la API REST v1 (GET /api/v1/system/info).
func (h SystemHandler) GetSystemInfo(_ context.Context, _ api.GetSystemInfoRequestObject) (api.GetSystemInfoResponseObject, error) {
	return api.GetSystemInfo200JSONResponse{
		Service:    h.service,
		Version:    h.appVersion,
		ApiVersion: apiVersionV1,
		Status:     api.ApiInfoStatusOk,
	}, nil
}
