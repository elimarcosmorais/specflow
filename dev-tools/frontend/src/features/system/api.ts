import type { components } from '../../shared/api/schema'
import { apiFetch } from '../../shared/api/client'

export type Health = components['schemas']['Health']

export function getHealth(): Promise<Health> {
  return apiFetch<Health>('/api/health')
}

export function openLogsFolder(): Promise<void> {
  return apiFetch('/api/system/open-logs-folder', { method: 'POST' })
}
