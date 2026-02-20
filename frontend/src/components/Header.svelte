<script>
  import FaIcon from './FaIcon.svelte'
  import { faBolt, faPen, faCheck, faUpload, faDownload, faPlus, faMoon, faSun } from '@fortawesome/free-solid-svg-icons'
  import { theme } from '../stores/theme.svelte.js'
  import { editMode } from '../stores/editMode.svelte.js'

  let { onimport, onexport, onadd } = $props()
</script>

<header class="app-header">
  <div class="logo">
    <FaIcon icon={faBolt} />
    SHELLY CONTROL
  </div>

  <label class="theme-toggle" title={theme.isDark ? 'Light Mode' : 'Dark Mode'}>
    <input type="checkbox" checked={theme.isDark} onchange={() => theme.toggle()} />
    <span class="track">
      <span class="track-icon track-icon--left"><FaIcon icon={faSun} /></span>
      <span class="thumb"></span>
      <span class="track-icon track-icon--right"><FaIcon icon={faMoon} /></span>
    </span>
  </label>

  <button
    class="btn btn-ghost"
    class:btn-edit-active={editMode.active}
    onclick={() => editMode.toggle()}
  >
    <FaIcon icon={editMode.active ? faCheck : faPen} />
    {editMode.active ? 'Fertig' : 'Bearbeiten'}
  </button>

  <button class="btn btn-ghost" onclick={() => onimport?.()}>
    <FaIcon icon={faUpload} /> Import
  </button>

  <button class="btn btn-ghost" onclick={() => onexport?.()}>
    <FaIcon icon={faDownload} /> Export
  </button>

  <button class="btn btn-primary" onclick={() => onadd?.()}>
    <FaIcon icon={faPlus} /> Gerät
  </button>
</header>
