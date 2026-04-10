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
      <q-btn
         color="primary"
         label="Search"
         @click="doSearch"
         :loading="loading"
      />

      <div class="results q-mt-md">
         <q-list v-if="results && results.length">
            <q-item
               v-for="item in results"
               :key="item.leagueId || item.LeagueID || JSON.stringify(item)"
            >
               <q-item-section>
                  <div v-if="item.error" class="text-negative">
                     {{ item.error }}
                  </div>
                  <div v-else>
                     <div class="text-subtitle2">
                        {{ item.name || item.Name }}
                     </div>
                     <div class="text-caption">
                        ID: {{ item.leagueId || item.LeagueID }}
                     </div>
                  </div>
               </q-item-section>
               <q-item-section side v-if="!item.error">
                  <router-link
                     :to="`/leagues/${item.leagueId || item.LeagueID}/dashboard`"
                     >Open Dashboard</router-link
                  >
               </q-item-section>
            </q-item>
         </q-list>
         <div v-else-if="!loading">No results</div>
      </div>
   </div>
</template>
