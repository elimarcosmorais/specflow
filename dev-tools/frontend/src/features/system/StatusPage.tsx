import { useEffect, useState } from 'react'
import type { BootstrapInfo } from '../../app/bootstrap'
import { useErrorMessage, useT } from '../../shared/i18n/I18nProvider'
import { getHealth, openLogsFolder } from './api'
import type { Health } from './api'

export function StatusPage({ info }: { info: BootstrapInfo }) {
  const t = useT()
  const errorMessage = useErrorMessage()
  const apiUp = info.apiStatus === 'ok'
  const [health, setHealth] = useState<Health | null>(null)
  const [healthError, setHealthError] = useState<string | null>(null)
  const [opening, setOpening] = useState(false)
  const [openError, setOpenError] = useState<string | null>(null)

  // Com a API no ar, a saúde vem de GET /api/health ao abrir a tela. Sem ela, o bootstrap
  // já tem versão e banco (§7.2).
  useEffect(() => {
    if (!apiUp) {
      return
    }
    let cancelled = false
    getHealth().then(
      (result) => !cancelled && setHealth(result),
      (err) => !cancelled && setHealthError(errorMessage(err)),
    )
    return () => {
      cancelled = true
    }
    // errorMessage fica fora das dependências: trocar o idioma não deve refazer o GET.
  }, [apiUp])

  async function handleOpenLogs() {
    setOpening(true)
    setOpenError(null)
    try {
      await openLogsFolder()
    } catch (err) {
      setOpenError(errorMessage(err))
    } finally {
      setOpening(false)
    }
  }

  const version = health?.version ?? info.version
  const database = apiUp ? health?.database : info.databaseStatus
  const warning =
    info.apiStatus === 'port_in_use'
      ? t('status.portInUse', { apiPort: info.apiPort, logsDir: info.logsDir })
      : info.apiStatus === 'failed'
        ? t('status.apiFailed', { logsDir: info.logsDir })
        : null

  return (
    <section>
      <h1>{t('status.title')}</h1>
      {warning && (
        <p className="alert alert-warning" role="alert">
          {warning}
        </p>
      )}
      {healthError && (
        <p className="alert alert-danger" role="alert">
          {healthError}
        </p>
      )}
      <dl className="details">
        <dt>{t('status.version')}</dt>
        <dd>{version}</dd>
        <dt>{t('status.database')}</dt>
        <dd className={database === 'unavailable' ? 'alert alert-warning' : undefined}>
          {database ? t(`status.database.${database}`) : healthError ? '—' : t('status.loading')}
        </dd>
        <dt>{t('status.apiAddress')}</dt>
        <dd>{info.apiBaseUrl}</dd>
        <dt>{t('status.docsAddress')}</dt>
        <dd>{`${info.apiBaseUrl}/docs`}</dd>
        <dt>{t('status.logsFolder')}</dt>
        <dd>{info.logsDir}</dd>
      </dl>
      <button type="button" disabled={!apiUp || opening} onClick={() => void handleOpenLogs()}>
        {t('status.openLogs')}
      </button>
      {openError && (
        <p className="alert alert-danger" role="alert">
          {openError}
        </p>
      )}
    </section>
  )
}
