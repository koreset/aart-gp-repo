<template>
  <div>
    <!-- Header actions -->
    <div class="d-flex justify-space-between align-center mb-4 flex-wrap gap-2">
      <div>
        <div class="text-subtitle-1 font-weight-medium">Payment Schedules</div>
        <div class="text-caption text-medium-emphasis" style="max-width: 420px">
          Select approved claims, generate a payment schedule, and export it.
          Generate ACB files for BankServ processing and reconcile bank
          responses.
        </div>
      </div>
      <div class="d-flex gap-3">
        <v-btn
          v-if="hasPermission('claims_pay:manage_bank_profiles')"
          rounded
          size="small"
          color="teal"
          variant="outlined"
          prepend-icon="mdi-bank"
          class="mr-3"
          @click="openBankProfilesDialog"
        >
          Bank Profiles
        </v-btn>
        <v-btn
          rounded
          size="small"
          color="primary"
          prepend-icon="mdi-plus"
          @click="openCreateDialog"
        >
          New Payment Schedule
        </v-btn>
      </div>
    </div>

    <!-- Schedules list -->
    <v-progress-linear
      v-if="loading"
      indeterminate
      color="primary"
      class="mb-2"
    />

    <v-alert
      v-if="!loading && schedules.length === 0"
      type="info"
      variant="tonal"
      class="mb-4"
    >
      No payment schedules yet. Create one by selecting approved claims above.
    </v-alert>

    <template v-if="schedules.length > 0">
      <!-- Horizontal scrolling card strip -->
      <div class="schedule-card-strip mb-4">
        <div class="schedule-card-strip__inner">
          <div
            v-for="schedule in schedules"
            :key="'card-' + schedule.id"
            class="schedule-card-strip__item"
          >
            <v-card
              rounded="lg"
              border
              :color="
                selectedScheduleId === schedule.id
                  ? 'primary'
                  : scheduleCardColor(schedule.status)
              "
              :variant="
                selectedScheduleId === schedule.id ? 'outlined' : 'tonal'
              "
              class="h-100 schedule-card-strip__card"
              @click="selectScheduleCard(schedule)"
            >
              <v-card-title
                class="text-subtitle-2 d-flex justify-space-between align-center pt-2 px-3 pb-0"
              >
                <span>{{ schedule.schedule_number }}</span>
                <div class="d-flex gap-1">
                  <v-chip
                    v-if="schedule.acb_file_generated"
                    color="teal"
                    size="x-small"
                    label
                    variant="flat"
                  >
                    ACB
                  </v-chip>
                  <v-chip
                    :color="statusColor(schedule.status)"
                    size="x-small"
                    label
                  >
                    {{ statusLabel(schedule.status) }}
                  </v-chip>
                </div>
              </v-card-title>
              <v-card-text class="px-3 py-1">
                <div class="d-flex justify-space-between">
                  <div>
                    <div class="text-caption text-medium-emphasis">Claims</div>
                    <div class="text-body-2 font-weight-medium">{{
                      schedule.claims_count
                    }}</div>
                  </div>
                  <div class="text-right">
                    <div class="text-caption text-medium-emphasis">Total</div>
                    <div class="text-body-2 font-weight-medium">{{
                      formatCurrency(schedule.total_amount)
                    }}</div>
                  </div>
                </div>
                <div class="text-caption text-medium-emphasis mt-1">
                  {{ formatDate(schedule.created_at) }}
                  <span v-if="schedule.created_by">
                    &middot; {{ schedule.created_by }}</span
                  >
                </div>
                <div
                  v-if="itemsMissingBanking(schedule).length > 0"
                  class="mt-1"
                >
                  <v-chip color="orange" size="x-small" variant="tonal">
                    {{ itemsMissingBanking(schedule).length }} missing bank
                  </v-chip>
                </div>
              </v-card-text>
            </v-card>
          </div>
        </div>
      </div>

      <!-- Data Grid list -->
      <DataGrid
        :row-data="schedules"
        :column-defs="scheduleColumnDefs"
        density="compact"
        :pagination="true"
        :pagination-page-size="20"
        @row-double-clicked="onScheduleRowClicked"
      />
    </template>

    <!-- ── Create Schedule Dialog ── -->
    <v-dialog v-model="createDialog" persistent max-width="800px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >New Payment Schedule</v-card-title
        >
        <v-card-text>
          <v-text-field
            v-model="newScheduleDescription"
            label="Description (optional)"
            variant="outlined"
            density="compact"
            class="mb-4"
          />

          <!-- Filter approved claims -->
          <div class="d-flex align-center justify-space-between mb-2">
            <span class="text-subtitle-2">Select Approved Claims</span>
            <v-chip size="small" color="primary" variant="tonal">
              {{ selectedClaimIDs.length }} selected
            </v-chip>
          </div>
          <v-row dense class="mb-2">
            <v-col cols="6">
              <v-text-field
                v-model="claimFilter"
                label="Filter by Claim Number"
                prepend-inner-icon="mdi-magnify"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="benefitFilter"
                label="Filter by Benefit"
                prepend-inner-icon="mdi-filter-outline"
                variant="outlined"
                density="compact"
                clearable
              />
            </v-col>
          </v-row>
          <DataGrid
            :row-data="filteredApprovedClaims"
            :column-defs="claimsColumnDefs"
            row-selection="multiple"
            density="compact"
            :pagination="false"
            @row-selection-changed="onClaimSelectionChanged"
          />
          <div class="text-subtitle-2 text-right mt-2">
            Subtotal: <strong>{{ formatCurrency(selectedTotal) }}</strong>
          </div>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="createDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="creating"
            :disabled="selectedClaimIDs.length === 0"
            @click="createSchedule"
          >
            Create Schedule
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── View Schedule Detail Dialog ── -->
    <v-dialog v-model="viewDialog" max-width="1000px">
      <v-card v-if="activeSchedule" rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2 d-flex align-center gap-2">
          {{ activeSchedule.schedule_number }}
          <v-chip
            :color="statusColor(activeSchedule.status)"
            size="small"
            label
          >
            {{ statusLabel(activeSchedule.status) }}
          </v-chip>
          <v-chip
            v-if="activeSchedule.acb_file_generated"
            color="teal"
            size="x-small"
            label
          >
            ACB Generated
          </v-chip>
        </v-card-title>
        <v-card-text>
          <v-row dense class="mb-3">
            <v-col cols="6" md="3">
              <div class="text-caption text-medium-emphasis">Claims</div>
              <div class="font-weight-medium">{{
                activeSchedule.claims_count
              }}</div>
            </v-col>
            <v-col cols="6" md="3">
              <div class="text-caption text-medium-emphasis">Total</div>
              <div class="font-weight-medium">{{
                formatCurrency(activeSchedule.total_amount)
              }}</div>
            </v-col>
            <v-col cols="6" md="3">
              <div class="text-caption text-medium-emphasis">Created By</div>
              <div class="font-weight-medium">{{
                activeSchedule.created_by
              }}</div>
            </v-col>
            <v-col cols="6" md="3">
              <div class="text-caption text-medium-emphasis">Created</div>
              <div class="font-weight-medium">{{
                formatDate(activeSchedule.created_at)
              }}</div>
            </v-col>
          </v-row>

          <v-tabs v-model="viewTab" density="compact" class="mb-3">
            <v-tab value="claims">Claims</v-tab>
            <v-tab value="acb">ACB Files</v-tab>
            <v-tab value="reconciliation">Reconciliation</v-tab>
            <v-tab value="proofs">Proof of Payment</v-tab>
          </v-tabs>

          <v-tabs-window v-model="viewTab">
            <!-- Claims Tab -->
            <v-tabs-window-item value="claims">
              <v-table density="compact" class="border rounded mb-4">
                <thead>
                  <tr>
                    <th>Claim #</th>
                    <th>Member</th>
                    <th>ID Number</th>
                    <th>Scheme</th>
                    <th>Benefit</th>
                    <th>Bank</th>
                    <th class="text-right">Amount</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in activeSchedule.items" :key="item.id">
                    <td>{{ item.claim_number }}</td>
                    <td>{{ item.member_name }}</td>
                    <td>{{ item.member_id_number }}</td>
                    <td>{{ item.scheme_name }}</td>
                    <td>{{ item.benefit_name }}</td>
                    <td>
                      <v-chip
                        v-if="item.bank_account_number"
                        size="x-small"
                        color="teal"
                        variant="tonal"
                      >
                        {{ item.bank_name || 'Set' }}
                      </v-chip>
                      <v-chip
                        v-else
                        size="x-small"
                        color="orange"
                        variant="tonal"
                      >
                        Missing
                      </v-chip>
                    </td>
                    <td class="text-right">{{
                      formatCurrency(item.claim_amount)
                    }}</td>
                  </tr>
                </tbody>
              </v-table>
            </v-tabs-window-item>

            <!-- ACB Files Tab -->
            <v-tabs-window-item value="acb">
              <v-progress-linear
                v-if="loadingACBFiles"
                indeterminate
                color="teal"
                class="mb-2"
              />
              <v-alert
                v-if="!loadingACBFiles && acbFiles.length === 0"
                type="info"
                variant="tonal"
                density="compact"
                class="mb-3"
              >
                No ACB files generated yet for this schedule.
              </v-alert>
              <v-list v-else density="compact" class="border rounded mb-3">
                <v-list-item
                  v-for="acb in acbFiles"
                  :key="acb.id"
                  :subtitle="`Generated by ${acb.generated_by} on ${formatDate(acb.generated_at)} | ${acb.transaction_count} transactions | ${formatCurrency(acb.total_amount)}`"
                  :title="acb.file_name"
                >
                  <template #prepend>
                    <v-icon
                      :color="acb.status === 'reconciled' ? 'success' : 'teal'"
                    >
                      {{
                        acb.status === 'reconciled'
                          ? 'mdi-check-circle'
                          : 'mdi-file-document'
                      }}
                    </v-icon>
                  </template>
                  <template #append>
                    <div class="d-flex gap-1">
                      <v-chip
                        v-if="acb.is_retry"
                        size="x-small"
                        color="orange"
                        variant="tonal"
                        class="mr-1"
                      >
                        Retry
                      </v-chip>
                      <v-chip
                        :color="
                          acb.status === 'reconciled' ? 'success' : 'grey'
                        "
                        size="x-small"
                        label
                      >
                        {{ acb.status }}
                      </v-chip>
                      <v-btn
                        size="x-small"
                        variant="text"
                        color="primary"
                        icon="mdi-download"
                        :loading="downloadingACB === acb.id"
                        @click="downloadACBFile(acb)"
                      />
                    </div>
                  </template>
                </v-list-item>
              </v-list>
            </v-tabs-window-item>

            <!-- Reconciliation Tab -->
            <v-tabs-window-item value="reconciliation">
              <v-progress-linear
                v-if="loadingRecon"
                indeterminate
                color="deep-purple"
                class="mb-2"
              />

              <!-- Reconciliation summary chips -->
              <div v-if="reconSummary" class="d-flex gap-2 mb-3 flex-wrap">
                <v-chip color="default" variant="tonal">
                  Total: {{ reconSummary.total_transactions }}
                </v-chip>
                <v-chip color="success" variant="tonal">
                  Paid: {{ reconSummary.paid }} ({{
                    formatCurrency(reconSummary.total_paid)
                  }})
                </v-chip>
                <v-chip color="error" variant="tonal">
                  Failed: {{ reconSummary.failed }} ({{
                    formatCurrency(reconSummary.total_failed)
                  }})
                </v-chip>
                <v-chip color="orange" variant="tonal">
                  Unmatched: {{ reconSummary.unmatched }}
                </v-chip>
              </div>

              <!-- Reconciliation results grid -->
              <v-alert
                v-if="!loadingRecon && reconResults.length === 0"
                type="info"
                variant="tonal"
                density="compact"
              >
                No reconciliation data yet. Upload a bank response to reconcile.
              </v-alert>
              <v-table v-else density="compact" class="border rounded mb-3">
                <thead>
                  <tr>
                    <th>Claim #</th>
                    <th>Account</th>
                    <th class="text-right">Amount</th>
                    <th>Status</th>
                    <th>Failure Reason</th>
                    <th>Bank Ref</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="r in reconResults" :key="r.id">
                    <td>{{ r.claim_number || '—' }}</td>
                    <td>{{ r.account_number }}</td>
                    <td class="text-right">{{ formatCurrency(r.amount) }}</td>
                    <td>
                      <v-chip
                        :color="reconStatusColor(r.status)"
                        size="x-small"
                        label
                        variant="flat"
                      >
                        {{ r.status }}
                      </v-chip>
                    </td>
                    <td class="text-caption">{{ r.failure_reason || '—' }}</td>
                    <td class="text-caption">{{ r.bank_reference || '—' }}</td>
                  </tr>
                </tbody>
              </v-table>

              <!-- Retry failed button -->
              <v-btn
                v-if="reconResults.some((r: any) => r.status === 'failed')"
                color="orange"
                size="small"
                variant="outlined"
                prepend-icon="mdi-refresh"
                :loading="retrying"
                @click="retryFailed"
              >
                Retry Failed Payments
              </v-btn>
            </v-tabs-window-item>

            <!-- Proofs Tab -->
            <v-tabs-window-item value="proofs">
              <v-alert
                v-if="
                  !activeSchedule.proof_of_payments ||
                  activeSchedule.proof_of_payments.length === 0
                "
                type="warning"
                variant="tonal"
                density="compact"
              >
                No proof of payment uploaded yet.
              </v-alert>
              <v-list v-else density="compact" class="border rounded">
                <v-list-item
                  v-for="proof in activeSchedule.proof_of_payments"
                  :key="proof.id"
                  :subtitle="`Uploaded by ${proof.uploaded_by} on ${formatDate(proof.uploaded_at)}`"
                  :title="proof.file_name"
                >
                  <template #append>
                    <v-btn
                      size="x-small"
                      variant="text"
                      color="primary"
                      icon="mdi-download"
                      :loading="downloadingProof === proof.id"
                      @click="downloadProof(proof)"
                    />
                  </template>
                </v-list-item>
              </v-list>
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0 flex-wrap gap-1">
          <v-btn
            variant="outlined"
            size="small"
            color="info"
            prepend-icon="mdi-download"
            :loading="exporting === activeSchedule.id"
            @click="exportSchedule(activeSchedule)"
          >
            Export CSV
          </v-btn>
          <v-btn
            v-if="hasPermission('claims_pay:generate_acb')"
            variant="outlined"
            size="small"
            color="teal"
            prepend-icon="mdi-file-document-outline"
            :disabled="itemsMissingBanking(activeSchedule).length > 0"
            @click="openACBDialog(activeSchedule)"
          >
            Generate ACB
          </v-btn>
          <v-btn
            v-if="
              activeSchedule.acb_file_generated &&
              hasPermission('claims_pay:upload_response')
            "
            variant="outlined"
            size="small"
            color="deep-purple"
            prepend-icon="mdi-file-upload"
            @click="openResponseDialog(activeSchedule)"
          >
            Upload Response
          </v-btn>
          <v-btn
            v-if="activeSchedule.status !== 'confirmed'"
            variant="outlined"
            size="small"
            color="success"
            prepend-icon="mdi-upload"
            @click="openProofDialog(activeSchedule)"
          >
            Upload Proof
          </v-btn>
          <v-spacer />
          <v-btn variant="text" size="small" @click="viewDialog = false"
            >Close</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── ACB Generate Dialog ── -->
    <v-dialog v-model="acbDialog" persistent max-width="500px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">Generate ACB File</v-card-title>
        <v-card-text>
          <p class="text-body-2 mb-3">
            Generate a BankServ ACB file for schedule
            <strong>{{ acbTargetSchedule?.schedule_number }}</strong>
            ({{ acbTargetSchedule?.claims_count }} claims,
            {{ formatCurrency(acbTargetSchedule?.total_amount ?? 0) }}).
          </p>
          <v-select
            v-model="acbProfileId"
            :items="bankProfiles"
            item-title="profile_name"
            item-value="id"
            label="Bank Profile *"
            variant="outlined"
            density="compact"
            class="mb-3"
            :loading="loadingProfiles"
          />
          <v-text-field
            v-model="acbActionDate"
            label="Action Date *"
            type="date"
            variant="outlined"
            density="compact"
            hint="Date the bank should process the payments"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="acbDialog = false">Cancel</v-btn>
          <v-btn
            color="teal"
            :loading="generatingACB"
            :disabled="!acbProfileId || !acbActionDate"
            @click="generateACB"
          >
            Generate ACB
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Bank Response Upload Dialog ── -->
    <v-dialog v-model="responseDialog" persistent max-width="500px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Upload Bank Response</v-card-title
        >
        <v-card-text>
          <p class="text-body-2 mb-3">
            Upload the bank response file for schedule
            <strong>{{ responseTargetSchedule?.schedule_number }}</strong> to
            reconcile payments.
          </p>
          <v-select
            v-model="responseACBFileId"
            :items="responseACBFiles"
            item-title="file_name"
            item-value="id"
            label="ACB File *"
            variant="outlined"
            density="compact"
            class="mb-3"
            :loading="loadingResponseACBFiles"
          />
          <v-file-input
            v-model="responseFile"
            label="Bank Response File"
            prepend-icon="mdi-file-upload"
            variant="outlined"
            density="compact"
            accept=".txt,.csv"
            hint="ACB response (.txt) or CSV format (.csv)"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="responseDialog = false">Cancel</v-btn>
          <v-btn
            color="deep-purple"
            :loading="processingResponse"
            :disabled="!responseACBFileId || !responseFile"
            @click="processResponse"
          >
            Reconcile
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Upload Proof Dialog ── -->
    <v-dialog v-model="proofDialog" persistent max-width="500px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2"
          >Upload Proof of Payment</v-card-title
        >
        <v-card-text>
          <p class="text-body-2 mb-3">
            Upload the proof-of-payment document for schedule
            <strong>{{ proofTargetSchedule?.schedule_number }}</strong
            >. Once uploaded, all claims in this schedule will be marked as
            <strong>Paid</strong> and the schedule will be confirmed.
          </p>
          <v-file-input
            v-model="proofFile"
            label="Proof of Payment Document"
            prepend-icon="mdi-file-upload"
            variant="outlined"
            density="compact"
            accept=".pdf,.csv,.xlsx,.xls,.png,.jpg,.jpeg"
            class="mb-3"
          />
          <v-textarea
            v-model="proofNotes"
            label="Notes (optional)"
            variant="outlined"
            density="compact"
            rows="3"
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="proofDialog = false">Cancel</v-btn>
          <v-btn
            color="success"
            :loading="uploadingProof"
            :disabled="!proofFile"
            @click="uploadProof"
          >
            Confirm Payment
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Bank Profiles Management Dialog ── -->
    <v-dialog v-model="profilesDialog" max-width="700px">
      <v-card rounded="lg">
        <v-card-title
          class="text-h6 pa-4 pb-2 d-flex justify-space-between align-center"
        >
          Bank Profiles
          <v-btn
            size="small"
            color="primary"
            prepend-icon="mdi-plus"
            @click="openCreateProfileDialog"
          >
            New Profile
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-progress-linear
            v-if="loadingProfiles"
            indeterminate
            color="teal"
            class="mb-2"
          />
          <v-alert
            v-if="!loadingProfiles && bankProfiles.length === 0"
            type="info"
            variant="tonal"
            density="compact"
          >
            No bank profiles configured yet.
          </v-alert>
          <v-list v-else density="compact" class="border rounded">
            <v-list-item
              v-for="profile in bankProfiles"
              :key="profile.id"
              :subtitle="`${profile.bank_name} | Account: ${profile.user_account_number} | Gen #${profile.generation_number}`"
              :title="profile.profile_name"
            >
              <template #prepend>
                <v-icon color="teal">mdi-bank</v-icon>
              </template>
              <template #append>
                <v-chip
                  :color="profile.is_active ? 'success' : 'grey'"
                  size="x-small"
                  label
                  class="mr-2"
                >
                  {{ profile.is_active ? 'Active' : 'Inactive' }}
                </v-chip>
                <v-btn
                  size="x-small"
                  variant="text"
                  color="primary"
                  icon="mdi-pencil"
                  @click="editProfile(profile)"
                />
                <v-btn
                  size="x-small"
                  variant="text"
                  color="error"
                  icon="mdi-delete"
                  @click="deleteProfile(profile)"
                />
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="profilesDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- ── Create/Edit Profile Dialog ── -->
    <v-dialog v-model="profileFormDialog" persistent max-width="550px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">
          {{ editingProfile ? 'Edit Bank Profile' : 'New Bank Profile' }}
        </v-card-title>
        <v-card-text>
          <v-text-field
            v-model="profileForm.profile_name"
            label="Profile Name *"
            variant="outlined"
            density="compact"
            class="mb-2"
          />
          <v-select
            v-model="profileForm.bank_name"
            :items="bankNameOptions"
            label="Bank *"
            variant="outlined"
            density="compact"
            class="mb-2"
            @update:model-value="onProfileBankSelected"
          />
          <v-text-field
            v-model="profileForm.user_code"
            label="BankServ User Code *"
            variant="outlined"
            density="compact"
            class="mb-2"
            hint="4-character code assigned by the bank"
            persistent-hint
          />
          <v-row dense>
            <v-col cols="6">
              <v-text-field
                v-model="profileForm.user_branch_code"
                label="Source Branch Code"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-text-field
                v-model="profileForm.user_account_number"
                label="Source Account Number *"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
          <v-row dense>
            <v-col cols="6">
              <v-select
                v-model="profileForm.user_account_type"
                :items="accountTypeOptions"
                label="Account Type"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="6">
              <v-select
                v-model="profileForm.service_type"
                :items="serviceTypeOptions"
                label="Service Type"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
          <v-text-field
            v-model="profileForm.bank_type_code"
            label="Bank Type Code"
            variant="outlined"
            density="compact"
            hint="Default: 04 (standard)"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="profileFormDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            :loading="savingProfile"
            :disabled="
              !profileForm.profile_name ||
              !profileForm.bank_name ||
              !profileForm.user_code
            "
            @click="saveProfile"
          >
            {{ editingProfile ? 'Update' : 'Create' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar -->
    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="4000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn variant="text" color="white" @click="snackbar = false"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import type { ColDef } from 'ag-grid-community'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

interface ScheduleItem {
  id: number
  claim_id: number
  claim_number: string
  member_name: string
  member_id_number: string
  benefit_name: string
  scheme_name: string
  scheme_id: number
  claim_amount: number
  bank_name?: string
  bank_branch_code?: string
  bank_account_number?: string
  bank_account_type?: string
  account_holder_name?: string
}

interface PaymentProof {
  id: number
  schedule_id: number
  file_name: string
  content_type: string
  size_bytes: number
  notes: string
  uploaded_by: string
  uploaded_at: string
}

interface PaymentSchedule {
  id: number
  schedule_number: string
  description: string
  status: string
  total_amount: number
  claims_count: number
  exported_at?: string
  exported_by?: string
  acb_file_generated?: boolean
  acb_generated_at?: string
  acb_generated_by?: string
  created_by: string
  created_at: string
  items: ScheduleItem[]
  proof_of_payments: PaymentProof[]
}

interface Claim {
  id: number
  claim_number: string
  member_name: string
  member_id_number: string
  scheme_name: string
  benefit_alias: string
  claim_amount: number
  status: string
}

interface BankProfile {
  id: number
  profile_name: string
  bank_name: string
  user_code: string
  user_branch_code: string
  user_account_number: string
  user_account_type: string
  bank_type_code: string
  service_type: string
  generation_number: number
  is_active: boolean
}

interface ACBFile {
  id: number
  schedule_id: number
  file_name: string
  action_date: string
  transaction_count: number
  total_amount: number
  status: string
  is_retry: boolean
  generated_by: string
  generated_at: string
}

interface ReconResult {
  id: number
  claim_number: string
  account_number: string
  amount: number
  status: string
  failure_reason: string
  bank_reference: string
}

const { hasPermission } = usePermissionCheck()

// ── State ──────────────────────────────────────────────
const loading = ref(false)
const schedules = ref<PaymentSchedule[]>([])

// Approved claims for selection
const approvedClaims = ref<Claim[]>([])
const loadingApproved = ref(false)
const claimFilter = ref('')
const benefitFilter = ref('')

// Create dialog
const createDialog = ref(false)
const newScheduleDescription = ref('')
const selectedClaimIDs = ref<number[]>([])
const creating = ref(false)

// View dialog
const viewDialog = ref(false)
const activeSchedule = ref<PaymentSchedule | null>(null)
const viewTab = ref('claims')

// Proof upload dialog
const proofDialog = ref(false)
const proofTargetSchedule = ref<PaymentSchedule | null>(null)
const proofFile = ref<File | null>(null)
const proofNotes = ref('')
const uploadingProof = ref(false)

// Export / download
const exporting = ref<number | null>(null)
const downloadingProof = ref<number | null>(null)

// ACB generation
const acbDialog = ref(false)
const acbTargetSchedule = ref<PaymentSchedule | null>(null)
const acbProfileId = ref<number | null>(null)
const acbActionDate = ref('')
const generatingACB = ref(false)

// ACB files in view dialog
const acbFiles = ref<ACBFile[]>([])
const loadingACBFiles = ref(false)
const downloadingACB = ref<number | null>(null)

// Bank response upload
const responseDialog = ref(false)
const responseTargetSchedule = ref<PaymentSchedule | null>(null)
const responseACBFileId = ref<number | null>(null)
const responseACBFiles = ref<ACBFile[]>([])
const loadingResponseACBFiles = ref(false)
const responseFile = ref<File | null>(null)
const processingResponse = ref(false)

// Reconciliation
const reconResults = ref<ReconResult[]>([])
const reconSummary = ref<any>(null)
const loadingRecon = ref(false)
const retrying = ref(false)

// Bank profiles
const profilesDialog = ref(false)
const bankProfiles = ref<BankProfile[]>([])
const loadingProfiles = ref(false)
const profileFormDialog = ref(false)
const editingProfile = ref<BankProfile | null>(null)
const savingProfile = ref(false)
const profileForm = ref({
  profile_name: '',
  bank_name: '',
  user_code: '',
  user_branch_code: '',
  user_account_number: '',
  user_account_type: '1',
  bank_type_code: '04',
  service_type: 'two_day'
})

// Snackbar
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

// ── Constants ─────────────────────────────────────────
const bankNameOptions = [
  'FNB',
  'Standard Bank',
  'ABSA',
  'Nedbank',
  'Capitec',
  'Investec',
  'African Bank',
  'TymeBank',
  'Discovery Bank',
  'Bank Zero'
]

const universalBranchCodes: Record<string, string> = {
  FNB: '250655',
  'Standard Bank': '051001',
  ABSA: '632005',
  Nedbank: '198765',
  Capitec: '470010',
  Investec: '580105',
  'African Bank': '430000',
  TymeBank: '678910',
  'Discovery Bank': '679000',
  'Bank Zero': '888000'
}

const accountTypeOptions = [
  { title: 'Current/Cheque', value: '1' },
  { title: 'Savings', value: '2' },
  { title: 'Transmission', value: '3' }
]

const serviceTypeOptions = [
  { title: 'Same Day', value: 'same_day' },
  { title: 'One Day', value: 'one_day' },
  { title: 'Two Day', value: 'two_day' }
]

// ── Computed ────────────────────────────────────────────
const filteredApprovedClaims = computed(() => {
  let result = approvedClaims.value
  if (claimFilter.value) {
    const q = claimFilter.value.toLowerCase()
    result = result.filter((c) => c.claim_number.toLowerCase().includes(q))
  }
  if (benefitFilter.value) {
    const b = benefitFilter.value.toLowerCase()
    result = result.filter((c) => c.benefit_alias?.toLowerCase().includes(b))
  }
  return result
})

const selectedTotal = computed(() => {
  return approvedClaims.value
    .filter((c) => selectedClaimIDs.value.includes(c.id))
    .reduce((sum, c) => sum + c.claim_amount, 0)
})

const claimsColumnDefs: ColDef<Claim>[] = [
  {
    checkboxSelection: true,
    headerCheckboxSelection: true,
    width: 48,
    minWidth: 48,
    maxWidth: 48,
    pinned: 'left',
    suppressMovable: true,
    resizable: false
  },
  {
    headerName: 'Claim #',
    field: 'claim_number',
    sortable: true,
    minWidth: 120
  },
  { headerName: 'Member', field: 'member_name', sortable: true, minWidth: 140 },
  { headerName: 'ID Number', field: 'member_id_number', minWidth: 120 },
  { headerName: 'Scheme', field: 'scheme_name', sortable: true, minWidth: 140 },
  { headerName: 'Benefit', field: 'benefit_alias', minWidth: 130 },
  {
    headerName: 'Amount',
    field: 'claim_amount',
    minWidth: 120,
    type: 'rightAligned',
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    headerName: 'Status',
    field: 'status',
    minWidth: 100,
    cellRenderer: (p: any) =>
      `<span class="v-chip v-chip--label v-chip--density-comfortable bg-${statusColor(p.value)}">${p.value}</span>`
  }
]

const scheduleColumnDefs: ColDef<PaymentSchedule>[] = [
  {
    headerName: 'Schedule #',
    field: 'schedule_number',
    sortable: true,
    minWidth: 160
  },
  {
    headerName: 'Status',
    field: 'status',
    sortable: true,
    minWidth: 130,
    cellRenderer: (p: any) => {
      const color = statusColor(p.value)
      const label = statusLabel(p.value)
      return `<span class="v-chip v-chip--label v-chip--size-small bg-${color}" style="font-size:11px;padding:0 8px;height:22px;display:inline-flex;align-items:center">${label}</span>`
    }
  },
  {
    headerName: 'Description',
    field: 'description',
    sortable: true,
    minWidth: 180,
    flex: 1
  },
  {
    headerName: 'Claims',
    field: 'claims_count',
    sortable: true,
    minWidth: 80,
    type: 'rightAligned'
  },
  {
    headerName: 'Total Amount',
    field: 'total_amount',
    sortable: true,
    minWidth: 140,
    type: 'rightAligned',
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    headerName: 'Created By',
    field: 'created_by',
    sortable: true,
    minWidth: 120
  },
  {
    headerName: 'Created',
    field: 'created_at',
    sortable: true,
    minWidth: 120,
    valueFormatter: (p) => formatDate(p.value)
  },
  {
    headerName: 'ACB',
    field: 'acb_file_generated',
    minWidth: 70,
    maxWidth: 70,
    cellRenderer: (p: any) =>
      p.value
        ? '<span class="v-chip v-chip--label v-chip--size-x-small bg-teal" style="font-size:10px;padding:0 6px;height:18px;display:inline-flex;align-items:center;color:#fff">ACB</span>'
        : ''
  },
  {
    headerName: 'Actions',
    minWidth: 100,
    maxWidth: 100,
    pinned: 'right',
    cellRenderer: () =>
      '<span style="cursor:pointer;color:#1976d2;font-size:12px;text-decoration:underline">View</span>',
    onCellClicked: (params: any) => {
      if (params.data) {
        selectedScheduleId.value = params.data.id
        viewSchedule(params.data as PaymentSchedule)
      }
    }
  }
]

const selectedScheduleId = ref<number | null>(null)

// ── Helpers ─────────────────────────────────────────────
function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function formatDate(val?: string) {
  if (!val) return '—'
  return new Date(val).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: '2-digit'
  })
}

function statusLabel(status: string) {
  const map: Record<string, string> = {
    submitted: 'Submitted for Payment',
    confirmed: 'Paid / Confirmed',
    draft: 'Draft'
  }
  return map[status] ?? status
}

function statusColor(status: string) {
  const map: Record<string, string> = {
    submitted: 'warning',
    confirmed: 'success',
    draft: 'grey',
    approved: 'info',
    submitted_for_payment: 'warning',
    paid: 'success',
    payment_failed: 'error',
    pending: 'default',
    declined: 'error'
  }
  return map[status] ?? 'default'
}

function reconStatusColor(status: string) {
  const map: Record<string, string> = {
    paid: 'success',
    failed: 'error',
    unmatched: 'orange'
  }
  return map[status] ?? 'default'
}

function scheduleCardColor(status: string) {
  if (status === 'confirmed') return 'success'
  if (status === 'submitted') return 'warning'
  return undefined
}

function itemsMissingBanking(schedule: PaymentSchedule): ScheduleItem[] {
  if (!schedule.items) return []
  return schedule.items.filter(
    (i) => !i.bank_account_number || !i.bank_branch_code
  )
}

function notify(message: string, color: string = 'success') {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

// Unwrap PremiumResponse { success, data } envelope if present
function unwrap(res: any) {
  const body = res?.data
  if (body && typeof body === 'object' && 'success' in body && 'data' in body) {
    return body.data
  }
  return body
}

function downloadBlob(data: any, filename: string, type: string) {
  const blob = new Blob([data], { type })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

// ── Data loading ────────────────────────────────────────
async function loadSchedules() {
  loading.value = true
  try {
    const res = await GroupPricingService.getPaymentSchedules()
    schedules.value = unwrap(res) ?? []
  } catch (e: any) {
    notify('Failed to load payment schedules', 'error')
  } finally {
    loading.value = false
  }
}

async function loadApprovedClaims() {
  loadingApproved.value = true
  try {
    const res = await GroupPricingService.getClaims()
    const all: Claim[] = unwrap(res) ?? []
    approvedClaims.value = all.filter((c) => c.status === 'approved')
  } catch (e: any) {
    notify('Failed to load claims', 'error')
  } finally {
    loadingApproved.value = false
  }
}

async function loadBankProfiles() {
  loadingProfiles.value = true
  try {
    const res = await GroupPricingService.getBankProfiles()
    bankProfiles.value = (unwrap(res) ?? []).filter(
      (p: BankProfile) => p.is_active
    )
  } catch (e: any) {
    notify('Failed to load bank profiles', 'error')
  } finally {
    loadingProfiles.value = false
  }
}

async function loadACBFiles(scheduleId: number) {
  loadingACBFiles.value = true
  try {
    const res = await GroupPricingService.getACBFileRecords(scheduleId)
    acbFiles.value = unwrap(res) ?? []
  } catch {
    acbFiles.value = []
  } finally {
    loadingACBFiles.value = false
  }
}

async function loadReconData(scheduleId: number) {
  loadingRecon.value = true
  try {
    const summaryRes =
      await GroupPricingService.getReconciliationSummary(scheduleId)
    reconSummary.value = unwrap(summaryRes)

    // Load results from all ACB files
    const filesRes = await GroupPricingService.getACBFileRecords(scheduleId)
    const files: ACBFile[] = unwrap(filesRes) ?? []
    const allResults: ReconResult[] = []
    for (const f of files) {
      if (f.status === 'reconciled') {
        const res = await GroupPricingService.getReconciliationResults(f.id)
        allResults.push(...(unwrap(res) ?? []))
      }
    }
    reconResults.value = allResults
  } catch {
    reconResults.value = []
    reconSummary.value = null
  } finally {
    loadingRecon.value = false
  }
}

// ── AG Grid handlers ─────────────────────────────────────
function onClaimSelectionChanged(rows: Claim[]) {
  selectedClaimIDs.value = rows.map((c) => c.id)
}

function selectScheduleCard(schedule: PaymentSchedule) {
  selectedScheduleId.value = schedule.id
  viewSchedule(schedule)
}

function onScheduleRowClicked(row: any) {
  const schedule = row as PaymentSchedule
  if (schedule?.id) {
    selectedScheduleId.value = schedule.id
    viewSchedule(schedule)
  }
}

// ── Actions ─────────────────────────────────────────────
async function openCreateDialog() {
  selectedClaimIDs.value = []
  newScheduleDescription.value = ''
  claimFilter.value = ''
  benefitFilter.value = ''
  await loadApprovedClaims()
  createDialog.value = true
}

async function createSchedule() {
  creating.value = true
  try {
    await GroupPricingService.createPaymentSchedule({
      claim_ids: selectedClaimIDs.value,
      description: newScheduleDescription.value
    })
    createDialog.value = false
    notify('Payment schedule created. Claims moved to "Submitted for Payment".')
    await loadSchedules()
  } catch (e: any) {
    notify(e?.response?.data ?? 'Failed to create schedule', 'error')
  } finally {
    creating.value = false
  }
}

function viewSchedule(schedule: PaymentSchedule) {
  activeSchedule.value = schedule
  viewTab.value = 'claims'
  viewDialog.value = true
}

// Load ACB files and recon data when switching tabs
watch(viewTab, async (tab) => {
  if (!activeSchedule.value) return
  if (tab === 'acb') {
    await loadACBFiles(activeSchedule.value.id)
  } else if (tab === 'reconciliation') {
    await loadReconData(activeSchedule.value.id)
  }
})

async function exportSchedule(schedule: PaymentSchedule) {
  exporting.value = schedule.id
  try {
    const res = await GroupPricingService.exportPaymentScheduleCSV(schedule.id)
    downloadBlob(
      res.data,
      `payment_schedule_${schedule.schedule_number}.csv`,
      'text/csv'
    )
    notify('Payment schedule exported.')
    await loadSchedules()
  } catch (e: any) {
    notify('Failed to export schedule', 'error')
  } finally {
    exporting.value = null
  }
}

// ── ACB Generation ──────────────────────────────────────
async function openACBDialog(schedule: PaymentSchedule) {
  acbTargetSchedule.value = schedule
  acbProfileId.value = null
  acbActionDate.value = new Date(Date.now() + 2 * 86400000)
    .toISOString()
    .split('T')[0]
  await loadBankProfiles()
  acbDialog.value = true
}

async function generateACB() {
  if (!acbTargetSchedule.value || !acbProfileId.value || !acbActionDate.value)
    return
  generatingACB.value = true
  try {
    const res = await GroupPricingService.generateACBFile(
      acbTargetSchedule.value.id,
      {
        bank_profile_id: acbProfileId.value,
        action_date: acbActionDate.value
      }
    )
    acbDialog.value = false
    notify('ACB file generated successfully.')

    // Download immediately
    const acbRecord = unwrap(res)
    const dlRes = await GroupPricingService.downloadACBFile(acbRecord.id)
    downloadBlob(dlRes.data, acbRecord.file_name, 'text/plain')

    await loadSchedules()
  } catch (e: any) {
    notify(
      e?.response?.data?.message ??
        e?.response?.data ??
        'Failed to generate ACB file',
      'error'
    )
  } finally {
    generatingACB.value = false
  }
}

async function downloadACBFile(acb: ACBFile) {
  downloadingACB.value = acb.id
  try {
    const res = await GroupPricingService.downloadACBFile(acb.id)
    downloadBlob(res.data, acb.file_name, 'text/plain')
  } catch {
    notify('Failed to download ACB file', 'error')
  } finally {
    downloadingACB.value = null
  }
}

// ── Bank Response Upload ────────────────────────────────
async function openResponseDialog(schedule: PaymentSchedule) {
  responseTargetSchedule.value = schedule
  responseACBFileId.value = null
  responseFile.value = null
  loadingResponseACBFiles.value = true
  responseDialog.value = true
  try {
    const res = await GroupPricingService.getACBFileRecords(schedule.id)
    responseACBFiles.value = (unwrap(res) ?? []).filter(
      (f: ACBFile) => f.status === 'generated'
    )
  } catch {
    responseACBFiles.value = []
  } finally {
    loadingResponseACBFiles.value = false
  }
}

async function processResponse() {
  if (!responseACBFileId.value || !responseFile.value) return
  processingResponse.value = true
  try {
    const formData = new FormData()
    formData.append('file', responseFile.value as File)
    const res = await GroupPricingService.processBankResponse(
      responseACBFileId.value,
      formData
    )
    responseDialog.value = false

    const s = unwrap(res)
    notify(
      `Reconciliation complete: ${s.paid} paid, ${s.failed} failed, ${s.unmatched} unmatched`,
      s.failed > 0 ? 'warning' : 'success'
    )
    await loadSchedules()
  } catch (e: any) {
    notify(
      e?.response?.data?.message ??
        e?.response?.data ??
        'Failed to process bank response',
      'error'
    )
  } finally {
    processingResponse.value = false
  }
}

// ── Retry Failed ────────────────────────────────────────
async function retryFailed() {
  if (!activeSchedule.value) return
  retrying.value = true
  try {
    // Find the most recent reconciled ACB file
    const filesRes = await GroupPricingService.getACBFileRecords(
      activeSchedule.value.id
    )
    const reconciledFile = (unwrap(filesRes) ?? []).find(
      (f: ACBFile) => f.status === 'reconciled'
    )
    if (!reconciledFile) {
      notify('No reconciled ACB file found to retry from', 'error')
      return
    }
    const res = await GroupPricingService.retryFailedPayments(
      reconciledFile.id,
      { item_ids: [] }
    )
    notify('Retry ACB file generated.')

    // Download immediately
    const retryRecord = unwrap(res)
    const dlRes = await GroupPricingService.downloadACBFile(retryRecord.id)
    downloadBlob(dlRes.data, retryRecord.file_name, 'text/plain')

    await loadReconData(activeSchedule.value.id)
    await loadACBFiles(activeSchedule.value.id)
    await loadSchedules()
  } catch (e: any) {
    notify(e?.response?.data ?? 'Failed to retry failed payments', 'error')
  } finally {
    retrying.value = false
  }
}

// ── Proof of Payment ────────────────────────────────────
function openProofDialog(schedule: PaymentSchedule) {
  proofTargetSchedule.value = schedule
  proofFile.value = null
  proofNotes.value = ''
  proofDialog.value = true
}

async function uploadProof() {
  if (!proofFile.value || !proofTargetSchedule.value) return
  uploadingProof.value = true
  try {
    const formData = new FormData()
    formData.append('file', proofFile.value as File)
    formData.append('notes', proofNotes.value)
    await GroupPricingService.uploadPaymentProof(
      proofTargetSchedule.value.id,
      formData
    )
    proofDialog.value = false
    viewDialog.value = false
    notify('Proof of payment uploaded. All claims marked as Paid.')
    await loadSchedules()
  } catch (e: any) {
    notify(e?.data ?? 'Failed to upload proof of payment', 'error')
  } finally {
    uploadingProof.value = false
  }
}

async function downloadProof(proof: PaymentProof) {
  downloadingProof.value = proof.id
  try {
    const res = await GroupPricingService.downloadPaymentProof(proof.id)
    downloadBlob(
      res.data,
      proof.file_name,
      proof.content_type || 'application/octet-stream'
    )
  } catch (e: any) {
    notify('Failed to download proof', 'error')
  } finally {
    downloadingProof.value = null
  }
}

// ── Bank Profile Management ──────────────────────────────
async function openBankProfilesDialog() {
  await loadBankProfiles()
  profilesDialog.value = true
}

function onProfileBankSelected(bankName: string) {
  profileForm.value.user_branch_code = universalBranchCodes[bankName] || ''
}

function openCreateProfileDialog() {
  editingProfile.value = null
  profileForm.value = {
    profile_name: '',
    bank_name: '',
    user_code: '',
    user_branch_code: '',
    user_account_number: '',
    user_account_type: '1',
    bank_type_code: '04',
    service_type: 'two_day'
  }
  profileFormDialog.value = true
}

function editProfile(profile: BankProfile) {
  editingProfile.value = profile
  profileForm.value = {
    profile_name: profile.profile_name,
    bank_name: profile.bank_name,
    user_code: profile.user_code,
    user_branch_code: profile.user_branch_code,
    user_account_number: profile.user_account_number,
    user_account_type: profile.user_account_type || '1',
    bank_type_code: profile.bank_type_code || '04',
    service_type: profile.service_type || 'two_day'
  }
  profileFormDialog.value = true
}

async function saveProfile() {
  savingProfile.value = true
  try {
    if (editingProfile.value) {
      await GroupPricingService.updateBankProfile(
        editingProfile.value.id,
        profileForm.value
      )
      notify('Bank profile updated.')
    } else {
      await GroupPricingService.createBankProfile(profileForm.value)
      notify('Bank profile created.')
    }
    profileFormDialog.value = false
    await loadBankProfiles()
  } catch (e: any) {
    notify(e?.response?.data ?? 'Failed to save bank profile', 'error')
  } finally {
    savingProfile.value = false
  }
}

async function deleteProfile(profile: BankProfile) {
  try {
    await GroupPricingService.deleteBankProfile(profile.id)
    notify('Bank profile deleted.')
    await loadBankProfiles()
  } catch (e: any) {
    notify('Failed to delete bank profile', 'error')
  }
}

onMounted(loadSchedules)
</script>

<style scoped>
.schedule-card-strip {
  overflow-x: auto;
  overflow-y: hidden;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
}

.schedule-card-strip__inner {
  display: flex;
  gap: 12px;
  padding-bottom: 4px;
}

.schedule-card-strip__item {
  flex: 0 0 320px;
  min-width: 320px;
}

.schedule-card-strip__card {
  cursor: pointer;
  transition: box-shadow 0.15s ease;
}

.schedule-card-strip__card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
</style>
