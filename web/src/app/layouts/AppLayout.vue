<script setup lang="ts">
import type { Component } from 'vue'
import {
  BarChart3,
  Bot,
  LayoutDashboard,
  LogOut,
  Menu,
  MessageCircle,
  Package,
  Settings,
  Users,
  X,
} from '@lucide/vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store'
import { Button } from '@/shared/ui/button'

interface NavItem {
  to: string
  label: string
  icon: Component
}

// PENDIENTE: filtrar por permisos, rubro y módulos del plan (FR-004).
const nav: NavItem[] = [
  { to: '/', label: 'Panel', icon: LayoutDashboard },
  { to: '/inbox', label: 'Conversaciones', icon: MessageCircle },
  { to: '/contacts', label: 'Contactos', icon: Users },
  { to: '/catalog', label: 'Catálogo', icon: Package },
  { to: '/automations', label: 'Automatizaciones', icon: Bot },
  { to: '/reports', label: 'Reportes', icon: BarChart3 },
  { to: '/settings', label: 'Configuración', icon: Settings },
]

const auth = useAuthStore()
const router = useRouter()
const mobileOpen = ref(false)

async function logout() {
  auth.logout()
  await router.replace({ name: 'login' })
}
</script>

<template>
  <div class="flex min-h-svh">
    <!-- Fondo oscuro del menú en móvil -->
    <div
      v-if="mobileOpen"
      class="fixed inset-0 z-30 bg-black/40 md:hidden"
      aria-hidden="true"
      @click="mobileOpen = false"
    />

    <aside
      :class="[
        'fixed inset-y-0 left-0 z-40 flex w-64 flex-col border-r bg-background transition-transform md:static md:translate-x-0',
        mobileOpen ? 'translate-x-0' : '-translate-x-full',
      ]"
      aria-label="Navegación principal"
    >
      <div class="flex h-14 items-center justify-between border-b px-4">
        <span class="font-semibold">{{ auth.tenant?.name ?? 'SaaS WhatsApp' }}</span>
        <Button variant="ghost" size="icon" class="md:hidden" aria-label="Cerrar menú" @click="mobileOpen = false">
          <X />
        </Button>
      </div>
      <nav class="flex-1 space-y-1 p-2">
        <RouterLink
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 rounded-md px-3 py-2 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground"
          exact-active-class="bg-accent text-accent-foreground font-medium"
          @click="mobileOpen = false"
        >
          <component :is="item.icon" class="size-4" aria-hidden="true" />
          {{ item.label }}
        </RouterLink>
      </nav>
      <div class="border-t p-2">
        <Button variant="ghost" class="w-full justify-start" @click="logout">
          <LogOut /> Cerrar sesión
        </Button>
      </div>
    </aside>

    <div class="flex min-w-0 flex-1 flex-col">
      <header class="flex h-14 items-center gap-2 border-b px-4">
        <Button variant="ghost" size="icon" class="md:hidden" aria-label="Abrir menú" @click="mobileOpen = true">
          <Menu />
        </Button>
        <h1 class="text-base font-semibold">{{ $route.meta.title }}</h1>
        <span class="ml-auto text-sm text-muted-foreground">{{ auth.user?.fullName }}</span>
      </header>
      <main class="flex-1 p-4 md:p-6">
        <RouterView />
      </main>
    </div>
  </div>
</template>
