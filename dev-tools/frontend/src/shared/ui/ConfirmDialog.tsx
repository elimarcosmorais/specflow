import { useEffect, useRef, useState } from 'react'
import { useErrorMessage, useT } from '../i18n/I18nProvider'

interface Props {
  title: string
  message: string
  confirmLabel: string
  busyLabel: string
  // Resolver significa que a ação deu certo e quem chamou fecha o diálogo; rejeitar mostra a
  // mensagem do erro no próprio diálogo, que continua aberto.
  onConfirm: () => Promise<void>
  onCancel: () => void
}

export function ConfirmDialog({ title, message, confirmLabel, busyLabel, onConfirm, onCancel }: Props) {
  const t = useT()
  const errorMessage = useErrorMessage()
  const dialogRef = useRef<HTMLDialogElement>(null)
  const cancelRef = useRef<HTMLButtonElement>(null)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const dialog = dialogRef.current
    dialog?.showModal()
    // O foco inicial fica no "Cancelar": confirmar uma remoção exige uma ação deliberada.
    cancelRef.current?.focus()
    return () => dialog?.close()
  }, [])

  async function handleConfirm() {
    setBusy(true)
    setError(null)
    try {
      await onConfirm()
    } catch (err) {
      setError(errorMessage(err))
      setBusy(false)
    }
  }

  return (
    <dialog
      ref={dialogRef}
      className="modal"
      aria-labelledby="confirm-dialog-title"
      onCancel={(event) => {
        event.preventDefault()
        if (!busy) {
          onCancel()
        }
      }}
    >
      <h2 id="confirm-dialog-title">{title}</h2>
      {error && (
        <p className="alert alert-danger" role="alert">
          {error}
        </p>
      )}
      <p>{message}</p>
      <div className="modal-actions">
        <button ref={cancelRef} type="button" disabled={busy} onClick={onCancel}>
          {t('common.cancel')}
        </button>
        <button
          type="button"
          className="button-danger"
          disabled={busy}
          onClick={() => void handleConfirm()}
        >
          {busy ? busyLabel : confirmLabel}
        </button>
      </div>
    </dialog>
  )
}
