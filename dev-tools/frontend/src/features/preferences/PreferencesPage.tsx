import { useState } from 'react'
import type { Language } from '../../shared/i18n/I18nProvider'
import { useErrorMessage, useLanguage, useT } from '../../shared/i18n/I18nProvider'
import { applyTheme } from '../../shared/theme/applyTheme'
import type { Theme } from '../../shared/theme/applyTheme'
import { savePreferences } from './api'
import type { Preferences } from './api'

const themes = ['light', 'dark', 'system'] as const
const languages: { value: Language; label: string }[] = [
  { value: 'pt-br', label: 'Português (Brasil)' },
  { value: 'en', label: 'English' },
]

interface Props {
  theme: Theme
  onThemeChange: (theme: Theme) => void
}

export function PreferencesPage({ theme, onThemeChange }: Props) {
  const t = useT()
  const errorMessage = useErrorMessage()
  const { language, setLanguage } = useLanguage()
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // As opções vêm do estado confirmado (tema do shell, idioma do provider): só mudam ao
  // receber 200. Em erro, nada mudou e a opção anterior continua marcada.
  async function save(next: Preferences) {
    setSaving(true)
    setError(null)
    try {
      const saved = await savePreferences(next)
      applyTheme(saved.theme)
      setLanguage(saved.language)
      onThemeChange(saved.theme)
    } catch (err) {
      setError(errorMessage(err))
    } finally {
      setSaving(false)
    }
  }

  return (
    <section>
      <h1>{t('preferences.title')}</h1>
      {error && (
        <p className="alert alert-danger" role="alert">
          {error}
        </p>
      )}
      <fieldset className="field-group" disabled={saving}>
        <legend>{t('preferences.theme')}</legend>
        {themes.map((option) => (
          <label key={option} className="option">
            <input
              type="radio"
              name="theme"
              value={option}
              checked={theme === option}
              onChange={() => void save({ theme: option, language })}
            />
            {t(`preferences.theme.${option}`)}
          </label>
        ))}
      </fieldset>
      <fieldset className="field-group" disabled={saving}>
        <legend>{t('preferences.language')}</legend>
        <select
          aria-label={t('preferences.language')}
          value={language}
          onChange={(event) => void save({ theme, language: event.target.value as Language })}
        >
          {languages.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </fieldset>
    </section>
  )
}
