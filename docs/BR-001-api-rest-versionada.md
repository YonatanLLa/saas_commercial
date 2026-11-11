# BR-001: API REST versionada (`/api/v1`) y fundamentos del router HTTP

## 1. Requisito atendido y alineación arquitectónica

- **Requisito funcional/técnico:** `BR-001` (*"Implementar una API REST versionada, por ejemplo `/api/v1`"*), junto con soporte transversal para `BR-004` (formato estándar de errores), `BR-005` (contrato OpenAPI) y `BR-169` (resolución segura de IP del cliente).
- **Alineación con `arquitectura_backend.md`:**
  - **Sección 5.5 (`platform/clientip`):** se eliminó `middleware.RealIP` de `chi` y se implementó `internal/platform/clientip` para obtener la IP del cliente únicamente cuando el peer TCP inmediato pertenece a un proxy/balanceador de confianza configurado.
  - **Sección 5.6 y Sección 9 (`internal/http/router.go`):** se separaron los grupos de rutas `/api/v1`, `/rt` y `/webhooks`, evitando que el timeout de 30 s de `/api/v1` afecte a tiempo real (`/rt`) o webhooks (`/webhooks`).
  - **Sección 6 (flujo spec-first):** se actualizó primero `api/openapi.yaml` y se regeneró `backend/internal/http/api/api.gen.go` mediante `oapi-codegen` (`chi-server`, `models`, `strict-server`).

---

## 2. Archivos creados y modificados

| Archivo | Estado | Descripción |
|---|---|---|
| `api/openapi.yaml` | Modificado | Se mantiene `openapi: 3.1.0`, `info.version: 0.1.0`, `/healthz`, `/readyz` y los esquemas base (`HealthStatus`, `Error`, `PageMeta`). Se añade `GET /api/v1/system/info` (`operationId: getSystemInfo`, `security: []`) referenciando `#/components/schemas/ApiInfo` (`200`) y `#/components/schemas/Error` (`500`). |
| `backend/go.mod` | Modificado | Se eliminó `github.com/oapi-codegen/runtime` del bloque `tool (...)`, manteniendo `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen` como herramienta y `github.com/oapi-codegen/runtime v1.9.0` en `require`. |
| `backend/internal/http/api/api.gen.go` | Regenerado | Regenerado con `go tool oapi-codegen -config oapi-codegen.yaml ../api/openapi.yaml`. Incluye `ApiInfo`, `ApiInfoStatusOk`, `Error`, `HealthStatus`, `ServerInterface` y `StrictServerInterface` con `GetSystemInfo`. |
| `backend/internal/platform/clientip/clientip.go` | Creado | Implementa `Resolver` y `NewResolver(trustedCIDRs, header)` usando `net/netip`. Resuelve `X-Forwarded-For` de derecha a izquierda (o cabecera configurada como `X-Real-IP`) solo si `r.RemoteAddr` pertenece a un CIDR confiable; en caso contrario retorna la IP directa de la conexión. |
| `backend/internal/platform/clientip/clientip_test.go` | Creado | Pruebas unitarias para ausencia de proxies confiables (ignora cabeceras falsificadas), proxies confiables por CIDR (incluyendo múltiples saltos y cabeceras personalizadas), proxies no confiables, cabeceras malformadas e IPv6 / IPv4-mapeado. |
| `backend/internal/http/server.go` | Creado | Define `Server` (composición modular de handlers) y `SystemHandler` implementando `api.StrictServerInterface`: `GetHealth`, `GetReadiness` (checks con timeout de 2 s) y `GetSystemInfo`. |
| `backend/internal/http/respond.go` | Modificado | Alinea las respuestas de error con el modelo generado `api.Error` (incluyendo `Details` con `FieldError` y `RequestId`). Añade `writeErrorWithDetails`, `requestErrorHandler`, `responseErrorHandler` y `newResponseErrorHandler` para evitar respuestas en texto plano de `http.Error`. |
| `backend/internal/http/router.go` | Modificado | Actualiza `Deps` manteniendo compatibilidad con `cmd/api/main.go` (`AppVersion`, `ClientIP`, `MountRealtime`, `MountWebhooks`). Quita `middleware.RealIP`, añade `requestIDHeader` para devolver `X-Request-Id` en la cabecera de respuesta, registra `client_ip` vía `clientip.Resolver`, agrega `Idempotency-Key` a CORS `AllowedHeaders`, monta OpenAPI con `onlyPathPrefix("/api/v1", ...)` (`middleware.Timeout(30 * time.Second)` y `X-API-Version: v1`), aísla `/rt` y `/webhooks`, y conserva `NotFound`/`MethodNotAllowed` en JSON con `X-API-Version: v1` bajo `/api/v1`. |
| `backend/internal/http/router_test.go` | Creado | Pruebas `httptest` para `/healthz`, `/readyz` (200/503), `/api/v1/system/info`, `404` y `405` estandarizados, prevención de spoofing de `X-Forwarded-For`, aislamiento de `/rt` y `/webhooks`, CORS y manejadores de error. |
| `docs/BR-001-api-rest-versionada.md` | Creado | Este documento técnico. |

---

## 3. Decisiones técnicas

1. **Contrato spec-first en `api/openapi.yaml`:**
   - El endpoint `GET /api/v1/system/info` inaugura el prefijo versionado `/api/v1` en el contrato OpenAPI.
   - Al referenciar `#/components/schemas/Error` en la respuesta `"500"` de `getSystemInfo`, `oapi-codegen` genera el tipo `api.Error` dentro de `api.gen.go` en lugar de podarlo por no estar referenciado.
2. **Conservación de `/healthz` y `/readyz` en la raíz `/`:**
   - Se mantienen en `/healthz` y `/readyz` (sin prefijo `/api/v1`) para preservar compatibilidad directa con el frontend existente (`web/src/modules/dashboard/pages/DashboardPage.vue` consume `api.GET('/readyz')`) y con sondas de infraestructura estándar.
3. **Separación de grupos `/api/v1`, `/rt` y `/webhooks`:**
   - Las rutas generadas por OpenAPI se montan mediante `api.HandlerWithOptions(api.NewStrictHandlerWithOptions(...))`, aplicando `middleware.Timeout(30 * time.Second)` y el encabezado `X-API-Version: v1` **únicamente** a rutas con prefijo exacto `/api/v1` (`onlyPathPrefix`, validando frontera de segmento para que `/api/v10` no coincida).
   - Los subrouters `/rt` (`MountRealtime`) y `/webhooks` (`MountWebhooks`) quedan fuera de ese timeout y de los encabezados/middlewares de `/api/v1`, cumpliendo las secciones 5.2 y 5.6 de `arquitectura_backend.md`.
   - `NotFound` y `MethodNotAllowed` devuelven JSON compatible con `api.Error` (con `requestId`) y añaden `X-API-Version: v1` cuando la ruta pertenece a `/api/v1`.
4. **Eliminación de `middleware.RealIP` y uso de `platform/clientip`:**
   - `middleware.RealIP` reescribía `r.RemoteAddr` confiando ciegamente en `X-Forwarded-For`. Se eliminó por completo para cumplir `BR-169`.
   - `requestLogger` calcula `client_ip` mediante `clientip.Resolver` (por defecto `&clientip.Resolver{}`, que usa directamente `r.RemoteAddr` sin confiar en cabeceras de reenvío cuando no hay proxies confiables configurados).
   - Cuando existe un proxy confiable en `X-Forwarded-For`, `Resolver` recorre la lista de derecha a izquierda para evitar la falsificación en el extremo izquierdo de la cabecera.
5. **Propagación de `X-Request-Id` en respuestas:**
   - Se incorporó `requestIDHeader` inmediatamente después de `middleware.RequestID` para reflejar el identificador de correlación en la cabecera HTTP de respuesta `X-Request-Id` (en línea con `ExposedHeaders` de CORS).
6. **Composición modular en `Server` (`internal/http/server.go`):**
   - `Server` implementa `api.StrictServerInterface` componiendo `SystemHandler`. Cuando se añadan nuevos módulos (`auth`, `tenants`, `inbox`, etc.), sus handlers se incrustarán en `Server` dentro de `server.go` sin reescribir `router.go`.

---

## 4. Cómo probarlo y resultados observados

Todos los comandos se ejecutaron desde `backend/` tras los cambios finales de código:

1. **Regeneración de código Go desde OpenAPI:**
   ```powershell
   go tool oapi-codegen -config oapi-codegen.yaml ../api/openapi.yaml
   ```
   Resultado: terminó con código de salida `0` y actualizó `backend/internal/http/api/api.gen.go`.

2. **Verificación de dependencias locales:**
   ```powershell
   go mod verify
   ```
   Resultado: `all modules verified` (código de salida `0`).
   *(Nota: `go mod tidy` intentó consultar `proxy.golang.org` para dependencias de prueba transitivas de herramientas en `tool` y falló sin red en el entorno aislado; `go.mod` ya conserva `github.com/oapi-codegen/runtime v1.9.0` en `require` y verificado en `go.sum`.)*

3. **Pruebas unitarias y HTTP (sin caché):**
   ```powershell
   go test -count=1 ./...
   ```
   Resultado (código de salida `0`):
   - `ok  	saas/backend/internal/http	1.087s`
   - `ok  	saas/backend/internal/platform/clientip	0.580s`
   - El resto de paquetes (`cmd/api`, `cmd/worker`, `internal/http/api`, `internal/jobs`, `internal/platform/cache`, `internal/platform/config`, `internal/platform/database`, `internal/platform/logger`) compilaron sin archivos de prueba.

   Ejecución detallada de ambos paquetes con pruebas (`go test -count=1 -v ./internal/platform/clientip` y `go test -count=1 -v ./internal/http`):
   - `TestResolver_NoTrustedProxiesIgnoresSpoofedHeaders` (3 subcasos): PASS
   - `TestResolver_TrustedProxyCIDR` (5 subcasos): PASS
   - `TestResolver_UntrustedProxyFallsBackToRemoteAddr`: PASS
   - `TestResolver_MalformedHeadersAndIPv6` (6 subcasos): PASS
   - `TestRouter_Healthz` (2 subcasos): PASS
   - `TestRouter_Readyz` (2 subcasos): PASS
   - `TestRouter_SystemInfoV1` (2 subcasos): PASS
   - `TestRouter_NotFoundAndMethodNotAllowed` (3 subcasos): PASS
   - `TestRouter_IgnoresSpoofedForwardedHeaders` (2 subcasos): PASS
   - `TestRouter_RealtimeAndWebhooksIsolation` (3 subcasos): PASS
   - `TestRouter_CORS` (2 subcasos): PASS
   - `TestRespond_ErrorHandlers` (2 subcasos): PASS

4. **Análisis estático y compilación:**
   ```powershell
   go vet ./...
   go build ./...
   ```
   Resultado: ambos finalizaron con código de salida `0` sin advertencias ni errores.

---

## 5. Nota para sincronización con el frontend (`web/`)

Cuando se trabaje en el frontend y se requiera reflejar `GET /api/v1/system/info` y el esquema `ApiInfo` en los tipos TypeScript (`web/src/shared/api/schema.d.ts`), se debe ejecutar desde la raíz del proyecto:

```bash
cd web && pnpm install && pnpm api:types
```
