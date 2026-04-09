<script setup lang="ts">
const search = defineModel({
   type: String,
   required: true,
});
</script>

<template>
   <q-input
      v-model="search"
      outlined
      dense
   />
</template>
