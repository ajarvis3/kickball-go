<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../../utility/fetchSearch";
import SearchItem from "./SearchItem.vue";

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
      <q-btn color="primary" label="Search" @click="doSearch" :loading="loading" />
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
