<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../utility/fetchSearch";

// version left as string here; parent can parse to int when making requests
const state = reactive({ leagueId: "", version: "" });
const { leagueId, version } = toRefs(state);

const results = ref<any[]>([]);
const loading = ref(false);

async function doSearch() {
   if (!leagueId.value) {
      results.value = [{ error: "leagueId required" }];
      return;
   }
   loading.value = true;
   try {
      const params = new URLSearchParams();
      params.append("leagueId", leagueId.value);
      if (version.value) params.append("version", version.value);
      results.value = await fetchSearch("leaguerules", params);
   } catch (e) {
      results.value = [{ error: String(e) }];
   } finally {
      loading.value = false;
   }
}
</script>

<template>
   <div class="search">
      <SearchItem v-model="leagueId" />
      <SearchItem v-model="version" />
      <button @click="doSearch" :disabled="loading">Search</button>
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
