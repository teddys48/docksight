<script>
  import { theme, activeTab, sseConnected } from '../stores.js';
  import { LayoutDashboard, Box, HardDrive, Database, Sun, Moon, Container } from 'lucide-svelte';

  function toggleTheme() {
    theme.update(t => (t === 'dark' ? 'light' : 'dark'));
  }
</script>

<header class="sticky top-0 z-40 w-full glass-card border-b border-slate-800/50 backdrop-blur-md">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
    <!-- Brand Logo -->
    <div class="flex items-center gap-3">
      <div class="p-2 rounded-xl bg-gradient-to-tr from-sky-500 to-blue-600 text-white shadow-lg shadow-sky-500/20">
        <Container class="w-6 h-6" />
      </div>
      <div>
        <h1 class="font-bold text-lg tracking-tight bg-gradient-to-r from-sky-400 to-blue-500 bg-clip-text text-transparent">
          docksight
        </h1>
        <div class="flex items-center gap-1.5 text-xs text-slate-400">
          <span class="relative flex h-2 w-2">
            {#if $sseConnected}
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
            {:else}
              <span class="relative inline-flex rounded-full h-2 w-2 bg-rose-500"></span>
            {/if}
          </span>
          <span>{$sseConnected ? 'Live Realtime' : 'Disconnected'}</span>
        </div>
      </div>
    </div>

    <!-- Navigation Tabs -->
    <nav class="hidden md:flex items-center gap-1 bg-slate-900/60 p-1.5 rounded-xl border border-slate-800">
      <button
        onclick={() => activeTab.set('dashboard')}
        class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200
          {$activeTab === 'dashboard' ? 'bg-sky-500 text-white shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}"
      >
        <LayoutDashboard class="w-4 h-4" />
        <span>Dashboard</span>
      </button>

      <button
        onclick={() => activeTab.set('containers')}
        class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200
          {$activeTab === 'containers' ? 'bg-sky-500 text-white shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}"
      >
        <Box class="w-4 h-4" />
        <span>Containers</span>
      </button>

      <button
        onclick={() => activeTab.set('images')}
        class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200
          {$activeTab === 'images' ? 'bg-sky-500 text-white shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}"
      >
        <HardDrive class="w-4 h-4" />
        <span>Images</span>
      </button>

      <button
        onclick={() => activeTab.set('volumes')}
        class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200
          {$activeTab === 'volumes' ? 'bg-sky-500 text-white shadow-md shadow-sky-500/20' : 'text-slate-400 hover:text-white hover:bg-slate-800/50'}"
      >
        <Database class="w-4 h-4" />
        <span>Volumes</span>
      </button>
    </nav>

    <!-- Theme Toggle -->
    <div class="flex items-center gap-2">
      <button
        onclick={toggleTheme}
        aria-label="Toggle theme"
        class="p-2.5 rounded-xl bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition-all border border-slate-700/50"
      >
        {#if $theme === 'dark'}
          <Sun class="w-5 h-5 text-amber-400" />
        {:else}
          <Moon class="w-5 h-5 text-indigo-400" />
        {/if}
      </button>
    </div>
  </div>

  <!-- Mobile Bottom Nav -->
  <div class="md:hidden flex border-t border-slate-800/60 bg-slate-900/90 p-2 justify-around">
    <button
      onclick={() => activeTab.set('dashboard')}
      class="flex flex-col items-center gap-1 p-2 rounded-lg text-xs font-medium {$activeTab === 'dashboard' ? 'text-sky-400' : 'text-slate-400'}"
    >
      <LayoutDashboard class="w-5 h-5" />
      <span>Dashboard</span>
    </button>
    <button
      onclick={() => activeTab.set('containers')}
      class="flex flex-col items-center gap-1 p-2 rounded-lg text-xs font-medium {$activeTab === 'containers' ? 'text-sky-400' : 'text-slate-400'}"
    >
      <Box class="w-5 h-5" />
      <span>Containers</span>
    </button>
    <button
      onclick={() => activeTab.set('images')}
      class="flex flex-col items-center gap-1 p-2 rounded-lg text-xs font-medium {$activeTab === 'images' ? 'text-sky-400' : 'text-slate-400'}"
    >
      <HardDrive class="w-5 h-5" />
      <span>Images</span>
    </button>
    <button
      onclick={() => activeTab.set('volumes')}
      class="flex flex-col items-center gap-1 p-2 rounded-lg text-xs font-medium {$activeTab === 'volumes' ? 'text-sky-400' : 'text-slate-400'}"
    >
      <Database class="w-5 h-5" />
      <span>Volumes</span>
    </button>
  </div>
</header>
