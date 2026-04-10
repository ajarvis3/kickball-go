<script setup lang="ts">
const props = defineProps<{ game: any }>();
</script>

<template>
   <q-expansion-item
      label="Game Summary"
      icon="sports_baseball"
      expand-separator
   >
      <div class="q-pa-sm">
         <div><strong>Game ID:</strong> {{ props.game?.gameId }}</div>
         <div><strong>League ID:</strong> {{ props.game?.leagueId }}</div>
         <div><strong>Home Team:</strong> {{ props.game?.homeTeamId }}</div>
         <div><strong>Away Team:</strong> {{ props.game?.awayTeamId }}</div>
         <div>
            <strong>Rules Version:</strong> {{ props.game?.rulesVersion }}
         </div>
         <div class="q-mt-sm"><strong>State:</strong></div>
         <div class="q-ml-sm">Inning: {{ props.game?.state?.inning }}</div>
         <div class="q-ml-sm">Half: {{ props.game?.state?.half }}</div>
         <div class="q-ml-sm">Outs: {{ props.game?.state?.outs }}</div>
         <div class="q-ml-sm">
            Home Score: {{ props.game?.state?.homeScore }}
         </div>
         <div class="q-ml-sm">
            Away Score: {{ props.game?.state?.awayScore }}
         </div>
      </div>
   </q-expansion-item>
</template>

<!-- removed component styles to rely on Quasar defaults -->
