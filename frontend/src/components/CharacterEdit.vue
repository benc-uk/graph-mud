<template>
  <div class="hello">
    <h1>Edit Your Character</h1>
    <button @click="save" class="golden-btn">SAVE</button>
    <button @click="cancel" class="golden-btn">CANCEL</button>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'CharacterEdit',

  methods: {
    save() {
      this.$router.push({ path: '/' })
    },

    cancel() {
      this.$router.push({ path: '/' })
    },
  },
})
</script>

<style scoped></style>
