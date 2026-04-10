<script setup lang="ts">
import { ref } from "vue";
import type { PropType } from "vue";
import OutcomeComponent from "./OutcomeComponent.vue";

const model = defineModel({
   type: Object as PropType<{
      playerId: string;
      rbi: number;
      selectedOutcome: string;
   }>,
   default: () => ({ playerId: "", rbi: 0, selectedOutcome: "" }),
});

const props = defineProps<{ currentTeamId?: string; posting?: boolean }>();
const emit = defineEmits(["confirm"]);

const showConfirmDialog = ref(false);

function onConfirm() {
   showConfirmDialog.value = false;
   emit("confirm");
}
</script>

<template>
   <section class="atbat-controls">
      <div>
         <label>Player ID:</label>
         <q-input v-model="model.playerId" placeholder="player id" clearable />
      </div>

      <div>
         <label>Team:</label>
         <div class="q-pa-sm">{{ props.currentTeamId || "TBD" }}</div>
      </div>

      <div>
         <label>RBI:</label>
         <q-input type="number" v-model.number="model.rbi" min="0" />
      </div>

      <OutcomeComponent v-model="model.selectedOutcome" />

      <div>
         <q-btn
            color="primary"
            @click="showConfirmDialog = true"
            :loading="props.posting"
            :disable="
               props.posting ||
               !model.selectedOutcome ||
               !props.currentTeamId
            "
            >Post Outcome</q-btn
         >
      </div>

      <q-dialog v-model="showConfirmDialog">
         <q-card style="min-width: 320px">
            <q-card-section>
               <div class="text-h6">Confirm At-Bat</div>
               <div>Player: {{ model.playerId }}</div>
               <div>Team: {{ props.currentTeamId }}</div>
               <div>Outcome: {{ model.selectedOutcome }}</div>
               <div>RBI: {{ model.rbi }}</div>
            </q-card-section>
            <q-card-actions align="right">
               <q-btn flat label="Cancel" color="primary" v-close-popup />
               <q-btn
                  flat
                  label="Confirm"
                  color="primary"
                  @click="onConfirm"
                  :loading="props.posting"
               />
            </q-card-actions>
         </q-card>
      </q-dialog>
   </section>
</template>

<!-- removed component styles to rely on Quasar defaults -->
