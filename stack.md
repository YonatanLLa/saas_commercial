# Stack técnico

## SaaS multiindustria de automatización comercial por WhatsApp

**Versión:** 1.0\
**Fecha:** 09 de noviembre de 2026\
**Relacionado con:** `requerimientos_frontend.md`, `requerimientos_backend.md`

Las versiones indicadas son las mínimas recomendadas. Al iniciar el
proyecto se usará la última versión estable de cada una.

**Criterio de licencias:** todas las dependencias centrales deben ser
open source con licencia permisiva (MIT, Apache 2.0, BSD o similar), sin
versiones "PRO" o comerciales obligatorias. Por este motivo se descartó
PrimeVue 5, que pasó a la licencia PrimeUI.

------------------------------------------------------------------------

# 1. Resumen

| Capa | Tecnología |
|---|---|
| Frontend | Vue 3 + Vite + TypeScript |
| UI | shadcn-vue (Reka UI) + Tailwind CSS 4 |
| Backend | Go 1.25+ (API y worker como binarios separados) |
| Contrato API | OpenAPI 3.1 *spec-first* (`oapi-codegen` + `openapi-typescript`) |
| Base de datos | PostgreSQL 17 con `tenant_id` + Row Level Security |
| Acceso a datos | `pgx` + `sqlc` |
| Colas | River (cola transaccional sobre PostgreSQL) |
| Pub/sub, caché, rate limit | Redis (o Valkey) |
| Tiempo real | WebSocket propio con `coder/websocket` + Redis pub/sub |
| Archivos | Almacenamiento compatible con S3 (AWS S3 / Cloudflare R2; SeaweedFS en local) |
| IA | Interfaz propia de proveedores; Anthropic (Claude) como proveedor inicial |
| Observabilidad | `slog` + OpenTelemetry + Sentry |
| Infraestructura | Docker; CI con GitHub Actions |

------------------------------------------------------------------------

# 2. Estructura del repositorio

```
saas/
├── api/
│   └── openapi.yaml          # Contrato único entre frontend y backend
├── backend/                  # Módulo Go
│   ├── cmd/
│   │   ├── api/              # Binario: API REST /api/v1 + WebSocket
│   │   └── worker/           # Binario: trabajos River (webhooks, envíos,
│   │                         #   automatizaciones, recordatorios, IA)
│   ├── internal/
│   │   ├── platform/         # Config, BD, Redis, logs, auth, tenant, RLS
│   │   ├── http/             # Handlers generados + middlewares
│   │   ├── realtime/         # Hub WebSocket
│   │   ├── integrations/     # whatsapp, ai, storage, email
│   │   └── modules/          # Dominio: auth, tenants, users, contacts,
│   │                         #   conversations, messages, automations,
│   │                         #   catalog, inventory, orders, retail,
│   │                         #   restaurant, clinic, quotations, reports,
│   │                         #   audit, plans
│   ├── db/
│   │   ├── migrations/       # goose
│   │   └── queries/          # SQL para sqlc
│   └── go.mod
├── web/                      # Aplicación Vue
├── infra/
│   └── docker/               # docker-compose local y Dockerfiles
└── docs/                     # Requerimientos, ADRs y runbooks
```

`api` y `worker` comparten el mismo código de dominio y se despliegan
como procesos separados para escalar por separado (RNF-002).

------------------------------------------------------------------------

# 3. Frontend (`web/`)

| Necesidad | Elección | Requerimientos |
|---|---|---|
| Framework | Vue 3 (Composition API, `<script setup>`) + TypeScript | FR-001 |
| Build | Vite | RNF-001 |
| Componentes | shadcn-vue (sobre Reka UI); el código vive en `src/shared/ui` | FR-008, RNF-008 |
| Estilos | Tailwind CSS 4 | FR-001 |
| Íconos | Lucide (`lucide-vue-next`) | — |
| Tablas | TanStack Table en un componente `DataTable` compartido (paginación, filtros y orden del lado del servidor) | FR-008, FR-036 |
| Notificaciones | vue-sonner | FR-006 |
| Ruteo | Vue Router con guards por permiso, rubro y módulo | FR-004 |
| Estado de cliente | Pinia (sesión, tenant y sucursal activos, preferencias) | FR-003, FR-005 |
| Estado del servidor | TanStack Query (`@tanstack/vue-query`) | FR-006, RNF-001 |
| Formularios | VeeValidate + Zod | FR-039, FR-046 |
| Cliente HTTP | `openapi-fetch` con tipos generados por `openapi-typescript` | BR-005 |
| Tiempo real | WebSocket nativo con composable propio (reconexión con backoff y resincronización) | FR-031 |
| Editor de automatizaciones | Vue Flow | FR-053 |
| Gráficos | ECharts (`vue-echarts`) | FR-015, FR-017 |
| Fechas y zonas horarias | date-fns + date-fns-tz | RNF-011 |
| Internacionalización | vue-i18n (es-PE por defecto) | RNF-011 |
| Utilidades | VueUse | — |
| WhatsApp Embedded Signup | SDK JavaScript de Meta, cargado solo en la pantalla de integración | FR-070 |
| Pruebas | Vitest + Vue Testing Library; Playwright para E2E | BR-155 |

**Estructura interna:**

```
web/src/
├── app/            # Arranque, router, plugins, layouts
├── modules/        # Un directorio por dominio: auth, inbox, contacts,
│                   #   catalog, automations, users, whatsapp, reports,
│                   #   settings, retail, restaurant, clinic, hardware
├── shared/
│   ├── ui/         # Componentes shadcn-vue
│   ├── components/ # DataTable, estados vacíos/error, filtros
│   ├── composables/# useRealtime, usePermissions, useTenant
│   └── api/        # Cliente generado desde openapi.yaml
└── i18n/
```

Cada módulo vertical se carga de forma diferida y solo se registra si
el tenant lo tiene habilitado (FR-004, FR-104).

**Sesión:** el access token vive solo en memoria. El refresh token va
en una cookie `httpOnly`, `Secure` y `SameSite=Strict`, que JavaScript
no puede leer (FR-007).

------------------------------------------------------------------------

# 4. Backend (`backend/`)

| Necesidad | Elección | Requerimientos |
|---|---|---|
| Lenguaje | Go 1.25+ | — |
| Router HTTP | `chi` (compatible con `net/http`) | BR-001 |
| Contrato | `openapi.yaml` escrito primero; `oapi-codegen` genera interfaces del servidor y tipos | BR-004, BR-005 |
| Validación | Validación del esquema OpenAPI (`kin-openapi`) + `go-playground/validator` para reglas de negocio | BR-003, BR-135 |
| Base de datos | PostgreSQL 17 + `pgx` | — |
| Consultas | `sqlc` (SQL escrito a mano → código Go tipado, siempre parametrizado) | BR-136 |
| Migraciones | `goose` | — |
| Aislamiento multi-tenant | `tenant_id` en todas las tablas + RLS de PostgreSQL; el tenant se fija por transacción con `set_config('app.tenant_id', …)` desde un middleware | BR-009 a BR-013 |
| Autenticación | JWT de acceso de 15 min (`golang-jwt`) + refresh token rotativo guardado como hash en BD | BR-016, BR-018 a BR-020 |
| Contraseñas | Argon2id (`golang.org/x/crypto/argon2`) | BR-017 |
| MFA | TOTP (`pquerna/otp`) para propietarios y administradores | BR-026 |
| Autorización | RBAC propio por permisos, aplicado con middleware en cada endpoint | BR-021, BR-022 |
| Rate limiting | `httprate` con almacenamiento en Redis | BR-024, BR-140 |
| Seguridad HTTP | CORS con orígenes explícitos (`go-chi/cors`), cabeceras de seguridad y CSRF en endpoints con cookie | BR-137 a BR-139 |
| Colas | River: el trabajo se encola en la **misma transacción** que guarda el dato; reintentos con backoff, trabajos descartados y claves únicas para idempotencia | BR-117 a BR-121 |
| Tareas programadas | Trabajos periódicos de River (recordatorios y vencimiento de cotizaciones) según zona horaria del tenant | BR-122 |
| WhatsApp | Cliente HTTP propio sobre Graph API (Cloud API) en `integrations/whatsapp`, aislado del dominio | BR-050 a BR-060 |
| Webhooks | Validan la firma `X-Hub-Signature-256`, guardan el evento crudo, encolan y responden `200` de inmediato; deduplicación por ID externo con índice único | BR-041 a BR-045 |
| Secretos de integraciones | Cifrado AES-256-GCM a nivel de aplicación, con clave en un gestor de secretos (AWS KMS o equivalente) | BR-052, BR-141 |
| Archivos | S3 con URLs prefirmadas (`aws-sdk-go-v2`); validación de tipo MIME y tamaño | BR-049 |
| IA | Interfaz `AIProvider` propia; adaptador inicial con el SDK oficial de Anthropic para Go. Las acciones propuestas por el modelo pasan por los mismos servicios, permisos y validaciones | BR-071 a BR-080 |
| Correo | Amazon SES o Resend; Mailpit en local | BR-023, BR-028 |
| Logs | `log/slog` en JSON, con redacción de tokens, contraseñas y datos sensibles | BR-142 |
| Trazas y métricas | OpenTelemetry; errores en Sentry | BR-148, RNF-005 |
| Pruebas | `testing` + testify; `testcontainers-go` con PostgreSQL y Redis reales | BR-149 a BR-155 |

**Módulos verticales:** se activan por tenant mediante un middleware
`RequireModule("retail")` que consulta el plan (BR-015, BR-126).

**Concurrencia:** bloqueos optimistas con columna `version` en
conversaciones y pedidos (BR-048). Para stock y citas, `SELECT … FOR
UPDATE` y restricciones de exclusión de PostgreSQL contra la doble
reserva (BR-085, BR-103).

**Auditoría:** tabla `audit_logs` de solo inserción. El usuario de BD
de la aplicación no tiene permisos `UPDATE` ni `DELETE` sobre ella
(BR-131).

------------------------------------------------------------------------

# 5. Tiempo real (WebSocket)

**Librería:** `coder/websocket` en Go; WebSocket nativo del navegador
en Vue. No se usa Socket.IO porque sus servidores en Go están poco
mantenidos.

**Flujo:**

```
WhatsApp → webhook → guarda en BD → encola (River)
                                       ↓
                            worker procesa el mensaje
                                       ↓
                     publica evento en Redis pub/sub
                                       ↓
        cada instancia de la API reenvía a sus WebSockets
               (solo a usuarios autorizados del tenant)
```

**Reglas:**

-   Autenticación al conectar con un token de corta duración, nunca con
    el refresh token.
-   Canales `tenant:{id}`, `user:{id}` y `conversation:{id}`. Antes de
    suscribir se valida membresía y permisos (BR-047).
-   Cada evento lleva un ID incremental. Al reconectar, el cliente pide
    por REST lo que se perdió desde su último evento.
-   Heartbeat con ping/pong para detectar conexiones muertas.
-   Presencia de agentes (FR-068) y avisos de edición concurrente
    (FR-035) sobre los mismos canales.
-   Redis pub/sub distribuye eventos entre instancias; no se requieren
    sticky sessions.

**Evolución:** si se necesita presencia o recuperación avanzada a gran
escala, se puede migrar a Centrifugo (núcleo MIT) sin cambios
importantes en el frontend.

------------------------------------------------------------------------

# 6. Flujo del contrato API

1.  Se modifica `api/openapi.yaml`.
2.  `oapi-codegen` regenera interfaces y tipos en Go; el compilador
    marca los handlers que faltan o no coinciden.
3.  `openapi-typescript` regenera los tipos del cliente en
    `web/src/shared/api`.
4.  CI falla si el código generado no está actualizado.

------------------------------------------------------------------------

# 7. Datos e infraestructura

| Componente | Local (docker-compose) | Producción (propuesta) |
|---|---|---|
| PostgreSQL 17 | Contenedor | Servicio administrado (AWS RDS, Neon o similar) con backups PITR |
| Redis / Valkey | Contenedor | Servicio administrado (ElastiCache, Upstash o similar) |
| Archivos | SeaweedFS (Apache 2.0) | AWS S3 o Cloudflare R2 |
| Correo | Mailpit | Amazon SES o Resend |
| API y worker | `go run` con recarga (`air`) | Contenedores Docker (AWS ECS Fargate, Fly.io o Railway) |
| Web | Vite dev server | Estático en Cloudflare Pages o S3 + CloudFront |
| Webhooks en local | Túnel (cloudflared o ngrok) | Dominio público con HTTPS |

**CI/CD (GitHub Actions):**

-   Backend: `golangci-lint` → `go vet` → verificación de código
    generado (`sqlc`, `oapi-codegen`) → pruebas unitarias → pruebas de
    integración con Testcontainers → `govulncheck` → build de imagen.
-   Frontend: `pnpm` con caché → ESLint → `vue-tsc` → Vitest → build.
-   E2E: Playwright contra el entorno levantado con docker-compose.
-   Dependencias: Dependabot para Go y npm (BR-146, BR-156).

**Entornos:** `local`, `staging` (con número de prueba de WhatsApp) y
`production`.

------------------------------------------------------------------------

# 8. Integraciones externas

| Integración | Elección |
|---|---|
| WhatsApp | WhatsApp Business Platform – Cloud API + Embedded Signup |
| IA | Anthropic (Claude) como proveedor inicial, detrás de `AIProvider` |
| Pagos (fase posterior) | Proveedor con presencia en Perú: Culqi, Mercado Pago o Niubiz; Stripe para clientes internacionales. Nunca se almacenan datos de tarjeta (BR-095) |
| Facturación electrónica (fase posterior) | Proveedor OSE/PSE autorizado por SUNAT |

------------------------------------------------------------------------

# 9. Decisiones cerradas

-   **Frontend:** Vue 3 + Vite + shadcn-vue.
-   **Backend:** Go, con API y worker separados.
-   **Contrato:** OpenAPI spec-first.
-   **Multi-tenant:** esquema compartido con `tenant_id` + RLS. Queda
    abierta la opción de una base dedicada para clientes que lo exijan.
-   **Sesiones:** JWT de acceso corto + refresh rotativo en cookie
    `httpOnly`.
-   **Colas:** River sobre PostgreSQL.
-   **Tiempo real:** WebSocket propio con `coder/websocket` + Redis.
-   **IA inicial:** Anthropic, con capa desacoplada.

**Siguen pendientes:** proveedor de hosting definitivo, límites por
plan, política de retención, pasarela de pagos y revisión legal de la
Ley N.º 29733 (protección de datos personales en Perú) antes del
vertical de clínicas.

------------------------------------------------------------------------

# 10. Alternativas descartadas

| Decisión | Alternativa | Motivo |
|---|---|---|
| UI | PrimeVue 5 | Cambió a licencia PrimeUI (de pago para empresas que superan los límites de la licencia Community) |
| UI | PrimeVue 4 | MIT, pero sin nuevas versiones mayores |
| Backend | NestJS (Node.js) | Permite compartir tipos con Vue y avanza más rápido en CRUD; se prefirió Go por estabilidad bajo carga, menor consumo y manejo de WebSocket |
| Frontend | Nuxt | El panel es una SPA privada; no necesita SSR ni SEO |
| Colas | Asynq (Redis) | River permite encolar en la misma transacción que el dato, lo que simplifica la idempotencia |
| Tiempo real | Socket.IO | Servidores Go poco mantenidos |
| Tiempo real | Centrifugo desde el inicio | Añade un servicio más; se reserva para cuando la escala lo justifique |
| Multi-tenant | Un schema por tenant | Migraciones y operación complejas con cientos de tenants |
