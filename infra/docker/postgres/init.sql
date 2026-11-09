-- Rol de la aplicación: no es superusuario ni dueño de las tablas,
-- por lo que las políticas de Row Level Security se le aplican siempre.
-- Solo para desarrollo local; en producción la contraseña viene de un gestor de secretos.
CREATE ROLE saas_app LOGIN PASSWORD 'saas_app' NOSUPERUSER NOBYPASSRLS;
GRANT CONNECT ON DATABASE saas TO saas_app;
GRANT USAGE ON SCHEMA public TO saas_app;

-- Las tablas que cree el rol dueño (saas) mediante migraciones quedan
-- accesibles para la aplicación. Tablas especiales como audit_logs
-- revocan UPDATE/DELETE en su propia migración.
ALTER DEFAULT PRIVILEGES FOR ROLE saas IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO saas_app;
ALTER DEFAULT PRIVILEGES FOR ROLE saas IN SCHEMA public
  GRANT USAGE, SELECT ON SEQUENCES TO saas_app;
ALTER DEFAULT PRIVILEGES FOR ROLE saas IN SCHEMA public
  GRANT EXECUTE ON FUNCTIONS TO saas_app;
