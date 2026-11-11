# SaaS WhatsApp

Plataforma SaaS multi-tenant de atención y automatización comercial por
WhatsApp para tiendas, restaurantes, clínicas y ferreterías.

- Requerimientos: [requerimientos_frontend.md](requerimientos_frontend.md), [requerimientos_backend.md](requerimientos_backend.md)
- Stack y decisiones: [stack.md](stack.md)

## Estructura

```
api/openapi.yaml   Contrato único entre frontend y backend
backend/           Go: cmd/api (API REST) y cmd/worker (trabajos River)
web/               Vue 3 + Vite + shadcn-vue + Tailwind
infra/docker/      PostgreSQL, Valkey (Redis), SeaweedFS (S3) y Mailpit para desarrollo
```

## Requisitos

| Herramienta | Versión |
|---|---|
| Go | 1.25 o superior |
| Node.js | 24 LTS recomendado (mínimo 22.12) |
| pnpm | 10 o superior |
| Docker Desktop | Reciente |

> Con Node 22.14 el CLI de shadcn-vue y jsdom 27+ no funcionan. Por eso
> `web/` usa jsdom 26. Al actualizar a Node 24 se pueden subir.

## Puesta en marcha

Los comandos son para Git Bash. Ejecutarlos desde la raíz del proyecto.

### 1. Servicios de desarrollo

```bash
cp .env.example .env
```

```bash
docker compose -f infra/docker/docker-compose.yml up -d
```

- PostgreSQL: `localhost:5434` (dueño `saas`, aplicación `saas_app`)
- Valkey/Redis: `localhost:6380`
- S3 (SeaweedFS): `http://localhost:8333`
- Mailpit: `http://localhost:8025`

### 2. Backend

Instalar dependencias y herramientas (solo la primera vez):

```bash
cd backend && go mod tidy
```

```bash
cd backend && go get -tool github.com/pressly/goose/v3/cmd/goose@latest github.com/riverqueue/river/cmd/river@latest github.com/sqlc-dev/sqlc/cmd/sqlc@latest github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Migraciones (con el rol dueño del esquema):

```bash
set -a; source .env; set +a; cd backend && go tool goose -dir db/migrations postgres "$MIGRATIONS_DATABASE_URL" up && go tool river migrate-up --database-url "$MIGRATIONS_DATABASE_URL"
```

Levantar la API (`:8080`) y el worker, cada uno en su terminal:

```bash
set -a; source .env; set +a; cd backend && go run ./cmd/api
```

```bash
set -a; source .env; set +a; cd backend && go run ./cmd/worker
```

Comprobar: `http://localhost:8080/readyz` debe responder `{"status":"ok",...}`.

### 3. Frontend

```bash
cd web && pnpm install && pnpm dev
```

Abrir `http://localhost:5173`. En desarrollo, el login acepta cualquier
correo y contraseña (sesión de demostración) hasta que exista
`/api/v1/auth/login`. El panel muestra el estado de `/readyz`.

## Flujo de trabajo

**Cambiar la API:**
1. Editar `api/openapi.yaml`.
2. Regenerar Go: `cd backend && go tool oapi-codegen -config oapi-codegen.yaml ../api/openapi.yaml`
3. Regenerar tipos de Vue: `cd web && pnpm api:types`

**Cambiar la base de datos:**
1. Crear la migración: `cd backend && go tool goose -dir db/migrations create nombre sql`
2. Toda tabla de negocio lleva `tenant_id` y política RLS (ver `00001_core_tenancy.sql`).
3. Escribir las consultas en `db/queries/` y generar con `cd backend && go tool sqlc generate`.

**Acceso a datos por tenant:** usar siempre `database.WithTenantTx`, que fija
`app.tenant_id` en la transacción para que RLS filtre las filas. El tenant
sale de la membresía autenticada, nunca de un parámetro del cliente.

**Componentes de UI:** los componentes de shadcn-vue viven en
`web/src/shared/ui`. Con Node 22.15 o superior se agregan con
`pnpm dlx shadcn-vue@latest add <componente>`.

## Comandos útiles

| Comando | Qué hace |
|---|---|
| `cd web && pnpm build` | Type-check y build de producción |
| `cd web && pnpm test` | Pruebas unitarias (Vitest) |
| `cd backend && go test ./...` | Pruebas del backend |
| `cd backend && go vet ./...` | Análisis estático |
