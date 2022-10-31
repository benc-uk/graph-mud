<template>
  <div id="centerBox">
    <div>
      <span class="mr-3">Client v{{ version }}</span> <span v-if="loggedIn">User: {{ username }}</span>
    </div>
    <div class="logo">🌍 Nano Realms ⚔️</div>

    <button v-if="!loggedIn" @click="login" class="golden-btn">LOGIN</button>

    <router-link v-if="loggedIn && player" to="/play" v-slot="{ href, navigate }">
      <button :href="href" @click="navigate" class="golden-btn mb-3">PLAY</button>
    </router-link>

    <div class="details" v-if="loggedIn && player">
      Signed-in: {{ username }}
      <hr />
      <ul>
        <li>Name: {{ player.name }}</li>
        <li>Class: {{ player.class }}</li>
      </ul>
    </div>

    <router-link v-if="loggedIn && !player" to="/character" v-slot="{ href, navigate }">
      <button :href="href" @click="navigate" class="golden-btn mb-3">NEW CHARACTER</button>
    </router-link>

    <button v-if="loggedIn" @click="logout" class="golden-btn mb-2">LOGOUT</button>

    <button v-if="loggedIn && player" @click="showDeleteDialog = true" class="golden-btn">DELETE CHARACTER</button>

    <button @click="showAboutDialog = true" class="golden-btn mb-2">ABOUT</button>

    <div class="dialog" v-if="showDeleteDialog">
      <h2>Delete Character</h2>
      <p>Are you sure you want to delete your character? This can not be undone!</p>
      <button @click="deleteCharacter" class="golden-btn">DELETE</button>
      <button @click="showDeleteDialog = false" class="golden-btn">CANCEL</button>
    </div>

    <div class="dialog" v-if="showAboutDialog">
      <h2 class="mb-4">Nano Realms</h2>
      <ul>
        <li>Client version: {{ version }}</li>
        <li>Server endpoint: {{ apiEndpoint }}</li>
        <li>Client ID: {{ clientId }}</li>
        <li>Server version: {{ serverInfo.version }}</li>
        <li>Server build: {{ serverInfo.buildInfo }}</li>
        <li>Server hostname: {{ serverInfo.hostname }}</li>
        <li>Server healthy: {{ serverInfo.healthy }}</li>
      </ul>
      <button @click="showAboutDialog = false" class="golden-btn mt-3">OK</button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { api, msalInstance } from '@/main'
import { getUsername, isLoggedIn, globalClientId } from '@/services/auth'
import { PlayerInfo, ServerInfo } from '@/services/api-client'
import * as PackageJSON from '../../package.json'

export default defineComponent({
  name: 'HomePage',

  data: () => ({
    loggedIn: isLoggedIn(),
    username: getUsername(),
    player: null as PlayerInfo | null,
    showDeleteDialog: false,
    showAboutDialog: false,
    version: PackageJSON.version,
    serverInfo: {} as ServerInfo,
    apiEndpoint: api.apiEndpoint,
    clientId: globalClientId || 'Not set, demo mode',
  }),

  async mounted() {
    // check user has a player in the world
    try {
      this.serverInfo = await api.serverStatus()

      if (isLoggedIn()) {
        this.player = await api.getPlayer()
      }
    } catch (error) {
      // That's ok
    }
  },

  methods: {
    async login() {
      try {
        const resp = await msalInstance.loginPopup()
        this.loggedIn = true
        if (!resp) {
          throw new Error('No response from login flow')
        }

        const allAccts = await msalInstance.getAllAccounts()
        if (allAccts.length > 0) {
          console.log('### Found existing account:', allAccts[0])
          await msalInstance.setActiveAccount(allAccts[0])
        }

        this.player = await api.getPlayer()
      } catch (err) {
        if (String(err).includes('Player not found')) {
          return
        }
        console.error('### Error login & getting accounts:', err)
      }
    },

    logout() {
      // remove msal keys from local storage
      for (const key of Object.keys(localStorage)) {
        if (key.startsWith('msal') || key.includes('login') || key.startsWith('fake')) {
          localStorage.removeItem(key)
        }
      }
      this.loggedIn = false
    },

    deleteCharacter() {
      api.deletePlayer()
      this.player = null
      this.showDeleteDialog = false
    },
  },
})
</script>

<style scoped>
button {
  width: 300px;
}

@font-face {
  font-family: gothic;
  src: url(/public/fonts/9577ba5901a597d0bf165f26338d6bd2.woff);
}

.logo {
  font-family: gothic;
  font-size: 80px;
  color: #f0d000;
}

.details {
  background-color: rgba(0, 0, 0, 0.2);
  padding: 1rem;
  margin-bottom: 4rem;
  box-shadow: inset 0 0 10px rgba(0, 0, 0, 0.5);
}

ul {
  text-align: left;
}
</style>
