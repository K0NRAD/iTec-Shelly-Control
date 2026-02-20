function getSystemPreference() {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

let current = $state(getSystemPreference())

$effect.root(() => {
  $effect(() => {
    document.documentElement.setAttribute('data-theme', current)
  })
})

export const theme = {
  get isDark() { return current === 'dark' },
  toggle() {
    current = current === 'dark' ? 'light' : 'dark'
  },
}
