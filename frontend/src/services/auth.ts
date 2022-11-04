import { PublicClientApplication, LogLevel, AccountInfo, AuthenticationResult } from '@azure/msal-browser'
import { msalInstance } from '@/main'
import { encode } from 'universal-base64url'

const LOG_LEVEL = LogLevel.Warning
export let globalClientId = ''

// Create and store a unique ID for this browser, used for demo accounts
if (!localStorage.getItem('browserId')) {
  localStorage.setItem('browserId', makeId(5))
}

// Config object to be passed to Msal on creation
export function msalInit(clientId: string) {
  globalClientId = clientId
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
        const acct = localStorage.getItem('fakeAccount')
        if (acct) return [JSON.parse(acct)]
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
        logLevel: LOG_LEVEL,
      },
    },
  }

  // Create the real MSAL application object
  return new PublicClientApplication(config)
}

export function isLoggedIn() {
  return msalInstance.getAllAccounts().length > 0
}

export function getUsername() {
  const account = msalInstance.getAllAccounts()[0]
  if (account) return account.username
  return null
}

const fakeAccount: AccountInfo = {
  username: fakeUsername(),
  tenantId: '00000000-0000-0000-0000-000000000000',
  nativeAccountId: '00000000-0000-0000-0000-000000000000',
  homeAccountId: '00000000-0000-0000-0000-000000000000',
  localAccountId: '00000000-0000-0000-0000-000000000000',
  name: `Demo User ${localStorage.getItem('browserId')}`,
  idToken: '',
  idTokenClaims: {},
  environment: fakeUsername(),
}

const fakeAuthRes: AuthenticationResult = {
  authority: 'https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000',
  uniqueId: '00000000-0000-0000-0000-000000000000',
  tenantId: '00000000-0000-0000-0000-000000000000',
  scopes: [],
  account: fakeAccount,
  idToken: '',
  idTokenClaims: {},
  accessToken: 'header.' + encode(`{"preferred_username": "${fakeUsername()}"}`) + '.signature',
  fromCache: false,
  expiresOn: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
  tokenType: 'Bearer',
  correlationId: '00000000-0000-0000-0000-000000000000',
  cloudGraphHostName: '',
  msGraphHost: '',
  state: '',
}

function makeId(tokenLen: number) {
  if (tokenLen == null) {
    tokenLen = 16
  }
  let text = ''
  const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  for (let i = 0; i < tokenLen; ++i) text += possible.charAt(Math.floor(Math.random() * possible.length))

  return text
}

function fakeUsername() {
  return `demo_${localStorage.getItem('browserId')}@example.net`
}
