<template>
  <div id="home">
    <div class="logo">🌐 Graph MUD ⚔️</div>

    <button v-if="!loggedIn" @click="login" class="golden-btn">LOGIN</button>

    <router-link v-if="loggedIn" to="/play" v-slot="{ href, navigate }">
      <button :href="href" @click="navigate" class="golden-btn">ENTER THE DUNGEON</button>
    </router-link>

    <router-link v-if="loggedIn" to="/character" v-slot="{ href, navigate }">
      <button :href="href" @click="navigate" class="golden-btn">EDIT CHARACTER</button>
    </router-link>

    <br />
    <br />
    <button v-if="loggedIn" @click="logout" class="golden-btn">LOGOUT</button>
  </div>
</template>

<script lang="ts">
import { AuthenticationResult } from '@azure/msal-common/dist/response/AuthenticationResult'
import { defineComponent } from 'vue'
import { isLoggedIn, msalInstance } from '../auth'

export default defineComponent({
  name: 'HomePage',

  data: () => ({
    loggedIn: isLoggedIn(),
  }),

  methods: {
    login() {
      msalInstance
        .loginPopup()
        .then((response: AuthenticationResult) => {
          console.log(response)
          this.loggedIn = isLoggedIn()
        })
        .catch((error: any) => {
          console.log(error)
        })
    },

    logout() {
      msalInstance.browserStorage.clear()
      this.loggedIn = false
      //msalInstance.logout()
    },
  },
})
</script>

<style scoped>
button {
  width: 300px;
}

#home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
}

@font-face {
  font-family: gothic;
  src: url(/public/fonts/9577ba5901a597d0bf165f26338d6bd2.woff);
}
.logo {
  font-family: gothic;
  font-size: 90px;
  color: #f0d000;
}
</style>
