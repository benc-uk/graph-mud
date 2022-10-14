<template>
  <main>
    <RouterView></RouterView>
  </main>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { RouterView } from 'vue-router'
export default defineComponent({
  components: {
    RouterView,
  },
  name: 'App',
})
</script>

<style>
html {
  background-color: rgb(49, 48, 47);
  color: #ddd;
  font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
  padding: 0;
  margin: 0;
  height: 100%;
}
main,
body {
  padding: 0;
  margin: 0;
  height: 100%;
}
</style>
