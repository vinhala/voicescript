import type { Opportunity } from '../../../../shared/types/api'

export default defineEventHandler(async (event): Promise<Opportunity[]> => {
  const companyId = getRouterParam(event, 'companyId')
  if (!companyId) {
    throw createError({ statusCode: 400, statusMessage: 'company id is required' })
  }

  try {
    return await $fetch<Opportunity[]>(
      `${backendBaseUrl(event)}/api/companies/${encodeURIComponent(companyId)}/opportunities`,
    )
  } catch (error) {
    throwBackendError(error, 'Opportunities could not be loaded.')
  }
})
