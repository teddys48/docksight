<script>
  import { onMount, onDestroy } from 'svelte';
  import { theme } from '../stores.js';
  import {
    Chart,
    LineController,
    LineElement,
    PointElement,
    LinearScale,
    CategoryScale,
    Title,
    Tooltip,
    Legend,
    Filler
  } from 'chart.js';

  Chart.register(
    LineController,
    LineElement,
    PointElement,
    LinearScale,
    CategoryScale,
    Title,
    Tooltip,
    Legend,
    Filler
  );

  export let title = 'Metric Chart';
  export let labels = [];
  export let data = [];
  export let color = '#38bdf8'; // Sky blue
  export let unit = '%';

  let canvasEl;
  let chartInstance;
  let unsubscribeTheme;

  $: if (chartInstance && data) {
    chartInstance.data.labels = labels;
    chartInstance.data.datasets[0].data = data;
    chartInstance.update('none');
  }

  function applyTheme(isLight) {
    if (!chartInstance) return;
    const gridColor = isLight ? 'rgba(203, 213, 225, 0.6)' : 'rgba(51, 65, 85, 0.3)';
    const tickColor = isLight ? '#475569' : '#64748b';

    chartInstance.options.scales.y.grid.color = gridColor;
    chartInstance.options.scales.x.ticks.color = tickColor;
    chartInstance.options.scales.y.ticks.color = tickColor;

    chartInstance.options.plugins.tooltip.backgroundColor = isLight ? 'rgba(255, 255, 255, 0.95)' : 'rgba(15, 23, 42, 0.9)';
    chartInstance.options.plugins.tooltip.titleColor = isLight ? '#475569' : '#94a3b8';
    chartInstance.options.plugins.tooltip.bodyColor = isLight ? '#0f172a' : '#f8fafc';
    chartInstance.options.plugins.tooltip.borderColor = isLight ? '#cbd5e1' : '#334155';

    chartInstance.update('none');
  }

  onMount(() => {
    const ctx = canvasEl.getContext('2d');

    // Create Gradient
    const gradient = ctx.createLinearGradient(0, 0, 0, 200);
    gradient.addColorStop(0, color + '55'); // Semi transparent
    gradient.addColorStop(1, color + '00'); // Transparent

    const isLight = $theme === 'light';

    chartInstance = new Chart(ctx, {
      type: 'line',
      data: {
        labels: labels,
        datasets: [
          {
            label: title,
            data: data,
            borderColor: color,
            borderWidth: 2.5,
            backgroundColor: gradient,
            fill: true,
            tension: 0.4,
            pointRadius: 0,
            pointHoverRadius: 5,
            pointHoverBackgroundColor: color,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: false,
        plugins: {
          legend: { display: false },
          tooltip: {
            mode: 'index',
            intersect: false,
            backgroundColor: isLight ? 'rgba(255, 255, 255, 0.95)' : 'rgba(15, 23, 42, 0.9)',
            titleColor: isLight ? '#475569' : '#94a3b8',
            bodyColor: isLight ? '#0f172a' : '#f8fafc',
            borderColor: isLight ? '#cbd5e1' : '#334155',
            borderWidth: 1,
            padding: 10,
            displayColors: false,
            callbacks: {
              label: (context) => `${title}: ${context.parsed.y.toFixed(1)}${unit}`,
            },
          },
        },
        scales: {
          x: {
            grid: { display: false },
            ticks: {
              color: isLight ? '#475569' : '#64748b',
              maxTicksLimit: 6,
              font: { size: 10 },
            },
          },
          y: {
            grid: { color: isLight ? 'rgba(203, 213, 225, 0.6)' : 'rgba(51, 65, 85, 0.3)' },
            ticks: {
              color: isLight ? '#475569' : '#64748b',
              font: { size: 10 },
              callback: (val) => `${val}${unit}`,
            },
            min: 0,
          },
        },
      },
    });

    unsubscribeTheme = theme.subscribe(val => {
      applyTheme(val === 'light');
    });
  });

  onDestroy(() => {
    if (unsubscribeTheme) unsubscribeTheme();
    if (chartInstance) {
      chartInstance.destroy();
    }
  });
</script>

<div class="glass-card p-5 rounded-2xl border border-slate-800 flex flex-col justify-between h-72">
  <div class="flex items-center justify-between mb-4">
    <h3 class="font-semibold text-sm text-slate-300 flex items-center gap-2">
      <span class="w-2.5 h-2.5 rounded-full" style="background-color: {color};"></span>
      {title}
    </h3>
    <span class="text-xs font-mono font-bold text-slate-400">
      {data.length > 0 ? `${data[data.length - 1]?.toFixed(1)}${unit}` : `0${unit}`}
    </span>
  </div>
  <div class="relative w-full flex-1">
    <canvas bind:this={canvasEl}></canvas>
  </div>
</div>
