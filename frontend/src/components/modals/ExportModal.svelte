<script>
  import FaIcon from '../FaIcon.svelte'
  import { faDownload } from '@fortawesome/free-solid-svg-icons'

  let { tabs = [], activeTabId = '', onclose } = $props()

  let format = $state('json')
  let scope = $state('all')

  function doExport() {
    const url = `/api/export?format=${format}&scope=${scope === 'active' ? activeTabId : 'all'}`
    window.location.href = url
    onclose?.()
  }
</script>

<div class="modal-backdrop" onclick={(e) => { if (e.target === e.currentTarget) onclose?.() }}>
  <div class="modal-box">
    <h2 class="modal-title">Konfiguration exportieren</h2>

    <div class="field">
      <label>Format</label>
      <select bind:value={format}>
        <option value="json">JSON</option>
        <option value="csv">CSV</option>
      </select>
    </div>

    <div class="field">
      <label>Umfang</label>
      <select bind:value={scope}>
        <option value="all">Alle Tabs</option>
        <option value="active">Nur aktiver Tab ({tabs.find(t => t.id === activeTabId)?.name ?? ''})</option>
      </select>
    </div>

    <div class="modal-actions">
      <button class="btn btn-secondary" onclick={onclose}>Abbrechen</button>
      <button class="btn btn-primary" onclick={doExport}><FaIcon icon={faDownload} /> Exportieren</button>
    </div>
  </div>
</div>
