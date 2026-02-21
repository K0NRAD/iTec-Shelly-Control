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
      onchange?.(tabId)
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

<style>
  .tab-bar {
    background: var(--header-bg);
    border-bottom: 1px solid var(--header-border);
    padding: 0 1rem;
    display: flex;
    align-items: flex-end;
    gap: 0;
    overflow-x: auto;
  }

  .tab-item {
    padding: 0.55rem 1.1rem;
    border: 1px solid transparent;
    border-radius: 6px 6px 0 0;
    cursor: pointer;
    white-space: nowrap;
    color: var(--text-secondary);
    background: var(--tab-inactive-bg);
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.9rem;
    margin-right: 3px;
    margin-bottom: -1px;
    transition: background 0.15s, color 0.15s;
  }

  .tab-item:hover {
    color: var(--text-primary);
    background: var(--tab-hover-bg);
  }

  .tab-item.active {
    color: var(--tab-active-color);
    font-weight: 600;
    background: var(--tab-active-bg);
    border-color: var(--tab-active-border);
    border-bottom-color: var(--tab-active-bg);
  }

  .tab-item.drag-over-tab {
    color: #ffffff;
    background: #3273dc;
    border-color: #2366d1;
    border-bottom-color: #3273dc;
    font-weight: 600;
  }

  .tab-actions {
    display: flex;
    gap: 0.2rem;
    margin-left: 0.3rem;
  }

  .tab-actions button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.1rem 0.2rem;
    color: inherit;
    font-size: 0.75rem;
    line-height: 1;
    opacity: 0.6;
  }

  .tab-actions button:hover { opacity: 1; }

  .add-tab-btn {
    padding: 0.6rem 0.75rem;
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 1.1rem;
    background: none;
    border: none;
    line-height: 1;
  }

  .add-tab-btn:hover { color: #3273dc; }
</style>
