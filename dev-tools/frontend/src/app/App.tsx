import { useEffect, useState } from 'react'
import { PreferencesPage } from '../features/preferences/PreferencesPage'
import { ProjectsPage } from '../features/projects/ProjectsPage'
import { StatusPage } from '../features/system/StatusPage'
import { configureApi } from '../shared/api/client'
import { useLanguage, useT } from '../shared/i18n/I18nProvider'
import { applyTheme } from '../shared/theme/applyTheme'
import type { Theme } from '../shared/theme/applyTheme'
import { bootstrap } from './bootstrap'
import type { BootstrapInfo } from './bootstrap'

type Screen = 'projects' | 'preferences' | 'status'

const screens: { screen: Screen; label: 'nav.projects' | 'nav.preferences' | 'nav.status' }[] = [
  { screen: 'projects', label: 'nav.projects' },
  { screen: 'preferences', label: 'nav.preferences' },
  { screen: 'status', label: 'nav.status' },
]

type Boot = { state: 'loading' } | { state: 'failed' } | { state: 'ready'; info: BootstrapInfo }

export function App() {
  const t = useT()
  const { setLanguage } = useLanguage()
  const [boot, setBoot] = useState<Boot>({ state: 'loading' })

  useEffect(() => {
    let cancelled = false
    bootstrap().then(
      (info) => {
        if (cancelled) {
          return
        }
        configureApi(info.apiBaseUrl, info.token)
        setLanguage(info.language)
        applyTheme(info.theme)
        setBoot({ state: 'ready', info })
      },
      () => !cancelled && setBoot({ state: 'failed' }),
    )
    return () => {
      cancelled = true
    }
  }, [setLanguage])

  if (boot.state === 'loading') {
    return null
  }
  if (boot.state === 'failed') {
    return (
      <p className="alert alert-danger" role="alert">
        {t('bootstrap.failed')}
      </p>
    )
  }
  return <Shell info={boot.info} />
}

function Shell({ info }: { info: BootstrapInfo }) {
  const t = useT()
  const apiUp = info.apiStatus === 'ok'
  const [screen, setScreen] = useState<Screen>(apiUp ? 'projects' : 'status')
  const [theme, setTheme] = useState<Theme>(info.theme)

  return (
    <div className="shell">
      <nav className="shell-menu">
        {screens.map(({ screen: item, label }) => {
          const needsApi = item !== 'status'
          const disabled = needsApi && !apiUp
          return (
            <button
              key={item}
              type="button"
              className="shell-menu-item"
              aria-current={screen === item ? 'page' : undefined}
              disabled={disabled}
              title={disabled ? t('nav.apiUnavailable') : undefined}
              onClick={() => setScreen(item)}
            >
              {t(label)}
            </button>
          )
        })}
      </nav>
      <main className="shell-content">
        {screen === 'projects' && <ProjectsPage />}
        {screen === 'preferences' && <PreferencesPage theme={theme} onThemeChange={setTheme} />}
        {screen === 'status' && <StatusPage info={info} />}
      </main>
    </div>
  )
}
