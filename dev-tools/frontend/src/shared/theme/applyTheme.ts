import type { components } from '../api/schema'

export type Theme = components['schemas']['Preferences']['theme']

export function applyTheme(theme: Theme): void {
  document.documentElement.dataset.theme = theme
}
