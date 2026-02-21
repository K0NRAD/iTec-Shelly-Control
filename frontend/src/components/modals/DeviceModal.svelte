<script>
  import { deviceStore } from '../../stores/devices.svelte.js'

  let { device = null, tabs = [], activeTabId = '', onclose } = $props()

  const isEdit = device !== null

  let form = $state({
    id:          device?.id ?? '',
    name:        device?.name ?? '',
    ip:          device?.ip ?? '',
    generation:  device?.generation ?? 1,
    tab_id:      device?.tab_id ?? activeTabId ?? (tabs[0]?.id ?? ''),
    description: device?.description ?? '',
  })

  let error = $state('')
  let saving = $state(false)
  let testing = $state(false)
  let testResult = $state(null)

  async function save() {
    if (!form.name.trim() || !form.ip.trim()) {
      error = 'Name und IP sind Pflichtfelder'
      return
    }
    if (!isEdit && !form.id.trim()) {
      error = 'ID ist ein Pflichtfeld'
      return
    }
    saving = true
    error = ''
    try {
      if (isEdit) {
        await deviceStore.updateDevice(device.id, { ...form, generation: Number(form.generation) })
      } else {
        await deviceStore.addDevice({ ...form, generation: Number(form.generation) })
      }
      onclose?.()
    } catch (e) {
      error = e.message
    } finally {
      saving = false
    }
  }

  async function testConnection() {
    if (!form.ip) return
    testing = true
    testResult = null
    try {
      // Direkter Test ohne gespeichertes Gerät → temporäres Gerät anlegen wäre aufwändig.
      // Stattdessen: bei existierendem Gerät /api/devices/{id}/test, sonst IP-Ping via Backend.
      // Für neue Geräte zeigen wir einen Hinweis.
      if (isEdit) {
        const res = await deviceStore.testDevice(device.id)
        testResult = res.ok ? 'Verbindung OK' : `Fehler: ${res.error}`
      } else {
        testResult = 'Test nur nach dem Anlegen möglich'
      }
    } finally {
      testing = false
    }
  }
</script>

<div class="modal-backdrop" onclick={(e) => { if (e.target === e.currentTarget) onclose?.() }}>
  <div class="modal-box">
    <h2 class="modal-title">{isEdit ? 'Gerät bearbeiten' : 'Gerät anlegen'}</h2>

    {#if !isEdit}
      <div class="field">
        <label>ID (unveränderlich)</label>
        <input type="text" bind:value={form.id} placeholder="z.B. s1r1" />
      </div>
    {/if}

    <div class="field">
      <label>Name</label>
      <input type="text" bind:value={form.name} placeholder="z.B. S1R1" autofocus={isEdit} />
    </div>

    <div class="field">
      <label>IP-Adresse</label>
      <input type="text" bind:value={form.ip} placeholder="192.168.1.50" />
    </div>

    <div class="field">
      <label>Generation</label>
      <select bind:value={form.generation}>
        <option value={1}>Gen 1 — Shelly Plug S</option>
        <option value={3}>Gen 3 — Shelly Plug S MTR</option>
        <option value={4}>Gen 4 — Shelly 1PM Gen4</option>
      </select>
    </div>

    <div class="field">
      <label>Tab</label>
      <select bind:value={form.tab_id}>
        {#each tabs as tab}
          <option value={tab.id}>{tab.name}</option>
        {/each}
      </select>
    </div>

    <div class="field">
      <label>Beschreibung</label>
      <input type="text" bind:value={form.description} placeholder="z.B. Dell PowerEdge R640" />
    </div>

    {#if error}
      <div class="error-banner">{error}</div>
    {/if}

    {#if testResult}
      <div class="error-banner" style={testResult.startsWith('Verbindung') ? 'border-color: #48c774; color: #48c774; background: #48c77410;' : ''}>
        {testResult}
      </div>
    {/if}

    <div class="modal-actions">
      <button class="btn btn-ghost" onclick={testConnection} disabled={testing} style="margin-right: auto;">
        {testing ? 'Teste…' : 'Verbindung testen'}
      </button>
      <button class="btn btn-secondary" onclick={onclose}>Abbrechen</button>
      <button class="btn btn-primary" onclick={save} disabled={saving}>
        {saving ? 'Speichern…' : 'Speichern'}
      </button>
    </div>
  </div>
</div>
