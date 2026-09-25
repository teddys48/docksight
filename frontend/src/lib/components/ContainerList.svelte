<script>
  import { performContainerAction, formatBytes, formatUptime } from '../api.js';
  import { addToast, selectedLogContainerId, activeTab } from '../stores.js';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import ContainerDetailModal from './ContainerDetailModal.svelte';
  import { Search, Play, Square, RotateCw, Trash2, Terminal, Box, ChevronLeft, ChevronRight, ArrowUp, ArrowDown, ArrowUpDown } from 'lucide-svelte';

  export let containers = [];
  export let loading = false;

  let searchQuery = '';
  let statusFilter = 'all'; // 'all' | 'running' | 'stopped'
  let selectedContainer = null;
  let actionLoading = {};

  // Sorting state
  let sortField = 'name'; // 'name' | 'status' | 'cpu' | 'ram' | 'network' | 'uptime'
  let sortDirection = 'asc'; // 'asc' | 'desc'

  // Pagination state
  let currentPage = 1;
  let pageSize = 10;

  function toggleSort(field) {
    if (sortField === field) {
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      sortField = field;
      sortDirection = (field === 'name' || field === 'status') ? 'asc' : 'desc';
    }
    currentPage = 1;
  }

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

  $: sortedContainers = [...filteredContainers].sort((a, b) => {
    let valA, valB;
    if (sortField === 'name') {
      valA = (a.names?.[0] || a.id).toLowerCase();
      valB = (b.names?.[0] || b.id).toLowerCase();
    } else if (sortField === 'status') {
      valA = (a.state || '').toLowerCase();
      valB = (b.state || '').toLowerCase();
    } else if (sortField === 'cpu') {
      valA = a.cpu_usage || 0;
      valB = b.cpu_usage || 0;
    } else if (sortField === 'ram') {
      valA = a.ram_usage || 0;
      valB = b.ram_usage || 0;
    } else if (sortField === 'network') {
      valA = (a.net_input || 0) + (a.net_output || 0);
      valB = (b.net_input || 0) + (b.net_output || 0);
    } else if (sortField === 'uptime') {
      valA = a.created || 0;
      valB = b.created || 0;
    }

    if (valA < valB) return sortDirection === 'asc' ? -1 : 1;
    if (valA > valB) return sortDirection === 'asc' ? 1 : -1;
    return 0;
  });

  $: totalPages = Math.ceil(sortedContainers.length / pageSize) || 1;
  $: if (currentPage > totalPages) currentPage = totalPages;

  $: paginatedContainers = sortedContainers.slice(
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

  function navigateToLogs(containerId) {
    selectedLogContainerId.set(containerId);
    activeTab.set('logs');
  }
</script>

<div class="space-y-6">
  <!-- Header & Toolbar -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-white flex items-center gap-3">
        <Box class="w-7 h-7 text-sky-400" /> Containers
      </h2>
      <p class="text-sm text-slate-400 mt-1">Manage, sort, and inspect Docker containers in real-time</p>
    </div>

    <!-- Filters & Search & Sort -->
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

      <!-- Quick Sort Selector -->
      <select
        bind:value={sortField}
        onchange={() => { sortDirection = (sortField === 'name' || sortField === 'status') ? 'asc' : 'desc'; currentPage = 1; }}
        class="bg-slate-900/80 border border-slate-800 rounded-xl px-3 py-2 text-xs text-slate-300 focus:outline-none focus:border-sky-500"
      >
        <option value="name">Sort: Name</option>
        <option value="status">Sort: Status</option>
        <option value="cpu">Sort: CPU %</option>
        <option value="ram">Sort: RAM Usage</option>
        <option value="network">Sort: Network I/O</option>
        <option value="uptime">Sort: Uptime</option>
      </select>

      <!-- Direction Toggle -->
      <button
        onclick={() => sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'}
        title="Toggle sort direction"
        class="p-2 rounded-xl bg-slate-900/80 border border-slate-800 text-slate-300 hover:text-white transition-colors"
      >
        {#if sortDirection === 'asc'}
          <ArrowUp class="w-4 h-4 text-sky-400" />
        {:else}
          <ArrowDown class="w-4 h-4 text-sky-400" />
        {/if}
      </button>

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
            <tr class="border-b border-slate-800/80 bg-slate-900/60 text-slate-400 uppercase tracking-wider font-semibold select-none">
              <!-- Name Header -->
              <th
                onclick={() => toggleSort('name')}
                class="py-4 px-5 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Name & Image</span>
                  {#if sortField === 'name'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- Status Header -->
              <th
                onclick={() => toggleSort('status')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Status</span>
                  {#if sortField === 'status'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- CPU Header -->
              <th
                onclick={() => toggleSort('cpu')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>CPU %</span>
                  {#if sortField === 'cpu'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- RAM Header -->
              <th
                onclick={() => toggleSort('ram')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>RAM Usage</span>
                  {#if sortField === 'ram'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- Network Header -->
              <th
                onclick={() => toggleSort('network')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Network I/O</span>
                  {#if sortField === 'network'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- Uptime Header -->
              <th
                onclick={() => toggleSort('uptime')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Uptime</span>
                  {#if sortField === 'uptime'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-sky-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

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
                      onclick={() => navigateToLogs(c.id)}
                      title="Drill-down Logs"
                      class="p-2 rounded-xl bg-slate-800 hover:bg-sky-950/40 text-slate-300 hover:text-sky-400 border border-slate-700/60 transition-all"
                    >
                      <Terminal class="w-3.5 h-3.5 text-sky-400" />
                    </button>

                    <button
                      disabled={actionLoading[c.id]}
                      onclick={() => openDetail(c)}
                      title="Container Details & JSON"
                      class="p-2 rounded-xl bg-slate-800 hover:bg-indigo-950/40 text-slate-300 hover:text-indigo-400 border border-slate-700/60 transition-all"
                    >
                      <Box class="w-3.5 h-3.5 text-indigo-400" />
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
            Showing {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, sortedContainers.length)} of {sortedContainers.length}
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
