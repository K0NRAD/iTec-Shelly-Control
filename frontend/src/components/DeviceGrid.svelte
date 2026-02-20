<script>
  import DeviceCard from './DeviceCard.svelte'

  let { devices = [], historyMinutes = 10, activeTabId, onedit, ondelete } = $props()

  let filtered = $derived(devices.filter(d => d.tab_id === activeTabId))
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
