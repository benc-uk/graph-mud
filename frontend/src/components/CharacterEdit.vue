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

    <div class="row">
      <button @click="save" class="golden-btn mr-2" v-show="formComplete">CREATE</button>
      <button @click="cancel" class="golden-btn">CANCEL</button>
    </div>
  </div>
</template>

<script lang="ts">
import { api } from '@/main'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'CharacterEdit',

  computed: {
    formComplete() {
      return this.name.trim() != '' && this.description1 != '' && this.description2 != '' && this.className != ''
    },
  },

  data: () => ({
    name: '',
    className: '',
    description1: '',
    description2: '',

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

    descriptionList: [
      'muddy',
      'grim',
      'huaghty',
      'tall',
      'short',
      'fat',
      'skinny',
      'ugly',
      'pretty',
      'handsome',
      'beautiful',
      'dashing',
      'dapper',
      'dainty',
      'dazzling',
      'delightful',
      'divine',
      'dopey',
      'dorky',
      'dramatic',
      'dreamy',
      'drunk',
      'dull',
      'dumb',
      'dusty',
      'dutiful',
      'eager',
      'earnest',
      'easy',
      'elegant',
      'embarrassed',
      'enchanting',
      'encouraging',
      'energetic',
      'enthusiastic',
      'envious',
      'evil',
      'excellent',
      'excited',
      'expensive',
      'exuberant',
      'fabulous',
      'fancy',
      'fantastic',
      'fierce',
      'filthy',
      'fine',
      'foolish',
      'fragile',
      'frail',
      'frantic',
      'friendly',
      'frightened',
      'funny',
      'fuzzy',
      'gentle',
      'giant',
      'giddy',
      'gigantic',
      'glamorous',
      'gleaming',
      'glorious',
      'good',
      'gorgeous',
      'graceful',
      'greasy',
    ],
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
