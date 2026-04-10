<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchSearch } from "../../../utility/fetchSearch";

const route = useRoute();
const atBatId = String(route.query.atBatId || "");
const gameId = String(route.query.gameId || "");
const playerId = String(route.query.playerId || "");
const item = ref<any>(null);

onMounted(async () => {
   try {
      if (!atBatId) {
         console.error("atBatId query parameter is required");
         return;
      }
      // Handler supports listing by gameId or playerId
      if (!gameId && !playerId) {
         console.error("gameId or playerId query parameter is required to locate at-bat");
         return;
      }
      const params = new URLSearchParams();
      if (gameId) params.set("gameId", gameId);
      if (playerId) params.set("playerId", playerId);
      const atbats = await fetchSearch(`/atbats`, params);
      item.value = atbats.find((a: any) => a.atBatId === atBatId) || null;
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page atbat-view-page">
      <h1>At Bat</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
