<script>
  import { onMount, onDestroy } from 'svelte';
  import { sseConnected, activeTab } from './lib/stores.js';
  import { fetchSystemHistory } from './lib/api.js';

  import Navbar from './lib/components/Navbar.svelte';
  import Dashboard from './lib/components/Dashboard.svelte';
  import ContainerList from './lib/components/ContainerList.svelte';
  import LogsExplorer from './lib/components/LogsExplorer.svelte';
  import ImageList from './lib/components/ImageList.svelte';
  import VolumeList from './lib/components/VolumeList.svelte';
  import Toast from './lib/components/Toast.svelte';

  let systemStats = null;
  let history = [];
  let loading = true;
  let eventSource = null;

  onMount(async () => {
    // 1. Fetch historical metrics for charts
    try {
      history = await fetchSystemHistory(40);
    } catch (err) {
      console.error('Failed to fetch history:', err);
    } finally {
      loading = false;
    }

    // 2. Connect to Real-time SSE stats stream
    connectSSE();
  });

  function connectSSE() {
    eventSource = new EventSource('/api/sse/stats');

    eventSource.onopen = () => {
      sseConnected.set(true);
    };

    eventSource.onmessage = (event) => {
      try {
        const stats = JSON.parse(event.data);
        systemStats = stats;

        // Append to history for real-time line charts
        if (stats && stats.timestamp) {
          const snapshot = {
            timestamp: stats.timestamp,
            cpu_usage: stats.host_cpu || 0,
            ram_usage: stats.host_ram?.used || 0,
            ram_total: stats.host_ram?.total || 1,
            disk_usage: stats.host_disk?.used || 0,
            disk_total: stats.host_disk?.total || 1,
            running_containers: stats.running_containers || 0,
            stopped_containers: stats.stopped_containers || 0,
          };

          history = [...history.slice(-59), snapshot];
        }
      } catch (err) {
        console.error('Failed to parse SSE stats:', err);
      }
    };

    eventSource.onerror = (err) => {
      sseConnected.set(false);
      console.warn('SSE stats connection error, reconnecting...', err);
    };
  }

  onDestroy(() => {
    if (eventSource) {
      eventSource.close();
    }
  });
</script>

<div class="min-h-screen flex flex-col font-sans text-slate-100 antialiased selection:bg-sky-500 selection:text-white">
  <Navbar />
  <Toast />

  <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if $activeTab === 'dashboard'}
      <Dashboard
        {systemStats}
        {history}
        {loading}
        onNavigateToContainers={() => activeTab.set('containers')}
      />
    {:else if $activeTab === 'containers'}
      <ContainerList
        containers={systemStats?.containers || []}
        {loading}
      />
    {:else if $activeTab === 'logs'}
      <LogsExplorer />
    {:else if $activeTab === 'images'}
      <ImageList />
    {:else if $activeTab === 'volumes'}
      <VolumeList />
    {/if}
  </main>

  <!-- Footer -->
  <footer class="border-t border-slate-800/60 py-6 text-center text-xs text-slate-500">
    <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2">
      <span class="font-mono">docksight &copy; {new Date().getFullYear()} — High Performance Docker Monitoring</span>
      <span class="text-slate-400 font-mono">Golang + Svelte 5 + Bun + SQLite + SSE</span>
    </div>
  </footer>
</div>
