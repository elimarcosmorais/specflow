import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import type { components } from '../api/schema'
import { ApiError } from '../api/client'
import { en } from './en'
import { ptBr } from './pt-br'
import type { Messages } from './pt-br'

export type Language = components['schemas']['Preferences']['language']

const catalogs: Record<Language, Messages> = { 'pt-br': ptBr, en }

type Params = Record<string, string | number>
type Translate = (key: keyof Messages, params?: Params) => string

interface I18nValue {
  language: Language
  setLanguage: (language: Language) => void
  t: Translate
}

const I18nContext = createContext<I18nValue | null>(null)

// Até o bootstrap entregar o idioma salvo, segue o idioma do WebView (mesma regra do motor).
function initialLanguage(): Language {
  return navigator.language.toLowerCase() === 'pt-br' ? 'pt-br' : 'en'
}

export function I18nProvider({ children }: { children: ReactNode }) {
  const [language, setLanguage] = useState<Language>(initialLanguage)

  useEffect(() => {
    document.documentElement.lang = language === 'pt-br' ? 'pt-BR' : 'en'
  }, [language])

  const t = useCallback<Translate>(
    (key, params) =>
      catalogs[language][key].replace(/\{(\w+)\}/g, (placeholder, name: string) =>
        params && name in params ? String(params[name]) : placeholder,
      ),
    [language],
  )

  const value = useMemo(() => ({ language, setLanguage, t }), [language, t])
  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>
}

function useI18n(): I18nValue {
  const value = useContext(I18nContext)
  if (!value) {
    throw new Error('useT/useLanguage fora do I18nProvider')
  }
  return value
}

export function useT(): Translate {
  return useI18n().t
}

export function useLanguage(): Pick<I18nValue, 'language' | 'setLanguage'> {
  const { language, setLanguage } = useI18n()
  return { language, setLanguage }
}

// Erro da API: a mensagem já vem traduzida pelo motor. Qualquer outra falha é do canal
// (o motor não respondeu), e o texto é do catálogo do front (ADR 005).
export function useErrorMessage(): (error: unknown) => string {
  const t = useT()
  return useCallback(
    (error) => (error instanceof ApiError ? error.message : t('errors.network')),
    [t],
  )
}
