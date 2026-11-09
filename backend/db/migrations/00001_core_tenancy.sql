-- +goose Up
-- Núcleo multi-tenant: tenants, usuarios globales, membresías y auditoría.
-- Patrón para todas las tablas de negocio: columna tenant_id + RLS.

CREATE TABLE tenants (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name          text        NOT NULL,
    slug          text        NOT NULL UNIQUE,
    industry      text        NOT NULL CHECK (industry IN ('retail', 'restaurant', 'clinic', 'hardware')),
    timezone      text        NOT NULL DEFAULT 'America/Lima',
    currency      char(3)     NOT NULL DEFAULT 'PEN',
    status        text        NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'cancelled')),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

-- Usuarios globales: una persona puede pertenecer a varios tenants (BR-027).
CREATE TABLE users (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email          text        NOT NULL,
    full_name      text        NOT NULL,
    password_hash  text,
    status         text        NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'pending', 'locked', 'disabled')),
    mfa_secret_enc bytea,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_key ON users (lower(email));

CREATE TABLE memberships (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  uuid        NOT NULL REFERENCES tenants (id) ON DELETE CASCADE,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role       text        NOT NULL CHECK (role IN ('owner', 'admin', 'supervisor', 'agent', 'operator')),
    status     text        NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'invited', 'revoked')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tenant_id, user_id)
);
CREATE INDEX memberships_user_idx ON memberships (user_id);

-- Antes de elegir tenant (al iniciar sesión), el usuario solo puede leer
-- sus propias membresías. Se fija con set_config('app.user_id', …, true).
-- Las políticas permisivas se combinan con OR.
CREATE POLICY own_memberships ON memberships
    FOR SELECT
    USING (user_id = NULLIF(current_setting('app.user_id', true), '')::uuid);

ALTER TABLE memberships ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON memberships
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid);

-- Bitácora de auditoría: solo inserción (BR-131).
CREATE TABLE audit_logs (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tenant_id    uuid        NOT NULL REFERENCES tenants (id),
    actor_id     uuid        REFERENCES users (id),
    action       text        NOT NULL,
    entity_type  text        NOT NULL,
    entity_id    text,
    metadata     jsonb       NOT NULL DEFAULT '{}',
    ip           inet,
    created_at   timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX audit_logs_tenant_created_idx ON audit_logs (tenant_id, created_at DESC);

ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON audit_logs
    USING (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid)
    WITH CHECK (tenant_id = NULLIF(current_setting('app.tenant_id', true), '')::uuid);

-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'saas_app') THEN
        REVOKE UPDATE, DELETE, TRUNCATE ON audit_logs FROM saas_app;
    END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down
DROP TABLE audit_logs;
DROP TABLE memberships;
DROP TABLE users;
DROP TABLE tenants;
