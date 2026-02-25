<script>
  import DeviceCard from './DeviceCard.svelte'

  let { devices = [], historyMinutes = 10, activeTabId, onedit, ondelete } = $props()

  let filtered = $derived(
    devices
      .filter(d => d.tab_id === activeTabId)
      .toSorted((a, b) => (a.order ?? 99) - (b.order ?? 99) || a.name.localeCompare(b.name))
  )
</script>

<div class="device-grid-container">
  {#if filtered.length === 0}
    <div class="empty-state">
      <p>Keine Geräte in diesem Tab.</p>
    </div>
  {:else}
    <div class="device-grid">
      {#each filtered as device (device.id)}
        <DeviceCard
          {device}
          {historyMinutes}
          onedit={onedit}
          ondelete={ondelete}
        />
      {/each}
    </div>
  {/if}
</div>

<style>
  .device-grid-container {
    padding: 1rem;
    flex: 1;
  }

  .device-grid {
    container-type: inline-size;
    container-name: cards-grid;
    display: grid;
    grid-template-columns: repeat(6, 1fr);
    gap: 0.75rem;
  }

  @container cards-grid (max-width: 1400px) {
    .device-grid { grid-template-columns: repeat(5, 1fr); }
  }
  @container cards-grid (max-width: 1100px) {
    .device-grid { grid-template-columns: repeat(4, 1fr); }
  }
  @container cards-grid (max-width: 800px) {
    .device-grid { grid-template-columns: repeat(3, 1fr); }
  }
  @container cards-grid (max-width: 550px) {
    .device-grid { grid-template-columns: repeat(2, 1fr); }
  }
</style>
