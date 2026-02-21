<script>
  import { onMount } from 'svelte'
  import Header from './components/Header.svelte'
  import TabBar from './components/TabBar.svelte'
  import DeviceGrid from './components/DeviceGrid.svelte'
  import DeviceModal from './components/modals/DeviceModal.svelte'
  import TabModal from './components/modals/TabModal.svelte'
  import ImportModal from './components/modals/ImportModal.svelte'
  import ExportModal from './components/modals/ExportModal.svelte'
  import ConfirmModal from './components/modals/ConfirmModal.svelte'
  import { deviceStore } from './stores/devices.svelte.js'

  let activeTabId = $state('')
  let historyMinutes = $state(5)

  // Modal state
  let modal = $state(null) // { type: 'device'|'tab'|'import'|'export'|'confirm', data: ... }

  onMount(async () => {
    await deviceStore.load()
    if (deviceStore.tabs.length > 0 && !activeTabId) {
      activeTabId = deviceStore.tabs[0].id
    }
    const es = deviceStore.connectSSE()

    // config.json power_history_minutes über API nicht direkt verfügbar,
    // aber wir können den Wert aus der History-Länge ableiten.
    // Für jetzt: Standardwert 10 Minuten.
    return () => es.close()
  })

  // Sicherstellen, dass activeTabId immer gültig ist
  $effect(() => {
    if (deviceStore.tabs.length > 0) {
      const valid = deviceStore.tabs.some(t => t.id === activeTabId)
      if (!valid) activeTabId = deviceStore.tabs[0].id
    }
  })

  function openAddDevice() {
    modal = { type: 'device', data: null }
  }

  function openEditDevice(device) {
    modal = { type: 'device', data: device }
  }

  function openDeleteDevice(device) {
    modal = {
      type: 'confirm',
      title: 'Gerät löschen',
      message: `Gerät „${device.name}" wirklich löschen?`,
      onconfirm: async () => {
        await deviceStore.deleteDevice(device.id)
        modal = null
      }
    }
  }

  function openAddTab() {
    modal = { type: 'tab', data: null }
  }

  function openEditTab(tab) {
    modal = { type: 'tab', data: tab }
  }

  function openDeleteTab(tab) {
    modal = {
      type: 'confirm',
      title: 'Tab löschen',
      message: `Tab „${tab.name}" wirklich löschen? (Nur möglich wenn leer)`,
      onconfirm: async () => {
        try {
          await deviceStore.deleteTab(tab.id)
        } catch (e) {
          alert(e.message)
        }
        modal = null
      }
    }
  }
</script>

<Header
  onimport={() => { modal = { type: 'import' } }}
  onexport={() => { modal = { type: 'export' } }}
  onadd={openAddDevice}
/>

<TabBar
  tabs={deviceStore.tabs}
  {activeTabId}
  onchange={(id) => { activeTabId = id }}
  onedit={openEditTab}
  ondelete={openDeleteTab}
  onadd={openAddTab}
/>

{#if deviceStore.loading}
  <div class="empty-state">Lade Geräte…</div>
{:else if deviceStore.error}
  <div class="device-grid-container">
    <div class="error-banner">Fehler: {deviceStore.error}</div>
  </div>
{:else}
  <DeviceGrid
    devices={deviceStore.devices}
    {activeTabId}
    {historyMinutes}
    onedit={openEditDevice}
    ondelete={openDeleteDevice}
  />
{/if}

<footer class="app-footer">
  <span><strong>iTec</strong> - learn together grow together</span>
</footer>

<!-- Modals -->

{#if modal?.type === 'device'}
  <DeviceModal
    device={modal.data}
    tabs={deviceStore.tabs}
    {activeTabId}
    onclose={() => { modal = null }}
  />
{/if}

{#if modal?.type === 'tab'}
  <TabModal
    tab={modal.data}
    onclose={() => { modal = null }}
  />
{/if}

{#if modal?.type === 'import'}
  <ImportModal
    onclose={() => { modal = null; deviceStore.load() }}
  />
{/if}

{#if modal?.type === 'export'}
  <ExportModal
    tabs={deviceStore.tabs}
    {activeTabId}
    onclose={() => { modal = null }}
  />
{/if}

{#if modal?.type === 'confirm'}
  <ConfirmModal
    title={modal.title}
    message={modal.message}
    onconfirm={modal.onconfirm}
    oncancel={() => { modal = null }}
  />
{/if}

<style>
  .app-footer {
    text-align: right;
    padding: 0.5rem 1rem;
    font-size: 0.75rem;
    color: var(--text-secondary);
    border-top: 1px solid var(--card-border);
  }
</style>
