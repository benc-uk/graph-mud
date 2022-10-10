import { createApp } from 'vue'
import * as VueRouter from 'vue-router'
import { msalInit, msalInstance } from './auth'
import App from './App.vue'
import { routes } from './routes'

import './assets/shared.css'

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

  msalInit(AUTH_CLIENT_ID)

  const router = VueRouter.createRouter({
    history: VueRouter.createWebHistory(),
    routes,
  })

  // api.configure(API_ENDPOINT, AUTH_CLIENT_ID, 'smilr.events')

  const app = createApp(App)
  app.use(router)
  app.mount('#app')
}
