import * as msal from '@azure/msal-browser'

// MSAL object used for signing in users with MS identity platform
let msalApp: any

export default {
  //
  // Configure with clientId or empty string/null to set in "demo" mode
  //
  async configure(clientId: string) {
    // Can only call configure once
    if (msalApp) {
      return
    }

    // If no clientId provided then create a mock MSAL UserAgentApplication
    // Allows us to run without Azure AD for demos & local dev
    if (!clientId) {
      console.log('### Azure AD sign-in: disabled. Will run in demo mode with dummy demo@example.net account')

      const dummyUser = {
        accountIdentifier: '00000000-0000-0000-0000-000000000000',
        homeAccountIdentifier: '',
        userName: 'demo@example.net',
        name: 'Demo User',
        idToken: null,
        idTokenClaims: null,
        sid: '',
        environment: '',
      }

      // Stub out all the functions we call and return static dummy user where required
      // Use localStorage to simulate MSAL caching and logging out
      msalApp = {
        clientId: null,

        loginPopup() {
          localStorage.setItem('dummyAccount', JSON.stringify(dummyUser))
          return new Promise<void>((resolve) => resolve())
        },
        logout() {
          localStorage.removeItem('dummyAccount')
          window.location.href = '/'
          return new Promise<void>((resolve) => resolve())
        },
        acquireTokenSilent() {
          return new Promise<void>((resolve) => resolve())
        },
        cacheStorage: {
          clear() {
            localStorage.removeItem('dummyAccount')
          },
        },
        getAccount() {
          const dummyAccount = localStorage.getItem('dummyAccount')
          if (dummyAccount) return JSON.parse(dummyAccount)
          return null
        },
      }
      return
    }

    const config = {
      auth: {
        clientId: clientId,
        redirectUri: window.location.origin,
      },
      cache: {
        cacheLocation: 'localStorage',
      },
      // Only uncomment when you *really* need to debug what is going on in MSAL
      /* system: {
        logger: new msal.Logger(
          (logLevel, msg) => { console.log(msg) },
          {
            level: msal.LogLevel.Verbose
          }
        )
      } */
    }
    console.log('### Azure AD sign-in: enabled\n', config)

    // Create our shared/static MSAL app object
    msalApp = new msal.PublicClientApplication(config)
  },

  //
  // Return the configured client id
  //
  clientId() {
    if (!msalApp) {
      return null
    }

    return msalApp.clientId
  },

  //
  // Login a user with a popup
  //
  async login() {
    if (!msalApp) {
      return
    }

    const LOGIN_SCOPES = ['user.read', 'openid', 'profile', 'email']
    await msalApp.loginPopup({
      scopes: LOGIN_SCOPES,
      prompt: 'select_account',
    })
  },

  //
  // Logout any stored user
  //
  async logout() {
    if (!msalApp) {
      return
    }

    await msalApp.logout()
  },

  //
  // Call to get user, probably cached and stored locally by MSAL
  //
  user() {
    if (!msalApp) {
      return null
    }

    return msalApp.getAccount()
  },

  //
  // Call through to acquireTokenSilent or acquireTokenPopup
  //
  async acquireToken(scopes = ['user.read']) {
    if (!msalApp) {
      return null
    }

    // Set scopes for token request
    const accessTokenRequest = {
      scopes,
    }

    let tokenResp
    try {
      // 1. Try to acquire token silently
      tokenResp = await msalApp.acquireTokenSilent(accessTokenRequest)
      // console.log('### MSAL acquireTokenSilent was successful')
    } catch (err) {
      // 2. Silent process might have failed so try via popup
      tokenResp = await msalApp.acquireTokenPopup(accessTokenRequest)
      // console.log('### MSAL acquireTokenPopup was successful')
    }

    // Just in case check, probably never triggers
    if (!tokenResp.accessToken) {
      throw new Error("### accessToken not found in response, that's bad")
    }

    return tokenResp.accessToken
  },

  //
  // Clear any stored/cached user
  //
  clearLocal() {
    if (msalApp) {
      msalApp.cacheStorage.clear()
    }
  },
}
