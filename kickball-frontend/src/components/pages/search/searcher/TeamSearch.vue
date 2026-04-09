<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../utility/fetchSearch";

const state = reactive({ leagueId: "" });
const { leagueId } = toRefs(state);

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
      results.value = await fetchSearch("teams", params);
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
      <button @click="doSearch" :disabled="loading">Search</button>
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
