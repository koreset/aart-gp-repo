// services/validations.ts
import * as yup from 'yup'

export const addressSchema = yup.object().shape({
  street: yup.string().required('Street is required'),
  city: yup.string().required('City is required'),
  postalCode: yup
    .string()
    .required('Postal code is required')
    .matches(/^\d{4,5}$/, 'Invalid postal code format'), // Ensure your regex matches your needs
  country: yup.string().required('Country is required')
})

export const userProfileSchema = yup.object().shape({
  firstName: yup.string().required('First name is required'),
  lastName: yup.string().required('Last name is required'),
  email: yup
    .string()
    .required('Email is required')
    .email('Invalid email address'),
  address: addressSchema.required(), // Mark the nested object itself as required if applicable
  preferences: yup
    .object()
    .shape({
      newsletter: yup.boolean().required(),
      language: yup.string().required('Language selection is required')
    })
    .required()
})

export const myFormValidationSchema = yup.object().shape({
  orderNumber: yup
    .string()
    .required('Order number is required')
    .matches(/^[A-Z0-9-]+$/, 'Invalid order number format'),
  user: userProfileSchema.required() // Mark the nested user object as required
})

// Type for the validated form data (inferred from Yup schema)
export type ValidatedMyFormData = yup.InferType<typeof myFormValidationSchema>
