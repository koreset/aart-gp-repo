import { ref } from 'vue'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import type { InsurerData } from '@/renderer/types/metadata'
import { VALIDATION_MESSAGES } from '@/renderer/constants/metadata'
import GroupPricingService from '@/renderer/api/GroupPricingService'

export function useInsurer() {
  const isLoading = ref(false)
  const isSaving = ref(false)

  const insurerData = ref<InsurerData>({
    name: '',
    address_line_1: '',
    address_line_2: '',
    city: '',
    province: '',
    post_code: '',
    country: '',
    telephone: '',
    email: '',
    logo: '',
    year_end_month: null,
    introductory_text: '',
    general_provisions_text: ''
  })

  const insurerSchema = yup.object({
    name: yup.string().trim().required('Insurer name is required'),
    address_line_1: yup.string().trim().required('Address line 1 is required'),
    city: yup.string().trim().required('City is required'),
    province: yup.string().trim().required('Province is required'),
    post_code: yup.string().trim().required('Postal code is required'),
    telephone: yup
      .string()
      .trim()
      .required(VALIDATION_MESSAGES.INSURER_PHONE_REQUIRED)
      .matches(/^0\d{9}$/, VALIDATION_MESSAGES.INVALID_PHONE),
    email: yup
      .string()
      .trim()
      .required(VALIDATION_MESSAGES.INSURER_EMAIL_REQUIRED)
      .email(VALIDATION_MESSAGES.INVALID_EMAIL)
  })

  const { defineField, handleSubmit, errors, meta } = useForm({
    validationSchema: insurerSchema,
    validateOnMount: false
  })

  const loadInsurer = async () => {
    try {
      isLoading.value = true
      const res = await GroupPricingService.getInsurer()
      if (res.data) {
        insurerData.value = res.data
        return res.data
      }
    } finally {
      isLoading.value = false
    }
  }

  const createInsurer = async (formData: FormData) => {
    try {
      isSaving.value = true
      const validFormData = await insurerSchema.validate(insurerData.value)
      formData.append('insurer', JSON.stringify(validFormData))

      const res = await GroupPricingService.createInsurer(formData)
      return res
    } finally {
      isSaving.value = false
    }
  }

  return {
    insurerData,
    isLoading,
    isSaving,
    errors,
    meta,
    defineField,
    handleSubmit,
    loadInsurer,
    createInsurer
  }
}
