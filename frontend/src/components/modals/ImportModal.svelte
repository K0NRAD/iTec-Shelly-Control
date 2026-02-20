<script>
  import FaIcon from '../FaIcon.svelte'
  import { faUpload } from '@fortawesome/free-solid-svg-icons'

  let { onclose } = $props()

  let format = $state('json')
  let mode = $state('add')
  let preview = $state(null)
  let loading = $state(false)
  let error = $state('')
  let dragActive = $state(false)
  let fileContent = $state('')
  let fileName = $state('')

  function handleDragOver(e) {
    e.preventDefault()
    dragActive = true
  }

  function handleDragLeave() {
    dragActive = false
  }

  function handleDrop(e) {
    e.preventDefault()
    dragActive = false
    const file = e.dataTransfer.files[0]
    if (file) loadFile(file)
  }

  function handleFileInput(e) {
    const file = e.target.files[0]
    if (file) loadFile(file)
  }

  function loadFile(file) {
    fileName = file.name
    if (file.name.endsWith('.csv')) format = 'csv'
    else format = 'json'
    const reader = new FileReader()
    reader.onload = (e) => { fileContent = e.target.result }
    reader.readAsText(file)
  }

  async function loadPreview() {
    if (!fileContent) { error = 'Bitte Datei auswählen'; return }
    loading = true
    error = ''
    try {
      const res = await fetch(`/api/import?format=${format}&mode=${mode}&preview=true`, {
        method: 'POST',
        headers: { 'Content-Type': format === 'csv' ? 'text/csv' : 'application/json' },
        body: fileContent,
      })
      if (!res.ok) {
        const e = await res.json()
        throw new Error(e.error)
      }
      preview = await res.json()
    } catch (e) {
      error = e.message
    } finally {
      loading = false
    }
  }

  async function doImport() {
    if (!fileContent) { error = 'Bitte Datei auswählen'; return }
    loading = true
    error = ''
    try {
      const res = await fetch(`/api/import?format=${format}&mode=${mode}`, {
        method: 'POST',
        headers: { 'Content-Type': format === 'csv' ? 'text/csv' : 'application/json' },
        body: fileContent,
      })
      if (!res.ok) {
        const e = await res.json()
        throw new Error(e.error)
      }
      onclose?.()
    } catch (e) {
      error = e.message
    } finally {
      loading = false
    }
  }

  const statusLabel = { new: 'Neu', update: 'Update', unchanged: 'Unverändert', error: 'Fehler' }
  const statusClass = { new: 'status-new', update: 'status-update', unchanged: 'status-unchanged', error: 'status-error' }
</script>

<div class="modal-backdrop" onclick={(e) => { if (e.target === e.currentTarget) onclose?.() }}>
  <div class="modal-box modal-wide">
    <h2 class="modal-title">Konfiguration importieren</h2>

    <div
      class="dropzone"
      class:drag-active={dragActive}
      ondragover={handleDragOver}
      ondragleave={handleDragLeave}
      ondrop={handleDrop}
      onclick={() => document.getElementById('import-file-input').click()}
      role="button"
      tabindex="0"
    >
      {#if fileName}
        📄 {fileName}
      {:else}
        Datei hierher ziehen oder klicken zum Auswählen
      {/if}
    </div>
    <input id="import-file-input" type="file" accept=".json,.csv" style="display:none" onchange={handleFileInput} />

    <div style="display: flex; gap: 1rem;">
      <div class="field" style="flex: 1;">
        <label>Format</label>
        <select bind:value={format}>
          <option value="json">JSON</option>
          <option value="csv">CSV</option>
        </select>
      </div>
      <div class="field" style="flex: 1;">
        <label>Modus</label>
        <select bind:value={mode}>
          <option value="add">Nur neue hinzufügen</option>
          <option value="overwrite">Bestehende überschreiben</option>
          <option value="replace">Alles ersetzen</option>
        </select>
      </div>
    </div>

    {#if error}
      <div class="error-banner">{error}</div>
    {/if}

    {#if preview}
      <div style="overflow-x: auto; max-height: 300px; overflow-y: auto;">
        <table class="import-table">
          <thead>
            <tr>
              <th>Status</th>
              <th>Name</th>
              <th>IP</th>
              <th>Gen</th>
              <th>Tab</th>
              <th>Beschreibung</th>
            </tr>
          </thead>
          <tbody>
            {#each preview as row}
              <tr>
                <td class={statusClass[row.status]}>{statusLabel[row.status]}{row.error ? `: ${row.error}` : ''}</td>
                <td>{row.name}</td>
                <td>{row.ip}</td>
                <td>{row.generation}</td>
                <td>{row.tab}</td>
                <td>{row.description}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}

    <div class="modal-actions">
      <button class="btn btn-ghost" onclick={loadPreview} disabled={loading || !fileContent} style="margin-right: auto;">
        Vorschau laden
      </button>
      <button class="btn btn-secondary" onclick={onclose}>Abbrechen</button>
      <button class="btn btn-primary" onclick={doImport} disabled={loading || !fileContent}>
        {#if !loading}<FaIcon icon={faUpload} />{/if}
        {loading ? 'Importiere…' : 'Importieren'}
      </button>
    </div>
  </div>
</div>
