import { BaseAuthProvider } from '../BaseAuthProvider'
import type { AuthResult, LoginCredentials } from '@/renderer/types/auth'

export class InternalAuthProvider extends BaseAuthProvider {
  /**
   * Login with username and password
   */
  async login(credentials: LoginCredentials): Promise<AuthResult> {
    if (!credentials?.username || !credentials?.password) {
      return {
        success: false,
        error: 'Username and password are required',
        provider: this.provider.id
      }
    }

    try {
      const response = await fetch(this.provider.endpoints.login, {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: credentials.username,
          password: credentials.password,
          provider: this.provider.id
        })
      })

      const data = await response.json()
      return this.handleAuthResponse(response, data)
    } catch (error) {
      return this.handleNetworkError(error)
    }
  }

  /**
   * Logout user
   */
  async logout(): Promise<void> {
    try {
      if (this.provider.endpoints.logout) {
        await fetch(this.provider.endpoints.logout, {
          method: 'POST',
          headers: {
            Accept: 'application/json'
          }
        })
      }
    } catch (error) {
      console.error('Logout error:', error)
    }
  }
}
