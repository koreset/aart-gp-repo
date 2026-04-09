import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useFlashStore } from '../flash'

describe('useFlashStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('initialises with empty hidden state', () => {
    const store = useFlashStore()
    expect(store.message).toBe('')
    expect(store.color).toBe('info')
    expect(store.visible).toBe(false)
  })

  it('setMessage updates message only', () => {
    const store = useFlashStore()
    store.setMessage('hello')
    expect(store.message).toBe('hello')
    expect(store.visible).toBe(false)
  })

  it('show sets message, color, and visible', () => {
    const store = useFlashStore()
    store.show('Success!', 'success')
    expect(store.message).toBe('Success!')
    expect(store.color).toBe('success')
    expect(store.visible).toBe(true)
  })

  it('show defaults color to info', () => {
    const store = useFlashStore()
    store.show('Info message')
    expect(store.color).toBe('info')
    expect(store.visible).toBe(true)
  })

  it('auto-hides after 4 seconds', () => {
    const store = useFlashStore()
    store.show('Temporary')
    expect(store.visible).toBe(true)
    vi.advanceTimersByTime(3999)
    expect(store.visible).toBe(true)
    vi.advanceTimersByTime(1)
    expect(store.visible).toBe(false)
  })

  it('hide sets visible to false', () => {
    const store = useFlashStore()
    store.show('Visible')
    store.hide()
    expect(store.visible).toBe(false)
  })
})
