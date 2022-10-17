<template>
  <div class="center">
    <h1 class="mb-3">New Character</h1>

    Name:
    <input type="text" v-model="name" /> <br />
    Description:
    <input type="text" v-model="description" /> <br />
    Class:
    <input type="text" v-model="className" /><br />

    <div class="row">
      <button @click="save" class="golden-btn">SAVE</button>
      <button @click="cancel" class="golden-btn">CANCEL</button>
    </div>
  </div>
</template>

<script lang="ts">
import { api } from '@/main'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'CharacterEdit',

  data: () => ({
    name: '',
    className: '',
    description: '',
  }),

  methods: {
    async save() {
      try {
        await api.createPlayer({
          name: this.name,
          className: this.className,
          description: this.description,
        })
        this.$router.push('/')
      } catch (err) {
        console.error('### Error creating player:', err)
      }
    },

    cancel() {
      this.$router.push({ path: '/' })
    },
  },
})
</script>

<style scoped></style>
