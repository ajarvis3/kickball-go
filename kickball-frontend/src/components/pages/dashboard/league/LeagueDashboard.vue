<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { fetchPost } from "../../../../utility/fetchPost";
import { useRoute } from "vue-router";
import { Notify } from "quasar";
import { fetchSearch } from "../../../../utility/fetchSearch";

const route = useRoute();
const leagueId = String(route.params.id || "");

const teams = ref<any[]>([]);
const games = ref<any[]>([]);
const loading = ref(false);

const newTeam = reactive({ name: "" });
const newGame = reactive({ homeTeamId: "", awayTeamId: "", date: "" });

async function fetchTeams() {
   if (!leagueId) return;
   const params = new URLSearchParams({ leagueId });
   teams.value = await fetchSearch(`/teams`, params);
}

async function fetchGames() {
   if (!leagueId) return;
   const params = new URLSearchParams({ leagueId });
   games.value = await fetchSearch(`/games`, params);
}

async function addTeam() {
   if (!newTeam.name) return;
   try {
      await fetchPost(
         "/teams",
         { "Content-Type": "application/json" },
         JSON.stringify({ leagueId, name: newTeam.name }),
      );
      await fetchTeams();
      newTeam.name = "";
      Notify.create({ type: "positive", message: "Team added" });
   } catch (err: any) {
      Notify.create({ type: "negative", message: String(err) });
   }
}

async function addGame() {
   if (!newGame.homeTeamId || !newGame.awayTeamId) return;
   try {
      // TODO add date/time for games
      const payload = {
         leagueId,
         homeTeamId: newGame.homeTeamId,
         awayTeamId: newGame.awayTeamId,
         name: `Game on ${newGame.date} ${newGame.homeTeamId} vs ${newGame.awayTeamId}`,
      };
      await fetchPost(
         "/games",
         { "Content-Type": "application/json" },
         JSON.stringify(payload),
      );
      await fetchGames();
      newGame.homeTeamId = "";
      newGame.awayTeamId = "";
      newGame.date = "";
      Notify.create({ type: "positive", message: "Game added" });
   } catch (err: any) {
      Notify.create({ type: "negative", message: String(err) });
   }
}

onMounted(async () => {
   loading.value = true;
   try {
      await Promise.all([fetchTeams(), fetchGames()]);
   } catch (e) {
      console.error(e);
   } finally {
      loading.value = false;
   }
});
</script>

<template>
   <div class="page league-dashboard">
      <h1>League Dashboard</h1>
      <div v-if="!leagueId">No league selected (route param :id required)</div>
      <div v-else>
         <section>
            <h2>Teams</h2>
            <div>
               <q-input v-model="newTeam.name" placeholder="New team name" />
               <q-btn color="primary" label="Add Team" @click="addTeam" />
            </div>
            <q-list bordered separator>
               <q-item v-for="t in teams" :key="t.teamId" clickable>
                  <q-item-section>
                     <div class="text-subtitle2">{{ t.name }}</div>
                     <div class="text-caption">ID: {{ t.teamId }}</div>
                  </q-item-section>
                  <q-item-section side>
                     <router-link
                        :to="{
                           path: `/teams/${t.teamId}`,
                           query: { teamId: t.teamId, leagueId },
                        }"
                        >View</router-link
                     >
                  </q-item-section>
               </q-item>
            </q-list>
         </section>

         <section class="q-mt-md">
            <h2>Games</h2>
            <div class="row q-gutter-sm">
               <q-select
                  v-model="newGame.homeTeamId"
                  :options="
                     teams.map((t) => ({ label: t.name, value: t.teamId }))
                  "
                  label="Home team"
                  emit-value
                  map-options
               />
               <q-select
                  v-model="newGame.awayTeamId"
                  :options="
                     teams.map((t) => ({ label: t.name, value: t.teamId }))
                  "
                  label="Away team"
                  emit-value
                  map-options
               />
               <q-input v-model="newGame.date" type="date" label="Date" />
               <q-btn color="primary" label="Add Game" @click="addGame" />
            </div>

            <q-list bordered separator class="q-mt-sm">
               <q-item v-for="g in games" :key="g.gameId" clickable>
                  <q-item-section>
                     <div class="text-subtitle2">Game: {{ g.gameId }}</div>
                     <div class="text-caption">
                        Home: {{ g.homeTeamId }} Away: {{ g.awayTeamId }}
                     </div>
                  </q-item-section>
                  <q-item-section side>
                     <router-link
                        :to="{
                           path: `/games/${g.gameId}/dashboard`,
                           query: { gameId: g.gameId, leagueId },
                        }"
                        >View</router-link
                     >
                  </q-item-section>
               </q-item>
            </q-list>
         </section>
      </div>
   </div>
</template>

<style scoped>
.row {
   display: flex;
   align-items: center;
   gap: 8px;
}
.text-caption {
   color: rgba(0, 0, 0, 0.6);
}
</style>
