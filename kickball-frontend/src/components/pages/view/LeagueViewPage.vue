<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchSearch } from "../../../utility/fetchSearch";

const route = useRoute();
const leagueId = String(route.query.leagueId || "");
const leagueName = String(route.query.leagueName || "");
const item = ref<any>(null);

onMounted(async () => {
   try {
      const params = new URLSearchParams();
      if (leagueId) params.set("leagueId", leagueId);
      if (leagueName) params.set("leagueName", leagueName);
      if (!leagueId && !leagueName) {
         console.error("leagueId or leagueName query parameter is required");
         return;
      }
      item.value = await fetchSearch(`/leagues`, params);
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page league-view-page">
      <h1>League</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
