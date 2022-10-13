import { msalInstance } from '@/main'
import { AuthenticationResult } from '@azure/msal-browser'

export class APIClient {
  public apiEndpoint: string
  public clientId: string
  public apiScopes: string[]

  constructor(apiEndpoint: string, clientId: string, apiScopes: string[]) {
    this.apiEndpoint = apiEndpoint
    this.clientId = clientId
    this.apiScopes = apiScopes
  }

  async getPlayer(playerId: string) {
    console.log('### getPlayer', playerId)

    return this.baseRequest('/players')
  }

  private async baseRequest(path: string, method = 'GET', body?: any): Promise<any> {
    let tokenRes: AuthenticationResult | null = null
    try {
      tokenRes = await msalInstance.acquireTokenSilent({
        scopes: this.apiScopes,
      })
    } catch (e) {
      tokenRes = await msalInstance.acquireTokenPopup({
        scopes: this.apiScopes,
      })
    }
    if (!tokenRes) throw new Error('Failed to get auth token')

    const headers = new Headers({ 'Content-Type': 'application/json' })
    if (tokenRes.accessToken) {
      headers.append('Authorization', `Bearer ${tokenRes.accessToken}`)
    }

    const response = await fetch(`${this.apiEndpoint}/${path}`, {
      method,
      body: body ? JSON.stringify(body) : undefined,
      headers,
    })

    if (!response.ok) {
      // Check if there is a JSON error message
      let data = null
      try {
        data = await response.json()
      } catch (e) {
        // Not JSON, throw a generic error
        throw new Error(response.statusText)
      }

      if (data.error !== undefined) {
        throw new Error(data.error)
      }
      throw new Error(response.statusText)
    }

    return await response.json()
  }
}
