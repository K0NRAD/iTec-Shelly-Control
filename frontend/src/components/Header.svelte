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

<style>
  .app-header {
    height: 56px;
    background: var(--header-bg);
    border-bottom: 1px solid var(--header-border);
    display: flex;
    align-items: center;
    padding: 0 1rem;
    gap: 0.75rem;
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .logo {
    font-weight: 700;
    font-size: 1.1rem;
    letter-spacing: 0.05em;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-right: auto;
  }

  .logo :global(svg) { color: #df376a; }

  /* ── Dark Mode Toggle ──────────────────────────────────── */
  .theme-toggle {
    display: flex;
    align-items: center;
    cursor: pointer;
    user-select: none;
  }

  .theme-toggle input[type="checkbox"] { display: none; }

  .track {
    width: 56px;
    height: 26px;
    background: var(--toggle-off);
    border-radius: 13px;
    position: relative;
    display: flex;
    align-items: center;
    transition: background 0.25s;
  }

  .theme-toggle input:checked + .track { background: #2c3e6b; }

  .thumb {
    position: absolute;
    top: 3px;
    left: 3px;
    width: 20px;
    height: 20px;
    border-radius: 50%;
    background: white;
    box-shadow: 0 1px 3px rgba(0,0,0,0.3);
    transition: transform 0.25s;
    z-index: 1;
  }

  .theme-toggle input:checked + .track .thumb { transform: translateX(30px); }

  .track-icon {
    position: absolute;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    font-size: 11px;
    transition: opacity 0.25s;
    pointer-events: none;
  }

  .track-icon--left  { left: 3px; }
  .track-icon--right { right: 3px; }

  /* Light Mode aktiv (unchecked): Sonne blau, Mond gedimmt aber sichtbar */
  .theme-toggle input:not(:checked) + .track .track-icon--left  { color: #3273dc; opacity: 1;   }
  .theme-toggle input:not(:checked) + .track .track-icon--right { color: #444;    opacity: 0.5; }

  /* Dark Mode aktiv (checked): Mond weiß, Sonne gedimmt aber sichtbar */
  .theme-toggle input:checked + .track .track-icon--left  { color: #fff;     opacity: 0.5; }
  .theme-toggle input:checked + .track .track-icon--right { color: #ffffff;  opacity: 1;   }
</style>
