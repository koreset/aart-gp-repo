import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStatusBarStore } from '../statusBar'

describe('useStatusBarStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initialises with empty items', () => {
    const store = useStatusBarStore()
    expect(store.items).toEqual([])
  })

  it('set replaces items', () => {
    const store = useStatusBarStore()
    const items = [
      { text: 'Connected', icon: 'mdi-wifi', severity: 'info' as const },
      { text: 'Error', severity: 'error' as const }
    ]
    store.set(items)
    expect(store.items).toEqual(items)
  })

  it('clear empties items', () => {
    const store = useStatusBarStore()
    store.set([{ text: 'X' }])
    store.clear()
    expect(store.items).toEqual([])
  })
})
