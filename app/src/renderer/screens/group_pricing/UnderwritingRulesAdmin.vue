<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Underwriting Rules</span>
          </template>

          <template #default>
            <!-- Action toolbar — moved out of the dark headline bar so
                 the buttons are clearly visible. -->
            <v-row class="mb-3">
              <v-col class="d-flex flex-wrap gap-2">
                <v-btn
                  size="small"
                  color="primary"
                  variant="flat"
                  prepend-icon="mdi-plus"
                  @click="newRuleSetDialog = true"
                  >New rule set</v-btn
                >
                <v-btn
                  size="small"
                  color="secondary"
                  variant="tonal"
                  prepend-icon="mdi-upload"
                  @click="importDialog = true"
                  >Import CSV</v-btn
                >
                <v-btn
                  v-if="!ruleSets.length"
                  size="small"
                  color="info"
                  variant="tonal"
                  prepend-icon="mdi-magic-staff"
                  :loading="seedingStarter"
                  @click="seedStarter"
                  >Load starter template</v-btn
                >
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="4">
                <v-card variant="outlined" rounded="lg">
                  <v-card-title class="font-weight-bold">Versions</v-card-title>
                  <v-card-text v-if="ruleSets.length">
                    <v-list density="compact" :selected="selectedSetId">
                      <v-list-item
                        v-for="rs in ruleSets"
                        :key="rs.id"
                        :value="rs.id"
                        :active="selectedSetId.includes(rs.id)"
                        @click="selectRuleSet(rs.id)"
                      >
                        <template #title>
                          {{ rs.name }} v{{ rs.version }}
                        </template>
                        <template #subtitle>
                          {{ rs.active ? 'Active' : 'Inactive' }} ·
                          {{ formatDate(rs.creation_date) }} ·
                          {{ (rs.rules || []).length }} rules
                        </template>
                        <template #append>
                          <div class="d-flex align-center gap-1">
                            <v-chip
                              v-if="rs.active"
                              color="success"
                              variant="tonal"
                              size="x-small"
                              >Live</v-chip
                            >
                            <v-btn
                              v-else
                              size="x-small"
                              variant="text"
                              color="primary"
                              @click.stop="activate(rs.id)"
                              >Activate</v-btn
                            >
                            <v-btn
                              size="x-small"
                              variant="text"
                              icon="mdi-content-copy"
                              title="Duplicate as new version"
                              @click.stop="duplicate(rs.id)"
                            />
                            <v-btn
                              v-if="!rs.active"
                              size="x-small"
                              variant="text"
                              color="error"
                              icon="mdi-delete-outline"
                              title="Delete (only if not referenced)"
                              @click.stop="confirmDeleteRuleSet(rs)"
                            />
                          </div>
                        </template>
                      </v-list-item>
                    </v-list>
                  </v-card-text>
                  <v-card-text v-else class="text-center py-6">
                    <v-icon size="48" color="grey-lighten-1" class="mb-3"
                      >mdi-clipboard-list-outline</v-icon
                    >
                    <p class="text-subtitle-2 font-weight-bold mb-1"
                      >No rule sets yet</p
                    >
                    <p class="text-caption text-grey mb-4"
                      >Pick one of the buttons above to get started — the
                      starter template is the fastest path.</p
                    >
                    <v-btn
                      size="small"
                      color="info"
                      variant="tonal"
                      block
                      prepend-icon="mdi-magic-staff"
                      :loading="seedingStarter"
                      class="mb-2"
                      @click="seedStarter"
                      >Load starter template</v-btn
                    >
                    <v-btn
                      size="small"
                      color="primary"
                      variant="text"
                      block
                      prepend-icon="mdi-plus"
                      @click="newRuleSetDialog = true"
                      >Create rule set</v-btn
                    >
                    <v-btn
                      size="small"
                      color="secondary"
                      variant="text"
                      block
                      prepend-icon="mdi-upload"
                      @click="importDialog = true"
                      >Import CSV</v-btn
                    >
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="8">
                <v-card v-if="selectedSet" variant="outlined" rounded="lg">
                  <v-card-title class="d-flex align-center flex-wrap gap-2"
                    ><span class="font-weight-bold"
                      >{{ selectedSet.name }} v{{
                        selectedSet.version
                      }}</span
                    >
                    <v-chip
                      v-if="selectedSet.active"
                      color="success"
                      variant="tonal"
                      size="x-small"
                      >Live</v-chip
                    >
                    <v-spacer />
                    <v-btn
                      size="small"
                      variant="text"
                      prepend-icon="mdi-download"
                      @click="exportCsv"
                      >Export CSV</v-btn
                    >
                    <v-btn
                      size="small"
                      color="primary"
                      variant="tonal"
                      prepend-icon="mdi-plus"
                      @click="openNewRuleDialog"
                      >Add rule</v-btn
                    >
                  </v-card-title>
                  <v-card-text>
                    <v-table density="compact">
                      <thead>
                        <tr>
                          <th>Category</th>
                          <th>Field</th>
                          <th>Op</th>
                          <th>Condition</th>
                          <th>Outcome</th>
                          <th class="text-right">Loading</th>
                          <th>Exclusion</th>
                          <th class="text-right">Priority</th>
                          <th></th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr
                          v-for="r in selectedSet.rules || []"
                          :key="r.id"
                        >
                          <td>{{ r.category }}</td>
                          <td>{{ r.field }}</td>
                          <td>{{ r.op }}</td>
                          <td class="text-caption text-grey">{{
                            r.condition_json
                          }}</td>
                          <td>
                            <v-chip
                              size="x-small"
                              :color="outcomeColor(r.outcome)"
                              variant="tonal"
                              >{{ r.outcome }}</v-chip
                            >
                          </td>
                          <td class="text-right">{{ r.loading_percent || '—' }}</td>
                          <td>{{ r.exclusion_code || '—' }}</td>
                          <td class="text-right">{{ r.priority }}</td>
                          <td>
                            <div class="d-flex">
                              <v-btn
                                size="x-small"
                                variant="text"
                                color="primary"
                                icon="mdi-pencil-outline"
                                title="Edit rule"
                                @click="openEditRuleDialog(r)"
                              />
                              <v-btn
                                size="x-small"
                                variant="text"
                                color="error"
                                icon="mdi-delete-outline"
                                title="Delete rule"
                                @click="deleteRule(r)"
                              />
                            </div>
                          </td>
                        </tr>
                        <tr v-if="!(selectedSet.rules || []).length">
                          <td colspan="9" class="text-grey text-center"
                            >No rules in this set.</td
                          >
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>

                <v-card
                  v-if="selectedSet"
                  variant="outlined"
                  rounded="lg"
                  class="mt-4"
                >
                  <v-card-title class="font-weight-bold"
                    >Dry-run</v-card-title
                  >
                  <v-card-text>
                    <p class="text-caption text-grey mb-2"
                      >Enter sample disclosure values and see how the engine
                      decides.</p
                    >
                    <v-row dense>
                      <v-col cols="6" md="3">
                        <v-text-field
                          v-model.number="dryRun.bmi"
                          label="BMI"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6" md="3">
                        <v-text-field
                          v-model.number="dryRun.age"
                          label="Age"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6" md="3">
                        <v-text-field
                          v-model.number="dryRun.occupation_class"
                          label="Occupation class"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6" md="3">
                        <v-checkbox
                          v-model="dryRun.smoker"
                          label="Smoker"
                          density="compact"
                          hide-details
                        />
                      </v-col>
                    </v-row>
                    <v-btn
                      color="primary"
                      :loading="dryRunBusy"
                      @click="runDryRun"
                      >Evaluate</v-btn
                    >
                    <div v-if="dryRunResult" class="mt-3">
                      <v-chip
                        :color="outcomeColor(dryRunResult.outcome)"
                        variant="tonal"
                        size="small"
                        class="mr-2"
                        >Outcome: {{ dryRunResult.outcome }}</v-chip
                      >
                      <v-chip
                        v-if="dryRunResult.max_loading"
                        color="warning"
                        variant="tonal"
                        size="small"
                        class="mr-2"
                        >Loading: {{ dryRunResult.max_loading }}%</v-chip
                      >
                      <v-chip
                        v-for="ex in dryRunResult.exclusions || []"
                        :key="ex"
                        color="grey"
                        variant="tonal"
                        size="small"
                        class="mr-2"
                        >{{ ex }}</v-chip
                      >
                      <v-table density="compact" class="mt-2">
                        <thead>
                          <tr>
                            <th>Rule</th>
                            <th>Field</th>
                            <th>Outcome</th>
                            <th class="text-right">Loading</th>
                            <th>Exclusion</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="o in dryRunResult.outcomes || []"
                            :key="o.rule_id"
                          >
                            <td>{{ o.rule_id }}</td>
                            <td>{{ o.field }}</td>
                            <td>{{ o.outcome }}</td>
                            <td class="text-right">{{ o.loading_percent || '—' }}</td>
                            <td>{{ o.exclusion_code || '—' }}</td>
                          </tr>
                          <tr v-if="!(dryRunResult.outcomes || []).length">
                            <td colspan="5" class="text-grey text-center"
                              >No rules matched.</td
                            >
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </v-card-text>
                </v-card>

                <empty-state
                  v-if="!selectedSet"
                  icon="mdi-clipboard-list-outline"
                  title="Pick a rule set"
                  message="Select a version on the left to edit rules and run dry-runs."
                />
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- New rule set -->
    <v-dialog v-model="newRuleSetDialog" max-width="480">
      <v-card>
        <v-card-title>New rule set</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newRuleSet.name"
            label="Name"
            density="compact"
            variant="outlined"
          />
          <v-text-field
            v-model.number="newRuleSet.version"
            label="Version"
            type="number"
            density="compact"
            variant="outlined"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="newRuleSetDialog = false">Cancel</v-btn>
          <v-btn color="primary" @click="createRuleSet">Create</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- New / Edit rule -->
    <v-dialog v-model="newRuleDialog" max-width="640">
      <v-card>
        <v-card-title>{{ editingRuleId ? 'Edit rule' : 'Add rule' }}</v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col cols="12">
              <v-select
                v-model="newRule.category"
                :items="categoryOptions"
                item-title="title"
                item-value="value"
                label="Category"
                density="compact"
                variant="outlined"
                @update:model-value="onCategoryChange"
              />
              <v-alert
                v-if="selectedCategoryInfo"
                type="info"
                variant="tonal"
                density="compact"
                class="mt-1 mb-1"
              >
                <div class="text-caption">
                  <strong>What this category is for: </strong>{{
                    selectedCategoryInfo.description
                  }}
                </div>
                <div class="text-caption mt-1">
                  <strong>Standard rate assumption: </strong>{{
                    selectedCategoryInfo.standardAssumption
                  }}
                </div>
              </v-alert>
            </v-col>
            <v-col cols="12">
              <v-combobox
                v-model="newRule.field"
                :items="fieldsForCategory"
                item-title="title"
                item-value="value"
                :return-object="false"
                label="Field"
                placeholder="Pick a field or type a custom one"
                density="compact"
                variant="outlined"
                hide-no-data
              />
              <v-alert
                v-if="selectedFieldInfo"
                type="info"
                variant="tonal"
                density="compact"
                class="mt-1 mb-1"
                icon="mdi-information-outline"
              >
                <div class="text-caption">
                  <strong>Standard reference: </strong>{{
                    selectedFieldInfo.standardReference
                  }}
                </div>
                <div class="text-caption mt-1">
                  <strong>Considered risky: </strong>{{
                    selectedFieldInfo.riskyThreshold
                  }}
                </div>
              </v-alert>
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="newRule.op"
                :items="operatorOptions"
                item-title="title"
                item-value="value"
                label="Operator"
                density="compact"
                variant="outlined"
                :hint="selectedOperatorInfo?.description || ''"
                persistent-hint
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="newRule.value_a"
                :label="
                  newRule.op === 'in'
                    ? 'Values (pipe-separated, e.g. 3|4)'
                    : newRule.op === 'between'
                      ? 'Value A (lower bound)'
                      : 'Value A'
                "
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col
              v-if="newRule.op === 'between'"
              cols="6"
            >
              <v-text-field
                v-model="newRule.value_b"
                label="Value B (upper bound)"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="newRule.outcome"
                :items="outcomeOptions"
                label="Outcome"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="newRule.loading_percent"
                label="Loading %"
                type="number"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="newRule.exclusion_code"
                label="Exclusion code"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model.number="newRule.priority"
                label="Priority"
                type="number"
                density="compact"
                variant="outlined"
              />
            </v-col>
            <v-col cols="12">
              <v-textarea
                v-model="newRule.notes"
                label="Notes"
                rows="2"
                density="compact"
                variant="outlined"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="newRuleDialog = false">Cancel</v-btn>
          <v-btn color="primary" @click="submitRule">{{
            editingRuleId ? 'Save changes' : 'Add'
          }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- CSV import -->
    <v-dialog v-model="importDialog" max-width="520">
      <v-card>
        <v-card-title>Import rules from CSV</v-card-title>
        <v-card-text>
          <p class="text-caption mb-2"
            >Header columns: rule_set_name, version, category, field, op,
            value_a, value_b, values, outcome, loading_percent, exclusion_code,
            priority, notes. `values` is pipe-separated for the `in`
            operator.</p
          >
          <v-file-input
            v-model="importFile"
            label="CSV file"
            density="compact"
            variant="outlined"
            accept=".csv"
            show-size
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="importDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="importBusy" @click="importCsv"
            >Import</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface UWRule {
  id: number
  rule_set_id: number
  category: string
  field: string
  op: string
  condition_json: string
  outcome: string
  loading_percent: number
  exclusion_code: string
  priority: number
  notes: string
}
interface UWRuleSet {
  id: number
  name: string
  version: number
  active: boolean
  creation_date: string
  rules?: UWRule[]
}

const ruleSets = ref<UWRuleSet[]>([])
const selectedSetId = ref<number[]>([])
const selectedSet = ref<UWRuleSet | null>(null)

const newRuleSetDialog = ref(false)
const newRuleDialog = ref(false)
const editingRuleId = ref<number | null>(null)
const importDialog = ref(false)
const importBusy = ref(false)
const dryRunBusy = ref(false)
const seedingStarter = ref(false)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const newRuleSet = ref({ name: '', version: 1 })
const newRule = ref({
  category: 'medical_condition',
  field: '',
  op: 'eq',
  value_a: '',
  value_b: '',
  outcome: 'refer',
  loading_percent: 0,
  exclusion_code: '',
  priority: 10,
  notes: ''
})
const importFile = ref<File[]>([])

const dryRun = ref<Record<string, any>>({
  bmi: undefined,
  age: undefined,
  occupation_class: undefined,
  smoker: false
})
const dryRunResult = ref<any>(null)

// Category dictionary — surfaces a description plus the "standard
// rates assume / considered risky" baseline so rule authors understand
// what a rule in each category is for.
interface CategoryInfo {
  title: string
  value: string
  description: string
  standardAssumption: string
}
const categoryOptions: CategoryInfo[] = [
  {
    title: 'Build',
    value: 'build',
    description:
      'Height-and-weight-derived measures: BMI, height, weight. Used to identify obesity-related mortality risk.',
    standardAssumption:
      'Standard rates assume BMI ≤ 30 (normal / overweight). BMI > 30 progressively increases mortality risk.'
  },
  {
    title: 'Lifestyle',
    value: 'lifestyle',
    description:
      'Self-reported lifestyle factors: smoking, alcohol use, hazardous hobbies.',
    standardAssumption:
      'Standard rates assume non-smoker, moderate alcohol use (< 14 units/week), no hazardous hobbies. Any of these flagged is typically loaded.'
  },
  {
    title: 'Occupation',
    value: 'occupation',
    description:
      'Occupation-class lookup. Class 1 = office / professional, Class 4 = high-risk manual / hazardous work.',
    standardAssumption:
      'Standard rates assume occupation class 1 or 2. Classes 3 and 4 carry higher injury / mortality risk.'
  },
  {
    title: 'Medical condition',
    value: 'medical_condition',
    description:
      'Disclosed medical conditions. Each disclosed condition becomes a boolean field on the case (e.g. condition_diabetes_type2 = true).',
    standardAssumption:
      'Standard rates assume no disclosed conditions. Any disclosed condition typically refers to a human underwriter for assessment.'
  },
  {
    title: 'Takeover',
    value: 'takeover',
    description:
      'Continuation-of-cover rules for prior-insurer schedules. References prior loadings, prior exclusions, in-force status.',
    standardAssumption:
      'Standard takeover assumption: members in-force with no prior loadings continue with no further evidence; members with prior loadings continue at the same terms; members not in-force or new to the scheme require fresh evidence.'
  }
]

// Field catalog per category. Each field includes a standard reference
// point (what value the standard rate assumes) and a risky threshold
// (where loadings / declines typically kick in) so the rule author has
// guidance without leaving the dialog. Authors can still type a custom
// field — the combobox falls back to free text.
interface FieldInfo {
  title: string
  value: string
  standardReference: string
  riskyThreshold: string
}
const fieldCatalog: Record<string, FieldInfo[]> = {
  build: [
    {
      title: 'bmi — Body Mass Index',
      value: 'bmi',
      standardReference: 'Standard rates assume BMI 18.5–30.',
      riskyThreshold:
        'BMI 30–35: typically +25%. 35–40: +50%. > 40: decline.'
    },
    {
      title: 'height_cm — Height in centimetres',
      value: 'height_cm',
      standardReference: 'Used to compute BMI; rarely a rule field on its own.',
      riskyThreshold: 'No standalone threshold — combine with weight via BMI.'
    },
    {
      title: 'weight_kg — Weight in kilograms',
      value: 'weight_kg',
      standardReference: 'Used to compute BMI; rarely a rule field on its own.',
      riskyThreshold: 'No standalone threshold — combine with height via BMI.'
    }
  ],
  lifestyle: [
    {
      title: 'smoker — Smoker (true/false)',
      value: 'smoker',
      standardReference: 'Standard rates assume non-smoker.',
      riskyThreshold: 'Any disclosed smoker: typical loading +15% to +50%.'
    },
    {
      title: 'cigarettes_per_day — Daily cigarette count',
      value: 'cigarettes_per_day',
      standardReference: 'Standard rates assume 0.',
      riskyThreshold:
        '1–10: low band. 11–20: medium. > 20: heavy smoker — high loading.'
    },
    {
      title: 'alcohol_units_per_week — Weekly alcohol intake',
      value: 'alcohol_units_per_week',
      standardReference: 'Standard rates assume < 14 units / week.',
      riskyThreshold:
        '14–28: moderate excess. > 28: refer for further evidence.'
    },
    {
      title: 'hazardous_hobbies — Has hazardous hobbies (true/false)',
      value: 'hazardous_hobbies',
      standardReference: 'Standard rates assume no hazardous hobbies.',
      riskyThreshold:
        'Scuba, climbing, motorsports, aviation: typically refer for specific exclusion.'
    }
  ],
  occupation: [
    {
      title: 'occupation_class — Class 1 (low risk) to 4 (high risk)',
      value: 'occupation_class',
      standardReference:
        'Standard rates assume class 1 (office / professional) or 2 (light manual).',
      riskyThreshold:
        'Class 3 (heavy manual): typical loading. Class 4 (hazardous, e.g. mining, offshore): high loading or decline.'
    }
  ],
  medical_condition: [
    {
      title: 'condition_diabetes_type1 — Disclosed Type 1 diabetes',
      value: 'condition_diabetes_type1',
      standardReference: 'Standard rates assume no disclosed diabetes.',
      riskyThreshold: 'Disclosed: refer to a human underwriter.'
    },
    {
      title: 'condition_diabetes_type2 — Disclosed Type 2 diabetes',
      value: 'condition_diabetes_type2',
      standardReference: 'Standard rates assume no disclosed diabetes.',
      riskyThreshold:
        'Disclosed: typical loading depends on control / HbA1c — refer.'
    },
    {
      title: 'condition_hypertension — Disclosed high blood pressure',
      value: 'condition_hypertension',
      standardReference: 'Standard rates assume no disclosed hypertension.',
      riskyThreshold:
        'Disclosed: controlled hypertension often accepts; uncontrolled refers.'
    },
    {
      title: 'condition_asthma — Disclosed asthma',
      value: 'condition_asthma',
      standardReference: 'Standard rates assume no asthma.',
      riskyThreshold: 'Mild controlled: often accept. Severe / hospitalised: refer.'
    },
    {
      title: 'condition_cancer_history — Disclosed cancer history',
      value: 'condition_cancer_history',
      standardReference: 'Standard rates assume no cancer history.',
      riskyThreshold:
        'Disclosed: refer; outcome depends on type, stage, time since treatment.'
    },
    {
      title: 'condition_cardiovascular_disease — Disclosed cardiovascular disease',
      value: 'condition_cardiovascular_disease',
      standardReference: 'Standard rates assume no cardiovascular history.',
      riskyThreshold: 'Disclosed: refer or decline.'
    },
    {
      title: 'condition_kidney_disease — Disclosed kidney disease',
      value: 'condition_kidney_disease',
      standardReference: 'Standard rates assume no kidney disease.',
      riskyThreshold: 'Disclosed: refer; severity-dependent outcome.'
    },
    {
      title: 'condition_mental_health_treatment — Disclosed mental health treatment',
      value: 'condition_mental_health_treatment',
      standardReference: 'Standard rates assume no current mental health treatment.',
      riskyThreshold:
        'Disclosed: refer; outcome depends on diagnosis and current treatment.'
    }
  ],
  takeover: [
    {
      title: 'in_force — Member in-force with prior insurer (true/false)',
      value: 'in_force',
      standardReference:
        'Standard takeover assumes prior in-force lives continue.',
      riskyThreshold:
        'Not in-force = new evidence required (treat as new business).'
    },
    {
      title: 'prior_loading_gla — Prior GLA loading %',
      value: 'prior_loading_gla',
      standardReference: 'Standard takeover assumes no prior loading.',
      riskyThreshold: '> 0: continue with the same loading.'
    },
    {
      title: 'prior_loading_ptd — Prior PTD loading %',
      value: 'prior_loading_ptd',
      standardReference: 'Standard takeover assumes no prior loading.',
      riskyThreshold: '> 0: continue with the same loading.'
    },
    {
      title: 'prior_loading_ci — Prior CI loading %',
      value: 'prior_loading_ci',
      standardReference: 'Standard takeover assumes no prior loading.',
      riskyThreshold: '> 0: continue with the same loading.'
    }
  ]
}

// Operator dictionary — descriptive labels + one-line semantics so the
// `eq` / `gte` shorthand is no longer a guessing game for the author.
interface OperatorInfo {
  title: string
  value: string
  description: string
}
const operatorOptions: OperatorInfo[] = [
  {
    title: 'Equals',
    value: 'eq',
    description: 'Match when the field equals Value A exactly.'
  },
  {
    title: 'Not equal to',
    value: 'ne',
    description: 'Match when the field does NOT equal Value A.'
  },
  {
    title: 'Greater than',
    value: 'gt',
    description: 'Match when the field is strictly greater than Value A.'
  },
  {
    title: 'Greater than or equal to',
    value: 'gte',
    description: 'Match when the field is ≥ Value A.'
  },
  {
    title: 'Less than',
    value: 'lt',
    description: 'Match when the field is strictly less than Value A.'
  },
  {
    title: 'Less than or equal to',
    value: 'lte',
    description: 'Match when the field is ≤ Value A.'
  },
  {
    title: 'Between (inclusive)',
    value: 'between',
    description:
      'Match when the field is between Value A and Value B inclusive (Value A ≤ field ≤ Value B).'
  },
  {
    title: 'Is one of',
    value: 'in',
    description:
      'Match when the field is any of the listed values. Use the pipe (|) separator, e.g. 3|4.'
  }
]
const outcomeOptions = ['accept', 'refer', 'decline']

// Currently-selected category info (for help text in the dialog).
const selectedCategoryInfo = computed<CategoryInfo | null>(
  () =>
    categoryOptions.find((c) => c.value === newRule.value.category) || null
)

// Fields available for the currently-selected category. Falls back to
// an empty list if the category isn't catalogued (the combobox still
// accepts free text in that case).
const fieldsForCategory = computed<FieldInfo[]>(
  () => fieldCatalog[newRule.value.category] || []
)

// Help block for the chosen field — appears below the field combobox.
const selectedFieldInfo = computed<FieldInfo | null>(
  () =>
    fieldsForCategory.value.find((f) => f.value === newRule.value.field) ||
    null
)

// Help line for the chosen operator — appears below the operator select.
const selectedOperatorInfo = computed<OperatorInfo | null>(
  () => operatorOptions.find((o) => o.value === newRule.value.op) || null
)

// When the category changes, clear the field if the previously-chosen
// field doesn't appear in the new category's catalog. Prevents stale
// pairs like (category=Lifestyle, field=bmi) lingering.
const onCategoryChange = () => {
  const fields = fieldCatalog[newRule.value.category] || []
  if (
    newRule.value.field &&
    !fields.some((f) => f.value === newRule.value.field)
  ) {
    newRule.value.field = ''
  }
}

const outcomeColor = (o: string) => {
  if (o === 'decline') return 'error'
  if (o === 'refer') return 'warning'
  return 'success'
}
const formatDate = (s: string) => (s ? new Date(s).toLocaleString() : '—')

const flash = (msg: string, color = 'success') => {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

const loadRuleSets = async () => {
  try {
    const res = await GroupPricingService.listUWRuleSets()
    ruleSets.value = res.data || []
    if (
      ruleSets.value.length &&
      (!selectedSet.value || !ruleSets.value.find((r) => r.id === selectedSet.value!.id))
    ) {
      selectRuleSet(ruleSets.value[0].id)
    }
  } catch (err) {
    flash('Failed to load rule sets', 'error')
  }
}

const selectRuleSet = async (id: number) => {
  selectedSetId.value = [id]
  try {
    const res = await GroupPricingService.getUWRuleSet(id)
    selectedSet.value = res.data
  } catch (err) {
    flash('Failed to load rule set', 'error')
  }
}

const createRuleSet = async () => {
  try {
    const res = await GroupPricingService.createUWRuleSet({
      name: newRuleSet.value.name,
      version: newRuleSet.value.version
    })
    newRuleSetDialog.value = false
    newRuleSet.value = { name: '', version: 1 }
    await loadRuleSets()
    if (res.data?.id) selectRuleSet(res.data.id)
    flash('Rule set created')
  } catch (err: any) {
    flash(err?.response?.data || 'Failed', 'error')
  }
}

const activate = async (id: number) => {
  try {
    await GroupPricingService.activateUWRuleSet(id)
    await loadRuleSets()
    flash('Rule set activated')
  } catch (err: any) {
    flash(err?.response?.data || 'Failed', 'error')
  }
}

const buildConditionJSON = (rule: any): string => {
  if (rule.op === 'in') {
    const parts = String(rule.value_a || '')
      .split('|')
      .map((p: string) => p.trim())
      .filter(Boolean)
      .map(scalar)
    return JSON.stringify({ values: parts })
  }
  if (rule.op === 'between') {
    return JSON.stringify({
      value_a: scalar(rule.value_a),
      value_b: scalar(rule.value_b)
    })
  }
  return JSON.stringify({ value_a: scalar(rule.value_a) })
}
const scalar = (s: any) => {
  if (s === 'true' || s === true) return true
  if (s === 'false' || s === false) return false
  const n = Number(s)
  if (s !== '' && s != null && !isNaN(n)) return n
  return s
}

// Reset the rule-form state to defaults. Used by both "Add" and "Cancel"
// so re-opening the dialog never inherits a previous edit's values.
const resetRuleForm = () => {
  newRule.value = {
    category: 'medical_condition',
    field: '',
    op: 'eq',
    value_a: '',
    value_b: '',
    outcome: 'refer',
    loading_percent: 0,
    exclusion_code: '',
    priority: 10,
    notes: ''
  }
  editingRuleId.value = null
}

const openNewRuleDialog = () => {
  resetRuleForm()
  newRuleDialog.value = true
}

// Decompose ConditionJSON back into the dialog's value_a / value_b /
// "pipe-separated values" fields so an edit pre-populates correctly.
const populateRuleFormFromCondition = (rule: any) => {
  try {
    const cond = JSON.parse(rule.condition_json || '{}')
    if (rule.op === 'in' && Array.isArray(cond.values)) {
      newRule.value.value_a = cond.values.join('|')
    } else {
      if (cond.value_a !== undefined && cond.value_a !== null) {
        newRule.value.value_a = String(cond.value_a)
      }
      if (cond.value_b !== undefined && cond.value_b !== null) {
        newRule.value.value_b = String(cond.value_b)
      }
    }
  } catch {
    // Malformed condition — leave the fields blank and let the admin
    // re-enter values.
  }
}

const openEditRuleDialog = (rule: any) => {
  editingRuleId.value = rule.id
  newRule.value = {
    category: rule.category,
    field: rule.field,
    op: rule.op,
    value_a: '',
    value_b: '',
    outcome: rule.outcome,
    loading_percent: rule.loading_percent || 0,
    exclusion_code: rule.exclusion_code || '',
    priority: rule.priority || 0,
    notes: rule.notes || ''
  }
  populateRuleFormFromCondition(rule)
  newRuleDialog.value = true
}

const submitRule = async () => {
  if (!selectedSet.value) return
  const payload = {
    category: newRule.value.category,
    field: newRule.value.field,
    op: newRule.value.op,
    condition_json: buildConditionJSON(newRule.value),
    outcome: newRule.value.outcome,
    loading_percent: newRule.value.loading_percent,
    exclusion_code: newRule.value.exclusion_code,
    priority: newRule.value.priority,
    notes: newRule.value.notes
  }
  try {
    if (editingRuleId.value) {
      await GroupPricingService.updateUWRule(editingRuleId.value, payload)
      flash('Rule updated')
    } else {
      await GroupPricingService.createUWRule(selectedSet.value.id, payload)
      flash('Rule added')
    }
    newRuleDialog.value = false
    resetRuleForm()
    await selectRuleSet(selectedSet.value.id)
  } catch (err: any) {
    flash(err?.response?.data || 'Failed', 'error')
  }
}

const deleteRule = async (rule: any) => {
  if (!selectedSet.value) return
  const op = rule.op || '?'
  const field = rule.field || '?'
  const category = rule.category || '?'
  const outcome = rule.outcome || '?'
  const loading = rule.loading_percent ? ` (+${rule.loading_percent}%)` : ''
  const summary = `${category} · ${field} ${op} → ${outcome}${loading}`
  if (
    !window.confirm(
      `Delete this rule?\n\n${summary}\n\nThis cannot be undone.`,
    )
  ) {
    return
  }
  try {
    await GroupPricingService.deleteUWRule(rule.id)
    await selectRuleSet(selectedSet.value.id)
    flash('Rule deleted')
  } catch (err: any) {
    flash(err?.response?.data || 'Failed', 'error')
  }
}

const importCsv = async () => {
  if (!importFile.value?.length) return
  importBusy.value = true
  try {
    const fd = new FormData()
    const file = Array.isArray(importFile.value)
      ? importFile.value[0]
      : importFile.value
    fd.append('file', file)
    const res = await GroupPricingService.importUWRulesCSV(fd)
    flash(`Imported ${res.data?.imported || 0} rules`)
    importDialog.value = false
    importFile.value = []
    await loadRuleSets()
  } catch (err: any) {
    flash(err?.response?.data?.error || 'Import failed', 'error')
  } finally {
    importBusy.value = false
  }
}

const runDryRun = async () => {
  if (!selectedSet.value) return
  dryRunBusy.value = true
  try {
    const ctx: Record<string, any> = {}
    for (const [k, v] of Object.entries(dryRun.value)) {
      if (v !== undefined && v !== null && v !== '') ctx[k] = v
    }
    const res = await GroupPricingService.dryRunUWRules({
      context: ctx,
      rule_set_id: selectedSet.value.id
    })
    dryRunResult.value = res.data
  } catch (err: any) {
    flash(err?.response?.data || 'Dry-run failed', 'error')
  } finally {
    dryRunBusy.value = false
  }
}

const seedStarter = async () => {
  seedingStarter.value = true
  try {
    const res = await GroupPricingService.seedStarterUWRuleSet()
    await loadRuleSets()
    if (res.data?.id) {
      await selectRuleSet(res.data.id)
    }
    flash('Starter rule set loaded')
  } catch (err: any) {
    flash(err?.response?.data || 'Failed to seed starter', 'error')
  } finally {
    seedingStarter.value = false
  }
}

const duplicate = async (id: number) => {
  try {
    const res = await GroupPricingService.duplicateUWRuleSet(id)
    await loadRuleSets()
    if (res.data?.id) {
      await selectRuleSet(res.data.id)
    }
    flash('Rule set duplicated')
  } catch (err: any) {
    flash(err?.response?.data || 'Failed to duplicate', 'error')
  }
}

const confirmDeleteRuleSet = async (rs: UWRuleSet) => {
  // Lightweight confirmation — the server enforces the safety rules
  // (refuses when active or referenced), so we trust the backend's NO
  // and just nudge the user.
  if (
    !window.confirm(
      `Delete rule set "${rs.name} v${rs.version}"?\n\n` +
        'This is allowed only if the set is inactive AND no underwriting cases reference it. ' +
        'The server will refuse otherwise.'
    )
  ) {
    return
  }
  try {
    await GroupPricingService.deleteUWRuleSet(rs.id)
    if (selectedSet.value?.id === rs.id) {
      selectedSet.value = null
      selectedSetId.value = []
    }
    await loadRuleSets()
    flash('Rule set deleted')
  } catch (err: any) {
    const msg =
      err?.response?.data?.error ||
      err?.response?.data ||
      'Failed to delete rule set'
    flash(msg, 'error')
  }
}

const exportCsv = async () => {
  if (!selectedSet.value) return
  try {
    const res = await GroupPricingService.exportUWRuleSetCSV(
      selectedSet.value.id
    )
    const blob = new Blob([res.data], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const filename = `${selectedSet.value.name.replace(/\s+/g, '_')}_v${selectedSet.value.version}.csv`
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  } catch (err: any) {
    flash(err?.response?.data || 'Export failed', 'error')
  }
}

onMounted(loadRuleSets)
</script>
