import { AuthConfigService } from '@/renderer/services/AuthConfigService'
import { InternalAuthProvider } from './providers/InternalAuthProvider'
import { SSOAuthProvider } from './providers/SSOAuthProvider'
import type {
  IAuthProvider,
  AuthResult,
  LoginCredentials
} from '@/renderer/types/auth'
import { BaseAuthProvider } from './BaseAuthProvider'

export class AuthProviderFactory {
  private static providers = new Map<string, BaseAuthProvider>()
  private static configService = AuthConfigService.getInstance()

  /**
   * Get authentication provider by ID
   */
  static getProvider(providerId: string): BaseAuthProvider | null {
    // Check if provider is already cached
    if (this.providers.has(providerId)) {
      return this.providers.get(providerId)!
    }

    // Get provider configuration
    const providerConfig = this.configService.getProvider(providerId)
    if (!providerConfig) {
      console.error(`Provider ${providerId} not found in configuration`)
      return null
    }

    // Create provider instance based on type
    const provider = this.createProvider(providerConfig)
    if (provider) {
      this.providers.set(providerId, provider)
    }

    return provider
  }

  /**
   * Create provider instance based on configuration
   */
  private static createProvider(
    config: IAuthProvider
  ): BaseAuthProvider | null {
    switch (config.type) {
      case 'internal':
        return new InternalAuthProvider(config)

      case 'saml':
      case 'oauth2':
      case 'oidc':
      case 'azure-ad':
      case 'google':
      case 'okta':
        return new SSOAuthProvider(config)

      case 'ldap':
        // LDAP could be implemented later as an extension of InternalAuthProvider
        return new InternalAuthProvider(config)

      default:
        console.error(`Unsupported provider type: ${config.type}`)
        return null
    }
  }

  /**
   * Get default provider
   */
  static getDefaultProvider(): BaseAuthProvider | null {
    const defaultProviderConfig = this.configService.getDefaultProvider()
    if (!defaultProviderConfig) {
      return null
    }
    return this.getProvider(defaultProviderConfig.id)
  }

  /**
   * Get all available providers
   */
  static getAllProviders(): BaseAuthProvider[] {
    const configs = this.configService.getProviders()
    const providers: BaseAuthProvider[] = []

    for (const config of configs) {
      const provider = this.getProvider(config.id)
      if (provider) {
        providers.push(provider)
      }
    }

    return providers
  }

  /**
   * Clear provider cache (useful for testing or configuration changes)
   */
  static clearCache(): void {
    this.providers.clear()
  }

  /**
   * Check if provider requires credentials
   */
  static providerRequiresCredentials(providerId: string): boolean {
    const provider = this.getProvider(providerId)
    return provider?.requiresCredentials() || false
  }

  /**
   * Get provider icon
   */
  static getProviderIcon(providerId: string): string {
    const provider = this.getProvider(providerId)
    return provider?.getIcon() || 'mdi-domain'
  }
}

// Main Authentication Service
export class AuthService {
  // eslint-disable-next-line no-use-before-define
  private static instance: AuthService
  private configService: AuthConfigService
  private currentProvider: BaseAuthProvider | null = null

  private constructor() {
    this.configService = AuthConfigService.getInstance()
  }

  // eslint-disable-next-line no-use-before-define
  public static getInstance(): AuthService {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService()
    }
    return AuthService.instance
  }

  /**
   * Initialize authentication service
   */
  async initialize(clientId?: string): Promise<void> {
    try {
      await this.configService.loadConfiguration(clientId)

      // Set default provider
      const defaultProvider = AuthProviderFactory.getDefaultProvider()
      if (defaultProvider) {
        this.currentProvider = defaultProvider
      }
    } catch (error) {
      console.error('Failed to initialize auth service:', error)
      throw new Error('Authentication service initialization failed')
    }
  }

  /**
   * Login with specified provider
   */
  async login(
    providerId: string,
    credentials?: LoginCredentials
  ): Promise<AuthResult> {
    const provider = AuthProviderFactory.getProvider(providerId)
    if (!provider) {
      return {
        success: false,
        error: `Provider ${providerId} not found`,
        provider: providerId
      }
    }

    this.currentProvider = provider
    return provider.login(credentials)
  }

  /**
   * Logout current user
   */
  async logout(): Promise<void> {
    if (this.currentProvider) {
      await this.currentProvider.logout()
      this.currentProvider = null
    }
  }

  /**
   * Get available providers
   */
  getAvailableProviders(): IAuthProvider[] {
    return this.configService.getProviders()
  }

  /**
   * Get current provider
   */
  getCurrentProvider(): BaseAuthProvider | null {
    return this.currentProvider
  }

  /**
   * Check if provider selection is allowed
   */
  allowsProviderSelection(): boolean {
    return this.configService.allowsProviderSelection()
  }

  /**
   * Get authentication configuration
   */
  getConfiguration() {
    return this.configService.getConfiguration()
  }

  /**
   * Handle SSO callback
   */
  async handleSSOCallback(
    providerId: string,
    callbackData: any
  ): Promise<AuthResult> {
    const provider = AuthProviderFactory.getProvider(providerId)
    if (!provider || !(provider instanceof SSOAuthProvider)) {
      return {
        success: false,
        error: 'Invalid SSO provider',
        provider: providerId
      }
    }

    return provider.handleCallback(callbackData)
  }
}
