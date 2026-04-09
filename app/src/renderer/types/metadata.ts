export interface InsurerData {
  name: string
  address_line_1: string
  address_line_2: string
  city: string
  province: string
  post_code: string
  country: string
  telephone: string
  email: string
  logo: string
  year_end_month: number | null
  introductory_text: string
  general_provisions_text: string
}

export interface BrokerData {
  id?: number
  name: string
  contact_email: string
  contact_number: string
  fsp_number?: string
  fsp_category?: string
  binder_agreement_ref?: string
  tied_agent_ref?: string
}

export interface CreateBrokerPayload {
  name: string
  contact_email: string
  contact_number: string
  fsp_number?: string
  fsp_category?: string
  binder_agreement_ref?: string
  tied_agent_ref?: string
}

export interface SchemeCategoryData {
  id?: number
  name: string
  description?: string
  active: boolean
}

export interface CreateSchemeCategoryPayload {
  name: string
  description?: string
  active: boolean
}

export interface YearEndMonth {
  name: string
  month_number: number
}

export interface ApiResponse<T> {
  data: T
  status: number
  message?: string
}

export interface ColumnDef {
  headerName: string
  field: string
  valueFormatter?: (params: any) => string
  minWidth: number
  sortable: boolean
  filter: boolean
  resizable: boolean
}
