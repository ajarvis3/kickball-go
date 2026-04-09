<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../utility/fetchSearch";

const state = reactive({ gameId: "", leagueId: "" });
const { gameId, leagueId } = toRefs(state);

const results = ref<any[]>([]);
const loading = ref(false);

async function doSearch() {
   if (!gameId.value && !leagueId.value) {
      results.value = [{ error: "gameId or leagueId required" }];
      return;
   }
   loading.value = true;
   try {
      const params = new URLSearchParams();
      if (gameId.value) params.append("gameId", gameId.value);
      if (leagueId.value) params.append("leagueId", leagueId.value);
      results.value = await fetchSearch("games", params);
   } catch (e) {
      results.value = [{ error: String(e) }];
   } finally {
      loading.value = false;
   }
}
</script>

<template>
   <div class="search">
      <SearchItem v-model="gameId" />
      <SearchItem v-model="leagueId" />
      <button @click="doSearch" :disabled="loading">Search</button>
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
