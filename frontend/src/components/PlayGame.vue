<template>
  <div class="row">
    <div class="colflex main">
      <input class="cmd fullwidth" type="text" spellcheck="false" />
      <textarea v-model="msgLog" class="mywidth" readonly></textarea>
      <br />
      <button class="golden-btn" @click="exit">EXIT</button>
    </div>
    <div class="colflex side">
      <h2>🧭 Exits</h2>
      <textarea v-model="exits" class="mywidth" readonly></textarea>

      <h2 class="topmargin">🗺️ Location</h2>
      <textarea v-model="location" class="mywidth" readonly></textarea>

      <h2 class="topmargin">💼 Inventory</h2>
      <textarea v-model="inventory" class="mywidth" readonly></textarea>

      <h2 class="topmargin">🧑 Player</h2>
      <textarea v-model="player" class="mywidth" readonly></textarea>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { WebSocketClient, ServerMessage } from '@/services/websockets'
import { api } from '@/main'

export default defineComponent({
  name: 'PlayGame',

  data() {
    return {
      msgLog: '',
      exits: 'North\nSouth',
      location: 'A Dark Room',
      player: 'You are a grubby, worrisome mage',
      inventory: 'A Sword\nA Shield',
    }
  },

  mounted() {
    const wsClient = new WebSocketClient(api.apiEndpoint)

    wsClient.addMessageCallback(this.readMessage)
  },

  methods: {
    exit() {
      this.$router.push('/')
    },

    readMessage(msg: ServerMessage) {
      console.log('got a message: ', msg)

      this.msgLog += msg.text + '\n'

      const textarea = document.querySelector('textarea')
      if (textarea) {
        textarea.scrollTop = textarea.scrollHeight
      }
    },
  },
})
</script>

<style scoped>
.cmd {
  background-color: rgb(15, 15, 15);
  color: #ddd;
  font-size: 18px;
  border: 2px solid #109abd;
  border-radius: 0.8rem;
  padding: 0.7rem;
  margin-bottom: 0.2rem;
  caret-color: rgb(8, 172, 95);
  box-shadow: inset 0px 0px 14px 9px #013a0f;
}
.cmd::spelling-error {
  background-color: rgb(179, 197, 17);
  color: #fff;
}
.cmd:focus {
  outline: none;
}
.main {
  flex: 3;
  padding: 1rem;
}
.fullwidth {
  width: 100%;
  -moz-box-sizing: border-box;
  -ms-box-sizing: border-box;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
}
.topmargin {
  margin-top: 0.9rem;
}
.row {
  height: 100vh;
}
.side {
  padding: 1rem;
  background-color: rgb(31, 31, 31);
  box-shadow: inset 0 0 15px rgb(0, 0, 0);
}

textarea {
  background-color: rgb(15, 15, 15);
  color: #fff;
  font-size: 1.2rem;
  border: 2px solid #0b6670;
  border-radius: 0.8rem;
  height: 100%;
  padding: 0.3rem;
  margin: 0.15rem;
  cursor: no-drop;
  box-shadow: inset 0px 0px 14px 4px #013a0f;
}

textarea:focus {
  outline: none;
}
</style>
