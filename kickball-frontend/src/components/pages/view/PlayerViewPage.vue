<script setup lang="ts">
import { onMounted, ref } from "vue"
import { useRoute } from "vue-router"

const route = useRoute()
const id = route.params.id as string
const item = ref<any>(null)

onMounted(async () => {
   try {
      const res = await fetch(`/players/${id}`)
      if (!res.ok) throw new Error(res.statusText)
      item.value = await res.json()
   } catch (err) {
      console.error(err)
   }
})
</script>

<template>
   <div class="page player-view-page">
      <h1>Player</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else>Loading...</div>
   </div>
</template>
