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

  <Sparkline history={history.length > 0 ? history : (device.online ? [{ timestamp: new Date().toISOString(), watt: device.watt }] : [])} {historyMinutes} />
</div>

<style>
  .device-card {
    background: var(--card-bg);
    border: 1px solid var(--card-border);
    border-radius: 6px;
    padding: 0.75rem;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    position: relative;
    transition: border-color 0.15s;
  }

  .device-card:hover  { border-color: #3273dc44; }
  .device-card.offline { opacity: 0.6; }
  .device-card.dragging { opacity: 0.25; }

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .card-name {
    font-weight: 600;
    font-size: 0.9rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .card-description {
    font-size: 0.78rem;
    color: var(--text-secondary);
    min-height: 1.1em;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .card-watt {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-primary);
  }

  .card-watt .unit {
    font-size: 0.75rem;
    font-weight: 400;
    color: var(--text-secondary);
    margin-left: 0.2rem;
  }

  /* ── Relay Toggle ───────────────────────────────────────── */
  .relay-toggle {
    display: flex;
    align-items: center;
    cursor: pointer;
    gap: 0.3rem;
    flex-shrink: 0;
  }

  .relay-toggle input[type="checkbox"] { display: none; }

  .track {
    width: 40px;
    height: 22px;
    background: var(--toggle-off);
    border-radius: 11px;
    position: relative;
    transition: background 0.2s;
  }

  .track::after {
    content: '';
    position: absolute;
    top: 3px;
    left: 3px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: white;
    transition: transform 0.2s;
    box-shadow: 0 1px 3px rgba(0,0,0,0.3);
  }

  .relay-toggle input:checked + .track { background: var(--toggle-on); }
  .relay-toggle input:checked + .track::after { transform: translateX(18px); }

  /* ── Edit Mode ──────────────────────────────────────────── */
  .card-action-btn {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.2rem 0.3rem;
    color: var(--text-secondary);
    font-size: 0.8rem;
    line-height: 1;
    border-radius: 3px;
    flex-shrink: 0;
    transition: color 0.15s;
  }

  .card-action-btn:hover { color: var(--text-primary); }
  .card-action-btn--danger:hover { color: #f14668; }

  .drag-handle {
    cursor: grab;
    color: var(--text-secondary);
    padding: 0 0.3rem 0 0;
    font-size: 0.9rem;
    align-self: center;
    flex-shrink: 0;
  }

  .drag-handle:active { cursor: grabbing; }
</style>
