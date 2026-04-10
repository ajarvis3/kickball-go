<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
   game: any;
   atbats: any[];
}>();

// compute inningsCount based on game state and atbats
const inningsCount = computed(() => {
   let max = 0;
   if (props.game?.state?.inning) max = props.game.state.inning;
   for (const a of props.atbats || []) {
      if (a?.inning > max) max = a.inning;
   }
   return Math.max(1, max);
});

const innings = computed(() =>
   Array.from({ length: inningsCount.value }, (_, i) => i + 1),
);

const computeRunsFromBackend = (isAway: boolean) => {
   const ir = props.game?.state?.inningRuns;
   const count = inningsCount.value;
   if (!Array.isArray(ir) || ir.length === 0) return null;
   const out = Array(count).fill(0);
   for (let i = 0; i < count; i++) {
      const idx = isAway ? i * 2 : i * 2 + 1;
      if (idx < ir.length) out[i] = ir[idx] || 0;
   }
   return out;
};

const computeRunsFromAtBats = (isAway: boolean) => {
   const count = inningsCount.value;
   const arr = Array(count).fill(0);
   for (const a of props.atbats || []) {
      const idx = Math.max(0, Math.min(count - 1, (a.inning || 1) - 1));
      if (!props.game || !a.teamId) continue;
      const matches = isAway
         ? a.teamId === props.game.awayTeamId
         : a.teamId === props.game.homeTeamId;
      if (matches) arr[idx] += a.rbi || 0;
   }
   return arr;
};

const awayRuns = computed(
   () => computeRunsFromBackend(true) ?? computeRunsFromAtBats(true),
);
const homeRuns = computed(
   () => computeRunsFromBackend(false) ?? computeRunsFromAtBats(false),
);

const awayTotal = computed(() =>
   awayRuns.value.reduce((s: number, v: number) => s + v, 0),
);
const homeTotal = computed(() =>
   homeRuns.value.reduce((s: number, v: number) => s + v, 0),
);

const columns = computed(() => {
   const cols: any[] = [{ name: "team", label: "" }];
   for (const i of innings.value)
      cols.push({ name: `i${i}`, label: String(i) });
   cols.push({ name: "R", label: "R" });
   return cols;
});

const rows = computed(() => {
   const away: any = { team: props.game?.awayTeamId || "Away" };
   const home: any = { team: props.game?.homeTeamId || "Home" };
   for (let i = 0; i < innings.value.length; i++) {
      away[`i${i + 1}`] = awayRuns.value[i] ?? 0;
      home[`i${i + 1}`] = homeRuns.value[i] ?? 0;
   }
   away.R = awayTotal.value;
   home.R = homeTotal.value;
   return [away, home];
});
</script>

<template>
   <div class="scoreboard-component">
      <q-table
         :columns="columns"
         :rows="rows"
         flat
         table-layout="auto"
         hide-bottom
      />
   </div>
</template>

<style scoped>
.scoreboard-component {
   width: 100%;
}
</style>
