<script>
  import { onMount, onDestroy } from 'svelte';
  import { fetchContainers } from '../api.js';
  import { selectedLogContainerId } from '../stores.js';
  import LogViewer from './LogViewer.svelte';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import { Terminal, Box, Search, Layers, ChevronDown, ChevronRight, Maximize2, RefreshCw, CheckCircle2, XCircle } from 'lucide-svelte';

  let containers = [];
  let loading = true;
  let searchQuery = '';
  let activeContainerId = null; // null = 'all' grouped mode
  let expandedAccordions = {}; // { [containerId]: boolean }

  onMount(async () => {
    await loadContainers();

    // Check if a container was pre-selected from ContainerList tab
    if ($selectedLogContainerId) {
      activeContainerId = $selectedLogContainerId;
      expandedAccordions[activeContainerId] = true;
    }
  });

  async function loadContainers() {
    loading = true;
    try {
      containers = await fetchContainers();

      // Default expand ONLY the first 2 running containers in accordion view
      // to avoid exceeding browser HTTP 1.1 concurrent SSE connection limits
      let expandedCount = 0;
      containers.forEach(c => {
        if (c.state === 'running' && !activeContainerId && expandedCount < 2) {
          expandedAccordions[c.id] = true;
          expandedCount++;
        }
      });
    } catch (err) {
      console.error('Failed to fetch containers for logs:', err);
    } finally {
      loading = false;
    }
  }

  function selectContainer(id) {
    activeContainerId = id;
    selectedLogContainerId.set(id);
    if (id) {
      expandedAccordions[id] = true;
    }
  }

  function toggleAccordion(id) {
    expandedAccordions[id] = !expandedAccordions[id];
    expandedAccordions = { ...expandedAccordions };
  }

  $: filteredContainers = containers.filter(c => {
    if (!searchQuery) return true;
    const q = searchQuery.toLowerCase();
    const name = (c.names?.[0] || '').toLowerCase();
    const id = (c.id || '').toLowerCase();
    const img = (c.image || '').toLowerCase();
    return name.includes(q) || id.includes(q) || img.includes(q);
  });

  $: currentContainer = containers.find(c => c.id === activeContainerId);
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-white flex items-center gap-3">
        <Terminal class="w-7 h-7 text-sky-400" /> Container Logs Explorer
      </h2>
      <p class="text-sm text-slate-400 mt-1">
        Grafana-style drill-down log streams grouped by container
      </p>
    </div>

    <button
      onclick={loadContainers}
      class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-slate-900/80 border border-slate-800 text-xs font-semibold text-slate-300 hover:text-white hover:border-slate-700 transition-all"
    >
      <RefreshCw class="w-4 h-4 text-sky-400" />
      <span>Refresh Containers</span>
    </button>
  </div>

  <!-- Grafana Drill-Down Layout Grid -->
  <div class="grid grid-cols-1 lg:grid-cols-4 gap-6 items-start relative">
    <!-- Left Sidebar: Sticky Container Group Selector List -->
    <div class="lg:sticky lg:top-20 glass-card rounded-2xl p-4 border border-slate-800/80 space-y-4 max-h-[calc(100vh-6rem)] flex flex-col shadow-xl">
      <div class="flex items-center justify-between pb-2 border-b border-slate-800/80 shrink-0">
        <span class="text-xs font-bold text-slate-400 uppercase tracking-wider flex items-center gap-2">
          <Layers class="w-4 h-4 text-sky-400" /> Container Groups
        </span>
        <span class="text-xs font-mono px-2 py-0.5 rounded-full bg-slate-800 text-slate-400">
          {containers.length}
        </span>
      </div>

      <!-- Quick Search Filter -->
      <div class="relative shrink-0">
        <Search class="w-3.5 h-3.5 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Filter containers..."
          class="w-full bg-slate-950/80 border border-slate-800 rounded-xl pl-8 pr-3 py-1.5 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:border-sky-500 transition-all"
        />
      </div>

      <!-- Container Selector List -->
      <div class="space-y-1.5 flex-1 overflow-y-auto pr-1">
        <!-- Option: All Containers Grouped Mode -->
        <button
          onclick={() => selectContainer(null)}
          class="w-full text-left p-3 rounded-xl border transition-all flex items-center justify-between group
            {activeContainerId === null ? 'sidebar-btn-active' : 'sidebar-btn-inactive'}"
        >
          <div class="flex items-center gap-2.5">
            <div class="p-1.5 rounded-lg bg-sky-500/20 text-sky-400">
              <Layers class="w-4 h-4" />
            </div>
            <div>
              <div class="font-bold text-xs font-mono sidebar-btn-title">All Containers</div>
              <div class="text-[10px] sidebar-btn-subtitle">Grouped Drill-Down</div>
            </div>
          </div>
          <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-sky-950 text-sky-400 border border-sky-800/50">
            Grouped
          </span>
        </button>

        <div class="pt-2 text-[11px] font-semibold text-slate-500 uppercase tracking-wider px-1">
          Individual Containers
        </div>

        {#if loading}
          <div class="p-4 text-center text-xs text-slate-500">Loading containers...</div>
        {:else if filteredContainers.length === 0}
          <div class="p-4 text-center text-xs text-slate-500 italic">No containers match search</div>
        {:else}
          {#each filteredContainers as c (c.id)}
            <button
              onclick={() => selectContainer(c.id)}
              class="w-full text-left p-2.5 rounded-xl border transition-all flex items-center justify-between group
                {activeContainerId === c.id ? 'sidebar-btn-active' : 'sidebar-btn-inactive'}"
            >
              <div class="flex items-center gap-2.5 min-w-0">
                <span class="w-2 h-2 rounded-full shrink-0
                  {c.state === 'running' ? 'bg-emerald-400 shadow-sm shadow-emerald-400' : 'bg-rose-500'}"
                ></span>
                <div class="truncate">
                  <div class="font-bold text-xs font-mono truncate sidebar-btn-title">
                    {c.names?.[0]?.replace(/^\//, '') || c.id.substring(0, 12)}
                  </div>
                  <div class="text-[10px] font-mono truncate sidebar-btn-subtitle">{c.image}</div>
                </div>
              </div>

              <span class="text-[10px] font-mono uppercase px-2 py-0.5 rounded-full shrink-0 font-semibold
                {c.state === 'running'
                  ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
                  : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}"
              >
                {c.state}
              </span>
            </button>
          {/each}
        {/if}
      </div>
    </div>

    <!-- Main Drill-Down Log View Panel -->
    <div class="lg:col-span-3 space-y-6">
      {#if activeContainerId !== null && currentContainer}
        <!-- Focused Single Container Drill-Down Header -->
        <div class="glass-card rounded-2xl p-4 border border-slate-800/80 flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <span class="w-3 h-3 rounded-full shrink-0
              {currentContainer.state === 'running' ? 'bg-emerald-400 shadow-md shadow-emerald-400' : 'bg-rose-500'}"
            ></span>
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-base text-white font-mono">
                  {currentContainer.names?.[0]?.replace(/^\//, '') || currentContainer.id.substring(0, 12)}
                </h3>
                <span class="text-xs font-mono px-2 py-0.5 rounded-full bg-slate-800 text-slate-400">
                  {currentContainer.id.substring(0, 12)}
                </span>
              </div>
              <p class="text-xs text-slate-400 font-mono mt-0.5">{currentContainer.image}</p>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <button
              onclick={() => selectContainer(null)}
              class="px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold border border-slate-700/60 transition-all"
            >
              Back to All Containers
            </button>
          </div>
        </div>

        <!-- Embedded Log Viewer -->
        <LogViewer containerId={activeContainerId} />
      {:else}
        <!-- Grouped Accordion Drill-Down View (Grafana Style) -->
        <div class="space-y-4">
          {#if loading}
            <SkeletonLoader type="card" rows={3} />
          {:else if containers.length === 0}
            <div class="glass-card rounded-2xl p-12 text-center space-y-3">
              <Box class="w-12 h-12 text-slate-600 mx-auto" />
              <h3 class="text-base font-semibold text-slate-300">No Containers Available</h3>
              <p class="text-xs text-slate-500">There are no Docker/Podman containers available to log.</p>
            </div>
          {:else}
            {#each filteredContainers as c (c.id)}
              <div class="glass-card rounded-2xl border border-slate-800/80 overflow-hidden shadow-xl transition-all">
                <!-- Accordion Container Header Bar -->
                <div class="p-4 bg-slate-900/60 border-b border-slate-800/80 flex items-center justify-between gap-4">
                  <button
                    onclick={() => toggleAccordion(c.id)}
                    class="flex items-center gap-3 text-left font-mono font-bold text-sm text-slate-100 hover:text-sky-400 transition-colors"
                  >
                    {#if expandedAccordions[c.id]}
                      <ChevronDown class="w-4 h-4 text-sky-400 shrink-0" />
                    {:else}
                      <ChevronRight class="w-4 h-4 text-slate-500 shrink-0" />
                    {/if}

                    <span class="w-2.5 h-2.5 rounded-full shrink-0
                      {c.state === 'running' ? 'bg-emerald-400 shadow-sm shadow-emerald-400' : 'bg-rose-500'}"
                    ></span>

                    <span class="text-slate-100 font-bold">
                      {c.names?.[0]?.replace(/^\//, '') || c.id.substring(0, 12)}
                    </span>

                    <span class="text-xs text-slate-500 font-mono font-normal truncate max-w-xs hidden sm:inline">
                      ({c.image})
                    </span>
                  </button>

                  <div class="flex items-center gap-2">
                    <span class="text-[11px] font-mono uppercase px-2.5 py-1 rounded-full font-semibold
                      {c.state === 'running' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}"
                    >
                      {c.state}
                    </span>

                    <button
                      onclick={() => selectContainer(c.id)}
                      title="Focus Drill-Down Log Stream"
                      class="p-2 rounded-xl bg-slate-800 hover:bg-sky-950/50 text-slate-300 hover:text-sky-400 border border-slate-700/60 transition-all flex items-center gap-1.5 text-xs font-semibold"
                    >
                      <Maximize2 class="w-3.5 h-3.5 text-sky-400" />
                      <span class="hidden sm:inline">Focus Stream</span>
                    </button>
                  </div>
                </div>

                <!-- Accordion Body: Log Viewer -->
                {#if expandedAccordions[c.id]}
                  <div class="p-2 bg-slate-950/50">
                    <LogViewer containerId={c.id} />
                  </div>
                {/if}
              </div>
            {/each}
          {/if}
        </div>
      {/if}
    </div>
  </div>
</div>
