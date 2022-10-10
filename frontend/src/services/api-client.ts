export class APIClient {
  public apiEndpoint: string
  public clientId: string
  public apiScope: string

  constructor(apiEndpoint: string, clientId: string, apiScope: string) {
    this.apiEndpoint = apiEndpoint
    this.clientId = clientId
    this.apiScope = apiScope
  }
}
