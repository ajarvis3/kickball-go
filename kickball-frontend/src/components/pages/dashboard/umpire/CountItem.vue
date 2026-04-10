<script setup lang="ts">
const model = defineModel({ type: Number, default: 0 });

const props = defineProps<{
  label?: string;
  min?: number;
}>();

const dec = () => {
  const current = model.value ?? 0;
  const next = Math.max(props.min ?? 0, current - 1);
  model.value = next;
};

const inc = () => {
  const current = model.value ?? 0;
  model.value = current + 1;
};
</script>

<template>
  <div class="count-item">
    <label>{{ label }}</label>
    <q-btn-group>
      <q-btn flat color="primary" icon="remove" @click="dec" />
      <div class="q-pa-xs">{{ model }}</div>
      <q-btn flat color="primary" icon="add" @click="inc" />
    </q-btn-group>
  </div>
</template>

<style scoped>
.count-item { display: inline-flex; align-items: center; gap: 8px; }
label { margin-right: 8px; }
</style>
