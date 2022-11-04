<template>
  <div class="row" style="overflow: hidden">
    <div class="colflex main">
      <input ref="cmdInput" v-model="cmd" class="cmd fullwidth" type="text" spellcheck="false" @keyup.enter="submitCmd" />

      <div class="textBox msgLog">
        <div v-for="m in msgLog" :key="`${m.timestamp}+${m.type}`" class="message" :class="`${m.source}-${m.type}`">{{ m.text }}</div>
      </div>

      <br />
      <button class="golden-btn" @click="exit">EXIT</button>
    </div>
    <div class="colflex side">
      <h2>🗺️ Location</h2>
      <div class="textBox mb-3 mt-1">
        {{ location }}
      </div>

      <h2>🧭 Exits</h2>
      <div class="textBox mb-3 mt-1">
        <div v-for="e in exits" :key="e">{{ e }}</div>
      </div>

      <h2>💼 Inventory</h2>
      <div class="textBox mb-3 mt-1">
        <div v-for="i in inventory" :key="i">{{ i }}</div>
      </div>

      <h2>🧑 Player</h2>
      <textarea v-model="player" class="textBox" readonly></textarea>
    </div>
    <div class="dialog" v-if="serverError">
      <h2>Connection Error 😢</h2>
      <p class="mb-4">{{ serverError }}</p>
      <button @click="reload()" class="golden-btn">RECONNECT</button>
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
      cmd: '',
      msgLog: [] as ServerMessage[],
      exits: Array<string>(),
      location: '',
      player: '',
      inventory: Array<string>(),
      serverError: '',
    }
  },

  async mounted() {
    const wsClient = new WebSocketClient(api.apiEndpoint)

    wsClient.addMessageCallback(this.handleMessage)
    wsClient.addClosedCallback(this.handleClosed)

    const p = await api.getPlayer()
    this.player = `You are ${p.name} a ${p.description} ${p.class}`

    await api.cmd('look')
    await this.update()

    let cmdInput: HTMLInputElement = this.$refs.cmdInput as HTMLInputElement
    cmdInput.focus()
  },

  methods: {
    exit() {
      this.$router.push('/')
    },

    handleMessage(msg: ServerMessage) {
      console.debug(JSON.stringify(msg))

      this.msgLog.push(msg)

      const textarea = document.querySelector('textarea')
      if (textarea) {
        textarea.scrollTop = textarea.scrollHeight
      }

      if (msg.type === 'move') {
        this.update()
      }
    },

    handleClosed() {
      this.msgLog.push({
        type: 'error',
        source: 'server',
        text: 'Connection to server lost',
        timestamp: new Date(),
      })
      this.serverError = 'Connection to server lost'
    },

    async submitCmd() {
      if (this.cmd === 'clear') {
        this.cmd = ''
        this.msgLog = []
        return
      }

      try {
        // Echo locally
        this.msgLog.push({ source: 'local', type: 'command', text: this.cmd, timestamp: new Date() })
        // Actually send the command to the server
        await api.cmd(this.cmd)
      } catch (err) {
        // Push the error to the message log
        this.msgLog.push({ source: 'server', type: 'error', text: `${err}`, timestamp: new Date() })
        console.error(err)
      }

      this.cmd = ''
    },

    async update() {
      const loc = await api.playerLocation()
      this.location = loc.description
      this.exits = loc.exits
      this.inventory = []
      this.inventory.push('Some cheese')
      this.inventory.push('A sword')
    },

    reload() {
      window.location.reload()
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
  border-radius: 0.6rem;
  padding: 0.7rem;
  margin-bottom: 0.2rem;
  caret-color: rgb(8, 172, 95);
  box-shadow: inset 0px 0px 14px 9px #013a0f;
  font-family: 'Overpass Mono', monospace;
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
.row {
  height: 100vh;
}
.side {
  padding: 1rem;
  background-color: rgb(31, 31, 31);
  box-shadow: inset 0 0 15px rgb(0, 0, 0);
  height: 100vh;
}

.msgLog {
  height: 80vh;
  overflow-y: scroll;
}

.textBox {
  background-color: rgb(15, 15, 15);
  color: #fff;
  font-size: 1.2rem;
  border: 2px solid #0b6670;
  border-radius: 0.4rem;
  padding: 0.5rem;
  box-shadow: inset 0px 0px 14px 4px #013a0f;
  font-family: 'Overpass Mono', monospace;
}

.textBox:focus {
  outline: none;
}

.message {
  padding-top: 2px;
  padding-bottom: 2px;
}

.look,
.server-move,
.command-look {
  color: rgb(148, 213, 233);
}
.server-connection {
  color: rgb(101, 209, 101);
}
.server-error {
  color: rgb(235, 25, 154);
  font-weight: 700;
}

.command-invalid {
  color: rgb(219, 195, 60);
}

.command-blocked {
  color: rgb(228, 104, 47);
}

.command-say {
  color: rgb(49, 228, 198);
}
</style>
