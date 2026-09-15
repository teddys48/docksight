<script>
  export let type = 'card'; // 'card' | 'table' | 'chart'
  export let rows = 4;
</script>

{#if type === 'card'}
  <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4 animate-pulse">
    {#each Array(rows) as _}
      <div class="p-6 rounded-2xl bg-slate-800/40 border border-slate-700/40 space-y-4">
        <div class="flex justify-between items-center">
          <div class="h-4 w-24 bg-slate-700/60 rounded"></div>
          <div class="h-8 w-8 bg-slate-700/60 rounded-xl"></div>
        </div>
        <div class="h-8 w-36 bg-slate-700/60 rounded"></div>
        <div class="h-2 w-full bg-slate-700/60 rounded-full"></div>
      </div>
    {/each}
  </div>
{:else if type === 'table'}
  <div class="w-full bg-slate-800/40 rounded-2xl border border-slate-700/40 p-4 space-y-4 animate-pulse">
    <div class="h-10 w-full bg-slate-700/60 rounded-xl"></div>
    {#each Array(rows) as _}
      <div class="h-14 w-full bg-slate-700/30 rounded-xl flex items-center px-4 gap-4">
        <div class="h-4 w-1/4 bg-slate-700/60 rounded"></div>
        <div class="h-4 w-1/6 bg-slate-700/60 rounded"></div>
        <div class="h-4 w-1/6 bg-slate-700/60 rounded"></div>
        <div class="h-4 w-1/4 bg-slate-700/60 rounded"></div>
      </div>
    {/each}
  </div>
{:else if type === 'chart'}
  <div class="h-64 w-full bg-slate-800/40 rounded-2xl border border-slate-700/40 p-6 animate-pulse flex flex-col justify-between">
    <div class="flex justify-between">
      <div class="h-5 w-32 bg-slate-700/60 rounded"></div>
      <div class="h-5 w-20 bg-slate-700/60 rounded"></div>
    </div>
    <div class="h-40 w-full bg-slate-700/20 rounded-xl flex items-end p-4 gap-2">
      <div class="h-1/2 w-full bg-slate-700/40 rounded"></div>
      <div class="h-3/4 w-full bg-slate-700/40 rounded"></div>
      <div class="h-2/3 w-full bg-slate-700/40 rounded"></div>
      <div class="h-full w-full bg-slate-700/40 rounded"></div>
      <div class="h-4/5 w-full bg-slate-700/40 rounded"></div>
    </div>
  </div>
{/if}
