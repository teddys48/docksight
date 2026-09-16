<script>
  import { onMount } from 'svelte';
  import { inspectContainer, formatBytes, formatDate, formatUptime } from '../api.js';
  import LogViewer from './LogViewer.svelte';
  import { X, Info, Activity, Terminal, Code, Cpu, HardDrive, Network, Layers } from 'lucide-svelte';

  export let container = null;
  export let onClose = () => {};

  let activeModalTab = 'overview'; // 'overview' | 'stats' | 'logs' | 'inspect'
  let inspectData = null;
  let loadingInspect = false;

  $: if (container && activeModalTab === 'inspect' && !inspectData) {
    loadInspect();
  }

  async function loadInspect() {
    loadingInspect = true;
    try {
      inspectData = await inspectContainer(container.id);
    } catch (err) {
      console.error('Inspect load error:', err);
    } finally {
      loadingInspect = false;
    }
  }
</script>

{#if container}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    <div class="glass-card w-full max-w-5xl rounded-3xl border border-slate-800 shadow-2xl flex flex-col max-h-[90vh] overflow-hidden">
      <!-- Modal Header -->
      <div class="flex items-center justify-between p-6 border-b border-slate-800">
        <div class="flex items-center gap-3">
          <div class="p-3 rounded-2xl bg-sky-500/10 border border-sky-500/20 text-sky-400">
            <Layers class="w-6 h-6" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <h2 class="text-xl font-bold text-white tracking-tight">
                {container.names?.[0]?.replace(/^\//, '') || container.id.substring(0, 12)}
              </h2>
              <span class="px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wider
                {container.state === 'running' ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' : 'bg-rose-500/20 text-rose-400 border border-rose-500/30'}">
                {container.state}
              </span>
            </div>
            <p class="text-xs text-slate-400 font-mono mt-0.5">{container.image}</p>
          </div>
        </div>

        <button
          onclick={onClose}
          class="p-2 rounded-xl bg-slate-800/80 text-slate-400 hover:text-white hover:bg-slate-700 transition-all border border-slate-700/50"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Modal Nav Tabs -->
      <div class="flex border-b border-slate-800 bg-slate-900/60 px-6 gap-2">
        <button
          onclick={() => activeModalTab = 'overview'}
          class="flex items-center gap-2 py-3 px-4 font-medium text-sm border-b-2 transition-all
            {activeModalTab === 'overview' ? 'border-sky-500 text-sky-400 font-semibold' : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          <Info class="w-4 h-4" />
          <span>Overview</span>
        </button>

        <button
          onclick={() => activeModalTab = 'stats'}
          class="flex items-center gap-2 py-3 px-4 font-medium text-sm border-b-2 transition-all
            {activeModalTab === 'stats' ? 'border-sky-500 text-sky-400 font-semibold' : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          <Activity class="w-4 h-4" />
          <span>Live Stats</span>
        </button>

        <button
          onclick={() => activeModalTab = 'logs'}
          class="flex items-center gap-2 py-3 px-4 font-medium text-sm border-b-2 transition-all
            {activeModalTab === 'logs' ? 'border-sky-500 text-sky-400 font-semibold' : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          <Terminal class="w-4 h-4" />
          <span>Realtime Logs</span>
        </button>

        <button
          onclick={() => activeModalTab = 'inspect'}
          class="flex items-center gap-2 py-3 px-4 font-medium text-sm border-b-2 transition-all
            {activeModalTab === 'inspect' ? 'border-sky-500 text-sky-400 font-semibold' : 'border-transparent text-slate-400 hover:text-slate-200'}"
        >
          <Code class="w-4 h-4" />
          <span>Inspect</span>
        </button>
      </div>

      <!-- Modal Body -->
      <div class="flex-1 overflow-y-auto p-6 space-y-6">
        {#if activeModalTab === 'overview'}
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Basic Info Card -->
            <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
              <h3 class="font-bold text-sm text-slate-300 uppercase tracking-wider">General Information</h3>
              <div class="space-y-3 text-xs">
                <div>
                  <span class="text-slate-500 block">Full ID</span>
                  <span class="font-mono text-slate-200 break-all select-all">{container.id}</span>
                </div>
                <div>
                  <span class="text-slate-500 block">Image</span>
                  <span class="font-mono text-sky-400">{container.image}</span>
                </div>
                <div>
                  <span class="text-slate-500 block">Command</span>
                  <span class="font-mono text-slate-300">{container.command || 'N/A'}</span>
                </div>
                <div>
                  <span class="text-slate-500 block">Created</span>
                  <span class="text-slate-300">{formatDate(container.created)} ({formatUptime(container.created)} ago)</span>
                </div>
                <div>
                  <span class="text-slate-500 block">Status</span>
                  <span class="text-slate-200 font-medium">{container.status}</span>
                </div>
              </div>
            </div>

            <!-- Ports Mapping Card -->
            <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-4">
              <h3 class="font-bold text-sm text-slate-300 uppercase tracking-wider">Port Forwarding</h3>
              {#if container.ports && container.ports.length > 0}
                <div class="space-y-2">
                  {#each container.ports as p}
                    <div class="flex items-center justify-between p-2.5 rounded-xl bg-slate-950 border border-slate-800/80 font-mono text-xs">
                      <span class="text-slate-400">{p.private_port}/{p.type}</span>
                      <span class="text-emerald-400 font-semibold">{p.public_port ? `${p.ip || '0.0.0.0'}:${p.public_port}` : 'Not bound'}</span>
                    </div>
                  {/each}
                </div>
              {:else}
                <p class="text-xs text-slate-500 italic">No public ports exposed</p>
              {/if}
            </div>
          </div>
        {:else if activeModalTab === 'stats'}
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- CPU Gauge -->
            <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center gap-2">
                  <Cpu class="w-4 h-4 text-sky-400" /> CPU Usage
                </span>
                <span class="font-mono font-bold text-sky-400">{container.cpu_usage?.toFixed(1) || 0}%</span>
              </div>
              <div class="w-full h-3 bg-slate-800 rounded-full overflow-hidden">
                <div
                  class="h-full bg-gradient-to-r from-sky-500 to-blue-600 transition-all duration-500"
                  style="width: {Math.min(container.cpu_usage || 0, 100)}%;"
                ></div>
              </div>
            </div>

            <!-- Memory Gauge -->
            <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center gap-2">
                  <HardDrive class="w-4 h-4 text-emerald-400" /> RAM Usage
                </span>
                <span class="font-mono font-bold text-emerald-400">
                  {formatBytes(container.ram_usage)} ({container.ram_percent?.toFixed(1) || 0}%)
                </span>
              </div>
              <div class="w-full h-3 bg-slate-800 rounded-full overflow-hidden">
                <div
                  class="h-full bg-gradient-to-r from-emerald-500 to-teal-500 transition-all duration-500"
                  style="width: {Math.min(container.ram_percent || 0, 100)}%;"
                ></div>
              </div>
            </div>

            <!-- Network I/O -->
            <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-3">
              <span class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center gap-2">
                <Network class="w-4 h-4 text-indigo-400" /> Network I/O
              </span>
              <div class="grid grid-cols-2 gap-4 pt-2">
                <div class="p-3 rounded-xl bg-slate-950 border border-slate-800">
                  <span class="text-slate-500 text-xs block">Received (RX)</span>
                  <span class="font-mono font-bold text-indigo-400 text-sm">{formatBytes(container.net_input)}</span>
                </div>
                <div class="p-3 rounded-xl bg-slate-950 border border-slate-800">
                  <span class="text-slate-500 text-xs block">Transmitted (TX)</span>
                  <span class="font-mono font-bold text-indigo-400 text-sm">{formatBytes(container.net_output)}</span>
                </div>
              </div>
            </div>

            <!-- Block I/O -->
            <div class="p-5 rounded-2xl bg-slate-900/60 border border-slate-800 space-y-3">
              <span class="text-xs font-bold uppercase tracking-wider text-slate-400 flex items-center gap-2">
                <HardDrive class="w-4 h-4 text-amber-400" /> Block I/O (Disk)
              </span>
              <div class="grid grid-cols-2 gap-4 pt-2">
                <div class="p-3 rounded-xl bg-slate-950 border border-slate-800">
                  <span class="text-slate-500 text-xs block">Read</span>
                  <span class="font-mono font-bold text-amber-400 text-sm">{formatBytes(container.block_read)}</span>
                </div>
                <div class="p-3 rounded-xl bg-slate-950 border border-slate-800">
                  <span class="text-slate-500 text-xs block">Written</span>
                  <span class="font-mono font-bold text-amber-400 text-sm">{formatBytes(container.block_write)}</span>
                </div>
              </div>
            </div>
          </div>
        {:else if activeModalTab === 'logs'}
          <LogViewer containerId={container.id} />
        {:else if activeModalTab === 'inspect'}
          <div class="inspect-container p-4 rounded-2xl border font-mono text-xs overflow-x-auto max-h-[500px]">
            {#if loadingInspect}
              <div class="p-8 text-center text-slate-500 italic">Loading inspect metadata...</div>
            {:else if inspectData}
              <pre class="inspect-json-text leading-relaxed">{JSON.stringify(inspectData, null, 2)}</pre>
            {:else}
              <div class="p-8 text-center text-rose-400">Failed to load container inspection.</div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}
