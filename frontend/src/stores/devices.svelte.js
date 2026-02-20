let tabs = $state([])
let devices = $state([])
let loading = $state(true)
let error = $state(null)

async function loadDevices() {
  try {
    const res = await fetch('/api/devices')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()
    tabs = data.tabs
    devices = data.devices
    loading = false
  } catch (e) {
    error = e.message
    loading = false
  }
}

function connectSSE() {
  const es = new EventSource('/api/events')
  es.addEventListener('device_update', (e) => {
    const update = JSON.parse(e.data)
    const idx = devices.findIndex(d => d.id === update.id)
    if (idx !== -1) {
      devices[idx] = { ...devices[idx], on: update.on, watt: update.watt, online: update.online }
    }
  })
  es.onerror = () => {
    // Reconnect automatisch durch EventSource
  }
  return es
}

export const deviceStore = {
  get tabs() { return tabs },
  get devices() { return devices },
  get loading() { return loading },
  get error() { return error },
  load: loadDevices,
  connectSSE,

  async toggle(id, on) {
    await fetch(`/api/devices/${id}/toggle`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ on }),
    })
  },

  async addDevice(device) {
    const res = await fetch('/api/devices', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(device),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }
    await loadDevices()
  },

  async updateDevice(id, device) {
    const res = await fetch(`/api/devices/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(device),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }
    await loadDevices()
  },

  async patchDevice(id, fields) {
    await fetch(`/api/devices/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(fields),
    })
    await loadDevices()
  },

  async deleteDevice(id) {
    const res = await fetch(`/api/devices/${id}`, { method: 'DELETE' })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }
    await loadDevices()
  },

  async testDevice(id) {
    const res = await fetch(`/api/devices/${id}/test`, { method: 'POST' })
    return await res.json()
  },

  async addTab(tab) {
    const res = await fetch('/api/tabs', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(tab),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }
    await loadDevices()
  },

  async updateTab(id, tab) {
    const res = await fetch(`/api/tabs/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(tab),
    })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }
    await loadDevices()
  },

  async deleteTab(id) {
    const res = await fetch(`/api/tabs/${id}`, { method: 'DELETE' })
    if (!res.ok) {
      const err = await res.json()
      throw new Error(err.error)
    }
    await loadDevices()
  },
}
