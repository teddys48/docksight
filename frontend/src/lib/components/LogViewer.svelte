<script>
  import { onMount, onDestroy, tick } from 'svelte';
  import { Play, Pause, Trash2, Search, ArrowDown, Clock, Terminal } from 'lucide-svelte';

  export let containerId;

  let logs = [];
  let isPaused = false;
  let autoScroll = true;
  let showTimestamps = true;
  let showStdout = true;
  let showStderr = true;
  let tailLines = '100';
  let searchQuery = '';

  let eventSource = null;
  let logContainerEl;

  $: if (containerId) {
    connectSSE();
  }

  function connectSSE() {
    if (eventSource) {
      eventSource.close();
    }
    logs = [];

    const url = `/api/sse/logs?id=${containerId}&tail=${tailLines}&stdout=${showStdout}&stderr=${showStderr}&timestamps=${showTimestamps}`;
    eventSource = new EventSource(url);

    eventSource.onmessage = async (e) => {
      if (isPaused) return;
      try {
        const payload = JSON.parse(e.data);
        if (payload.log) {
          logs = [...logs, payload.log];
          if (logs.length > 2000) {
            logs = logs.slice(-1500); // Prevent memory bloat
          }
          if (autoScroll) {
            await tick();
            scrollToBottom();
          }
        }
      } catch (err) {
        console.error('Log parse error:', err);
      }
    };

    eventSource.onerror = (err) => {
      console.warn('SSE log connection error:', err);
    };
  }

  function togglePause() {
    isPaused = !isPaused;
  }

  function clearLogs() {
    logs = [];
  }

  function scrollToBottom() {
    if (logContainerEl) {
      logContainerEl.scrollTop = logContainerEl.scrollHeight;
    }
  }

  function applyFilterChange() {
    connectSSE();
  }

  $: filteredLogs = logs.filter(line => {
    if (!searchQuery) return true;
    return line.toLowerCase().includes(searchQuery.toLowerCase());
  });

  onMount(() => {
    connectSSE();
  });

  onDestroy(() => {
    if (eventSource) {
      eventSource.close();
    }
  });
</script>

<div class="flex flex-col h-[520px] rounded-2xl bg-slate-950 border border-slate-800 overflow-hidden shadow-2xl">
  <!-- Log Toolbar -->
  <div class="flex flex-wrap items-center justify-between p-3 bg-slate-900 border-b border-slate-800 gap-3 text-xs">
    <!-- Left Controls -->
    <div class="flex items-center gap-2">
      <div class="flex items-center gap-1 bg-slate-950 px-2.5 py-1.5 rounded-lg border border-slate-800 text-slate-300">
        <Terminal class="w-3.5 h-3.5 text-sky-400" />
        <span class="font-mono font-bold text-slate-200">Logs</span>
      </div>

      <button
        onclick={togglePause}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium
          {isPaused ? 'bg-amber-500/20 border-amber-500/40 text-amber-300 hover:bg-amber-500/30' : 'bg-slate-800 border-slate-700 text-slate-300 hover:text-white'}"
      >
        {#if isPaused}
          <Play class="w-3.5 h-3.5 fill-current" />
          <span>Resume</span>
        {:else}
          <Pause class="w-3.5 h-3.5" />
          <span>Pause</span>
        {/if}
      </button>

      <button
        onclick={clearLogs}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-rose-950/40 hover:text-rose-300 border border-slate-700 text-slate-300 transition-all font-medium"
      >
        <Trash2 class="w-3.5 h-3.5" />
        <span>Clear</span>
      </button>

      <button
        onclick={() => { autoScroll = !autoScroll; if (autoScroll) scrollToBottom(); }}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium
          {autoScroll ? 'bg-sky-500/20 border-sky-500/40 text-sky-300' : 'bg-slate-800 border-slate-700 text-slate-400'}"
      >
        <ArrowDown class="w-3.5 h-3.5" />
        <span>Auto-scroll</span>
      </button>

      <button
        onclick={() => { showTimestamps = !showTimestamps; applyFilterChange(); }}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium
          {showTimestamps ? 'bg-indigo-500/20 border-indigo-500/40 text-indigo-300' : 'bg-slate-800 border-slate-700 text-slate-400'}"
      >
        <Clock class="w-3.5 h-3.5" />
        <span>Timestamps</span>
      </button>
    </div>

    <!-- Right Controls: Filter & Search -->
    <div class="flex items-center gap-2">
      <!-- stdout / stderr toggle -->
      <div class="flex items-center bg-slate-950 p-0.5 rounded-lg border border-slate-800 font-mono">
        <button
          onclick={() => { showStdout = !showStdout; applyFilterChange(); }}
          class="px-2 py-1 rounded-md text-[11px] font-semibold transition-all {showStdout ? 'bg-emerald-500/20 text-emerald-400' : 'text-slate-500'}"
        >
          STDOUT
        </button>
        <button
          onclick={() => { showStderr = !showStderr; applyFilterChange(); }}
          class="px-2 py-1 rounded-md text-[11px] font-semibold transition-all {showStderr ? 'bg-rose-500/20 text-rose-400' : 'text-slate-500'}"
        >
          STDERR
        </button>
      </div>

      <!-- Tail selector -->
      <select
        bind:value={tailLines}
        onchange={applyFilterChange}
        class="bg-slate-950 border border-slate-800 text-slate-300 rounded-lg px-2 py-1.5 font-mono text-xs focus:outline-none"
      >
        <option value="50">50 lines</option>
        <option value="100">100 lines</option>
        <option value="200">200 lines</option>
        <option value="500">500 lines</option>
      </select>

      <!-- Search Input -->
      <div class="relative">
        <Search class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-400" />
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Filter log output..."
          class="bg-slate-950 border border-slate-800 rounded-lg pl-8 pr-3 py-1 text-slate-200 placeholder-slate-500 text-xs focus:outline-none focus:border-sky-500 w-44"
        />
      </div>
    </div>
  </div>

  <!-- Log View Container -->
  <div
    bind:this={logContainerEl}
    class="flex-1 p-4 font-mono text-xs overflow-y-auto space-y-1 select-text bg-slate-950/90 leading-relaxed"
  >
    {#if filteredLogs.length === 0}
      <div class="h-full flex items-center justify-center text-slate-600 italic">
        {searchQuery ? 'No log lines matching search query' : 'Waiting for logs...'}
      </div>
    {:else}
      {#each filteredLogs as line, i}
        <div class="flex items-start hover:bg-slate-900/60 rounded px-1 group">
          <span class="text-slate-600 select-none w-10 text-right pr-3 shrink-0 font-light">{i + 1}</span>
          <span class="break-all
            {line.includes('[STDERR]') ? 'text-rose-400 font-semibold' : 'text-slate-300'}">
            {line}
          </span>
        </div>
      {/each}
    {/if}
  </div>
</div>
