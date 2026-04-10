<script setup lang="ts">
import { onMounted, reactive, ref, toRefs } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const gameId = route.params.id as string;

const game = ref<any>(null);
const atbats = ref<any[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);

const countState = reactive({ balls: 0, strikes: 0, fouls: 0 });
const { balls, strikes, fouls } = toRefs(countState);

const playerId = ref("");
const teamId = ref("");
const rbi = ref(0);
const selectedOutcome = ref("");
const posting = ref(false);

// Outcomes mirror the Go game engine's result values
const outcomes = [
   "single",
   "double",
   "triple",
   "homerun",
   "walk",
   "sacrifice",
   "error",
   "fielderschoice",
   "out",
   "strikeout",
   "doubleplay",
   "tripleplay",
];

async function postAtBat() {
   if (
      !selectedOutcome.value ||
      !playerId.value ||
      !teamId.value ||
      !game.value
   )
      return;
   posting.value = true;
   const payload = {
      gameId: gameId,
      leagueId: game.value.leagueId,
      playerId: playerId.value,
      teamId: teamId.value,
      strikes: strikes.value,
      balls: balls.value,
      fouls: fouls.value,
      result: selectedOutcome.value,
      rbi: rbi.value,
   };
   try {
      const res = await fetch(`/atbats`, {
         method: "POST",
         headers: { "Content-Type": "application/json" },
         body: JSON.stringify(payload),
      });
      if (!res.ok) {
         const errText = await res.text();
         throw new Error(errText || res.statusText);
      }
      // success — refresh game and at-bats and reset local counters
      const created = await res.json();
      // refresh resources
      const [g, a] = await Promise.all([fetchGame(), fetchAtBats()]);
      game.value = g;
      atbats.value = a;
      // reset local state
      balls.value = 0;
      strikes.value = 0;
      fouls.value = 0;
      selectedOutcome.value = "";
      rbi.value = 0;
      playerId.value = "";
      teamId.value = "";
   } catch (err: any) {
      console.error(err);
      error.value = err?.message || String(err);
   } finally {
      posting.value = false;
   }
}

async function fetchGame() {
   const res = await fetch(`/games/${gameId}`);
   if (!res.ok) throw new Error(res.statusText);
   return res.json();
}

async function fetchAtBats() {
   const params = new URLSearchParams({ gameId: gameId });
   const res = await fetch(`/atbats?` + params.toString());
   if (!res.ok) throw new Error(res.statusText);
   return res.json();
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

      <div v-if="loading">Loading game and at-bats...</div>
      <div v-else-if="error">Error: {{ error }}</div>

      <section v-else>
         <h2>Game</h2>
         <pre>{{ game }}</pre>

         <h2>Count:</h2>
         <div class="count-row">
            <label>Balls:</label>
            <button @click="balls.value = Math.max(0, balls.value - 1)">
               -
            </button>
            <span>{{ balls }}</span>
            <button @click="balls.value = balls.value + 1">+</button>

            <label>Strikes:</label>
            <button @click="strikes.value = Math.max(0, strikes.value - 1)">
               -
            </button>
            <span>{{ strikes }}</span>
            <button @click="strikes.value = strikes.value + 1">+</button>

            <label>Fouls:</label>
            <button @click="fouls.value = Math.max(0, fouls.value - 1)">
               -
            </button>
            <span>{{ fouls }}</span>
            <button @click="fouls.value = fouls.value + 1">+</button>
         </div>

         <section class="umpire-controls">
            <div>
               <label>Player ID:</label>
               <input v-model="playerId" placeholder="player id" />
            </div>
            <div>
               <label>Team ID:</label>
               <input v-model="teamId" placeholder="team id" />
            </div>
            <div>
               <label>RBI:</label>
               <input type="number" v-model.number="rbi" min="0" />
            </div>

            <div>
               <label>Outcome:</label>
               <select v-model="selectedOutcome">
                  <option value="">-- select outcome --</option>
                  <option v-for="o in outcomes" :key="o" :value="o">
                     {{ o }}
                  </option>
               </select>
            </div>

            <div>
               <button
                  @click="postAtBat"
                  :disabled="
                     posting || !selectedOutcome || !playerId || !teamId
                  "
               >
                  Post Outcome
               </button>
            </div>
         </section>

         <h2>At-Bats</h2>
         <div v-if="atbats.length === 0">No at-bats found.</div>
         <ul>
            <li v-for="ab in atbats" :key="ab.atBatId">{{ ab }}</li>
         </ul>
      </section>
   </div>
</template>
