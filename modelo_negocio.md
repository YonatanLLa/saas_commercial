# Modelo de negocio

## SaaS multiindustria de automatización comercial por WhatsApp

**Versión:** 1.0\
**Fecha:** 10 de noviembre de 2026\
**Mercado inicial:** Perú (moneda PEN)\
**Documentos relacionados:** [requerimientos_backend.md](requerimientos_backend.md),
[requerimientos_frontend.md](requerimientos_frontend.md), [stack.md](stack.md)

> **Estado de este documento.** Los precios, límites, costos y metas son
> **hipótesis de partida**, no datos validados. Están marcados como
> *(supuesto)* y se listan en la sección 14 para validarlos con clientes
> reales antes de fijarlos en el módulo `plans`.

------------------------------------------------------------------------

# 1. Resumen

Vendemos a pequeñas y medianas empresas peruanas una suscripción mensual
a una plataforma que ordena su WhatsApp: una bandeja compartida para
todo el equipo, respuestas automáticas y la gestión de lo que nace de
cada conversación (pedidos, citas o cotizaciones), adaptada a su rubro.

- **Quién paga:** el dueño del negocio.
- **Por qué paga:** deja de perder ventas por mensajes sin responder y
  de llevar pedidos en cuadernos, Excel o un solo celular.
- **Cómo cobramos:** suscripción por plan, con límites de usuarios,
  números, conversaciones y uso de IA, más complementos.
- **Por qué nosotros:** un producto por rubro, en español, con precio en
  soles y con la API oficial de Meta, frente a herramientas genéricas en
  dólares o soluciones informales con riesgo de bloqueo del número.

------------------------------------------------------------------------

# 2. Problema

Las pymes peruanas venden por WhatsApp, pero lo operan de forma
artesanal:

1. **Un solo celular para todo.** Varios empleados comparten un teléfono
   o WhatsApp Web; nadie sabe quién respondió qué.
2. **Mensajes perdidos.** Fuera de horario o en horas pico, las
   consultas quedan sin respuesta y el cliente compra en otro lado.
3. **Trabajo repetido.** Las mismas preguntas (precio, stock, horario,
   dirección) se responden a mano decenas de veces al día.
4. **Operación fuera del chat.** El pedido, la cita o la cotización se
   anota en otro sitio, con errores y sin seguimiento.
5. **Sin datos.** El dueño no sabe cuántas conversaciones recibe, cuánto
   demora en responder ni cuántas terminan en venta.
6. **Riesgo.** Las herramientas no oficiales exponen al negocio a perder
   su número.

------------------------------------------------------------------------

# 3. Segmentos de clientes

## 3.1 Perfil común

Negocios de 2 a 30 personas, con 1 a 5 sucursales, que ya reciben por
WhatsApp un volumen que no pueden atender bien con un solo teléfono
(referencia: más de 30 conversaciones diarias *(supuesto)*).

## 3.2 Verticales y orden de entrada

El orden sigue la priorización del MVP en los requerimientos.

| Orden | Vertical | Dolor principal | Lo que resuelve el módulo |
|---|---|---|---|
| 1 | Tiendas de ropa y zapatillas | "¿Tienes talla 40 en negro?" respondido a mano; stock desactualizado | Catálogo con variantes, stock por sucursal, pedidos y seguimiento |
| 2 | Restaurantes | Pedidos por chat en hora pico; errores en modificadores | Menú, extras, disponibilidad, delivery o recojo |
| 3 | Ferreterías | Cotizaciones largas hechas a mano | Catálogo con unidades, cotizaciones con vigencia, conversión a pedido |
| 4 | Clínicas y consultorios | Inasistencias y agenda por teléfono | Agenda, citas, confirmaciones y recordatorios |

Clínicas va al final porque exige una evaluación previa de protección
de datos personales y de salud.

## 3.3 Quién decide y quién usa

| Persona | Rol en la compra | Qué le importa |
|---|---|---|
| Dueño o gerente | Decide y paga | Vender más, precio claro en soles, no perder el número |
| Administrador | Recomienda y configura | Que sea fácil de poner en marcha |
| Agente u operador | Usa a diario | Rapidez, no repetir respuestas |

------------------------------------------------------------------------

# 4. Propuesta de valor

**"Todo tu WhatsApp de ventas en un solo lugar, hecho para tu rubro."**

| Beneficio | Cómo lo entregamos |
|---|---|
| Ningún cliente sin respuesta | Bandeja compartida, asignación a agentes, respuestas automáticas y mensajes fuera de horario |
| Menos trabajo repetido | Automatizaciones por reglas y palabras clave; IA opcional con transferencia a una persona |
| La venta se cierra en el chat | Pedidos, citas o cotizaciones creados desde la conversación |
| Hecho para el rubro | Módulos verticales: tallas y colores, modificadores de menú, unidades de medida, agenda |
| Control para el dueño | Métricas de atención y conversión, auditoría de acciones |
| Número seguro | Integración oficial con WhatsApp Business Platform |
| Sin sorpresas | Precio fijo en soles, límites visibles y avisos antes de alcanzarlos |

------------------------------------------------------------------------

# 5. Modelo de ingresos

## 5.1 Fuentes de ingreso

| Fuente | Tipo | Peso esperado *(supuesto)* |
|---|---|---|
| Suscripción por plan | Recurrente mensual o anual | 80 % |
| Complementos (usuarios, números, IA) | Recurrente | 12 % |
| Puesta en marcha asistida | Pago único | 6 % |
| Plan Empresa a medida | Recurrente, contrato | 2 % al inicio |

## 5.2 Planes *(precios y límites son supuestos)*

Precios mensuales en soles, sin IGV.

| | **Emprendedor** | **Negocio** | **Pro** | **Empresa** |
|---|---|---|---|---|
| Precio mensual | S/ 99 | S/ 249 | S/ 549 | A medida |
| Usuarios incluidos | 2 | 5 | 15 | A medida |
| Números de WhatsApp | 1 | 1 | 3 | A medida |
| Conversaciones al mes | 1 000 | 5 000 | 20 000 | A medida |
| Sucursales | 1 | 2 | 10 | A medida |
| Bandeja, CRM, etiquetas, notas | Sí | Sí | Sí | Sí |
| Automatizaciones activas | 3 | 15 | Sin tope práctico | A medida |
| Módulo vertical del rubro | No | Sí | Sí | Sí |
| Respuestas de IA al mes | No | 500 | 3 000 | A medida |
| Reportes | Básicos | Completos | Completos y exportación | Completos |
| Auditoría | 30 días | 90 días | 12 meses | A medida |
| Soporte | Correo | Correo y WhatsApp | Prioritario | Gestor asignado |

**El plan Negocio es el plan objetivo:** es el primero que incluye el
módulo vertical, que es donde está la diferenciación.

## 5.3 Complementos *(supuestos)*

| Complemento | Precio |
|---|---|
| Usuario adicional | S/ 29 al mes |
| Número de WhatsApp adicional | S/ 79 al mes |
| Paquete de 1 000 respuestas de IA | S/ 49 |
| Paquete de 5 000 conversaciones | S/ 99 al mes |
| Puesta en marcha asistida (carga de catálogo, flujos iniciales, capacitación) | S/ 300, pago único |

## 5.4 Reglas comerciales *(supuestos)*

- **Prueba gratuita:** 14 días del plan Negocio, sin tarjeta.
- **Pago anual:** 2 meses de descuento (se pagan 10).
- **Sin plan gratuito permanente:** cada tenant tiene un costo real de
  infraestructura y soporte.
- **Al alcanzar un límite:** aviso al 80 % y al 100 %. No se corta la
  recepción de mensajes; se bloquea lo nuevo (más usuarios, más
  automatizaciones, más respuestas de IA) hasta subir de plan o comprar
  un paquete.
- **Falta de pago:** periodo de gracia de 7 días, luego modo de solo
  lectura, luego suspensión. Los datos se conservan 60 días.

## 5.5 Costo de mensajería de Meta

Meta cobra por los mensajes de plantilla iniciados por el negocio
(marketing, utilidad, autenticación). Las respuestas dentro de la
ventana de atención de 24 horas no tienen costo por mensaje.

**Decisión recomendada:** cada cliente vincula su propia cuenta de
WhatsApp Business y **paga a Meta directamente** con su medio de pago.

| Opción | Ventaja | Desventaja |
|---|---|---|
| **El cliente paga a Meta (recomendada)** | Sin riesgo financiero ni de cobranza; precio transparente | Un paso más en el alta |
| Nosotros pagamos y refacturamos con margen | Ingreso adicional | Adelantamos dinero, riesgo de impago y de tipo de cambio |

Las tarifas de Meta cambian por país y categoría; deben consultarse en
su lista vigente antes de comunicar cifras a clientes.

## 5.6 Relación con el producto

Los planes ya están previstos en los requerimientos:

| Regla comercial | Requerimiento |
|---|---|
| Planes, funcionalidades y límites | BR-123 |
| Suscripción asociada al tenant | BR-124 |
| Medición de conversaciones, agentes, números, automatizaciones e IA | BR-125 |
| Límites aplicados en backend | BR-126 |
| Estados de suscripción y periodo de gracia | BR-127 |
| Eventos de facturación y pagos | BR-128 |
| Plan, consumo y avisos visibles para el cliente | FR-101, FR-102, FR-104 |
| Varios números según plan | BR-059 |
| Presupuesto de IA por tenant | BR-076 |

------------------------------------------------------------------------

# 6. Estructura de costos

## 6.1 Costos variables (crecen con los clientes)

| Costo | Qué lo mueve | Cómo se controla |
|---|---|---|
| Infraestructura (cómputo, base de datos, almacenamiento) | Tenants, mensajes, archivos | Límite de conversaciones y política de retención |
| IA | Respuestas generadas | Cupo por plan y paquetes prepagados |
| Comisión de la pasarela de pagos | Monto cobrado | Incentivo al pago anual |
| Soporte | Número de clientes | Puesta en marcha guiada, plantillas por rubro, centro de ayuda |
| Correo transaccional | Invitaciones y avisos | Marginal |

## 6.2 Costos fijos

- Equipo de desarrollo y producto (el mayor costo).
- Ventas y marketing.
- Monitoreo y herramientas (Sentry, observabilidad).
- Contabilidad, facturación electrónica y asesoría legal (protección de
  datos personales).

El stack usa componentes open source con licencia permisiva
([stack.md](stack.md)), por lo que no hay costo de licencias por
usuario ni por tenant.

------------------------------------------------------------------------

# 7. Economía unitaria *(ejemplo ilustrativo, todos supuestos)*

| Concepto | Valor |
|---|---|
| Ingreso promedio por cliente al mes (ARPA) | S/ 220 |
| Costo variable por cliente al mes (infraestructura, IA, pasarela, soporte) | S/ 45 |
| Margen bruto | 80 % |
| Bajas mensuales (churn) | 4 % |
| Vida promedio del cliente | 25 meses |
| Valor de vida (LTV) = 220 × 80 % × 25 | S/ 4 400 |
| Costo de adquisición objetivo (CAC) | S/ 600 o menos |
| LTV / CAC | 7 o más |
| Recuperación del CAC | 3 a 4 meses |

**Punto de equilibrio:** con costos fijos de S/ 40 000 al mes y S/ 175
de margen por cliente, se necesitan cerca de **230 clientes activos**.
Con S/ 25 000 de costos fijos, cerca de **145**.

------------------------------------------------------------------------

# 8. Canales y estrategia de salida al mercado

## 8.1 Etapas

| Etapa | Meta | Canal principal |
|---|---|---|
| 1. Validación | 10 tiendas de ropa o zapatillas usando el producto, 5 pagando | Venta directa del fundador: visitas a galerías y emporios (por ejemplo, Gamarra), contactos propios |
| 2. Primeros 100 clientes | Repetir la venta en un solo vertical | Anuncios de Meta con botón a WhatsApp, TikTok e Instagram con casos reales, referidos |
| 3. Expansión | Nuevos verticales y ciudades | Alianzas y contenido por rubro |

## 8.2 Canales

- **Venta directa y demostración por WhatsApp.** Usamos nuestro propio
  producto para vender: es la mejor demostración.
- **Anuncios de clic a WhatsApp** segmentados por rubro.
- **Referidos:** un mes gratis para quien refiere y para el referido.
- **Aliados:** contadores, agencias de marketing, proveedores de punto
  de venta y de facturación electrónica, gremios y cámaras de comercio.
  Comisión recurrente del 20 % durante 12 meses *(supuesto)*.
- **Contenido:** guías cortas por rubro ("cómo tomar pedidos por
  WhatsApp sin perder ninguno").

## 8.3 Activación

El cliente percibe valor cuando, en su primera semana:

1. Vincula su número.
2. Invita al menos a un compañero.
3. Activa una respuesta automática.
4. Registra su primer pedido, cita o cotización desde un chat.

Estos cuatro pasos son la medida de activación.

------------------------------------------------------------------------

# 9. Relación con el cliente

- **Alta guiada** dentro de la aplicación, con plantillas de catálogo y
  de automatizaciones por rubro.
- **Puesta en marcha asistida** opcional y de pago.
- **Soporte por WhatsApp** en horario comercial; prioritario en Pro.
- **Seguimiento proactivo:** contacto cuando un cliente deja de usar la
  bandeja o se acerca a un límite.
- **Retención:** cuanto más historial, catálogo y automatizaciones tiene
  el cliente dentro de la plataforma, más le cuesta irse.

------------------------------------------------------------------------

# 10. Recursos, actividades y socios clave

| | |
|---|---|
| **Recursos clave** | Plataforma multi-tenant; cuenta de proveedor tecnológico aprobada por Meta; conocimiento de cada rubro; plantillas por vertical |
| **Actividades clave** | Desarrollo del producto; alta y acompañamiento de clientes; mantener el cumplimiento de las políticas de Meta; seguridad y aislamiento de datos |
| **Socios clave** | Meta (WhatsApp Business Platform); pasarela de pagos con presencia en Perú (Culqi, Mercado Pago o Niubiz); proveedor de facturación electrónica autorizado por SUNAT; proveedor de IA; proveedor de hosting |

------------------------------------------------------------------------

# 11. Competencia y diferenciación

| Alternativa | Fortaleza | Debilidad frente a nosotros |
|---|---|---|
| WhatsApp Business (aplicación gratuita) | Gratis y conocida | Pocos dispositivos, sin asignación, sin pedidos ni reportes |
| Bandejas multicanal internacionales | Producto maduro | Precio en dólares por agente, genéricas, sin lógica de rubro |
| Constructores de chatbots | Automatización potente | Requieren configuración técnica; no gestionan la operación |
| Herramientas no oficiales | Baratas | Riesgo de bloqueo del número |
| Sistemas de punto de venta o ERP | Gestionan inventario y ventas | WhatsApp no es su canal central |

**Nuestra posición:** conversación y operación en un mismo lugar, por
rubro, en soles y con la API oficial.

**Ventaja defendible con el tiempo:** plantillas y flujos por vertical,
datos históricos del cliente dentro de la plataforma y red de aliados
locales.

------------------------------------------------------------------------

# 12. Métricas del negocio

| Área | Métrica | Meta inicial *(supuesto)* |
|---|---|---|
| Ingresos | Ingreso recurrente mensual (MRR) | Crecimiento mensual de 15 % en el primer año |
| Ingresos | Ingreso promedio por cliente | S/ 220 o más |
| Adquisición | Conversión de prueba a pago | 20 % o más |
| Adquisición | Costo de adquisición | S/ 600 o menos |
| Activación | Clientes que completan los 4 pasos en 7 días | 60 % o más |
| Retención | Bajas mensuales de clientes | 4 % o menos |
| Retención | Retención neta de ingresos | 100 % o más |
| Rentabilidad | Margen bruto | 75 % o más |
| Uso | Clientes con al menos una automatización activa | 70 % o más |

------------------------------------------------------------------------

# 13. Riesgos

| Riesgo | Impacto | Mitigación |
|---|---|---|
| Dependencia de Meta (precios, políticas, aprobación) | Alto | Usar solo la API oficial; el cliente paga a Meta; la arquitectura ya separa canales para sumar otros |
| El cliente no quiere pagar por algo que hoy usa gratis | Alto | Prueba gratuita, vender ahorro de tiempo y ventas recuperadas, plan de entrada accesible |
| Vinculación del número compleja para una pyme | Medio | Alta guiada y puesta en marcha asistida |
| Abarcar cuatro rubros a la vez | Alto | Un vertical por vez; no abrir el siguiente hasta validar el anterior |
| Costo de IA mayor al previsto | Medio | Cupos por plan, paquetes prepagados, reglas deterministas como base |
| Fuga de datos entre negocios | Crítico | Aislamiento por tenant con RLS y pruebas automatizadas |
| Datos de salud en el vertical de clínicas | Alto | Revisión legal previa; sin historia clínica en el MVP |
| Bajas altas en pymes | Medio | Seguimiento proactivo, pago anual con descuento |

------------------------------------------------------------------------

# 14. Supuestos por validar

| # | Supuesto | Cómo validarlo |
|---|---|---|
| 1 | Una tienda de ropa paga S/ 249 al mes por bandeja, automatización y pedidos | 15 entrevistas y preventa a 5 negocios |
| 2 | Los límites de conversaciones por plan corresponden al volumen real | Medir el consumo de los primeros 10 clientes |
| 3 | El dueño acepta pagar a Meta por separado | Probar el alta con clientes piloto |
| 4 | 14 días bastan para activar al cliente | Medir el tiempo hasta el primer pedido registrado |
| 5 | El costo variable por cliente es cercano a S/ 45 | Medir infraestructura e IA en producción |
| 6 | El CAC se mantiene por debajo de S/ 600 | Campaña piloto con presupuesto acotado |
| 7 | La IA justifica un precio mayor | Ofrecerla como complemento y medir adopción |

------------------------------------------------------------------------

# 15. Hoja de ruta comercial

Alineada con las fases de [requerimientos_backend.md](requerimientos_backend.md).

| Fase del producto | Qué se puede vender | Acción comercial |
|---|---|---|
| 1. Núcleo SaaS | Bandeja compartida y CRM | Pilotos gratuitos con 5 a 10 tiendas |
| 2. Automatización | Plan Emprendedor | Primeros cobros; validar precio |
| 3. Vertical tiendas | Plan Negocio | Lanzamiento de pago en un rubro; cobro en línea y facturación electrónica |
| 4. IA y métricas | Plan Pro y paquetes de IA | Subida de plan de clientes existentes |
| 5. Otros verticales | Restaurantes, ferreterías, clínicas | Repetir la estrategia por rubro; programa de aliados |

------------------------------------------------------------------------

# 16. Decisiones pendientes

Corresponden a las decisiones abiertas en los requerimientos y en
[stack.md](stack.md):

- Precios y límites definitivos por plan.
- Pasarela de pagos y proveedor de facturación electrónica.
- Quién asume el costo de mensajería de Meta (recomendado: el cliente).
- Proveedor de IA y presupuesto por tenant.
- Duración de la prueba, periodo de gracia y retención de datos tras la baja.
- Política de retención de mensajes y archivos por plan.
- Revisión legal de protección de datos personales en Perú.
