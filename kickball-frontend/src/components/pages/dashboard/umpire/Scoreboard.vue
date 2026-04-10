<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
   game: any;
   atbats: any[];
}>();

// compute inningsCount based on game state and atbats
const inningsCount = computed(() => {
   return props.game?.state.inningRuns.length / 2;
});

const innings = computed(() =>
   Array.from({ length: inningsCount.value }, (_, i) => i + 1),
);

const awayRuns = computed(() =>
   props.game?.state?.inningRuns?.filter((_, index: number) => index % 2 === 0),
);
const homeRuns = computed(() =>
   props.game?.state?.inningRuns?.filter((_, index: number) => index % 2 === 1),
);

const awayTotal = computed(() =>
   awayRuns.value.reduce((s: number, v: number) => s + v, 0),
);
const homeTotal = computed(() =>
   homeRuns.value.reduce((s: number, v: number) => s + v, 0),
);

const columns = computed(() => {
   const cols: any[] = [
      {
         name: "team",
         label: "Team",
         field: "team",
         style: "width: 180px; max-width: 180px;",
         headerStyle: "width: 180px;",
      },
   ];
   for (const i of innings.value)
      cols.push({ name: `i${i}`, label: String(i), field: `i${i}` });
   cols.push({ name: "Final", label: "Final", field: "Final" });
   return cols;
});

const rows = computed(() => {
   const away: any = { team: props.game?.awayTeamId || "Away" };
   const home: any = { team: props.game?.homeTeamId || "Home" };
   console.log(
      `away: ${JSON.stringify(awayRuns.value)}, home: ${JSON.stringify(homeRuns.value)}`,
   );
   for (let i = 0; i < innings.value.length; i++) {
      console.log(
         `i ${i + 1}: away ${awayRuns.value[i] ?? 0}, home ${homeRuns.value[i] ?? 0}`,
      );
      away[`i${i + 1}`] = awayRuns.value[i] ?? 0;
      home[`i${i + 1}`] = homeRuns.value[i] ?? 0;
   }
   away.Final = awayTotal.value;
   home.Final = homeTotal.value;
   console.log(away);
   return [away, home];
});
</script>

<template>
   <div class="scoreboard-component">
      <q-table
         :columns="columns"
         :rows="rows"
         row-key="team"
         flat
         table-layout="fixed"
         hide-bottom
      />
   </div>
</template>

<!-- removed component styles to rely on Quasar defaults -->
