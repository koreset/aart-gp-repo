import { defineStore } from 'pinia'

export const useGroupPricingStore = defineStore('groupPricing', {
  state: () => ({
    quoteTypes: ['New Business', 'Renewal'],
    obligationTypes: ['Voluntary', 'Compulsory'],
    currencies: ['USD', 'ZAR'],
    terminalIllnessBenefits: ['Yes', 'No'],
    yesNoItems: ['Yes', 'No'],
    claimsExperiences: ['Yes', 'No', 'Override'],
    riskTypes: ['All Causes', 'Accidental'],
    benefitTypes: ['Standalone', 'Accelerated'],
    disabilityDefinitions: ['Own Occupation', 'Any Occupation'],
    phiEscalationPercentages: ['0%', '5%', '7.5%', '10%'],
    productTypes: ['PHI', 'TTD', 'PHI & TTD'],
    benefitStructures: ['Standalone', 'Accelerated'],
    distributionChannels: [
      { title: 'Broker', value: 'broker' },
      { title: 'Direct', value: 'direct' },
      { title: 'Binder', value: 'binder' },
      { title: 'Tied Agent', value: 'tied_agent' }
    ],
    brokers: [],
    industries: [],
    groupSchemes: [],
    selectedQuote: null,
    memberDataCount: 0,
    claimsExperienceCount: 0,
    scheme_category_template: {
      scheme_category: null,
      ptd_benefit: false,
      gla_benefit: false,
      ci_benefit: false,
      sgla_benefit: false,
      phi_benefit: false,
      ttd_benefit: false,
      family_funeral_benefit: false,
      gla_salary_multiple: 0,
      gla_terminal_illness_benefit: null,
      gla_waiting_period: null,
      gla_educator_benefit: null,
      gla_educator_benefit_type: null,
      gla_benefit_type: null,
      gla_conversion_on_withdrawal: false,
      gla_conversion_on_retirement: false,
      additional_accidental_gla_benefit: false,
      additional_accidental_gla_benefit_type: null,
      tax_saver_benefit: false,
      // Conversion / continuity flags. Default off so pricing is unchanged
      // until the scheme actuary populates the matching general_loadings
      // rate column. The four GLA/PTD/CI base-conversion flags already exist
      // above via gla_conversion_on_withdrawal / gla_conversion_on_retirement
      // / ptd_conversion_on_withdrawal / ci_conversion_on_withdrawal.
      gla_educator_conversion_on_withdrawal: false,
      gla_educator_conversion_on_retirement: false,
      gla_educator_continuity_during_disability: false,
      ptd_educator_conversion_on_withdrawal: false,
      ptd_educator_conversion_on_retirement: false,
      phi_conversion_on_withdrawal: false,
      sgla_conversion_on_withdrawal: false,
      fun_conversion_on_withdrawal: false,
      ttd_conversion_on_withdrawal: false,
      gla_continuity_during_disability: false,
      additional_gla_cover_benefit: false,
      additional_gla_cover_age_band_source: 'standard',
      additional_gla_cover_age_band_type: '',
      additional_gla_cover_custom_age_bands: [] as Array<{
        min_age: number
        max_age: number
      }>,
      additional_gla_cover_band_rates: [] as Array<{
        min_age: number
        max_age: number
        risk_rate_per1000: number
        risk_rate_per1000_male: number
        risk_rate_per1000_female: number
        binder_fee_per1000: number
        binder_fee_per1000_male: number
        binder_fee_per1000_female: number
        outsource_fee_per1000: number
        outsource_fee_per1000_male: number
        outsource_fee_per1000_female: number
        commission_per1000: number
        commission_per1000_male: number
        commission_per1000_female: number
        office_rate_per1000: number
        office_rate_per1000_male: number
        office_rate_per1000_female: number
        male_prop_used: number
        weighted_office_rate_per1000?: number | null
        weighted_office_rate_per1000_male?: number | null
        weighted_office_rate_per1000_female?: number | null
        original_office_rate_per1000?: number | null
        original_office_rate_per1000_male?: number | null
        original_office_rate_per1000_female?: number | null
        smoothed_office_rate_per1000?: number | null
        smoothed_office_rate_per1000_male?: number | null
        smoothed_office_rate_per1000_female?: number | null
        smoothing_factor?: number | null
        smoothing_factor_male?: number | null
        smoothing_factor_female?: number | null
      }>,
      ptd_risk_type: null,
      ptd_benefit_type: null,
      ptd_salary_multiple: 0,
      ptd_deferred_period: null,
      ptd_disability_definition: null,
      ptd_educator_benefit: null,
      ptd_educator_benefit_type: null,
      ptd_conversion_on_withdrawal: false,
      ci_benefit_structure: null,
      ci_benefit_definition: null,
      ci_critical_illness_salary_multiple: 0,
      ci_conversion_on_withdrawal: false,
      sgla_salary_multiple: 0,
      phi_risk_type: null,
      phi_income_replacement_percentage: 0,
      phi_use_tiered_income_replacement_ratio: false,
      phi_tiered_income_replacement_type: 'standard',
      phi_premium_waiver: null,
      phi_medical_aid_premium_waiver: null,
      phi_benefit_escalation: null,
      phi_waiting_period: null,
      phi_deferred_period: null,
      phi_disability_definition: null,
      phi_normal_retirement_age: 0,
      ttd_risk_type: null,
      ttd_income_replacement_percentage: 0,
      ttd_use_tiered_income_replacement_ratio: false,
      ttd_tiered_income_replacement_type: 'standard',
      ttd_premium_waiver_percentage: 0,
      ttd_waiting_period: null,
      ttd_deferred_period: null,
      ttd_disability_definition: null,
      family_funeral_main_member_funeral_sum_assured: 0,
      family_funeral_spouse_funeral_sum_assured: 0,
      family_funeral_children_funeral_sum_assured: 0,
      family_funeral_adult_dependant_sum_assured: 0,
      family_funeral_parent_funeral_sum_assured: 0,
      family_funeral_max_number_children: 0,
      family_funeral_max_number_adult_dependants: 0,
      extended_family_benefit: false,
      extended_family_age_band_source: 'standard',
      extended_family_age_band_type: '',
      extended_family_custom_age_bands: [] as Array<{
        min_age: number
        max_age: number
      }>,
      extended_family_pricing_method: 'rate_per_1000',
      extended_family_sums_assured: [] as Array<{
        min_age: number
        max_age: number
        sum_assured: number
      }>,
      extended_family_band_rates: [] as Array<{
        min_age: number
        max_age: number
        average_rate: number
        sum_assured?: number
        monthly_premium: number
      }>,
      region: ''
    },
    scheme_category: {},
    group_pricing_quote: {
      reviewer: null,
      quote_id: 0,
      creation_date: null,
      quote_type: '',
      scheme_id: 0,
      scheme_name: null,
      scheme_contact: null,
      scheme_email: null,
      quote_broker: null,
      distribution_channel: null,
      obligation_type: '',
      commencement_date: null,
      industry: '',
      categories: [],
      selected_scheme_categories: [],
      scheme_categories: [],
      occupation_class: 0,
      enforce_fcl: false,
      free_cover_limit: 0,
      currency: null,
      exchangeRate: 0,
      normal_retirement_age: 0,
      experience_rating: '',
      use_global_salary_multiple: false,
      basis: null,
      risk_rate_code: null,
      edit_mode: false,
      uploadData: {
        member_data_file: null,
        claims_experience_file: null
      }
    }
  }),
  actions: {
    // setGroupPricing(groupPricing) {
    //   this.groupPricing = groupPricing
    // }
    // reset group_pricing_quote
    updateGroupPricingQuote(payload: any) {
      this.group_pricing_quote = { ...this.group_pricing_quote, ...payload }
    },
    getInitialQuoteData() {
      return {
        quote_type: null,
        scheme_name: undefined,
        scheme_contact: undefined,
        scheme_email: undefined,
        quote_broker: undefined,
        distribution_channel: undefined,
        obligation_type: undefined,
        commencement_date: undefined,
        industry: undefined,
        categories: undefined,
        currency: undefined,
        exchangeRate: undefined,
        experience_rating: undefined,
        free_cover_limit: undefined,
        normal_retirement_age: undefined
      }
    },
    resetGroupPricingQuote() {
      this.group_pricing_quote = {
        reviewer: null,
        quote_id: 0,
        creation_date: null,
        quote_type: '',
        scheme_name: null,
        scheme_contact: null,
        scheme_email: null,
        scheme_id: 0,
        quote_broker: null,
        distribution_channel: null,
        obligation_type: '',
        commencement_date: null,
        industry: '',
        categories: [],
        selected_scheme_categories: [],
        scheme_categories: [],
        occupation_class: 0,
        enforce_fcl: false,
        free_cover_limit: 0,
        currency: null,
        exchangeRate: 0,
        normal_retirement_age: 0,
        experience_rating: '',
        basis: null,
        risk_rate_code: null,
        use_global_salary_multiple: false,
        edit_mode: false,
        uploadData: {
          member_data_file: null,
          claims_experience_file: null
        }
      }
    }
  }
})
