<template>
  <div>
    <div class="d-flex justify-space-between align-center mb-3">
      <div class="text-subtitle-1">
        Rules are evaluated highest-priority first; the first matching rule
        overrides the GLM result.
      </div>
      <v-btn color="primary" prepend-icon="mdi-plus" @click="openCreate">
        Add Rule
      </v-btn>
    </div>

    <v-alert
      v-if="loadError"
      type="error"
      variant="tonal"
      density="compact"
      class="mb-3"
    >
      {{ loadError }}
    </v-alert>

    <v-table density="compact">
      <thead>
        <tr>
          <th>Name</th>
          <th>Risk level</th>
          <th class="text-right">Priority</th>
          <th>Enabled</th>
          <th>Updated by</th>
          <th>Updated at</th>
          <th class="text-right">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="!rules.length">
          <td colspan="7" class="text-center text-medium-emphasis py-4">
            No rules configured yet. Add one to start overriding the GLM.
          </td>
        </tr>
        <tr v-for="rule in rules" :key="rule.id">
          <td>
            <div>{{ rule.name }}</div>
            <div
              v-if="rule.description"
              class="text-caption text-medium-emphasis"
            >
              {{ rule.description }}
            </div>
          </td>
          <td>
            <v-chip :color="riskColor(rule.risk_level)" size="small" label>
              {{ rule.risk_level }}
            </v-chip>
          </td>
          <td class="text-right">{{ rule.priority }}</td>
          <td>
            <v-switch
              :model-value="rule.enabled"
              color="success"
              density="compact"
              hide-details
              @update:model-value="toggleEnabled(rule, $event)"
            />
          </td>
          <td>{{ rule.updated_by || '—' }}</td>
          <td>{{ formatDate(rule.updated_at) }}</td>
          <td class="text-right">
            <v-btn
              size="small"
              variant="text"
              icon="mdi-pencil"
              @click="openEdit(rule)"
            />
            <v-btn
              size="small"
              variant="text"
              icon="mdi-delete"
              color="error"
              @click="confirmDelete(rule)"
            />
          </td>
        </tr>
      </tbody>
    </v-table>

    <FraudRiskRuleEditor
      v-if="editorOpen"
      v-model:open="editorOpen"
      :rule="editingRule"
      :features="features"
      @saved="onSaved"
    />

    <v-dialog v-model="deleteDialog" max-width="400px">
      <v-card>
        <v-card-title>Delete rule</v-card-title>
        <v-card-text>
          Are you sure you want to delete
          <strong>{{ deletingRule?.name }}</strong
          >? This cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" :loading="deleting" @click="doDelete">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import FraudRiskService, {
  type FraudRiskRule,
  type FraudFeatureSpec
} from '@/renderer/api/FraudRiskService'
import FraudRiskRuleEditor from '@/renderer/screens/group_pricing/administration/components/FraudRiskRuleEditor.vue'

const rules = ref<FraudRiskRule[]>([])
const features = ref<FraudFeatureSpec[]>([])
const loadError = ref<string | null>(null)
const editorOpen = ref(false)
const editingRule = ref<FraudRiskRule | null>(null)
const deleteDialog = ref(false)
const deletingRule = ref<FraudRiskRule | null>(null)
const deleting = ref(false)

function riskColor(level: string) {
  return (
    {
      low: 'success',
      medium: 'warning',
      high: 'orange-darken-2',
      critical: 'error'
    }[level] || 'grey'
  )
}

function formatDate(iso: string | null | undefined) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString()
}

async function load() {
  loadError.value = null
  try {
    const [rulesRes, featuresRes] = await Promise.all([
      FraudRiskService.listRules(),
      FraudRiskService.getFeatureCatalogue()
    ])
    rules.value = rulesRes.data
    features.value = featuresRes.data
  } catch (err: any) {
    loadError.value =
      err?.response?.data || err?.message || 'Failed to load rules'
  }
}

function openCreate() {
  editingRule.value = null
  editorOpen.value = true
}

function openEdit(rule: FraudRiskRule) {
  editingRule.value = rule
  editorOpen.value = true
}

async function onSaved() {
  editorOpen.value = false
  await load()
}

async function toggleEnabled(rule: FraudRiskRule, enabled: boolean | null) {
  const next = enabled === true
  try {
    await FraudRiskService.updateRule(rule.id, { ...rule, enabled: next })
    rule.enabled = next
  } catch (err: any) {
    loadError.value =
      err?.response?.data || err?.message || 'Failed to toggle rule'
  }
}

function confirmDelete(rule: FraudRiskRule) {
  deletingRule.value = rule
  deleteDialog.value = true
}

async function doDelete() {
  if (!deletingRule.value) return
  deleting.value = true
  try {
    await FraudRiskService.deleteRule(deletingRule.value.id)
    deleteDialog.value = false
    deletingRule.value = null
    await load()
  } catch (err: any) {
    loadError.value =
      err?.response?.data || err?.message || 'Failed to delete rule'
  } finally {
    deleting.value = false
  }
}

onMounted(load)
</script>
