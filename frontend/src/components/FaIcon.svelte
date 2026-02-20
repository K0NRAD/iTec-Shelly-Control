<script>
  import { icon as faIcon } from '@fortawesome/fontawesome-svg-core'

  let { icon, class: cls = '', style = '' } = $props()

  let rendered = $derived.by(() => {
    const result = faIcon(icon)
    if (!result) return null
    const [width, height, , , path] = result.icon
    const pathD = Array.isArray(path) ? path.join(' ') : path
    return { width, height, pathD }
  })
</script>

{#if rendered}
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 {rendered.width} {rendered.height}"
    width="1em"
    height="1em"
    fill="currentColor"
    aria-hidden="true"
    class={cls}
    {style}
  >
    <path d={rendered.pathD} />
  </svg>
{/if}
