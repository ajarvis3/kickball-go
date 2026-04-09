<script setup lang="ts">
const search = defineModel({
   type: String,
   required: true,
});

const update = (target: EventTarget | null) => {
   if (target === null) return;
   search.value = (target as HTMLTextAreaElement).value;
};
</script>

<template>
   <textarea
      class="input"
      :value="search"
      @input="update($event.target)"
   ></textarea>
</template>
