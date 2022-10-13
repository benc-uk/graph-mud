<template>
  <div class="hello">
    <h1>Edit Your Character</h1>
    <button @click="test">TEST</button>
    <div>{{ p.p1 }}</div>
    <RouterLink class="button" to="/home">EXIT!</RouterLink>
  </div>
</template>

<script lang="ts">
import { api } from '@/main'
import { defineComponent } from 'vue'
import { RouterLink } from 'vue-router'

export default defineComponent({
  name: 'CharacterEdit',
  components: {
    RouterLink,
  },

  data() {
    return {
      p: '',
    }
  },

  methods: {
    async test() {
      try {
        this.p = await api.getPlayer('')
      } catch (e) {
        console.log('API ERROR', e)
      }
    },
  },
})
</script>

<style scoped></style>
