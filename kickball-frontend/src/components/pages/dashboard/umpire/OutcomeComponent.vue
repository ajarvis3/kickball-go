<script setup lang="ts">
const model = defineModel({ type: String, default: "" });

const defaultOutcomes = [
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

const props = defineProps<{
   options?: string[];
   label?: string;
}>();
</script>

<template>
   <div>
      <label>{{ props.label ?? "Outcome" }}:</label>
      <q-select
         v-model="model"
         :options="props.options ?? defaultOutcomes"
         emit-value
         map-options
      />
   </div>
</template>
