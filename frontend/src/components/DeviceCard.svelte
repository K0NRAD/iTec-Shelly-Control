<script>
  import Sparkline from './Sparkline.svelte'
  import FaIcon from './FaIcon.svelte'
  import { faGripVertical, faPen, faTrash } from '@fortawesome/free-solid-svg-icons'
  import { deviceStore } from '../stores/devices.svelte.js'
  import { editMode } from '../stores/editMode.svelte.js'

  let { device, historyMinutes = 10, onedit, ondelete } = $props()

  let history = $state([])
  let toggling = $state(false)
  let dragging = $state(false)
  let dragEnabled = $state(false)

  async function loadHistory() {
    try {
      const res = await fetch(`/api/devices/${device.id}/history`)
      if (res.ok) history = await res.json()
    } catch {}
  }

  $effect(() => {
    loadHistory()
    const interval = setInterval(loadHistory, 5000)
    return () => clearInterval(interval)
  })

  async function handleToggle(e) {
    if (toggling) return
    toggling = true
    try {
      await deviceStore.toggle(device.id, e.target.checked)
    } finally {
      toggling = false
    }
  }

  function handleDragStart(e) {
    if (!dragEnabled) { e.preventDefault(); return }
    e.dataTransfer.setData('text/plain', device.id)
    e.dataTransfer.effectAllowed = 'move'
    dragging = true
  }

  function handleDragEnd() {
    dragging = false
    dragEnabled = false
  }

</script>

<div
  class="device-card"
  class:offline={!device.online}
  class:dragging={dragging}
  draggable={editMode.active}
  ondragstart={handleDragStart}
  ondragend={handleDragEnd}
>
  <div class="card-header">
    {#if editMode.active}
      <span
        class="drag-handle"
        title="Verschieben"
        onpointerdown={() => { dragEnabled = true }}
        onpointerup={() => { dragEnabled = false }}
      ><FaIcon icon={faGripVertical} /></span>
    {/if}
    <span class="card-name" title={device.name}>{device.name}</span>
    <label class="relay-toggle" title={device.on ? 'Ausschalten' : 'Einschalten'}>
      <input
        type="checkbox"
        checked={device.on}
        disabled={toggling || !device.online}
        onchange={handleToggle}
      />
      <span class="track"></span>
    </label>
    {#if editMode.active}
      <button class="card-action-btn" onclick={() => onedit?.(device)} title="Bearbeiten"><FaIcon icon={faPen} /></button>
      <button class="card-action-btn card-action-btn--danger" onclick={() => ondelete?.(device)} title="Löschen"><FaIcon icon={faTrash} /></button>
    {/if}
  </div>

  <div class="card-description" title={device.description}>
    {device.description || '\u00a0'}
  </div>

  <div class="card-watt">
    {device.online ? device.watt.toFixed(1) : '—'}
    <span class="unit">W</span>
  </div>

  <Sparkline {history} {historyMinutes} />
</div>
