import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useGroupPricingStore } from '../group_pricing'

describe('useGroupPricingStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initialises with default lookup arrays', () => {
    const store = useGroupPricingStore()
    expect(store.quoteTypes).toEqual(['New Business', 'Renewal'])
    expect(store.obligationTypes).toEqual(['Voluntary', 'Compulsory'])
    expect(store.currencies).toEqual(['USD', 'ZAR'])
    expect(store.riskTypes).toEqual(['All Causes', 'Accidental'])
  })

  it('initialises group_pricing_quote with defaults', () => {
    const store = useGroupPricingStore()
    expect(store.group_pricing_quote.quote_id).toBe(0)
    expect(store.group_pricing_quote.quote_type).toBe('')
    expect(store.group_pricing_quote.scheme_name).toBeNull()
    expect(store.group_pricing_quote.categories).toEqual([])
    expect(store.group_pricing_quote.edit_mode).toBe(false)
  })

  it('updateGroupPricingQuote merges partial updates', () => {
    const store = useGroupPricingStore()
    store.updateGroupPricingQuote({
      scheme_name: 'Test Scheme',
      quote_type: 'Renewal'
    })
    expect(store.group_pricing_quote.scheme_name).toBe('Test Scheme')
    expect(store.group_pricing_quote.quote_type).toBe('Renewal')
    // Other fields remain default
    expect(store.group_pricing_quote.quote_id).toBe(0)
  })

  it('updateGroupPricingQuote does not lose existing fields', () => {
    const store = useGroupPricingStore()
    store.updateGroupPricingQuote({ scheme_name: 'A' })
    store.updateGroupPricingQuote({ quote_type: 'New Business' })
    expect(store.group_pricing_quote.scheme_name).toBe('A')
    expect(store.group_pricing_quote.quote_type).toBe('New Business')
  })

  it('resetGroupPricingQuote restores defaults', () => {
    const store = useGroupPricingStore()
    store.updateGroupPricingQuote({
      scheme_name: 'Modified',
      quote_type: 'Renewal',
      edit_mode: true
    })
    store.resetGroupPricingQuote()
    expect(store.group_pricing_quote.scheme_name).toBeNull()
    expect(store.group_pricing_quote.quote_type).toBe('')
    expect(store.group_pricing_quote.edit_mode).toBe(false)
  })

  it('getInitialQuoteData returns skeleton with undefined values', () => {
    const store = useGroupPricingStore()
    const initial = store.getInitialQuoteData()
    expect(initial.quote_type).toBeNull()
    expect(initial.scheme_name).toBeUndefined()
    expect(initial.industry).toBeUndefined()
  })

  it('scheme_category_template has expected benefit flags', () => {
    const store = useGroupPricingStore()
    const tpl = store.scheme_category_template
    expect(tpl.gla_benefit).toBe(false)
    expect(tpl.ptd_benefit).toBe(false)
    expect(tpl.ci_benefit).toBe(false)
    expect(tpl.phi_benefit).toBe(false)
    expect(tpl.family_funeral_benefit).toBe(false)
    expect(tpl.gla_salary_multiple).toBe(0)
  })

  it('distributionChannels have title/value structure', () => {
    const store = useGroupPricingStore()
    expect(store.distributionChannels).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ title: 'Broker', value: 'broker' })
      ])
    )
    expect(store.distributionChannels).toHaveLength(4)
  })
})
