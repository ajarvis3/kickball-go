<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchSearch } from "../../../utility/fetchSearch";

const route = useRoute();
const playerId = String(route.query.playerId || "");
const teamId = String(route.query.teamId || "");
const item = ref<any>(null);

onMounted(async () => {
   try {
      if (!playerId || !teamId) {
         console.error("playerId and teamId query parameters are required");
         return;
      }
      const params = new URLSearchParams({ teamId });
      const players = await fetchSearch(`/players`, params);
      item.value = players.find((p: any) => p.playerId === playerId) || null;
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page player-view-page">
      <h1>Player</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
