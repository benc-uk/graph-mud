<template>
  <div>
    <h1>Test Harness</h1>
    <button class="golden-btn" @click="test">TEST</button>

    <div>{{ player }}</div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { api } from '@/main'

export default defineComponent({
  name: 'CharacterEdit',

  data() {
    return {
      player: null,
    }
  },

  methods: {
    async test() {
      try {
        this.player = await api.getPlayer()
      } catch (err) {
        console.log(err)
      }
    },
  },
})
</script>

<style scoped></style>
