<script setup lang="ts">
import type { AnalysisResultResponse, Company, Opportunity } from '../../shared/types/api'

const companies = ref<Company[]>([])
const opportunities = ref<Opportunity[]>([])
const selectedCompanyId = ref('')
const selectedOpportunityId = ref('')
const recording = ref<File | null>(null)
const loadingCompanies = ref(true)
const loadingOpportunities = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const analysisResult = ref<AnalysisResultResponse | null>(null)

const canSubmit = computed(() => {
  return Boolean(selectedCompanyId.value && selectedOpportunityId.value && recording.value && !submitting.value)
})

onMounted(async () => {
  try {
    companies.value = await $fetch<Company[]>('/api/companies')
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, 'Companies could not be loaded.')
  } finally {
    loadingCompanies.value = false
  }
})

watch(selectedCompanyId, async (companyId) => {
  selectedOpportunityId.value = ''
  opportunities.value = []
  successMessage.value = ''
  analysisResult.value = null

  if (!companyId) {
    return
  }

  loadingOpportunities.value = true
  errorMessage.value = ''

  try {
    opportunities.value = await $fetch<Opportunity[]>(`/api/companies/${encodeURIComponent(companyId)}/opportunities`)
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, 'Opportunities could not be loaded.')
  } finally {
    loadingOpportunities.value = false
  }
})

function onRecordingChange(event: Event) {
  errorMessage.value = ''
  successMessage.value = ''
  analysisResult.value = null

  const input = event.target as HTMLInputElement
  const file = input.files?.[0] ?? null
  if (!file) {
    recording.value = null
    return
  }

  const extension = file.name.split('.').pop()?.toLowerCase()
  if (extension !== 'mp3' && extension !== 'mp4') {
    recording.value = null
    input.value = ''
    errorMessage.value = 'Please select an MP3 or MP4 recording.'
    return
  }

  recording.value = file
}

async function submitAnalysis() {
  if (!recording.value || !canSubmit.value) {
    return
  }

  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  analysisResult.value = null

  const formData = new FormData()
  formData.append('companyId', selectedCompanyId.value)
  formData.append('opportunityId', selectedOpportunityId.value)
  formData.append('recording', recording.value)

  try {
    const response = await $fetch<AnalysisResultResponse>('/api/analyses', {
      method: 'POST',
      body: formData,
    })
    analysisResult.value = response
    successMessage.value = `Analysis ${response.analysisId} completed and was stored as note ${response.noteId}.`
  } catch (error) {
    errorMessage.value = apiErrorMessage(error, 'Analysis could not be started.')
  } finally {
    submitting.value = false
  }
}

function apiErrorMessage(error: unknown, fallback: string) {
  if (typeof error === 'object' && error && 'data' in error) {
    const data = (error as { data?: { error?: string, data?: { error?: string } } }).data
    if (data?.error || data?.data?.error) {
      return data.error ?? data.data?.error ?? fallback
    }
  }

  return fallback
}
</script>

<template>
  <main class="page-shell">
    <section class="workspace">
      <header class="header">
        <p class="eyebrow">DieStimme</p>
        <h1>Start a requirements session analysis</h1>
        <p class="subtitle">
          Select a CRM company and opportunity, attach the session recording, and generate the
          onboarding questionnaire in Twenty CRM.
        </p>
      </header>

      <form class="panel" @submit.prevent="submitAnalysis">
        <div class="form-grid">
          <div class="field">
            <label for="company">Company</label>
            <select id="company" v-model="selectedCompanyId" :disabled="loadingCompanies || submitting">
              <option value="">
                {{ loadingCompanies ? 'Loading companies...' : 'Select a company' }}
              </option>
              <option v-for="company in companies" :key="company.id" :value="company.id">
                {{ company.name }}
              </option>
            </select>
            <span class="hint">Companies are loaded from Twenty CRM.</span>
          </div>

          <div class="field">
            <label for="opportunity">Opportunity</label>
            <select
              id="opportunity"
              v-model="selectedOpportunityId"
              :disabled="!selectedCompanyId || loadingOpportunities || submitting"
            >
              <option value="">
                {{ loadingOpportunities ? 'Loading opportunities...' : 'Select an opportunity' }}
              </option>
              <option v-for="opportunity in opportunities" :key="opportunity.id" :value="opportunity.id">
                {{ opportunity.name }}
              </option>
            </select>
            <span class="hint">The opportunity list is scoped to the selected company.</span>
          </div>

          <div class="field field-wide">
            <label for="recording">Recording</label>
            <input
              id="recording"
              type="file"
              accept=".mp3,.mp4,audio/mpeg,audio/mp3,audio/mp4,video/mp4"
              :disabled="submitting"
              @change="onRecordingChange"
            >
            <span class="hint">Upload an MP3 or MP4 recording from the requirements elicitation session.</span>
          </div>
        </div>

        <div class="actions">
          <div
            class="status"
            :class="{ 'status-error': errorMessage, 'status-success': successMessage }"
            aria-live="polite"
          >
            {{ errorMessage || successMessage || (submitting ? 'Analyzing the recording...' : 'Ready to analyze once all fields are complete.') }}
          </div>
          <button class="submit-button" type="submit" :disabled="!canSubmit">
            {{ submitting ? 'Analyzing...' : 'Start analysis' }}
          </button>
        </div>
      </form>

      <section v-if="analysisResult" class="result-panel" aria-live="polite">
        <div class="result-header">
          <div>
            <p class="eyebrow">Completed questionnaire</p>
            <h2>Stored in Twenty CRM</h2>
          </div>
          <span class="note-id">Note {{ analysisResult.noteId }}</span>
        </div>
        <pre class="markdown-output">{{ analysisResult.questionnaireMarkdown }}</pre>
      </section>
    </section>
  </main>
</template>
