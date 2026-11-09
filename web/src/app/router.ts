import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store'

declare module 'vue-router' {
  interface RouteMeta {
    public?: boolean
    title?: string
  }
}

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/modules/auth/pages/LoginPage.vue'),
      meta: { public: true, title: 'Iniciar sesión' },
    },
    {
      path: '/',
      component: () => import('@/app/layouts/AppLayout.vue'),
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('@/modules/dashboard/pages/DashboardPage.vue'),
          meta: { title: 'Panel' },
        },
        // Módulos aún no construidos: comparten una página "en construcción".
        ...[
          ['inbox', 'Conversaciones'],
          ['contacts', 'Contactos'],
          ['catalog', 'Catálogo'],
          ['automations', 'Automatizaciones'],
          ['reports', 'Reportes'],
          ['settings', 'Configuración'],
        ].map(([path, title]) => ({
          path: path!,
          name: path!,
          component: () => import('@/shared/components/ComingSoonPage.vue'),
          meta: { title },
        })),
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.isAuthenticated) {
    return { name: 'dashboard' }
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} · SaaS WhatsApp` : 'SaaS WhatsApp'
})
