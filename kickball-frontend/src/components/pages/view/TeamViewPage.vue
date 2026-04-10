<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchSearch } from "../../../utility/fetchSearch";

const route = useRoute();
const teamId = String(route.query.teamId || "");
const leagueId = String(route.query.leagueId || "");
const item = ref<any>(null);

onMounted(async () => {
   try {
      if (!teamId || !leagueId) {
         console.error("teamId and leagueId query parameters are required");
         return;
      }
      const params = new URLSearchParams({ leagueId });
      const teams = await fetchSearch(`/teams`, params);
      item.value = teams.find((t: any) => t.teamId === teamId) || null;
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page team-view-page">
      <h1>Team</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
