<template>
  <div class="hello">
    <h1>Edit Your Character</h1>
    <RouterLink class="button" to="/home">EXIT</RouterLink>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { RouterLink } from 'vue-router'

export default defineComponent({
  name: 'CharacterEdit',
  components: {
    RouterLink,
  },

  methods: {},
})
</script>

<style scoped></style>
