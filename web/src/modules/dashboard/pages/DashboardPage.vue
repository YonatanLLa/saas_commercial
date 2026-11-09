<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { api } from '@/shared/api/client'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'

// Primera integración real con el backend: estado de /readyz.
const { data, isPending, isError } = useQuery({
  queryKey: ['system', 'readiness'],
  queryFn: async () => {
    // /readyz responde 503 con el mismo cuerpo cuando alguna dependencia falla.
    const { data, error, response } = await api.GET('/readyz')
    if (data) return data
    if (response.status === 503 && error) return error
    throw new Error(`HTTP ${response.status}`)
  },
  refetchInterval: 15_000,
})
</script>

<template>
  <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
    <Card>
      <CardHeader>
        <CardTitle>Estado del sistema</CardTitle>
        <CardDescription>Conexión con la API, la base de datos y Redis.</CardDescription>
      </CardHeader>
      <CardContent>
        <p v-if="isPending" class="text-sm text-muted-foreground">Consultando…</p>
        <p v-else-if="isError || !data" class="text-sm text-destructive">
          No se pudo contactar a la API. ¿Está corriendo el backend en :8080?
        </p>
        <ul v-else class="space-y-1 text-sm">
          <li class="font-medium">
            General:
            <span :class="data.status === 'ok' ? 'text-green-600' : 'text-destructive'">{{ data.status }}</span>
          </li>
          <li v-for="(value, name) in data.checks" :key="name" class="text-muted-foreground">
            {{ name }}: {{ value }}
          </li>
        </ul>
      </CardContent>
    </Card>
  </div>
</template>
