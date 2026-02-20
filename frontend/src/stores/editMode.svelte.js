let active = $state(false)

export const editMode = {
  get active() { return active },
  toggle() { active = !active },
  deactivate() { active = false },
}
