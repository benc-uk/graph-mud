import { createApp } from 'vue'
import * as VueRouter from 'vue-router'
import { msalInit } from '@/services/auth'
import App from './App.vue'
import { routes } from './routes'

import './assets/shared.css'
import { APIClient } from './services/api-client'

export let api: APIClient
export let msalInstance: any
const SCOPES = ['User.Read']

appStartup()

async function appStartup() {
  console.log(`### App running in ${process.env.NODE_ENV} mode`)

  // Take local defaults from .env.development or .env.development.local
  // Or fall back to internal defaults
  let API_ENDPOINT = process.env.VUE_APP_API_ENDPOINT || '/'
  let AUTH_CLIENT_ID = process.env.VUE_APP_AAD_CLIENT_ID || ''

  // Load config at runtime from special `/config` endpoint on frontend-host
  try {
    const configResp = await fetch('.config/API_ENDPOINT,AAD_CLIENT_ID')
    if (configResp.ok) {
      const config = await configResp.json()
      API_ENDPOINT = config.API_ENDPOINT
      AUTH_CLIENT_ID = config.AAD_CLIENT_ID
      console.log('### Config loaded:', config)
    }
  } catch (err) {
    console.warn("### Failed to fetch '/.config' endpoint. Defaults will be used")
  }
  console.log('### API_ENDPOINT:', API_ENDPOINT)
  console.log('### AUTH_CLIENT_ID:', AUTH_CLIENT_ID)

  msalInstance = msalInit(AUTH_CLIENT_ID)
  api = new APIClient(API_ENDPOINT, AUTH_CLIENT_ID, SCOPES)

  const router = VueRouter.createRouter({
    history: VueRouter.createWebHistory(),
    routes,
  })

  const allAccts = await msalInstance.getAllAccounts()
  if (allAccts.length > 0) {
    console.log('### Found existing account:', allAccts[0])
    await msalInstance.setActiveAccount(allAccts[0])
  }

  const app = createApp(App)
  app.use(router)
  app.mount('#app')
}
