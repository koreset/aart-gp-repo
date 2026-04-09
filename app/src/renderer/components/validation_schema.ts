import * as yup from 'yup'

export const groupPricingGeneralInputSchema = yup.object({
  group_pricing_quote: yup.object().shape({
    quote_type: yup.string().required('Quote Type is required')
  })
})
