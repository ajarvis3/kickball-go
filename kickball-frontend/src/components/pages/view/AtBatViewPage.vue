<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const id = route.params.id as string;
const item = ref<any>(null);

onMounted(async () => {
   try {
      const res = await fetch(`/atbats/${id}`);
      if (!res.ok) throw new Error(res.statusText);
      item.value = await res.json();
   } catch (err) {
      console.error(err);
   }
});
</script>

<template>
   <div class="page atbat-view-page">
      <h1>At Bat</h1>
      <div v-if="item">{{ item }}</div>
      <div v-else class="flex flex-center q-pa-md">
         <q-spinner size="3em" color="primary" />
      </div>
   </div>
</template>
