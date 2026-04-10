<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import Scoreboard from "../umpire/Scoreboard.vue";
import GameData from "../umpire/GameData.vue";
import AtBatsView from "../umpire/AtBatsView.vue";
import { fetchSearch } from "../../../../utility/fetchSearch";

const route = useRoute();
const gameId = String(route.query.gameId || route.params.id || "");
const leagueId = String(route.query.leagueId || "");

const game = ref<any>(null);
const atbats = ref<any[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

async function fetchGame() {
   const params = new URLSearchParams();
   if (gameId) params.set("gameId", gameId);
   if (leagueId) params.set("leagueId", leagueId);
   return fetchSearch(`/games`, params);
}

async function fetchAtBats() {
   const params = new URLSearchParams({ gameId });
   return fetchSearch(`/atbats`, params);
}

onMounted(async () => {
   loading.value = true;
   error.value = null;
   try {
      const [g, a] = await Promise.all([fetchGame(), fetchAtBats()]);
      game.value = g;
      atbats.value = a;
   } catch (err: any) {
      console.error(err);
      error.value = err?.message || String(err);
   } finally {
      loading.value = false;
   }
});
</script>

<template>
   <div class="page game-dashboard">
      <h1>Game Dashboard</h1>

      <div v-if="loading" class="global-spinner-overlay">
         <q-spinner size="3em" color="primary" />
      </div>
      <div v-else-if="error">Error: {{ error }}</div>

      <section v-else>
         <div class="q-mb-md">
            <router-link
               :to="{
                  path: `/games/${gameId}/umpire`,
                  query: { gameId, leagueId },
               }"
            >
               Open Umpire Dashboard
            </router-link>
         </div>

         <GameData :game="game" />

         <h2>Scoreboard</h2>
         <div class="scoreboard">
            <Scoreboard :game="game" :atbats="atbats" />
         </div>

         <AtBatsView :atbats="atbats" />
      </section>
   </div>
</template>

<!-- removed component styles to rely on Quasar defaults -->
