<script>
  import { formatBytes } from '../api.js';
  import MetricChart from './MetricChart.svelte';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import { Cpu, HardDrive, Box, Play, Square, Network, LayoutDashboard, ArrowUpRight } from 'lucide-svelte';

  export let systemStats = null;
  export let history = [];
  export let loading = false;
  export let onNavigateToContainers = () => {};

  $: chartLabels = history.map(h => {
    if (!h.timestamp) return '';
    const date = new Date(h.timestamp);
    return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`;
  });

  $: cpuData = history.map(h => h.cpu_usage || 0);
  $: ramData = history.map(h => {
    if (!h.ram_total || h.ram_total === 0) return 0;
    return (h.ram_usage / h.ram_total) * 100;
  });
</script>

<div class="space-y-8">
  <!-- Top Stat Cards -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
    <!-- Host CPU Card -->
    <div class="glass-card p-6 rounded-3xl border border-slate-800/80 shadow-lg space-y-4 hover:border-sky-500/40 transition-all">
      <div class="flex items-center justify-between">
        <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Host CPU</span>
        <div class="p-2.5 rounded-2xl bg-sky-500/10 text-sky-400 border border-sky-500/20">
          <Cpu class="w-5 h-5" />
        </div>
      </div>
      <div>
        <div class="text-3xl font-extrabold tracking-tight text-white font-mono">
          {systemStats?.host_cpu?.toFixed(1) || 0}%
        </div>
        <div class="w-full h-2 bg-slate-800 rounded-full mt-3 overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-sky-500 to-blue-600 transition-all duration-500"
            style="width: {Math.min(systemStats?.host_cpu || 0, 100)}%;"
          ></div>
        </div>
      </div>
    </div>

    <!-- Host RAM Card -->
    <div class="glass-card p-6 rounded-3xl border border-slate-800/80 shadow-lg space-y-4 hover:border-emerald-500/40 transition-all">
      <div class="flex items-center justify-between">
        <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Host RAM</span>
        <div class="p-2.5 rounded-2xl bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
          <HardDrive class="w-5 h-5" />
        </div>
      </div>
      <div>
        <div class="flex items-baseline gap-2">
          <div class="text-3xl font-extrabold tracking-tight text-white font-mono">
            {systemStats?.host_ram?.used_percent?.toFixed(1) || 0}%
          </div>
          <span class="text-xs text-slate-400 font-mono">
            {formatBytes(systemStats?.host_ram?.used)} / {formatBytes(systemStats?.host_ram?.total)}
          </span>
        </div>
        <div class="w-full h-2 bg-slate-800 rounded-full mt-3 overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-emerald-500 to-teal-500 transition-all duration-500"
            style="width: {Math.min(systemStats?.host_ram?.used_percent || 0, 100)}%;"
          ></div>
        </div>
      </div>
    </div>

    <!-- Host Disk Card -->
    <div class="glass-card p-6 rounded-3xl border border-slate-800/80 shadow-lg space-y-4 hover:border-indigo-500/40 transition-all">
      <div class="flex items-center justify-between">
        <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Host Disk</span>
        <div class="p-2.5 rounded-2xl bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
          <HardDrive class="w-5 h-5" />
        </div>
      </div>
      <div>
        <div class="flex items-baseline gap-2">
          <div class="text-3xl font-extrabold tracking-tight text-white font-mono">
            {systemStats?.host_disk?.used_percent?.toFixed(1) || 0}%
          </div>
          <span class="text-xs text-slate-400 font-mono">
            {formatBytes(systemStats?.host_disk?.used)} / {formatBytes(systemStats?.host_disk?.total)}
          </span>
        </div>
        <div class="w-full h-2 bg-slate-800 rounded-full mt-3 overflow-hidden">
          <div
            class="h-full bg-gradient-to-r from-indigo-500 to-purple-600 transition-all duration-500"
            style="width: {Math.min(systemStats?.host_disk?.used_percent || 0, 100)}%;"
          ></div>
        </div>
      </div>
    </div>

    <!-- Container Summary Card -->
    <div class="glass-card p-6 rounded-3xl border border-slate-800/80 shadow-lg space-y-4 hover:border-amber-500/40 transition-all">
      <div class="flex items-center justify-between">
        <span class="text-xs font-bold uppercase tracking-wider text-slate-400">Containers</span>
        <div class="p-2.5 rounded-2xl bg-amber-500/10 text-amber-400 border border-amber-500/20">
          <Box class="w-5 h-5" />
        </div>
      </div>
      <div>
        <div class="text-3xl font-extrabold tracking-tight text-white font-mono">
          {systemStats?.total_containers || 0}
        </div>
        <div class="flex items-center gap-4 mt-3 text-xs font-semibold">
          <span class="flex items-center gap-1 text-emerald-400">
            <Play class="w-3 h-3 fill-current" /> {systemStats?.running_containers || 0} Running
          </span>
          <span class="flex items-center gap-1 text-rose-400">
            <Square class="w-3 h-3 fill-current" /> {systemStats?.stopped_containers || 0} Stopped
          </span>
        </div>
      </div>
    </div>
  </div>

  <!-- Realtime Metric Charts -->
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
    <MetricChart
      title="Host CPU Usage History"
      labels={chartLabels}
      data={cpuData}
      color="#38bdf8"
      unit="%"
    />
    <MetricChart
      title="Host RAM Usage History"
      labels={chartLabels}
      data={ramData}
      color="#34d399"
      unit="%"
    />
  </div>

  <!-- Active Containers Quick Preview Table -->
  <div class="glass-card rounded-3xl border border-slate-800/80 p-6 space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="font-bold text-lg text-white">Active Containers Preview</h3>
        <p class="text-xs text-slate-400">Real-time status of top running containers</p>
      </div>
      <button
        onclick={onNavigateToContainers}
        class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-sky-400 transition-all border border-slate-700/60"
      >
        <span>View All Containers</span>
        <ArrowUpRight class="w-4 h-4" />
      </button>
    </div>

    {#if !systemStats?.containers || systemStats.containers.length === 0}
      <div class="py-8 text-center text-slate-500 italic text-sm">No containers currently running</div>
    {:else}
      <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="border-b border-slate-800/80 text-slate-400 uppercase tracking-wider font-semibold">
              <th class="py-3 px-4">Container</th>
              <th class="py-3 px-4">State</th>
              <th class="py-3 px-4">CPU %</th>
              <th class="py-3 px-4">RAM Usage</th>
              <th class="py-3 px-4">Net I/O</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/40 text-slate-300 font-medium">
            {#each systemStats.containers.slice(0, 5) as c (c.id)}
              <tr class="hover:bg-slate-800/30 transition-colors">
                <td class="py-3.5 px-4 font-mono font-bold text-slate-100">
                  {c.names?.[0]?.replace(/^\//, '') || c.id.substring(0, 12)}
                </td>
                <td class="py-3.5 px-4">
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider
                    {c.state === 'running' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-rose-500/20 text-rose-400'}">
                    {c.state}
                  </span>
                </td>
                <td class="py-3.5 px-4 font-mono text-sky-400">{c.cpu_usage?.toFixed(1) || 0}%</td>
                <td class="py-3.5 px-4 font-mono text-emerald-400">{formatBytes(c.ram_usage)}</td>
                <td class="py-3.5 px-4 font-mono text-slate-400">↓ {formatBytes(c.net_input)} / ↑ {formatBytes(c.net_output)}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>
