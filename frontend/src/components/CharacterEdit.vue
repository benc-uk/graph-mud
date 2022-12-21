<template>
  <div id="centerBox">
    <h1 class="mb-3">New Character</h1>

    Name
    <input type="text" v-model="name" class="mb-3 fancyInput" /> <br />
    Description 1
    <select class="fancyInput" v-model="description1">
      <option v-for="d in descriptionList" :value="d" :key="d">{{ d }}</option>
    </select>
    Description 2
    <select class="fancyInput" v-model="description2">
      <option v-for="d in descriptionList" :value="d" :key="d">{{ d }}</option>
    </select>
    Class
    <select class="fancyInput" v-model="className">
      <option v-for="c in classList" :value="c.value" :key="c.value">{{ c.name }}</option>
    </select>
    <br />

    <div v-if="error" class="error">
      {{ error }}
    </div>

    <div class="row">
      <button @click="save" class="golden-btn mr-2" v-show="formComplete">CREATE</button>
      <button @click="cancel" class="golden-btn">CANCEL</button>
    </div>
  </div>
</template>

<script lang="ts">
import { api } from '@/main'
import { defineComponent } from 'vue'
import { adjectives, names } from '@/language'

export default defineComponent({
  name: 'CharacterEdit',

  computed: {
    formComplete() {
      return this.name.trim() != '' && this.description1 != '' && this.description2 != '' && this.className != ''
    },
  },

  data: () => ({
    name: names[Math.floor(Math.random() * names.length)] + ' ' + names[Math.floor(Math.random() * adjectives.length)],
    className: '',
    description1: adjectives[Math.floor(Math.random() * adjectives.length)],
    description2: adjectives[Math.floor(Math.random() * adjectives.length)],

    classList: [
      { name: 'Warrior', value: 'warrior' },
      { name: 'Rogue', value: 'rogue' },
      { name: 'Cleric', value: 'cleric' },
      { name: 'Paladin', value: 'paladin' },
      { name: 'Ranger', value: 'ranger' },
      { name: 'Bard', value: 'bard' },
      { name: 'Druid', value: 'druid' },
      { name: 'Monk', value: 'monk' },
      { name: 'Barbarian', value: 'barbarian' },
      { name: 'Sorcerer', value: 'sorcerer' },
      { name: 'Wizard', value: 'wizard' },
    ],

    descriptionList: adjectives,

    error: '',
  }),

  methods: {
    async save() {
      try {
        await api.createPlayer({
          name: this.name.trim(),
          class: this.className,
          description: `${this.description1} ${this.description2}`,
        })
        this.$router.push('/')
      } catch (err) {
        console.error('### Error creating player:', err)
        this.error = err as any
      }
    },

    cancel() {
      this.$router.push({ path: '/' })
    },
  },
})
</script>

<style scoped>
select {
  width: 220px;
}
</style>
