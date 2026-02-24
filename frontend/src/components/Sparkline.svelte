<script>
  let { history = [], historyMinutes = 10 } = $props()

  const WIDTH = 119
  const HEIGHT = 32
  const BAR_GAP = 1

  const SLOTS = 60
  const BAR_WIDTH = (WIDTH - (SLOTS - 1) * BAR_GAP) / SLOTS

  let bars = $derived.by(() => {
    const recent = history.slice(-SLOTS)
    const padded = [...Array(SLOTS - recent.length).fill({ watt: 0 }), ...recent]
    const maxWatt = Math.max(...padded.map(s => s.watt), 1)

    return padded.map((sample, i) => {
      const barHeight = sample.watt === 0 ? 1 : Math.max(2, (sample.watt / maxWatt) * (HEIGHT - 2))
      return {
        x: i * (BAR_WIDTH + BAR_GAP),
        y: HEIGHT - barHeight,
        w: BAR_WIDTH,
        h: barHeight,
      }
    })
  })
</script>

<div class="sparkline-wrapper">
  <svg
    viewBox="0 0 {WIDTH} {HEIGHT}"
    width="100%"
    height="{HEIGHT}"
    preserveAspectRatio="none"
    aria-hidden="true"
  >
    {#each bars as bar}
      {#if bar.h > 0}
        <rect
          x={bar.x}
          y={bar.y}
          width={bar.w}
          height={bar.h}
          fill="var(--sparkline-bar)"
          rx="1"
        />
      {/if}
    {/each}
  </svg>
  <span class="sparkline-label">← {historyMinutes} min</span>
</div>

<style>
  .sparkline-wrapper {
    display: flex;
    flex-direction: column;
  }

  .sparkline-label {
    text-align: right;
    font-size: 0.65rem;
    color: var(--text-secondary);
  }
</style>
