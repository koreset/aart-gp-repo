import type { YearEndMonth } from '@/renderer/types/metadata'

export const YEAR_END_MONTHS: readonly YearEndMonth[] = [
  { name: 'January', month_number: 1 },
  { name: 'February', month_number: 2 },
  { name: 'March', month_number: 3 },
  { name: 'April', month_number: 4 },
  { name: 'May', month_number: 5 },
  { name: 'June', month_number: 6 },
  { name: 'July', month_number: 7 },
  { name: 'August', month_number: 8 },
  { name: 'September', month_number: 9 },
  { name: 'October', month_number: 10 },
  { name: 'November', month_number: 11 },
  { name: 'December', month_number: 12 }
] as const

export const VALIDATION_MESSAGES = {
  BROKER_NAME_REQUIRED: 'Broker name is required',
  BROKER_EMAIL_REQUIRED: 'Broker contact email is required',
  BROKER_PHONE_REQUIRED: 'Broker contact number is required',
  INSURER_PHONE_REQUIRED: 'Insurer telephone number is required',
  INSURER_EMAIL_REQUIRED: 'Insurer email is required',
  INVALID_EMAIL: 'Must be a valid email address',
  INVALID_PHONE: 'Must be a valid phone number (format: 0123456789)',
  SERVER_ERROR: 'Server Error: Something went wrong. Please try again later.',
  DUPLICATE_EMAIL: 'This email is already registered',
  BROKER_ADDED: 'Broker added successfully',
  INSURER_SAVED: 'Insurer details saved successfully',
  BENEFITS_SAVED: 'Custom benefit names saved successfully',
  BENEFITS_FAILED: 'Failed to save custom benefit names',
  BROKER_DELETED: 'Broker deletion operation successful',
  BROKER_DELETE_ERROR: 'Error deleting broker. Please try again.'
} as const
