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

  let connectionStatus = 'connecting'; // 'connecting' | 'connected' | 'ended' | 'error'
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
    connectionStatus = 'connecting';

    const url = `/api/sse/logs?id=${containerId}&tail=${tailLines}&stdout=${showStdout}&stderr=${showStderr}&timestamps=${showTimestamps}`;
    eventSource = new EventSource(url);

    eventSource.onopen = () => {
      connectionStatus = 'connected';
    };

    eventSource.onmessage = async (e) => {
      if (isPaused) return;
      try {
        const payload = JSON.parse(e.data);
        if (payload.log) {
          connectionStatus = 'connected';
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

    eventSource.addEventListener('end', () => {
      connectionStatus = 'ended';
      if (eventSource) eventSource.close();
    });

    eventSource.onerror = (err) => {
      if (logs.length === 0) {
        connectionStatus = 'error';
      }
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

<div class="terminal-container flex flex-col h-[520px] rounded-2xl border overflow-hidden shadow-2xl">
  <!-- Log Toolbar -->
  <div class="terminal-toolbar flex flex-wrap items-center justify-between p-3 border-b gap-3 text-xs">
    <!-- Left Controls -->
    <div class="flex items-center gap-2">
      <div class="terminal-badge flex items-center gap-1 px-2.5 py-1.5 rounded-lg border text-xs">
        <Terminal class="w-3.5 h-3.5 text-sky-400" />
        <span class="font-mono font-bold">Logs</span>
      </div>

      <button
        onclick={togglePause}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium text-xs
          {isPaused ? 'terminal-btn-amber' : 'terminal-btn-inactive'}"
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
        class="terminal-btn-inactive hover:terminal-btn-rose flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium text-xs"
      >
        <Trash2 class="w-3.5 h-3.5" />
        <span>Clear</span>
      </button>

      <button
        onclick={() => { autoScroll = !autoScroll; if (autoScroll) scrollToBottom(); }}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium text-xs
          {autoScroll ? 'terminal-btn-sky' : 'terminal-btn-inactive'}"
      >
        <ArrowDown class="w-3.5 h-3.5" />
        <span>Auto-scroll</span>
      </button>

      <button
        onclick={() => { showTimestamps = !showTimestamps; applyFilterChange(); }}
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border transition-all font-medium text-xs
          {showTimestamps ? 'terminal-btn-indigo' : 'terminal-btn-inactive'}"
      >
        <Clock class="w-3.5 h-3.5" />
        <span>Timestamps</span>
      </button>
    </div>

    <!-- Right Controls: Filter & Search -->
    <div class="flex items-center gap-2">
      <!-- stdout / stderr toggle -->
      <div class="terminal-badge flex items-center p-0.5 rounded-lg border font-mono">
        <button
          onclick={() => { showStdout = !showStdout; applyFilterChange(); }}
          class="px-2 py-1 rounded-md text-[11px] font-semibold transition-all {showStdout ? 'terminal-btn-emerald' : 'text-slate-500'}"
        >
          STDOUT
        </button>
        <button
          onclick={() => { showStderr = !showStderr; applyFilterChange(); }}
          class="px-2 py-1 rounded-md text-[11px] font-semibold transition-all {showStderr ? 'terminal-btn-rose' : 'text-slate-500'}"
        >
          STDERR
        </button>
      </div>

      <!-- Tail selector -->
      <select
        bind:value={tailLines}
        onchange={applyFilterChange}
        class="terminal-input border rounded-lg px-2 py-1.5 font-mono text-xs focus:outline-none"
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
          class="terminal-input border rounded-lg pl-8 pr-3 py-1 text-xs focus:outline-none w-44"
        />
      </div>
    </div>
  </div>

  <!-- Log View Container -->
  <div
    bind:this={logContainerEl}
    class="terminal-body flex-1 p-4 font-mono text-xs overflow-y-auto space-y-1 select-text leading-relaxed"
  >
    {#if filteredLogs.length === 0}
      <div class="h-full flex flex-col items-center justify-center text-slate-500 italic font-mono text-xs p-6 text-center space-y-2">
        {#if searchQuery}
          <span>No log lines matching "{searchQuery}"</span>
        {:else if connectionStatus === 'connecting'}
          <div class="flex items-center gap-2 text-sky-400 font-semibold not-italic">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-sky-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-sky-500"></span>
            </span>
            <span>Connecting to container log stream...</span>
          </div>
        {:else if connectionStatus === 'ended'}
          <span class="text-slate-400">Container log stream completed / reached end of output</span>
        {:else if connectionStatus === 'error'}
          <span class="text-rose-400">Unable to stream logs (Connection error or container stopped)</span>
        {:else}
          <span>No log output recorded for this container</span>
        {/if}
      </div>
    {:else}
      {#each filteredLogs as line, i}
        <div class="terminal-line-hover flex items-start rounded px-1 group">
          <span class="terminal-line-num select-none w-10 text-right pr-3 shrink-0 font-light">{i + 1}</span>
          <span class="break-all {line.includes('[STDERR]') ? 'terminal-line-stderr' : 'terminal-line-stdout'}">
            {line}
          </span>
        </div>
      {/each}
    {/if}
  </div>
</div>
