<script setup lang="ts">

const props = defineProps<{
   atbats?: Array<any>;
}>();
</script>

<template>
   <div>
      <h2>At-Bats</h2>
      <div v-if="!(props.atbats && props.atbats.length)">No at-bats found.</div>
      <q-list v-else bordered separator>
         <q-item v-for="ab in props.atbats" :key="ab.atBatId" clickable>
            <q-item-section>
               <div class="text-subtitle2">
                  {{ ab.playerId }} — {{ ab.result }}
               </div>
               <div class="text-caption">
                  Inning {{ ab.inning }} · RBI {{ ab.rbi }} · Team
                  {{ ab.teamId }}
               </div>
            </q-item-section>
         </q-item>
      </q-list>
   </div>
</template>

<style scoped>
.text-caption {
   color: var(--text);
}
.text-subtitle2 {
   color: var(--text);
}
</style>
