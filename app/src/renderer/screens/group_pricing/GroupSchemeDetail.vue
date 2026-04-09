<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card v-if="scheme" :show-actions="false">
          <template #header>
            <span class="headline">{{ scheme.name }}</span>
          </template>
          <template #default>
            <v-row class="mb-n5">
              <v-col>
                <v-btn variant="plain" @click="goBack">{{ backText }}</v-btn>
              </v-col>
            </v-row>
            <v-row align="center" class="mt-2 mb-4">
              <v-col cols="12" md="8">
                <h1 class="text-h4 font-weight-bold text-grey-darken-4">{{
                  scheme.name
                }}</h1>
              </v-col>
              <v-col cols="12" md="4" class="text-md-end">
                <v-chip
                  :color="statusColor"
                  label
                  size="large"
                  class="font-weight-bold"
                >
                  {{ formatStatus(scheme.status) }}
                </v-chip>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="3">
                <v-card variant="tonal" class="pa-2">
                  <v-card-title
                    class="text-subtitle-1 text-grey-darken-1 d-flex justify-space-between"
                  >
                    <span>Member Count</span>
                    <v-icon>mdi-account-group-outline</v-icon>
                  </v-card-title>
                  <v-card-text
                    class="text-h6 font-weight-bold text-grey-darken-4 pt-1"
                  >
                    {{
                      new Intl.NumberFormat('fr-FR').format(scheme.member_count)
                    }}
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="3">
                <v-card variant="tonal" class="pa-2">
                  <v-card-title
                    class="text-subtitle-1 text-grey-darken-1 d-flex justify-space-between"
                  >
                    <span>Annual Premium</span>
                    <v-icon>mdi-cash-multiple</v-icon>
                  </v-card-title>
                  <v-card-text
                    class="text-h6 font-weight-bold text-grey-darken-4 pt-1"
                  >
                    R
                    {{ roundUpToTwoDecimalsAccounting(scheme.annual_premium) }}
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="3">
                <v-card variant="tonal" class="pa-2">
                  <v-card-title
                    class="text-subtitle-1 text-grey-darken-1 d-flex justify-space-between"
                  >
                    <span>Actual Claims</span>
                    <v-icon>mdi-file-chart-outline</v-icon>
                  </v-card-title>
                  <v-card-text
                    class="text-h6 font-weight-bold text-grey-darken-4 pt-1"
                  >
                    R{{ roundUpToTwoDecimalsAccounting(scheme.actual_claims) }}
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="3">
                <v-card variant="tonal" class="pa-2">
                  <v-card-title
                    class="text-subtitle-1 text-grey-darken-1 d-flex justify-space-between"
                  >
                    <span>Expected Claims</span>
                    <v-icon>mdi-file-chart-outline</v-icon>
                  </v-card-title>
                  <v-card-text
                    class="text-h6 font-weight-bold text-grey-darken-4 pt-1"
                  >
                    R{{
                      roundUpToTwoDecimalsAccounting(scheme.expected_claims)
                    }}
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col>
                <v-card class="mt-6" variant="outlined">
                  <v-card-title
                    class="d-flex align-center justify-space-between"
                  >
                    <span class="text-h6">Scheme Details</span>
                    <v-spacer></v-spacer>
                    <v-btn
                      size="small"
                      variant="text"
                      color="primary"
                      @click="openManageStatusDialog"
                      >Manage Status</v-btn
                    >
                    <v-btn
                      size="small"
                      variant="text"
                      color="info"
                      @click="openStatusHistoryDialog"
                      >View Status History</v-btn
                    >
                    <v-btn
                      size="small"
                      variant="text"
                      color="secondary"
                      @click="openQuotesDialog"
                      >View Associated Quotes</v-btn
                    >

                    <v-btn
                      size="small"
                      variant="text"
                      color="primary"
                      @click="openEditEndDateDialog"
                      >Update Cover Date</v-btn
                    >
                    <v-btn
                      size="small"
                      variant="text"
                      color="red"
                      @click="deleteScheme"
                      >Delete Scheme</v-btn
                    >
                  </v-card-title>
                  <v-divider></v-divider>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Cover Period</div
                        >
                        <div class="text-body-1 font-weight-medium"
                          >{{
                            formatDateString(
                              scheme.cover_start_date,
                              true,
                              true,
                              true
                            )
                          }}
                          to
                          {{
                            formatDateString(
                              scheme.cover_end_date,
                              true,
                              true,
                              true
                            )
                          }}</div
                        >
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Broker</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          scheme.broker.name
                        }}</div>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1">Type</div>
                        <div class="text-body-1 font-weight-medium">{{
                          scheme.quote.obligation_type
                        }}</div>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Expected Claims</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          roundUpToTwoDecimalsAccounting(scheme.expected_claims)
                        }}</div>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Expected Claims Ratio</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          roundUpToTwoDecimalsAccounting(
                            scheme.expected_claims_ratio
                          )
                        }}</div>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Commission</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          roundUpToTwoDecimalsAccounting(scheme.commission)
                        }}</div>
                      </v-col>

                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Contact Person</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          scheme.contact_person
                        }}</div>
                        <a
                          :href="`mailto:${scheme.contact_email}`"
                          class="text-body-2 text-primary"
                          >{{ scheme.contact_email }}</a
                        >
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Created By</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          scheme.created_by
                        }}</div>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <div class="text-caption text-grey-darken-1"
                          >Creation Date</div
                        >
                        <div class="text-body-1 font-weight-medium">{{
                          formatDateString(
                            scheme.creation_date,
                            true,
                            true,
                            true
                          )
                        }}</div>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Contribution Config Tab Section -->
            <v-row class="mt-6">
              <v-col>
                <v-tabs v-model="activeSchemeTab" color="primary" class="mb-5">
                  <v-tab value="contribution-config">Contribution Config</v-tab>
                </v-tabs>
                <v-window v-model="activeSchemeTab">
                  <v-window-item value="contribution-config">
                    <contribution-config v-if="scheme" :scheme-id="scheme.id" />
                  </v-window-item>
                </v-window>
              </v-col>
            </v-row>

            <v-row
              v-for="item in relatedTables"
              :key="item.table_type"
              align="center"
              class="mt-6"
            >
              <v-col cols="12" md="5">
                <h2 class="text-h6 font-weight-bold">Member Data</h2>
                <div class="text-body-3 text-grey-darken-1"
                  >Manage members associated with this scheme.</div
                >
              </v-col>
              <v-col
                cols="12"
                md="7"
                class="d-flex justify-start justify-md-end flex-wrap ga-2"
              >
                <v-btn
                  size="small"
                  variant="outlined"
                  @click.stop="loadMembersForScheme"
                  >View Members</v-btn
                >
                <v-btn
                  color="primary"
                  flat
                  size="small"
                  prepend-icon="mdi-plus"
                  @click.stop="openMemberDialog(item)"
                  >Add Member</v-btn
                >
              </v-col>
            </v-row>

            <!-- Member Management Section -->
            <v-row v-if="showMemberManagement" class="mt-6">
              <v-col>
                <v-card variant="outlined">
                  <v-card-title
                    class="d-flex justify-space-between align-center"
                  >
                    <span class="text-h6">Scheme Members</span>
                    <v-btn
                      size="small"
                      variant="outlined"
                      @click="hideMemberManagement"
                    >
                      Hide Members
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <!-- Empty state when no members loaded -->
                    <v-row
                      v-if="
                        !loadingMembers &&
                        members.length === 0 &&
                        !hasMemberSearched
                      "
                      class="mb-4"
                    >
                      <v-col cols="12" class="text-center py-8">
                        <v-icon size="64" color="grey-lighten-1" class="mb-4"
                          >mdi-account-search</v-icon
                        >
                        <h3 class="text-h6 text-grey-darken-1 mb-2"
                          >No Members Loaded</h3
                        >
                        <p class="text-body-1 text-grey mb-4">
                          Use the search and filters above, then click "Search
                          Members" to load member data.
                          <br />
                          <small
                            >This prevents loading potentially large datasets
                            automatically.</small
                          >
                        </p>
                      </v-col>
                    </v-row>

                    <!-- Loading Progress Bar -->
                    <v-row v-if="loadingMembers" class="mb-4">
                      <v-col cols="12">
                        <v-card variant="outlined" class="pa-4">
                          <div
                            class="d-flex justify-space-between align-center mb-2"
                          >
                            <span>{{ memberLoadingMessage }}</span>
                            <v-btn
                              size="small"
                              color="error"
                              variant="outlined"
                              @click="cancelMemberLoading"
                            >
                              Cancel
                            </v-btn>
                          </div>
                          <v-progress-linear
                            :model-value="memberLoadingProgress"
                            color="primary"
                            height="8"
                            rounded
                          />
                        </v-card>
                      </v-col>
                    </v-row>

                    <!-- Search and Filter Bar -->
                    <v-row class="mb-4">
                      <v-col cols="12" md="4">
                        <v-text-field
                          v-model="memberSearchQuery"
                          label="Search Members"
                          prepend-inner-icon="mdi-magnify"
                          variant="outlined"
                          density="compact"
                          clearable
                          @input="searchSchemeMembers"
                        />
                      </v-col>
                      <v-col cols="12" md="3">
                        <v-select
                          v-model="selectedMemberStatus"
                          :items="memberStatuses"
                          label="Filter by Status"
                          variant="outlined"
                          density="compact"
                          clearable
                          @update:model-value="filterMembersByStatus"
                        />
                      </v-col>
                      <v-col cols="12" md="5" class="d-flex gap-2">
                        <v-btn
                          class="mr-2 mt-1"
                          rounded
                          size="small"
                          color="primary"
                          variant="outlined"
                          :loading="loadingMembers"
                          @click="reloadSchemeMembers"
                        >
                          Search Members
                        </v-btn>
                        <v-btn
                          class="mt-1"
                          rounded
                          size="small"
                          color="info"
                          variant="outlined"
                          @click="exportSchemeMembers"
                        >
                          Export
                        </v-btn>
                      </v-col>
                    </v-row>

                    <!-- Members Data Grid -->
                    <v-row v-if="members.length > 0 || loadingMembers">
                      <v-col>
                        <data-grid
                          :columnDefs="memberColumnDefs"
                          :rowData="members"
                          :pagination="false"
                          :loading="loadingMembers"
                          @row-double-clicked="handleMemberRowClick"
                        />
                      </v-col>
                    </v-row>

                    <!-- Pagination Controls -->
                    <v-row v-if="members.length > 0" class="mt-4">
                      <v-col cols="12" md="6">
                        <v-card variant="outlined" class="pa-3">
                          <div class="text-body-2 text-medium-emphasis">
                            Showing
                            {{ memberPaginationInfo.displayedMembers }} of
                            {{ memberPaginationInfo.totalMembers }} members
                          </div>
                        </v-card>
                      </v-col>
                      <v-col cols="12" md="6" class="d-flex justify-end">
                        <v-btn
                          v-if="memberPaginationInfo.hasMore"
                          rounded
                          size="small"
                          :loading="loadingMembers"
                          :disabled="loadingMembers"
                          color="primary"
                          variant="outlined"
                          @click="loadMoreSchemeMembers"
                        >
                          Load More Members
                        </v-btn>
                        <v-btn
                          v-if="members.length > 0"
                          rounded
                          size="small"
                          :loading="loadingMembers"
                          :disabled="loadingMembers"
                          color="info"
                          variant="outlined"
                          class="ml-2"
                          @click="reloadSchemeMembers"
                        >
                          Refresh
                        </v-btn>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-snackbar
      v-model="snackbar"
      centered
      :timeout="timeout"
      :multi-line="true"
    >
      {{ snackbarText }}
      <v-btn rounded color="red" variant="text" @click="snackbar = false"
        >Close</v-btn
      >
    </v-snackbar>
    <v-row>
      <v-col>
        <file-upload-dialog
          :yearLabel="yearLabel"
          :isDialogOpen="isDialogOpen"
          :showModelPoint="showModelPoint"
          :mpLabel="mpLabel"
          :table="'undefined'"
          :uploadTitle="uploadTitle"
          :years="years"
          @upload="handleUpload"
          @update:isDialogOpen="updateDialog"
        />
      </v-col>
    </v-row>
    <confirm-dialog ref="confirmAction" />
    <v-dialog v-model="removeMemberDialog" persistent max-width="1024px">
      <base-card>
        <template #header>
          <span class="headline">Remove Member</span>
        </template>
        <template #default>
          <v-row>
            <v-col>
              <v-autocomplete
                v-model="selectedMember"
                variant="outlined"
                density="compact"
                :items="items"
                :loading="loading"
                :search="search"
                label="Search"
                item-title="member_name"
                item-value="id"
                no-data-text="No results found"
                hide-no-data
                hide-selected
                return-object
                @update:model-value="displayUser"
                @update:search="onSearchUpdate"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="4">
              <v-text-field
                v-model="removeMemberName"
                variant="outlined"
                density="compact"
                label="Member Name"
                placeholder="Enter member name"
              ></v-text-field>
            </v-col>
            <v-col cols="4">
              <v-date-input
                v-model="removeMemberDateOfBirth"
                hide-actions
                locale="en-ZA"
                view-mode="month"
                prepend-icon=""
                prepend-inner-icon="$calendar"
                variant="outlined"
                density="compact"
                label="Date of Birth"
                placeholder="Select a date"
              ></v-date-input>
            </v-col>
            <v-col cols="4">
              <v-text-field
                v-model="removeMemberIdNumber"
                variant="outlined"
                density="compact"
                label="ID/Passport Number"
              ></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="4">
              <v-date-input
                v-model="selectedEffectiveExitDate"
                hide-actions
                locale="en-ZA"
                view-mode="month"
                prepend-icon=""
                prepend-inner-icon="$calendar"
                variant="outlined"
                density="compact"
                label="Effective Exit Date"
                placeholder="Select an exit date"
              ></v-date-input>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-btn color="primary" @click="removeMember"
            >Remove from Scheme</v-btn
          >
          <v-btn color="red" @click="handleCloseRemoveDialog">Cancel</v-btn>
        </template>
      </base-card>
    </v-dialog>
    <v-dialog v-model="editEndDateDialog" persistent max-width="400px">
      <base-card>
        <template #header>
          <span class="headline">Edit Cover End Date</span>
        </template>
        <template #default>
          <!-- <v-date-input
            v-model="newCoverEndDate"
            hide-actions
            locale="en-ZA"
            view-mode="month"
            variant="outlined"
            density="compact"
            label="Cover End Date"
            :min="scheme.cover_start_date"
          /> Remember to revert back to this: Note to future Motlatsi and self -->

          <v-date-input
            v-model="newCoverEndDate"
            hide-actions
            locale="en-ZA"
            view-mode="month"
            variant="outlined"
            density="compact"
            label="Cover End Date"
          />
        </template>
        <template #actions>
          <v-btn color="primary" @click="saveCoverEndDate">Save</v-btn>
          <v-btn color="red" @click="editEndDateDialog = false">Cancel</v-btn>
        </template>
      </base-card>
    </v-dialog>
    <v-dialog v-model="manageStatusDialog" persistent max-width="500px">
      <base-card>
        <template #header>
          <span class="headline">Manage Scheme Status</span>
        </template>
        <template #default>
          <v-row>
            <v-col cols="12">
              <v-card variant="tonal" class="pa-3 mb-4">
                <v-card-title class="text-subtitle-1 text-grey-darken-1">
                  Current Status
                </v-card-title>
                <v-card-text
                  class="text-h6 font-weight-bold text-grey-darken-4 pt-1"
                >
                  <v-chip
                    :color="statusColor"
                    label
                    size="large"
                    class="font-weight-bold"
                  >
                    {{ formatStatus(scheme.status) }}
                  </v-chip>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-select
                v-model="selectedStatus"
                variant="outlined"
                density="compact"
                label="New Status"
                placeholder="Select new status"
                :items="statusOptions"
                item-title="title"
                item-value="value"
                :rules="[(v) => !!v || 'Status is required']"
              ></v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-textarea
                v-model="schemeStatusMessage"
                variant="outlined"
                density="compact"
                label="Status Change Details"
                placeholder="Enter details about the status change..."
                rows="4"
                :rules="[(v) => !!v || 'Details are required']"
              ></v-textarea>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-btn color="primary" @click="updateSchemeStatus"
            >Update Status</v-btn
          >
          <v-btn color="red" @click="closeManageStatusDialog">Cancel</v-btn>
        </template>
      </base-card>
    </v-dialog>
    <v-dialog v-model="statusHistoryDialog" persistent max-width="900px">
      <base-card>
        <template #header>
          <span class="headline">Scheme Status History</span>
        </template>
        <template #default>
          <v-row v-if="loadingStatusHistory">
            <v-col class="text-center">
              <v-progress-circular
                indeterminate
                color="primary"
              ></v-progress-circular>
              <div class="mt-2">Loading status history...</div>
            </v-col>
          </v-row>
          <v-row v-else-if="statusHistory.length === 0">
            <v-col class="text-center">
              <v-icon size="64" color="grey">mdi-history</v-icon>
              <div class="text-h6 text-grey mt-2">No status changes found</div>
              <div class="text-body-2 text-grey"
                >This scheme has no recorded status changes.</div
              >
            </v-col>
          </v-row>
          <v-row v-else>
            <v-col>
              <v-data-table
                :headers="statusHistoryHeaders"
                :items="statusHistory"
                item-key="id"
                class="elevation-0"
                density="compact"
                :items-per-page="10"
              >
                <template #[`item.changed_at`]="{ item }">
                  {{ formatDateTime((item as any).changed_at) }}
                </template>
                <template #[`item.old_status`]="{ item }">
                  <v-chip
                    size="small"
                    :color="getStatusColor((item as any).old_status)"
                    variant="tonal"
                  >
                    {{ formatStatus((item as any).old_status) }}
                  </v-chip>
                </template>
                <template #[`item.new_status`]="{ item }">
                  <v-chip
                    size="small"
                    :color="getStatusColor((item as any).new_status)"
                    variant="tonal"
                  >
                    {{ formatStatus((item as any).new_status) }}
                  </v-chip>
                </template>
                <template #[`item.status_message`]="{ item }">
                  <div class="text-wrap" style="max-width: 300px">
                    {{ (item as any).status_message }}
                  </div>
                </template>
              </v-data-table>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-btn color="primary" @click="statusHistoryDialog = false"
            >Close</v-btn
          >
        </template>
      </base-card>
    </v-dialog>
    <v-dialog v-model="quotesDialog" persistent max-width="1200px">
      <base-card>
        <template #header>
          <span class="headline">Associated Quotes</span>
        </template>
        <template #default>
          <v-row v-if="loadingQuotes">
            <v-col class="text-center">
              <v-progress-circular
                indeterminate
                color="primary"
              ></v-progress-circular>
              <div class="mt-2">Loading quotes...</div>
            </v-col>
          </v-row>
          <v-row v-else-if="quotes.length === 0">
            <v-col class="text-center">
              <v-icon size="64" color="grey"
                >mdi-file-document-multiple-outline</v-icon
              >
              <div class="text-h6 text-grey mt-2">No quotes found</div>
              <div class="text-body-2 text-grey"
                >This scheme has no associated quotes.</div
              >
            </v-col>
          </v-row>
          <v-row v-else>
            <v-col>
              <v-data-table
                :headers="quotesHeaders"
                :items="quotes"
                item-key="id"
                class="elevation-0"
                density="compact"
                :items-per-page="10"
                style="cursor: pointer"
                @click:row="navigateToQuoteDetails"
              >
                <template #[`item.creation_date`]="{ item }">
                  {{ formatDate((item as any).creation_date) }}
                </template>
                <template #[`item.status`]="{ item }">
                  <v-chip
                    size="small"
                    :color="getQuoteStatusColor((item as any).status)"
                    variant="tonal"
                  >
                    {{ formatStatus((item as any).status) }}
                  </v-chip>
                </template>
                <template #[`item.quote_broker.name`]="{ item }">
                  {{ (item as any).quote_broker?.name || 'N/A' }}
                </template>
                <template #[`item.basis`]="{ item }">
                  <v-chip size="small" variant="outlined">
                    {{ formatBasis((item as any).basis) }}
                  </v-chip>
                </template>
                <template #[`item.member_data_count`]="{ item }">
                  <div class="text-center">
                    {{ (item as any).member_data_count }}
                  </div>
                </template>
              </v-data-table>
            </v-col>
          </v-row>
        </template>
        <template #actions>
          <v-btn color="primary" @click="quotesDialog = false">Close</v-btn>
        </template>
      </base-card>
    </v-dialog>

    <!-- Member Details Dialog -->
    <v-dialog v-model="memberDetailsDialog" persistent max-width="1200px">
      <base-card>
        <template #header>
          <div class="d-flex justify-space-between align-center">
            <span class="headline"
              >Member Details - {{ selectedSchemeMember?.member_name }}</span
            >
            <div>
              <v-btn
                size="small"
                rounded
                color="white"
                variant="outlined"
                class="mr-2"
                @click="editSchemeMember"
              >
                Edit Member
              </v-btn>
              <v-btn
                rounded
                size="small"
                color="error"
                variant="outlined"
                @click="deactivateSchemeMember"
              >
                Deactivate
              </v-btn>
            </div>
          </div>
        </template>
        <template #default>
          <div v-if="selectedSchemeMember" class="pa-4">
            <v-row>
              <v-col cols="12" md="6">
                <h3 class="text-h6 mb-3">Personal Information</h3>
                <div class="mb-2"
                  ><strong>Name:</strong>
                  {{ selectedSchemeMember.member_name }}</div
                >
                <div class="mb-2"
                  ><strong>ID Number:</strong>
                  {{ selectedSchemeMember.member_id_number }}</div
                >
                <div class="mb-2"
                  ><strong>Employee Number:</strong>
                  {{ selectedSchemeMember.employee_number || 'N/A' }}</div
                >
                <div class="mb-2"
                  ><strong>Gender:</strong>
                  {{ selectedSchemeMember.gender }}</div
                >
                <div class="mb-2"
                  ><strong>Date of Birth:</strong>
                  {{ formatDate(selectedSchemeMember.date_of_birth) }}</div
                >
                <div class="mb-2"
                  ><strong>Email:</strong>
                  {{ selectedSchemeMember.email || 'N/A' }}</div
                >
                <div class="mb-2"
                  ><strong>Phone:</strong>
                  {{ selectedSchemeMember.phone_number || 'N/A' }}</div
                >
              </v-col>
              <v-col cols="12" md="6">
                <h3 class="text-h6 mb-3">Scheme Information</h3>
                <div class="mb-2"
                  ><strong>Scheme:</strong>
                  {{ selectedSchemeMember.scheme_name }}</div
                >
                <div class="mb-2"
                  ><strong>Status:</strong>
                  <v-chip
                    size="small"
                    :color="
                      getMemberStatusColor(
                        selectedSchemeMember.status || 'active'
                      )
                    "
                  >
                    {{
                      (selectedSchemeMember.status || 'active').toUpperCase()
                    }}
                  </v-chip>
                </div>
                <div class="mb-2"
                  ><strong>Entry Date:</strong>
                  {{ formatDate(selectedSchemeMember.entry_date) }}</div
                >
                <div class="mb-2"
                  ><strong>Annual Salary:</strong>
                  {{ formatValues(selectedSchemeMember.annual_salary) }}</div
                >
                <div class="mb-2"
                  ><strong>Occupation:</strong>
                  {{ selectedSchemeMember.occupation || 'N/A' }}</div
                >
                <div class="mb-2"
                  ><strong>Occupational Class:</strong>
                  {{ selectedSchemeMember.occupational_class || 'N/A' }}</div
                >
              </v-col>
            </v-row>
          </div>
        </template>
        <template #actions>
          <v-btn color="grey" @click="memberDetailsDialog = false">Close</v-btn>
        </template>
      </base-card>
    </v-dialog>
  </v-container>
</template>
<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import GroupPricingService from '@/renderer/api/GroupPricingService'
// import formatValues from '@/renderer/utils/format_values'
import {
  formatValues,
  roundUpToTwoDecimalsAccounting
} from '@/renderer/utils/format_values'
// import { useForm } from 'vee-validate'
// import * as yup from 'yup'

import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import FileUploadDialog from '@/renderer/components/FileUploadDialog.vue'
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
import { VDateInput } from 'vuetify/labs/VDateInput'
import formatDateString from '@/renderer/utils/helpers.js'
import ContributionConfig from './premiums/ContributionConfig.vue'

const confirmAction = ref()
const backText = ref('< Back to listing')
const route = useRoute()
const router = useRouter()
const scheme: any = ref(null)
const schemes: any = ref([])
const activeSchemeTab = ref('contribution-config')
// const tableData: any = ref([])
// const selectedTable = ref('')
// const columnDefs: any = ref([])
// const genderItems = ['Male', 'Female']

const showModelPoint = ref(false)
const yearLabel = ref('') // 'Select a year'
const uploadTitle = ref('')
const mpLabel = ref('')
const isDialogOpen = ref(false)
const removeMemberDialog = ref(false)
const editEndDateDialog = ref(false)
const manageStatusDialog = ref(false)
const statusHistoryDialog = ref(false)
const quotesDialog = ref(false)
const selectedMember: any = ref(null)
const selectedEffectiveExitDate = ref<string | null>(null)

// Local state for remove member dialog
const removeMemberName = ref('')
const removeMemberDateOfBirth = ref(null)
const removeMemberIdNumber = ref('')
const newCoverEndDate = ref(scheme.value ? scheme.value.cover_end_date : null)
const selectedStatus = ref('')
const schemeStatusMessage = ref('')
const statusHistory: any = ref([])
const loadingStatusHistory = ref(false)
const quotes: any = ref([])
const loadingQuotes = ref(false)
const statusOptions = [
  { title: 'In Force', value: 'in_force' },
  { title: 'Terminated', value: 'terminated' },
  { title: 'Suspended', value: 'suspended' },
  { title: 'Lapsed', value: 'lapsed' },
  { title: 'Not Taken Up', value: 'not_taken_up' },
  { title: 'Cancelled', value: 'cancelled' },
  { title: 'Accepted', value: 'accepted' }
]

const statusHistoryHeaders = [
  { title: 'Date Changed', value: 'changed_at', key: 'changed_at' },
  { title: 'From Status', value: 'old_status', key: 'old_status' },
  { title: 'To Status', value: 'new_status', key: 'new_status' },
  { title: 'Details', value: 'status_message', key: 'status_message' },
  { title: 'Changed By', value: 'changed_by', key: 'changed_by' }
]

const quotesHeaders = [
  { title: 'Quote ID', value: 'id', key: 'id' },
  { title: 'Quote Type', value: 'quote_type', key: 'quote_type' },
  { title: 'Basis', value: 'basis', key: 'basis' },
  { title: 'Status', value: 'status', key: 'status' },
  {
    title: 'Scheme Status',
    value: 'scheme_quote_status',
    key: 'scheme_quote_status'
  },
  { title: 'Creation Date', value: 'creation_date', key: 'creation_date' },
  { title: 'Created By', value: 'created_by', key: 'created_by' },
  {
    title: 'Member Count',
    value: 'member_data_count',
    key: 'member_data_count'
  },
  {
    title: 'Channel',
    value: 'distribution_channel',
    key: 'distribution_channel'
  },
  { title: 'Broker', value: 'quote_broker.name', key: 'broker' }
]

// Member management state (similar to MemberManagement.vue)
const showMemberManagement = ref(false)
const loadingMembers = ref(false)
const memberLoadingProgress = ref(0)
const memberLoadingMessage = ref('')
const memberSearchQuery = ref('')
const selectedMemberStatus = ref<string | null>(null)
const members = ref<any[]>([])
const selectedSchemeMember = ref<any | null>(null)
const memberDetailsDialog = ref(false)
const totalSchemeMembers = ref(0)
const hasMoreSchemeMembers = ref(true)
const hasMemberSearched = ref(false)
const memberLoadAbortController = ref<AbortController | null>(null)

// Server-side filters for members
const memberServerFilters = ref({
  search: '',
  status: null as string | null,
  page: 1,
  pageSize: 100
})

// Filter options for members
const memberStatuses = [
  { title: 'ACTIVE', value: 'ACTIVE' },
  { title: 'INACTIVE', value: 'INACTIVE' },
  { title: 'PENDING', value: 'PENDING' },
  { title: 'SUSPENDED', value: 'SUSPENDED' }
]

// Column definitions for member data grid
const memberColumnDefs = [
  {
    headerName: 'Member Name',
    field: 'member_name',
    filter: true,
    sortable: true,
    minWidth: 200
  },
  {
    headerName: 'Employee Number',
    field: 'employee_number',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'ID Number',
    field: 'member_id_number',
    filter: true,
    sortable: true,
    minWidth: 150
  },
  {
    headerName: 'Status',
    field: 'status',
    filter: true,
    sortable: true,
    minWidth: 120,
    cellRenderer: (params: any) => {
      const status = params.value || 'active'
      const color = getMemberStatusColor(status)
      return `<v-chip size="small" color="${color}">${status.toUpperCase()}</v-chip>`
    }
  },
  {
    headerName: 'Annual Salary',
    field: 'annual_salary',
    filter: true,
    sortable: true,
    valueFormatter: formatValues,
    minWidth: 150
  },
  {
    headerName: 'Entry Date',
    field: 'entry_date',
    filter: true,
    sortable: true,
    minWidth: 120,
    valueFormatter: (params: any) => {
      return params.value ? new Date(params.value).toLocaleDateString() : ''
    }
  },
  {
    headerName: 'Gender',
    field: 'gender',
    filter: true,
    sortable: true,
    minWidth: 100
  },
  {
    headerName: 'Actions',
    field: 'actions',
    sortable: false,
    filter: false,
    minWidth: 120,
    cellRenderer: () => {
      return `
        <v-btn size="small" color="primary" variant="text">
          View Details
        </v-btn>
      `
    }
  }
]

// Computed property for member pagination info
const memberPaginationInfo = computed(() => ({
  currentPage: Math.ceil(
    members.value.length / memberServerFilters.value.pageSize
  ),
  totalPages: Math.ceil(
    totalSchemeMembers.value / memberServerFilters.value.pageSize
  ),
  totalMembers: totalSchemeMembers.value,
  displayedMembers: members.value.length,
  hasMore: hasMoreSchemeMembers.value
}))
// Member form removed - now handled in MemberManagement.vue
const years = ref<number[]>(
  Array.from({ length: 10 }, (v, k) => new Date().getFullYear() - k)
)
const updateDialog = (value: boolean) => {
  isDialogOpen.value = value
}

const snackbar = ref(false)
const timeout = 2000
const snackbarText = ref('')

// const parseDateString = (dateString) => {
//   const date = new Date(dateString)
//   const formattedDate = date.toISOString().split('T')[0]
//   return formattedDate
// }

const relatedTables = computed(() => {
  const tables: any = []
  tables.push({
    table_type: 'Member Data',
    value: 'member_data',
    populated: true
  })
  return tables
})

const goBack = () => {
  router.go(-1)
}

// A computed property to dynamically set the status chip color.
const statusColor = computed(() => {
  switch (scheme.value.status) {
    case 'Out Of Force':
      return 'red-darken-1'
    case 'Active':
      return 'green-darken-1'
    case 'Pending':
      return 'orange-darken-1'
    default:
      return 'grey'
  }
})

const removeMember = async () => {
  if (!selectedMember.value) {
    snackbarText.value = 'Please select a member to deactivate'
    snackbar.value = true
    return
  }

  // Show confirmation dialog
  const result = await confirmAction.value.open(
    'Deactivate Scheme Member',
    `Are you sure you want to deactivate ${selectedMember.value.member_name} with effective exit date ${selectedEffectiveExitDate.value || 'not specified'}?`
  )

  // If user canceled, don't proceed
  if (!result) {
    return
  }

  console.log(
    'Deactivating member:',
    selectedMember.value,
    'with exit date:',
    selectedEffectiveExitDate.value
  )
  selectedMember.value.effective_exit_date = selectedEffectiveExitDate.value
  GroupPricingService.removeMemberFromScheme(
    scheme.value.id,
    selectedMember.value
  )
    .then((res) => {
      snackbarText.value = 'Member deactivated successfully'
      snackbar.value = true
      items.value = []
      selectedMember.value = null
      removeMemberName.value = ''
      removeMemberDateOfBirth.value = null
      removeMemberIdNumber.value = ''
      selectedEffectiveExitDate.value = null
      removeMemberDialog.value = false

      // Refresh the scheme data
      GroupPricingService.getScheme(route.params.id).then((response) => {
        scheme.value = response.data
        schemes.value = [scheme.value]
      })
    })
    .catch((error) => {
      console.log('Error:', error)
      snackbarText.value = 'Failed to deactivate member'
      snackbar.value = true
      removeMemberDialog.value = false
    })
}

// Load members for the current scheme using sophisticated approach
const loadMembersForScheme = () => {
  showMemberManagement.value = true
  // Don't auto-load members, let user search/filter first
}

const hideMemberManagement = () => {
  showMemberManagement.value = false
  members.value = []
  hasMemberSearched.value = false
}

const loadSchemeMembers = async (append: boolean = false) => {
  // Cancel any existing load operation
  if (memberLoadAbortController.value) {
    memberLoadAbortController.value.abort()
  }

  memberLoadAbortController.value = new AbortController()
  loadingMembers.value = true
  memberLoadingProgress.value = 0
  memberLoadingMessage.value = 'Loading scheme members...'

  try {
    let membersResponse
    let newMembers, total, hasMore

    try {
      // Try the new paginated API first for this specific scheme
      membersResponse = await GroupPricingService.getMembersPaginated({
        page: memberServerFilters.value.page,
        pageSize: memberServerFilters.value.pageSize,
        search: memberServerFilters.value.search,
        schemeId: scheme.value.id, // Filter by current scheme
        status: memberServerFilters.value.status,
        signal: memberLoadAbortController.value.signal
      })

      const response = membersResponse.data
      newMembers = response.data
      total = response.total
      hasMore = response.hasMore
    } catch (paginationError) {
      // Fallback to scheme-specific member loading
      console.warn(
        'Paginated API not available, falling back to scheme member API'
      )
      memberLoadingMessage.value = 'Loading members for scheme...'

      const response = await GroupPricingService.getMembersInForce(
        scheme.value.id
      )
      let allMembers = response.data.map((member: any) => ({
        ...member,
        scheme_name: scheme.value.name,
        scheme_id: scheme.value.id,
        status: member.status || 'active'
      }))

      // Apply client-side filtering
      if (memberServerFilters.value.search) {
        const query = memberServerFilters.value.search.toLowerCase()
        allMembers = allMembers.filter(
          (member: any) =>
            member.member_name?.toLowerCase().includes(query) ||
            member.member_id_number?.toLowerCase().includes(query)
        )
      }

      if (memberServerFilters.value.status) {
        allMembers = allMembers.filter(
          (member: any) => member.status === memberServerFilters.value.status
        )
      }

      // Apply pagination to client-filtered results
      const startIndex =
        (memberServerFilters.value.page - 1) *
        memberServerFilters.value.pageSize
      const endIndex = startIndex + memberServerFilters.value.pageSize

      newMembers = allMembers.slice(startIndex, endIndex)
      total = allMembers.length
      hasMore = endIndex < allMembers.length
    }

    if (append) {
      members.value = [...members.value, ...newMembers]
    } else {
      members.value = newMembers
    }

    totalSchemeMembers.value = total
    hasMoreSchemeMembers.value = hasMore
    memberLoadingProgress.value = 100
    memberLoadingMessage.value = 'Loading complete'
  } catch (error: any) {
    if (error.name === 'AbortError') {
      memberLoadingMessage.value = 'Loading cancelled'
      return
    }
    console.error('Error loading scheme members:', error)
    snackbarText.value = 'Error loading scheme members'
    snackbar.value = true
  } finally {
    loadingMembers.value = false
    memberLoadAbortController.value = null
  }
}

// Load more members (for pagination)
const loadMoreSchemeMembers = async () => {
  if (!hasMoreSchemeMembers.value || loadingMembers.value) return

  memberServerFilters.value.page += 1
  await loadSchemeMembers(true)
}

// Reset and reload with new filters
const reloadSchemeMembers = async () => {
  memberServerFilters.value.page = 1
  hasMemberSearched.value = true
  await loadSchemeMembers(false)
}

// Cancel current loading operation
const cancelMemberLoading = () => {
  if (memberLoadAbortController.value) {
    memberLoadAbortController.value.abort()
    loadingMembers.value = false
    memberLoadingMessage.value = 'Loading cancelled'
  }
}

// Handle member row click
const handleMemberRowClick = (event: any) => {
  console.log('Member row clicked:', event.data)
  selectedSchemeMember.value = event.data
  memberDetailsDialog.value = true
}

// Edit scheme member
const editSchemeMember = () => {
  if (!selectedSchemeMember.value) return
  router.push({
    name: 'group-pricing-member-management',
    query: {
      schemeId: scheme.value.id,
      openMemberDetails: selectedSchemeMember.value.id
    }
  })
}

// Deactivate scheme member
const deactivateSchemeMember = async () => {
  if (!selectedSchemeMember.value) return

  const result = await confirmAction.value.open(
    'Deactivate Scheme Member',
    `Are you sure you want to deactivate ${selectedSchemeMember.value.member_name}?`
  )

  if (!result) return

  try {
    const memberToUpdate = {
      ...selectedSchemeMember.value,
      effective_exit_date: new Date().toISOString().split('T')[0],
      status: 'INACTIVE'
    }
    await GroupPricingService.removeMemberFromScheme(
      scheme.value.id,
      memberToUpdate
    )
    snackbarText.value = 'Member deactivated successfully'
    snackbar.value = true
    await reloadSchemeMembers()
    memberDetailsDialog.value = false
  } catch (error) {
    console.error('Error deactivating member:', error)
    snackbarText.value = 'Error deactivating member'
    snackbar.value = true
  }
}

// Debounced search for members
let memberSearchTimeout: ReturnType<typeof setTimeout> | null = null
const searchSchemeMembers = () => {
  if (memberSearchTimeout) clearTimeout(memberSearchTimeout)

  memberSearchTimeout = setTimeout(async () => {
    memberServerFilters.value.search = memberSearchQuery.value
    await reloadSchemeMembers()
  }, 500) // 500ms delay
}

const filterMembersByStatus = async () => {
  memberServerFilters.value.status = selectedMemberStatus.value
  await reloadSchemeMembers()
}

const exportSchemeMembers = () => {
  // Implementation for exporting scheme members to CSV/Excel
  const csvData = members.value.map((member) => ({
    'Member Name': member.member_name,
    'ID Number': member.member_id_number,
    'Employee Number': member.employee_number || '',
    Status: member.status,
    'Annual Salary': member.annual_salary,
    'Entry Date': member.entry_date,
    Gender: member.gender
  }))

  const csvContent =
    'data:text/csv;charset=utf-8,' +
    Object.keys(csvData[0]).join(',') +
    '\n' +
    csvData.map((row) => Object.values(row).join(',')).join('\n')

  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute(
    'download',
    `scheme_${scheme.value.id}_members_export_${new Date().toISOString().split('T')[0]}.csv`
  )
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const getMemberStatusColor = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'error'
    case 'pending':
      return 'warning'
    case 'suspended':
      return 'orange'
    default:
      return 'grey'
  }
}

const openMemberDialog = (item: any) => {
  router.push({
    name: 'group-pricing-member-management',
    query: { schemeId: scheme.value.id }
  })
}

// const openRemoveMemberDialog = (item: any) => {
//   // Reset member values
//   selectedMember.value = null
//   removeMemberName.value = ''
//   removeMemberDateOfBirth.value = null
//   removeMemberIdNumber.value = ''

//   removeMemberDialog.value = true
// }

// function parseDateStringToDate(dateString) {
//   if (!dateString) return null
//   // Handles 'yyyy-mm-dd' format
//   const [year, month, day] = dateString.split('-').map(Number)
//   return new Date(year, month - 1, day)
// }

function parseDateStringToDate(dateString) {
  if (!dateString) return null
  const [year, month, day] = dateString.split('-').map(Number)
  return new Date(year, month - 1, day)
}

function formatDateToYMD(dateObj) {
  if (!dateObj) return ''
  const year = dateObj.getFullYear()
  const month = String(dateObj.getMonth() + 1).padStart(2, '0')
  const day = String(dateObj.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const openEditEndDateDialog = () => {
  newCoverEndDate.value = parseDateStringToDate(scheme.value.cover_end_date)
  editEndDateDialog.value = true
}

async function saveCoverEndDate() {
  if (!newCoverEndDate.value) return
  try {
    const formattedDate = formatDateToYMD(newCoverEndDate.value)
    await GroupPricingService.updateSchemeCoverEndDate(
      scheme.value.id,
      formattedDate
    )
    scheme.value.cover_end_date = formattedDate // update value on screen
    snackbarText.value = 'Cover End Date updated successfully'
    snackbar.value = true
    editEndDateDialog.value = false
  } catch (err) {
    console.error('Error deleting broker:', err)
    snackbarText.value = 'Failed to Update Cover End Date'
    snackbar.value = true
  }
}

const openManageStatusDialog = () => {
  selectedStatus.value = scheme.value.status
  schemeStatusMessage.value = ''
  manageStatusDialog.value = true
}

const closeManageStatusDialog = () => {
  selectedStatus.value = ''
  schemeStatusMessage.value = ''
  manageStatusDialog.value = false
}

const updateSchemeStatus = async () => {
  if (!selectedStatus.value || !schemeStatusMessage.value) {
    snackbarText.value = 'Please select a status and provide details'
    snackbar.value = true
    return
  }

  if (selectedStatus.value === scheme.value.status) {
    snackbarText.value = 'The selected status is the same as the current status'
    snackbar.value = true
    return
  }

  try {
    const statusUpdatePayload = {
      status: selectedStatus.value,
      scheme_status_message: schemeStatusMessage.value,
      status_change_date: new Date().toISOString().split('T')[0]
    }

    await GroupPricingService.updateSchemeStatus(
      scheme.value.id,
      statusUpdatePayload
    )

    // Update the scheme status on the UI
    scheme.value.status = selectedStatus.value

    snackbarText.value = 'Scheme status updated successfully'
    snackbar.value = true
    closeManageStatusDialog()
  } catch (err) {
    console.error('Failed to update scheme status:', err)
    snackbarText.value = 'Failed to update scheme status'
    snackbar.value = true
  }
}

const openStatusHistoryDialog = async () => {
  statusHistoryDialog.value = true
  loadingStatusHistory.value = true

  try {
    const response = await GroupPricingService.getSchemeStatusHistory(
      scheme.value.id
    )
    statusHistory.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch status history:', err)
    snackbarText.value = 'Failed to load status history'
    snackbar.value = true
    statusHistory.value = []
  } finally {
    loadingStatusHistory.value = false
  }
}

const formatDateTime = (dateTimeString: string) => {
  if (!dateTimeString) return ''
  const date = new Date(dateTimeString)
  return date.toLocaleString('en-ZA', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusColor = (status: string) => {
  switch (status) {
    case 'in_force':
      return 'green-darken-1'
    case 'terminated':
      return 'red-darken-1'
    case 'suspended':
      return 'orange-darken-1'
    case 'lapsed':
      return 'red-lighten-1'
    case 'not_taken_up':
      return 'grey-darken-1'
    case 'cancelled':
      return 'red-darken-2'
    case 'accepted':
      return 'blue-darken-1'
    case 'out_of_force':
      return 'grey'
    default:
      return 'grey'
  }
}

const openQuotesDialog = async () => {
  quotesDialog.value = true
  loadingQuotes.value = true

  try {
    const response = await GroupPricingService.getSchemeQuotes(scheme.value.id)
    quotes.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch quotes:', err)
    snackbarText.value = 'Failed to load quotes'
    snackbar.value = true
    quotes.value = []
  } finally {
    loadingQuotes.value = false
  }
}

const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const getQuoteStatusColor = (status: string) => {
  switch (status) {
    case 'in_force':
      return 'green-darken-1'
    case 'draft':
      return 'orange-darken-1'
    case 'submitted':
      return 'blue-darken-1'
    case 'approved':
      return 'green'
    case 'rejected':
      return 'red-darken-1'
    case 'expired':
      return 'grey-darken-1'
    default:
      return 'grey'
  }
}

const formatBasis = (basis: string) => {
  if (!basis) return ''
  return basis
    .replace(/_/g, ' ')
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}

const navigateToQuoteDetails = (event: any, { item }: { item: any }) => {
  // Close the quotes dialog first
  quotesDialog.value = false

  // Navigate to quote details page with the quote ID
  router.push({
    name: 'group-pricing-scheme-details',
    params: { id: item.id }
  })
}

onMounted(() => {
  GroupPricingService.getScheme(route.params.id).then((response) => {
    scheme.value = response.data
    schemes.value.push(scheme.value)
    console.log('Loaded scheme:', scheme.value)
  })
})

// const clearData = () => {
//   tableData.value = []
//   selectedTable.value = ''
// }

const handleUpload = async (payload: any) => {
  // const formdata = new FormData()
  // formdata.append('file', payload.file)
  // formdata.append('quote_id', quote.value.id)
  // formdata.append('table_type', selectedTable.value.table_type)
  // GroupPricingService.uploadQuoteTable(formdata)
  //   .then((res) => {
  //     console.log('Response:', res.data)
  //     const count = res.data
  //     snackbarText.value = 'Upload Successful'
  //     snackbar.value = true
  //     if (selectedTable.value.table_type === 'Member Data') {
  //       quote.value.member_data_count = count
  //     } else if (selectedTable.value.table_type === 'Claims Experience') {
  //       quote.value.claims_experience_count = count
  //     } else if (selectedTable.value.table_type === 'Member Rating Results') {
  //       quote.value.member_rating_result_count = count
  //     } else if (selectedTable.value.table_type === 'Member Premium Schedules') {
  //       quote.value.member_premium_schedule_count = count
  //     }
  //   })
  //   .catch((error) => {
  //     console.log('Error:', error)
  //     snackbarText.value = 'Upload Failed'
  //     snackbar.value = true
  //   })
}

// const createColumnDefs = (data: any) => {
//   columnDefs.value = []
//   Object.keys(data[0]).forEach((element) => {
//     const header: any = {}
//     header.headerName = element
//     header.field = element
//     header.valueFormatter = formatValues
//     header.minWidth = 200
//     header.sortable = true
//     header.filter = true
//     header.resizable = true
//     columnDefs.value.push(header)
//   })
// }

/// test code
// State
const search = ref('')
const items = ref<any[]>([])
const loading = ref(false)
let searchTimeout: any = null

// Watch search input and debounce API call
const onSearchUpdate = (val: string) => {
  search.value = val
  if (searchTimeout) clearTimeout(searchTimeout)

  // Debounce API call
  searchTimeout = setTimeout(() => {
    fetchItems(val)
  }, 500)
}

// Simulated API call
const fetchItems = async (query: string) => {
  if (!query) {
    items.value = []
    return
  }

  loading.value = true
  try {
    // Replace this with your actual API call
    const response = await GroupPricingService.searchMembers(
      scheme.value.id,
      scheme.value.quote_id,
      query
    )
    const result = response.data
    items.value = result || []
  } catch (err) {
    console.error('Error fetching items:', err)
    items.value = []
  } finally {
    loading.value = false
  }
}

const displayUser = (val: any) => {
  if (val) {
    removeMemberName.value = val.member_name
    removeMemberDateOfBirth.value = val.date_of_birth
    removeMemberIdNumber.value = val.member_id_number
    // Set current date as default effective exit date
    selectedEffectiveExitDate.value = new Date().toISOString().substr(0, 10)
  }
}

const handleCloseRemoveDialog = () => {
  console.log('handleCloseRemoveDialog called')
  // Reset all fields
  selectedMember.value = null
  removeMemberName.value = ''
  removeMemberDateOfBirth.value = null
  removeMemberIdNumber.value = ''

  // Close the dialog
  removeMemberDialog.value = false
  console.log(
    'Dialog should be closed now, removeMemberDialog =',
    removeMemberDialog.value
  )
}

const deleteScheme = async () => {
  const result = await confirmAction.value.open(
    'Delete Scheme',
    `Are you sure you want to delete the scheme "${scheme.value.name}"? This action cannot be undone.`
  )

  if (!result) {
    return
  }

  try {
    await GroupPricingService.deleteScheme(scheme.value.id)
    snackbarText.value = 'Scheme deleted successfully'
    snackbar.value = true
    setTimeout(() => {
      router.push('/group-pricing/schemes')
    }, 1000)
  } catch (err) {
    console.error('Failed to delete scheme:', err)
    snackbarText.value = 'Failed to delete scheme'
    snackbar.value = true
  }
}

// Add this method in <script setup>
function formatStatus(status) {
  if (!status) return ''
  return status
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ')
}
</script>
<style scoped>
.table-row {
  white-space: nowrap;
}

::v-deep(.v-data-table thead th) {
  background-color: #223f54 !important;
  color: white;
  text-align: center;
  font-weight: bold;
  white-space: nowrap;
  min-width: 150px;
}

.search-box {
  width: 100%;
}
.v-table__wrapper > table > thead {
  background-color: #223f54 !important;
  color: white;
  white-space: nowrap;
}

.value-text {
  display: inline-block;
  text-align: left;
  vertical-align: top;
  min-width: 1px;
  width: auto;
  white-space: pre-line;
}
</style>
