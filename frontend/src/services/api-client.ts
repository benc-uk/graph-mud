import { msalInstance } from '@/main'
import { AuthenticationResult } from '@azure/msal-browser'
import { getUsername } from './auth'

export interface PlayerInfo {
  name: string
  class: string
  description: string
}

export interface LocationInfo {
  name: string
  description: string
  gameEntry: boolean
  exits: string[]
}

export class APIClient {
  public apiEndpoint: string
  public apiScopes: string[]

  constructor(apiEndpoint: string, apiScopes: string[]) {
    this.apiEndpoint = apiEndpoint.replace(/\/$/, '')
    this.apiScopes = apiScopes
  }

  async getPlayer(): Promise<PlayerInfo> {
    return this.baseRequest('player')
  }

  async createPlayer(newPlayer: any) {
    console.log('createPlayer', newPlayer)

    return this.baseRequest('player', 'POST', newPlayer)
  }

  async deletePlayer() {
    return this.baseRequest(`player`, 'DELETE')
  }

  async playerLocation(): Promise<LocationInfo> {
    return this.baseRequest('player/location')
  }

  async cmd(cmdText: string) {
    return this.baseRequest('cmd', 'POST', { text: cmdText })
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
    if (getUsername()) {
      headers.append('X-Username', getUsername())
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

      if (data.title !== undefined) {
        throw new Error(`${data.title}(${data.instance}): ${data.detail}`)
      }
      throw new Error(response.statusText)
    }

    // Return unmarshalled object if response is JSON
    const contentType = response.headers.get('content-type')
    if (contentType && contentType.indexOf('application/json') !== -1) {
      return await response.json()
    }
    return await response.text()
  }
}
