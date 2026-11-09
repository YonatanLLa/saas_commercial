<script setup lang="ts">
import { MessageCircle } from '@lucide/vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/modules/auth/store'
import { Button } from '@/shared/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const email = ref('')
const password = ref('')
const error = ref<string | null>(null)
const loading = ref(false)

async function onSubmit() {
  error.value = null
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    // Solo rutas internas, para evitar redirecciones abiertas.
    await router.replace(redirect.startsWith('/') && !redirect.startsWith('//') ? redirect : '/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'No se pudo iniciar sesión.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <main class="flex min-h-svh items-center justify-center bg-muted/40 p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <div class="mb-2 flex size-10 items-center justify-center rounded-lg bg-primary text-primary-foreground">
          <MessageCircle aria-hidden="true" />
        </div>
        <CardTitle class="text-xl">Iniciar sesión</CardTitle>
        <CardDescription>Ingresa a tu panel de atención por WhatsApp.</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="grid gap-4" novalidate @submit.prevent="onSubmit">
          <div class="grid gap-2">
            <Label for="email">Correo electrónico</Label>
            <Input id="email" v-model="email" type="email" autocomplete="username" required />
          </div>
          <div class="grid gap-2">
            <Label for="password">Contraseña</Label>
            <Input id="password" v-model="password" type="password" autocomplete="current-password" required />
          </div>
          <p v-if="error" role="alert" class="text-sm text-destructive">{{ error }}</p>
          <Button type="submit" :disabled="loading || !email || !password">
            {{ loading ? 'Ingresando…' : 'Ingresar' }}
          </Button>
        </form>
      </CardContent>
    </Card>
  </main>
</template>
