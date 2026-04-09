<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../../../utility/fetchSearch";
import SearchItem from "./SearchItem.vue";

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
      <q-btn color="primary" label="Search" @click="doSearch" :loading="loading" />
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
