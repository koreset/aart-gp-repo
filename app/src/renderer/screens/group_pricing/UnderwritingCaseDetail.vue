<template>
  <v-container v-if="caseFile" fluid>
    <!-- Hero header — at-a-glance case identity, status and ratio. -->
    <v-card flat rounded="lg" class="case-hero mb-4">
      <v-card-text class="pa-5">
        <v-row align="center" no-gutters>
          <v-col cols="auto" class="pr-3">
            <v-btn
              size="small"
              variant="text"
              prepend-icon="mdi-arrow-left"
              @click="router.go(-1)"
              >Back</v-btn
            >
          </v-col>
          <v-col>
            <div class="text-overline text-grey-darken-1 mb-1">
              <v-icon size="small" class="mr-1">mdi-pound</v-icon>
              Case {{ caseFile.id }} ·
              <v-icon size="small" class="mx-1"
                >mdi-file-document-outline</v-icon
              >
              Quote {{ caseFile.quote_id }} ·
              {{ caseFile.category }}
            </div>
            <div class="text-h5 font-weight-bold text-grey-darken-4">
              {{ caseFile.member_name }}
            </div>
            <div class="text-caption text-grey-darken-1 mt-1">
              ID
              <span class="font-weight-medium">{{
                caseFile.member_id_number || '—'
              }}</span>
            </div>
          </v-col>
          <v-col cols="12" md="auto" class="case-hero-meta">
            <div class="d-flex flex-wrap justify-md-end gap-2 mb-3">
              <v-chip
                :color="tierColor(caseFile.tier)"
                variant="flat"
                size="small"
                label
                prepend-icon="mdi-shield-alert-outline"
                >Tier {{ caseFile.tier }} —
                {{ tierLabel(caseFile.tier) }}</v-chip
              >
              <v-chip
                :color="statusColor(caseFile.status)"
                variant="flat"
                size="small"
                label
                :prepend-icon="statusIcon(caseFile.status)"
                >{{ statusLabel(caseFile.status) }}</v-chip
              >
              <v-chip
                v-if="caseFile.engine_outcome"
                color="indigo"
                variant="tonal"
                size="small"
                label
                prepend-icon="mdi-robot-outline"
                >Engine: {{ caseFile.engine_outcome
                }}<span v-if="caseFile.engine_loading">
                  · +{{ caseFile.engine_loading }}%</span
                ></v-chip
              >
            </div>
            <div class="case-hero-ratio">
              <div class="text-overline text-grey-darken-1">SA ÷ FCL</div>
              <div
                class="text-h4 font-weight-bold"
                :class="`text-${tierColor(caseFile.tier)}`"
                >{{
                  caseFile.fcl_excess_ratio
                    ? caseFile.fcl_excess_ratio.toFixed(2) + '×'
                    : '—'
                }}</div
              >
            </div>
          </v-col>
        </v-row>

        <v-divider class="my-4" />

        <!-- Workflow progress — simple pill stepper tied to case.status. -->
        <div class="d-flex align-center flex-wrap workflow-pills">
          <div
            v-for="(step, i) in workflowSteps"
            :key="step.key"
            class="d-flex align-center"
          >
            <v-chip
              :variant="step.state === 'pending' ? 'outlined' : 'flat'"
              :color="step.color"
              :prepend-icon="step.icon"
              size="small"
              label
              class="workflow-pill"
              >{{ step.label }}</v-chip
            >
            <v-icon
              v-if="i < workflowSteps.length - 1"
              size="small"
              class="mx-2 text-grey"
              >mdi-chevron-right</v-icon
            >
          </div>
        </div>
      </v-card-text>
    </v-card>

    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #default>
            <v-row>
              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary"
                      >mdi-account-circle-outline</v-icon
                    >
                    Member
                  </v-card-title>
                  <v-card-text>
                    <div class="key-value-row">
                      <span class="key-label">Full name</span>
                      <span class="key-value">{{ caseFile.member_name }}</span>
                    </div>
                    <div class="key-value-row">
                      <span class="key-label">ID number</span>
                      <span class="key-value">{{
                        caseFile.member_id_number || '—'
                      }}</span>
                    </div>
                    <div class="key-value-row">
                      <span class="key-label">Scheme category</span>
                      <span class="key-value">{{ caseFile.category }}</span>
                    </div>
                    <div class="key-value-row">
                      <span class="key-label">Quote</span>
                      <span class="key-value">#{{ caseFile.quote_id }}</span>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary"
                      >mdi-shield-half-full</v-icon
                    >
                    Cover vs FCL
                  </v-card-title>
                  <v-card-text>
                    <div class="key-value-row">
                      <span class="key-label">Free Cover Limit</span>
                      <span class="key-value font-weight-medium">{{
                        format(caseFile.free_cover_limit)
                      }}</span>
                    </div>
                    <v-divider class="my-2" />
                    <div class="key-value-row">
                      <span class="key-label">GLA</span>
                      <span
                        class="key-value"
                        :class="
                          caseFile.gla_sum_assured > caseFile.free_cover_limit
                            ? 'text-warning font-weight-medium'
                            : ''
                        "
                        >{{ format(caseFile.gla_sum_assured) }}</span
                      >
                    </div>
                    <div class="key-value-row">
                      <span class="key-label">PTD</span>
                      <span class="key-value">{{
                        format(caseFile.ptd_sum_assured)
                      }}</span>
                    </div>
                    <div class="key-value-row">
                      <span class="key-label">CI</span>
                      <span class="key-value">{{
                        format(caseFile.ci_sum_assured)
                      }}</span>
                    </div>
                    <div class="key-value-row">
                      <span class="key-label">Spouse GLA</span>
                      <span class="key-value">{{
                        format(caseFile.spouse_gla_sum_assured)
                      }}</span>
                    </div>
                    <v-divider class="my-2" />
                    <div
                      class="ratio-callout"
                      :class="`ratio-${tierColor(caseFile.tier)}`"
                    >
                      <div class="text-caption text-grey-darken-1"
                        >Highest SA as a multiple of FCL</div
                      >
                      <div
                        class="text-h5 font-weight-bold"
                        :class="`text-${tierColor(caseFile.tier)}`"
                        >{{
                          caseFile.fcl_excess_ratio
                            ? caseFile.fcl_excess_ratio.toFixed(2) + '×'
                            : '—'
                        }}</div
                      >
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary"
                      >mdi-briefcase-clock-outline</v-icon
                    >
                    Workflow
                  </v-card-title>
                  <v-card-text>
                    <v-autocomplete
                      v-model="assigneeInput"
                      :items="underwriterOptions"
                      item-title="display"
                      item-value="email"
                      label="Assigned underwriter"
                      density="compact"
                      variant="outlined"
                      :loading="busy || loadingUnderwriters"
                      clearable
                      hide-no-data
                      append-inner-icon="mdi-content-save"
                      :placeholder="
                        underwriterOptions.length
                          ? 'Choose an underwriter'
                          : 'Loading users…'
                      "
                      @click:append-inner="assignCase"
                    />
                    <v-divider class="my-3" />
                    <div class="d-flex flex-wrap gap-2">
                      <v-btn
                        v-for="t in allowedTransitions"
                        :key="t.value"
                        size="small"
                        :color="t.color"
                        variant="tonal"
                        :loading="busy"
                        @click="transition(t.value)"
                        >{{ t.label }}</v-btn
                      >
                    </div>
                    <p v-if="!allowedTransitions.length" class="text-grey mt-2"
                      >Terminal state — no further transitions.</p
                    >
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary"
                      >mdi-robot-outline</v-icon
                    >
                    <span>Rules engine recommendation</span>
                    <v-spacer />
                    <v-btn
                      size="x-small"
                      variant="text"
                      :loading="engineBusy"
                      @click="runEngine"
                      >Re-evaluate</v-btn
                    >
                  </v-card-title>
                  <v-card-text>
                    <p v-if="!engineResult" class="text-grey"
                      >No evaluation yet. Add hypothetical disclosure values
                      below and re-evaluate.</p
                    >
                    <div v-else>
                      <v-chip
                        :color="outcomeColor(engineResult.outcome)"
                        variant="tonal"
                        size="small"
                        class="mr-2"
                        >Engine: {{ engineResult.outcome }}</v-chip
                      >
                      <v-chip
                        v-if="engineResult.max_loading"
                        color="warning"
                        variant="tonal"
                        size="small"
                        class="mr-2"
                        >+{{ engineResult.max_loading }}% loading</v-chip
                      >
                      <v-chip
                        v-for="ex in engineResult.exclusions || []"
                        :key="ex"
                        color="grey"
                        variant="tonal"
                        size="small"
                        class="mr-2 mt-1"
                        >{{ ex }}</v-chip
                      >
                      <p
                        v-if="engineResult.rule_set_id"
                        class="text-caption text-grey mt-2"
                        >Rule set #{{ engineResult.rule_set_id }} v{{
                          engineResult.rule_set_version
                        }}.</p
                      >
                    </div>
                    <v-divider class="my-3" />
                    <v-row dense>
                      <v-col cols="6">
                        <v-text-field
                          v-model.number="engineCtx.bmi"
                          label="BMI"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-text-field
                          v-model.number="engineCtx.age"
                          label="Age"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-text-field
                          v-model.number="engineCtx.occupation_class"
                          label="Occupation class"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-checkbox
                          v-model="engineCtx.smoker"
                          label="Smoker"
                          density="compact"
                          hide-details
                        />
                      </v-col>
                    </v-row>
                    <v-btn
                      size="small"
                      color="primary"
                      variant="tonal"
                      :disabled="!engineResult"
                      @click="prefillDecisionFromEngine"
                      >Use as decision</v-btn
                    >
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary">mdi-gavel</v-icon>
                    Record decision
                  </v-card-title>
                  <v-card-text>
                    <v-row dense>
                      <v-col cols="6">
                        <v-select
                          v-model="decision.benefit_type"
                          :items="benefitOptions"
                          label="Benefit"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-select
                          v-model="decision.outcome"
                          :items="outcomeOptions"
                          label="Outcome"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-text-field
                          v-model.number="decision.loading_percent"
                          label="Loading %"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-text-field
                          v-model.number="decision.cover_cap"
                          label="Cover cap"
                          type="number"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="12">
                        <v-text-field
                          v-model="decision.exclusion_code"
                          label="Exclusion code"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                      <v-col cols="12">
                        <v-textarea
                          v-model="decision.notes"
                          label="Notes"
                          rows="2"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                    </v-row>
                    <v-btn
                      color="primary"
                      :loading="busy"
                      :disabled="!decision.benefit_type || !decision.outcome"
                      @click="submitDecision"
                      >Save decision</v-btn
                    >
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary"
                      >mdi-clipboard-text-clock-outline</v-icon
                    >
                    Decisions on this case
                  </v-card-title>
                  <v-card-text>
                    <v-list density="compact">
                      <v-list-item
                        v-for="d in caseFile.decisions || []"
                        :key="d.id"
                        :title="`${d.benefit_type.toUpperCase()} — ${d.outcome}`"
                        :subtitle="
                          [
                            d.loading_percent
                              ? `loading ${d.loading_percent}%`
                              : '',
                            d.exclusion_code ? `exc ${d.exclusion_code}` : '',
                            d.cover_cap ? `cap ${format(d.cover_cap)}` : ''
                          ]
                            .filter(Boolean)
                            .join(' · ')
                        "
                      >
                        <template #append>
                          <span class="text-caption text-grey">{{
                            formatDate(d.creation_date)
                          }}</span>
                        </template>
                      </v-list-item>
                    </v-list>
                    <p
                      v-if="
                        !caseFile.decisions || caseFile.decisions.length === 0
                      "
                      class="text-grey"
                      >No decisions recorded yet.</p
                    >
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="6">
                <v-card variant="outlined" rounded="lg" class="h-100">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary">mdi-paperclip</v-icon>
                    Attachments
                  </v-card-title>
                  <v-card-text>
                    <v-row dense>
                      <v-col cols="7">
                        <v-file-input
                          v-model="uploadFile"
                          label="File"
                          density="compact"
                          variant="outlined"
                          show-size
                          chips
                        />
                      </v-col>
                      <v-col cols="5">
                        <v-select
                          v-model="uploadKind"
                          :items="attachmentKindOptions"
                          label="Kind"
                          density="compact"
                          variant="outlined"
                        />
                      </v-col>
                    </v-row>
                    <v-btn
                      color="primary"
                      :loading="busy"
                      :disabled="!uploadFile"
                      @click="upload"
                      >Upload</v-btn
                    >
                    <v-divider class="my-3" />
                    <v-list density="compact">
                      <v-list-item
                        v-for="a in caseFile.attachments || []"
                        :key="a.id"
                        :title="a.filename"
                        :subtitle="`${a.kind} · ${formatBytes(a.size_bytes)}`"
                      >
                        <template #append>
                          <v-btn
                            size="x-small"
                            variant="text"
                            :href="a.viewer_url"
                            target="_blank"
                            prepend-icon="mdi-download"
                            >Open</v-btn
                          >
                        </template>
                      </v-list-item>
                    </v-list>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="6">
                <DisclosureForm
                  :case-id="caseFile.id"
                  @submitted="onDisclosureSubmitted"
                />
              </v-col>
              <v-col cols="12" md="6">
                <ActivelyAtWorkAttestation
                  :case-id="caseFile.id"
                  :quote-id="caseFile.quote_id"
                  :member-name="caseFile.member_name"
                  :member-id-number="caseFile.member_id_number"
                  @submitted="load"
                />
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12">
                <v-card variant="outlined" rounded="lg">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary"
                      >mdi-cloud-sync-outline</v-icon
                    >
                    External services
                    <v-spacer />
                    <span class="text-caption text-grey"
                      >Pathology, GP records, e-sign and SMS via vendor
                      adapters.</span
                    >
                  </v-card-title>
                  <v-card-text>
                    <v-row dense>
                      <v-col
                        v-for="action in vendorActions"
                        :key="action.kind"
                        cols="6"
                        md="3"
                      >
                        <v-btn
                          block
                          variant="tonal"
                          :color="action.color"
                          :prepend-icon="action.icon"
                          :loading="vendorBusy === action.kind"
                          @click="openVendorDialog(action.kind)"
                          >{{ action.label }}</v-btn
                        >
                      </v-col>
                    </v-row>

                    <v-divider v-if="vendorRequests.length" class="my-3" />
                    <v-list v-if="vendorRequests.length" density="compact">
                      <v-list-item
                        v-for="r in vendorRequests"
                        :key="r.id"
                        :title="`${vendorKindLabel(r.kind)} — ${r.subject || '—'}`"
                        :subtitle="`${r.provider} · ${formatDate(r.requested_at)} · ${
                          r.cost_cents
                            ? 'R' + (r.cost_cents / 100).toFixed(2)
                            : 'cost N/A'
                        }`"
                      >
                        <template #append>
                          <v-chip
                            :color="vendorStatusColor(r.status)"
                            variant="tonal"
                            size="x-small"
                            class="mr-2"
                            >{{ r.status.replace(/_/g, ' ') }}</v-chip
                          >
                          <v-btn
                            v-if="
                              r.status === 'awaiting_response' &&
                              r.provider === 'mock'
                            "
                            size="x-small"
                            variant="text"
                            color="primary"
                            @click="fireMock(r.id)"
                            >Simulate</v-btn
                          >
                        </template>
                      </v-list-item>
                    </v-list>
                    <p v-else class="text-grey text-caption"
                      >No external requests yet.</p
                    >
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-dialog v-model="vendorDialogOpen" max-width="520">
              <v-card>
                <v-card-title
                  >Request {{ vendorKindLabel(vendorDialogKind) }}</v-card-title
                >
                <v-card-text>
                  <v-text-field
                    v-model="vendorDialog.subject"
                    :label="vendorSubjectLabel"
                    density="compact"
                    variant="outlined"
                    autofocus
                  />
                  <v-textarea
                    v-model="vendorDialog.body"
                    label="Body / instruction"
                    rows="3"
                    density="compact"
                    variant="outlined"
                  />
                </v-card-text>
                <v-card-actions>
                  <v-spacer />
                  <v-btn variant="text" @click="vendorDialogOpen = false"
                    >Cancel</v-btn
                  >
                  <v-btn
                    color="primary"
                    :loading="!!vendorBusy"
                    @click="submitVendor"
                    >Send</v-btn
                  >
                </v-card-actions>
              </v-card>
            </v-dialog>

            <v-row v-if="caseFile.engine_outcome">
              <v-col cols="12">
                <v-alert
                  type="info"
                  variant="tonal"
                  density="compact"
                  icon="mdi-robot-outline"
                >
                  <strong>Latest engine snapshot:</strong>
                  {{ caseFile.engine_outcome }}
                  <span v-if="caseFile.engine_loading"
                    >· loading {{ caseFile.engine_loading }}%</span
                  >
                  <span class="text-caption text-grey ml-2"
                    >Evaluated
                    {{ formatDate(caseFile.engine_evaluated_at || '') }}</span
                  >
                </v-alert>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12">
                <v-card variant="outlined" rounded="lg">
                  <v-card-title class="d-flex align-center font-weight-bold">
                    <v-icon class="mr-2" color="primary">mdi-history</v-icon>
                    Audit timeline
                  </v-card-title>
                  <v-card-text>
                    <v-timeline density="compact" align="start">
                      <v-timeline-item
                        v-for="e in (caseFile.events || []).slice().reverse()"
                        :key="e.id"
                        size="x-small"
                        dot-color="primary"
                      >
                        <div class="text-body-2"
                          ><b>{{ e.event_type.replace(/_/g, ' ') }}</b> ·
                          {{ e.actor || 'system' }}</div
                        >
                        <div class="text-caption text-grey">{{
                          formatDate(e.creation_date)
                        }}</div>
                        <pre
                          v-if="e.payload"
                          class="text-caption mt-1"
                          style="white-space: pre-wrap"
                          >{{ e.payload }}</pre
                        >
                      </v-timeline-item>
                    </v-timeline>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import DisclosureForm from './DisclosureForm.vue'
import ActivelyAtWorkAttestation from './ActivelyAtWorkAttestation.vue'
import { useAppStore } from '@/renderer/store/app'

const route = useRoute()
const router = useRouter()
const caseId = computed(() => String(route.params.caseId))

interface UnderwritingCase {
  id: number
  quote_id: number
  member_name: string
  member_id_number: string
  category: string
  tier: number
  fcl_excess_ratio: number
  gla_sum_assured: number
  ptd_sum_assured: number
  ci_sum_assured: number
  spouse_gla_sum_assured: number
  free_cover_limit: number
  status: string
  assigned_underwriter_email: string
  engine_outcome?: string
  engine_loading?: number
  engine_exclusions?: string
  engine_evaluated_at?: string
  decisions?: Array<{
    id: number
    benefit_type: string
    outcome: string
    loading_percent: number
    exclusion_code: string
    cover_cap: number
    creation_date: string
  }>
  events?: Array<{
    id: number
    event_type: string
    actor: string
    payload: string
    creation_date: string
  }>
  attachments?: Array<{
    id: number
    filename: string
    kind: string
    size_bytes: number
    viewer_url: string
  }>
}

const caseFile = ref<UnderwritingCase | null>(null)
const busy = ref(false)
const assigneeInput = ref<string | null>('')

// Org users sourced from the existing /org-users endpoint so the assignee
// dropdown matches the rest of the app (Reviewer pickers etc.). When the
// org has a "GP role" model in place a future refinement can filter to
// users with underwriting permissions only — for now we surface every
// user and let the assigner pick. The current assignee is added to the
// option list defensively (in case they were removed from the org or
// have a different role tag than expected).
const appStore = useAppStore()
const underwriterOptions = ref<
  Array<{ name: string; email: string; display: string }>
>([])
const loadingUnderwriters = ref(false)
const loadUnderwriters = async () => {
  const orgName = appStore.getOrganisationName
  if (!orgName) return
  loadingUnderwriters.value = true
  try {
    const res = await GroupPricingService.getOrgUsers({ name: orgName })
    if (Array.isArray(res?.data)) {
      const seen = new Set<string>()
      underwriterOptions.value = res.data
        .filter((u: any) => u?.email)
        .filter((u: any) => {
          const k = String(u.email).toLowerCase()
          if (seen.has(k)) return false
          seen.add(k)
          return true
        })
        .map((u: any) => ({
          name: u.name || u.email,
          email: u.email,
          display: u.name ? `${u.name} (${u.email})` : u.email
        }))
    }
  } catch (err) {
    console.warn('Failed to load underwriters', err)
  } finally {
    loadingUnderwriters.value = false
  }
}
// Ensure the current assignee is selectable even if their org row is
// missing or filtered out. Falls back to email-as-display so we never
// orphan a stored assignment.
const ensureCurrentAssigneeOption = (email: string) => {
  if (!email) return
  if (underwriterOptions.value.some((o) => o.email === email)) return
  underwriterOptions.value = [
    {
      name: email,
      email,
      display: email
    },
    ...underwriterOptions.value
  ]
}

const engineBusy = ref(false)
const engineResult = ref<any>(null)
const engineCtx = ref<Record<string, any>>({
  bmi: undefined,
  age: undefined,
  occupation_class: undefined,
  smoker: false
})

const outcomeColor = (o: string) => {
  if (o === 'decline') return 'error'
  if (o === 'refer') return 'warning'
  return 'success'
}

const runEngine = async () => {
  if (!caseFile.value) return
  engineBusy.value = true
  try {
    const overrides: Record<string, any> = {}
    for (const [k, v] of Object.entries(engineCtx.value)) {
      if (v !== undefined && v !== null && v !== '') overrides[k] = v
    }
    const res = await GroupPricingService.dryRunUWRulesForCase(
      caseFile.value.id,
      overrides
    )
    engineResult.value = res.data
  } catch (err: any) {
    flash(err?.response?.data || 'Engine evaluation failed', 'error')
  } finally {
    engineBusy.value = false
  }
}

interface VendorRequest {
  id: number
  kind: string
  provider: string
  subject: string
  status: string
  cost_cents: number
  requested_at: string
}

type VendorKind = 'pathology' | 'gp_records' | 'esign' | 'sms'
const vendorRequests = ref<VendorRequest[]>([])
const vendorBusy = ref<string | null>(null)
const vendorDialogOpen = ref(false)
const vendorDialogKind = ref<VendorKind>('pathology')
const vendorDialog = ref({ subject: '', body: '' })

const vendorActions: Array<{
  kind: VendorKind
  label: string
  icon: string
  color: string
}> = [
  {
    kind: 'pathology',
    label: 'Request pathology',
    icon: 'mdi-test-tube',
    color: 'primary'
  },
  {
    kind: 'gp_records',
    label: 'Request GP records',
    icon: 'mdi-file-document-multiple-outline',
    color: 'info'
  },
  {
    kind: 'esign',
    label: 'Send e-sign link',
    icon: 'mdi-draw',
    color: 'success'
  },
  {
    kind: 'sms',
    label: 'Send SMS',
    icon: 'mdi-message-text-outline',
    color: 'secondary'
  }
]

const vendorSubjectLabel = computed(() => {
  switch (vendorDialogKind.value) {
    case 'pathology':
    case 'gp_records':
      return 'Member ID / reference'
    case 'esign':
      return 'Recipient email'
    case 'sms':
      return 'Recipient phone (+27...)'
    default:
      return 'Subject'
  }
})

const vendorKindLabel = (k: string) => {
  if (k === 'pathology') return 'Pathology'
  if (k === 'gp_records') return 'GP records'
  if (k === 'esign') return 'E-signature'
  if (k === 'sms') return 'SMS'
  return k
}
const vendorStatusColor = (s: string) => {
  if (s === 'complete') return 'success'
  if (s === 'failed') return 'error'
  if (s === 'awaiting_response') return 'warning'
  if (s === 'in_flight') return 'info'
  return 'grey'
}

const openVendorDialog = (kind: VendorKind) => {
  vendorDialogKind.value = kind
  vendorDialog.value = { subject: '', body: defaultBodyFor(kind) }
  vendorDialogOpen.value = true
}
const defaultBodyFor = (kind: VendorKind): string => {
  if (kind === 'pathology') return 'Full blood count, cholesterol panel.'
  if (kind === 'gp_records') return 'Five years of consultation records.'
  if (kind === 'esign') return 'Please sign the attached consent form.'
  if (kind === 'sms')
    return 'Action required: complete your underwriting questionnaire at the link sent to your email.'
  return ''
}

const loadVendorRequests = async () => {
  if (!caseFile.value) return
  try {
    const res = await GroupPricingService.listVendorRequestsForCase(
      caseFile.value.id
    )
    vendorRequests.value = res.data || []
  } catch (err) {
    console.warn('Failed to load vendor requests', err)
  }
}

const submitVendor = async () => {
  if (!caseFile.value) return
  vendorBusy.value = vendorDialogKind.value
  try {
    await GroupPricingService.submitVendorRequest(vendorDialogKind.value, {
      case_id: caseFile.value.id,
      quote_id: caseFile.value.quote_id,
      subject: vendorDialog.value.subject,
      body: vendorDialog.value.body
    })
    vendorDialogOpen.value = false
    await loadVendorRequests()
    flash('Vendor request sent')
  } catch (err: any) {
    flash(err?.response?.data?.error || 'Vendor request failed', 'error')
  } finally {
    vendorBusy.value = null
  }
}

const fireMock = async (requestId: number) => {
  try {
    await GroupPricingService.fireMockVendorWebhook(requestId, {
      status: 'complete',
      filename: 'mock-delivery.txt',
      body_text: 'Mock delivery body',
      attach_kind: 'medical_report'
    })
    await loadVendorRequests()
    await load()
    flash('Mock webhook fired')
  } catch (err: any) {
    flash(err?.response?.data?.error || 'Mock webhook failed', 'error')
  }
}

const onDisclosureSubmitted = async () => {
  await load()
  // Refresh the engine recommendation card if it's been used at least once
  // so the underwriter sees the latest snapshot reflected immediately.
  if (engineResult.value) {
    await runEngine()
  }
}

const prefillDecisionFromEngine = () => {
  if (!engineResult.value) return
  const map: Record<string, string> = {
    accept: 'accept',
    refer: 'postpone',
    decline: 'decline'
  }
  decision.value.outcome = map[engineResult.value.outcome] || 'accept'
  decision.value.loading_percent = engineResult.value.max_loading || 0
  decision.value.exclusion_code =
    (engineResult.value.exclusions || []).join(', ') || ''
}

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const uploadFile = ref<File[]>([])
const uploadKind = ref('medical_report')
const attachmentKindOptions = [
  { title: 'Medical report', value: 'medical_report' },
  { title: 'Actively at work', value: 'actively_at_work' },
  { title: 'Consent', value: 'consent' },
  { title: 'Disclosure', value: 'disclosure' }
]

const decision = ref({
  benefit_type: 'gla',
  outcome: 'accept',
  loading_percent: 0,
  loading_flat_amount: 0,
  exclusion_code: '',
  exclusion_text: '',
  cover_cap: 0,
  notes: ''
})
const benefitOptions = [
  { title: 'GLA', value: 'gla' },
  { title: 'PTD', value: 'ptd' },
  { title: 'CI', value: 'ci' },
  { title: 'Spouse GLA', value: 'sgla' }
]
const outcomeOptions = [
  { title: 'Accept', value: 'accept' },
  { title: 'Postpone', value: 'postpone' },
  { title: 'Decline', value: 'decline' }
]

const transitions: Record<
  string,
  Array<{ label: string; value: string; color: string }>
> = {
  pending_evidence: [
    { label: 'Start review', value: 'in_review', color: 'info' },
    { label: 'Postpone', value: 'postponed', color: 'grey' },
    { label: 'Decline', value: 'declined', color: 'error' }
  ],
  in_review: [
    { label: 'Mark decided', value: 'decided', color: 'success' },
    { label: 'Decline', value: 'declined', color: 'error' },
    { label: 'Postpone', value: 'postponed', color: 'grey' },
    { label: 'Need more evidence', value: 'pending_evidence', color: 'warning' }
  ],
  postponed: [
    { label: 'Resume review', value: 'in_review', color: 'info' },
    { label: 'Request evidence', value: 'pending_evidence', color: 'warning' }
  ],
  decided: [],
  declined: [],
  auto_accepted: []
}
const allowedTransitions = computed(() =>
  caseFile.value ? transitions[caseFile.value.status] || [] : []
)

const tierLabel = (tier: number) => {
  if (tier === 2) return 'Full UW'
  if (tier === 1) return 'Short-form'
  return 'Within FCL'
}
const tierColor = (tier: number) => {
  if (tier === 2) return 'error'
  if (tier === 1) return 'warning'
  return 'success'
}
const statusLabel = (status: string) =>
  status.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
const statusIcon = (status: string) => {
  if (status === 'decided') return 'mdi-check-circle-outline'
  if (status === 'declined') return 'mdi-close-octagon-outline'
  if (status === 'in_review') return 'mdi-eye-outline'
  if (status === 'postponed') return 'mdi-pause-circle-outline'
  if (status === 'auto_accepted') return 'mdi-shield-check-outline'
  return 'mdi-clipboard-clock-outline'
}

// Workflow pill stepper. Visualises the canonical journey
// pending_evidence → in_review → decided (with declined/postponed shown
// as the third pill when terminal). The current step gets the tier
// colour so the header reads as a single coherent status.
const workflowSteps = computed(() => {
  const status = caseFile.value?.status || 'pending_evidence'
  const isPending = status === 'pending_evidence'
  const isReview = status === 'in_review'
  const isDecided = status === 'decided'
  const isDeclined = status === 'declined'
  const isPostponed = status === 'postponed'
  const isAutoAccepted = status === 'auto_accepted'

  // Pill state: 'done' | 'active' | 'pending'.
  const pillFor = (state: 'done' | 'active' | 'pending') => ({
    state,
    color: state === 'pending' ? 'grey' : 'primary'
  })

  let third: any
  if (isDecided) third = { state: 'active', color: 'success' }
  else if (isDeclined) third = { state: 'active', color: 'error' }
  else if (isPostponed) third = { state: 'active', color: 'grey-darken-1' }
  else if (isAutoAccepted) third = { state: 'active', color: 'grey-darken-1' }
  else third = pillFor('pending')

  return [
    {
      key: 'pending_evidence',
      label: 'Pending evidence',
      icon: 'mdi-clipboard-clock-outline',
      ...(isPending ? { state: 'active', color: 'warning' } : pillFor('done'))
    },
    {
      key: 'in_review',
      label: 'In review',
      icon: 'mdi-eye-outline',
      ...(isReview
        ? { state: 'active', color: 'info' }
        : isPending
          ? pillFor('pending')
          : pillFor('done'))
    },
    {
      key: 'terminal',
      label: isDeclined
        ? 'Declined'
        : isPostponed
          ? 'Postponed'
          : isAutoAccepted
            ? 'Auto-accepted'
            : isDecided
              ? 'Decided'
              : 'Decided',
      icon: isDeclined
        ? 'mdi-close-octagon-outline'
        : isPostponed
          ? 'mdi-pause-circle-outline'
          : isAutoAccepted
            ? 'mdi-shield-check-outline'
            : 'mdi-check-circle-outline',
      ...third
    }
  ]
})
const statusColor = (status: string) => {
  if (status === 'decided') return 'success'
  if (status === 'declined') return 'error'
  if (status === 'in_review') return 'info'
  if (status === 'postponed') return 'grey'
  if (status === 'auto_accepted') return 'grey'
  return 'warning'
}
const format = (n: number) =>
  n ? Number(n).toLocaleString(undefined, { maximumFractionDigits: 0 }) : '—'
const formatDate = (s: string) => (s ? new Date(s).toLocaleString() : '—')
const formatBytes = (n: number) => {
  if (!n) return '0 B'
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / 1024 / 1024).toFixed(1)} MB`
}

const flash = (message: string, color = 'success') => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

const load = async () => {
  try {
    const res = await GroupPricingService.getUnderwritingCase(caseId.value)
    caseFile.value = res.data
    assigneeInput.value = res.data?.assigned_underwriter_email || ''
    ensureCurrentAssigneeOption(assigneeInput.value || '')
    await loadVendorRequests()
  } catch (err) {
    console.error(err)
    flash('Failed to load case', 'error')
  }
}

const assignCase = async () => {
  if (!caseFile.value) return
  busy.value = true
  try {
    await GroupPricingService.assignUnderwritingCase(
      caseFile.value.id,
      assigneeInput.value || ''
    )
    flash('Assignee saved')
    await load()
  } catch (err: any) {
    flash(err?.response?.data || 'Failed to assign', 'error')
  } finally {
    busy.value = false
  }
}

const transition = async (status: string) => {
  if (!caseFile.value) return
  busy.value = true
  try {
    await GroupPricingService.transitionUnderwritingCase(
      caseFile.value.id,
      status
    )
    flash('Status updated')
    await load()
  } catch (err: any) {
    flash(err?.response?.data || 'Transition rejected', 'error')
  } finally {
    busy.value = false
  }
}

const submitDecision = async () => {
  if (!caseFile.value) return
  busy.value = true
  try {
    await GroupPricingService.createUnderwritingDecision(
      caseFile.value.id,
      decision.value
    )
    flash('Decision recorded')
    await load()
    decision.value.notes = ''
    decision.value.loading_percent = 0
    decision.value.cover_cap = 0
    decision.value.exclusion_code = ''
  } catch (err: any) {
    flash(err?.response?.data || 'Failed to record decision', 'error')
  } finally {
    busy.value = false
  }
}

const upload = async () => {
  if (!caseFile.value || !uploadFile.value?.length) return
  busy.value = true
  try {
    const formData = new FormData()
    const files = Array.isArray(uploadFile.value)
      ? uploadFile.value
      : [uploadFile.value]
    for (const f of files) {
      formData.append('file', f)
      formData.append('kind', uploadKind.value)
    }
    await GroupPricingService.uploadUnderwritingCaseAttachments(
      caseFile.value.id,
      formData
    )
    flash('Uploaded')
    uploadFile.value = []
    await load()
  } catch (err: any) {
    flash(err?.response?.data || 'Upload failed', 'error')
  } finally {
    busy.value = false
  }
}

onMounted(() => {
  loadUnderwriters()
  load()
})
</script>

<style scoped>
/* Hero banner — subtle elevation + tinted top border so it reads as a
   page-level header rather than just another card. */
.case-hero {
  background: linear-gradient(135deg, #f7f9fc 0%, #ffffff 100%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-top: 3px solid rgb(var(--v-theme-primary));
}

.case-hero-meta {
  min-width: 220px;
}

.case-hero-ratio {
  text-align: right;
}

/* Workflow pill row — left-aligned step indicator under the hero. */
.workflow-pills {
  gap: 4px;
}
.workflow-pill {
  font-weight: 500;
}

/* Generic key / value row used in Member + Cover cards for a clean,
   right-aligned value column. Single-line by default; wraps gracefully
   on narrow viewports. */
.key-value-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  padding: 4px 0;
  gap: 12px;
  font-size: 0.875rem;
  flex-wrap: wrap;
}
.key-label {
  color: rgba(0, 0, 0, 0.6);
}
.key-value {
  color: rgba(0, 0, 0, 0.87);
  text-align: right;
}

/* Callout block for the headline SA/FCL ratio at the bottom of the
   Cover vs FCL card. Tinted by tier severity. */
.ratio-callout {
  margin-top: 6px;
  padding: 10px 12px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(0, 0, 0, 0.02);
}
.ratio-callout.ratio-success {
  background: rgba(76, 175, 80, 0.08);
  border-color: rgba(76, 175, 80, 0.2);
}
.ratio-callout.ratio-warning {
  background: rgba(255, 152, 0, 0.08);
  border-color: rgba(255, 152, 0, 0.25);
}
.ratio-callout.ratio-error {
  background: rgba(244, 67, 54, 0.08);
  border-color: rgba(244, 67, 54, 0.25);
}

/* Card titles with leading icons should sit closer to their content. */
.v-card-title.d-flex .v-icon {
  flex-shrink: 0;
}
</style>
