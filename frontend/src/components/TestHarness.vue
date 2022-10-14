<template>
  <div>
    <h1>Test Harness</h1>
    <button class="golden-btn">TEST</button>

    <div>{{ msg }}</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

import { WebSocketClient } from '@/services/web-sockets'
import { ServerMessage } from '@/services/web-sockets'
import { api } from '@/main'

export default defineComponent({
  name: 'CharacterEdit',

  data() {
    return {
      msg: '',
    }
  },

  mounted() {
    const wsClient = new WebSocketClient(api.apiEndpoint)

    wsClient.addMessageCallback((msg: ServerMessage) => {
      console.log('got a message: ', msg)

      this.msg = msg.text
    })
  },
})
</script>

<style scoped></style>
