import type {
  IAuthProvider,
  AuthResult,
  LoginCredentials,
  UserProfile
} from '@/renderer/types/auth'

export abstract class BaseAuthProvider {
  protected provider: IAuthProvider

  constructor(provider: IAuthProvider) {
    this.provider = provider
  }

  abstract login(credentials?: LoginCredentials): Promise<AuthResult>
  abstract logout(): Promise<void>

  /**
   * Get provider information
   */
  getProviderInfo(): IAuthProvider {
    return this.provider
  }

  /**
   * Check if provider requires credentials
   */
  requiresCredentials(): boolean {
    return this.provider.type === 'internal' || this.provider.type === 'ldap'
  }

  /**
   * Get provider icon based on type
   */
  getIcon(): string {
    const icons = {
      saml: 'mdi-shield-key',
      oauth2: 'mdi-account-key',
      oidc: 'mdi-openid',
      ldap: 'mdi-server-network',
      'azure-ad': 'mdi-microsoft',
      google: 'mdi-google',
      okta: 'mdi-shield-account',
      internal: 'mdi-account'
    }
    return icons[this.provider.type] || 'mdi-domain'
  }

  /**
   * Handle authentication response
   */
  protected handleAuthResponse(response: Response, data: any): AuthResult {
    if (response.ok && (response.status === 200 || response.status === 201)) {
      return {
        success: true,
        token: data.access_token || data.token,
        user: this.mapUserProfile(data.user || data),
        provider: this.provider.id
      }
    } else if (response.status === 401) {
      return {
        success: false,
        error: 'Invalid credentials',
        provider: this.provider.id
      }
    } else {
      return {
        success: false,
        error: data.message || 'Authentication failed',
        provider: this.provider.id
      }
    }
  }

  /**
   * Map provider-specific user data to standard UserProfile
   */
  protected mapUserProfile(userData: any): UserProfile | undefined {
    if (!userData) return undefined

    return {
      id: userData.id || userData.sub || userData.user_id,
      email: userData.email || userData.preferred_username,
      name:
        userData.name ||
        userData.display_name ||
        `${userData.given_name || ''} ${userData.family_name || ''}`.trim(),
      firstName: userData.given_name || userData.first_name,
      lastName: userData.family_name || userData.last_name,
      roles: userData.roles || userData.groups,
      permissions: userData.permissions,
      avatar: userData.picture || userData.avatar_url,
      provider: this.provider.id
    }
  }

  /**
   * Handle network errors
   */
  protected handleNetworkError(error: any): AuthResult {
    console.error('Network error during authentication:', error)
    return {
      success: false,
      error: 'Network error. Please check your connection and try again.',
      provider: this.provider.id
    }
  }
}
