import type { Company } from '../../shared/types/api'

export default defineEventHandler(async (event): Promise<Company[]> => {
  try {
    return await $fetch<Company[]>(`${backendBaseUrl(event)}/api/companies`)
  } catch (error) {
    throwBackendError(error, 'Companies could not be loaded.')
  }
})
