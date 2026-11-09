# Requerimientos del Frontend

## SaaS multiindustria de automatización comercial por WhatsApp

**Versión:** 1.0\
**Fecha:** 09 de noviembre de 2026\
**Alcance:** Frontend (aplicación web responsive)\
**Documento origen:** `requerimientos_saas_whatsapp_frontend_backend.md`\
**Documento relacionado:** `requerimientos_backend.md`

------------------------------------------------------------------------

# 1. Visión general del producto

## 1.1 Objetivo

Desarrollar una plataforma SaaS que permita a pequeñas y medianas
empresas centralizar la atención por WhatsApp, automatizar
conversaciones y gestionar los procesos comerciales que nacen de esas
conversaciones.

La plataforma tendrá un núcleo común para todos los negocios y módulos
especializados según el rubro contratado.

## 1.2 Sectores contemplados

1.  **Tiendas de ropa y zapatillas:** catálogo, variantes, tallas,
    colores, inventario, pedidos y seguimiento.
2.  **Restaurantes:** menú, modificadores, disponibilidad, pedidos,
    delivery o recojo y seguimiento.
3.  **Clínicas y consultorios:** servicios, profesionales, horarios,
    citas, confirmaciones y recordatorios.
4.  **Ferreterías:** catálogo, unidades de medida, inventario,
    cotizaciones, pedidos y seguimiento.

La plataforma no debe obligar a todos los negocios a utilizar los mismos
procesos. Cada tenant tendrá un rubro y un conjunto de módulos
habilitados.

## 1.3 Principios de arquitectura

-   Aplicación multi-tenant: cada negocio administra sus propios datos.
-   Aislamiento estricto entre tenants en frontend, API, base de datos,
    archivos, eventos y procesos en segundo plano.
-   Núcleo común reutilizable y módulos verticales activables por
    configuración o plan.
-   Integración oficial con WhatsApp Business Platform / Cloud API.
-   Procesamiento asíncrono para webhooks, mensajes, notificaciones y
    automatizaciones.
-   Trazabilidad de acciones relevantes.
-   Diseño responsive y accesible.
-   API versionada y contratos documentados.
-   La IA será complementaria: las operaciones críticas deben poder
    funcionar mediante reglas deterministas.

## 1.4 Roles del sistema

  -----------------------------------------------------------------------
  Rol                                 Alcance
  ----------------------------------- -----------------------------------
  Superadministrador SaaS             Administra la plataforma, tenants,
                                      planes, límites y soporte.

  Propietario del negocio             Administra su negocio,
                                      configuración, integraciones,
                                      usuarios y reportes.

  Administrador del negocio           Gestiona operaciones y usuarios
                                      según permisos.

  Supervisor                          Supervisa conversaciones, agentes y
                                      métricas.

  Agente                              Atiende conversaciones asignadas y
                                      consulta información autorizada.

  Operador                            Gestiona pedidos, citas,
                                      cotizaciones u operaciones según el
                                      rubro.
  -----------------------------------------------------------------------

Los permisos deben ser configurables y aplicarse en backend, no
solamente ocultando opciones en la interfaz.

------------------------------------------------------------------------

# 2. Requerimientos del Frontend

## 2.1 Aplicación y estructura general

**FR-001.** El frontend debe ofrecer una aplicación web responsive para
escritorio, tablet y móvil.

**FR-002.** Debe contar con una pantalla de autenticación.

**FR-003.** Después del inicio de sesión, debe identificar el tenant y
el contexto de negocio autorizado para el usuario.

**FR-004.** Debe mostrar navegación y módulos según el rol, permisos,
rubro y funcionalidades habilitadas.

**FR-005.** Debe permitir cambiar de sucursal o unidad de negocio cuando
el usuario tenga acceso a más de una.

**FR-006.** Debe mostrar estados de carga, vacío, error, éxito, conexión
y permisos insuficientes.

**FR-007.** Debe manejar expiración de sesión y errores de autenticación
de manera controlada.

**FR-008.** Debe centralizar componentes compartidos: tablas,
formularios, filtros, diálogos, paginación, notificaciones y estados.

## 2.2 Autenticación y cuenta

**FR-009.** Permitir iniciar sesión con las credenciales habilitadas por
el sistema.

**FR-010.** Permitir cerrar sesión.

**FR-011.** Permitir recuperar o restablecer contraseña mediante un
flujo seguro.

**FR-012.** Mostrar perfil del usuario y permitir actualizar los campos
autorizados.

**FR-013.** Mostrar sesiones o dispositivos activos si el backend ofrece
esta capacidad.

**FR-014.** Gestionar estados de cuenta: activa, pendiente, bloqueada o
deshabilitada.

## 2.3 Panel principal / Dashboard

**FR-015.** Mostrar indicadores del periodo seleccionado: -
Conversaciones recibidas. - Conversaciones pendientes. - Tiempo de
primera respuesta. - Conversaciones atendidas. - Clientes nuevos. -
Leads u oportunidades generadas. - Pedidos, citas o cotizaciones según
el rubro. - Tasa de conversión cuando exista información suficiente.

**FR-016.** Permitir filtrar indicadores por fecha, sucursal, canal y
agente cuando corresponda.

**FR-017.** Mostrar gráficos y tendencias sin mezclar información de
otros tenants.

**FR-018.** Mostrar tareas o conversaciones que requieren atención.

**FR-019.** Diferenciar datos reales, datos incompletos y métricas
todavía no disponibles.

## 2.4 Bandeja de conversaciones (Inbox)

**FR-020.** Mostrar conversaciones de los canales conectados, comenzando
por WhatsApp.

**FR-021.** Permitir filtrar por estado, asignación, agente, etiqueta,
fecha y búsqueda de texto.

**FR-022.** Mostrar datos del contacto, última interacción, canal,
estado y agente asignado.

**FR-023.** Mostrar el historial de mensajes en orden cronológico.

**FR-024.** Permitir enviar mensajes de texto y los tipos de contenido
soportados por la integración habilitada.

**FR-025.** Mostrar estados del mensaje cuando estén disponibles: en
cola, enviado, entregado, leído o fallido.

**FR-026.** Permitir asignar y reasignar una conversación a un agente
autorizado.

**FR-027.** Permitir transferir una conversación a otro equipo o cola.

**FR-028.** Permitir cambiar el estado: abierta, pendiente, resuelta o
cerrada, según las reglas del producto.

**FR-029.** Permitir agregar notas internas que no se envíen al cliente.

**FR-030.** Permitir agregar y quitar etiquetas.

**FR-031.** Mostrar indicadores de mensajes no leídos y nuevas
conversaciones en tiempo real.

**FR-032.** Mostrar avisos cuando una conversación requiera una
plantilla aprobada o exista una restricción de ventana de atención de
WhatsApp.

**FR-033.** Permitir consultar el perfil y las operaciones relacionadas
del contacto sin abandonar la conversación.

**FR-034.** Permitir que un agente solicite asistencia o transferencia a
un supervisor.

**FR-035.** Evitar que dos agentes realicen cambios incompatibles sin
advertencia cuando exista concurrencia.

## 2.5 CRM y contactos

**FR-036.** Mostrar listado paginado de contactos.

**FR-037.** Permitir buscar por nombre, teléfono, correo u otros campos
autorizados.

**FR-038.** Mostrar ficha del contacto con datos, etiquetas,
consentimiento, origen y actividad.

**FR-039.** Permitir crear y editar contactos.

**FR-040.** Permitir asociar conversaciones, pedidos, citas y
cotizaciones al contacto.

**FR-041.** Permitir clasificar contactos mediante etiquetas y
segmentos.

**FR-042.** Mostrar historial de interacciones.

**FR-043.** Identificar posibles duplicados y permitir su revisión por
usuarios autorizados.

**FR-044.** Mostrar y administrar preferencias de comunicación y
consentimiento cuando corresponda.

## 2.6 Catálogo y productos compartidos

**FR-045.** Permitir listar, crear, editar, archivar y consultar
productos o servicios.

**FR-046.** Permitir registrar nombre, descripción, código/SKU,
imágenes, precio, estado y categoría.

**FR-047.** Permitir administrar categorías y subcategorías.

**FR-048.** Permitir gestionar variantes y atributos según el rubro.

**FR-049.** Permitir visualizar disponibilidad e inventario cuando el
módulo esté habilitado.

**FR-050.** Permitir seleccionar productos o servicios desde una
conversación para compartir información al cliente.

## 2.7 Automatizaciones y chatbot

**FR-051.** Mostrar un listado de automatizaciones con nombre, estado,
canal y última modificación.

**FR-052.** Permitir crear, editar, duplicar, activar, pausar y archivar
automatizaciones.

**FR-053.** Permitir configurar disparadores, condiciones y acciones
mediante un editor visual o formulario guiado.

**FR-054.** Permitir configurar respuestas a preguntas frecuentes y
palabras clave.

**FR-055.** Permitir definir horarios de atención y mensajes fuera de
horario.

**FR-056.** Permitir derivar conversaciones a un agente o equipo según
reglas.

**FR-057.** Permitir configurar mensajes de bienvenida y cierre.

**FR-058.** Permitir probar una automatización en un entorno de prueba
antes de activarla.

**FR-059.** Mostrar historial de ejecuciones, errores y resultados.

**FR-060.** Permitir activar asistencia de IA con límites, instrucciones
y fuentes de información autorizadas.

**FR-061.** Permitir definir cuándo la IA debe abstenerse y transferir
la conversación a una persona.

**FR-062.** Mostrar claramente si un mensaje fue generado por una
automatización, IA o agente.

## 2.8 Usuarios, equipos y permisos

**FR-063.** Permitir invitar usuarios al tenant.

**FR-064.** Permitir asignar roles y permisos.

**FR-065.** Permitir activar, desactivar y revocar acceso a usuarios.

**FR-066.** Permitir crear equipos o colas de atención.

**FR-067.** Permitir asignar agentes a equipos.

**FR-068.** Mostrar disponibilidad de agentes si esa funcionalidad está
implementada.

## 2.9 Integración de WhatsApp

**FR-069.** Mostrar el estado de conexión de cada número o cuenta de
WhatsApp.

**FR-070.** Permitir iniciar el proceso de vinculación mediante el flujo
oficial de Meta habilitado para la aplicación.

**FR-071.** Mostrar nombre del número, identificadores no sensibles,
estado y calidad de conexión cuando estén disponibles.

**FR-072.** Permitir administrar plantillas mediante los estados
devueltos por Meta.

**FR-073.** Mostrar errores de configuración y pasos de resolución sin
exponer tokens ni secretos.

**FR-074.** Permitir desconectar una integración a usuarios con permiso.

## 2.10 Módulos verticales

### A. Tiendas de ropa y zapatillas

**FR-075.** Administrar variantes de talla, color, modelo y otros
atributos.

**FR-076.** Mostrar stock por variante y sucursal.

**FR-077.** Crear y consultar pedidos.

**FR-078.** Mostrar estados del pedido: recibido, confirmado,
preparando, enviado, entregado o cancelado, configurables.

**FR-079.** Permitir registrar datos de entrega o recojo.

### B. Restaurantes

**FR-080.** Administrar categorías del menú, platos, bebidas y
complementos.

**FR-081.** Permitir configurar modificadores, extras y observaciones.

**FR-082.** Administrar disponibilidad por horario y producto.

**FR-083.** Registrar modalidad de atención: delivery, recojo o consumo
en local, según el negocio.

**FR-084.** Crear y actualizar pedidos y sus estados.

**FR-085.** Mostrar notas y restricciones alimentarias comunicadas por
el cliente, sin convertirlas en recomendaciones médicas.

### C. Clínicas y consultorios

**FR-086.** Administrar servicios, profesionales, sedes y horarios.

**FR-087.** Mostrar disponibilidad de citas.

**FR-088.** Permitir registrar, confirmar, reprogramar y cancelar citas
conforme a permisos y políticas.

**FR-089.** Mostrar agenda por profesional y sede.

**FR-090.** Permitir configurar recordatorios y confirmaciones.

**FR-091.** Limitar la información clínica visible según permisos y
necesidad de acceso. El MVP debe evitar almacenar historias clínicas si
no son necesarias para la gestión de citas.

### D. Ferreterías

**FR-092.** Administrar productos con SKU, marca, unidad de medida y
categorías.

**FR-093.** Mostrar existencias por almacén o sucursal.

**FR-094.** Permitir crear cotizaciones con productos, cantidades,
precios y vigencia.

**FR-095.** Permitir convertir cotizaciones aceptadas en pedidos.

**FR-096.** Mostrar estado de cotización y pedido.

## 2.11 Reportes y auditoría

**FR-097.** Mostrar reportes de conversaciones, agentes,
automatizaciones y resultados comerciales.

**FR-098.** Permitir filtrar y exportar reportes según permisos.

**FR-099.** Mostrar historial de acciones relevantes: usuario, acción,
fecha y entidad afectada.

**FR-100.** No permitir que el frontend altere registros de auditoría.

## 2.12 Suscripción y configuración del negocio

**FR-101.** Mostrar plan actual, límites y consumo del tenant.

**FR-102.** Mostrar avisos al acercarse a límites de usuarios,
conversaciones, números o funcionalidades.

**FR-103.** Permitir administrar datos del negocio, zona horaria,
moneda, horarios y sucursales.

**FR-104.** Mostrar módulos disponibles y módulos no incluidos en el
plan.

**FR-105.** Permitir gestionar preferencias de notificaciones.

------------------------------------------------------------------------

# 3. Requerimientos no funcionales (Frontend)

**RNF-001. Rendimiento:** las vistas principales deben cargar de forma
eficiente y utilizar paginación para listados grandes.

**RNF-006. Mantenibilidad:** utilizar módulos desacoplados, componentes
compartidos, contratos claros con la API y documentación.

**RNF-007. Compatibilidad:** soportar navegadores modernos y diseño
responsive.

**RNF-008. Accesibilidad:** formularios, navegación y controles deben
ser utilizables mediante teclado y tecnologías de asistencia.

**RNF-009. Privacidad:** mostrar únicamente los datos necesarios según
rol y permisos; no persistir datos sensibles en el navegador.

**RNF-011. Internacionalización:** presentar fechas, moneda y zona
horaria según la configuración del negocio.

**RNF-012. Seguridad multi-tenant:** ningún fallo de interfaz debe
permitir consultar datos de otro negocio. La interfaz nunca reemplaza la
validación del backend.

------------------------------------------------------------------------

# 4. Módulos visibles por rubro

La navegación debe mostrar los módulos según el rubro y el plan del tenant:

  Módulo              Tiendas       Restaurantes   Clínicas       Ferreterías
  --------------- ---------------- -------------- ----------- -------------------
  Catálogo               Sí             Menú       Servicios          Sí
  Variantes        Tallas/colores      Extras      No aplica   Unidades/formatos
  Inventario             Sí           Opcional     No aplica          Sí
  Pedidos                Sí              Sí        No aplica          Sí
  Cotizaciones        Opcional        Opcional     No aplica          Sí
  Agenda             No aplica        Opcional        Sí           No aplica
  Profesionales      No aplica       No aplica        Sí           No aplica
  Delivery            Opcional           Sí        No aplica       Opcional

------------------------------------------------------------------------

# 5. Endpoints a consumir (propuesta inicial)

Los endpoints son una propuesta inicial. Deben ajustarse al diseño final
del dominio y a los contratos de integración.

  Dominio         Método     Ruta
  --------------- ---------- -----------------------------------------
  Auth            POST       `/api/v1/auth/login`
  Auth            POST       `/api/v1/auth/refresh`
  Auth            POST       `/api/v1/auth/logout`
  Auth            POST       `/api/v1/auth/forgot-password`
  Auth            POST       `/api/v1/auth/reset-password`
  Users           GET        `/api/v1/users`
  Users           POST       `/api/v1/users/invitations`
  Tenants         GET        `/api/v1/tenants/current`
  Tenants         PATCH      `/api/v1/tenants/current`
  Conversations   GET        `/api/v1/conversations`
  Conversations   GET        `/api/v1/conversations/:id`
  Conversations   PATCH      `/api/v1/conversations/:id`
  Messages        GET        `/api/v1/conversations/:id/messages`
  Messages        POST       `/api/v1/conversations/:id/messages`
  Contacts        GET        `/api/v1/contacts`
  Contacts        POST       `/api/v1/contacts`
  Contacts        GET        `/api/v1/contacts/:id`
  Contacts        PATCH      `/api/v1/contacts/:id`
  Automations     GET        `/api/v1/automations`
  Automations     POST       `/api/v1/automations`
  Automations     PATCH      `/api/v1/automations/:id`
  Automations     POST       `/api/v1/automations/:id/activate`
  Automations     POST       `/api/v1/automations/:id/pause`
  Products        GET        `/api/v1/products`
  Products        POST       `/api/v1/products`
  Products        PATCH      `/api/v1/products/:id`
  Orders          GET        `/api/v1/orders`
  Orders          POST       `/api/v1/orders`
  Orders          GET        `/api/v1/orders/:id`
  WhatsApp        GET        `/api/v1/integrations/whatsapp`
  WhatsApp        POST       `/api/v1/integrations/whatsapp/connect`
  WhatsApp        DELETE     `/api/v1/integrations/whatsapp/:id`
  Webhooks        GET/POST   `/api/v1/webhooks/whatsapp`
  Reports         GET        `/api/v1/reports/overview`
  Audit           GET        `/api/v1/audit-logs`

Los endpoints verticales se incorporarán según el módulo:

-   `/api/v1/restaurant/menu`
-   `/api/v1/restaurant/orders`
-   `/api/v1/clinic/appointments`
-   `/api/v1/clinic/professionals`
-   `/api/v1/retail/variants`
-   `/api/v1/inventory`
-   `/api/v1/quotations`

Además, el frontend debe consumir un canal de eventos en tiempo real
(WebSocket o SSE) para nuevas conversaciones, mensajes y cambios de
estado.

------------------------------------------------------------------------

# 6. Priorización del MVP (Frontend)

## Fase 1 --- Núcleo SaaS

-   Login, logout, recuperación de contraseña y perfil.
-   Selección de tenant/sucursal y navegación por permisos.
-   Configuración del negocio y gestión de usuarios.
-   Pantalla de conexión de WhatsApp.
-   Inbox: listado, historial, envío de mensajes y estados.
-   Contactos / CRM básico.
-   Asignación, etiquetas y notas internas.
-   Vista básica de auditoría.

## Fase 2 --- Automatización

-   Listado y editor de automatizaciones.
-   Respuestas automáticas, palabras clave y horarios.
-   Reglas de derivación a agentes.
-   Historial de ejecuciones.
-   Gestión de plantillas de WhatsApp.

## Fase 3 --- Primer vertical: tiendas de ropa y zapatillas

-   Catálogo con variantes de tallas y colores.
-   Vista de inventario.
-   Pedidos y estados de pedido.
-   Compartir productos desde la conversación.

## Fase 4 --- IA y métricas

-   Configuración del asistente IA y reglas de transferencia.
-   Indicador del origen del mensaje (IA, automatización o agente).
-   Visualización de consumo y límites.
-   Dashboard y reportes de atención, conversión y automatizaciones.

## Fase 5 --- Otros verticales

1.  Restaurantes.
2.  Ferreterías.
3.  Clínicas y consultorios.

------------------------------------------------------------------------

# 7. Fuera del alcance inicial del MVP

-   Aplicación móvil nativa (se cubre con la web responsive).
-   Marketplace público de productos.
-   Interfaz para todos los canales sociales.
-   Interfaz de historia clínica electrónica.

------------------------------------------------------------------------

# 8. Criterios de aceptación (Frontend)

El frontend del MVP no se considera listo para producción hasta
verificar que:

1.  Un usuario solo ve los tenants y sucursales de los que es miembro.
2.  La navegación y las acciones respetan rol, permisos, rubro y plan.
3.  Un agente puede recibir, responder, asignar y cerrar conversaciones
    desde el Inbox.
4.  Los mensajes fallidos se identifican visualmente y pueden
    gestionarse.
5.  Una automatización puede crearse, probarse, activarse y pausarse
    desde la interfaz.
6.  Los secretos y tokens de integraciones nunca se muestran ni se
    almacenan en el navegador.
7.  La expiración de sesión y los errores de autenticación se manejan
    sin pantallas rotas.
8.  Las vistas son responsive y navegables por teclado.
9.  Existen pruebas end-to-end de los flujos críticos.

------------------------------------------------------------------------

# 9. Decisiones pendientes antes de implementar (Frontend)

-   Framework definitivo del frontend.
-   Modelo de autenticación y sesiones (acordar con backend).
-   Tecnología de tiempo real (WebSocket o SSE).
-   Librería de componentes / sistema de diseño.
-   Flujo de onboarding y vinculación de WhatsApp (registro embebido de Meta).
-   Estrategia de internacionalización.

**Principio de implementación:** primero construir un núcleo
multi-tenant seguro y funcional; después incorporar automatización y un
sector comercial inicial. Los demás sectores deben reutilizar el núcleo
sin introducir dependencias innecesarias entre sus reglas de negocio.
