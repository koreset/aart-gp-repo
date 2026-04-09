import { BaseAuthProvider } from '../BaseAuthProvider'
import type { AuthResult, LoginCredentials } from '@/renderer/types/auth'

export class SSOAuthProvider extends BaseAuthProvider {
  /**
   * Initiate SSO login (redirect to external provider)
   */
  async login(credentials?: LoginCredentials): Promise<AuthResult> {
    try {
      // For SSO providers, we need to redirect to external authentication
      if (
        this.provider.type === 'saml' ||
        this.provider.type === 'oauth2' ||
        this.provider.type === 'oidc'
      ) {
        // Generate state parameter for security
        const state = this.generateState()
        sessionStorage.setItem('auth_state', state)
        sessionStorage.setItem('auth_provider', this.provider.id)

        // Construct SSO URL
        let ssoUrl = this.provider.endpoints.login

        // Add parameters for OAuth2/OIDC
        if (this.provider.type === 'oauth2' || this.provider.type === 'oidc') {
          const params = new URLSearchParams({
            client_id: this.provider.config.clientId || '',
            response_type: 'code',
            redirect_uri: this.provider.endpoints.callback || '',
            scope:
              this.provider.config.scopes?.join(' ') || 'openid profile email',
            state
          })
          ssoUrl += `?${params.toString()}`
        } else if (this.provider.type === 'saml') {
          // For SAML, add state as RelayState
          const params = new URLSearchParams({
            RelayState: state
          })
          ssoUrl += `?${params.toString()}`
        }

        // Open external authentication in system browser
        await window.mainApi?.send('msgOpenExternalLink', ssoUrl)

        // Return pending result - actual authentication will be handled via callback
        return {
          success: false,
          error: 'SSO_REDIRECT',
          redirectUrl: ssoUrl,
          provider: this.provider.id
        }
      } else {
        return {
          success: false,
          error: 'Unsupported SSO provider type',
          provider: this.provider.id
        }
      }
    } catch (error) {
      return this.handleNetworkError(error)
    }
  }

  /**
   * Handle SSO callback
   */
  async handleCallback(callbackData: {
    code?: string
    state?: string
    error?: string
  }): Promise<AuthResult> {
    try {
      // Verify state parameter
      const storedState = sessionStorage.getItem('auth_state')
      if (!storedState || storedState !== callbackData.state) {
        return {
          success: false,
          error: 'Invalid state parameter',
          provider: this.provider.id
        }
      }

      // Clear stored state
      sessionStorage.removeItem('auth_state')
      sessionStorage.removeItem('auth_provider')

      // Handle error from provider
      if (callbackData.error) {
        return {
          success: false,
          error: callbackData.error,
          provider: this.provider.id
        }
      }

      // Exchange authorization code for token
      if (callbackData.code) {
        const tokenResponse = await this.exchangeCodeForToken(callbackData.code)
        return tokenResponse
      }

      return {
        success: false,
        error: 'No authorization code received',
        provider: this.provider.id
      }
    } catch (error) {
      return this.handleNetworkError(error)
    }
  }

  /**
   * Exchange authorization code for access token
   */
  private async exchangeCodeForToken(code: string): Promise<AuthResult> {
    try {
      const response = await fetch(
        `${this.provider.endpoints.login}/callback`,
        {
          method: 'POST',
          headers: {
            Accept: 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            code,
            provider: this.provider.id,
            client_id: this.provider.config.clientId,
            redirect_uri: this.provider.endpoints.callback
          })
        }
      )

      const data = await response.json()
      return this.handleAuthResponse(response, data)
    } catch (error) {
      return this.handleNetworkError(error)
    }
  }

  /**
   * Logout user (may involve SSO provider logout)
   */
  async logout(): Promise<void> {
    try {
      // Logout from local application
      if (this.provider.endpoints.logout) {
        await fetch(this.provider.endpoints.logout, {
          method: 'POST',
          headers: {
            Accept: 'application/json'
          }
        })
      }

      // For some SSO providers, we might want to logout from the provider as well
      // This would require redirecting to the provider's logout endpoint
    } catch (error) {
      console.error('SSO logout error:', error)
    }
  }

  /**
   * Generate secure state parameter
   */
  private generateState(): string {
    const array = new Uint32Array(4)
    crypto.getRandomValues(array)
    return Array.from(array, (byte) => byte.toString(16).padStart(2, '0')).join(
      ''
    )
  }

  /**
   * SSO providers don't require local credentials
   */
  requiresCredentials(): boolean {
    return false
  }
}
