// Authentication Provider Types
export interface IAuthProvider {
  id: string
  name: string
  type:
    | 'saml'
    | 'oauth2'
    | 'oidc'
    | 'ldap'
    | 'azure-ad'
    | 'google'
    | 'okta'
    | 'internal'
  displayName: string
  description?: string
  config: AuthProviderConfig
  endpoints: AuthEndpoints
}

export interface AuthProviderConfig {
  clientId?: string
  clientSecret?: string
  domain?: string
  tenantId?: string
  scopes?: string[]
  metadata?: any
}

export interface AuthEndpoints {
  login: string
  logout: string
  callback?: string
  metadata?: string
  userInfo?: string
}

export interface AuthConfiguration {
  defaultProvider: string
  allowProviderSelection: boolean
  providers: IAuthProvider[]
  branding?: AuthBranding
  clientId?: string
}

export interface AuthBranding {
  logo?: string
  primaryColor?: string
  companyName?: string
  loginTitle?: string
  loginSubtitle?: string
}

export interface AuthResult {
  success: boolean
  token?: string
  user?: UserProfile
  redirectUrl?: string
  error?: string
  provider?: string
}

export interface UserProfile {
  id: string
  email: string
  name: string
  firstName?: string
  lastName?: string
  roles?: string[]
  permissions?: string[]
  avatar?: string
  provider?: string
}

export interface LoginCredentials {
  username: string
  password: string
  provider?: string
  rememberMe?: boolean
}

export interface SSOCallbackData {
  code?: string
  state?: string
  token?: string
  error?: string
  provider: string
}
