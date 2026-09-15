<script>
  import { performContainerAction, formatBytes, formatUptime } from '../api.js';
  import { addToast } from '../stores.js';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import ContainerDetailModal from './ContainerDetailModal.svelte';
  import { Search, Play, Square, RotateCw, Trash2, Terminal, Box, ChevronLeft, ChevronRight } from 'lucide-svelte';

  export let containers = [];
  export let loading = false;

  let searchQuery = '';
  let statusFilter = 'all'; // 'all' | 'running' | 'stopped'
  let selectedContainer = null;
  let actionLoading = {};

  // Pagination state
  let currentPage = 1;
  let pageSize = 10;

  $: filteredContainers = containers.filter(c => {
    // Filter status
    if (statusFilter === 'running' && c.state !== 'running') return false;
    if (statusFilter === 'stopped' && c.state === 'running') return false;

    // Search query
    if (!searchQuery) return true;
    const q = searchQuery.toLowerCase();
    const name = (c.names?.[0] || '').toLowerCase();
    const id = (c.id || '').toLowerCase();
    const img = (c.image || '').toLowerCase();

    return name.includes(q) || id.includes(q) || img.includes(q);
  });

  $: totalPages = Math.ceil(filteredContainers.length / pageSize) || 1;
  $: if (currentPage > totalPages) currentPage = totalPages;

  $: paginatedContainers = filteredContainers.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );

  async function handleAction(containerId, action, force = false) {
    actionLoading[containerId] = action;
    try {
      await performContainerAction(containerId, action, force);
      addToast(`Container ${action}ed successfully!`, 'success');
    } catch (err) {
      addToast(`Failed to ${action} container: ${err.message}`, 'error');
    } finally {
      delete actionLoading[containerId];
      actionLoading = { ...actionLoading };
    }
  }

  function openDetail(container) {
    selectedContainer = container;
  }
</script>

<div class="space-y-6">
  <!-- Header & Toolbar -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-white flex items-center gap-3">
        <Box class="w-7 h-7 text-sky-400" /> Containers
      </h2>
      <p class="text-sm text-slate-400 mt-1">Manage and inspect Docker containers in real-time</p>
    </div>

    <!-- Filters & Search -->
    <div class="flex flex-wrap items-center gap-3">
      <!-- Status Filter Tabs -->
      <div class="flex bg-slate-900/80 p-1 rounded-xl border border-slate-800">
        <button
          onclick={() => { statusFilter = 'all'; currentPage = 1; }}
          class="px-3 py-1.5 rounded-lg text-xs font-semibold transition-all {statusFilter === 'all' ? 'bg-sky-500 text-white shadow-md' : 'text-slate-400 hover:text-white'}"
        >
          All ({containers.length})
        </button>
        <button
          onclick={() => { statusFilter = 'running'; currentPage = 1; }}
          class="px-3 py-1.5 rounded-lg text-xs font-semibold transition-all {statusFilter === 'running' ? 'bg-emerald-500 text-white shadow-md' : 'text-slate-400 hover:text-white'}"
        >
          Running ({containers.filter(c => c.state === 'running').length})
        </button>
        <button
          onclick={() => { statusFilter = 'stopped'; currentPage = 1; }}
          class="px-3 py-1.5 rounded-lg text-xs font-semibold transition-all {statusFilter === 'stopped' ? 'bg-rose-500 text-white shadow-md' : 'text-slate-400 hover:text-white'}"
        >
          Stopped ({containers.filter(c => c.state !== 'running').length})
        </button>
      </div>

      <!-- Search Input -->
      <div class="relative w-64">
        <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
        <input
          type="text"
          bind:value={searchQuery}
          oninput={() => currentPage = 1}
          placeholder="Search containers..."
          class="w-full bg-slate-900/80 border border-slate-800 rounded-xl pl-9 pr-4 py-2 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-sky-500 transition-all"
        />
      </div>
    </div>
  </div>

  <!-- Container Table / Skeleton -->
  {#if loading && containers.length === 0}
    <SkeletonLoader type="table" rows={5} />
  {:else if filteredContainers.length === 0}
    <div class="glass-card rounded-2xl p-12 text-center space-y-3">
      <Box class="w-12 h-12 text-slate-600 mx-auto" />
      <h3 class="text-base font-semibold text-slate-300">No Containers Found</h3>
      <p class="text-xs text-slate-500">No containers match the current filter or search criteria.</p>
    </div>
  {:else}
    <div class="glass-card rounded-2xl border border-slate-800/80 overflow-hidden shadow-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b border-slate-800/80 bg-slate-900/60 text-slate-400 uppercase tracking-wider font-semibold">
              <th class="py-4 px-5">Name & Image</th>
              <th class="py-4 px-4">Status</th>
              <th class="py-4 px-4">CPU %</th>
              <th class="py-4 px-4">RAM Usage</th>
              <th class="py-4 px-4">Network I/O</th>
              <th class="py-4 px-4">Uptime</th>
              <th class="py-4 px-5 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/50 text-slate-300 font-medium">
            {#each paginatedContainers as c (c.id)}
              <tr class="hover:bg-slate-800/30 transition-colors group">
                <!-- Name & Image -->
                <td class="py-4 px-5">
                  <div class="flex items-center gap-3">
                    <span class="w-2.5 h-2.5 rounded-full shrink-0
                      {c.state === 'running' ? 'bg-emerald-400 shadow-sm shadow-emerald-400' : 'bg-rose-500'}"></span>
                    <div>
                      <button
                        onclick={() => openDetail(c)}
                        class="font-bold text-sm text-slate-100 hover:text-sky-400 transition-colors text-left font-mono"
                      >
                        {c.names?.[0]?.replace(/^\//, '') || c.id.substring(0, 12)}
                      </button>
                      <div class="text-[11px] text-slate-500 font-mono mt-0.5 truncate max-w-xs">{c.image}</div>
                    </div>
                  </div>
                </td>

                <!-- Status Badge -->
                <td class="py-4 px-4 whitespace-nowrap">
                  <span class="px-2.5 py-1 rounded-full text-[11px] font-semibold uppercase tracking-wide
                    {c.state === 'running' ? 'bg-emerald-500/10 text-emerald-400 border border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border border-rose-500/20'}">
                    {c.state}
                  </span>
                </td>

                <!-- CPU % -->
                <td class="py-4 px-4 font-mono font-semibold">
                  {#if c.state === 'running'}
                    <span class="text-sky-400">{c.cpu_usage?.toFixed(1) || 0}%</span>
                  {:else}
                    <span class="text-slate-600">-</span>
                  {/if}
                </td>

                <!-- RAM Usage -->
                <td class="py-4 px-4 font-mono font-semibold">
                  {#if c.state === 'running'}
                    <span class="text-emerald-400">{formatBytes(c.ram_usage)}</span>
                  {:else}
                    <span class="text-slate-600">-</span>
                  {/if}
                </td>

                <!-- Network I/O -->
                <td class="py-4 px-4 font-mono text-slate-400">
                  {#if c.state === 'running'}
                    <span>↓ {formatBytes(c.net_input)} / ↑ {formatBytes(c.net_output)}</span>
                  {:else}
                    <span class="text-slate-600">-</span>
                  {/if}
                </td>

                <!-- Uptime -->
                <td class="py-4 px-4 whitespace-nowrap text-slate-400 font-mono">
                  {formatUptime(c.created)}
                </td>

                <!-- Actions -->
                <td class="py-4 px-5 text-right whitespace-nowrap">
                  <div class="flex items-center justify-end gap-1.5">
                    {#if c.state === 'running'}
                      <button
                        disabled={actionLoading[c.id]}
                        onclick={() => handleAction(c.id, 'stop')}
                        title="Stop Container"
                        class="p-2 rounded-xl bg-slate-800 hover:bg-amber-950/40 text-slate-300 hover:text-amber-400 border border-slate-700/60 transition-all disabled:opacity-50"
                      >
                        <Square class="w-3.5 h-3.5 fill-current" />
                      </button>
                      <button
                        disabled={actionLoading[c.id]}
                        onclick={() => handleAction(c.id, 'restart')}
                        title="Restart Container"
                        class="p-2 rounded-xl bg-slate-800 hover:bg-sky-950/40 text-slate-300 hover:text-sky-400 border border-slate-700/60 transition-all disabled:opacity-50"
                      >
                        <RotateCw class="w-3.5 h-3.5" />
                      </button>
                    {:else}
                      <button
                        disabled={actionLoading[c.id]}
                        onclick={() => handleAction(c.id, 'start')}
                        title="Start Container"
                        class="p-2 rounded-xl bg-slate-800 hover:bg-emerald-950/40 text-slate-300 hover:text-emerald-400 border border-slate-700/60 transition-all disabled:opacity-50"
                      >
                        <Play class="w-3.5 h-3.5 fill-current" />
                      </button>
                    {/if}

                    <button
                      disabled={actionLoading[c.id]}
                      onclick={() => openDetail(c)}
                      title="Container Details & Logs"
                      class="p-2 rounded-xl bg-slate-800 hover:bg-indigo-950/40 text-slate-300 hover:text-indigo-400 border border-slate-700/60 transition-all"
                    >
                      <Terminal class="w-3.5 h-3.5" />
                    </button>

                    <button
                      disabled={actionLoading[c.id]}
                      onclick={() => handleAction(c.id, 'remove', true)}
                      title="Remove Container"
                      class="p-2 rounded-xl bg-slate-800 hover:bg-rose-950/40 text-slate-300 hover:text-rose-400 border border-slate-700/60 transition-all disabled:opacity-50"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Pagination Controls Bar -->
      <div class="flex flex-wrap items-center justify-between p-4 border-t border-slate-800 bg-slate-900/40 text-xs text-slate-400 gap-4">
        <div class="flex items-center gap-2">
          <span>Containers per page:</span>
          <select
            bind:value={pageSize}
            onchange={() => currentPage = 1}
            class="bg-slate-950 border border-slate-800 rounded-lg px-2 py-1 text-slate-200 text-xs focus:outline-none"
          >
            <option value={5}>5</option>
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
          </select>
          <span class="text-slate-500 font-mono ml-2">
            Showing {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, filteredContainers.length)} of {filteredContainers.length}
          </span>
        </div>

        <div class="flex items-center gap-2">
          <button
            disabled={currentPage === 1}
            onclick={() => currentPage--}
            class="p-1.5 rounded-lg bg-slate-800 border border-slate-700/60 text-slate-300 hover:text-white disabled:opacity-40 disabled:cursor-not-allowed transition-all"
          >
            <ChevronLeft class="w-4 h-4" />
          </button>
          <span class="font-mono text-slate-200 font-semibold px-2">Page {currentPage} of {totalPages}</span>
          <button
            disabled={currentPage >= totalPages}
            onclick={() => currentPage++}
            class="p-1.5 rounded-lg bg-slate-800 border border-slate-700/60 text-slate-300 hover:text-white disabled:opacity-40 disabled:cursor-not-allowed transition-all"
          >
            <ChevronRight class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

{#if selectedContainer}
  <ContainerDetailModal
    container={selectedContainer}
    onClose={() => selectedContainer = null}
  />
{/if}
