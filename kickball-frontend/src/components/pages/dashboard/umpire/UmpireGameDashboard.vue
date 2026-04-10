<script setup lang="ts">
import { onMounted, reactive, ref, toRefs, type Ref } from "vue";
import { useRoute } from "vue-router";
import { Notify } from "quasar";
import Scoreboard from "./Scoreboard.vue";
import GameData from "./GameData.vue";
import CountItem from "./CountItem.vue";
import AtBatControls from "./AtBatControls.vue";
import AtBatsView from "./AtBatsView.vue";
import { fetchPost } from "../../../../utility/fetchPost";
import { fetchSearch } from "../../../../utility/fetchSearch";

const route = useRoute();
const gameId = String(route.query.gameId || route.params.id || "");
const leagueId = String(route.query.leagueId || "");

const game = ref<any>(null);
const atbats = ref<any[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

const countState = reactive({ balls: 0, strikes: 0, fouls: 0 });
const { balls, strikes, fouls } = toRefs(countState);

const posting = ref(false);
const atBatForm = reactive({ playerId: "", rbi: 0, selectedOutcome: "" });

import { computed } from "vue";

// currentTeamId is derived from the game's half: top -> away team, bottom -> home team
const currentTeamId = computed(() => {
   const half = String(game.value?.state?.half || "").toLowerCase();
   if (!game.value) return "";
   if (half === "top" || half === "halftop" || half === "top")
      return game.value?.awayTeamId || "";
   if (half === "bottom" || half === "halfbottom" || half === "bottom")
      return game.value?.homeTeamId || "";
   // fallback: if inning half not present, prefer awayTeam
   return game.value?.awayTeamId || "";
});

function resetState(
   balls: Ref<number>,
   strikes: Ref<number>,
   fouls: Ref<number>,
   atBatForm: { selectedOutcome: string; rbi: number; playerId: string },
) {
   balls.value = 0;
   strikes.value = 0;
   fouls.value = 0;
   atBatForm.selectedOutcome = "";
   atBatForm.rbi = 0;
   atBatForm.playerId = "";
}

async function postAtBat() {
   if (!atBatForm.selectedOutcome || !game.value || !currentTeamId.value)
      return;
   posting.value = true;
   const payload = {
      gameId: gameId,
      leagueId: game.value.leagueId,
      playerId: atBatForm.playerId,
      teamId: currentTeamId.value,
      strikes: strikes.value,
      balls: balls.value,
      fouls: fouls.value,
      result: atBatForm.selectedOutcome,
      rbi: atBatForm.rbi,
   };
   try {
      await fetchPost(
         `/atbats`,
         { "Content-Type": "application/json" },
         JSON.stringify(payload),
      );
      // refresh resources
      const [g, a] = await Promise.all([fetchGame(), fetchAtBats()]);
      game.value = g;
      atbats.value = a;
      resetState(balls, strikes, fouls, atBatForm);
      Notify.create({ type: "positive", message: "At-bat posted" });
   } catch (err: any) {
      console.error(err);
      error.value = err?.message || String(err);
      Notify.create({ type: "negative", message: error.value || "" });
   } finally {
      posting.value = false;
   }
}

function postAtBatConfirmed() {
   postAtBat();
}

async function fetchGame() {
   const params = new URLSearchParams();
   if (gameId) params.set("gameId", gameId);
   if (leagueId) params.set("leagueId", leagueId);
   return fetchSearch(`/games`, params);
}

async function fetchAtBats() {
   const params = new URLSearchParams({ gameId: gameId });
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
   <div class="page umpire-game-dashboard">
      <h1>Umpire — Game Dashboard</h1>

      <div v-if="loading" class="global-spinner-overlay">
         <q-spinner size="4em" color="primary" />
      </div>
      <div v-else-if="error">Error: {{ error }}</div>

      <section v-else>
         <h2>Game</h2>
         <GameData :game="game" />

         <h2>Count:</h2>
         <div class="count-row">
            <CountItem v-model="balls" label="Balls" :min="0" />
            <CountItem v-model="strikes" label="Strikes" :min="0" />
            <CountItem v-model="fouls" label="Fouls" :min="0" />
         </div>

         <section class="umpire-controls">
            <AtBatControls
               v-model="atBatForm"
               :currentTeamId="currentTeamId"
               :posting="posting"
               @confirm="postAtBatConfirmed"
            />
         </section>

         <h2>Scoreboard</h2>
         <div class="scoreboard">
            <Scoreboard :game="game" :atbats="atbats" />
         </div>

         <AtBatsView :atbats="atbats" />
      </section>
   </div>
</template>

<!-- removed component styles to rely on Quasar defaults -->
