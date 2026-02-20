<script>
  let { title = 'Bestätigen', message = '', onconfirm, oncancel } = $props()
</script>

<div class="modal-backdrop" onclick={(e) => { if (e.target === e.currentTarget) oncancel?.() }}>
  <div class="modal-box">
    <h2 class="modal-title">{title}</h2>
    <p style="color: var(--text-secondary); margin: 0;">{message}</p>
    <div class="modal-actions">
      <button class="btn btn-secondary" onclick={oncancel}>Abbrechen</button>
      <button class="btn btn-danger" onclick={onconfirm}>Löschen</button>
    </div>
  </div>
</div>
