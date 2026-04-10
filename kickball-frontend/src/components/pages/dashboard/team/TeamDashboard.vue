<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { Notify } from "quasar";
import { fetchSearch } from "../../../../utility/fetchSearch";

const route = useRoute();
const teamId = String(route.query.teamId || route.params.id || "");
const leagueId = String(route.query.leagueId || "");

const team = ref<any | null>(null);
const players = ref<any[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

async function fetchTeam() {
   if (!teamId && !leagueId) return null;
   const params = new URLSearchParams();
   if (teamId) params.set("teamId", teamId);
   if (leagueId) params.set("leagueId", leagueId);
   try {
      const t = await fetchSearch("/teams", params);
      // fetchSearch returns list for list endpoints; pick first if array
      if (Array.isArray(t)) return t[0] || null;
      return t;
   } catch (e) {
      throw e;
   }
}

async function fetchPlayers() {
   if (!teamId) return [];
   const params = new URLSearchParams();
   params.set("teamId", teamId);
   return fetchSearch("/players", params);
}

onMounted(async () => {
   loading.value = true;
   error.value = null;
   try {
      team.value = await fetchTeam();
      players.value = await fetchPlayers();
   } catch (err: any) {
      console.error(err);
      error.value = err?.message || String(err);
      Notify.create({ type: "negative", message: error.value });
   } finally {
      loading.value = false;
   }
});
</script>

<template>
   <div class="page team-dashboard">
      <h1>Team Dashboard</h1>

      <div v-if="loading">
         <q-spinner color="primary" />
      </div>

      <div v-else-if="error">Error: {{ error }}</div>

      <section v-else>
         <h2>Team</h2>
         <div v-if="!team">No team found.</div>
         <div v-else>
            <q-card class="q-pa-md q-mb-md">
               <div class="text-h6">{{ team.name || "Unnamed Team" }}</div>
               <div>ID: {{ team.teamId }}</div>
               <div>League: {{ team.leagueId }}</div>
            </q-card>

            <h2>Players</h2>
            <div v-if="!players || players.length === 0">No players found.</div>
            <q-list v-else bordered separator>
               <q-item v-for="p in players" :key="p.playerId" clickable>
                  <q-item-section>
                     <div class="text-subtitle2">
                        {{ p.name || p.playerId }}
                     </div>
                     <div class="text-caption">ID: {{ p.playerId }}</div>
                  </q-item-section>
                  <q-item-section side>
                     <router-link
                        :to="{
                           path: `/players/${p.playerId}`,
                           query: {
                              playerId: p.playerId,
                              teamId: team.teamId,
                              leagueId,
                           },
                        }"
                        >View</router-link
                     >
                  </q-item-section>
               </q-item>
            </q-list>
         </div>
      </section>
   </div>
</template>
