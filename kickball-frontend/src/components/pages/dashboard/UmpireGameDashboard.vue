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

         <h2>
            Count:
            <span>Balls: {{ balls }}</span>
            <span>Strikes: {{ strikes }}</span>
            <span>Fouls: {{ fouls }}</span>
         </h2>

         <h2>At-Bats</h2>
         <div v-if="atbats.length === 0">No at-bats found.</div>
         <ul>
            <li v-for="ab in atbats" :key="ab.atBatId">{{ ab }}</li>
         </ul>
      </section>
   </div>
</template>
