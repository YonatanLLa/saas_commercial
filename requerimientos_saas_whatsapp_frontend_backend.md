# Requerimientos funcionales y técnicos

## SaaS multiindustria de automatización comercial por WhatsApp

**Versión:** 1.0\
**Fecha:** 09 de noviembre de 2026\
**Alcance:** Frontend y Backend separados\
**Modelo:** SaaS multi-tenant para tiendas de ropa y zapatillas,
restaurantes, clínicas y consultorios, y ferreterías.

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

# 3. Requerimientos del Backend

## 3.1 Arquitectura y API

**BR-001.** Implementar una API REST versionada, por ejemplo `/api/v1`.

**BR-002.** Organizar el backend por módulos de dominio y
responsabilidades.

**BR-003.** Validar todas las entradas mediante DTOs o esquemas
equivalentes.

**BR-004.** Estandarizar respuestas, errores, paginación, filtros y
ordenamiento.

**BR-005.** Documentar endpoints mediante OpenAPI/Swagger.

**BR-006.** Separar controladores, servicios de aplicación, dominio,
persistencia e integraciones cuando la complejidad lo justifique.

**BR-007.** No confiar en identificadores de tenant enviados por el
cliente sin verificar su pertenencia y autorización.

## 3.2 Multi-tenancy

**BR-008.** Implementar entidades de tenant, usuarios, membresías,
sucursales y planes.

**BR-009.** Resolver el tenant a partir de una identidad y membresía
autenticadas.

**BR-010.** Aplicar aislamiento por tenant en todas las consultas y
comandos.

**BR-011.** Evitar accesos cruzados mediante identificadores manipulados
(IDOR/BOLA).

**BR-012.** Asociar conversaciones, contactos, productos, pedidos,
citas, automatizaciones y archivos con su tenant.

**BR-013.** Definir estrategia de aislamiento de datos: esquema
compartido con `tenant_id` y controles estrictos, o aislamiento físico
para clientes que lo requieran.

**BR-014.** Crear pruebas automatizadas que intenten acceder a datos de
otro tenant.

**BR-015.** Mantener configuración, integraciones, límites y módulos
habilitados por tenant.

## 3.3 Autenticación, autorización y sesiones

**BR-016.** Implementar autenticación segura.

**BR-017.** Almacenar contraseñas con algoritmo de hash robusto y
parámetros actualizados.

**BR-018.** Implementar sesiones o tokens de acceso de corta duración y
renovación segura.

**BR-019.** Proteger refresh tokens mediante rotación, revocación y
almacenamiento seguro.

**BR-020.** Implementar cierre de sesión y revocación de sesiones.

**BR-021.** Aplicar autorización basada en roles y permisos (RBAC).

**BR-022.** Verificar permisos en cada endpoint y operación sensible.

**BR-023.** Implementar recuperación de contraseña con tokens de un solo
uso y expiración.

**BR-024.** Aplicar rate limiting a login, recuperación y endpoints
sensibles.

**BR-025.** Registrar intentos relevantes de autenticación y eventos de
seguridad.

**BR-026.** Considerar MFA para propietarios y administradores.

## 3.4 Usuarios, equipos y permisos

**BR-027.** Gestionar usuarios globales y membresías por tenant.

**BR-028.** Permitir invitaciones con expiración y aceptación segura.

**BR-029.** Gestionar roles, permisos y equipos.

**BR-030.** Verificar que un usuario pertenezca al tenant y equipo
requerido antes de consultar o modificar recursos.

**BR-031.** Permitir revocar acceso sin eliminar el historial operativo.

## 3.5 Contactos y CRM

**BR-032.** Implementar CRUD de contactos por tenant.

**BR-033.** Normalizar números telefónicos y almacenar el identificador
del canal cuando corresponda.

**BR-034.** Gestionar etiquetas, segmentos y origen de captación.

**BR-035.** Asociar contactos con conversaciones y operaciones
comerciales.

**BR-036.** Detectar posibles duplicados según reglas configurables.

**BR-037.** Registrar preferencias y consentimiento de comunicación.

**BR-038.** Mantener historial de actividad y cambios relevantes.

## 3.6 Conversaciones y mensajes

**BR-039.** Implementar entidades de canales, conversaciones,
participantes, mensajes, asignaciones, etiquetas y notas internas.

**BR-040.** Persistir mensajes entrantes y salientes con identificadores
externos.

**BR-041.** Procesar webhooks entrantes de manera idempotente.

**BR-042.** Validar la autenticidad de los webhooks según el mecanismo
oficial del proveedor.

**BR-043.** Responder rápidamente al webhook y procesar tareas pesadas
mediante una cola.

**BR-044.** Mantener estados de envío y errores de mensajes.

**BR-045.** Evitar duplicados por reintentos del proveedor.

**BR-046.** Gestionar asignación, transferencia, cierre y reapertura de
conversaciones.

**BR-047.** Publicar eventos de actualización en tiempo real únicamente
a usuarios autorizados del tenant.

**BR-048.** Aplicar controles de concurrencia para evitar asignaciones o
cambios inconsistentes.

**BR-049.** Almacenar referencias a medios y metadatos de forma segura;
validar tipo y tamaño de archivos.

## 3.7 Integración con WhatsApp Business Platform

**BR-050.** Integrar mediante las APIs oficiales de Meta.

**BR-051.** Implementar el proceso de vinculación de cuentas y números
conforme al flujo autorizado de Meta.

**BR-052.** Gestionar credenciales y tokens cifrados en reposo y nunca
exponerlos al frontend.

**BR-053.** Implementar verificación de webhooks.

**BR-054.** Procesar mensajes entrantes, estados de mensajes y eventos
de plantillas.

**BR-055.** Enviar mensajes mediante un servicio de integración aislado
del dominio comercial.

**BR-056.** Validar las reglas de mensajería, plantillas aprobadas y
ventanas de atención vigentes.

**BR-057.** Gestionar límites, errores, reintentos con backoff y fallos
permanentes.

**BR-058.** Registrar identificadores externos y costos o métricas de
consumo cuando estén disponibles.

**BR-059.** Permitir múltiples números o cuentas por tenant si el plan y
la integración lo admiten.

**BR-060.** Implementar alertas ante desconexión, token inválido o
errores repetidos.

## 3.8 Motor de automatizaciones

**BR-061.** Modelar automatizaciones mediante disparadores, condiciones
y acciones.

**BR-062.** Soportar disparadores como mensaje recibido, palabra clave,
horario, cambio de estado y evento comercial.

**BR-063.** Soportar condiciones sobre horario, etiquetas, campos del
contacto, canal y contexto autorizado.

**BR-064.** Soportar acciones como enviar respuesta, asignar agente,
agregar etiqueta, actualizar estado y crear una operación.

**BR-065.** Permitir activar, pausar y versionar automatizaciones.

**BR-066.** Evitar ejecuciones duplicadas mediante claves de
idempotencia y controles de concurrencia.

**BR-067.** Definir límites de ejecución, profundidad y tiempo para
evitar ciclos infinitos.

**BR-068.** Registrar ejecuciones, resultados, errores y duración.

**BR-069.** Permitir probar flujos sin afectar conversaciones reales.

**BR-070.** Implementar mecanismos de pausa y transferencia a agente
humano.

## 3.9 Inteligencia artificial

**BR-071.** Integrar modelos de IA mediante una capa desacoplada del
proveedor.

**BR-072.** Configurar instrucciones, contexto y fuentes de conocimiento
por tenant.

**BR-073.** Recuperar únicamente información autorizada del negocio
correspondiente.

**BR-074.** No permitir que el modelo ejecute operaciones sensibles sin
validación de reglas y permisos del backend.

**BR-075.** Validar las respuestas y acciones propuestas antes de
ejecutarlas.

**BR-076.** Definir límites de consumo, presupuesto, frecuencia y
longitud de contexto.

**BR-077.** Registrar proveedor, modelo, consumo y resultado técnico,
evitando almacenar datos sensibles innecesarios.

**BR-078.** Permitir desactivar la IA y mantener respuestas basadas en
reglas.

**BR-079.** Transferir a una persona cuando no exista información
suficiente, haya baja confianza o el cliente lo solicite.

**BR-080.** Proteger el sistema frente a instrucciones maliciosas
recibidas en mensajes o documentos de clientes.

## 3.10 Catálogo e inventario común

**BR-081.** Gestionar productos, servicios, categorías, precios e
imágenes por tenant.

**BR-082.** Gestionar variantes y atributos configurables.

**BR-083.** Implementar inventario por sucursal o almacén cuando
corresponda.

**BR-084.** Registrar movimientos de inventario y su motivo.

**BR-085.** Evitar modificaciones concurrentes que generen stock
inconsistente.

**BR-086.** Mantener historial de precios y cambios relevantes cuando
sea necesario.

**BR-087.** Permitir productos archivados sin perder su historial en
pedidos.

## 3.11 Pedidos y ventas

**BR-088.** Implementar pedidos asociados a tenant y contacto.

**BR-089.** Registrar detalle, cantidades, precios aplicados, descuentos
y totales.

**BR-090.** Mantener estados de pedido mediante transiciones válidas.

**BR-091.** Registrar dirección, modalidad de entrega y observaciones
cuando corresponda.

**BR-092.** Mantener historial de cambios de estado.

**BR-093.** Evitar duplicar pedidos cuando se reintente una operación.

**BR-094.** Separar registro del pedido de la confirmación de pago.

**BR-095.** No almacenar datos de tarjetas ni información de pago
sensible; delegar pagos a proveedores especializados.

## 3.12 Módulo de restaurantes

**BR-096.** Gestionar menú, categorías, modificadores, extras y
disponibilidad.

**BR-097.** Gestionar horarios de atención y modalidades de servicio.

**BR-098.** Registrar pedidos y sus estados operativos.

**BR-099.** Gestionar zonas, tarifas y condiciones de delivery si el
alcance lo incluye.

**BR-100.** Validar disponibilidad antes de confirmar pedidos.

## 3.13 Módulo de clínicas y consultorios

**BR-101.** Gestionar profesionales, especialidades, servicios, sedes y
horarios.

**BR-102.** Gestionar disponibilidad y reservas de citas.

**BR-103.** Evitar doble reserva mediante transacciones y restricciones
adecuadas.

**BR-104.** Gestionar confirmación, reprogramación y cancelación.

**BR-105.** Registrar recordatorios y resultados de notificación.

**BR-106.** Aplicar controles reforzados de acceso, minimización y
protección de datos personales y de salud.

**BR-107.** Separar agenda administrativa de historia clínica. La
historia clínica queda fuera del MVP salvo que se defina expresamente su
alcance y cumplimiento normativo.

## 3.14 Módulo de ferreterías

**BR-108.** Gestionar marcas, unidades de medida, categorías y códigos
SKU.

**BR-109.** Gestionar inventario por almacén.

**BR-110.** Crear cotizaciones con vigencia, detalle, cantidades y
precios.

**BR-111.** Gestionar estados de cotización: borrador, enviada,
aceptada, rechazada o vencida.

**BR-112.** Convertir una cotización aceptada en pedido sin duplicar
información.

## 3.15 Módulo de tiendas de ropa y zapatillas

**BR-113.** Gestionar productos con variantes de talla, color, modelo y
SKU.

**BR-114.** Mantener stock por variante y sucursal.

**BR-115.** Validar disponibilidad antes de confirmar pedidos.

**BR-116.** Gestionar estados de preparación, despacho, entrega y
cancelación.

## 3.16 Notificaciones y tareas en segundo plano

**BR-117.** Implementar colas para envío de mensajes, automatizaciones,
recordatorios y tareas pesadas.

**BR-118.** Reintentar fallos transitorios con backoff y límite de
intentos.

**BR-119.** Enviar tareas fallidas a una cola de fallos (dead-letter) o
mecanismo equivalente.

**BR-120.** Evitar ejecuciones duplicadas.

**BR-121.** Permitir consultar estado y trazabilidad de tareas.

**BR-122.** Gestionar zona horaria por tenant para horarios y
recordatorios.

## 3.17 Planes, suscripciones y límites

**BR-123.** Gestionar planes, funcionalidades y límites.

**BR-124.** Asociar suscripciones a tenants.

**BR-125.** Medir consumo de conversaciones, agentes, números,
automatizaciones y uso de IA según el modelo comercial.

**BR-126.** Aplicar límites en backend, incluso si el frontend no los
muestra.

**BR-127.** Gestionar estados de suscripción y periodos de gracia.

**BR-128.** Registrar eventos de facturación y pagos mediante un
proveedor cuando se incorpore el módulo.

## 3.18 Reportes y auditoría

**BR-129.** Generar métricas por tenant y periodo.

**BR-130.** Calcular indicadores de atención y operación con reglas
documentadas.

**BR-131.** Registrar acciones sensibles en una bitácora de auditoría
protegida contra modificaciones ordinarias.

**BR-132.** Permitir consultar auditoría según permisos.

**BR-133.** Evitar exponer información personal innecesaria en
exportaciones.

## 3.19 Seguridad

**BR-134.** Aplicar HTTPS en todos los entornos desplegados.

**BR-135.** Aplicar validación y sanitización de entradas.

**BR-136.** Proteger contra SQL injection mediante consultas
parametrizadas u ORM seguro.

**BR-137.** Configurar CORS con orígenes explícitos.

**BR-138.** Aplicar protección CSRF cuando la estrategia de
autenticación utilice cookies.

**BR-139.** Aplicar políticas de seguridad de contenido y cabeceras HTTP
apropiadas.

**BR-140.** Aplicar rate limiting y controles antiabuso.

**BR-141.** Cifrar secretos y credenciales de integraciones.

**BR-142.** No registrar contraseñas, tokens, secretos ni contenido
sensible innecesario en logs.

**BR-143.** Validar permisos y pertenencia al tenant en cada recurso.

**BR-144.** Implementar copias de seguridad y pruebas de restauración.

**BR-145.** Definir políticas de retención, eliminación y exportación de
datos.

**BR-146.** Mantener dependencias actualizadas y análisis de
vulnerabilidades.

**BR-147.** Proteger endpoints de webhooks contra falsificación y
repetición cuando aplique.

**BR-148.** Implementar monitoreo, alertas y registro de errores sin
filtrar datos sensibles.

## 3.20 Pruebas y calidad

**BR-149.** Pruebas unitarias de reglas de negocio.

**BR-150.** Pruebas de integración de base de datos, colas y proveedores
externos.

**BR-151.** Pruebas de autorización por rol y tenant.

**BR-152.** Pruebas de idempotencia para webhooks, mensajes, pedidos y
automatizaciones.

**BR-153.** Pruebas de concurrencia para inventario y reservas.

**BR-154.** Pruebas de contrato para integraciones externas.

**BR-155.** Pruebas end-to-end de flujos críticos.

**BR-156.** Pipeline CI con lint, type-check, pruebas y build.

------------------------------------------------------------------------

# 4. Requerimientos no funcionales

Estos requerimientos aplican tanto al frontend como al backend.

**RNF-001. Rendimiento:** las vistas principales deben cargar de forma
eficiente y utilizar paginación para listados grandes.

**RNF-002. Escalabilidad:** el procesamiento de mensajes y
automatizaciones debe poder escalar independientemente de la API web.

**RNF-003. Disponibilidad:** definir objetivos de disponibilidad y
recuperación antes de producción.

**RNF-004. Consistencia:** pedidos, inventario, pagos y reservas deben
protegerse mediante transacciones y reglas de dominio.

**RNF-005. Observabilidad:** disponer de logs estructurados, métricas,
trazas y alertas.

**RNF-006. Mantenibilidad:** utilizar módulos desacoplados, contratos
claros y documentación.

**RNF-007. Compatibilidad:** soportar navegadores modernos y diseño
responsive.

**RNF-008. Accesibilidad:** formularios, navegación y controles deben
ser utilizables mediante teclado y tecnologías de asistencia.

**RNF-009. Privacidad:** recopilar únicamente los datos necesarios y
definir su finalidad, acceso y conservación.

**RNF-010. Recuperación:** establecer RPO, RTO, backups y procedimientos
de restauración.

**RNF-011. Internacionalización:** almacenar fechas de forma consistente
y presentar fechas, moneda y zona horaria según la configuración del
negocio.

**RNF-012. Seguridad multi-tenant:** ningún fallo de interfaz o
manipulación de solicitudes debe permitir consultar datos de otro
negocio.

------------------------------------------------------------------------

# 5. Modelo de módulos sugerido

## 5.1 Núcleo común

-   Auth
-   Tenants
-   Memberships
-   Users
-   Roles y permisos
-   Branches
-   WhatsApp
-   Channels
-   Conversations
-   Messages
-   Contacts / CRM
-   Inbox
-   Automation Engine
-   Notifications
-   Reports
-   Audit
-   Plans / Subscriptions
-   Files

## 5.2 Módulos verticales

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

# 6. API: agrupación inicial de endpoints

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

------------------------------------------------------------------------

# 7. Priorización del MVP

No conviene construir los cuatro sectores completos desde el primer
lanzamiento.

## Fase 1 --- Núcleo SaaS

-   Autenticación y autorización.
-   Multi-tenancy.
-   Gestión de negocio y usuarios.
-   Integración inicial con WhatsApp.
-   Webhooks y recepción de mensajes.
-   Inbox y envío de mensajes.
-   Contactos / CRM básico.
-   Asignación de conversaciones.
-   Etiquetas y notas internas.
-   Logs y auditoría básica.

## Fase 2 --- Automatización

-   Respuestas automáticas.
-   Reglas y palabras clave.
-   Horarios de atención.
-   Derivación a agentes.
-   Historial de ejecuciones.
-   Plantillas de WhatsApp.
-   Colas, reintentos e idempotencia.

## Fase 3 --- Primer vertical: tiendas de ropa y zapatillas

-   Catálogo.
-   Variantes de tallas y colores.
-   Inventario.
-   Pedidos.
-   Estados de pedido.
-   Consulta de productos desde WhatsApp.

## Fase 4 --- IA y métricas

-   Asistente IA con información del negocio.
-   Transferencia a humano.
-   Límites de consumo.
-   Reportes de atención y conversión.
-   Métricas de automatizaciones.

## Fase 5 --- Otros verticales

1.  Restaurantes.
2.  Ferreterías.
3.  Clínicas y consultorios.

La incorporación de clínicas debe contemplar una evaluación específica
de privacidad y protección de datos personales y de salud.

------------------------------------------------------------------------

# 8. Fuera del alcance inicial del MVP

Para evitar que el proyecto crezca demasiado antes de validar el
mercado, se recomienda dejar fuera inicialmente:

-   Aplicación móvil nativa.
-   Integración con todos los canales sociales.
-   ERP contable completo.
-   Nómina y recursos humanos.
-   Historia clínica electrónica.
-   Pasarela de pagos propia.
-   Marketplace público de productos.
-   Automatizaciones ilimitadas.
-   Agentes de IA autónomos con capacidad de ejecutar operaciones
    críticas sin aprobación.
-   Integraciones con todos los sistemas POS del mercado.

------------------------------------------------------------------------

# 9. Criterios de aceptación globales

El MVP no se considera listo para producción hasta verificar que:

1.  Un usuario solo puede acceder a los tenants de los que es miembro.
2.  Los permisos se validan en backend.
3.  Un mensaje entrante de WhatsApp se procesa una sola vez aunque el
    webhook se reintente.
4.  Un agente puede recibir, responder, asignar y cerrar conversaciones.
5.  Los mensajes fallidos pueden identificarse y gestionarse.
6.  Una automatización puede activarse, ejecutarse, pausarse y
    auditarse.
7.  Los datos de un negocio no son visibles para otro.
8.  Los pedidos y movimientos críticos mantienen consistencia.
9.  Los errores de proveedores externos no bloquean la API principal.
10. Los secretos y tokens nunca se exponen al navegador ni a logs.
11. Existen pruebas automatizadas para los flujos críticos.
12. Existe documentación de instalación, configuración, despliegue y
    recuperación.

------------------------------------------------------------------------

# 10. Decisiones pendientes antes de implementar

Estas decisiones deben cerrarse en el diseño técnico:

-   Framework definitivo del frontend.
-   Framework definitivo del backend.
-   Estrategia de base de datos multi-tenant.
-   Proveedor de hosting y almacenamiento.
-   Modelo de autenticación y sesiones.
-   Proceso de onboarding y vinculación de WhatsApp.
-   Número de cuentas y números permitidos por plan.
-   Política de retención de mensajes y archivos.
-   Proveedor de IA y presupuesto por tenant.
-   Pasarela de pagos y facturación.
-   Alcance de sucursales y almacenes.
-   Reglas de suscripción, límites y periodos de gracia.
-   Requisitos legales aplicables en Perú, especialmente protección de
    datos personales y datos de salud.

**Principio de implementación:** primero construir un núcleo
multi-tenant seguro y funcional; después incorporar automatización y un
sector comercial inicial. Los demás sectores deben reutilizar el núcleo
sin introducir dependencias innecesarias entre sus reglas de negocio.
