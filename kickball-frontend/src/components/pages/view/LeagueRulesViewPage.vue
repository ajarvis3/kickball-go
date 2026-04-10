<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchSearch } from "../../../utility/fetchSearch";

const route = useRoute();
const leagueId = String(route.query.leagueId || "");
const version = String(route.query.version || "");
const item = ref<any>(null);

onMounted(async () => {
   try {
      if (!leagueId) {
         console.error("leagueId query parameter is required");
         return;
      }
      const params = new URLSearchParams({ leagueId });
      if (version) params.set("version", version);
      item.value = await fetchSearch(`/leaguerules`, params);
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page leaguerules-view-page">
      <h1>League Rules</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
