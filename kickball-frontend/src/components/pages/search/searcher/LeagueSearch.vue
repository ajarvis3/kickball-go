<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../../utility/fetchSearch";
import SearchItem from "./SearchItem.vue";

const state = reactive({ leagueId: "", leagueName: "" });
const { leagueId, leagueName } = toRefs(state);

const results = ref<any[]>([]);
const loading = ref(false);

async function doSearch() {
   // if neither provided, list all leagues
   loading.value = true;
   try {
      const params = new URLSearchParams();
      if (leagueId.value) params.append("leagueId", leagueId.value);
      if (leagueName.value) params.append("leagueName", leagueName.value);
      results.value = await fetchSearch("leagues", params);
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
      <SearchItem v-model="leagueName" />
      <q-btn color="primary" label="Search" @click="doSearch" :loading="loading" />
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
