<script>
  import FaIcon from './FaIcon.svelte'
  import { faPen, faTrash, faPlus } from '@fortawesome/free-solid-svg-icons'
  import { editMode } from '../stores/editMode.svelte.js'
  import { deviceStore } from '../stores/devices.svelte.js'

  let { tabs = [], activeTabId, onchange, onedit, ondelete, onadd } = $props()

  let dragOverTabId = $state(null)

  function handleDragOver(e, tabId) {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    dragOverTabId = tabId
  }

  function handleDragLeave(e) {
    // Nur zurücksetzen wenn wir das Tab-Element wirklich verlassen (nicht Kindelemente)
    if (!e.currentTarget.contains(e.relatedTarget)) {
      dragOverTabId = null
    }
  }

  async function handleDrop(e, tabId) {
    e.preventDefault()
    dragOverTabId = null
    const deviceId = e.dataTransfer.getData('text/plain')
    if (deviceId) {
      await deviceStore.patchDevice(deviceId, { tab_id: tabId })
      onchange?.(tabId) // Tab erst nach erfolgreichem Drop wechseln
    }
  }
</script>

<div class="tab-bar">
  {#each tabs as tab (tab.id)}
    <div
      class="tab-item"
      class:active={tab.id === activeTabId}
      class:drag-over-tab={dragOverTabId === tab.id}
      onclick={() => onchange?.(tab.id)}
      role="tab"
      tabindex="0"
      onkeydown={(e) => e.key === 'Enter' && onchange?.(tab.id)}
      ondragover={(e) => handleDragOver(e, tab.id)}
      ondragleave={handleDragLeave}
      ondrop={(e) => handleDrop(e, tab.id)}
    >
      {tab.name}
      {#if editMode.active}
        <span class="tab-actions" onclick={(e) => e.stopPropagation()}>
          <button onclick={() => onedit?.(tab)} title="Tab umbenennen"><FaIcon icon={faPen} /></button>
          <button onclick={() => ondelete?.(tab)} title="Tab löschen"><FaIcon icon={faTrash} /></button>
        </span>
      {/if}
    </div>
  {/each}

  {#if editMode.active}
    <button class="add-tab-btn" onclick={() => onadd?.()} title="Tab hinzufügen">
      <FaIcon icon={faPlus} />
    </button>
  {/if}
</div>
