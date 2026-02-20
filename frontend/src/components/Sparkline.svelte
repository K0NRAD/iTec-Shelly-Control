<script>
  let { history = [], historyMinutes = 10 } = $props()

  const WIDTH = 200
  const HEIGHT = 32
  const BAR_GAP = 1

  let bars = $derived.by(() => {
    if (history.length === 0) return []

    const maxWatt = Math.max(...history.map(s => s.watt), 1)
    const count = Math.max(history.length, 1)
    const barWidth = Math.max(1, (WIDTH - (count - 1) * BAR_GAP) / count)

    return history.map((sample, i) => {
      const barHeight = Math.max(2, (sample.watt / maxWatt) * (HEIGHT - 2))
      return {
        x: i * (barWidth + BAR_GAP),
        y: HEIGHT - barHeight,
        w: barWidth,
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
      <rect
        x={bar.x}
        y={bar.y}
        width={bar.w}
        height={bar.h}
        fill="var(--sparkline-bar)"
        rx="1"
      />
    {/each}
  </svg>
  <span class="sparkline-label">← {historyMinutes} min</span>
</div>
