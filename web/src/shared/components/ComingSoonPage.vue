<script setup lang="ts">
import { Construction } from '@lucide/vue'
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-3 rounded-xl border border-dashed p-12 text-center">
    <Construction class="size-8 text-muted-foreground" aria-hidden="true" />
    <p class="font-medium">{{ $route.meta.title }}</p>
    <p class="text-sm text-muted-foreground">Este módulo todavía está en construcción.</p>
  </div>
</template>
