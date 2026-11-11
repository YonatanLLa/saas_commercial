package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"saas/backend/internal/http/api"
)

// FieldError describe un error de validación asociado a un campo o parámetro (BR-004).
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("escribiendo respuesta JSON", "err", err)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	writeErrorWithDetails(w, r, status, code, message, nil)
}

// writeErrorWithDetails responde usando el esquema estándar Error de api/openapi.yaml (BR-004).
func writeErrorWithDetails(w http.ResponseWriter, r *http.Request, status int, code, message string, details []FieldError) {
	var body api.Error
	body.Error.Code = code
	body.Error.Message = message

	if r != nil {
		if reqID := middleware.GetReqID(r.Context()); reqID != "" {
			body.Error.RequestId = &reqID
		}
	}

	if len(details) > 0 {
		items := make([]struct {
			Field   string `json:"field"`
			Message string `json:"message"`
		}, len(details))
		for i, d := range details {
			items[i].Field = d.Field
			items[i].Message = d.Message
		}
		body.Error.Details = &items
	}

	writeJSON(w, status, body)
}

// requestErrorHandler sustituye http.Error en oapi-codegen para errores de entrada (400).
func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if details := extractParamFieldError(err); len(details) > 0 {
		writeErrorWithDetails(w, r, http.StatusBadRequest, "validation_error", "Los datos enviados no son válidos.", details)
		return
	}
	writeError(w, r, http.StatusBadRequest, "bad_request", "La solicitud no es válida.")
}

// responseErrorHandler sustituye http.Error en oapi-codegen para fallos internos (500).
func responseErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	writeInternalError(slog.Default(), w, r, err)
}

func newResponseErrorHandler(log *slog.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	if log == nil {
		return responseErrorHandler
	}
	return func(w http.ResponseWriter, r *http.Request, err error) {
		writeInternalError(log, w, r, err)
	}
}

func writeInternalError(log *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var method, path, reqID string
	if r != nil {
		method = r.Method
		path = r.URL.Path
		reqID = middleware.GetReqID(r.Context())
	}
	log.Error("error procesando respuesta HTTP",
		"err", err,
		"method", method,
		"path", path,
		"request_id", reqID,
	)
	writeError(w, r, http.StatusInternalServerError, "internal_error", "Error interno del servidor.")
}

func extractParamFieldError(err error) []FieldError {
	if err == nil {
		return nil
	}
	var (
		unescaped  *api.UnescapedCookieParamError
		unmarshal  *api.UnmarshalingParamError
		reqParam   *api.RequiredParamError
		reqHeader  *api.RequiredHeaderError
		badFormat  *api.InvalidParamFormatError
		tooManyVal *api.TooManyValuesForParamError
	)
	var field string
	switch {
	case errors.As(err, &unescaped):
		field = unescaped.ParamName
	case errors.As(err, &unmarshal):
		field = unmarshal.ParamName
	case errors.As(err, &reqParam):
		field = reqParam.ParamName
	case errors.As(err, &reqHeader):
		field = reqHeader.ParamName
	case errors.As(err, &badFormat):
		field = badFormat.ParamName
	case errors.As(err, &tooManyVal):
		field = tooManyVal.ParamName
	default:
		return nil
	}
	return []FieldError{{
		Field:   field,
		Message: err.Error(),
	}}
}
