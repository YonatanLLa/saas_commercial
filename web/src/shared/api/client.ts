import createClient from 'openapi-fetch'
import type { paths } from './schema'

/**
 * Cliente HTTP tipado a partir de api/openapi.yaml.
 * Regenerar los tipos con `pnpm api:types` cuando cambie el contrato.
 *
 * El access token se guarda solo en memoria; el refresh token viaja en una
 * cookie httpOnly que JavaScript no puede leer (FR-007).
 */
let accessToken: string | null = null

export function setAccessToken(token: string | null) {
  accessToken = token
}

export const api = createClient<paths>({ baseUrl: '/', credentials: 'include' })

api.use({
  onRequest({ request }) {
    if (accessToken) request.headers.set('Authorization', `Bearer ${accessToken}`)
    return request
  },
})

export type { components } from './schema'
