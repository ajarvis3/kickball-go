<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Notify } from 'quasar';
import { fetchPost } from "../../../../utility/fetchPost";
import { fetchSearch } from "../../../../utility/fetchSearch";

const leagues = ref<any[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

const form = ref({ name: '', description: '' });

async function loadLeagues() {
  loading.value = true;
  error.value = null;
  try {
    // fetchSearch returns an array for list endpoints
    const res = await fetchSearch('/leagues', new URLSearchParams());
    leagues.value = Array.isArray(res) ? res : (res ? [res] : []);
  } catch (e: any) {
    error.value = e?.message || String(e);
  } finally {
    loading.value = false;
  }
}

async function createLeague() {
  if (!form.value.name) {
    Notify.create({ type: 'negative', message: 'League name is required' });
    return;
  }

  try {
    await fetchPost(
      '/leagues',
      { 'Content-Type': 'application/json' },
      JSON.stringify({ name: form.value.name, description: form.value.description }),
    );
    Notify.create({ type: 'positive', message: 'League created' });
    form.value.name = '';
    form.value.description = '';
    await loadLeagues();
  } catch (e: any) {
    Notify.create({ type: 'negative', message: e?.message || String(e) });
  }
}

onMounted(() => void loadLeagues());
</script>

<template>
  <div class="page admin-dashboard">
    <h1>Admin — Leagues</h1>

    <section>
      <h2>Create League</h2>
      <div class="row q-gutter-sm q-mb-md">
        <q-input v-model="form.name" label="League name" />
        <q-input v-model="form.description" label="Description" />
        <q-btn color="primary" label="Create" @click="createLeague" />
      </div>
    </section>

    <section>
      <h2>Existing Leagues</h2>
      <div v-if="loading">Loading leagues...</div>
      <div v-else-if="error">Error: {{ error }}</div>
      <div v-else>
        <div v-if="!leagues || leagues.length === 0">No leagues found.</div>
        <q-list v-else bordered separator>
          <q-item v-for="l in leagues" :key="l.leagueId" clickable>
            <q-item-section>
              <div class="text-subtitle2">{{ l.name || l.leagueId }}</div>
              <div class="text-caption">ID: {{ l.leagueId }}</div>
            </q-item-section>
            <q-item-section side>
              <router-link :to="{ path: `/leagues/${l.leagueId}`, query: { leagueId: l.leagueId } }">View</router-link>
            </q-item-section>
          </q-item>
        </q-list>
      </div>
    </section>
  </div>
</template>
