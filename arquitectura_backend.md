# Arquitectura del Backend

## SaaS multiindustria de automatización comercial por WhatsApp

Este documento define la organización de archivos, las convenciones y
los flujos críticos de `backend/`. Se basa en `requerimientos_backend.md`
(BR / RNF) y `stack.md`. Complementa a `arquitectura_frontend.md`.

> **Revisión 2.** Incorpora la revisión de diseño contra el código
> existente: búsquedas previas al tenant, orden y concurrencia de
> mensajes, IP de cliente confiable, tiempo real, colas separadas,
> límites entre módulos y endurecimiento de RLS y auditoría. Los cambios
> que esto implica en el código actual se listan en la sección 9.

------------------------------------------------------------------------

# 1. Estilo de arquitectura

**Monolito modular** en Go, desplegado como **dos procesos** que
comparten el mismo código de dominio:

```
                    ┌────────────────────────────────────┐
  Frontend  ───────►│  cmd/api      REST /api/v1         │
  Agentes   ═══════►│               + tiempo real /rt    │
  Webhook Meta ────►│               + /webhooks          │
                    └───────────────┬────────────────────┘
                                    │ ambos usan el mismo código
                    ┌───────────────▼────────────────────┐
                    │  internal/modules/*   (dominio)    │
                    │  internal/platform/*  (BD, auth,   │
                    │                        tenant, RLS)│
                    └───────────────▲────────────────────┘
                                    │
                    ┌───────────────┴────────────────────┐
   Colas (River) ──►│  cmd/worker   mensajes, envíos,    │
                    │   automatizaciones, IA, recordat.  │
                    └────────────────────────────────────┘
                                    │
                PostgreSQL (RLS + River + LISTEN/NOTIFY) + Redis
```

- **`api`** atiende al frontend, mantiene las conexiones de tiempo real
  de los agentes y recibe los webhooks de WhatsApp. Responde rápido y
  delega el trabajo pesado (BR-043).
- **`worker`** procesa en segundo plano. Escala por separado (RNF-002).
- No son microservicios: un módulo puede separarse más adelante si hace
  falta, pero no se parte desde el inicio.

------------------------------------------------------------------------

# 2. Stack

| Necesidad | Elección | Requerimientos |
|---|---|---|
| Lenguaje | Go 1.25+ | — |
| Router HTTP | `chi` (compatible con `net/http`) | BR-001 |
| Contrato | `openapi.yaml` primero; `oapi-codegen` (`chi-server`, `strict-server`, `models`) | BR-004, BR-005 |
| Validación | Esquema OpenAPI (`kin-openapi`) + `go-playground/validator` | BR-003, BR-135 |
| Base de datos | PostgreSQL 17 + `pgx` | — |
| Consultas | `sqlc` (SQL escrito a mano → código Go tipado), **un paquete por módulo**; `squirrel` solo para filtros muy dinámicos | BR-136 |
| Migraciones | `goose` (esquema) + `river migrate` (tablas de la cola) | — |
| Aislamiento multi-tenant | Esquema compartido, `tenant_id` + RLS forzado; tenant fijado por transacción | BR-009 a BR-013 |
| Búsquedas sin tenant | Funciones `SECURITY DEFINER` acotadas | BR-157, BR-158 |
| Autenticación | JWT de acceso de 15 min + refresh token rotativo (hash en BD) | BR-016, BR-018 a BR-020 |
| Contraseñas | Argon2id | BR-017 |
| MFA | TOTP para propietarios y administradores | BR-026 |
| Autorización | RBAC por permisos, middleware en cada endpoint | BR-021, BR-022 |
| Cola de trabajos | River (sobre PostgreSQL), colas separadas por prioridad | BR-117 a BR-121, BR-164 |
| Tiempo real | `pg_notify` dentro de la transacción + hub en la API | BR-047, BR-163 |
| Rate limit y caché | Redis | BR-024, BR-140 |
| Logs | `slog` estructurado | RNF-005 |

**Sin ORM.** El SQL se escribe a mano y `sqlc` genera los structs y
funciones. Motivos: ver exactamente qué consulta corre dentro de qué
transacción con RLS, y detectar errores de columna o tipo al compilar.

------------------------------------------------------------------------

# 3. Capas y reglas de dependencia

Dentro de cada módulo, la dependencia va en un solo sentido:

```
handler  ──►  service  ──►  repository  ──►  store del módulo (sqlc)  ──►  PostgreSQL
                  │
                  ├──►  jobs (encolar en la misma transacción)
                  └──►  realtime.Notify (pg_notify en la misma transacción)
```

| Capa | Responsabilidad | No debe |
|---|---|---|
| `handler` | Implementa la interfaz generada por `oapi-codegen`; toma tenant y usuario del contexto; traduce errores a respuestas HTTP | Contener reglas de negocio ni acceder a la BD |
| `service` | Reglas de negocio, transacción (`WithTenantTx`), encolado de trabajos, eventos de tiempo real, auditoría | Conocer HTTP |
| `repository` | Envuelve el `store` del módulo y traduce errores de BD a errores del dominio | Abrir transacciones por su cuenta |
| `store` | Código generado por `sqlc` desde `db/queries/<modulo>.sql`; no se edita | — |

**Reglas generales:**

- `platform/` no importa de `modules/`.
- **Cada módulo tiene su propio paquete `store`** generado solo con sus
  consultas. Un módulo no puede llamar consultas de otro porque no
  existen en su paquete: el límite lo garantiza el compilador, no la
  revisión de código. Costo aceptado: `sqlc` repite los structs de
  tablas en cada paquete.
- Si un módulo necesita datos de otro, usa el `service` del módulo dueño
  o una interfaz pequeña definida en el módulo consumidor.
- Los módulos verticales (`retail`, `restaurant`, `clinic`, `hardware`)
  reutilizan los módulos comunes; los comunes no dependen de los
  verticales.
- El `tenant_id` y el `user_id` salen **siempre** del JWT y de la
  membresía autenticada, nunca del cuerpo de la petición (BR-007).

------------------------------------------------------------------------

# 4. Distribución de archivos

Los archivos marcados con **✓** ya existen; los marcados con **✎**
existen pero deben ajustarse (ver sección 9).

```
backend/
├── go.mod ✓
├── Dockerfile ✓
├── sqlc.yaml ✎                            # Una entrada por módulo
├── oapi-codegen.yaml ✓
│
├── cmd/
│   ├── api/main.go ✎                      # API + tiempo real + webhooks
│   └── worker/main.go ✎                   # Worker con colas separadas
│
├── db/
│   ├── migrations/                        # goose (las corre el rol dueño)
│   │   ├── 00001_core_tenancy.sql ✎       # tenants, users, memberships, audit_logs
│   │   ├── 00002_lookup_functions.sql     # Funciones SECURITY DEFINER sin tenant
│   │   ├── 00003_auth_sessions.sql        # refresh tokens, MFA, invitaciones, reset
│   │   ├── 00004_rbac.sql                 # roles, permisos, equipos
│   │   ├── 00005_contacts.sql
│   │   ├── 00006_whatsapp.sql             # cuentas, números, plantillas, webhook_events
│   │   ├── 00007_inbox.sql                # conversaciones, mensajes, notas, etiquetas
│   │   ├── 00008_automations.sql
│   │   ├── 00009_catalog_inventory.sql
│   │   ├── 00010_orders.sql
│   │   ├── 00011_plans_usage.sql
│   │   └── ...                            # Verticales: restaurant, clinic, quotations
│   └── queries/                           # Un archivo (o carpeta) por módulo
│       ├── tenants.sql ✓
│       ├── auth.sql
│       ├── users.sql
│       ├── contacts.sql
│       ├── whatsapp.sql
│       ├── inbox.sql
│       ├── automations.sql
│       ├── catalog.sql
│       ├── inventory.sql
│       ├── orders.sql
│       ├── plans.sql
│       ├── audit.sql
│       └── reports.sql
│
├── internal/
│   │
│   ├── platform/                          # ── TRANSVERSAL ──
│   │   ├── config/config.go ✓             # Variables de entorno
│   │   ├── database/
│   │   │   ├── postgres.go ✓              # Pool + WithTenantTx (RLS)
│   │   │   ├── lookup.go                  # Llamadas a las funciones sin tenant
│   │   │   └── guard.go                   # Verificación de privilegios al arrancar
│   │   ├── cache/redis.go ✓
│   │   ├── logger/logger.go ✓
│   │   ├── auth/
│   │   │   ├── jwt.go                     # Emisión y verificación del access token
│   │   │   ├── password.go                # Argon2id
│   │   │   ├── mfa.go                     # TOTP
│   │   │   ├── session.go                 # Sesión autenticada en el contexto
│   │   │   └── permissions.go             # Catálogo de permisos (RBAC)
│   │   ├── tenant/
│   │   │   ├── tenant.go                  # Tenant de la sesión
│   │   │   └── modules.go                 # Módulos habilitados por plan y rubro
│   │   ├── clientip/clientip.go           # IP real solo desde proxies de confianza
│   │   ├── crypto/crypto.go               # Cifrado de secretos en reposo (BR-052, BR-141)
│   │   ├── idempotency/idempotency.go     # Claves de idempotencia (BR-093)
│   │   ├── ratelimit/ratelimit.go         # Por IP y por cuenta (BR-024, BR-140)
│   │   ├── pagination/pagination.go       # Paginación por cursor
│   │   └── apperr/errors.go               # Errores de dominio → códigos HTTP
│   │
│   ├── http/                              # ── CAPA HTTP ──
│   │   ├── router.go ✎                    # Grupos: /api/v1, /rt, /webhooks
│   │   ├── respond.go ✓
│   │   ├── middleware/
│   │   │   ├── authenticate.go            # Valida JWT y carga la sesión
│   │   │   ├── authorize.go               # Verifica permisos por endpoint
│   │   │   ├── require_module.go          # Rechaza módulos no habilitados (BR-166)
│   │   │   ├── ratelimit.go
│   │   │   └── idempotency.go
│   │   └── api/api.gen.go                 # Generado por oapi-codegen (no se edita)
│   │
│   ├── realtime/                          # ── TIEMPO REAL ──
│   │   ├── notify.go                      # realtime.Notify(tx, evento) → pg_notify
│   │   ├── listener.go                    # LISTEN en una conexión dedicada por instancia
│   │   ├── hub.go                         # Conexiones de agentes agrupadas por tenant
│   │   └── events.go                      # Tipos de eventos (solo IDs y tipo)
│   │
│   ├── jobs/                              # ── COLAS (River) ──
│   │   ├── jobs.go ✓                      # Registro de workers
│   │   ├── queues.go                      # Colas, prioridades y concurrencia
│   │   ├── policy.go                      # Reintentos, backoff, descarte (BR-118, BR-119)
│   │   └── workers/
│   │       ├── process_webhook.go         # cola: inbound
│   │       ├── send_message.go            # cola: outbound
│   │       ├── run_automation.go          # cola: automations
│   │       ├── ai_reply.go                # cola: ai
│   │       ├── send_reminder.go           # cola: scheduled
│   │       └── cleanup.go                 # cola: maintenance (retención)
│   │
│   ├── integrations/                      # ── SERVICIOS EXTERNOS ──
│   │   ├── whatsapp/                      # Cliente técnico de Meta (sin reglas de negocio)
│   │   │   ├── client.go                  # Cloud API: envíos y plantillas
│   │   │   ├── webhook.go                 # Verificación de firma + respuesta rápida
│   │   │   ├── events.go                  # Parseo: un POST → varios eventos
│   │   │   └── onboarding.go              # Embedded Signup
│   │   ├── ai/client.go                   # Proveedor de IA + presupuesto por tenant
│   │   ├── storage/storage.go             # Archivos y medios
│   │   └── email/email.go
│   │
│   └── modules/                           # ── DOMINIO ──
│       │
│       │   ── Núcleo ──
│       ├── auth/                          # Login, refresh, logout, recuperación, MFA
│       ├── tenants/                       # Datos del negocio, sucursales, configuración
│       ├── users/                         # Usuarios, invitaciones, roles, equipos
│       ├── plans/                         # Plan, módulos, límites, consumo, gracia
│       ├── audit/                         # Bitácora de auditoría (solo inserción)
│       ├── contacts/                      # CRM y contactos
│       ├── inbox/                         # Conversaciones + mensajes + asignación,
│       │                                  #   estados, etiquetas y notas internas
│       ├── whatsapp/                      # Cuentas, números y plantillas (dominio)
│       ├── automations/                   # Reglas, ejecuciones, historial
│       ├── reports/                       # Métricas y reportes
│       │
│       │   ── Comunes a varios rubros ──
│       ├── catalog/                       # Productos, categorías, variantes
│       ├── inventory/                     # Existencias por sucursal/almacén
│       ├── orders/                        # Pedidos y estados
│       ├── quotations/                    # Cotizaciones
│       │
│       │   ── Verticales (solo si el tenant tiene el rubro) ──
│       ├── retail/                        # Tallas, colores, stock por variante
│       ├── restaurant/                    # Menú, modificadores, modalidades
│       ├── clinic/                        # Citas, profesionales, agenda, recordatorios
│       └── hardware/                      # Unidades de medida, almacenes
│
└── test/                                  # ── PRUEBAS DE INTEGRACIÓN ──
    ├── testdb.go                          # PostgreSQL real con testcontainers
    ├── tenant_isolation_test.go           # Un tenant nunca ve datos de otro (BR-014)
    ├── db_privileges_test.go              # RLS forzado, auditoría inmutable (BR-159, BR-168)
    ├── webhook_idempotency_test.go        # Reintentos de Meta no duplican (BR-152)
    └── webhook_ordering_test.go           # Eventos desordenados y concurrentes (BR-170)
```

## Molde de un módulo

```
modules/<dominio>/
├── handler.go       # Implementa la interfaz generada por oapi-codegen
├── service.go       # Reglas de negocio + WithTenantTx + auditoría + encolado + eventos
├── repository.go    # Envuelve el store del módulo
├── types.go         # Tipos, entradas y errores del módulo
├── store/           # Generado por sqlc desde db/queries/<dominio>.sql (no se edita)
└── service_test.go  # Pruebas
```

Los módulos pequeños pueden omitir `repository.go` si el `service` usa
su `store` directamente.

**Separación `integrations/whatsapp` vs `modules/whatsapp`:** la
integración solo sabe hablar con Meta (HTTP, firmas, formatos). El
módulo guarda y aplica las reglas de negocio (qué números tiene el
tenant, qué plantillas están aprobadas, límites del plan).

------------------------------------------------------------------------

# 5. Flujos transversales

## 5.1 Aislamiento multi-tenant (BR-009 a BR-014, RNF-012)

**Tablas de negocio.** Todas tienen `tenant_id`, RLS habilitado **y
forzado**, e índices que empiezan por `tenant_id`:

```sql
ALTER TABLE contacts ENABLE ROW LEVEL SECURITY;
ALTER TABLE contacts FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON contacts
    USING      (tenant_id = (SELECT NULLIF(current_setting('app.tenant_id', true), '')::uuid))
    WITH CHECK (tenant_id = (SELECT NULLIF(current_setting('app.tenant_id', true), '')::uuid));
```

- El `(SELECT ...)` hace que PostgreSQL evalúe el tenant una vez por
  consulta y no una vez por fila.
- `FORCE` aplica la política incluso al dueño de la tabla, como defensa
  si alguien conecta la aplicación con el rol equivocado.
- Toda operación de negocio corre dentro de `database.WithTenantTx`, que
  fija `app.tenant_id` solo para esa transacción.
- El `INSERT` toma el `tenant_id` de la transacción, no del cliente.
- Si una consulta corre sin la transacción, RLS devuelve cero filas en
  lugar de datos de otro tenant.

**Tablas de identidad (`tenants`, `users`, `memberships`).** También
tienen RLS:

| Tabla | Visible para la aplicación cuando… |
|---|---|
| `tenants` | `id = app.tenant_id`, o el usuario (`app.user_id`) tiene membresía activa en él |
| `users` | `id = app.user_id`, o el usuario es miembro del tenant actual |
| `memberships` | Pertenecen al tenant actual, o son del usuario actual (solo lectura) |

**Búsquedas sin tenant (BR-157).** Hay flujos que ocurren antes de
conocer el tenant. No se resuelven abriendo tablas, sino con funciones
`SECURITY DEFINER` (migración `00002_lookup_functions.sql`) que hacen
**una sola búsqueda exacta** y devuelven **lo mínimo**:

| Función | Usada por | Devuelve |
|---|---|---|
| `auth_find_user_by_email(email)` | Login | `id`, `password_hash`, `status` |
| `auth_find_reset_token(token_hash)` | Recuperar contraseña | `user_id`, `expires_at` |
| `auth_find_invitation(token_hash)` | Aceptar invitación | `tenant_id`, `email`, `role`, `expires_at` |
| `auth_find_refresh_token(token_hash)` | Refresh | `user_id`, `tenant_id`, `expires_at`, `revoked` |
| `whatsapp_resolve_tenant(phone_number_id)` | Webhook | `tenant_id` |

Reglas para estas funciones:

- Fijan `SET search_path = public` y solo aceptan búsqueda por igualdad
  exacta (nunca `LIKE`, listados ni paginación).
- Los tokens se buscan por su hash, nunca en claro.
- `saas_app` solo tiene `EXECUTE` sobre ellas; no tiene acceso directo
  a las filas de otros tenants.
- Agregar una función nueva de este tipo requiere revisión de seguridad.

**Tablas globales (BR-158).** Inventario de tablas sin `tenant_id`:

| Tabla | Por qué es global | Acceso de `saas_app` |
|---|---|---|
| `users` | Una persona pertenece a varios tenants (BR-027) | RLS + funciones de búsqueda |
| `tenants` | Es la propia raíz del aislamiento | RLS |
| `webhook_events` | Llega antes de saber el tenant | `INSERT` en la API; `SELECT`/`UPDATE` en el worker |
| `river_job` y demás tablas de River | Infraestructura de la cola | Lectura y escritura; los trabajos solo llevan IDs (BR-165) |
| `plans` | Catálogo comercial | Solo lectura |

**Verificación al arrancar (BR-159).** `platform/database/guard.go`
detiene la API y el worker si la conexión:

- es superusuario o tiene `BYPASSRLS`;
- es dueña de alguna tabla de negocio;
- tiene `UPDATE` o `DELETE` sobre `audit_logs`.

## 5.2 Webhooks de WhatsApp (BR-041 a BR-045, BR-053, BR-147, BR-160 a BR-162)

**En la API** (debe tardar milisegundos):

```
Meta ──POST /webhooks/whatsapp──►
   1. Verifica la firma con el secreto de la app          (BR-042, BR-147)
   2. En UNA transacción:
        INSERT webhook_events (payload, payload_sha256 UNIQUE, received_at)
        encola process_webhook(event_id) en la cola "inbound"
      Si el sha256 ya existe (reintento idéntico): no hace nada
   3. Responde 200                                         (BR-043)
```

**En el worker** (`process_webhook`):

```
   4. Parsea el payload: un POST puede traer VARIOS mensajes y estados
   5. Para cada evento:
        a. tenant := whatsapp_resolve_tenant(phone_number_id)
        b. WithTenantTx(tenant):
             - pg_advisory_xact_lock(hash(tenant, wa_id del cliente))  (BR-162)
             - upsert contacto y conversación abierta
             - mensaje:  INSERT ... ON CONFLICT (tenant_id, provider_message_id) DO NOTHING  (BR-161)
             - estado:   UPDATE ... WHERE status_rank < nuevo_rank                         (BR-160)
             - encola automatizaciones aplicables (cola "automations")
             - realtime.Notify(tx, "message.created", ids)                                (BR-163)
           COMMIT  → recién aquí llega el evento a los agentes
   6. Marca webhook_events como procesado
```

Puntos clave:

- **Deduplicación en dos niveles.** El hash del POST evita reprocesar
  reintentos idénticos; la restricción única por `provider_message_id`
  (`wamid`) evita duplicar mensajes aunque lleguen en POST distintos.
- **Orden.** El historial se ordena por la marca de tiempo del
  proveedor, no por la hora de inserción. Los estados tienen un rango
  (`sent=1, delivered=2, read=3`; `failed` es terminal) y solo avanzan.
- **Concurrencia.** El bloqueo consultivo por cliente serializa los
  eventos de una misma conversación sin frenar a las demás. Así no se
  crean dos conversaciones abiertas para el mismo cliente.
- **Número desconocido.** Si `whatsapp_resolve_tenant` no encuentra el
  número, el evento queda marcado como huérfano y genera una alerta; no
  se reintenta indefinidamente.

## 5.3 Colas de trabajo (BR-057, BR-117 a BR-121, BR-164, BR-165)

Colas separadas para que lo lento no frene lo urgente:

| Cola | Trabajos | Prioridad | Concurrencia orientativa |
|---|---|---|---|
| `inbound` | Procesar webhooks | Máxima | Alta |
| `outbound` | Enviar mensajes a Meta | Alta | Media (respeta límites de Meta) |
| `automations` | Ejecutar automatizaciones | Media | Media |
| `ai` | Respuestas y clasificación con IA | Media | Baja (costo y latencia) |
| `scheduled` | Recordatorios, mensajes programados | Media | Baja |
| `maintenance` | Limpieza por retención, reportes | Baja | Baja |

- Los trabajos se encolan dentro de la misma transacción que los datos
  que los originan (River `InsertTx`): o se guarda todo o nada.
- **Los argumentos de un trabajo solo llevan IDs** (`tenant_id`,
  `message_id`, …). El contenido se lee de las tablas con RLS al
  ejecutar. La tabla de la cola no guarda datos personales (BR-165).
- Cada trabajo corre su lógica dentro de `WithTenantTx` con el
  `tenant_id` de sus argumentos.
- `jobs/policy.go` define por tipo: intentos máximos, backoff
  exponencial y qué errores no se reintentan (p. ej., número bloqueado).
- Los trabajos agotados quedan en estado descartado: es la cola de
  fallos (BR-119). Se exponen en un panel interno y generan alerta.
- Los trabajos son idempotentes; además se usan trabajos únicos de
  River cuando el mismo trabajo no debe encolarse dos veces (BR-120).

## 5.4 Autenticación y autorización

- Access token JWT de 15 minutos; refresh token rotativo guardado como
  hash y entregado en cookie `httpOnly`, `Secure`, `SameSite=Strict`
  (BR-016, BR-018 a BR-020).
- Login: `auth_find_user_by_email` → verificar Argon2id → listar
  membresías con `app.user_id` fijado → emitir tokens para el tenant
  elegido.
- El middleware `authenticate` carga la sesión; `authorize` verifica el
  permiso de cada endpoint (BR-021, BR-022).
- MFA con TOTP para propietarios y administradores (BR-026).
- **Revocación:** al revocar una membresía, el acceso se corta como
  máximo al vencer el access token (15 min). Las operaciones sensibles
  (cambiar roles, desconectar WhatsApp, exportar datos) vuelven a
  comprobar la membresía en la base de datos.

## 5.5 IP del cliente y rate limiting (BR-024, BR-140, BR-169)

- La IP del cliente solo se toma de la cabecera que fija el
  balanceador o proxy de confianza, y solo si la petición viene de una
  IP de ese proxy (`platform/clientip`). Sin proxy configurado se usa la
  IP de la conexión.
- No se usa `middleware.RealIP` de chi, porque confía en
  `X-Forwarded-For` venga de quien venga.
- El rate limit del login y la recuperación se aplica **por IP y por
  cuenta (email)**, así que cambiar de IP no permite probar contraseñas
  sin límite.

## 5.6 Tiempo real (FR-031, BR-047, BR-163)

```
service / worker                        API (cada instancia)
   WithTenantTx {                          listener: LISTEN realtime
     ... cambios ...                            │
     realtime.Notify(tx, evento)  ──COMMIT──►   ▼
   }                                       hub: envía solo a los agentes
                                                del tenant con permiso
```

- `pg_notify` es transaccional: el evento **solo se entrega si la
  transacción hace commit**. Así nunca llega a los agentes un mensaje
  que se revirtió.
- El evento lleva tipo e IDs, no el contenido. El frontend vuelve a
  pedir los datos por la API (que aplica permisos y RLS).
- Cada instancia de la API mantiene una conexión dedicada con `LISTEN`
  y reparte a sus agentes conectados. Esto evita necesitar Redis pub/sub
  al inicio.
- **Rutas fuera de los timeouts.** El grupo `/rt` se monta fuera de
  `/api/v1` (que tiene timeout de 30 s) y el handler extiende su plazo
  de escritura con `http.ResponseController`; si no, el `WriteTimeout`
  del servidor cortaría la conexión.
- Al conectar, el agente se autentica con su token y queda suscrito
  solo a su tenant (BR-047).
- **Cuándo cambiar a Redis pub/sub:** si el volumen de eventos supera lo
  que `LISTEN/NOTIFY` maneja con holgura (cientos por segundo
  sostenidos) o si se necesita enviar contenido más grande.

## 5.7 Planes, módulos y límites (BR-015, BR-123 a BR-127, BR-166)

- `require_module.go` se aplica por grupo de rutas: los endpoints de
  `restaurant`, `clinic`, `retail`, `hardware`, `quotations`, `ai` y
  `automations` responden 403 si el tenant no tiene el módulo por plan o
  rubro.
- Los límites (usuarios, números, conversaciones, consumo de IA) se
  comprueban en el `service`, dentro de la transacción, no solo en el
  frontend (BR-126).
- Los trabajos del worker también comprueban el módulo: si un tenant
  pierde la IA, `ai_reply` se descarta sin llamar al proveedor.

## 5.8 Pedidos e inventario (RNF-004, BR-085, BR-093, BR-100)

- Los pedidos aceptan `Idempotency-Key` para que un reintento no cree un
  duplicado.
- El descuento de stock se hace dentro de la transacción, con
  `SELECT ... FOR UPDATE` o una restricción `CHECK (stock >= 0)`.
- El registro del pedido se separa de la confirmación de pago (BR-094).

## 5.9 Secretos y datos sensibles (BR-052, BR-141, BR-142)

- Los tokens de WhatsApp, IA y otras credenciales se guardan cifrados en
  reposo con `platform/crypto`.
- La clave de cifrado viene de una variable de entorno al inicio y de un
  KMS en producción.
- Los secretos nunca aparecen en logs, respuestas de la API ni
  argumentos de trabajos.

## 5.10 Auditoría (BR-025, BR-131, BR-167, BR-168)

- `audit_logs` es solo de inserción. El `REVOKE UPDATE, DELETE, TRUNCATE`
  se aplica **sin condiciones**: la migración falla si el rol
  `saas_app` no existe, en vez de saltarse la protección en silencio.
- `tenant_id` admite nulo para eventos de seguridad sin tenant (login
  fallido, token inválido). Esas filas no son visibles por la
  aplicación (la política exige igualdad con el tenant actual); solo
  las consulta el equipo de operación.
- El registro se escribe en la misma transacción que la acción
  auditada.
- `test/db_privileges_test.go` y `guard.go` comprueban que la protección
  sigue aplicada.

## 5.11 Observabilidad (RNF-005, BR-148)

- Logs estructurados con `slog`, con `request_id`, `tenant_id` y
  `user_id` en cada línea; los campos sensibles se ocultan.
- Métricas: latencia por endpoint, profundidad y antigüedad de cada
  cola, trabajos descartados, webhooks huérfanos, eventos de tiempo real
  enviados.
- Trazas con OpenTelemetry que unen la petición, el trabajo encolado y
  la llamada a Meta.
- Endpoints `/healthz` y `/readyz` (ya existen en `router.go`).

------------------------------------------------------------------------

# 6. Cómo escribir un módulo nuevo

1. **Contrato:** agregar los endpoints en `api/openapi.yaml` y regenerar
   con `go tool oapi-codegen -config oapi-codegen.yaml ../api/openapi.yaml`.
   Después regenerar los tipos del frontend con `pnpm api:types`.
2. **Migración:** crear `db/migrations/NNNNN_<modulo>.sql` con la tabla,
   `tenant_id`, índices que empiecen por `tenant_id`, RLS habilitado y
   forzado, y su `Down`.
3. **Consultas:** crear `db/queries/<modulo>.sql`, agregar la entrada
   del módulo en `sqlc.yaml` y ejecutar `sqlc generate`.
4. **Módulo:** crear `internal/modules/<modulo>/` con `handler.go`,
   `service.go`, `repository.go`, `types.go` y su `store/` generado.
5. **Registro:** conectar el handler en `internal/http/router.go`, con
   su permiso por endpoint y `require_module` si es opcional por plan o
   rubro.
6. **Auditoría, trabajos y eventos:** auditar las acciones relevantes,
   encolar lo lento en la cola adecuada (solo con IDs) y emitir eventos
   de tiempo real con `realtime.Notify` dentro de la transacción.
7. **Pruebas:** pruebas del `service` y una comprobación en
   `test/tenant_isolation_test.go`.

------------------------------------------------------------------------

# 7. Decisiones de diseño

| Decisión | Elección | Alternativa descartada | Motivo |
|---|---|---|---|
| Despliegue | Monolito modular, `api` + `worker` | Microservicios | Equipo pequeño; se puede separar después |
| Aislamiento | Esquema compartido + RLS forzado | Esquema o BD por tenant | Menor costo operativo; BR-013 lo permite |
| Búsquedas sin tenant | Funciones `SECURITY DEFINER` | Tablas sin RLS / rol con `BYPASSRLS` | Exponen solo un dato puntual, no tablas completas |
| Acceso a datos | `sqlc`, un paquete por módulo | ORM / un solo paquete `store` | SQL visible con RLS; límites entre módulos garantizados al compilar |
| Cola | River sobre PostgreSQL | Redis (asynq) | Encolado en la misma transacción que el dato |
| Tiempo real | `pg_notify` transaccional | Redis pub/sub | Solo entrega tras commit; una pieza menos al inicio |
| Orden de mensajes | Bloqueo consultivo por cliente + estados que solo avanzan | Cola FIFO única | Ordena por conversación sin frenar al resto |
| Conversaciones y mensajes | Un módulo `inbox` | Módulos separados | Siempre cambian juntos y en la misma transacción |

------------------------------------------------------------------------

# 8. Decisiones pendientes que afectan la estructura

Tomadas de la sección 9 de `requerimientos_backend.md`, más las que
surgieron en la revisión:

| Decisión | Qué cambia en la estructura |
|---|---|
| **WebSocket o SSE** | Define `realtime/hub.go`. SSE es más simple (solo servidor → cliente); WebSocket si el agente debe enviar por el mismo canal (p. ej., "escribiendo…") |
| **Proxy o balanceador de confianza** | Define qué cabecera y qué rangos de IP usa `platform/clientip` |
| **Aislamiento físico para clientes grandes** | Si algún cliente exige BD propia, cambian `database/` y `tenant/` |
| **Política de retención de mensajes y archivos** | Define los trabajos de la cola `maintenance` |
| **Proveedor de IA y presupuesto por tenant** | Define `integrations/ai/` y los límites en `plans/` |
| **Pasarela de pagos y facturación** | Probablemente requiere un módulo `billing/` |
| **Reglas de suscripción y periodo de gracia** | Define el módulo `plans/` |
| **Requisitos legales en Perú** (datos personales y de salud) | Afecta cifrado, retención y acceso en `clinic/` |
| **Hosting y almacenamiento** | Define `integrations/storage/` y el despliegue |

------------------------------------------------------------------------

# 9. Ajustes pendientes en el código existente

El código actual es el esqueleto inicial. Para alinearlo con este
documento hay que ajustar lo siguiente (todavía no se ha cambiado):

| Archivo | Ajuste | Motivo |
|---|---|---|
| `db/migrations/00001_core_tenancy.sql` | RLS en `tenants` y `users`; `FORCE ROW LEVEL SECURITY` en todas las tablas; políticas con `(SELECT current_setting(...))`; `audit_logs.tenant_id` nulo permitido; `REVOKE` sin condición | 5.1, 5.10. Se puede editar porque aún no está en producción |
| `db/queries/tenants.sql` | `GetTenantByID` debe depender de RLS, no del ID recibido | 5.1 |
| `internal/http/router.go` | Quitar `middleware.RealIP`; separar grupos `/api/v1`, `/rt` y `/webhooks` | 5.5, 5.6 |
| `cmd/api/main.go` | Mantener `WriteTimeout`, pero excluir las rutas de tiempo real con `ResponseController`; iniciar el `listener` de `LISTEN` | 5.6 |
| `cmd/worker/main.go` | Reemplazar la cola única (`MaxWorkers: 50`) por las colas de la sección 5.3; verificación de privilegios al arrancar | 5.3, 5.1 |
| `sqlc.yaml` | Una entrada por módulo con salida en `internal/modules/<modulo>/store` | 3 |
| `internal/platform/database/` | Agregar `guard.go` y `lookup.go` | 5.1 |

------------------------------------------------------------------------

# 10. Orden de construcción

Siguiendo el principio de `requerimientos_backend.md`: primero un núcleo
multi-tenant seguro; después automatización y un sector inicial.

| Orden | Alcance | Módulos y piezas |
|---|---|---|
| 0 | Ajustes del esqueleto | Sección 9 completa, `guard.go`, `test/db_privileges_test.go` |
| 1 | Núcleo seguro | `platform/*`, funciones de búsqueda, `auth`, `tenants`, `users`, `audit`, pruebas de aislamiento |
| 2 | WhatsApp y atención | `integrations/whatsapp`, `whatsapp`, `contacts`, `inbox`, `realtime`, colas `inbound` y `outbound`, pruebas de idempotencia y orden |
| 3 | Automatización | `automations`, cola `automations` |
| 4 | Primer vertical (tiendas) | `catalog`, `inventory`, `orders`, `retail` |
| 5 | IA, planes y métricas | `integrations/ai`, cola `ai`, `plans`, `reports` |
| 6 | Otros verticales | `restaurant`, `hardware`, `quotations`, `clinic` |
