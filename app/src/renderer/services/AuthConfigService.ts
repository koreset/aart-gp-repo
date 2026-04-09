import type { AuthConfiguration, IAuthProvider } from '@/renderer/types/auth'

export class AuthConfigService {
  // eslint-disable-next-line no-use-before-define
  private static instance: AuthConfigService
  private config: AuthConfiguration | null = null
  private authUrl: string

  private constructor() {
    this.authUrl = import.meta.env.VITE_APP_AUTH_URL
  }

  // eslint-disable-next-line no-use-before-define
  public static getInstance(): AuthConfigService {
    if (!AuthConfigService.instance) {
      AuthConfigService.instance = new AuthConfigService()
    }
    return AuthConfigService.instance
  }

  /**
   * Load authentication configuration from server or fallback to default
   */
  async loadConfiguration(clientId?: string): Promise<AuthConfiguration> {
    try {
      // Try to detect client ID from various sources
      const detectedClientId = clientId || (await this.detectClient())

      // Construct config URL
      const configUrl = `${this.authUrl}/auth/config${detectedClientId ? `?client=${detectedClientId}` : ''}`

      const response = await fetch(configUrl, {
        method: 'GET',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json'
        }
      })

      if (response.ok) {
        this.config = await response.json()
        return this.config!
      } else {
        console.warn('Failed to load auth config from server, using fallback')
        this.config = this.getFallbackConfiguration()
        return this.config
      }
    } catch (error) {
      console.error('Error loading auth configuration:', error)
      this.config = this.getFallbackConfiguration()
      return this.config
    }
  }

  /**
   * Detect client ID from various sources
   */
  private async detectClient(): Promise<string | null> {
    try {
      // Method 1: Check for client ID in environment
      const envClientId = import.meta.env.VITE_APP_CLIENT_ID
      if (envClientId && envClientId !== 'default') {
        return envClientId
      }

      // Method 2: Check subdomain
      const hostname = window.location.hostname
      const parts = hostname.split('.')
      if (parts.length > 2 && parts[0] !== 'www') {
        return parts[0] // Use subdomain as client ID
      }

      // Method 3: Check for stored client preference
      const storedClient = localStorage.getItem('aart_client_id')
      if (storedClient) {
        return storedClient
      }

      // Method 4: Check URL parameters
      const urlParams = new URLSearchParams(window.location.search)
      const urlClient = urlParams.get('client')
      if (urlClient) {
        return urlClient
      }

      return null
    } catch (error) {
      console.error('Error detecting client:', error)
      return null
    }
  }

  /**
   * Get fallback configuration for internal authentication
   */
  private getFallbackConfiguration(): AuthConfiguration {
    return {
      defaultProvider: 'internal',
      allowProviderSelection: false,
      providers: [
        {
          id: 'internal',
          name: 'internal',
          type: 'internal',
          displayName: 'AART Authentication',
          description: 'Internal authentication system',
          config: {},
          endpoints: {
            login: `${this.authUrl}/login`,
            logout: `${this.authUrl}/logout`
          }
        }
      ],
      branding: {
        companyName: 'AART Enterprise',
        loginTitle: 'Welcome Back',
        loginSubtitle: 'Sign in to your account'
      }
    }
  }

  /**
   * Get current configuration
   */
  getConfiguration(): AuthConfiguration | null {
    return this.config
  }

  /**
   * Get available providers
   */
  getProviders(): IAuthProvider[] {
    return this.config?.providers || []
  }

  /**
   * Get provider by ID
   */
  getProvider(providerId: string): IAuthProvider | null {
    return this.config?.providers.find((p) => p.id === providerId) || null
  }

  /**
   * Get default provider
   */
  getDefaultProvider(): IAuthProvider | null {
    if (!this.config) return null
    return this.getProvider(this.config.defaultProvider)
  }

  /**
   * Check if provider selection is allowed
   */
  allowsProviderSelection(): boolean {
    return this.config?.allowProviderSelection || false
  }

  /**
   * Set client preference
   */
  setClientPreference(clientId: string): void {
    localStorage.setItem('aart_client_id', clientId)
  }

  /**
   * Clear client preference
   */
  clearClientPreference(): void {
    localStorage.removeItem('aart_client_id')
  }
}
