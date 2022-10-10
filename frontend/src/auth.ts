import { PublicClientApplication, LogLevel, AccountInfo, AuthenticationResult } from '@azure/msal-browser'

export let msalInstance: any

// Config object to be passed to Msal on creation
export function msalInit(clientId: string) {
  // If we have no clientId, we pretend we have an instance and mock it
  if (!clientId) {
    console.log('### Azure AD sign-in: disabled. Will run in demo mode with dummy user')

    // Stub out all the functions we call and return static dummy user where required
    // Use localStorage to simulate MSAL caching and logging out
    msalInstance = {
      clientId: null,

      loginPopup() {
        localStorage.setItem('dummyAccount', JSON.stringify(dummyAccount))
        return new Promise<AuthenticationResult>((resolve) => resolve(dummyAuthRes))
      },
      logout() {
        localStorage.removeItem('dummyAccount')
        window.location.href = '/'
        return new Promise<void>((resolve) => resolve())
      },
      acquireTokenSilent() {
        return new Promise<AuthenticationResult>((resolve) => resolve(dummyAuthRes))
      },
      browserStorage: {
        clear() {
          localStorage.removeItem('dummyAccount')
        },
      },
      getAllAccounts(): AccountInfo[] {
        const dummyAccount = localStorage.getItem('dummyAccount')
        if (dummyAccount) return [JSON.parse(dummyAccount)]
        return []
      },
    }

    return
  }

  // If we have a clientId, then we're using Azure AD for auth
  const config = {
    auth: {
      clientId: clientId,
      redirectUri: window.location.origin,
    },
    cache: {
      cacheLocation: 'localStorage',
    },
    system: {
      loggerOptions: {
        loggerCallback: (level: LogLevel, message: string) => {
          console.log(`### MSAL (${level}): ${message}`)
        },
        logLevel: LogLevel.Verbose,
      },
    },
  }

  // Create the real MSAL application object
  msalInstance = new PublicClientApplication(config)
}

export function isLoggedIn() {
  return msalInstance.getAllAccounts().length > 0
}

const dummyAccount: AccountInfo = {
  username: 'demo@example.net',
  tenantId: '00000000-0000-0000-0000-000000000000',
  nativeAccountId: '00000000-0000-0000-0000-000000000000',
  homeAccountId: '00000000-0000-0000-0000-000000000000',
  localAccountId: '00000000-0000-0000-0000-000000000000',
  name: 'Demo User',
  idToken: '',
  idTokenClaims: {},
  environment: 'dummy.example.net',
}

const dummyAuthRes: AuthenticationResult = {
  authority: 'https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000',
  uniqueId: '00000000-0000-0000-0000-000000000000',
  tenantId: '00000000-0000-0000-0000-000000000000',
  scopes: [],
  account: dummyAccount,
  idToken: '',
  idTokenClaims: {},
  accessToken: '',
  fromCache: false,
  expiresOn: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
  tokenType: 'Bearer',
  correlationId: '00000000-0000-0000-0000-000000000000',
  cloudGraphHostName: '',
  msGraphHost: '',
  state: '',
}
