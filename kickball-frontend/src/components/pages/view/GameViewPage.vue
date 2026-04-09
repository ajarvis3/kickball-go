<script setup lang="ts">
import { useRoute } from "vue-router";
const route = useRoute();
const id = route.params.id as string;
</script>

<template>
   <div class="page game-view-page">
      <h1>Game — {{ id }}</h1>
      <p>Details will be implemented.</p>
   </div>
</template>
