<template>
  <div>
    <v-alert
      v-if="!model?.trained_at"
      type="info"
      variant="tonal"
      density="compact"
      class="mb-4"
    >
      The GLM has not been trained yet. Click <strong>Refit Model</strong> once
      there are at least 100 closed assessments (30 flagged as fraud).
    </v-alert>

    <v-card v-else variant="outlined" class="mb-4">
      <v-card-text>
        <v-row>
          <v-col cols="12" md="3">
            <div class="text-caption text-medium-emphasis">Trained at</div>
            <div>{{ formatDate(model.trained_at) }}</div>
          </v-col>
          <v-col cols="12" md="3">
            <div class="text-caption text-medium-emphasis">Trained by</div>
            <div>{{ model.trained_by || '—' }}</div>
          </v-col>
          <v-col cols="12" md="2">
            <div class="text-caption text-medium-emphasis">Sample size</div>
            <div>{{ model.sample_size }}</div>
          </v-col>
          <v-col cols="12" md="2">
            <div class="text-caption text-medium-emphasis">Positives</div>
            <div>{{ model.positive_count }}</div>
          </v-col>
          <v-col cols="12" md="2">
            <div class="text-caption text-medium-emphasis">Training AUC</div>
            <div>{{ model.auc?.toFixed(3) ?? '—' }}</div>
          </v-col>
        </v-row>
      </v-card-text>
    </v-card>

    <v-alert
      v-if="refitError"
      type="error"
      variant="tonal"
      density="compact"
      class="mb-3"
    >
      {{ refitError }}
    </v-alert>

    <div class="d-flex justify-space-between align-center mb-3">
      <div class="text-subtitle-1">Coefficients</div>
      <v-btn
        color="primary"
        :loading="refitting"
        prepend-icon="mdi-refresh"
        @click="refit"
      >
        Refit Model
      </v-btn>
    </div>

    <v-table density="compact">
      <thead>
        <tr>
          <th>Feature</th>
          <th>Description</th>
          <th class="text-right">β (coefficient)</th>
          <th>Direction</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td><em>(intercept)</em></td>
          <td>Base log-odds</td>
          <td class="text-right">{{ formatBeta(model?.intercept) }}</td>
          <td>—</td>
        </tr>
        <tr v-for="row in coefficientRows" :key="row.name">
          <td>{{ row.name }}</td>
          <td class="text-medium-emphasis">{{ row.description }}</td>
          <td class="text-right">{{ formatBeta(row.beta) }}</td>
          <td>
            <v-icon
              v-if="row.beta != null"
              :color="
                row.beta > 0 ? 'error' : row.beta < 0 ? 'success' : 'grey'
              "
              size="small"
            >
              {{
                row.beta > 0
                  ? 'mdi-arrow-up-bold'
                  : row.beta < 0
                    ? 'mdi-arrow-down-bold'
                    : 'mdi-minus'
              }}
            </v-icon>
          </td>
        </tr>
      </tbody>
    </v-table>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import FraudRiskService, {
  type FraudRiskModel,
  type FraudFeatureSpec
} from '@/renderer/api/FraudRiskService'

const model = ref<FraudRiskModel | null>(null)
const features = ref<FraudFeatureSpec[]>([])
const refitting = ref(false)
const refitError = ref<string | null>(null)

const coefficientRows = computed(() => {
  return features.value
    .filter((f) => f.used_by_glm)
    .map((f) => ({
      name: f.name,
      description: f.description,
      beta: model.value?.coefficients?.[f.name] ?? null
    }))
})

function formatDate(iso: string | null) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString()
}

function formatBeta(b: number | null | undefined) {
  if (b == null) return '—'
  return b.toFixed(4)
}

async function load() {
  const [modelRes, featuresRes] = await Promise.all([
    FraudRiskService.getModel(),
    FraudRiskService.getFeatureCatalogue()
  ])
  model.value = modelRes.data
  features.value = featuresRes.data
}

async function refit() {
  refitting.value = true
  refitError.value = null
  try {
    await FraudRiskService.refitModel()
    await load()
  } catch (err: any) {
    refitError.value = err?.response?.data || err?.message || 'Refit failed'
  } finally {
    refitting.value = false
  }
}

onMounted(load)
</script>
