<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline">Bordereaux Reconciliation</span>
                <p class="text-subtitle-1 text-medium-emphasis mt-2">
                  Match submissions with scheme confirmations and resolve
                  discrepancies
                </p>
              </div>
              <v-btn
                color="grey"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="$router.push('/group-pricing/bordereaux-management')"
              >
                Back to Dashboard
              </v-btn>
            </div>
          </template>
          <template #default>
            <!-- Reconciliation Summary -->
            <v-row class="mb-6">
              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Pending</p
                        >
                        <p class="text-h5 font-weight-bold text-warning">{{
                          reconciliationStats.pending
                        }}</p>
                        <p class="text-caption text-warning">Reconciliation</p>
                      </div>
                      <v-icon size="40" color="warning"
                        >mdi-clock-outline</v-icon
                      >
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Matched</p
                        >
                        <p class="text-h5 font-weight-bold text-success">{{
                          reconciliationStats.matched
                        }}</p>
                        <p class="text-caption text-success">This Week</p>
                      </div>
                      <v-icon size="40" color="success"
                        >mdi-check-circle</v-icon
                      >
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Discrepancies</p
                        >
                        <p class="text-h5 font-weight-bold text-error">{{
                          reconciliationStats.discrepancies
                        }}</p>
                        <p class="text-caption text-error">Require Attention</p>
                      </div>
                      <v-icon size="40" color="error">mdi-alert-circle</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Auto-Match</p
                        >
                        <p class="text-h5 font-weight-bold text-info"
                          >{{ reconciliationStats.autoMatch }}%</p
                        >
                        <p class="text-caption text-info">Success Rate</p>
                      </div>
                      <v-icon size="40" color="info">mdi-robot</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Action Buttons -->
            <v-row class="mb-4">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    <v-icon class="me-2">mdi-cog</v-icon>
                    Reconciliation Actions
                  </v-card-title>
                  <v-card-text>
                    <div class="d-flex flex-wrap gap-3">
                      <v-btn
                        class="mr-2"
                        color="primary"
                        prepend-icon="mdi-robot"
                        :loading="autoReconciling"
                        @click="runAutoReconciliation"
                      >
                        Run Auto-Reconciliation
                      </v-btn>
                      <v-btn
                        class="mr-2"
                        color="info"
                        prepend-icon="mdi-upload"
                        @click="showImportDialog = true"
                      >
                        Import Confirmations
                      </v-btn>
                      <v-btn
                        class="mr-2"
                        color="success"
                        prepend-icon="mdi-refresh"
                        :loading="loading"
                        @click="refreshReconciliation"
                      >
                        Refresh Status
                      </v-btn>
                      <v-btn
                        color="orange"
                        prepend-icon="mdi-file-document-multiple"
                        @click="generateReconciliationReport"
                      >
                        Generate Report
                      </v-btn>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Reconciliation Items -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="d-flex align-center justify-space-between"
                  >
                    <span class="text-h6 font-weight-bold"
                      >Reconciliation Items</span
                    >
                    <v-btn-toggle
                      v-model="viewMode"
                      mandatory
                      variant="outlined"
                      density="compact"
                    >
                      <v-btn value="pending" size="small">
                        <v-icon>mdi-clock</v-icon>
                        Pending
                      </v-btn>
                      <v-btn value="discrepancies" size="small">
                        <v-icon>mdi-alert</v-icon>
                        Discrepancies
                      </v-btn>
                      <v-btn value="matched" size="small">
                        <v-icon>mdi-check</v-icon>
                        Matched
                      </v-btn>
                      <v-btn value="all" size="small">
                        <v-icon>mdi-all-inclusive</v-icon>
                        All
                      </v-btn>
                    </v-btn-toggle>
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="headers"
                      :items="filteredItems"
                      :loading="loading"
                      :items-per-page="15"
                      item-key="id"
                    >
                      <template
                        #[`item.generated_bordereaux_id`]="{
                          item
                        }: {
                          item: any
                        }"
                      >
                        <v-chip color="primary" size="small" variant="tonal">
                          {{ item.generated_bordereaux_id }}
                        </v-chip>
                      </template>

                      <template #[`item.type`]="{ item }: { item: any }">
                        <v-chip
                          :color="getBordereauTypeColor(item.type)"
                          size="small"
                          variant="tonal"
                        >
                          {{ formatBordereauType(item.type) }}
                        </v-chip>
                      </template>

                      <template #[`item.status`]="{ item }: { item: any }">
                        <div class="d-flex align-center">
                          <v-chip
                            :color="getReconciliationStatusColor(item.status)"
                            size="small"
                            class="me-2"
                          >
                            {{ formatReconciliationStatus(item.status) }}
                          </v-chip>
                          <v-btn
                            v-if="item.status === 'discrepancy'"
                            size="x-small"
                            color="error"
                            variant="text"
                            icon="mdi-alert-circle"
                            @click="viewDiscrepancies(item)"
                          />
                        </div>
                      </template>

                      <template #[`item.match_score`]="{ item }: { item: any }">
                        <div class="match-score-container">
                          <v-progress-linear
                            :model-value="item.match_score"
                            :color="getMatchScoreColor(item.match_score)"
                            height="6"
                            rounded
                          />
                          <span class="text-caption"
                            >{{ item.match_score }}%</span
                          >
                        </div>
                      </template>

                      <template
                        #[`item.submitted_amount`]="{ item }: { item: any }"
                      >
                        <span class="font-mono">{{
                          formatCurrency(item.submitted_amount)
                        }}</span>
                      </template>

                      <template
                        #[`item.confirmed_amount`]="{ item }: { item: any }"
                      >
                        <span class="font-mono" :class="getAmountClass(item)">
                          {{
                            item.confirmed_amount
                              ? formatCurrency(item.confirmed_amount)
                              : '-'
                          }}
                        </span>
                      </template>

                      <template #[`item.variance`]="{ item }: { item: any }">
                        <v-chip
                          v-if="item.variance"
                          :color="item.variance === 0 ? 'success' : 'error'"
                          size="small"
                          variant="tonal"
                        >
                          {{ formatCurrency(item.variance) }}
                        </v-chip>
                        <span v-else>-</span>
                      </template>

                      <template
                        #[`item.last_reconciled`]="{ item }: { item: any }"
                      >
                        {{
                          item.last_reconciled
                            ? formatDateTime(item.last_reconciled)
                            : 'Not reconciled'
                        }}
                      </template>

                      <template #[`item.actions`]="{ item }: { item: any }">
                        <div class="d-flex align-center gap-1">
                          <v-btn
                            size="small"
                            class="mr-1"
                            color="info"
                            variant="tonal"
                            icon="mdi-eye"
                            @click="viewReconciliationDetail(item)"
                          />
                          <v-btn
                            v-if="item.status === 'pending'"
                            size="small"
                            class="mr-1"
                            color="primary"
                            variant="tonal"
                            icon="mdi-link"
                            @click="manualMatch(item)"
                          />
                          <v-btn
                            v-if="item.status === 'discrepancy'"
                            size="small"
                            class="mr-1"
                            color="warning"
                            variant="tonal"
                            icon="mdi-wrench"
                            @click="resolveDiscrepancy(item)"
                          />
                          <v-btn
                            v-if="item.status === 'matched'"
                            size="small"
                            class="mr-1"
                            color="success"
                            variant="tonal"
                            icon="mdi-check-circle"
                            @click="confirmReconciliation(item)"
                          />
                          <v-btn
                            size="small"
                            class="mr-1"
                            color="error"
                            variant="tonal"
                            icon="mdi-delete"
                            @click="confirmDelete(item)"
                          />
                          <v-menu>
                            <template #activator="{ props }">
                              <v-btn
                                size="small"
                                color="grey"
                                variant="tonal"
                                icon="mdi-dots-vertical"
                                v-bind="props"
                              />
                            </template>
                            <v-list density="compact">
                              <v-list-item @click="exportReconciliation(item)">
                                <template #prepend>
                                  <v-icon>mdi-export</v-icon>
                                </template>
                                <v-list-item-title>Export</v-list-item-title>
                              </v-list-item>
                              <v-list-item @click="reprocessItem(item)">
                                <template #prepend>
                                  <v-icon>mdi-refresh</v-icon>
                                </template>
                                <v-list-item-title>Reprocess</v-list-item-title>
                              </v-list-item>
                              <v-list-item @click="addNote(item)">
                                <template #prepend>
                                  <v-icon>mdi-note-plus</v-icon>
                                </template>
                                <v-list-item-title>Add Note</v-list-item-title>
                              </v-list-item>
                            </v-list>
                          </v-menu>
                        </div>
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Import Dialog -->
    <v-dialog v-model="showImportDialog" max-width="600">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold bg-info text-white">
          <v-icon class="me-2">mdi-upload</v-icon>
          Import Scheme Confirmations
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="info" variant="tonal" class="mb-4">
            <p class="font-weight-bold">Upload Scheme Confirmation Files</p>
            <p class="mb-2"
              >Import scheme confirmation files to facilitate reconciliation
              matching</p
            >
            <p class="text-caption mb-0">
              <v-icon size="small" class="me-1">mdi-information</v-icon>
              Large files (>1MB) will be automatically compressed using gzip to
              speed up upload
            </p>
          </v-alert>

          <v-row>
            <v-col cols="12">
              <v-select
                v-model="importData.file_type"
                :items="fileTypes"
                label="File Type *"
                variant="outlined"
                density="compact"
                required
              />
            </v-col>
            <v-col cols="12">
              <v-file-input
                v-model="importData.files"
                label="Confirmation Files"
                variant="outlined"
                density="compact"
                multiple
                accept=".xlsx,.csv,.pdf"
                show-size
                prepend-icon="mdi-paperclip"
              />
            </v-col>
            <v-col cols="12">
              <v-checkbox
                v-model="importData.auto_process"
                label="Auto-process after import"
                density="compact"
              />
            </v-col>

            <!-- Compression Progress -->
            <v-col v-if="compressing" cols="12">
              <v-card variant="outlined" color="primary">
                <v-card-text>
                  <div class="d-flex align-center mb-2">
                    <v-icon class="me-2" color="primary">mdi-file-zip</v-icon>
                    <span class="font-weight-bold"
                      >Compressing files for faster upload...</span
                    >
                  </div>
                  <v-progress-linear
                    :model-value="compressionProgress"
                    color="primary"
                    height="8"
                    rounded
                  />
                  <p class="text-caption mt-2 mb-0"
                    >{{ compressionProgress }}% complete</p
                  >
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="showImportDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            color="info"
            variant="flat"
            :loading="importing || compressing"
            :disabled="
              schemesLoading ||
              importing ||
              compressing ||
              !importData.file_type ||
              !importData.files?.length
            "
            @click="performImport"
          >
            <v-icon v-if="compressing" class="me-2">mdi-file-zip</v-icon>
            <v-icon v-else class="me-2">mdi-upload</v-icon>
            {{ compressing ? 'Compressing...' : 'Import Files' }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Reconciliation Detail Dialog -->
    <v-dialog v-model="showDetailDialog" max-width="1000">
      <v-card v-if="selectedItem">
        <v-card-title class="text-h6 font-weight-bold bg-primary text-white">
          <v-icon class="me-2">mdi-file-document-outline</v-icon>
          Reconciliation Details - {{ selectedItem.generated_bordereaux_id }}
        </v-card-title>
        <v-card-text class="pt-6">
          <v-row>
            <v-col cols="12" md="6">
              <h3 class="text-h6 font-weight-bold mb-3">Submission Details</h3>
              <v-list density="compact">
                <v-list-item>
                  <v-list-item-title>Bordereaux Type</v-list-item-title>
                  <v-list-item-subtitle>{{
                    formatBordereauType(selectedItem.type)
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Scheme</v-list-item-title>
                  <v-list-item-subtitle>{{
                    selectedItem.scheme_name
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Scheme</v-list-item-title>
                  <v-list-item-subtitle>{{
                    selectedItem.scheme_name
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Submitted Amount</v-list-item-title>
                  <v-list-item-subtitle>{{
                    formatCurrency(selectedItem.submitted_amount)
                  }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Record Count</v-list-item-title>
                  <v-list-item-subtitle>{{
                    selectedItem.record_count?.toLocaleString()
                  }}</v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-col>
            <v-col cols="12" md="6">
              <h3 class="text-h6 font-weight-bold mb-3"
                >Confirmation Details</h3
              >
              <v-list density="compact">
                <v-list-item>
                  <v-list-item-title>Status</v-list-item-title>
                  <v-list-item-subtitle>
                    <v-chip
                      :color="getReconciliationStatusColor(selectedItem.status)"
                      size="small"
                    >
                      {{ formatReconciliationStatus(selectedItem.status) }}
                    </v-chip>
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Confirmed Amount</v-list-item-title>
                  <v-list-item-subtitle>
                    {{
                      selectedItem.confirmed_amount
                        ? formatCurrency(selectedItem.confirmed_amount)
                        : 'Not confirmed'
                    }}
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Variance</v-list-item-title>
                  <v-list-item-subtitle>
                    <span
                      :class="
                        selectedItem.variance === 0
                          ? 'text-success'
                          : 'text-error'
                      "
                    >
                      {{
                        selectedItem.variance !== null
                          ? formatCurrency(selectedItem.variance)
                          : 'N/A'
                      }}
                    </span>
                  </v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Match Score</v-list-item-title>
                  <v-list-item-subtitle
                    >{{ selectedItem.match_score }}%</v-list-item-subtitle
                  >
                </v-list-item>
                <v-list-item>
                  <v-list-item-title>Last Reconciled</v-list-item-title>
                  <v-list-item-subtitle>
                    {{
                      selectedItem.last_reconciled
                        ? formatDateTime(selectedItem.last_reconciled)
                        : 'Not reconciled'
                    }}
                  </v-list-item-subtitle>
                </v-list-item>
              </v-list>
            </v-col>
          </v-row>

          <v-divider class="my-4"></v-divider>

          <h3 class="text-h6 font-weight-bold mb-3">Reconciliation History</h3>
          <v-timeline
            v-if="selectedItem.reconciliation_history?.length"
            density="compact"
          >
            <v-timeline-item
              v-for="event in selectedItem.reconciliation_history"
              :key="event.date"
              size="small"
              :dot-color="getEventColor(event.type)"
            >
              <template #opposite>
                <span class="text-caption">{{
                  formatDateTime(event.date)
                }}</span>
              </template>
              <div>
                <p class="text-body-2 font-weight-bold">{{ event.title }}</p>
                <p class="text-caption text-medium-emphasis">{{
                  event.description
                }}</p>
              </div>
            </v-timeline-item>
          </v-timeline>
          <p v-else class="text-medium-emphasis"
            >No reconciliation history available</p
          >
        </v-card-text>
        <v-card-actions>
          <v-btn
            v-if="selectedItem.status === 'pending'"
            color="primary"
            variant="outlined"
            prepend-icon="mdi-link"
            @click="manualMatch(selectedItem)"
          >
            Manual Match
          </v-btn>
          <v-btn
            v-if="selectedItem.status === 'discrepancy'"
            color="warning"
            variant="outlined"
            prepend-icon="mdi-wrench"
            @click="resolveDiscrepancy(selectedItem)"
          >
            Resolve
          </v-btn>
          <v-btn
            color="success"
            variant="outlined"
            prepend-icon="mdi-download"
            @click="exportReconciliation(selectedItem)"
          >
            Export
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="showDetailDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Discrepancy Detail Dialog -->
    <v-dialog v-model="showDiscrepancyDialog" max-width="900">
      <v-card v-if="selectedDiscrepancyItem">
        <v-card-title class="text-h6 font-weight-bold bg-error text-white">
          <v-icon class="me-2">mdi-alert-circle</v-icon>
          Discrepancy Details -
          {{ selectedDiscrepancyItem.generated_bordereaux_id }}
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="error" variant="tonal" class="mb-4">
            <div class="d-flex align-center">
              <v-icon class="me-2">mdi-alert-triangle</v-icon>
              <div>
                <p class="font-weight-bold mb-1"
                  >Reconciliation Discrepancy Detected</p
                >
                <p class="mb-0"
                  >Review the differences below and take appropriate action</p
                >
              </div>
            </div>
          </v-alert>

          <v-row>
            <v-col cols="12" md="6">
              <v-card variant="outlined" color="info">
                <v-card-title
                  class="text-subtitle-1 font-weight-bold bg-info text-white"
                >
                  <v-icon class="me-2">mdi-upload</v-icon>
                  Submitted Values
                </v-card-title>
                <v-card-text>
                  <v-list density="compact">
                    <v-list-item>
                      <v-list-item-title>Amount</v-list-item-title>
                      <v-list-item-subtitle class="font-weight-bold text-info">
                        {{
                          formatCurrency(
                            selectedDiscrepancyItem.submitted_amount
                          )
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Record Count</v-list-item-title>
                      <v-list-item-subtitle>
                        {{
                          selectedDiscrepancyItem.record_count?.toLocaleString()
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Submission Date</v-list-item-title>
                      <v-list-item-subtitle>
                        {{
                          selectedDiscrepancyItem.submitted_date
                            ? formatDateTime(
                                selectedDiscrepancyItem.submitted_date
                              )
                            : 'N/A'
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Reference Number</v-list-item-title>
                      <v-list-item-subtitle>
                        {{
                          selectedDiscrepancyItem.submission_reference || 'N/A'
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>

            <v-col cols="12" md="6">
              <v-card variant="outlined" color="warning">
                <v-card-title
                  class="text-subtitle-1 font-weight-bold bg-warning text-white"
                >
                  <v-icon class="me-2">mdi-download</v-icon>
                  Confirmed Values
                </v-card-title>
                <v-card-text>
                  <v-list density="compact">
                    <v-list-item>
                      <v-list-item-title>Amount</v-list-item-title>
                      <v-list-item-subtitle
                        class="font-weight-bold text-warning"
                      >
                        {{
                          selectedDiscrepancyItem.confirmed_amount
                            ? formatCurrency(
                                selectedDiscrepancyItem.confirmed_amount
                              )
                            : 'Not confirmed'
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Confirmed Count</v-list-item-title>
                      <v-list-item-subtitle>
                        {{
                          selectedDiscrepancyItem.confirmed_count?.toLocaleString() ||
                          'N/A'
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Confirmation Date</v-list-item-title>
                      <v-list-item-subtitle>
                        {{
                          selectedDiscrepancyItem.confirmed_date
                            ? formatDateTime(
                                selectedDiscrepancyItem.confirmed_date
                              )
                            : 'N/A'
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title
                        >Confirmation Reference</v-list-item-title
                      >
                      <v-list-item-subtitle>
                        {{
                          selectedDiscrepancyItem.confirmation_reference ||
                          'N/A'
                        }}
                      </v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <v-divider class="my-4"></v-divider>

          <v-row>
            <v-col cols="12">
              <h3 class="text-h6 font-weight-bold mb-3 text-error">
                <v-icon class="me-2">mdi-compare</v-icon>
                Identified Discrepancies
              </h3>
              <v-card variant="outlined" color="error">
                <v-card-text>
                  <v-row>
                    <v-col cols="12" sm="4">
                      <div class="text-center">
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Amount Variance</p
                        >
                        <p
                          class="text-h5 font-weight-bold"
                          :class="
                            selectedDiscrepancyItem.variance === 0
                              ? 'text-success'
                              : 'text-error'
                          "
                        >
                          {{
                            selectedDiscrepancyItem.variance !== null
                              ? formatCurrency(selectedDiscrepancyItem.variance)
                              : 'N/A'
                          }}
                        </p>
                      </div>
                    </v-col>
                    <v-col cols="12" sm="4">
                      <div class="text-center">
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Record Variance</p
                        >
                        <p class="text-h5 font-weight-bold text-error">
                          {{ selectedDiscrepancyItem.record_variance || 0 }}
                        </p>
                      </div>
                    </v-col>
                    <v-col cols="12" sm="4">
                      <div class="text-center">
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Match Score</p
                        >
                        <p
                          class="text-h5 font-weight-bold"
                          :class="
                            getMatchScoreTextColor(
                              selectedDiscrepancyItem.match_score
                            )
                          "
                        >
                          {{ selectedDiscrepancyItem.match_score }}%
                        </p>
                      </div>
                    </v-col>
                  </v-row>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <v-divider class="my-4"></v-divider>

          <h3 class="text-h6 font-weight-bold mb-3"
            >Possible Causes & Actions</h3
          >
          <v-list>
            <v-list-item>
              <template #prepend>
                <v-icon color="info">mdi-information</v-icon>
              </template>
              <v-list-item-title>Data Entry Error</v-list-item-title>
              <v-list-item-subtitle
                >Check for typing mistakes in amounts or member
                details</v-list-item-subtitle
              >
            </v-list-item>
            <v-list-item>
              <template #prepend>
                <v-icon color="warning">mdi-clock</v-icon>
              </template>
              <v-list-item-title>Timing Mismatch</v-list-item-title>
              <v-list-item-subtitle
                >Confirm submission and confirmation are for the same
                period</v-list-item-subtitle
              >
            </v-list-item>
            <v-list-item>
              <template #prepend>
                <v-icon color="error">mdi-file-remove</v-icon>
              </template>
              <v-list-item-title>Missing Records</v-list-item-title>
              <v-list-item-subtitle
                >Some member records may not have been included in
                confirmation</v-list-item-subtitle
              >
            </v-list-item>
          </v-list>

          <v-divider class="my-4"></v-divider>

          <h3 class="text-h6 font-weight-bold mb-3">
            <v-icon class="me-2">mdi-table</v-icon>
            Detailed Discrepancy Records
          </h3>

          <v-card variant="outlined">
            <v-card-text>
              <v-data-table
                :headers="discrepancyHeaders"
                :items="discrepancyDetails"
                :loading="loadingDiscrepancyDetails"
                :items-per-page="10"
                item-key="id"
                density="compact"
              >
                <template #[`item.member_name`]="{ item }: { item: any }">
                  <span class="font-weight-medium">{{ item.member_name }}</span>
                </template>

                <template #[`item.record_id`]="{ item }: { item: any }">
                  <v-chip color="primary" size="x-small" variant="outlined">
                    {{ item.record_id }}
                  </v-chip>
                </template>

                <template #[`item.field`]="{ item }: { item: any }">
                  <v-chip color="info" size="small" variant="tonal">
                    {{ item.field }}
                  </v-chip>
                </template>

                <template #[`item.expected_value`]="{ item }: { item: any }">
                  <span class="font-mono text-success">
                    {{
                      item.field === 'Amount'
                        ? formatCurrency(parseFloat(item.expected_value))
                        : item.expected_value
                    }}
                  </span>
                </template>

                <template #[`item.actual_value`]="{ item }: { item: any }">
                  <span class="font-mono text-warning">
                    {{
                      item.field === 'Amount'
                        ? formatCurrency(parseFloat(item.actual_value))
                        : item.actual_value
                    }}
                  </span>
                </template>

                <template #[`item.variance`]="{ item }: { item: any }">
                  <v-chip
                    :color="
                      item.variance === 0
                        ? 'success'
                        : item.variance > 0
                          ? 'error'
                          : 'warning'
                    "
                    size="small"
                    variant="tonal"
                  >
                    {{
                      item.field === 'Amount'
                        ? formatCurrency(item.variance)
                        : item.variance
                    }}
                  </v-chip>
                </template>

                <template #[`item.status`]="{ item }: { item: any }">
                  <v-chip
                    :color="item.is_resolved ? 'success' : 'error'"
                    size="small"
                  >
                    {{ item.is_resolved ? 'Resolved' : 'Outstanding' }}
                  </v-chip>
                </template>

                <template #[`item.comments`]="{ item }: { item: any }">
                  <div
                    v-if="item.comments && item.comments.trim()"
                    class="d-flex align-center"
                  >
                    <span
                      class="text-caption text-truncate me-2"
                      style="max-width: 100px"
                    >
                      {{ item.comments.substring(0, 20)
                      }}{{ item.comments.length > 20 ? '...' : '' }}
                    </span>
                    <v-btn
                      size="x-small"
                      color="info"
                      variant="text"
                      icon="mdi-eye"
                      @click="viewComment(item)"
                    />
                  </div>
                  <div v-else class="d-flex align-center">
                    <span class="text-caption text-medium-emphasis me-2"
                      >No comments</span
                    >
                    <v-btn
                      size="x-small"
                      color="primary"
                      variant="text"
                      icon="mdi-comment-plus"
                      @click="addDiscrepancyComment(item)"
                    />
                  </div>
                </template>

                <template #[`item.actions`]="{ item }: { item: any }">
                  <div class="d-flex align-center gap-1">
                    <v-btn
                      size="x-small"
                      color="success"
                      variant="tonal"
                      icon="mdi-check"
                      :disabled="item.is_resolved"
                      @click="resolveDetailedDiscrepancy(item)"
                    />
                    <v-btn
                      size="x-small"
                      color="info"
                      variant="tonal"
                      icon="mdi-comment-plus"
                      @click="addDiscrepancyComment(item)"
                    />
                  </div>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="warning"
            variant="outlined"
            prepend-icon="mdi-wrench"
            @click="resolveDiscrepancy(selectedDiscrepancyItem)"
          >
            Resolve Discrepancy
          </v-btn>
          <v-btn
            color="info"
            variant="outlined"
            prepend-icon="mdi-account-supervisor"
            @click="escalateDiscrepancy(selectedDiscrepancyItem)"
          >
            Escalate
          </v-btn>
          <v-btn
            color="success"
            variant="outlined"
            prepend-icon="mdi-download"
            @click="exportDiscrepancyReport(selectedDiscrepancyItem)"
          >
            Export Report
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="closeDiscrepancyDialog"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete Confirmation Dialog -->
    <v-dialog v-model="showDeleteDialog" max-width="500">
      <v-card v-if="itemToDelete">
        <v-card-title class="text-h6 font-weight-bold bg-error text-white">
          <v-icon class="me-2">mdi-delete-alert</v-icon>
          Confirm Deletion
        </v-card-title>
        <v-card-text class="pt-6">
          <v-alert color="error" variant="tonal" class="mb-4">
            <div class="d-flex align-center">
              <v-icon class="me-2">mdi-alert-circle</v-icon>
              <div>
                <p class="font-weight-bold mb-1"
                  >Warning: This action cannot be undone</p
                >
                <p class="mb-0"
                  >You are about to permanently delete this reconciliation
                  item</p
                >
              </div>
            </div>
          </v-alert>

          <v-row>
            <v-col cols="12">
              <v-card variant="outlined" class="bg-grey-lighten-5">
                <v-card-text>
                  <v-list density="compact">
                    <v-list-item>
                      <v-list-item-title>Bordereaux ID</v-list-item-title>
                      <v-list-item-subtitle class="font-weight-bold">
                        {{ itemToDelete.generated_bordereaux_id }}
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Type</v-list-item-title>
                      <v-list-item-subtitle>{{
                        formatBordereauType(itemToDelete.type)
                      }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Scheme</v-list-item-title>
                      <v-list-item-subtitle>{{
                        itemToDelete.scheme_name
                      }}</v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Status</v-list-item-title>
                      <v-list-item-subtitle>
                        <v-chip
                          :color="
                            getReconciliationStatusColor(itemToDelete.status)
                          "
                          size="small"
                        >
                          {{ formatReconciliationStatus(itemToDelete.status) }}
                        </v-chip>
                      </v-list-item-subtitle>
                    </v-list-item>
                    <v-list-item>
                      <v-list-item-title>Submitted Amount</v-list-item-title>
                      <v-list-item-subtitle class="font-weight-bold">
                        {{ formatCurrency(itemToDelete.submitted_amount) }}
                      </v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <v-divider class="my-4"></v-divider>

          <p class="text-body-2 text-medium-emphasis">
            Are you sure you want to delete this reconciliation item? This will
            permanently remove all associated data including:
          </p>
          <v-list density="compact">
            <v-list-item>
              <template #prepend>
                <v-icon size="small" color="error">mdi-bullet</v-icon>
              </template>
              <v-list-item-title class="text-body-2"
                >Reconciliation history and notes</v-list-item-title
              >
            </v-list-item>
            <v-list-item>
              <template #prepend>
                <v-icon size="small" color="error">mdi-bullet</v-icon>
              </template>
              <v-list-item-title class="text-body-2"
                >Discrepancy details and comments</v-list-item-title
              >
            </v-list-item>
            <v-list-item>
              <template #prepend>
                <v-icon size="small" color="error">mdi-bullet</v-icon>
              </template>
              <v-list-item-title class="text-body-2"
                >Matching information and scores</v-list-item-title
              >
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="showDeleteDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            color="error"
            variant="flat"
            :loading="deleting"
            prepend-icon="mdi-delete"
            @click="deleteReconciliationItem"
          >
            Delete Item
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Comment Viewing Dialog -->
    <v-dialog v-model="showCommentDialog" max-width="600">
      <v-card v-if="selectedComment">
        <v-card-title class="text-h6 font-weight-bold bg-info text-white">
          <v-icon class="me-2">mdi-comment-text</v-icon>
          Discrepancy Comment
        </v-card-title>
        <v-card-subtitle class="pt-3">
          <div class="d-flex align-center text-medium-emphasis">
            <v-icon class="me-1" size="small">mdi-account</v-icon>
            <span class="me-3">{{ selectedComment.member_name }}</span>
            <v-icon class="me-1" size="small">mdi-identifier</v-icon>
            <span class="me-3">{{ selectedComment.record_id }}</span>
            <v-icon class="me-1" size="small">mdi-calendar</v-icon>
            <span>{{ formatDateTime(selectedComment.created_at) }}</span>
          </div>
        </v-card-subtitle>
        <v-card-text class="pt-4">
          <v-card variant="outlined" class="bg-grey-lighten-5">
            <v-card-text>
              <p class="text-body-2 mb-0" style="white-space: pre-wrap">{{
                selectedComment.comments
              }}</p>
            </v-card-text>
          </v-card>
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            variant="outlined"
            prepend-icon="mdi-pencil"
            @click="editComment(selectedComment)"
          >
            Edit Comment
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn
            color="grey"
            variant="outlined"
            @click="showCommentDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as ExcelJS from 'exceljs'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import * as pako from 'pako'
import { useFlashStore } from '@/renderer/store/flash'

const flash = useFlashStore()

// Reactive data
const loading = ref(false)
const autoReconciling = ref(false)
const importing = ref(false)
const compressing = ref(false)
const compressionProgress = ref(0)
const showImportDialog = ref(false)
const showDetailDialog = ref(false)
const showDiscrepancyDialog = ref(false)
const showCommentDialog = ref(false)
const showDeleteDialog = ref(false)
const selectedItem: any = ref(null)
const itemToDelete: any = ref(null)
const deleting = ref(false)
const selectedDiscrepancyItem: any = ref(null)
const selectedComment: any = ref(null)
const discrepancyDetails: any = ref([])
const loadingDiscrepancyDetails = ref(false)
const viewMode = ref('all')

const reconciliationStats = ref({
  pending: 7,
  matched: 18,
  discrepancies: 3,
  autoMatch: 85
})

const importData: any = ref({
  scheme_id: null,
  file_type: '',
  files: [],
  auto_process: true
})

// Static data
const headers = [
  { title: 'Bordereaux ID', key: 'generated_bordereaux_id', sortable: true },
  { title: 'Type', key: 'type', sortable: true },
  { title: 'Scheme', key: 'scheme_name', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Match Score', key: 'match_score', sortable: true },
  { title: 'Submitted Amount', key: 'submitted_amount', sortable: true },
  { title: 'Confirmed Amount', key: 'confirmed_amount', sortable: true },
  { title: 'Variance', key: 'variance', sortable: true },
  { title: 'Last Reconciled', key: 'last_reconciled', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false, width: '200px' }
]

const discrepancyHeaders = [
  { title: 'Member Name', key: 'member_name', sortable: true },
  { title: 'Record ID', key: 'record_id', sortable: true },
  { title: 'Field', key: 'field', sortable: true },
  { title: 'Expected Value', key: 'expected_value', sortable: true },
  { title: 'Actual Value', key: 'actual_value', sortable: true },
  { title: 'Variance', key: 'variance', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Comments', key: 'comments', sortable: false },
  { title: 'Actions', key: 'actions', sortable: false, width: '120px' }
]

const schemes = ref([])
const schemesLoading = ref(false)

const fileTypes = [
  { title: 'Premium Confirmation', value: 'premium_confirmation' },
  { title: 'Claims Settlement', value: 'claims_settlement' },
  { title: 'Member Acknowledgment', value: 'member_acknowledgment' },
  { title: 'Payment Advice', value: 'payment_advice' }
]

const reconciliationItems: any = ref([])

// Computed properties
const filteredItems = computed(() => {
  console.log('Filtering items for view mode:', viewMode.value)
  console.log('Total reconciliation items:', reconciliationItems.value.length)
  console.log('Reconciliation itemsv2:', reconciliationItems.value)
  switch (viewMode.value) {
    case 'pending':
      return reconciliationItems.value.filter(
        (item) => item.status === 'pending'
      )
    case 'discrepancies':
      return reconciliationItems.value.filter(
        (item) => item.status === 'discrepancy'
      )
    case 'matched':
      return reconciliationItems.value.filter(
        (item) => item.status === 'matched' || item.status === 'Matched'
      )
    default:
      return reconciliationItems.value
  }
})

// Methods
const getBordereauTypeColor = (type: string): string => {
  const colors: Record<string, string> = {
    member: 'blue',
    premium: 'green',
    claims: 'orange',
    benefits: 'purple'
  }
  return colors[type] || 'grey'
}

const formatBordereauType = (type: string): string => {
  const types: Record<string, string> = {
    member: 'Member',
    premium: 'Premium',
    claims: 'Claims',
    benefits: 'Benefits'
  }
  return types[type] || type
}

const getReconciliationStatusColor = (status: string): string => {
  const colors: Record<string, string> = {
    pending: 'warning',
    matched: 'success',
    Matched: 'success',
    discrepancy: 'error',
    investigating: 'info',
    resolved: 'success'
  }
  return colors[status] || 'grey'
}

const formatReconciliationStatus = (status: string): string => {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const getMatchScoreColor = (score: number): string => {
  if (score === 100) return 'success'
  if (score >= 80) return 'info'
  if (score >= 60) return 'warning'
  return 'error'
}

const getMatchScoreTextColor = (score: number): string => {
  if (score === 100) return 'text-success'
  if (score >= 80) return 'text-info'
  if (score >= 60) return 'text-warning'
  return 'text-error'
}

const getAmountClass = (item: any): string => {
  if (!item.confirmed_amount || !item.submitted_amount) return ''
  return item.variance === 0 ? 'text-success' : 'text-error'
}

const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    minimumFractionDigits: 2
  }).format(amount)
}

const formatDateTime = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getEventColor = (type: string): string => {
  const colors: Record<string, string> = {
    submitted: 'blue',
    confirmed: 'teal',
    matched: 'success',
    discrepancy: 'error',
    resolved: 'success'
  }
  return colors[type] || 'grey'
}

const runAutoReconciliation = async () => {
  autoReconciling.value = true
  try {
    const response = await GroupPricingService.runAutoReconciliation()

    if (response.data.success) {
      await fetchReconciliationItems()
      flash.show('Auto-reconciliation complete', 'success')
    } else {
      flash.show(response.data.error || 'Auto-reconciliation failed', 'error')
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to run auto-reconciliation',
      'error'
    )
  } finally {
    autoReconciling.value = false
  }
}

const fetchReconciliationItems = async () => {
  loading.value = true
  try {
    const response = await GroupPricingService.getReconciliationItems()
    // console.log('Full API response:', response)
    // console.log('Response data:', response.data)
    // console.log('Response data type:', typeof response.data)
    // console.log('Is response.data an array?', Array.isArray(response.data))

    // Handle different possible response structures
    let items = []
    if (Array.isArray(response.data)) {
      // Direct array response
      items = response.data
    } else if (response.data && Array.isArray(response.data.items)) {
      // Nested items array
      items = response.data.items
    } else if (response.data && Array.isArray(response.data.results)) {
      // Nested results array
      items = response.data.results
    } else if (response.data && Array.isArray(response.data.data)) {
      // Double nested data array
      items = response.data.data
    } else {
      console.warn('Unexpected API response structure:', response.data)
      items = []
    }

    reconciliationItems.value = items

    // Update stats based on fetched data
    updateReconciliationStats()
  } catch (error) {
    console.error('Error fetching reconciliation items:', error)
    reconciliationItems.value = []
  } finally {
    loading.value = false
  }
}

const refreshReconciliation = async () => {
  await fetchReconciliationItems()
}

const updateReconciliationStats = () => {
  const stats = {
    pending: 0,
    matched: 0,
    discrepancies: 0,
    total: reconciliationItems.value.length
  }

  reconciliationItems.value.forEach((item) => {
    switch (item.status) {
      case 'pending':
        stats.pending++
        break
      case 'matched':
        stats.matched++
        break
      case 'discrepancy':
        stats.discrepancies++
        break
    }
  })

  // Calculate auto-match percentage
  const autoMatch =
    stats.total > 0 ? Math.round((stats.matched / stats.total) * 100) : 0

  reconciliationStats.value = {
    pending: stats.pending,
    matched: stats.matched,
    discrepancies: stats.discrepancies,
    autoMatch
  }
}

// Helper function to compress a file using gzip
const compressFile = async (file: File): Promise<File> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = () => {
      try {
        const arrayBuffer = reader.result as ArrayBuffer
        const uint8Array = new Uint8Array(arrayBuffer)

        // Compress using gzip
        const compressed = pako.gzip(uint8Array)

        // Create a new compressed file with .gz extension
        const compressedBlob = new Blob([compressed], {
          type: 'application/gzip'
        })
        const compressedFile = new File([compressedBlob], `${file.name}.gz`, {
          type: 'application/gzip',
          lastModified: file.lastModified
        })

        resolve(compressedFile)
      } catch (error) {
        reject(new Error(`Failed to compress ${file.name}: ${error}`))
      }
    }

    reader.onerror = () => {
      reject(new Error(`Failed to read ${file.name}`))
    }

    reader.readAsArrayBuffer(file)
  })
}

const performImport = async () => {
  if (!importData.value.file_type || !importData.value.files?.length) {
    return
  }

  importing.value = true
  compressing.value = true
  compressionProgress.value = 0

  try {
    // Compress files with progress tracking
    const compressedFiles: any[] = []
    const totalFiles = importData.value.files.length

    for (let i = 0; i < totalFiles; i++) {
      const file = importData.value.files[i]

      // Only compress files larger than 1MB to improve efficiency
      const shouldCompress = file.size > 1024 * 1024 // 1MB threshold

      if (shouldCompress) {
        console.log(
          `Compressing ${file.name} (${(file.size / (1024 * 1024)).toFixed(2)} MB)...`
        )
        const compressedFile = await compressFile(file)
        const compressionRatio = (
          (1 - compressedFile.size / file.size) *
          100
        ).toFixed(1)
        console.log(`Compressed ${file.name}: ${compressionRatio}% reduction`)
        compressedFiles.push(compressedFile)
      } else {
        // Use original file for small files
        compressedFiles.push(file)
      }

      // Update progress
      compressionProgress.value = Math.round(((i + 1) / totalFiles) * 100)
    }

    compressing.value = false

    // Create FormData to handle file uploads
    const formData = new FormData()
    formData.append('file_type', importData.value.file_type)
    formData.append('auto_process', importData.value.auto_process.toString())
    formData.append('files_compressed', 'true') // Flag to indicate compression was used

    // Add all compressed files
    for (const file of compressedFiles) {
      formData.append('confirmation_files', file)
    }

    // Call the import API
    const response =
      await GroupPricingService.importSchemeConfirmations(formData)

    if (response.data.success) {
      if (response.data.imported_count) {
        reconciliationStats.value.pending += response.data.imported_count
      }

      if (importData.value.auto_process && response.data.auto_process_results) {
        reconciliationStats.value.matched +=
          response.data.auto_process_results.matched || 0
        reconciliationStats.value.discrepancies +=
          response.data.auto_process_results.discrepancies || 0
        reconciliationStats.value.pending -=
          (response.data.auto_process_results.matched || 0) +
          (response.data.auto_process_results.discrepancies || 0)
      }

      importData.value = {
        scheme_id: null,
        file_type: '',
        files: [],
        auto_process: true
      }

      showImportDialog.value = false
      await refreshReconciliation()

      const importedCount = response.data.imported_count ?? 0
      flash.show(
        `Imported ${importedCount} confirmation file${importedCount === 1 ? '' : 's'}`,
        'success'
      )
    } else {
      flash.show(response.data.error || 'Confirmation import failed', 'error')
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to import confirmation files',
      'error'
    )
  } finally {
    importing.value = false
    compressing.value = false
    compressionProgress.value = 0
  }
}

const viewReconciliationDetail = (item: any) => {
  selectedItem.value = item
  showDetailDialog.value = true
}

const viewDiscrepancies = async (item: any) => {
  selectedDiscrepancyItem.value = item
  showDiscrepancyDialog.value = true

  // Fetch detailed discrepancy data
  await fetchDiscrepancyDetails(item.id)
}

const fetchDiscrepancyDetails = async (reconciliationItemId: number) => {
  loadingDiscrepancyDetails.value = true
  try {
    const response =
      await GroupPricingService.getDiscrepancyDetails(reconciliationItemId)
    console.log('Fetched discrepancy details:', response.data)

    // Handle different possible response structures
    if (Array.isArray(response.data)) {
      discrepancyDetails.value = response.data
    } else if (response.data && Array.isArray(response.data.details)) {
      discrepancyDetails.value = response.data.details
    } else if (response.data && Array.isArray(response.data.items)) {
      discrepancyDetails.value = response.data.items
    } else {
      console.warn(
        'Unexpected discrepancy details response structure:',
        response.data
      )
      discrepancyDetails.value = []
    }
  } catch (error) {
    console.error('Error fetching discrepancy details:', error)
    discrepancyDetails.value = []
  } finally {
    loadingDiscrepancyDetails.value = false
  }
}

const resolveDetailedDiscrepancy = async (detailItem: any) => {
  try {
    detailItem.is_resolved = true
    detailItem.status = 'resolved'

    await GroupPricingService.resolveDiscrepancy(detailItem.id, {
      resolution: 'accept_expected',
      notes: 'Resolved from discrepancy detail view'
    })
    flash.show('Discrepancy resolved', 'success')
  } catch (error: any) {
    detailItem.is_resolved = false
    detailItem.status = 'discrepancy'
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to resolve discrepancy',
      'error'
    )
  }
}

const addDiscrepancyComment = async (detailItem: any) => {
  const comment = prompt('Enter comment:')
  if (!comment) return
  try {
    const response = await GroupPricingService.addDiscrepancyComment(
      detailItem.id,
      { comment }
    )
    if (response.data?.data) {
      detailItem.comments = response.data.data.comments
    }
    flash.show('Comment added', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error || error.message || 'Failed to add comment',
      'error'
    )
  }
}

const viewComment = (detailItem: any) => {
  selectedComment.value = detailItem
  showCommentDialog.value = true
}

const editComment = async (detailItem: any) => {
  const comment = prompt('Edit comment:', detailItem.comments || '')
  if (comment === null) return
  try {
    await GroupPricingService.addDiscrepancyComment(detailItem.id, { comment })
    detailItem.comments = comment
    flash.show('Comment updated', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to update comment',
      'error'
    )
  }
  showCommentDialog.value = false
}

const closeDiscrepancyDialog = () => {
  showDiscrepancyDialog.value = false
  // Clear detailed data when closing dialog
  discrepancyDetails.value = []
  selectedDiscrepancyItem.value = null
}

const escalateDiscrepancy = async (item: any) => {
  const reason = prompt('Reason for escalation:')
  if (!reason) return
  try {
    await GroupPricingService.escalateDiscrepancy(
      item.id || selectedDiscrepancyItem.value?.id,
      {
        escalate_to: 'manager',
        reason: reason,
        priority: 'high'
      }
    )
    showDiscrepancyDialog.value = false
    await refreshReconciliation()
    flash.show('Discrepancy escalated', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to escalate discrepancy',
      'error'
    )
  }
}

const exportDiscrepancyReport = async (item: any) => {
  if (!item || !discrepancyDetails.value.length) {
    flash.show('No discrepancy data to export', 'warning')
    return
  }

  const workbook = new ExcelJS.Workbook()
  const worksheet = workbook.addWorksheet('Discrepancy Report')

  // --- Report Header ---
  worksheet.mergeCells('A1:F1')
  const titleCell = worksheet.getCell('A1')
  titleCell.value = `Discrepancy Report - Bordereaux ID: ${item.generated_bordereaux_id}`
  titleCell.font = { name: 'Arial', size: 16, bold: true }
  titleCell.alignment = { vertical: 'middle', horizontal: 'center' }
  worksheet.getRow(1).height = 40

  // --- Summary Information ---
  const summaryData = [
    { label: 'Scheme Name', value: item.scheme_name },
    { label: 'Bordereaux Type', value: formatBordereauType(item.type) },
    { label: 'Submitted Amount', value: formatCurrency(item.submitted_amount) },
    { label: 'Confirmed Amount', value: formatCurrency(item.confirmed_amount) },
    { label: 'Amount Variance', value: formatCurrency(item.variance) },
    { label: 'Record Variance', value: item.record_variance || 0 },
    { label: 'Export Date', value: new Date().toLocaleString('en-ZA') }
  ]

  let currentRow = 3
  summaryData.forEach((data) => {
    worksheet.mergeCells(`A${currentRow}:B${currentRow}`)
    worksheet.mergeCells(`C${currentRow}:D${currentRow}`)
    const labelCell = worksheet.getCell(`A${currentRow}`)
    const valueCell = worksheet.getCell(`C${currentRow}`)
    labelCell.value = data.label
    labelCell.font = { bold: true }
    valueCell.value = data.value
    currentRow++
  })

  // --- Discrepancy Details Table ---
  currentRow += 2 // Add a gap before the table
  const headerRow = worksheet.getRow(currentRow)
  headerRow.values = discrepancyHeaders.map((h) => h.title)
  headerRow.font = { bold: true, color: { argb: 'FFFFFFFF' } }
  headerRow.eachCell((cell) => {
    cell.fill = {
      type: 'pattern',
      pattern: 'solid',
      fgColor: { argb: 'FF4472C4' } // Blue background
    }
    cell.border = {
      top: { style: 'thin' },
      left: { style: 'thin' },
      bottom: { style: 'thin' },
      right: { style: 'thin' }
    }
  })

  // Add data rows
  discrepancyDetails.value.forEach((detail: any) => {
    currentRow++
    const row = worksheet.getRow(currentRow)
    row.values = [
      detail.member_name,
      detail.record_id,
      detail.field,
      detail.field === 'Amount'
        ? parseFloat(detail.expected_value)
        : detail.expected_value,
      detail.field === 'Amount'
        ? parseFloat(detail.actual_value)
        : detail.actual_value,
      detail.variance,
      detail.is_resolved ? 'Resolved' : 'Outstanding',
      detail.comments || '',
      '' // Actions column left blank
    ]

    // Apply formatting
    const expectedCell = row.getCell(4)
    const actualCell = row.getCell(5)
    const varianceCell = row.getCell(6)

    if (detail.field === 'Amount') {
      expectedCell.numFmt = '"R"#,##0.00'
      actualCell.numFmt = '"R"#,##0.00'
      varianceCell.numFmt = '"R"#,##0.00'
    }

    if (detail.variance !== 0) {
      varianceCell.font = { color: { argb: 'FFFF0000' } } // Red for variance
    }
  })

  // --- Column Widths ---
  worksheet.columns.forEach((column) => {
    if (!column) return
    let maxLength = 0
    if (column.eachCell) {
      column.eachCell({ includeEmpty: true }, (cell) => {
        const columnLength = cell.value ? cell.value.toString().length : 10
        if (columnLength > maxLength) {
          maxLength = columnLength
        }
      })
    }
    column.width = maxLength < 15 ? 15 : maxLength + 2
  })

  // --- Generate and Download File ---
  try {
    const buffer = await workbook.xlsx.writeBuffer()
    const blob = new Blob([buffer], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `Discrepancy_Report_${item.generated_bordereaux_id}_${new Date().toISOString().slice(0, 10)}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    flash.show('Discrepancy report generated', 'success')
  } catch (error: any) {
    flash.show(
      error.message || 'Failed to generate discrepancy report',
      'error'
    )
  }
}

const manualMatch = async (item: any) => {
  try {
    await GroupPricingService.resolveDiscrepancy(item.id, {
      resolution: 'manual_override',
      notes: 'Manually matched'
    })
    await refreshReconciliation()
    flash.show('Marked as manually matched', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to mark as manually matched',
      'error'
    )
  }
}

const resolveDiscrepancy = async (item: any) => {
  try {
    await GroupPricingService.resolveDiscrepancy(item.id, {
      resolution: 'accept_expected',
      notes: 'Discrepancy resolved'
    })
    await refreshReconciliation()
    flash.show('Discrepancy resolved', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to resolve discrepancy',
      'error'
    )
  }
}

const confirmReconciliation = async (item: any) => {
  try {
    await GroupPricingService.confirmReconciliation(item.id)
    await refreshReconciliation()
    flash.show('Reconciliation confirmed', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to confirm reconciliation',
      'error'
    )
  }
}

const exportReconciliation = async (item: any) => {
  try {
    const response = await GroupPricingService.getDiscrepancyDetails(item.id)
    const data = response.data?.data || response.data || []
    if (!data.length) {
      flash.show('No reconciliation data to export', 'warning')
      return
    }
    // Re-use the existing Excel export from exportDiscrepancyReport
    await exportDiscrepancyReport(item)
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to export reconciliation',
      'error'
    )
  }
}

const reprocessItem = async (item: any) => {
  try {
    await GroupPricingService.reprocessReconciliation(item.id)
    await refreshReconciliation()
    flash.show('Reconciliation reprocessed', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to reprocess reconciliation',
      'error'
    )
  }
}

const addNote = async (item: any) => {
  const note = prompt('Enter note:')
  if (!note) return
  try {
    await GroupPricingService.addReconciliationNote(item.id, { note })
    flash.show('Note added', 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error || error.message || 'Failed to add note',
      'error'
    )
  }
}

const confirmDelete = (item: any) => {
  itemToDelete.value = item
  showDeleteDialog.value = true
}

const deleteReconciliationItem = async () => {
  if (!itemToDelete.value) return

  deleting.value = true
  try {
    // Call the API to delete the reconciliation item
    const response = await GroupPricingService.deleteReconciliationItem(
      itemToDelete.value.id
    )

    if (response.status === 200) {
      const index = reconciliationItems.value.findIndex(
        (item) => item.id === itemToDelete.value.id
      )
      if (index !== -1) {
        reconciliationItems.value.splice(index, 1)
      }

      updateReconciliationStats()

      showDeleteDialog.value = false
      itemToDelete.value = null

      flash.show('Reconciliation item deleted', 'success')
    } else {
      flash.show(
        response.data.error || 'Failed to delete reconciliation item',
        'error'
      )
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to delete reconciliation item',
      'error'
    )
  } finally {
    deleting.value = false
  }
}

const generateReconciliationReport = async () => {
  try {
    // Export all reconciliation items as an Excel report using the existing pattern
    if (!reconciliationItems.value.length) return
    const workbook = new ExcelJS.Workbook()
    const worksheet = workbook.addWorksheet('Reconciliation Report')
    worksheet.columns = [
      { header: 'Bordereaux ID', key: 'generated_bordereaux_id', width: 20 },
      { header: 'Scheme', key: 'scheme_name', width: 25 },
      { header: 'Type', key: 'type', width: 15 },
      { header: 'Status', key: 'status', width: 15 },
      { header: 'Submitted', key: 'submitted_amount', width: 18 },
      { header: 'Confirmed', key: 'confirmed_amount', width: 18 },
      { header: 'Variance', key: 'variance', width: 18 },
      { header: 'Match Score', key: 'match_score', width: 12 },
      { header: 'Matched', key: 'matched_count', width: 10 },
      { header: 'Discrepancies', key: 'discrepancy_count', width: 14 }
    ]
    reconciliationItems.value.forEach((item: any) => worksheet.addRow(item))
    const buffer = await workbook.xlsx.writeBuffer()
    const blob = new Blob([buffer], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `Reconciliation_Report_${new Date().toISOString().slice(0, 10)}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (error) {
    console.error('Error generating report:', error)
  }
}

const fetchSchemes = async () => {
  schemesLoading.value = true
  try {
    const response = await GroupPricingService.getSchemesInforce()
    schemes.value = response.data
  } catch (error) {
    console.error('Error fetching schemes:', error)
    schemes.value = []
  } finally {
    schemesLoading.value = false
  }
}

onMounted(async () => {
  await fetchSchemes()
  await fetchReconciliationItems()
})
</script>

<style scoped>
.match-score-container {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 100px;
}

.font-mono {
  font-family: 'Courier New', monospace;
}
</style>
