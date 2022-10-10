<template>
  <div class="hello">
    <h1>PLAY GAME</h1>
    <RouterLink class="button" to="/home">EXIT</RouterLink>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { RouterLink } from 'vue-router'

export default defineComponent({
  name: 'PlayGame',
  components: {
    RouterLink,
  },

  methods: {},
})
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
