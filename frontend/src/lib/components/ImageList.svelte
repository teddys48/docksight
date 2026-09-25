<script>
  import { onMount } from 'svelte';
  import { fetchImages, formatBytes, formatDate } from '../api.js';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import { HardDrive, Search, Filter, AlertTriangle, Tag, Clock, ChevronLeft, ChevronRight, ArrowUp, ArrowDown, ArrowUpDown } from 'lucide-svelte';

  let images = [];
  let loading = true;
  let searchQuery = '';
  let filterDangling = false;

  // Sorting state
  let sortField = 'tag'; // 'tag' | 'id' | 'size' | 'created' | 'status'
  let sortDirection = 'asc'; // 'asc' | 'desc'

  // Pagination states
  let currentPage = 1;
  let pageSize = 10;

  onMount(async () => {
    await loadImages();
  });

  async function loadImages() {
    loading = true;
    try {
      images = await fetchImages();
    } catch (err) {
      console.error('Failed to load images:', err);
    } finally {
      loading = false;
    }
  }

  function toggleSort(field) {
    if (sortField === field) {
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      sortField = field;
      sortDirection = (field === 'tag' || field === 'id') ? 'asc' : 'desc';
    }
    currentPage = 1;
  }

  $: filteredImages = images.filter(img => {
    if (filterDangling && !img.is_dangling) return false;

    if (!searchQuery) return true;
    const q = searchQuery.toLowerCase();
    const tags = (img.repo_tags || []).join(' ').toLowerCase();
    const id = (img.id || '').toLowerCase();

    return tags.includes(q) || id.includes(q);
  });

  $: sortedImages = [...filteredImages].sort((a, b) => {
    let valA, valB;
    if (sortField === 'tag') {
      valA = (a.repo_tags || []).join(' ').toLowerCase();
      valB = (b.repo_tags || []).join(' ').toLowerCase();
    } else if (sortField === 'id') {
      valA = (a.id || '').toLowerCase();
      valB = (b.id || '').toLowerCase();
    } else if (sortField === 'size') {
      valA = a.size || 0;
      valB = b.size || 0;
    } else if (sortField === 'created') {
      valA = a.created || 0;
      valB = b.created || 0;
    } else if (sortField === 'status') {
      valA = a.is_dangling ? 1 : 0;
      valB = b.is_dangling ? 1 : 0;
    }

    if (valA < valB) return sortDirection === 'asc' ? -1 : 1;
    if (valA > valB) return sortDirection === 'asc' ? 1 : -1;
    return 0;
  });

  $: totalPages = Math.ceil(sortedImages.length / pageSize) || 1;
  $: if (currentPage > totalPages) currentPage = totalPages;

  $: paginatedImages = sortedImages.slice(
    (currentPage - 1) * pageSize,
    currentPage * pageSize
  );
</script>

<div class="space-y-6">
  <!-- Header & Toolbar -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2 class="text-2xl font-bold tracking-tight text-white flex items-center gap-3">
        <HardDrive class="w-7 h-7 text-indigo-400" /> Docker Images
      </h2>
      <p class="text-sm text-slate-400 mt-1">Explore local images, sizes, creation dates, and dangling tags</p>
    </div>

    <div class="flex flex-wrap items-center gap-3">
      <!-- Dangling Filter Toggle -->
      <button
        onclick={() => { filterDangling = !filterDangling; currentPage = 1; }}
        class="flex items-center gap-2 px-3.5 py-2 rounded-xl text-xs font-semibold border transition-all
          {filterDangling ? 'bg-amber-500/20 border-amber-500/40 text-amber-300' : 'bg-slate-900/80 border-slate-800 text-slate-400 hover:text-white'}"
      >
        <AlertTriangle class="w-4 h-4 text-amber-400" />
        <span>Dangling Only ({images.filter(i => i.is_dangling).length})</span>
      </button>

      <!-- Quick Sort Selector -->
      <select
        bind:value={sortField}
        onchange={() => { sortDirection = (sortField === 'tag' || sortField === 'id') ? 'asc' : 'desc'; currentPage = 1; }}
        class="bg-slate-900/80 border border-slate-800 rounded-xl px-3 py-2 text-xs text-slate-300 focus:outline-none focus:border-indigo-500"
      >
        <option value="tag">Sort: Tag / Repo</option>
        <option value="id">Sort: Image ID</option>
        <option value="size">Sort: Size</option>
        <option value="created">Sort: Created Date</option>
        <option value="status">Sort: Dangling Status</option>
      </select>

      <!-- Direction Toggle -->
      <button
        onclick={() => sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'}
        title="Toggle sort direction"
        class="p-2 rounded-xl bg-slate-900/80 border border-slate-800 text-slate-300 hover:text-white transition-colors"
      >
        {#if sortDirection === 'asc'}
          <ArrowUp class="w-4 h-4 text-indigo-400" />
        {:else}
          <ArrowDown class="w-4 h-4 text-indigo-400" />
        {/if}
      </button>

      <!-- Search Input -->
      <div class="relative w-64">
        <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
        <input
          type="text"
          bind:value={searchQuery}
          oninput={() => currentPage = 1}
          placeholder="Search images or tags..."
          class="w-full bg-slate-900/80 border border-slate-800 rounded-xl pl-9 pr-4 py-2 text-sm text-slate-200 placeholder-slate-500 focus:outline-none focus:border-indigo-500 transition-all"
        />
      </div>
    </div>
  </div>

  <!-- Images Table / Skeleton -->
  {#if loading}
    <SkeletonLoader type="table" rows={4} />
  {:else if filteredImages.length === 0}
    <div class="glass-card rounded-2xl p-12 text-center space-y-3">
      <HardDrive class="w-12 h-12 text-slate-600 mx-auto" />
      <h3 class="text-base font-semibold text-slate-300">No Docker Images Found</h3>
      <p class="text-xs text-slate-500">No images match your search criteria or dangling filter.</p>
    </div>
  {:else}
    <div class="glass-card rounded-2xl border border-slate-800/80 overflow-hidden shadow-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b border-slate-800/80 bg-slate-900/60 text-slate-400 uppercase tracking-wider font-semibold select-none">
              <!-- Tag Header -->
              <th
                onclick={() => toggleSort('tag')}
                class="py-4 px-5 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Repository & Tag</span>
                  {#if sortField === 'tag'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- ID Header -->
              <th
                onclick={() => toggleSort('id')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Image ID</span>
                  {#if sortField === 'id'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- Size Header -->
              <th
                onclick={() => toggleSort('size')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Size</span>
                  {#if sortField === 'size'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- Created Header -->
              <th
                onclick={() => toggleSort('created')}
                class="py-4 px-4 cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center gap-1.5">
                  <span>Created Date</span>
                  {#if sortField === 'created'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>

              <!-- Status Header -->
              <th
                onclick={() => toggleSort('status')}
                class="py-4 px-5 text-right cursor-pointer hover:text-white transition-colors group"
              >
                <div class="flex items-center justify-end gap-1.5">
                  <span>Status</span>
                  {#if sortField === 'status'}
                    {#if sortDirection === 'asc'}
                      <ArrowUp class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {:else}
                      <ArrowDown class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                    {/if}
                  {:else}
                    <ArrowUpDown class="w-3 h-3 text-slate-600 group-hover:text-slate-400 shrink-0 opacity-40" />
                  {/if}
                </div>
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/50 text-slate-300 font-medium">
            {#each paginatedImages as img (img.id)}
              <tr class="hover:bg-slate-800/30 transition-colors">
                <!-- Repository & Tags -->
                <td class="py-4 px-5">
                  <div class="space-y-1">
                    {#if img.repo_tags && img.repo_tags.length > 0}
                      {#each img.repo_tags as tag}
                        <div class="flex items-center gap-2">
                          <Tag class="w-3.5 h-3.5 text-indigo-400 shrink-0" />
                          <span class="font-bold text-slate-100 font-mono text-sm">{tag}</span>
                        </div>
                      {/each}
                    {:else}
                      <span class="text-amber-400 font-mono italic">&lt;none&gt;:&lt;none&gt;</span>
                    {/if}
                  </div>
                </td>

                <!-- ID -->
                <td class="py-4 px-4 font-mono text-slate-400 select-all">
                  {img.id.replace('sha256:', '').substring(0, 12)}
                </td>

                <!-- Size -->
                <td class="py-4 px-4 font-mono font-bold text-sky-400">
                  {formatBytes(img.size)}
                </td>

                <!-- Created Date -->
                <td class="py-4 px-4 text-slate-400 font-mono whitespace-nowrap">
                  {formatDate(img.created)}
                </td>

                <!-- Status / Dangling Badge -->
                <td class="py-4 px-5 text-right whitespace-nowrap">
                  {#if img.is_dangling}
                    <span class="px-2.5 py-1 rounded-full text-[11px] font-semibold bg-amber-500/20 text-amber-300 border border-amber-500/40 inline-flex items-center gap-1.5">
                      <AlertTriangle class="w-3 h-3 text-amber-400" /> Dangling
                    </span>
                  {:else}
                    <span class="px-2.5 py-1 rounded-full text-[11px] font-semibold bg-slate-800 text-slate-400 border border-slate-700/60">
                      Active Tag
                    </span>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Pagination Bar -->
      <div class="flex flex-wrap items-center justify-between p-4 border-t border-slate-800 bg-slate-900/40 text-xs text-slate-400 gap-4">
        <div class="flex items-center gap-2">
          <span>Images per page:</span>
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
            Showing {(currentPage - 1) * pageSize + 1} - {Math.min(currentPage * pageSize, sortedImages.length)} of {sortedImages.length}
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
