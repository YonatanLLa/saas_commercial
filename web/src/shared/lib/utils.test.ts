import { describe, expect, it } from 'vitest'
import { cn } from './utils'

describe('cn', () => {
  it('resuelve conflictos de clases de Tailwind', () => {
    expect(cn('px-2 text-sm', 'px-4')).toBe('text-sm px-4')
  })

  it('ignora valores falsos', () => {
    expect(cn('a', false && 'b', undefined, 'c')).toBe('a c')
  })
})
