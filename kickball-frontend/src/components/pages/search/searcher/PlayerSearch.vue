<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../utility/fetchSearch";

const state = reactive({ teamId: "" });
const { teamId } = toRefs(state);

const results = ref<any[]>([]);
const loading = ref(false);

async function doSearch() {
   if (!teamId.value) {
      results.value = [{ error: "teamId required" }];
      return;
   }
   loading.value = true;
   try {
      const params = new URLSearchParams();
      params.append("teamId", teamId.value);
      results.value = await fetchSearch("players", params);
   } catch (e) {
      results.value = [{ error: String(e) }];
   } finally {
      loading.value = false;
   }
}
</script>

<template>
   <div class="search">
      <SearchItem v-model="teamId" />
      <button @click="doSearch" :disabled="loading">Search</button>
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
