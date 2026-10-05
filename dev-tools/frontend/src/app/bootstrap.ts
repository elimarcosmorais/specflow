import { Bootstrap } from '../../wailsjs/go/desktop/Bridge'
import type { Language } from '../shared/i18n/I18nProvider'
import type { Theme } from '../shared/theme/applyTheme'

// §3.3.1: o binding gerado tipa os campos de estado como string; aqui eles ganham os valores do contrato.
export interface BootstrapInfo {
  version: string
  apiBaseUrl: string
  token: string
  apiStatus: 'ok' | 'port_in_use' | 'failed'
  apiPort: number
  databaseStatus: 'available' | 'unavailable'
  logsDir: string
  theme: Theme
  language: Language
}

let pending: Promise<BootstrapInfo> | undefined

// Única chamada pelos bindings do Wails (ADR 004). Memoizada para o StrictMode do React
// (que roda o efeito duas vezes em desenvolvimento) não chamá-la de novo.
export function bootstrap(): Promise<BootstrapInfo> {
  pending ??= Bootstrap() as Promise<BootstrapInfo>
  return pending
}
