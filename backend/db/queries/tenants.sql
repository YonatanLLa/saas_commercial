-- Consultas para sqlc. Ejecutar `sqlc generate` en backend/ para crear el código Go.

-- name: GetTenantByID :one
SELECT * FROM tenants WHERE id = $1;

-- name: ListMembershipsForUser :many
-- Se usa al iniciar sesión para saber a qué tenants pertenece el usuario (BR-009).
-- Requiere app.user_id fijado en la transacción (política own_memberships).
SELECT m.tenant_id, m.role, t.name AS tenant_name, t.industry
FROM memberships m
JOIN tenants t ON t.id = m.tenant_id
WHERE m.user_id = $1 AND m.status = 'active' AND t.status = 'active'
ORDER BY t.name;
