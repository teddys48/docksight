<script>
  import { toasts } from '../stores.js';
  import { AlertCircle, CheckCircle2, Info, X } from 'lucide-svelte';
</script>

<div class="fixed top-5 right-5 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
  {#each $toasts as toast (toast.id)}
    <div
      class="pointer-events-auto flex items-center justify-between p-4 rounded-xl shadow-xl border backdrop-blur-md transition-all duration-300 transform translate-y-0
        {toast.type === 'error' ? 'bg-rose-950/80 border-rose-500/50 text-rose-200' : ''}
        {toast.type === 'success' ? 'bg-emerald-950/80 border-emerald-500/50 text-emerald-200' : ''}
        {toast.type === 'info' ? 'bg-slate-900/90 border-slate-700 text-slate-100' : ''}"
    >
      <div class="flex items-center gap-3">
        {#if toast.type === 'error'}
          <AlertCircle class="w-5 h-5 text-rose-400 shrink-0" />
        {:else if toast.type === 'success'}
          <CheckCircle2 class="w-5 h-5 text-emerald-400 shrink-0" />
        {:else}
          <Info class="w-5 h-5 text-sky-400 shrink-0" />
        {/if}
        <span class="text-sm font-medium leading-snug">{toast.message}</span>
      </div>
      <button
        onclick={() => toasts.update(all => all.filter(t => t.id !== toast.id))}
        class="text-slate-400 hover:text-white p-1 rounded-lg transition-colors"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  {/each}
</div>
