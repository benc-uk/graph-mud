import { PublicClientApplication, LogLevel, AccountInfo, AuthenticationResult } from '@azure/msal-browser'
import { msalInstance } from '@/main'

// Config object to be passed to Msal on creation
export function msalInit(clientId: string) {
  // If we have no clientId, we pretend we have an instance and mock it
  if (!clientId) {
    console.log('### Azure AD sign-in: disabled. Will run in demo mode with dummy user')

    // Stub out all the functions we call and return static dummy user where required
    // Use localStorage to simulate MSAL caching and logging out
    return {
      clientId: null,

      loginPopup() {
        localStorage.setItem('fakeAccount', JSON.stringify(fakeAccount))
        return new Promise<AuthenticationResult>((resolve) => resolve(fakeAuthRes))
      },
      logout() {
        localStorage.removeItem('fakeAccount')
        window.location.href = '/'
        return new Promise<void>((resolve) => resolve())
      },
      acquireTokenSilent() {
        return new Promise<AuthenticationResult>((resolve) => resolve(fakeAuthRes))
      },
      getCacheManager() {
        return {
          clear() {
            localStorage.removeItem('fakeAccount')
          },
        }
      },
      getAllAccounts(): AccountInfo[] {
        const dummyAccount = localStorage.getItem('fakeAccount')
        if (dummyAccount) return [JSON.parse(dummyAccount)]
        return []
      },
      setActiveAccount() {
        return
      },
    }
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
          console.log(`🔑 (${level}): ${message}`)
        },
        logLevel: LogLevel.Info,
      },
    },
  }

  // Create the real MSAL application object
  return new PublicClientApplication(config)
}

export function isLoggedIn() {
  return msalInstance.getAllAccounts().length > 0
}

const fakeAccount: AccountInfo = {
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

const fakeJWT =
  'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwODAvYmxhaCIsInN1YiI6Im5hbm8tcmVhbG1zIiwibmJmIjoxNjY1Njk0NTQwLCJpYXQiOjE2NjU2OTQ1NDAsImV4cCI6MTY5NzIzMDU0MCwianRpIjoiaTl6b21xNzIiLCJhdWQiOlsibmFuby1yZWFsbXMiXSwiYXV0aF90aW1lIjoxNjk3MjMwNTQwfQ.nsX2Yyi19HAEjL41MmVhwHR6wIOVEUfdXVAXmjJQwxMdX0Puuueao6x8ES9ENLdX-YcbshFTQpxQXP7oG_BYowiuGQ7wZVZ00L3DklbSmo4XYtSDH5G1ndAHu3EtModbhdiuRYHLkF5SAWGBau2OwplLjS0-KslyhB2MpvT7rqNURyBT0hlWy8sxV2GU15fc08MmU8liUbTLp-SWzreJZjcJPv1idK2pqbX_f1U5uSgYCd334esgOdiUVAgOwuDTUD9W6BCX5FDIX6YAi3W0WX33e3_JSqJVUEkKwa0hqqqGUhA5gfxiOrZ6ZTeWh0V9MdXcR34ndIrBm_jkK4tJPA'

const fakeAuthRes: AuthenticationResult = {
  authority: 'https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000',
  uniqueId: '00000000-0000-0000-0000-000000000000',
  tenantId: '00000000-0000-0000-0000-000000000000',
  scopes: [],
  account: fakeAccount,
  idToken: '',
  idTokenClaims: {},
  accessToken: fakeJWT,
  fromCache: false,
  expiresOn: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
  tokenType: 'Bearer',
  correlationId: '00000000-0000-0000-0000-000000000000',
  cloudGraphHostName: '',
  msGraphHost: '',
  state: '',
}
