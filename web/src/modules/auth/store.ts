import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { setAccessToken } from '@/shared/api/client'

export interface SessionUser {
  id: string
  fullName: string
  email: string
}

export interface TenantContext {
  id: string
  name: string
  industry: 'retail' | 'restaurant' | 'clinic' | 'hardware'
  role: 'owner' | 'admin' | 'supervisor' | 'agent' | 'operator'
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<SessionUser | null>(null)
  const tenant = ref<TenantContext | null>(null)

  const isAuthenticated = computed(() => user.value !== null)

  async function login(email: string, _password: string) {
    // PENDIENTE: llamar a POST /api/v1/auth/login cuando exista en openapi.yaml.
    // Mientras tanto, solo en desarrollo se crea una sesión de demostración.
    if (!import.meta.env.DEV) {
      throw new Error('El inicio de sesión todavía no está implementado.')
    }
    user.value = { id: 'demo', fullName: 'Usuario demo', email }
    tenant.value = { id: 'demo', name: 'Negocio demo', industry: 'retail', role: 'owner' }
  }

  function logout() {
    setAccessToken(null)
    user.value = null
    tenant.value = null
  }

  return { user, tenant, isAuthenticated, login, logout }
})
