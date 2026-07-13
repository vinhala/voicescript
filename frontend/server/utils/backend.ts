import { createError, type H3Event } from 'h3'

export function backendBaseUrl(event: H3Event) {
  const config = useRuntimeConfig(event)
  return String(config.backendBaseUrl).replace(/\/$/, '')
}

export function throwBackendError(error: unknown, fallback: string): never {
  const backendError = error as {
    status?: number
    statusCode?: number
    data?: {
      error?: string
      message?: string
      statusMessage?: string
      data?: { error?: string }
    }
    message?: string
  }
  const message = backendError.data?.error
    ?? backendError.data?.data?.error
    ?? backendError.data?.message
    ?? backendError.data?.statusMessage
    ?? backendError.message
    ?? fallback

  throw createError({
    statusCode: backendError.statusCode ?? backendError.status ?? 500,
    statusMessage: fallback,
    message,
    data: { error: message },
  })
}
