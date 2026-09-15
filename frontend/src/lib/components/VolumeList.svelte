<script>
  import { onMount } from 'svelte';
  import { fetchVolumes, formatDate } from '../api.js';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import { Database, Search, ChevronLeft, ChevronRight, HardDrive } from 'lucide-svelte';

  let volumes = [];
  let loading = true;
  let searchQuery = '';

  // Pagination states
  let currentPage = 1;
  let pageSize = 10;

  onMount(async () => {
    await loadVolumes();
  });

  async function loadVolumes() {
    loading = true;
    try {
      volumes = await fetchVolumes();
    } catch (err) {
      console.error('Failed to load volumes:', err);
      volumes = [];
    } finally {
      loading = false;
    }
  }

  $: filteredVolumes = volumes.filter(v => {
    if (!searchQuery) return true;
    const q = searchQuery.toLowerCase();
    const name = (v.name || '').toLowerCase();
    const mp = (v.mountpoint || '').toLowerCase();
    const driver = (v.driver || '').toLowerCase();
    return name.includes(q) || mp.includes(q) || driver.includes(q);
  });

  $: totalPages = Math.ceil(filteredVolumes.length / pageSize) || 1;
  $: if (currentPage > totalPages) currentPage = totalPages;

  $: paginatedVolumes = filteredVolumes.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );
</script>

<div class="space-y-6">
  <!-- Header & Search -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-white flex items-center gap-3">
        <Database class="w-7 h-7 text-purple-400" /> Docker Volumes
      </h2>
      <p class="text-sm text-slate-400 mt-1">Explore host volume mounts and persistent storage</p>
    </div>

    <div class="flex flex-wrap items-center gap-3">
      <!-- Search Input -->
      <div class="relative w-64">
        <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
        <input
          type="text"
          bind:value={searchQuery}
          oninput={() => currentPage = 1}
          placeholder="Search volume name..."
          class="w-full bg-slate-900/80 border border-slate-800 rounded-xl pl-9 pr-4 py-2 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-purple-500 transition-all"
        />
      </div>
    </div>
  </div>

  <!-- Table / Skeleton -->
  {#if loading}
    <SkeletonLoader type="table" rows={5} />
  {:else if filteredVolumes.length === 0}
    <div class="glass-card rounded-2xl p-12 text-center space-y-3">
      <Database class="w-12 h-12 text-slate-600 mx-auto" />
      <h3 class="text-base font-semibold text-slate-300">No Volumes Found</h3>
      <p class="text-xs text-slate-500">No Docker volumes match your search criteria.</p>
    </div>
  {:else}
    <div class="glass-card rounded-2xl border border-slate-800/80 overflow-hidden shadow-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b border-slate-800/80 bg-slate-900/60 text-slate-400 uppercase tracking-wider font-semibold">
              <th class="py-4 px-5">Volume Name</th>
              <th class="py-4 px-4">Driver</th>
              <th class="py-4 px-4">Mount Point</th>
              <th class="py-4 px-4">Scope</th>
              <th class="py-4 px-5 text-right">Created At</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/50 text-slate-300 font-medium">
            {#each paginatedVolumes as v (v.name)}
              <tr class="hover:bg-slate-800/30 transition-colors">
                <td class="py-4 px-5 font-mono font-bold text-slate-100 flex items-center gap-2">
                  <Database class="w-4 h-4 text-purple-400 shrink-0" />
                  <span class="truncate max-w-xs">{v.name}</span>
                </td>
                <td class="py-4 px-4 font-mono">
                  <span class="px-2.5 py-1 rounded-full text-[11px] font-semibold bg-purple-500/10 text-purple-300 border border-purple-500/20">
                    {v.driver || 'local'}
                  </span>
                </td>
                <td class="py-4 px-4 font-mono text-slate-400 select-all truncate max-w-md">
                  {v.mountpoint}
                </td>
                <td class="py-4 px-4 font-mono text-slate-400">
                  {v.scope || 'local'}
                </td>
                <td class="py-4 px-5 text-right font-mono text-slate-400 whitespace-nowrap">
                  {formatDate(v.created)}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar -->
      <div class="flex flex-wrap items-center justify-between p-4 border-t border-slate-800 bg-slate-900/40 text-xs text-slate-400 gap-4">
        <div class="flex items-center gap-2">
          <span>Rows per page:</span>
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
            Showing {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, filteredVolumes.length)} of {filteredVolumes.length}
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
