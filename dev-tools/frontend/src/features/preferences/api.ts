import type { components } from '../../shared/api/schema'
import { apiFetch } from '../../shared/api/client'

export type Preferences = components['schemas']['Preferences']

export function savePreferences(preferences: Preferences): Promise<Preferences> {
  return apiFetch<Preferences>('/api/preferences', { method: 'PUT', body: preferences })
}
