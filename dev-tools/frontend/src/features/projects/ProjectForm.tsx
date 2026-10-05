import { useEffect, useRef, useState } from 'react'
import type { FormEvent } from 'react'
import { ApiError } from '../../shared/api/client'
import { useErrorMessage, useT } from '../../shared/i18n/I18nProvider'
import { createProject, updateProject } from './api'
import type { Project, ProjectType } from './api'

const types: ProjectType[] = ['web', 'android', 'windows']
type Field = 'name' | 'path' | 'type'

function isField(field: string | undefined): field is Field {
  return field === 'name' || field === 'path' || field === 'type'
}

interface Props {
  // null = novo projeto; senão, edição desse projeto.
  project: Project | null
  onSaved: () => void
  // O projeto sumiu entre a lista e o salvar (`project.not_found`): quem chamou mostra a
  // mensagem e recarrega a lista.
  onNotFound: (message: string) => void
  onCancel: () => void
}

export function ProjectForm({ project, onSaved, onNotFound, onCancel }: Props) {
  const t = useT()
  const errorMessage = useErrorMessage()
  const dialogRef = useRef<HTMLDialogElement>(null)
  const [name, setName] = useState(project?.name ?? '')
  const [path, setPath] = useState(project?.path ?? '')
  const [type, setType] = useState<ProjectType | ''>(project?.type ?? '')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<unknown>(null)

  useEffect(() => {
    const dialog = dialogRef.current
    dialog?.showModal()
    return () => dialog?.close()
  }, [])

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setSaving(true)
    setError(null)
    const input = { name, path, type }
    try {
      if (project) {
        await updateProject(project.id, input)
      } else {
        await createProject(input)
      }
    } catch (err) {
      if (err instanceof ApiError && err.code === 'project.not_found') {
        onNotFound(err.message)
        return
      }
      setError(err)
      setSaving(false)
      return
    }
    onSaved()
  }

  const fieldError = error instanceof ApiError && isField(error.field) ? error : null
  const generalError = error !== null && !fieldError ? errorMessage(error) : null
  const nameError = fieldError?.field === 'name' ? fieldError.message : null
  const pathError = fieldError?.field === 'path' ? fieldError.message : null
  const typeError = fieldError?.field === 'type' ? fieldError.message : null

  return (
    <dialog
      ref={dialogRef}
      className="modal"
      aria-labelledby="project-form-title"
      onCancel={(event) => {
        event.preventDefault()
        if (!saving) {
          onCancel()
        }
      }}
    >
      <h2 id="project-form-title">
        {t(project ? 'projects.form.editTitle' : 'projects.form.newTitle')}
      </h2>
      {generalError && (
        <p className="alert alert-danger" role="alert">
          {generalError}
        </p>
      )}
      <form onSubmit={(event) => void handleSubmit(event)} noValidate>
        <fieldset className="form-fields" disabled={saving}>
          <div className="form-field">
            <label htmlFor="project-name">{t('projects.field.name')}</label>
            <input
              id="project-name"
              type="text"
              value={name}
              aria-invalid={nameError !== null}
              aria-describedby={nameError ? 'project-name-error' : undefined}
              onChange={(event) => setName(event.target.value)}
            />
            {nameError && (
              <p id="project-name-error" className="field-error">
                {nameError}
              </p>
            )}
          </div>
          <div className="form-field">
            <label htmlFor="project-path">{t('projects.field.path')}</label>
            <input
              id="project-path"
              type="text"
              value={path}
              aria-invalid={pathError !== null}
              aria-describedby={pathError ? 'project-path-hint project-path-error' : 'project-path-hint'}
              onChange={(event) => setPath(event.target.value)}
            />
            <p id="project-path-hint" className="field-hint">
              {t('projects.field.pathHint')}
            </p>
            {pathError && (
              <p id="project-path-error" className="field-error">
                {pathError}
              </p>
            )}
          </div>
          <div className="form-field">
            <label htmlFor="project-type">{t('projects.field.type')}</label>
            <select
              id="project-type"
              value={type}
              aria-invalid={typeError !== null}
              aria-describedby={typeError ? 'project-type-error' : undefined}
              onChange={(event) => setType(event.target.value as ProjectType | '')}
            >
              <option value="">{t('projects.type.placeholder')}</option>
              {types.map((option) => (
                <option key={option} value={option}>
                  {t(`projects.type.${option}`)}
                </option>
              ))}
            </select>
            {typeError && (
              <p id="project-type-error" className="field-error">
                {typeError}
              </p>
            )}
          </div>
          <div className="modal-actions">
            <button type="button" onClick={onCancel}>
              {t('common.cancel')}
            </button>
            <button type="submit" className="button-primary">
              {saving ? t('projects.form.saving') : t('projects.form.save')}
            </button>
          </div>
        </fieldset>
      </form>
    </dialog>
  )
}
