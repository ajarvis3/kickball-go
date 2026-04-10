<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchSearch } from "../../../utility/fetchSearch";

const route = useRoute();
const gameId = String(route.query.gameId || "");
const leagueId = String(route.query.leagueId || "");
const item = ref<any>(null);

onMounted(async () => {
   try {
      if (!gameId && !leagueId) {
         console.error("gameId or leagueId query parameter is required");
         return;
      }
      const params = new URLSearchParams();
      if (gameId) params.set("gameId", gameId);
      if (leagueId) params.set("leagueId", leagueId);
      item.value = await fetchSearch(`/games`, params);
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page game-view-page">
      <h1>Game</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
