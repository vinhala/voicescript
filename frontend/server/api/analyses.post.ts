import type { AnalysisResultResponse } from '../../shared/types/api'

export default defineEventHandler(async (event): Promise<AnalysisResultResponse> => {
  const parts = await readMultipartFormData(event)
  if (!parts) {
    throw createError({ statusCode: 400, statusMessage: 'multipart form data is required' })
  }

  const formData = new FormData()

  for (const part of parts) {
    if (!part.name) {
      continue
    }

    if (part.filename) {
      formData.append(
        part.name,
        new Blob([part.data], { type: part.type || 'application/octet-stream' }),
        part.filename,
      )
      continue
    }

    formData.append(part.name, part.data.toString('utf8'))
  }

  try {
    return await $fetch<AnalysisResultResponse>(`${backendBaseUrl(event)}/api/analyses`, {
      method: 'POST',
      body: formData,
    })
  } catch (error) {
    throwBackendError(error, 'Analysis could not be completed.')
  }
})
