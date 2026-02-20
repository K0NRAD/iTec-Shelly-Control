<script>
  import { deviceStore } from '../../stores/devices.svelte.js'

  let { tab = null, onclose } = $props()

  const isEdit = tab !== null
  let name = $state(tab?.name ?? '')
  let id = $state(tab?.id ?? '')
  let error = $state('')
  let saving = $state(false)

  async function save() {
    if (!name.trim()) { error = 'Name ist Pflicht'; return }
    saving = true
    error = ''
    try {
      if (isEdit) {
        await deviceStore.updateTab(tab.id, { id: tab.id, name: name.trim(), order: tab.order })
      } else {
        if (!id.trim()) { error = 'ID ist Pflicht'; saving = false; return }
        await deviceStore.addTab({ id: id.trim(), name: name.trim(), order: 99 })
      }
      onclose?.()
    } catch (e) {
      error = e.message
    } finally {
      saving = false
    }
  }
</script>

<div class="modal-backdrop" onclick={(e) => { if (e.target === e.currentTarget) onclose?.() }}>
  <div class="modal-box">
    <h2 class="modal-title">{isEdit ? 'Tab umbenennen' : 'Tab anlegen'}</h2>

    {#if !isEdit}
      <div class="field">
        <label>ID (unveränderlich)</label>
        <input type="text" bind:value={id} placeholder="z.B. rack3" />
      </div>
    {/if}

    <div class="field">
      <label>Name</label>
      <input type="text" bind:value={name} placeholder="z.B. Rack 3" autofocus />
    </div>

    {#if error}
      <div class="error-banner">{error}</div>
    {/if}

    <div class="modal-actions">
      <button class="btn btn-secondary" onclick={onclose}>Abbrechen</button>
      <button class="btn btn-primary" onclick={save} disabled={saving}>
        {saving ? 'Speichern…' : 'Speichern'}
      </button>
    </div>
  </div>
</div>
