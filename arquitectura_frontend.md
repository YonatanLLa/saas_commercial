# Arquitectura del Frontend

## SaaS multiindustria de automatización comercial por WhatsApp

Este documento define la organización de archivos y las convenciones de
`web/`. Se basa en `requerimientos_frontend.md` (FR / RNF) y `stack.md`.

------------------------------------------------------------------------

# 1. Estilo de arquitectura

SPA en **Vue 3 + TypeScript + Vite** con una **arquitectura modular por
dominio** (feature-based / vertical slices), organizada en tres capas:

```
app/  ──►  modules/  ──►  shared/
```

| Capa | Responsabilidad | Puede importar de |
|---|---|---|
| `app/` | Arranque, router, guards, plugins, layouts, navegación | `modules/`, `shared/` |
| `modules/` | Un directorio por dominio de negocio, autocontenido | `shared/` y la API pública de otros módulos (componentes o composables exportados) |
| `shared/` | Piezas reutilizables sin lógica de dominio | Solo `shared/` |

**Reglas:**

- `shared/` nunca importa de `modules/` ni de `app/`.
- Un módulo no accede a la store, la API ni los componentes internos de
  otro módulo. Si necesita algo, el módulo dueño lo expone (por ejemplo,
  `catalog/components/ProductPicker.vue` es usado por `inbox`, `orders` y
  `quotations`).
- Los módulos verticales (`retail`, `restaurant`, `clinic`, `hardware`)
  extienden a los módulos comunes; los comunes no dependen de los
  verticales.

------------------------------------------------------------------------

# 2. Stack

| Necesidad | Elección |
|---|---|
| Framework | Vue 3 + TypeScript |
| Build | Vite, pnpm |
| UI | Tailwind CSS v4 + shadcn-vue (reka-ui), iconos Lucide |
| Notificaciones | vue-sonner |
| Ruteo | Vue Router con guards por permiso, rubro y plan |
| Estado de cliente | Pinia (sesión, tenant, sucursal, preferencias) |
| Estado del servidor | TanStack Query (`@tanstack/vue-query`) |
| Formularios | VeeValidate + Zod |
| Cliente HTTP | `openapi-fetch` + tipos de `openapi-typescript` |
| Tiempo real | WebSocket nativo con reconexión y resincronización |
| Editor de automatizaciones | Vue Flow |
| Gráficos | ECharts (`vue-echarts`) |
| Fechas | date-fns + date-fns-tz |
| i18n | vue-i18n (es-PE por defecto) |
| Utilidades | VueUse |
| Pruebas | Vitest + Vue Testing Library; Playwright para E2E |

------------------------------------------------------------------------

# 3. Convención interna de un módulo

Todos los módulos siguen el mismo molde (RNF-006). Solo se crean los
archivos que el módulo realmente necesita.

```
modules/<dominio>/
├── routes.ts        # Rutas lazy + meta { permission, industries, feature }
├── api.ts           # Llamadas openapi-fetch + query keys
├── queries.ts       # useXxxQuery / useXxxMutation (TanStack Query)
├── store.ts         # Pinia, solo si hay estado de cliente
├── schemas.ts       # Esquemas Zod para formularios
├── types.ts         # Tipos propios del dominio
├── pages/           # Vistas enrutadas (sufijo Page.vue)
├── components/      # Componentes del dominio
└── composables/     # Lógica reutilizable dentro del módulo
```

**Nombres:**

- Páginas: `XxxPage.vue`. Diálogos: `XxxDialog.vue`. Formularios:
  `XxxForm.vue`.
- Composables: `useXxx.ts`.
- Pruebas junto al archivo probado: `xxx.test.ts`.

------------------------------------------------------------------------

# 4. Distribución de archivos

Los archivos marcados con **✓** ya existen.

```
web/
├── index.html ✓
├── package.json ✓
├── vite.config.ts ✓
├── components.json ✓
├── tsconfig*.json ✓
├── e2e/                                   # Playwright
│   ├── auth.spec.ts
│   ├── inbox.spec.ts
│   └── tenant-isolation.spec.ts           # RNF-012
└── src/
    ├── main.ts ✓
    ├── App.vue ✓
    ├── style.css ✓
    ├── env.d.ts ✓
    │
    ├── app/                               # ── CAPA DE APLICACIÓN ──
    │   ├── router/
    │   │   ├── index.ts                   # createRouter (hoy: app/router.ts ✓)
    │   │   ├── routes.ts                  # Junta los routes.ts de cada módulo
    │   │   └── guards.ts                  # Auth, permiso, rubro, plan (FR-004)
    │   ├── plugins/
    │   │   ├── pinia.ts
    │   │   ├── vue-query.ts
    │   │   └── i18n.ts
    │   ├── navigation/
    │   │   └── menu.ts                    # Ítems del menú con permiso/rubro/feature
    │   ├── layouts/
    │   │   ├── AppLayout.vue ✓
    │   │   └── AuthLayout.vue             # Login, recuperación de contraseña
    │   └── components/
    │       ├── AppSidebar.vue
    │       ├── AppTopbar.vue
    │       ├── BranchSwitcher.vue         # FR-005
    │       ├── UserMenu.vue
    │       ├── RealtimeStatusBanner.vue   # FR-006 (estado de conexión)
    │       └── PlanLimitBanner.vue        # FR-102
    │
    ├── modules/                           # ── DOMINIOS ──
    │   │
    │   ├── auth/                          # FR-002, FR-007, FR-009 a FR-011, FR-014
    │   │   ├── routes.ts
    │   │   ├── api.ts
    │   │   ├── store.ts ✓                 # Sesión, tenant y sucursal activos
    │   │   ├── schemas.ts
    │   │   ├── composables/useSessionRefresh.ts
    │   │   └── pages/
    │   │       ├── LoginPage.vue ✓
    │   │       ├── ForgotPasswordPage.vue
    │   │       ├── ResetPasswordPage.vue
    │   │       └── AccountStatusPage.vue  # Pendiente / bloqueada / suspendida
    │   │
    │   ├── account/                       # FR-012, FR-013
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/ProfilePage.vue
    │   │   └── components/
    │   │       ├── ProfileForm.vue
    │   │       └── ActiveSessionsList.vue
    │   │
    │   ├── dashboard/                     # FR-015 a FR-019
    │   │   ├── routes.ts · api.ts · queries.ts
    │   │   ├── pages/DashboardPage.vue ✓
    │   │   └── components/
    │   │       ├── KpiCard.vue
    │   │       ├── DashboardFilters.vue   # Fecha, sucursal, canal, agente
    │   │       ├── TrendChart.vue         # ECharts
    │   │       ├── AttentionQueue.vue     # FR-018
    │   │       └── DataQualityHint.vue    # FR-019 (dato real / incompleto / estimado)
    │   │
    │   ├── inbox/                         # FR-020 a FR-035, FR-062
    │   │   ├── routes.ts · api.ts · queries.ts · store.ts · types.ts
    │   │   ├── composables/
    │   │   │   ├── useConversationRealtime.ts  # FR-031
    │   │   │   ├── useConversationLock.ts      # FR-035 (conflictos entre agentes)
    │   │   │   └── useMessageComposer.ts
    │   │   ├── pages/InboxPage.vue        # Lista + chat + panel de contacto
    │   │   └── components/
    │   │       ├── ConversationList.vue
    │   │       ├── ConversationListItem.vue
    │   │       ├── ConversationFilters.vue      # FR-021
    │   │       ├── ConversationHeader.vue       # Estado, asignación, transferencia
    │   │       ├── MessageThread.vue            # FR-023
    │   │       ├── MessageBubble.vue
    │   │       ├── MessageStatusIcon.vue        # FR-025
    │   │       ├── MessageOriginBadge.vue       # FR-062 (IA / automatización / agente)
    │   │       ├── MessageComposer.vue          # FR-024
    │   │       ├── TemplatePickerDialog.vue     # FR-032 (fuera de la ventana de 24 h)
    │   │       ├── ProductShareDialog.vue       # FR-050 (usa catalog)
    │   │       ├── InternalNoteForm.vue         # FR-029
    │   │       ├── TagEditor.vue                # FR-030
    │   │       ├── AssignAgentDialog.vue        # FR-026
    │   │       ├── TransferDialog.vue           # FR-027, FR-034
    │   │       └── ContactSidePanel.vue         # FR-022, FR-033
    │   │
    │   ├── contacts/                      # FR-036 a FR-044
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts · types.ts
    │   │   ├── pages/
    │   │   │   ├── ContactsListPage.vue
    │   │   │   ├── ContactDetailPage.vue
    │   │   │   └── DuplicatesReviewPage.vue     # FR-043
    │   │   └── components/
    │   │       ├── ContactForm.vue
    │   │       ├── ContactTimeline.vue          # FR-042
    │   │       ├── ContactRelations.vue         # Pedidos, citas, cotizaciones (FR-040)
    │   │       ├── ContactSegments.vue          # FR-041
    │   │       └── CommunicationPreferences.vue # FR-044 (opt-in / opt-out)
    │   │
    │   ├── catalog/                       # FR-045 a FR-050 (base común a todos los rubros)
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts · types.ts
    │   │   ├── pages/
    │   │   │   ├── ProductsListPage.vue
    │   │   │   ├── ProductEditPage.vue
    │   │   │   └── CategoriesPage.vue           # FR-047
    │   │   └── components/
    │   │       ├── ProductForm.vue
    │   │       ├── ProductImageUploader.vue
    │   │       ├── CategoryTree.vue
    │   │       ├── VariantEditor.vue            # Punto de extensión por rubro (FR-048)
    │   │       └── ProductPicker.vue            # Reutilizado por inbox, orders, quotations
    │   │
    │   ├── inventory/                     # FR-049, FR-076, FR-093
    │   │   ├── routes.ts · api.ts · queries.ts
    │   │   ├── pages/InventoryPage.vue
    │   │   └── components/
    │   │       ├── StockTable.vue               # Por variante / sucursal / almacén
    │   │       └── StockAdjustDialog.vue
    │   │
    │   ├── orders/                        # FR-077 a FR-079, FR-084 (tienda, restaurante, ferretería)
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts · types.ts
    │   │   ├── pages/
    │   │   │   ├── OrdersListPage.vue
    │   │   │   ├── OrderDetailPage.vue
    │   │   │   └── OrderCreatePage.vue
    │   │   └── components/
    │   │       ├── OrderForm.vue
    │   │       ├── OrderItemsTable.vue
    │   │       ├── OrderStatusStepper.vue       # FR-078
    │   │       └── DeliveryInfoForm.vue         # FR-079, FR-083
    │   │
    │   ├── quotations/                    # FR-094 a FR-096 (ferretería; opcional en otros rubros)
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/
    │   │   │   ├── QuotationsListPage.vue
    │   │   │   └── QuotationEditPage.vue
    │   │   └── components/
    │   │       ├── QuotationForm.vue
    │   │       ├── QuotationTotals.vue          # Descuentos, impuestos, vigencia
    │   │       └── ConvertToOrderDialog.vue     # FR-095
    │   │
    │   ├── automations/                   # FR-051 a FR-059
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts · types.ts
    │   │   ├── pages/
    │   │   │   ├── AutomationsListPage.vue
    │   │   │   ├── AutomationEditorPage.vue     # Vue Flow (FR-053)
    │   │   │   ├── AutomationRunsPage.vue       # FR-059
    │   │   │   ├── FaqResponsesPage.vue         # FR-054
    │   │   │   ├── BusinessHoursPage.vue        # FR-055
    │   │   │   └── RoutingRulesPage.vue         # FR-056
    │   │   └── components/
    │   │       ├── flow/
    │   │       │   ├── TriggerNode.vue
    │   │       │   ├── ConditionNode.vue
    │   │       │   ├── ActionNode.vue
    │   │       │   └── NodeInspector.vue
    │   │       ├── WelcomeClosingMessages.vue   # FR-057
    │   │       └── AutomationTestSandbox.vue    # FR-058
    │   │
    │   ├── ai-assistant/                  # FR-060, FR-061
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/AiAssistantSettingsPage.vue
    │   │   └── components/
    │   │       ├── AiInstructionsForm.vue
    │   │       ├── AiKnowledgeSources.vue
    │   │       └── AiHandoffRules.vue
    │   │
    │   ├── users/                         # FR-063 a FR-068
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/
    │   │   │   ├── UsersPage.vue
    │   │   │   ├── RolesPage.vue
    │   │   │   └── TeamsPage.vue
    │   │   └── components/
    │   │       ├── InviteUserDialog.vue
    │   │       ├── RolePermissionsMatrix.vue
    │   │       ├── TeamMembersEditor.vue
    │   │       └── AgentAvailabilityBadge.vue   # FR-068
    │   │
    │   ├── whatsapp/                      # FR-069 a FR-074
    │   │   ├── routes.ts · api.ts · queries.ts
    │   │   ├── composables/useMetaEmbeddedSignup.ts  # El SDK de Meta se carga solo aquí
    │   │   ├── pages/
    │   │   │   ├── WhatsAppIntegrationPage.vue
    │   │   │   └── TemplatesPage.vue            # FR-072
    │   │   └── components/
    │   │       ├── NumberConnectionCard.vue     # FR-069, FR-071
    │   │       ├── ConnectWizard.vue            # FR-070
    │   │       ├── IntegrationErrorHelp.vue     # FR-073
    │   │       └── TemplateForm.vue
    │   │
    │   ├── reports/                       # FR-097, FR-098
    │   │   ├── routes.ts · api.ts · queries.ts
    │   │   ├── pages/
    │   │   │   ├── ReportsOverviewPage.vue
    │   │   │   ├── AgentsReportPage.vue
    │   │   │   ├── ConversionReportPage.vue
    │   │   │   └── AutomationsReportPage.vue
    │   │   └── components/ExportButton.vue
    │   │
    │   ├── audit/                         # FR-099, FR-100 (solo lectura)
    │   │   ├── routes.ts · api.ts · queries.ts
    │   │   ├── pages/AuditLogPage.vue
    │   │   └── components/AuditDiffViewer.vue
    │   │
    │   ├── settings/                      # FR-101 a FR-105
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/
    │   │   │   ├── BusinessSettingsPage.vue     # Datos, zona horaria, moneda (FR-103)
    │   │   │   ├── BranchesPage.vue
    │   │   │   ├── SubscriptionPage.vue         # Plan, límites, consumo (FR-101)
    │   │   │   ├── ModulesPage.vue              # Incluidos / no incluidos (FR-104)
    │   │   │   └── NotificationsPage.vue        # FR-105
    │   │   └── components/
    │   │       ├── UsageMeter.vue
    │   │       └── ModuleCard.vue
    │   │
    │   │   ── VERTICALES (se registran solo si el tenant tiene el rubro) ──
    │   │
    │   ├── retail/                        # FR-075 a FR-079
    │   │   ├── routes.ts · api.ts · queries.ts
    │   │   └── components/
    │   │       ├── SizeColorVariantEditor.vue   # Se inyecta en catalog/VariantEditor
    │   │       └── VariantStockGrid.vue         # Talla × color × sucursal
    │   │
    │   ├── restaurant/                    # FR-080 a FR-085
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/
    │   │   │   ├── MenuPage.vue
    │   │   │   ├── ModifiersPage.vue            # FR-081
    │   │   │   ├── MenuAvailabilityPage.vue     # FR-082
    │   │   │   └── KitchenBoardPage.vue         # Pedidos en curso (FR-084)
    │   │   └── components/
    │   │       ├── ModifierGroupEditor.vue
    │   │       ├── ServiceModeSelect.vue        # Delivery / recojo / en local (FR-083)
    │   │       └── DietaryNotesBadge.vue        # FR-085
    │   │
    │   ├── clinic/                        # FR-086 a FR-091
    │   │   ├── routes.ts · api.ts · queries.ts · schemas.ts
    │   │   ├── pages/
    │   │   │   ├── AgendaPage.vue               # FR-089
    │   │   │   ├── AppointmentDetailPage.vue
    │   │   │   ├── ProfessionalsPage.vue
    │   │   │   ├── ServicesPage.vue
    │   │   │   └── RemindersSettingsPage.vue    # FR-090
    │   │   └── components/
    │   │       ├── AppointmentCalendar.vue
    │   │       ├── AvailabilitySlots.vue        # FR-087
    │   │       ├── AppointmentForm.vue          # FR-088
    │   │       └── ClinicalInfoGate.vue         # FR-091 (oculta datos sensibles)
    │   │
    │   └── hardware/                      # FR-092, FR-093
    │       ├── routes.ts · api.ts · queries.ts
    │       └── components/
    │           ├── UnitOfMeasureEditor.vue
    │           └── WarehouseStockTable.vue
    │
    ├── shared/                            # ── REUTILIZABLE, SIN DOMINIO ──
    │   ├── api/
    │   │   ├── client.ts ✓                # openapi-fetch + Bearer
    │   │   ├── schema.d.ts ✓              # Generado (pnpm api:types)
    │   │   ├── refresh.ts                 # Reintenta tras un 401 con /auth/refresh (FR-007)
    │   │   └── errors.ts                  # Normaliza errores de la API para la UI
    │   ├── realtime/
    │   │   ├── socket.ts                  # WebSocket con backoff y resincronización
    │   │   └── events.ts                  # Tipos de eventos (mensaje, conversación…)
    │   ├── composables/
    │   │   ├── useRealtime.ts
    │   │   ├── usePermissions.ts          # can('contacts.edit')
    │   │   ├── useTenant.ts               # Rubro, plan, módulos habilitados
    │   │   ├── useBranch.ts
    │   │   ├── useFormatters.ts           # Fecha, moneda y zona horaria del tenant (RNF-011)
    │   │   ├── usePaginatedQuery.ts       # RNF-001
    │   │   ├── useUrlFilters.ts           # Filtros sincronizados con la query string
    │   │   └── useConfirm.ts
    │   ├── components/
    │   │   ├── data-table/
    │   │   │   ├── DataTable.vue
    │   │   │   ├── DataTableToolbar.vue
    │   │   │   └── DataTablePagination.vue
    │   │   ├── states/                    # FR-006
    │   │   │   ├── LoadingState.vue
    │   │   │   ├── EmptyState.vue
    │   │   │   └── ErrorState.vue
    │   │   ├── gates/
    │   │   │   ├── PermissionGate.vue     # Oculta según permiso (RNF-009)
    │   │   │   └── FeatureGate.vue        # Oculta según plan o rubro (FR-104)
    │   │   ├── PageHeader.vue
    │   │   ├── FilterBar.vue
    │   │   ├── DateRangePicker.vue
    │   │   ├── StatusBadge.vue
    │   │   ├── ConfirmDialog.vue
    │   │   ├── UpgradePrompt.vue
    │   │   └── ComingSoonPage.vue ✓
    │   ├── ui/                            # shadcn-vue
    │   │   ├── button/ ✓  card/ ✓  input/ ✓  label/ ✓
    │   │   ├── dialog/  sheet/  dropdown-menu/  select/  command/  popover/
    │   │   ├── table/  tabs/  badge/  avatar/  tooltip/  skeleton/
    │   │   ├── form/  checkbox/  switch/  textarea/  calendar/
    │   │   └── sonner/  scroll-area/  separator/  pagination/
    │   ├── constants/
    │   │   ├── permissions.ts             # Catálogo de permisos (RBAC)
    │   │   ├── roles.ts                   # owner, admin, supervisor, agent, operator
    │   │   ├── industries.ts              # retail, restaurant, clinic, hardware
    │   │   └── module-matrix.ts           # Tabla módulo × rubro (sección 4 de requerimientos)
    │   ├── lib/
    │   │   ├── utils.ts ✓  utils.test.ts ✓
    │   │   ├── format.ts                  # date-fns-tz + Intl
    │   │   └── zod.ts                     # Validaciones comunes (teléfono E.164, RUC…)
    │   └── types/
    │       └── common.ts                  # Paginated<T>, ID, etc.
    │
    └── i18n/
        ├── index.ts
        └── locales/
            └── es-PE.json
```

------------------------------------------------------------------------

# 5. Flujos transversales

## 5.1 Visibilidad por rol, rubro y plan (FR-004, FR-104)

Cada `routes.ts` declara en `meta`:

```ts
meta: {
  title: 'Pedidos',
  permission: 'orders.read',
  industries: ['retail', 'restaurant', 'hardware'],
  feature: 'orders',
}
```

- `app/router/guards.ts` bloquea la navegación si falta el permiso, el
  rubro no corresponde o el plan no incluye el módulo.
- `app/navigation/menu.ts` usa los mismos datos para construir el menú.
- `shared/components/gates/` oculta acciones puntuales dentro de una vista.
- La interfaz nunca reemplaza la validación del backend (RNF-012).

## 5.2 Sesión (FR-007, RNF-009)

- El access token vive solo en memoria (`shared/api/client.ts`).
- El refresh token viaja en una cookie `httpOnly`, `Secure`,
  `SameSite=Strict`.
- Ante un 401, `shared/api/refresh.ts` llama a `/api/v1/auth/refresh` y
  reintenta una vez; si falla, cierra la sesión y redirige al login.
- No se guardan datos sensibles en `localStorage` ni `sessionStorage`.

## 5.3 Datos del servidor

- Toda lectura y escritura pasa por `api.ts` + `queries.ts` del módulo.
- Las query keys incluyen el `tenantId` y la sucursal activa, para que un
  cambio de tenant o sucursal nunca muestre datos cacheados de otro
  (RNF-012).
- Los listados grandes usan paginación del servidor (RNF-001).

## 5.4 Tiempo real (FR-031)

- `shared/realtime/socket.ts` mantiene una única conexión por sesión,
  con reconexión por backoff exponencial.
- Al reconectar, invalida las queries afectadas para resincronizar.
- Los módulos se suscriben con `useRealtime()`; por ejemplo,
  `inbox/composables/useConversationRealtime.ts` actualiza la caché de
  conversaciones y mensajes.

## 5.5 Módulos verticales

- Los módulos comunes (`catalog`, `inventory`, `orders`, `quotations`)
  contienen la lógica compartida por varios rubros.
- Los verticales (`retail`, `restaurant`, `clinic`, `hardware`) solo
  agregan lo específico y se conectan a puntos de extensión de los
  comunes (por ejemplo, `catalog/components/VariantEditor.vue` usa
  `retail/components/SizeColorVariantEditor.vue` cuando el rubro es
  retail).
- Sus rutas se cargan de forma diferida y solo se registran si el tenant
  tiene el rubro habilitado.

------------------------------------------------------------------------

# 6. Módulos por fase del MVP

| Fase | Alcance | Módulos |
|---|---|---|
| 1 | Núcleo SaaS | `auth`, `account`, `settings`, `users`, `whatsapp`, `inbox`, `contacts`, `audit` |
| 2 | Automatización | `automations` (+ plantillas en `whatsapp`) |
| 3 | Tiendas de ropa y zapatillas | `catalog`, `inventory`, `orders`, `retail` |
| 4 | IA y métricas | `ai-assistant`, `reports`, `dashboard` completo |
| 5 | Otros verticales | `restaurant`, `hardware`, `quotations`, `clinic` |

------------------------------------------------------------------------

# 7. Cómo agregar un módulo nuevo

1. Crear `src/modules/<dominio>/` con el molde de la sección 3.
2. Definir los endpoints en `api/openapi.yaml` y regenerar los tipos con
   `pnpm api:types`.
3. Escribir `api.ts` y `queries.ts` usando el cliente tipado.
4. Declarar las rutas en `routes.ts` con `permission`, `industries` y
   `feature`, y registrarlas en `app/router/routes.ts`.
5. Agregar el ítem de menú en `app/navigation/menu.ts`.
6. Agregar los textos en `i18n/locales/es-PE.json`.
7. Escribir pruebas unitarias junto a los archivos y, si es un flujo
   crítico, una prueba en `web/e2e/`.
