<script setup lang="ts">
import { reactive, toRefs, ref } from "vue";
import { fetchSearch } from "../../utility/fetchSearch";

const state = reactive({ gameId: "", playerId: "" });
const { gameId, playerId } = toRefs(state);

const results = ref<any[]>([]);
const loading = ref(false);

async function doSearch() {
   if (!gameId.value && !playerId.value) {
      results.value = [{ error: "gameId or playerId required" }];
      return;
   }
   loading.value = true;
   try {
      const params = new URLSearchParams();
      if (gameId.value) params.append("gameId", gameId.value);
      if (playerId.value) params.append("playerId", playerId.value);
      results.value = await fetchSearch("atbats", params);
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
      <SearchItem v-model="playerId" />
      <button @click="doSearch" :disabled="loading">Search</button>
      <pre>{{ JSON.stringify(results, null, 2) }}</pre>
   </div>
</template>
