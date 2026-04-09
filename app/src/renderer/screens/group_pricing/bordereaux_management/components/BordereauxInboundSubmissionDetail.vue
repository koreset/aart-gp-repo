<template>
  <v-container fluid>
    <base-card :show-actions="true">
      <template #header>
        <h3 class="mb-0">Employer Submission</h3>
      </template>
      <template #default>
        <page-header
          title="Employer Submission"
          :subtitle="
            submission
              ? `${submission.scheme_name} — ${submission.month}/${submission.year}`
              : ''
          "
          icon="mdi-inbox-arrow-down"
          :breadcrumbs="[
            {
              text: 'Inbound Submissions',
              to: { name: 'group-pricing-bordereaux-inbound' }
            },
            { text: `Submission #${submissionId}` }
          ]"
        />
        <v-row v-if="submission" class="mb-2" align="center">
          <v-col cols="auto">
            <v-chip
              :color="statusColor(submission.status)"
              size="small"
              variant="tonal"
            >
              {{ (submission.status ?? '').replace(/_/g, ' ') }}
            </v-chip>
            <v-chip
              v-if="submission.is_retro"
              color="deep-purple"
              class="ml-2"
              size="small"
              variant="tonal"
              prepend-icon="mdi-history"
            >
              Retrospective
            </v-chip>
          </v-col>
        </v-row>

        <div v-if="loading" class="text-center pa-8">
          <v-progress-circular indeterminate color="primary" />
        </div>

        <template v-else-if="submission">
          <!-- Meta card -->
          <v-row class="mb-3">
            <v-col cols="12" md="6">
              <v-card variant="outlined">
                <v-card-text>
                  <div class="d-flex flex-wrap gap-4 text-body-2">
                    <span class="mr-3"
                      ><strong>Scheme:</strong>
                      {{ submission.scheme_name }}</span
                    >
                    <span class="mr-3"
                      ><strong>Period:</strong> {{ submission.month }}/{{
                        submission.year
                      }}</span
                    >
                    <span v-if="submission.due_date" class="mr-3"
                      ><strong>Due:</strong> {{ submission.due_date }}</span
                    >
                    <span v-if="submission.received_date" class="mr-3"
                      ><strong>Received:</strong>
                      {{ fmtDate(submission.received_date) }}</span
                    >
                    <span v-if="submission.submitted_by" class="mr-3"
                      ><strong>Submitted by:</strong>
                      {{ submission.submitted_by }}</span
                    >
                    <span
                      v-if="
                        submission.is_retro && submission.retro_effective_date
                      "
                      class="mr-3"
                    >
                      <strong>Retro Effective:</strong>
                      {{ submission.retro_effective_date }}
                    </span>
                    <span v-if="submission.exits_synced_at" class="mr-3">
                      <v-icon size="14" color="teal" class="mr-1"
                        >mdi-account-off</v-icon
                      >
                      <strong>Exits Applied:</strong>
                      {{ fmtDate(submission.exits_synced_at) }}
                    </span>
                    <span v-if="submission.amendments_synced_at" class="mr-3">
                      <v-icon size="14" color="teal" class="mr-1"
                        >mdi-account-edit</v-icon
                      >
                      <strong>Amendments Applied:</strong>
                      {{ fmtDate(submission.amendments_synced_at) }}
                    </span>
                    <span v-if="submission.new_joiners_synced_at" class="mr-3">
                      <v-icon size="14" color="teal" class="mr-1"
                        >mdi-account-plus</v-icon
                      >
                      <strong>New Joiners Synced:</strong>
                      {{ fmtDate(submission.new_joiners_synced_at) }}
                    </span>
                  </div>
                  <div v-if="submission.notes" class="mt-2 text-body-2">
                    <strong>Notes:</strong> {{ submission.notes }}
                  </div>
                </v-card-text>
              </v-card>
            </v-col>

            <!-- Validation summary (shown after upload) -->
            <v-col v-if="submission.record_count > 0" cols="12" md="6">
              <v-card variant="outlined">
                <v-card-text>
                  <div class="d-flex gap-3">
                    <v-chip
                      class="mr-3"
                      color="primary"
                      variant="tonal"
                      size="small"
                    >
                      Total: {{ submission.record_count }}
                    </v-chip>
                    <v-chip
                      class="mr-3"
                      color="success"
                      variant="tonal"
                      size="small"
                    >
                      Valid: {{ submission.valid_count }}
                    </v-chip>
                    <v-chip
                      class="mr-3"
                      color="error"
                      variant="tonal"
                      size="small"
                    >
                      Invalid: {{ submission.invalid_count }}
                    </v-chip>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Invalid Records (shown after upload when there are validation failures) -->
          <v-row v-if="invalidRecords.length > 0" class="mb-3">
            <v-col cols="12">
              <v-card variant="outlined" border="error">
                <v-card-title
                  class="text-subtitle-1 d-flex align-center justify-space-between"
                >
                  <span>
                    <v-icon color="error" size="18" class="mr-1"
                      >mdi-alert-circle-outline</v-icon
                    >
                    Invalid Records ({{ invalidRecords.length }})
                  </span>
                  <v-btn
                    size="small"
                    variant="outlined"
                    color="error"
                    prepend-icon="mdi-file-download-outline"
                    @click="downloadInvalidRecords"
                  >
                    Export CSV
                  </v-btn>
                </v-card-title>
                <v-card-subtitle class="pb-2">
                  These rows could not be processed. Correct them in your source
                  file and re-upload, or contact the employer to provide updated
                  data.
                </v-card-subtitle>
                <v-card-text class="pa-0">
                  <ag-grid-vue
                    class="ag-theme-balham"
                    :style="{ height: invalidGridHeight, width: '100%' }"
                    :column-defs="invalidRecordColDefs"
                    :row-data="invalidRecords"
                    :default-col-def="recordDefaultColDef"
                  />
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Upload section (pending_receipt or received) -->
          <v-row
            v-if="
              submission.status === 'pending_receipt' ||
              submission.status === 'received'
            "
            class="mb-3"
          >
            <v-col cols="12">
              <v-card variant="outlined">
                <v-card-title class="text-subtitle-1"
                  >Upload Member Data File</v-card-title
                >
                <v-card-text>
                  <v-row align="center" class="mb-2">
                    <v-col cols="12" md="6">
                      <v-file-input
                        v-model="uploadFile"
                        label="Select file (.csv or .xlsx)"
                        accept=".csv,.xlsx"
                        variant="outlined"
                        density="compact"
                        prepend-inner-icon="mdi-file-upload"
                        prepend-icon=""
                        clearable
                      />
                    </v-col>
                    <v-col cols="auto">
                      <v-btn
                        color="primary"
                        :loading="uploading"
                        :disabled="!uploadFile"
                        prepend-icon="mdi-upload"
                        @click="handleUpload"
                      >
                        Upload
                      </v-btn>
                    </v-col>
                    <v-col cols="auto">
                      <v-btn
                        variant="outlined"
                        color="secondary"
                        prepend-icon="mdi-file-download-outline"
                        @click="downloadTemplate"
                      >
                        Download Template
                      </v-btn>
                    </v-col>
                  </v-row>

                  <!-- Column guide -->
                  <v-expansion-panels variant="accordion" class="mt-1">
                    <v-expansion-panel>
                      <v-expansion-panel-title
                        class="text-body-2 font-weight-medium"
                      >
                        <v-icon size="16" class="mr-2" color="info"
                          >mdi-information-outline</v-icon
                        >
                        Required file format &amp; column reference
                      </v-expansion-panel-title>
                      <v-expansion-panel-text>
                        <p class="text-body-2 mb-3">
                          Your file must be <strong>.csv</strong> or
                          <strong>.xlsx</strong> with a header row. Column names
                          are matched case-insensitively. Each row represents
                          one member. The table below shows every recognised
                          column.
                        </p>

                        <v-table density="compact" class="text-body-2 mb-3">
                          <thead>
                            <tr>
                              <th>Column</th>
                              <th>Required</th>
                              <th>Accepted header names</th>
                              <th>Notes</th>
                            </tr>
                          </thead>
                          <tbody>
                            <tr v-for="col in columnGuide" :key="col.name">
                              <td>
                                <code>{{ col.name }}</code>
                              </td>
                              <td>
                                <v-chip
                                  :color="
                                    col.required === 'Required'
                                      ? 'error'
                                      : col.required === 'Recommended'
                                        ? 'warning'
                                        : 'default'
                                  "
                                  size="x-small"
                                  variant="tonal"
                                >
                                  {{ col.required }}
                                </v-chip>
                              </td>
                              <td class="text-caption text-medium-emphasis">
                                {{ col.aliases.join(', ') }}
                              </td>
                              <td class="text-caption">{{ col.notes }}</td>
                            </tr>
                          </tbody>
                        </v-table>

                        <v-alert
                          type="info"
                          variant="tonal"
                          density="compact"
                          class="text-body-2"
                        >
                          <strong>Register diff tip:</strong> Include
                          <code>employee_no</code> or <code>id_number</code> so
                          the system can reliably identify new joiners, leavers,
                          and amendments when comparing against the live member
                          register. Include <code>gender</code> to track gender
                          amendments.
                        </v-alert>
                      </v-expansion-panel-text>
                    </v-expansion-panel>
                  </v-expansion-panels>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Query Notes (queries_raised+) -->
          <v-row v-if="submission.query_notes" class="mb-3">
            <v-col cols="12">
              <v-alert type="warning" variant="tonal" icon="mdi-message-alert">
                <strong>Query Notes:</strong> {{ submission.query_notes }}
              </v-alert>
            </v-col>
          </v-row>

          <!-- Audit trail -->
          <v-row
            v-if="
              submission.reviewed_by ||
              submission.accepted_by ||
              submission.rejected_by
            "
            class="mb-3"
          >
            <v-col cols="12">
              <v-card variant="outlined" density="compact">
                <v-card-text class="py-2">
                  <div class="d-flex flex-wrap ga-4 text-caption">
                    <span v-if="submission.reviewed_by">
                      <strong>Reviewed by:</strong> {{ submission.reviewed_by }}
                      <span v-if="submission.reviewed_at">
                        on {{ fmtDate(submission.reviewed_at) }}</span
                      >
                    </span>
                    <span v-if="submission.accepted_by">
                      <strong>Accepted by:</strong> {{ submission.accepted_by }}
                      <span v-if="submission.accepted_at">
                        on {{ fmtDate(submission.accepted_at) }}</span
                      >
                    </span>
                    <span v-if="submission.rejected_by">
                      <strong>Rejected by:</strong> {{ submission.rejected_by }}
                      <span v-if="submission.rejected_at">
                        on {{ fmtDate(submission.rejected_at) }}</span
                      >
                      <span v-if="submission.rejection_reason">
                        — {{ submission.rejection_reason }}</span
                      >
                    </span>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Rejection banner -->
          <v-row v-if="submission.status === 'rejected'" class="mb-3">
            <v-col cols="12">
              <v-alert type="error" variant="tonal">
                This submission was <strong>rejected</strong>.
                <span v-if="submission.rejection_reason">
                  Reason: {{ submission.rejection_reason }}</span
                >
              </v-alert>
            </v-col>
          </v-row>

          <!-- Linked Schedule (accepted) -->
          <v-row v-if="submission.status === 'accepted'" class="mb-3">
            <v-col cols="12">
              <v-card variant="outlined">
                <v-card-title class="text-subtitle-1">
                  {{
                    submission.is_retro
                      ? 'Supplementary Retro Schedule'
                      : 'Premium Schedule'
                  }}
                </v-card-title>
                <v-card-text>
                  <template v-if="submission.linked_premium_schedule_id">
                    <v-chip color="success" variant="tonal" class="mr-3">
                      Schedule #{{ submission.linked_premium_schedule_id }}
                    </v-chip>
                    <v-btn
                      size="small"
                      variant="outlined"
                      color="primary"
                      prepend-icon="mdi-arrow-right"
                      @click="viewSchedule"
                    >
                      View Schedule
                    </v-btn>
                  </template>
                  <template v-else>
                    <p class="text-body-2 mb-3">
                      <template v-if="submission.is_retro">
                        This retrospective submission has been accepted.
                        Generating the schedule will create a supplementary
                        catch-up schedule and invoice for the period from
                        <strong>{{ submission.retro_effective_date }}</strong>
                        to
                        <strong
                          >{{ submission.month }}/{{ submission.year }}</strong
                        >.
                      </template>
                      <template v-else>
                        This submission has been accepted. You can now generate
                        the premium schedule.
                      </template>
                    </p>
                    <v-btn
                      color="primary"
                      :loading="generatingSchedule"
                      :prepend-icon="
                        submission.is_retro
                          ? 'mdi-history'
                          : 'mdi-calendar-plus'
                      "
                      @click="handleGenerateSchedule"
                    >
                      {{
                        submission.is_retro
                          ? 'Generate Retro Supplementary Schedule'
                          : 'Generate Premium Schedule'
                      }}
                    </v-btn>
                  </template>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Records Grid -->
          <v-row v-if="(submission.records ?? []).length > 0" class="mb-3">
            <v-col cols="12">
              <v-card variant="outlined">
                <v-card-title class="text-subtitle-1"
                  >Member Records</v-card-title
                >
                <v-card-text class="pa-0">
                  <ag-grid-vue
                    class="ag-theme-balham"
                    style="height: 420px; width: 100%"
                    :column-defs="recordColumnDefs"
                    :row-data="submission.records"
                    :default-col-def="recordDefaultColDef"
                  />
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Register Diff Preview (shown whenever records have been uploaded) -->
          <v-row v-if="submission.record_count > 0" class="mb-3">
            <v-col cols="12">
              <v-card variant="outlined">
                <v-card-title
                  class="text-subtitle-1 d-flex align-center justify-space-between"
                >
                  <span>Member Register Diff</span>
                  <div class="d-flex align-center gap-2">
                    <!-- Provenance badge -->
                    <v-chip
                      v-if="registerDiff?.is_snapshot"
                      color="teal"
                      variant="tonal"
                      size="small"
                      prepend-icon="mdi-camera-outline"
                    >
                      Snapshot ·
                      {{ formatSnapshotDate(registerDiff.snapshot_at) }}
                    </v-chip>
                    <v-chip
                      v-else-if="registerDiff && !registerDiff.is_snapshot"
                      color="orange"
                      variant="tonal"
                      size="small"
                      prepend-icon="mdi-lightning-bolt"
                    >
                      Live view
                    </v-chip>
                    <!-- Refresh — only shown for non-accepted (pre-snapshot) submissions -->
                    <v-btn
                      v-if="submission.status !== 'accepted'"
                      size="small"
                      color="primary"
                      variant="tonal"
                      prepend-icon="mdi-refresh"
                      :loading="registerDiffLoading"
                      @click="loadRegisterDiff"
                    >
                      Refresh
                    </v-btn>
                    <!-- Re-snapshot — shown for accepted submissions so admins can backfill -->
                    <v-btn
                      v-else
                      size="small"
                      color="secondary"
                      variant="tonal"
                      prepend-icon="mdi-camera-retake"
                      :loading="registerDiffLoading"
                      @click="handleReSnapshot"
                    >
                      Re-snapshot
                    </v-btn>
                  </div>
                </v-card-title>
                <v-card-text>
                  <div v-if="registerDiffLoading" class="text-center py-4">
                    <v-progress-circular
                      indeterminate
                      color="primary"
                      size="24"
                    />
                  </div>
                  <template v-else-if="registerDiff">
                    <!-- Snapshot provenance detail -->
                    <v-alert
                      v-if="registerDiff.is_snapshot"
                      type="info"
                      variant="tonal"
                      density="compact"
                      class="text-body-2 mb-3"
                      prepend-icon="mdi-history"
                    >
                      Historical snapshot captured on
                      <strong>{{
                        formatSnapshotDate(registerDiff.snapshot_at)
                      }}</strong>
                      <span v-if="registerDiff.snapshot_by">
                        by {{ registerDiff.snapshot_by }}</span
                      >. This record is frozen and will not change even if the
                      member register is later updated.
                    </v-alert>
                    <v-alert
                      v-else
                      type="warning"
                      variant="tonal"
                      density="compact"
                      class="text-body-2 mb-3"
                      prepend-icon="mdi-alert-outline"
                    >
                      This is a <strong>live view</strong> against the current
                      member register. Accept the submission to lock in a
                      permanent historical snapshot.
                    </v-alert>
                    <!-- Summary chips -->
                    <div class="d-flex flex-wrap gap-3 mb-4">
                      <v-chip
                        color="success"
                        variant="tonal"
                        size="small"
                        prepend-icon="mdi-account-plus"
                      >
                        New Joiners: {{ registerDiff.new_joiners?.length ?? 0 }}
                      </v-chip>
                      <v-chip
                        color="error"
                        variant="tonal"
                        size="small"
                        prepend-icon="mdi-account-off"
                      >
                        Exits: {{ registerDiff.exits?.length ?? 0 }}
                      </v-chip>
                      <v-chip
                        color="warning"
                        variant="tonal"
                        size="small"
                        prepend-icon="mdi-account-edit"
                      >
                        Amendments: {{ registerDiff.amendments?.length ?? 0 }}
                      </v-chip>
                      <v-chip color="primary" variant="tonal" size="small">
                        Continuing: {{ registerDiff.continuing ?? 0 }}
                      </v-chip>
                    </div>

                    <v-expansion-panels variant="accordion">
                      <!-- New Joiners panel -->
                      <v-expansion-panel>
                        <v-expansion-panel-title
                          class="text-body-2 font-weight-medium"
                        >
                          <v-icon size="16" class="mr-2" color="success"
                            >mdi-account-plus</v-icon
                          >
                          New Joiners ({{
                            registerDiff.new_joiners?.length ?? 0
                          }})
                        </v-expansion-panel-title>
                        <v-expansion-panel-text>
                          <v-alert
                            type="info"
                            variant="tonal"
                            density="compact"
                            class="text-body-2 mb-3"
                          >
                            Full details required for these members — see the
                            New Joiner Details section below to upload and enrol
                            them.
                          </v-alert>
                          <ag-grid-vue
                            v-if="(registerDiff.new_joiners?.length ?? 0) > 0"
                            class="ag-theme-balham"
                            style="height: 280px; width: 100%"
                            :column-defs="diffNewJoinerColDefs"
                            :row-data="registerDiff.new_joiners"
                            :default-col-def="recordDefaultColDef"
                          />
                          <p v-else class="text-body-2 text-medium-emphasis"
                            >No new joiners identified.</p
                          >
                        </v-expansion-panel-text>
                      </v-expansion-panel>

                      <!-- Exits panel -->
                      <v-expansion-panel>
                        <v-expansion-panel-title
                          class="text-body-2 font-weight-medium"
                        >
                          <v-icon size="16" class="mr-2" color="error"
                            >mdi-account-off</v-icon
                          >
                          Exits ({{ registerDiff.exits?.length ?? 0 }})
                        </v-expansion-panel-title>
                        <v-expansion-panel-text>
                          <v-alert
                            type="warning"
                            variant="tonal"
                            density="compact"
                            class="text-body-2 mb-3"
                          >
                            These members will be deactivated in the live
                            register when "Apply Exits" is run.
                          </v-alert>
                          <ag-grid-vue
                            v-if="(registerDiff.exits?.length ?? 0) > 0"
                            class="ag-theme-balham"
                            style="height: 280px; width: 100%"
                            :column-defs="diffExitColDefs"
                            :row-data="registerDiff.exits"
                            :default-col-def="recordDefaultColDef"
                          />
                          <p v-else class="text-body-2 text-medium-emphasis"
                            >No exits identified.</p
                          >
                        </v-expansion-panel-text>
                      </v-expansion-panel>

                      <!-- Amendments panel -->
                      <v-expansion-panel>
                        <v-expansion-panel-title
                          class="text-body-2 font-weight-medium"
                        >
                          <v-icon size="16" class="mr-2" color="warning"
                            >mdi-account-edit</v-icon
                          >
                          Amendments ({{
                            registerDiff.amendments?.length ?? 0
                          }})
                        </v-expansion-panel-title>
                        <v-expansion-panel-text>
                          <v-alert
                            type="info"
                            variant="tonal"
                            density="compact"
                            class="text-body-2 mb-3"
                          >
                            These field changes will be applied to the live
                            register when "Apply Amendments" is run.
                          </v-alert>
                          <ag-grid-vue
                            v-if="(registerDiff.amendments?.length ?? 0) > 0"
                            class="ag-theme-balham"
                            style="height: 280px; width: 100%"
                            :column-defs="diffAmendmentColDefs"
                            :row-data="registerDiff.amendments"
                            :default-col-def="recordDefaultColDef"
                          />
                          <p v-else class="text-body-2 text-medium-emphasis"
                            >No amendments identified.</p
                          >
                        </v-expansion-panel-text>
                      </v-expansion-panel>
                    </v-expansion-panels>
                  </template>
                  <div v-else class="text-body-2 text-medium-emphasis">
                    Register diff will load automatically after upload.
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Register Sync Actions -->
          <template
            v-if="
              submission.status === 'accepted' ||
              submission.status === 'received'
            "
          >
            <!-- Card A: Apply Exits -->
            <v-row v-if="(registerDiff?.exits?.length ?? 0) > 0" class="mb-3">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="text-subtitle-1 d-flex align-center gap-2"
                  >
                    <v-icon color="error">mdi-account-off</v-icon>
                    Apply Exits
                    <v-chip
                      v-if="submission.exits_synced_at"
                      color="success"
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      Applied {{ fmtDate(submission.exits_synced_at) }}
                    </v-chip>
                  </v-card-title>
                  <v-card-text>
                    <p class="text-body-2 mb-3">
                      Deactivate {{ registerDiff.exits.length }} member(s) with
                      past exit dates in the live register.
                    </p>
                    <div v-if="exitResult" class="d-flex flex-wrap gap-2 mb-3">
                      <v-chip color="error" variant="tonal" size="small"
                        >{{ exitResult.deactivated }} deactivated</v-chip
                      >
                      <v-chip color="grey" variant="tonal" size="small"
                        >{{ exitResult.skipped }} skipped</v-chip
                      >
                      <v-chip
                        v-if="exitResult.errors?.length > 0"
                        color="warning"
                        variant="tonal"
                        size="small"
                      >
                        {{ exitResult.errors.length }} errors
                      </v-chip>
                    </div>
                    <v-btn
                      color="error"
                      variant="tonal"
                      prepend-icon="mdi-account-off"
                      :loading="applyingExits"
                      :disabled="!!submission.exits_synced_at"
                      @click="exitConfirmDialog = true"
                    >
                      {{
                        submission.exits_synced_at
                          ? 'Exits Applied'
                          : 'Apply Exits'
                      }}
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Card B: Apply Amendments -->
            <v-row
              v-if="(registerDiff?.amendments?.length ?? 0) > 0"
              class="mb-3"
            >
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="text-subtitle-1 d-flex align-center gap-2"
                  >
                    <v-icon color="warning">mdi-account-edit</v-icon>
                    Apply Amendments
                    <v-chip
                      v-if="submission.amendments_synced_at"
                      color="success"
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      Applied {{ fmtDate(submission.amendments_synced_at) }}
                    </v-chip>
                  </v-card-title>
                  <v-card-text>
                    <p class="text-body-2 mb-3">
                      Apply {{ registerDiff.amendments.length }} field change(s)
                      to existing members in the live register.
                    </p>
                    <div
                      v-if="amendmentResult"
                      class="d-flex flex-wrap gap-2 mb-3"
                    >
                      <v-chip color="warning" variant="tonal" size="small"
                        >{{ amendmentResult.updated }} updated</v-chip
                      >
                      <v-chip color="grey" variant="tonal" size="small"
                        >{{ amendmentResult.skipped }} skipped</v-chip
                      >
                      <v-chip
                        v-if="amendmentResult.errors?.length > 0"
                        color="error"
                        variant="tonal"
                        size="small"
                      >
                        {{ amendmentResult.errors.length }} errors
                      </v-chip>
                    </div>
                    <v-btn
                      color="warning"
                      variant="tonal"
                      prepend-icon="mdi-account-edit"
                      :loading="applyingAmendments"
                      :disabled="!!submission.amendments_synced_at"
                      @click="amendmentConfirmDialog = true"
                    >
                      {{
                        submission.amendments_synced_at
                          ? 'Amendments Applied'
                          : 'Apply Amendments'
                      }}
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Card C: New Joiner Details -->
            <v-row
              v-if="(registerDiff?.new_joiners?.length ?? 0) > 0"
              class="mb-3"
            >
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="text-subtitle-1 d-flex align-center gap-2"
                  >
                    <v-icon color="success">mdi-account-plus</v-icon>
                    New Joiner Details
                    <v-chip
                      v-if="submission.new_joiners_synced_at"
                      color="success"
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      Synced {{ fmtDate(submission.new_joiners_synced_at) }}
                    </v-chip>
                  </v-card-title>
                  <v-card-text>
                    <p class="text-body-2 mb-3">
                      {{ registerDiff.new_joiners.length }} new joiner(s)
                      identified. Upload a detailed file to enrol them in the
                      member register.
                    </p>
                    <v-row align="center" class="mb-2">
                      <v-col cols="12" md="6">
                        <v-file-input
                          v-model="newJoinerDetailFile"
                          label="New joiner detail file (.csv or .xlsx)"
                          accept=".csv,.xlsx"
                          variant="outlined"
                          density="compact"
                          prepend-inner-icon="mdi-file-upload"
                          prepend-icon=""
                          clearable
                        />
                      </v-col>
                      <v-col cols="auto">
                        <v-btn
                          color="primary"
                          :loading="uploadingNewJoinerDetails"
                          :disabled="!newJoinerDetailFile"
                          prepend-icon="mdi-upload"
                          @click="handleUploadNewJoinerDetails"
                        >
                          Upload Details
                        </v-btn>
                      </v-col>
                      <v-col cols="auto">
                        <v-btn
                          variant="outlined"
                          color="secondary"
                          prepend-icon="mdi-file-download-outline"
                          @click="downloadNewJoinerTemplate"
                        >
                          Download Template
                        </v-btn>
                      </v-col>
                    </v-row>

                    <div v-if="newJoinerDetails.length > 0" class="mb-3">
                      <v-chip
                        color="primary"
                        variant="tonal"
                        size="small"
                        class="mb-2"
                      >
                        {{ newJoinerDetails.length }} records staged ({{
                          newJoinerDetails.filter(
                            (r) => r.validation_status === 'valid'
                          ).length
                        }}
                        valid)
                      </v-chip>
                    </div>

                    <div
                      v-if="newJoinerSyncResult"
                      class="d-flex flex-wrap gap-2 mb-3"
                    >
                      <v-chip color="success" variant="tonal" size="small"
                        >{{ newJoinerSyncResult.added }} enrolled</v-chip
                      >
                      <v-chip color="grey" variant="tonal" size="small"
                        >{{ newJoinerSyncResult.skipped }} skipped</v-chip
                      >
                      <v-chip
                        v-if="newJoinerSyncResult.errors?.length > 0"
                        color="warning"
                        variant="tonal"
                        size="small"
                      >
                        {{ newJoinerSyncResult.errors.length }} errors
                      </v-chip>
                    </div>

                    <v-btn
                      v-if="newJoinerDetails.length > 0"
                      color="success"
                      variant="tonal"
                      prepend-icon="mdi-account-plus"
                      :loading="syncingNewJoiners"
                      :disabled="!!submission.new_joiners_synced_at"
                      @click="newJoinerSyncConfirmDialog = true"
                    >
                      {{
                        submission.new_joiners_synced_at
                          ? 'New Joiners Synced'
                          : 'Sync New Joiners'
                      }}
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </template>
      </template>
      <template #actions>
        <!-- Action bar (outside base-card) -->
        <v-row v-if="submission" class="mt-3 px-4">
          <v-col cols="12" class="d-flex gap-2 flex-wrap">
            <v-btn
              v-if="submission.status === 'received'"
              color="info"
              variant="tonal"
              prepend-icon="mdi-eye-check"
              :loading="actionLoading"
              @click="handleReview"
            >
              Mark as Under Review
            </v-btn>

            <template v-if="submission.status === 'under_review'">
              <v-btn
                color="warning"
                variant="tonal"
                prepend-icon="mdi-message-question"
                @click="raiseQueryDialog = true"
              >
                Raise Query
              </v-btn>
              <v-btn
                color="success"
                variant="tonal"
                prepend-icon="mdi-check"
                @click="confirmAcceptDialog = true"
              >
                Accept
              </v-btn>
              <v-btn
                color="error"
                variant="tonal"
                prepend-icon="mdi-close"
                @click="rejectDialog = true"
              >
                Reject
              </v-btn>
            </template>

            <template v-if="submission.status === 'queries_raised'">
              <v-btn
                color="success"
                variant="tonal"
                prepend-icon="mdi-check"
                @click="confirmAcceptDialog = true"
              >
                Accept
              </v-btn>
              <v-btn
                color="error"
                variant="tonal"
                prepend-icon="mdi-close"
                @click="rejectDialog = true"
              >
                Reject
              </v-btn>
            </template>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Raise Query Dialog -->
    <v-dialog v-model="raiseQueryDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Raise Query</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="queryNotes"
            label="Query Notes *"
            variant="outlined"
            density="compact"
            rows="3"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="raiseQueryDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="warning"
            :loading="actionLoading"
            :disabled="!queryNotes.trim()"
            @click="handleRaiseQuery"
          >
            Raise Query
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Confirm Accept Dialog -->
    <v-dialog v-model="confirmAcceptDialog" max-width="400">
      <v-card>
        <v-card-title>Confirm Acceptance</v-card-title>
        <v-card-text>
          Accept this employer submission? This will allow a premium schedule to
          be generated.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="confirmAcceptDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="success" :loading="actionLoading" @click="handleAccept"
            >Accept</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Reject Dialog -->
    <v-dialog v-model="rejectDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Reject Submission</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="rejectionReason"
            label="Rejection Reason *"
            variant="outlined"
            density="compact"
            rows="3"
            auto-grow
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="rejectDialog = false">Cancel</v-btn>
          <v-btn
            color="error"
            :loading="actionLoading"
            :disabled="!rejectionReason.trim()"
            @click="handleReject"
          >
            Reject
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Apply Exits Confirmation Dialog -->
    <v-dialog v-model="exitConfirmDialog" max-width="460" persistent>
      <v-card>
        <v-card-title class="d-flex align-center gap-2">
          <v-icon color="error">mdi-account-off</v-icon>
          Apply Exits
        </v-card-title>
        <v-card-text>
          Deactivate {{ registerDiff?.exits?.length ?? 0 }} member(s) with past
          exit dates in the live member register? This action updates the
          register immediately.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="exitConfirmDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="error"
            :loading="applyingExits"
            @click="handleApplyExits"
            >Apply Exits</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Apply Amendments Confirmation Dialog -->
    <v-dialog v-model="amendmentConfirmDialog" max-width="460" persistent>
      <v-card>
        <v-card-title class="d-flex align-center gap-2">
          <v-icon color="warning">mdi-account-edit</v-icon>
          Apply Amendments
        </v-card-title>
        <v-card-text>
          Apply {{ registerDiff?.amendments?.length ?? 0 }} field change(s) to
          existing members in the live register? Only changed fields will be
          updated.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="amendmentConfirmDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="warning"
            :loading="applyingAmendments"
            @click="handleApplyAmendments"
          >
            Apply Amendments
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Sync New Joiners Confirmation Dialog -->
    <v-dialog v-model="newJoinerSyncConfirmDialog" max-width="460" persistent>
      <v-card>
        <v-card-title class="d-flex align-center gap-2">
          <v-icon color="success">mdi-account-plus</v-icon>
          Sync New Joiners
        </v-card-title>
        <v-card-text>
          Enrol
          {{
            newJoinerDetails.filter((r) => r.validation_status === 'valid')
              .length
          }}
          new member(s) into the live member register using the uploaded detail
          records?
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="newJoinerSyncConfirmDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="success"
            :loading="syncingNewJoiners"
            @click="handleSyncNewJoiners"
          >
            Sync New Joiners
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" timeout="4000">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import { AgGridVue } from 'ag-grid-vue3'
import BaseCard from '@/renderer/components/BaseCard.vue'
import PageHeader from '@/renderer/components/PageHeader.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

const invalidGridHeight = useGridHeight(500)

const props = defineProps<{ submissionId: string | number }>()
const router = useRouter()

const loading = ref(false)
const uploading = ref(false)
const actionLoading = ref(false)
const generatingSchedule = ref(false)
const submission = ref<any>(null)
const uploadFile = ref<File | null>(null)

// Register diff
const registerDiff: any = ref(null)
const registerDiffLoading = ref(false)

// Per-category sync state
const applyingExits = ref(false)
const applyingAmendments = ref(false)
const exitResult = ref<any>(null)
const amendmentResult = ref<any>(null)

// New joiner details
const newJoinerDetailFile = ref<File | null>(null)
const uploadingNewJoinerDetails = ref(false)
const newJoinerDetails = ref<any[]>([])
const syncingNewJoiners = ref(false)
const newJoinerSyncResult = ref<any>(null)

// Dialogs
const raiseQueryDialog = ref(false)
const confirmAcceptDialog = ref(false)
const rejectDialog = ref(false)
const exitConfirmDialog = ref(false)
const amendmentConfirmDialog = ref(false)
const newJoinerSyncConfirmDialog = ref(false)

const queryNotes = ref('')
const rejectionReason = ref('')
const snackbar = ref({ show: false, message: '', color: 'success' })

const invalidRecords = computed(() =>
  (submission.value?.records ?? []).filter(
    (r: any) => r.validation_status !== 'valid'
  )
)

const statusColor = (s: string) => {
  const map: Record<string, string> = {
    pending_receipt: 'grey',
    received: 'blue',
    under_review: 'orange',
    queries_raised: 'amber',
    accepted: 'success',
    rejected: 'error'
  }
  return map[s] ?? 'grey'
}

const fmtDate = (v: string | null) => {
  if (!v) return ''
  return new Date(v).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const fmtCurrency = (v: number) =>
  v == null
    ? ''
    : new Intl.NumberFormat('en-ZA', {
        style: 'currency',
        currency: 'ZAR',
        minimumFractionDigits: 2
      }).format(v)

const recordDefaultColDef = { sortable: true, filter: true, resizable: true }

const columnGuide = [
  {
    name: 'member_name',
    required: 'Required',
    aliases: ['member_name', 'name', 'full_name', 'fullname'],
    notes: 'Full name of the member. Row is rejected if blank.'
  },
  {
    name: 'employee_no',
    required: 'Recommended',
    aliases: [
      'employee_no',
      'employee_number',
      'emp_no',
      'emp_number',
      'staff_no'
    ],
    notes:
      'Primary key used for delta tracking. Strongly recommended to include.'
  },
  {
    name: 'id_number',
    required: 'Recommended',
    aliases: ['id_number', 'id_no', 'sa_id', 'id'],
    notes:
      'RSA ID (13 digits, Luhn-validated) or passport number (5–20 alphanumeric chars). Type is auto-detected; optionally supply id_type to override.'
  },
  {
    name: 'gender',
    required: 'Optional',
    aliases: ['gender', 'sex'],
    notes: 'M or F. Used for amendment tracking against the live register.'
  },
  {
    name: 'dob',
    required: 'Optional',
    aliases: ['dob', 'date_of_birth', 'birth_date'],
    notes: 'Date of birth. Accepted formats: YYYY-MM-DD, DD/MM/YYYY.'
  },
  {
    name: 'salary',
    required: 'Optional',
    aliases: ['salary', 'annual_salary', 'gross_salary'],
    notes: 'Annual gross salary in ZAR. Used for benefit multiple calculations.'
  },
  {
    name: 'benefit_code',
    required: 'Optional',
    aliases: ['benefit_code', 'benefit', 'benefit_type'],
    notes: 'Benefit type code, e.g. GLA, PTD, CI, TTD, PHI, GFF.'
  },
  {
    name: 'premium_amount',
    required: 'Optional',
    aliases: ['premium_amount', 'premium', 'monthly_premium'],
    notes: 'Monthly premium in ZAR. Must not be negative.'
  },
  {
    name: 'entry_date',
    required: 'Optional',
    aliases: ['entry_date', 'commencement_date', 'start_date'],
    notes: 'Date the member joined the scheme. Format: YYYY-MM-DD.'
  },
  {
    name: 'exit_date',
    required: 'Optional',
    aliases: ['exit_date', 'termination_date', 'end_date'],
    notes: 'Date the member left (leavers only). Format: YYYY-MM-DD.'
  }
]

const downloadTemplate = () => {
  const headers = [
    'member_name',
    'employee_no',
    'id_number',
    'id_type',
    'gender',
    'dob',
    'salary',
    'benefit_code',
    'premium_amount',
    'entry_date',
    'exit_date'
  ]
  const exampleRow = [
    'John Smith',
    'EMP001',
    '8001015009087',
    'rsa_id',
    'M',
    '1980-01-01',
    '240000',
    'GLA',
    '850.00',
    '2024-01-01',
    ''
  ]
  const csv = [headers.join(','), exampleRow.join(',')].join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'member_submission_template.csv'
  a.click()
  URL.revokeObjectURL(url)
}

const downloadNewJoinerTemplate = () => {
  const headers = [
    'member_name',
    'employee_no',
    'id_number',
    'id_type',
    'gender',
    'dob',
    'salary',
    'scheme_category',
    'benefits_gla_multiple',
    'benefits_ptd_multiple',
    'benefits_ci_multiple',
    'benefits_ttd_multiple',
    'benefits_phi_multiple',
    'benefits_sgla_multiple',
    'address_line1',
    'address_line2',
    'city',
    'province',
    'postal_code',
    'phone',
    'email',
    'occupation',
    'occupational_class',
    'entry_date'
  ]
  const exampleRow = [
    'Jane Doe',
    'EMP002',
    '9001020084082',
    'RSA_ID',
    'F',
    '1990-01-02',
    '180000',
    '3',
    '0.5',
    '0.3',
    '0.2',
    '0.1',
    '2',
    '12 Main St',
    '',
    'Johannesburg',
    'Gauteng',
    '2001',
    '0821234567',
    'jane@example.com',
    'Clerk',
    'Class 1',
    '2024-03-01'
  ]
  const csv = [headers.join(','), exampleRow.join(',')].join('\n')
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'new_joiner_detail_template.csv'
  a.click()
  URL.revokeObjectURL(url)
}

const validationBadge = (value: string) => {
  const colors: Record<string, string> = {
    valid: '#4CAF50',
    id_invalid: '#F44336',
    amount_invalid: '#FF9800',
    missing_data: '#FF5722',
    excluded: '#9E9E9E'
  }
  const c = colors[value] ?? '#9E9E9E'
  return `<span style="background:${c}22;color:${c};padding:2px 8px;border-radius:4px;font-size:12px">${(value ?? '').replace(/_/g, ' ')}</span>`
}

const maskID = (v: string) => {
  if (!v) return ''
  const s = v.trim()
  if (/^\d{13}$/.test(s)) {
    // RSA ID — reveal first 6 digits and last 2
    return s.substring(0, 6) + '*****' + s.substring(11)
  }
  // Passport — reveal first 2 chars and last 1, mask the rest
  if (s.length <= 3) return '*'.repeat(s.length)
  return (
    s.substring(0, 2) + '*'.repeat(s.length - 3) + s.substring(s.length - 1)
  )
}

// AG Grid column defs for the register diff panels
const diffNewJoinerColDefs = [
  { headerName: 'Row', field: 'row_number', width: 70 },
  { headerName: 'Member Name', field: 'member_name', flex: 1, minWidth: 130 },
  { headerName: 'Emp No', field: 'employee_number', width: 110 },
  {
    headerName: 'ID Number',
    field: 'id_number',
    width: 140,
    valueFormatter: (p: any) => maskID(p.value)
  }
]

const diffExitColDefs = [
  { headerName: 'Row', field: 'row_number', width: 70 },
  { headerName: 'Member Name', field: 'member_name', flex: 1, minWidth: 130 },
  { headerName: 'Emp No', field: 'employee_number', width: 110 },
  {
    headerName: 'ID Number',
    field: 'id_number',
    width: 140,
    valueFormatter: (p: any) => maskID(p.value)
  },
  { headerName: 'Exit Date', field: 'exit_date', width: 120 }
]

const diffAmendmentColDefs = [
  { headerName: 'Row', field: 'row_number', width: 70 },
  { headerName: 'Member Name', field: 'member_name', flex: 1, minWidth: 130 },
  { headerName: 'Emp No', field: 'employee_number', width: 110 },
  {
    headerName: 'Changed Fields',
    field: 'changed_fields',
    flex: 1,
    minWidth: 220,
    valueFormatter: (p: any) => {
      if (!p.value) return ''
      try {
        const obj = typeof p.value === 'string' ? JSON.parse(p.value) : p.value
        return Object.entries(obj)
          .map(([k, v]: any) => `${k}: ${v[0]} → ${v[1]}`)
          .join('; ')
      } catch {
        return JSON.stringify(p.value)
      }
    }
  }
]

const invalidRecordColDefs = [
  { headerName: 'Row', field: 'row_number', width: 65 },
  { headerName: 'Member Name', field: 'member_name', flex: 1, minWidth: 130 },
  { headerName: 'Emp No', field: 'employee_number', width: 100 },
  {
    headerName: 'ID / Passport',
    field: 'id_number',
    width: 140,
    valueFormatter: (p: any) => maskID(p.value ?? '')
  },
  {
    headerName: 'Issue',
    field: 'validation_status',
    width: 120,
    cellRenderer: (p: any) => validationBadge(p.value)
  },
  { headerName: 'Reason', field: 'exclusion_reason', flex: 1, minWidth: 180 }
]

const downloadInvalidRecords = () => {
  const rows = invalidRecords.value
  if (!rows.length) return
  const headers = [
    'row_number',
    'member_name',
    'employee_number',
    'id_number',
    'validation_status',
    'exclusion_reason'
  ]
  const lines = [
    headers.join(','),
    ...rows.map((r: any) =>
      headers.map((h) => JSON.stringify(r[h] ?? '')).join(',')
    )
  ]
  const blob = new Blob([lines.join('\n')], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `invalid_records_submission_${props.submissionId}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

const recordColumnDefs = [
  { headerName: 'Row', field: 'row_number', width: 70 },
  { headerName: 'Member Name', field: 'member_name', flex: 1, minWidth: 130 },
  { headerName: 'Emp No', field: 'employee_number', width: 100 },
  {
    headerName: 'ID Number',
    field: 'id_number',
    width: 140,
    valueFormatter: (p: any) => maskID(p.value ?? '')
  },
  {
    headerName: 'ID Type',
    field: 'id_type',
    width: 90,
    valueFormatter: (p: any) =>
      p.value === 'rsa_id' ? 'RSA ID' : p.value === 'passport' ? 'Passport' : ''
  },
  {
    headerName: 'Salary',
    field: 'salary',
    width: 110,
    type: 'numericColumn',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  { headerName: 'Benefit', field: 'benefit_code', width: 100 },
  {
    headerName: 'Premium',
    field: 'premium_amount',
    width: 110,
    type: 'numericColumn',
    valueFormatter: (p: any) => fmtCurrency(p.value)
  },
  { headerName: 'Entry Date', field: 'entry_date', width: 110 },
  { headerName: 'Exit Date', field: 'exit_date', width: 110 },
  {
    headerName: 'Status',
    field: 'validation_status',
    width: 130,
    cellRenderer: (p: any) => validationBadge(p.value)
  }
]

const loadSubmission = async () => {
  loading.value = true
  try {
    const res = await GroupPricingService.getEmployerSubmission(
      props.submissionId
    )
    submission.value = res.data?.data ?? null
    console.log('Loaded submission:', submission.value)
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to load submission',
      color: 'error'
    }
  } finally {
    loading.value = false
  }
}

const handleUpload = async () => {
  if (!uploadFile.value) return
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', uploadFile.value)
    await GroupPricingService.uploadEmployerSubmission(props.submissionId, fd)
    uploadFile.value = null
    snackbar.value = {
      show: true,
      message: 'File uploaded and parsed successfully',
      color: 'success'
    }
    await loadSubmission()
    await loadRegisterDiff()
  } catch {
    snackbar.value = { show: true, message: 'Upload failed', color: 'error' }
  } finally {
    uploading.value = false
  }
}

const handleReview = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.reviewEmployerSubmission(props.submissionId, {})
    snackbar.value = {
      show: true,
      message: 'Submission moved to Under Review',
      color: 'success'
    }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

const handleRaiseQuery = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.raiseSubmissionQuery(props.submissionId, {
      query_notes: queryNotes.value
    })
    raiseQueryDialog.value = false
    queryNotes.value = ''
    snackbar.value = { show: true, message: 'Query raised', color: 'warning' }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

const handleAccept = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.acceptEmployerSubmission(props.submissionId, {})
    confirmAcceptDialog.value = false
    snackbar.value = {
      show: true,
      message: 'Submission accepted',
      color: 'success'
    }
    await loadSubmission()
    // Reload the register diff so the UI immediately shows the persisted snapshot
    // (is_snapshot: true, teal chip, provenance banner) that was auto-created on acceptance.
    await loadRegisterDiff()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

const handleReject = async () => {
  actionLoading.value = true
  try {
    await GroupPricingService.rejectEmployerSubmission(props.submissionId, {
      reason: rejectionReason.value
    })
    rejectDialog.value = false
    rejectionReason.value = ''
    snackbar.value = {
      show: true,
      message: 'Submission rejected',
      color: 'error'
    }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Action failed',
      color: 'error'
    }
  } finally {
    actionLoading.value = false
  }
}

const handleGenerateSchedule = async () => {
  generatingSchedule.value = true
  try {
    await GroupPricingService.generateScheduleFromSubmission(props.submissionId)
    snackbar.value = {
      show: true,
      message: 'Premium schedule generated',
      color: 'success'
    }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Generation failed',
      color: 'error'
    }
  } finally {
    generatingSchedule.value = false
  }
}

const viewSchedule = () => {
  if (submission.value?.linked_premium_schedule_id) {
    router.push({
      name: 'group-pricing-premium-schedule-detail',
      params: { schedule_id: submission.value.linked_premium_schedule_id }
    })
  }
}

const loadRegisterDiff = async () => {
  if (!submission.value || submission.value.record_count === 0) return
  registerDiffLoading.value = true
  try {
    const res = await GroupPricingService.computeRegisterDiff(
      props.submissionId
    )
    registerDiff.value = res.data?.data ?? null
  } catch {
    // non-fatal
  } finally {
    registerDiffLoading.value = false
  }
}

const handleReSnapshot = async () => {
  registerDiffLoading.value = true
  try {
    const res = await GroupPricingService.snapshotRegisterDiff(
      props.submissionId
    )
    registerDiff.value = res.data?.data ?? null
    snackbar.value = {
      show: true,
      message: 'Register diff snapshot updated',
      color: 'success'
    }
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Snapshot failed',
      color: 'error'
    }
  } finally {
    registerDiffLoading.value = false
  }
}

const formatSnapshotDate = (iso: string | null | undefined): string => {
  if (!iso) return ''
  const d = new Date(iso)
  return d.toLocaleDateString(undefined, {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

const handleApplyExits = async () => {
  exitConfirmDialog.value = false
  applyingExits.value = true
  try {
    const res = await GroupPricingService.applySubmissionExits(
      props.submissionId
    )
    exitResult.value = res.data?.data ?? null
    snackbar.value = {
      show: true,
      message: 'Exits applied to member register',
      color: 'success'
    }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Apply exits failed',
      color: 'error'
    }
  } finally {
    applyingExits.value = false
  }
}

const handleApplyAmendments = async () => {
  amendmentConfirmDialog.value = false
  applyingAmendments.value = true
  try {
    const res = await GroupPricingService.applySubmissionAmendments(
      props.submissionId
    )
    amendmentResult.value = res.data?.data ?? null
    snackbar.value = {
      show: true,
      message: 'Amendments applied to member register',
      color: 'success'
    }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Apply amendments failed',
      color: 'error'
    }
  } finally {
    applyingAmendments.value = false
  }
}

const loadNewJoinerDetails = async () => {
  try {
    const res = await GroupPricingService.getNewJoinerDetails(
      props.submissionId
    )
    newJoinerDetails.value = res.data?.data ?? []
  } catch {
    // non-fatal
  }
}

const handleUploadNewJoinerDetails = async () => {
  if (!newJoinerDetailFile.value) return
  uploadingNewJoinerDetails.value = true
  try {
    const fd = new FormData()
    fd.append('file', newJoinerDetailFile.value)
    const res = await GroupPricingService.uploadNewJoinerDetails(
      props.submissionId,
      fd
    )
    newJoinerDetails.value = res.data?.data ?? []
    newJoinerDetailFile.value = null
    snackbar.value = {
      show: true,
      message: `${newJoinerDetails.value.length} new joiner records staged`,
      color: 'success'
    }
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Upload failed',
      color: 'error'
    }
  } finally {
    uploadingNewJoinerDetails.value = false
  }
}

const handleSyncNewJoiners = async () => {
  newJoinerSyncConfirmDialog.value = false
  syncingNewJoiners.value = true
  try {
    const res = await GroupPricingService.syncNewJoiners(props.submissionId)
    newJoinerSyncResult.value = res.data?.data ?? null
    snackbar.value = {
      show: true,
      message: 'New joiners enrolled in member register',
      color: 'success'
    }
    await loadSubmission()
  } catch (e: any) {
    snackbar.value = {
      show: true,
      message: e?.response?.data?.message ?? 'Sync new joiners failed',
      color: 'error'
    }
  } finally {
    syncingNewJoiners.value = false
  }
}

onMounted(async () => {
  await loadSubmission()
  if (submission.value && submission.value.record_count > 0) {
    await Promise.all([loadRegisterDiff(), loadNewJoinerDetails()])
  }
})
</script>
