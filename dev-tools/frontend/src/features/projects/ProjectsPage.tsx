import { useEffect, useState } from 'react'
import { ApiError } from '../../shared/api/client'
import { useErrorMessage, useT } from '../../shared/i18n/I18nProvider'
import { ConfirmDialog } from '../../shared/ui/ConfirmDialog'
import { deleteProject, listProjects } from './api'
import type { Project } from './api'
import { ProjectForm } from './ProjectForm'

type Load =
  | { state: 'loading' }
  | { state: 'error'; error: unknown }
  | { state: 'ready'; projects: Project[] }

// Formulário aberto: 'new' = cadastro; um projeto = edição dele.
type FormTarget = 'new' | Project

export function ProjectsPage() {
  const t = useT()
  const errorMessage = useErrorMessage()
  const [load, setLoad] = useState<Load>({ state: 'loading' })
  const [reloads, setReloads] = useState(0)
  const [form, setForm] = useState<FormTarget | null>(null)
  const [removing, setRemoving] = useState<Project | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  // Cada mudança em `reloads` refaz o GET. Durante um recarregamento a lista anterior
  // continua na tela; só "Tentar de novo" volta ao estado de carregando.
  useEffect(() => {
    let cancelled = false
    listProjects().then(
      (projects) => !cancelled && setLoad({ state: 'ready', projects }),
      (error: unknown) => !cancelled && setLoad({ state: 'error', error }),
    )
    return () => {
      cancelled = true
    }
  }, [reloads])

  const reload = () => setReloads((n) => n + 1)

  function retry() {
    setLoad({ state: 'loading' })
    reload()
  }

  // Projeto removido por outro caminho entre a lista e a ação (§6.2): mostra a mensagem da
  // API e recarrega a lista.
  function handleNotFound(message: string) {
    setForm(null)
    setRemoving(null)
    setNotice(message)
    reload()
  }

  function handleSaved() {
    setForm(null)
    setNotice(null)
    reload()
  }

  function openForm(target: FormTarget) {
    setNotice(null)
    setForm(target)
  }

  async function confirmRemove(project: Project) {
    try {
      await deleteProject(project.id)
    } catch (err) {
      if (err instanceof ApiError && err.code === 'project.not_found') {
        handleNotFound(err.message)
        return
      }
      throw err
    }
    setRemoving(null)
    setNotice(null)
    reload()
  }

  return (
    <section>
      <div className="page-header">
        <h1>{t('projects.title')}</h1>
        {load.state === 'ready' && (
          <button type="button" className="button-primary" onClick={() => openForm('new')}>
            {t('projects.new')}
          </button>
        )}
      </div>
      {notice && (
        <p className="alert alert-danger" role="alert">
          {notice}
        </p>
      )}
      {load.state === 'loading' && <p>{t('projects.loading')}</p>}
      {load.state === 'error' && (
        <>
          <p className="alert alert-danger" role="alert">
            {errorMessage(load.error)}
          </p>
          <button type="button" onClick={retry}>
            {t('projects.retry')}
          </button>
        </>
      )}
      {load.state === 'ready' && load.projects.length === 0 && <p>{t('projects.empty')}</p>}
      {load.state === 'ready' && load.projects.length > 0 && (
        <table className="table">
          <thead>
            <tr>
              <th>{t('projects.column.name')}</th>
              <th>{t('projects.column.type')}</th>
              <th>{t('projects.column.path')}</th>
              <th>
                <span className="visually-hidden">{t('projects.column.actions')}</span>
              </th>
            </tr>
          </thead>
          <tbody>
            {load.projects.map((project) => (
              <tr key={project.id}>
                <td>{project.name}</td>
                <td>{t(`projects.type.${project.type}`)}</td>
                <td className="cell-path">
                  {project.path}
                  {!project.folderExists && (
                    <span className="badge badge-warning">{t('projects.folderMissing')}</span>
                  )}
                </td>
                <td className="cell-actions">
                  <button type="button" onClick={() => openForm(project)}>
                    {t('projects.edit')}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setNotice(null)
                      setRemoving(project)
                    }}
                  >
                    {t('projects.remove')}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {form && (
        <ProjectForm
          project={form === 'new' ? null : form}
          onSaved={handleSaved}
          onNotFound={handleNotFound}
          onCancel={() => setForm(null)}
        />
      )}
      {removing && (
        <ConfirmDialog
          title={t('projects.confirm.title')}
          message={t('projects.confirm.message', { name: removing.name })}
          confirmLabel={t('projects.remove')}
          busyLabel={t('projects.confirm.removing')}
          onConfirm={() => confirmRemove(removing)}
          onCancel={() => setRemoving(null)}
        />
      )}
    </section>
  )
}
