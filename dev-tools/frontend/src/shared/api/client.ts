import type { components } from './schema'

type ErrorBody = components['schemas']['ErrorResponse']

export class ApiError extends Error {
  readonly status: number
  readonly code: string
  readonly field?: string

  constructor(status: number, code: string, message: string, field?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.field = field
  }
}

// O motor não respondeu (ou respondeu algo fora do contrato da API): a tela mostra
// `errors.network`, do próprio catálogo.
export class NetworkError extends Error {
  constructor(cause: unknown) {
    super('network error', { cause })
    this.name = 'NetworkError'
  }
}

let config: { baseUrl: string; token: string } | null = null

export function configureApi(baseUrl: string, token: string): void {
  config = { baseUrl, token }
}

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: unknown
}

export async function apiFetch<T = void>(path: string, options: RequestOptions = {}): Promise<T> {
  if (!config) {
    throw new Error('apiFetch chamado antes de configureApi')
  }
  const hasBody = options.body !== undefined

  let response: Response
  try {
    response = await fetch(config.baseUrl + path, {
      method: options.method ?? 'GET',
      headers: {
        Authorization: `Bearer ${config.token}`,
        ...(hasBody && { 'Content-Type': 'application/json' }),
      },
      body: hasBody ? JSON.stringify(options.body) : undefined,
    })
  } catch (cause) {
    throw new NetworkError(cause)
  }

  if (response.status >= 400) {
    throw await toApiError(response)
  }
  if (response.status === 204) {
    return undefined as T
  }
  try {
    return (await response.json()) as T
  } catch (cause) {
    throw new NetworkError(cause)
  }
}

async function toApiError(response: Response): Promise<ApiError | NetworkError> {
  let body: Partial<ErrorBody> | null
  try {
    body = await response.json()
  } catch (cause) {
    return new NetworkError(cause)
  }
  const error = body?.error
  if (typeof error?.code !== 'string' || typeof error.message !== 'string') {
    return new NetworkError(new Error(`resposta ${response.status} fora do formato de erro da API`))
  }
  return new ApiError(
    response.status,
    error.code,
    error.message,
    typeof error.field === 'string' ? error.field : undefined,
  )
}
