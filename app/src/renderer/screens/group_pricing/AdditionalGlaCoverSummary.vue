<template>
  <base-card :show-actions="false">
    <template #header>
      <div class="d-flex align-center" style="width: 100%; gap: 16px">
        <span class="headline">{{ benefitTitle }} Summary</span>
        <v-tabs
          v-if="categories.length > 0"
          v-model="viewTab"
          color="white"
          slider-color="white"
          density="compact"
          class="agla-header-tabs"
          hide-slider
        >
          <v-tab value="smoothing">Smoothing</v-tab>
          <v-tab value="final">Final results</v-tab>
        </v-tabs>
        <v-spacer />
        <v-btn
          v-if="hasAnyRates"
          size="small"
          variant="flat"
          color="white"
          class="text-primary"
          prepend-icon="mdi-download"
          @click="exportSummaryToExcel"
        >
          Export All
        </v-btn>
      </div>
    </template>
    <template #default>
      <v-alert
        v-if="categories.length === 0"
        type="info"
        density="compact"
        variant="tonal"
      >
        {{ benefitTitle }} is not enabled on any scheme category for this quote.
      </v-alert>
      <template v-else>
        <v-tabs
          v-model="activeTab"
          color="primary"
          density="compact"
          show-arrows
        >
          <v-tab
            v-for="(rs, i) in categories"
            :key="'agla-cat-tab-' + i"
            :value="rs.category"
          >
            {{ rs.category }}
          </v-tab>
        </v-tabs>
        <v-window v-model="activeTab" class="mt-3">
          <v-window-item
            v-for="(rs, i) in categories"
            :key="'agla-cat-pane-' + i"
            :value="rs.category"
          >
            <template
              v-if="(rs.additional_gla_cover_band_rates?.length ?? 0) > 0"
            >
              <div class="text-caption mb-2">
                <span v-if="rs.additional_gla_cover_age_band_type">
                  Age bands:
                  <em>{{ rs.additional_gla_cover_age_band_type }}</em>
                </span>
                <span v-if="rs.additional_gla_cover_male_prop_used != null">
                  &middot; Male proportion used:
                  {{
                    (
                      (rs.additional_gla_cover_male_prop_used ?? 0) * 100
                    ).toFixed(1)
                  }}%
                </span>
              </div>

              <template v-if="viewTab === 'smoothing'">
                <!-- ===== Top: Smoothed comparison table ===== -->
                <v-alert
                  v-if="isLocked"
                  type="info"
                  density="compact"
                  variant="tonal"
                  icon="mdi-lock"
                  class="mt-2 mb-2"
                >
                  Smoothed rates are locked because this quote is
                  <strong>{{ formatStatus(quoteStatus) }}</strong
                  >. Any further changes require reverting the quote status.
                </v-alert>
                <div
                  v-if="getMetadata(rs.category)"
                  class="text-caption text-medium-emphasis mt-2"
                >
                  <v-icon size="x-small" class="mr-1">mdi-history</v-icon>
                  Last updated
                  <span v-if="getMetadata(rs.category)?.updatedBy">
                    by
                    <strong>{{ getMetadata(rs.category)?.updatedBy }}</strong>
                  </span>
                  <span v-if="getMetadata(rs.category)?.updatedAt">
                    at
                    {{ formatTimestamp(getMetadata(rs.category)?.updatedAt) }}
                  </span>
                </div>
                <div
                  class="d-flex align-center mt-2 mb-1 flex-wrap"
                  style="gap: 8px"
                >
                  <div class="text-subtitle-2">
                    Smoothed Office Rate / 1,000 — comparison
                  </div>
                  <v-chip
                    v-if="hasUnsavedChanges(rs.category)"
                    size="x-small"
                    color="warning"
                    variant="tonal"
                    >Unsaved changes</v-chip
                  >
                  <v-spacer />
                  <v-btn
                    size="x-small"
                    variant="text"
                    prepend-icon="mdi-download"
                    @click="downloadTemplate(rs)"
                  >
                    Download Template
                  </v-btn>
                  <v-btn
                    size="x-small"
                    variant="text"
                    prepend-icon="mdi-upload"
                    :disabled="isReadonly(rs.category)"
                    @click="triggerUpload(rs.category)"
                  >
                    Upload Smoothed
                  </v-btn>
                  <input
                    :ref="(el) => bindUploadInput(rs.category, el as any)"
                    type="file"
                    accept=".xlsx"
                    style="display: none"
                    @change="onUploadFile($event, rs.category)"
                  />
                  <v-btn
                    size="x-small"
                    variant="text"
                    color="warning"
                    prepend-icon="mdi-restore"
                    :disabled="
                      isReadonly(rs.category) || !hasUnsavedChanges(rs.category)
                    "
                    @click="resetCategory(rs.category)"
                  >
                    Reset
                  </v-btn>
                  <!-- Save/Edit toggle: once a category has persisted smoothed
                     values, the table drops to read-only and the primary
                     button reads "Edit". Clicking it re-opens the inputs;
                     the next Save commits and flips back. -->
                  <v-btn
                    v-if="isReadonly(rs.category) && !isLocked"
                    size="x-small"
                    variant="flat"
                    color="primary"
                    prepend-icon="mdi-pencil"
                    @click="enterEditMode(rs.category)"
                  >
                    Edit
                  </v-btn>
                  <v-btn
                    v-else
                    size="x-small"
                    variant="flat"
                    color="primary"
                    prepend-icon="mdi-content-save"
                    :disabled="
                      isLocked ||
                      !hasUnsavedChanges(rs.category) ||
                      saving[rs.category]
                    "
                    :loading="saving[rs.category]"
                    @click="saveCategory(rs.category)"
                  >
                    Save
                  </v-btn>
                </div>

                <div style="overflow-x: auto">
                  <v-table density="compact" class="agla-summary-table">
                    <thead>
                      <tr>
                        <th rowspan="2" class="agla-age-band">Age Band</th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Office Rate / 1,000
                        </th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Weighted by GLA Covered SA
                        </th>
                        <th colspan="3" class="text-center">
                          Smoothed Office Rate / 1,000
                        </th>
                      </tr>
                      <tr>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th class="text-right agla-combined-cell">Combined</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(r, j) in rs.additional_gla_cover_band_rates"
                        :key="'agla-cmp-' + i + '-' + j"
                      >
                        <td class="agla-age-band">{{ formatBandLabel(r) }}</td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              getOfficeRate(r, 'male')
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              getOfficeRate(r, 'female')
                            )
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            roundUpToTwoDecimalsAccounting(
                              getOfficeRate(r, 'combined')
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            formatNullable(r.weighted_office_rate_per1000_male)
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            formatNullable(
                              r.weighted_office_rate_per1000_female
                            )
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{ formatNullable(r.weighted_office_rate_per1000) }}
                        </td>
                        <td class="text-right agla-mf-cell pa-0">
                          <input
                            class="smoothed-input"
                            type="number"
                            step="0.01"
                            :readonly="isReadonly(rs.category)"
                            :value="effectiveSmoothed(rs.category, j, 'male')"
                            @input="
                              (e) =>
                                onSmoothedInput(
                                  rs.category,
                                  j,
                                  'male',
                                  (e.target as HTMLInputElement).value
                                )
                            "
                          />
                        </td>
                        <td class="text-right agla-mf-cell pa-0">
                          <input
                            class="smoothed-input"
                            type="number"
                            step="0.01"
                            :readonly="isReadonly(rs.category)"
                            :value="effectiveSmoothed(rs.category, j, 'female')"
                            @input="
                              (e) =>
                                onSmoothedInput(
                                  rs.category,
                                  j,
                                  'female',
                                  (e.target as HTMLInputElement).value
                                )
                            "
                          />
                        </td>
                        <td class="text-right agla-combined-cell pa-0">
                          <input
                            class="smoothed-input"
                            type="number"
                            step="0.01"
                            :readonly="isReadonly(rs.category)"
                            :value="
                              effectiveSmoothed(rs.category, j, 'combined')
                            "
                            @input="
                              (e) =>
                                onSmoothedInput(
                                  rs.category,
                                  j,
                                  'combined',
                                  (e.target as HTMLInputElement).value
                                )
                            "
                          />
                        </td>
                      </tr>
                      <tr class="factor-row">
                        <td class="agla-age-band text-medium-emphasis">
                          Smoothing factor
                        </td>
                        <td colspan="6" class="text-medium-emphasis">
                          Apply a multiplier to OfficeRate/1000 (1.00 = no
                          change). Direct edits to a Smoothed cell override the
                          factor for that cell.
                        </td>
                        <td class="text-right pa-0">
                          <input
                            class="factor-input"
                            type="number"
                            step="0.01"
                            :readonly="isReadonly(rs.category)"
                            :value="
                              getCategoryFactor(rs.category, 'male') ?? ''
                            "
                            placeholder="1.00"
                            @input="
                              (e) =>
                                onFactorInput(
                                  rs.category,
                                  'male',
                                  (e.target as HTMLInputElement).value
                                )
                            "
                          />
                        </td>
                        <td class="text-right pa-0">
                          <input
                            class="factor-input"
                            type="number"
                            step="0.01"
                            :readonly="isReadonly(rs.category)"
                            :value="
                              getCategoryFactor(rs.category, 'female') ?? ''
                            "
                            placeholder="1.00"
                            @input="
                              (e) =>
                                onFactorInput(
                                  rs.category,
                                  'female',
                                  (e.target as HTMLInputElement).value
                                )
                            "
                          />
                        </td>
                        <td class="text-right pa-0">
                          <input
                            class="factor-input"
                            type="number"
                            step="0.01"
                            :readonly="isReadonly(rs.category)"
                            :value="
                              getCategoryFactor(rs.category, 'combined') ?? ''
                            "
                            placeholder="1.00"
                            @input="
                              (e) =>
                                onFactorInput(
                                  rs.category,
                                  'combined',
                                  (e.target as HTMLInputElement).value
                                )
                            "
                          />
                        </td>
                      </tr>
                    </tbody>
                  </v-table>
                </div>

                <!-- ===== Middle: comparison chart ===== -->
                <div class="d-flex align-center mt-4 mb-1">
                  <div class="text-subtitle-2">
                    Office Rate / 1,000 — Original vs Smoothed
                  </div>
                  <v-spacer />
                  <ChartMenu
                    :chart-ref="chartRefs[rs.category]"
                    :title="`${benefitTitle} — ${rs.category} — Office vs Smoothed`"
                    :data="chartDataFor(rs)"
                  />
                </div>
                <ag-charts
                  :ref="(el) => bindChartRef(rs.category, el)"
                  :options="chartOptionsFor(rs)"
                />

                <!-- ===== Bottom: existing detailed breakdown ===== -->
                <div class="text-subtitle-2 mt-4 mb-1">
                  Detailed breakdown (current)
                </div>
                <div style="overflow-x: auto">
                  <v-table density="compact" class="agla-summary-table">
                    <thead>
                      <tr>
                        <th rowspan="2" class="agla-age-band">Age Band</th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Risk Rate / 1,000
                        </th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Binder Fee / 1,000
                        </th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Outsource Fee / 1,000
                        </th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Commission / 1,000
                        </th>
                        <th colspan="3" class="text-center">
                          Office Rate / 1,000
                        </th>
                      </tr>
                      <tr>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th class="text-right agla-combined-cell">Combined</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(r, j) in rs.additional_gla_cover_band_rates"
                        :key="'agla-detail-' + i + '-' + j"
                      >
                        <td class="agla-age-band">{{ formatBandLabel(r) }}</td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.risk_rate_per1000_male ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.risk_rate_per1000_female ?? 0
                            )
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            roundUpToTwoDecimalsAccounting(r.risk_rate_per1000)
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.binder_fee_per1000_male ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.binder_fee_per1000_female ?? 0
                            )
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.binder_fee_per1000 ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.outsource_fee_per1000_male ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.outsource_fee_per1000_female ?? 0
                            )
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.outsource_fee_per1000 ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              (r.office_rate_per1000_male ?? 0) *
                                schemeCommissionRate
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              (r.office_rate_per1000_female ?? 0) *
                                schemeCommissionRate
                            )
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            roundUpToTwoDecimalsAccounting(
                              (r.office_rate_per1000 ?? 0) *
                                schemeCommissionRate
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.office_rate_per1000_male ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.office_rate_per1000_female ?? 0
                            )
                          }}
                        </td>
                        <td class="text-right agla-combined-cell">
                          {{
                            roundUpToTwoDecimalsAccounting(
                              r.office_rate_per1000
                            )
                          }}
                        </td>
                      </tr>
                    </tbody>
                  </v-table>
                </div>
              </template>

              <!-- ===== Final results view ===== -->
              <template v-else>
                <div class="text-subtitle-2 mt-2 mb-1">
                  Final results (Smoothed Office Rate / 1,000 applied)
                </div>
                <div class="text-caption text-medium-emphasis mb-2">
                  Binder and outsource per 1,000 are projected by applying
                  each band's category-level fee rate to the Smoothed Office
                  Rate. Commission per 1,000 is the Smoothed Office Rate
                  multiplied by the scheme-wide ratio of final commission to
                  final premium, so the rate reconciles to the scheme totals
                  shown elsewhere in the quote.
                </div>
                <div style="overflow-x: auto">
                  <v-table density="compact" class="agla-summary-table">
                    <thead>
                      <tr>
                        <th rowspan="2" class="agla-age-band">Age Band</th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Binder Fee / 1,000
                        </th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Outsource Fee / 1,000
                        </th>
                        <th colspan="3" class="text-center agla-group-divider">
                          Commission / 1,000
                        </th>
                        <th colspan="3" class="text-center">
                          Final Office Rate / 1,000
                        </th>
                      </tr>
                      <tr>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          Combined
                        </th>
                        <th class="text-right agla-mf-cell">M</th>
                        <th class="text-right agla-mf-cell">F</th>
                        <th class="text-right agla-combined-cell">
                          Combined
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(r, j) in rs.additional_gla_cover_band_rates"
                        :key="'agla-final-' + i + '-' + j"
                      >
                        <td class="agla-age-band">
                          {{ formatBandLabel(r) }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{ formatFinal(rs.category, j, 'binder', 'male') }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{ formatFinal(rs.category, j, 'binder', 'female') }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            formatFinal(rs.category, j, 'binder', 'combined')
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{ formatFinal(rs.category, j, 'outsource', 'male') }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            formatFinal(rs.category, j, 'outsource', 'female')
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            formatFinal(rs.category, j, 'outsource', 'combined')
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            formatFinal(rs.category, j, 'commission', 'male')
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{
                            formatFinal(rs.category, j, 'commission', 'female')
                          }}
                        </td>
                        <td
                          class="text-right agla-combined-cell agla-group-divider"
                        >
                          {{
                            formatFinal(
                              rs.category,
                              j,
                              'commission',
                              'combined'
                            )
                          }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{ formatFinal(rs.category, j, 'office', 'male') }}
                        </td>
                        <td class="text-right agla-mf-cell">
                          {{ formatFinal(rs.category, j, 'office', 'female') }}
                        </td>
                        <td class="text-right agla-combined-cell">
                          {{
                            formatFinal(rs.category, j, 'office', 'combined')
                          }}
                        </td>
                      </tr>
                    </tbody>
                  </v-table>
                </div>
              </template>
            </template>
            <v-alert v-else type="warning" density="compact" variant="tonal">
              {{ benefitTitle }} is enabled for
              <strong>{{ rs.category }}</strong> but per-band rates could not be
              computed. Re-run the quote calculation after fixing the
              configuration.
            </v-alert>
          </v-window-item>
        </v-window>
      </template>
    </template>
  </base-card>
  <v-snackbar
    :model-value="!!saveError"
    color="error"
    location="bottom"
    :timeout="6000"
    @update:model-value="(v) => (v ? null : (saveError = null))"
  >
    {{ saveError }}
    <template #actions>
      <v-btn variant="text" @click="saveError = null">Dismiss</v-btn>
    </template>
  </v-snackbar>
  <v-snackbar
    :model-value="!!saveSuccess"
    color="success"
    location="bottom"
    :timeout="3500"
    @update:model-value="(v) => (v ? null : (saveSuccess = null))"
  >
    <v-icon class="mr-2">mdi-check-circle</v-icon>
    {{ saveSuccess }}
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted } from 'vue'
import * as XLSX from 'xlsx'
import { AgCharts } from 'ag-charts-vue3'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ChartMenu from '@/renderer/components/charts/ChartMenu.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { roundUpToTwoDecimalsAccounting } from '@/renderer/utils/format_values'

type Gender = 'male' | 'female' | 'combined'

interface BandRate {
  min_age: number
  max_age: number
  office_rate_per1000: number
  office_rate_per1000_male: number
  office_rate_per1000_female: number
  risk_rate_per1000: number
  risk_rate_per1000_male: number
  risk_rate_per1000_female: number
  binder_fee_per1000: number
  binder_fee_per1000_male: number
  binder_fee_per1000_female: number
  outsource_fee_per1000: number
  outsource_fee_per1000_male: number
  outsource_fee_per1000_female: number
  commission_per1000: number
  commission_per1000_male: number
  commission_per1000_female: number
  weighted_office_rate_per1000?: number | null
  weighted_office_rate_per1000_male?: number | null
  weighted_office_rate_per1000_female?: number | null
  original_office_rate_per1000?: number | null
  original_office_rate_per1000_male?: number | null
  original_office_rate_per1000_female?: number | null
  smoothed_office_rate_per1000?: number | null
  smoothed_office_rate_per1000_male?: number | null
  smoothed_office_rate_per1000_female?: number | null
  smoothing_factor?: number | null
  smoothing_factor_male?: number | null
  smoothing_factor_female?: number | null
}

interface CategorySummary {
  category: string
  additional_gla_cover_benefit: boolean | number | string
  additional_gla_cover_age_band_type?: string
  additional_gla_cover_male_prop_used?: number | null
  additional_gla_cover_band_rates?: BandRate[]
}

const props = defineProps<{
  resultSummaries: any[]
  quote?: any
}>()

const emit = defineEmits<{
  (e: 'quote-updated'): void
}>()

const categories = computed<CategorySummary[]>(() =>
  (props.resultSummaries ?? []).filter((rs: any) => {
    const v = rs?.additional_gla_cover_benefit
    return v === true || v === 1 || v === '1' || v === 'true'
  })
)

// Scheme-wide blended commission rate, derived from the post-commission
// totals on the response so the commission shown reconciles to the scheme
// final commission and final premium the rest of the quote displays.
// final_scheme_total_commission is mirrored on every summary; final_total_
// annual_premium is per-category, so it is summed across summaries.
const schemeCommissionRate = computed<number>(() => {
  const sums = (props.resultSummaries ?? []) as any[]
  if (sums.length === 0) return 0
  const schemeCommission = Number(sums[0]?.final_scheme_total_commission ?? 0)
  const schemePremium = sums.reduce(
    (acc: number, s: any) => acc + Number(s?.final_total_annual_premium ?? 0),
    0
  )
  if (!schemePremium) return 0
  return schemeCommission / schemePremium
})

// Smoothed rates feed downstream pricing, so once the quote leaves the
// underwriter's hands (approved / accepted / on risk) the table drops to
// read-only. Status strings come from the Go models.Status enum.
const quoteStatus = computed<string>(() =>
  String((props.quote as any)?.status ?? '').toLowerCase()
)
const isLocked = computed<boolean>(() =>
  ['approved', 'accepted', 'in_force'].includes(quoteStatus.value)
)
const formatStatus = (s: string) => {
  if (!s) return ''
  return s.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())
}
const formatTimestamp = (iso?: string | null) => {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleString()
}

// Audit metadata is sourced from props.quote.scheme_categories (matched by
// category name) on first render, then overlaid with anything the save
// endpoint returned this session so the banner updates without waiting for
// the parent to refetch.
type AuditMetadata = { updatedAt?: string | null; updatedBy?: string | null }
const localAudit = reactive<Record<string, AuditMetadata>>({})
const getMetadata = (cat: string): AuditMetadata | null => {
  if (localAudit[cat]) return localAudit[cat]
  const sc = ((props.quote as any)?.scheme_categories ?? []).find(
    (c: any) => c?.scheme_category === cat
  )
  if (!sc) return null
  const at = sc.additional_gla_smoothed_updated_at
  const by = sc.additional_gla_smoothed_updated_by
  if (!at && !by) return null
  return { updatedAt: at ?? null, updatedBy: by ?? null }
}

const hasAnyRates = computed<boolean>(() =>
  categories.value.some(
    (rs) => (rs.additional_gla_cover_band_rates?.length ?? 0) > 0
  )
)

const activeTab = ref<string>('')
watch(
  categories,
  (cats) => {
    if (cats.length === 0) {
      activeTab.value = ''
      return
    }
    const stillValid = cats.some((rs) => rs.category === activeTab.value)
    if (!stillValid) activeTab.value = cats[0]?.category ?? ''
  },
  { immediate: true }
)

// Top-level view tab — Smoothing (default) vs Final results. Sits above the
// per-category tabs so switching views applies across every category.
const viewTab = ref<'smoothing' | 'final'>('smoothing')

// Customised benefit name from the GroupPricing benefit-map (resolved from
// the AGLA code on mount). Falls back to "Additional GLA Cover" when the
// API isn't reachable so the UI stays sensible.
const benefitTitle = ref('Additional GLA Cover')
onMounted(async () => {
  try {
    const res = await GroupPricingService.getBenefitMaps()
    const aglaBenefit: any = (res?.data ?? []).find(
      (item: any) => item?.benefit_code === 'AGLA'
    )
    if (aglaBenefit) {
      benefitTitle.value =
        aglaBenefit.benefit_alias ||
        aglaBenefit.benefit_name ||
        benefitTitle.value
    }
  } catch (e) {
    // Non-fatal — keep the default label
  }
})

const formatBandLabel = (r: { min_age: number; max_age: number }) => {
  if (r.max_age >= 150) return `${r.min_age}+`
  return `${r.min_age}–${r.max_age}`
}

const formatNullable = (v: number | null | undefined) =>
  v == null ? '—' : roundUpToTwoDecimalsAccounting(v)

// Per-category drafts: smoothed values and factors that the user has edited
// but not yet saved. Keyed by category, then by band index for smoothed
// overrides; factors are stored at the category level.
type SmoothedDraft = {
  smoothedByBand: Record<number, Partial<Record<Gender, number | null>>>
  factor: Partial<Record<Gender, number | null>>
}
const drafts = reactive<Record<string, SmoothedDraft>>({})
const saving = reactive<Record<string, boolean>>({})

const ensureDraft = (cat: string): SmoothedDraft => {
  if (!drafts[cat]) {
    drafts[cat] = { smoothedByBand: {}, factor: {} }
  }
  return drafts[cat]
}

const findCategory = (cat: string): CategorySummary | undefined =>
  categories.value.find((c) => c.category === cat)

// Prefers the pre-smoothing snapshot (original_office_rate_per1000_*) when
// present so the smoothing comparison and chart keep showing the unsmoothed
// reference even after office_rate_per1000_* has been overwritten with the
// smoothed-derived final value. Falls back to office_rate_per1000_* for old
// data written before the snapshot field existed.
const getOfficeRate = (band: BandRate, g: Gender): number => {
  if (g === 'male') {
    return (
      band.original_office_rate_per1000_male ?? band.office_rate_per1000_male
    )
  }
  if (g === 'female') {
    return (
      band.original_office_rate_per1000_female ??
      band.office_rate_per1000_female
    )
  }
  return band.original_office_rate_per1000 ?? band.office_rate_per1000
}

const getPersistedSmoothed = (
  band: BandRate,
  g: Gender
): number | null | undefined => {
  if (g === 'male') return band.smoothed_office_rate_per1000_male
  if (g === 'female') return band.smoothed_office_rate_per1000_female
  return band.smoothed_office_rate_per1000
}

const getCategoryFactor = (cat: string, g: Gender): number | null => {
  const d = drafts[cat]
  if (d && g in d.factor && d.factor[g] !== undefined) {
    return d.factor[g] ?? null
  }
  // Fall back to persisted factor on first band (factors are category-wide).
  const rs = findCategory(cat)
  const first = rs?.additional_gla_cover_band_rates?.[0]
  if (!first) return null
  if (g === 'male') return first.smoothing_factor_male ?? null
  if (g === 'female') return first.smoothing_factor_female ?? null
  return first.smoothing_factor ?? null
}

// Resolves the value to display in the smoothed input cell, in priority:
// 1. user's draft direct edit, 2. user's draft factor × office, 3. persisted
// smoothed, 4. office rate (default).
const effectiveSmoothed = (cat: string, bandIdx: number, g: Gender): string => {
  const v = effectiveSmoothedNumber(cat, bandIdx, g)
  if (v == null) return ''
  // Round-up to two decimals so an unsmoothed cell displays identically to
  // its OfficeRate/1000 source, which uses roundUpToTwoDecimalsAccounting.
  return (Math.ceil(v * 100) / 100).toFixed(2)
}

const effectiveSmoothedNumber = (
  cat: string,
  bandIdx: number,
  g: Gender
): number | null => {
  const rs = findCategory(cat)
  const band = rs?.additional_gla_cover_band_rates?.[bandIdx]
  if (!band) return null
  const d = drafts[cat]
  if (d?.smoothedByBand[bandIdx] && g in d.smoothedByBand[bandIdx]) {
    const v = d.smoothedByBand[bandIdx][g]
    return v ?? null
  }
  const factor = d?.factor?.[g]
  if (factor != null && Number.isFinite(factor)) {
    return getOfficeRate(band, g) * factor
  }
  const persisted = getPersistedSmoothed(band, g)
  if (persisted != null) return persisted
  return getOfficeRate(band, g)
}

const onSmoothedInput = (
  cat: string,
  bandIdx: number,
  g: Gender,
  raw: string
) => {
  const d = ensureDraft(cat)
  if (!d.smoothedByBand[bandIdx]) d.smoothedByBand[bandIdx] = {}
  if (raw === '' || raw == null) {
    d.smoothedByBand[bandIdx][g] = null
  } else {
    const n = Number(raw)
    if (Number.isFinite(n)) d.smoothedByBand[bandIdx][g] = n
  }
}

const onFactorInput = (cat: string, g: Gender, raw: string) => {
  const d = ensureDraft(cat)
  if (raw === '' || raw == null) {
    d.factor[g] = null
  } else {
    const n = Number(raw)
    if (Number.isFinite(n)) d.factor[g] = n
  }
  // Clearing per-band overrides on factor change so the factor takes effect.
  d.smoothedByBand = {}
}

const hasUnsavedChanges = (cat: string): boolean => {
  const d = drafts[cat]
  if (!d) return false
  if (Object.keys(d.smoothedByBand).length > 0) return true
  for (const k of ['male', 'female', 'combined'] as Gender[]) {
    if (k in d.factor) return true
  }
  return false
}

const resetCategory = (cat: string) => {
  delete drafts[cat]
}

const buildPayloadRows = (
  cat: string
): Parameters<
  typeof GroupPricingService.saveAdditionalGlaSmoothedRates
>[1]['rows'] => {
  const rs = findCategory(cat)
  const bands = rs?.additional_gla_cover_band_rates ?? []
  const d = drafts[cat]
  const rows: any[] = []
  bands.forEach((band, idx) => {
    const row: any = { min_age: band.min_age, max_age: band.max_age }
    let touched = false
    for (const g of ['male', 'female', 'combined'] as Gender[]) {
      const directEdit = d?.smoothedByBand?.[idx]?.[g]
      if (directEdit !== undefined) {
        const fieldKey =
          g === 'combined'
            ? 'smoothed_office_rate_per1000'
            : `smoothed_office_rate_per1000_${g}`
        if (directEdit === null) {
          row[g === 'combined' ? 'clear_smoothed' : `clear_smoothed_${g}`] =
            true
        } else {
          row[fieldKey] = directEdit
        }
        touched = true
      } else if (d?.factor && g in d.factor) {
        const factor = d.factor[g]
        const fieldKey =
          g === 'combined'
            ? 'smoothed_office_rate_per1000'
            : `smoothed_office_rate_per1000_${g}`
        if (factor == null) {
          row[g === 'combined' ? 'clear_smoothed' : `clear_smoothed_${g}`] =
            true
        } else {
          row[fieldKey] = getOfficeRate(band, g) * factor
        }
        touched = true
      }
    }
    if (touched) rows.push(row)
  })
  // Encode the category-wide factor on every row so it persists alongside.
  if (d?.factor) {
    for (const g of ['male', 'female', 'combined'] as Gender[]) {
      if (!(g in d.factor)) continue
      const factor = d.factor[g]
      const fieldKey =
        g === 'combined' ? 'smoothing_factor' : `smoothing_factor_${g}`
      const clearKey = g === 'combined' ? 'clear_factor' : `clear_factor_${g}`
      bands.forEach((band, idx) => {
        let row = rows.find(
          (r) => r.min_age === band.min_age && r.max_age === band.max_age
        )
        if (!row) {
          row = { min_age: band.min_age, max_age: band.max_age }
          rows.push(row)
        }
        if (factor == null) row[clearKey] = true
        else
          row[fieldKey] = factor
          // Mark the band index for stable ordering when sorting later
        ;(row as any).__idx = idx
      })
    }
  }
  return rows.map((r) => {
    const { __idx, ...rest } = r as any
    return rest
  })
}

const saveError = ref<string | null>(null)
const saveSuccess = ref<string | null>(null)

// Edit-mode toggle. Once a category has saved smoothed values, the inputs
// drop to read-only and the Save button reads "Edit". Clicking Edit flips
// editingMode[cat] to true so the user can amend; the next Save persists
// and flips back. New (never-saved) categories start in edit mode so the
// underwriter can enter values without an extra click.
const editingMode = reactive<Record<string, boolean>>({})

const hasSavedSmoothed = (cat: string): boolean => {
  const rs = findCategory(cat)
  const bands = rs?.additional_gla_cover_band_rates ?? []
  for (const b of bands) {
    if (
      b.smoothed_office_rate_per1000 != null ||
      b.smoothed_office_rate_per1000_male != null ||
      b.smoothed_office_rate_per1000_female != null ||
      b.smoothing_factor != null ||
      b.smoothing_factor_male != null ||
      b.smoothing_factor_female != null
    ) {
      return true
    }
  }
  return false
}

const isReadonly = (cat: string): boolean => {
  if (isLocked.value) return true
  if (!hasSavedSmoothed(cat)) return false
  return editingMode[cat] !== true
}

const enterEditMode = (cat: string) => {
  if (isLocked.value) return
  editingMode[cat] = true
}

const saveCategory = async (cat: string) => {
  const rows = buildPayloadRows(cat)
  if (rows.length === 0) return
  const quoteId = (props.quote as any)?.id
  if (!quoteId) return
  saving[cat] = true
  saveError.value = null
  try {
    const resp = await GroupPricingService.saveAdditionalGlaSmoothedRates(
      quoteId,
      { category: cat, rows }
    )
    const updated: BandRate[] | undefined =
      resp?.data?.additional_gla_cover_band_rates
    if (updated) {
      const rs = findCategory(cat)
      if (rs) rs.additional_gla_cover_band_rates = updated
    }
    localAudit[cat] = {
      updatedAt: resp?.data?.additional_gla_smoothed_updated_at ?? null,
      updatedBy: resp?.data?.additional_gla_smoothed_updated_by ?? null
    }
    delete drafts[cat]
    editingMode[cat] = false
    saveSuccess.value = `Smoothed rates for "${cat}" saved successfully.`
    emit('quote-updated')
  } catch (e: any) {
    if (e?.response?.status === 409) {
      saveError.value =
        e?.response?.data?.error ??
        'Smoothed rates are locked because this quote has been approved or accepted.'
    } else {
      saveError.value =
        e?.response?.data?.error ??
        e?.message ??
        'Failed to save smoothed rates.'
    }
    console.error('Failed to save smoothed AGLA rates', e)
  } finally {
    saving[cat] = false
  }
}

// ===== Template download =====
const downloadTemplate = (rs: CategorySummary) => {
  const bands = rs.additional_gla_cover_band_rates ?? []
  if (bands.length === 0) return
  const wb = XLSX.utils.book_new()
  const aoa: any[][] = []
  aoa.push([`${benefitTitle.value} Smoothed Rates — ${rs.category}`])
  aoa.push([
    'Edit the Smoothed columns and re-upload. MinAge / MaxAge identify each band.'
  ])
  aoa.push([])
  aoa.push([
    'Age Band',
    'MinAge',
    'MaxAge',
    'OfficeRate M',
    'OfficeRate F',
    'OfficeRate Combined',
    'Smoothed M',
    'Smoothed F',
    'Smoothed Combined'
  ])
  bands.forEach((b) => {
    aoa.push([
      formatBandLabel(b),
      b.min_age,
      b.max_age,
      Number(b.office_rate_per1000_male ?? 0),
      Number(b.office_rate_per1000_female ?? 0),
      Number(b.office_rate_per1000 ?? 0),
      Number(
        b.smoothed_office_rate_per1000_male ?? b.office_rate_per1000_male ?? 0
      ),
      Number(
        b.smoothed_office_rate_per1000_female ??
          b.office_rate_per1000_female ??
          0
      ),
      Number(b.smoothed_office_rate_per1000 ?? b.office_rate_per1000 ?? 0)
    ])
  })
  const ws = XLSX.utils.aoa_to_sheet(aoa)
  ws['!cols'] = [
    { wch: 14 },
    { wch: 8 },
    { wch: 8 },
    { wch: 14 },
    { wch: 14 },
    { wch: 18 },
    { wch: 14 },
    { wch: 14 },
    { wch: 18 }
  ]
  const sheetName = String(rs.category ?? 'Category')
    .replace(/[:\\/?*[\]]/g, '_')
    .slice(0, 31)
  XLSX.utils.book_append_sheet(wb, ws, sheetName || 'Category')
  const schemeName = (props.quote as any)?.scheme_name ?? 'quote'
  XLSX.writeFile(
    wb,
    `additional_gla_smoothed_template_${schemeName}_${sheetName}.xlsx`
  )
}

// ===== Template upload =====
const uploadInputs: Record<string, HTMLInputElement | null> = {}
const bindUploadInput = (cat: string, el: HTMLInputElement | null) => {
  uploadInputs[cat] = el
}
const triggerUpload = (cat: string) => {
  if (isLocked.value) return
  const el = uploadInputs[cat]
  if (el) el.click()
}
const onUploadFile = async (e: Event, cat: string) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  try {
    const buffer = await file.arrayBuffer()
    const wb = XLSX.read(buffer, { type: 'array' })
    const ws = wb.Sheets[wb.SheetNames[0]]
    const rows: any[] = XLSX.utils.sheet_to_json(ws, { defval: null })
    const rs = findCategory(cat)
    const bands = rs?.additional_gla_cover_band_rates ?? []
    const d = ensureDraft(cat)
    // Match by MinAge/MaxAge.
    rows.forEach((row) => {
      const minAge = Number(row.MinAge ?? row.min_age)
      const maxAge = Number(row.MaxAge ?? row.max_age)
      if (!Number.isFinite(minAge) || !Number.isFinite(maxAge)) return
      const idx = bands.findIndex(
        (b) => b.min_age === minAge && b.max_age === maxAge
      )
      if (idx < 0) return
      if (!d.smoothedByBand[idx]) d.smoothedByBand[idx] = {}
      const m = row['Smoothed M']
      const f = row['Smoothed F']
      const c = row['Smoothed Combined']
      if (m != null && Number.isFinite(Number(m))) {
        d.smoothedByBand[idx].male = Number(m)
      }
      if (f != null && Number.isFinite(Number(f))) {
        d.smoothedByBand[idx].female = Number(f)
      }
      if (c != null && Number.isFinite(Number(c))) {
        d.smoothedByBand[idx].combined = Number(c)
      }
    })
  } catch (err) {
    console.error('Failed to parse uploaded smoothed rates xlsx', err)
  } finally {
    if (input) input.value = ''
  }
}

// ===== Final results: smoothed office rate × scheme-level fee rates =====
//
// Binder and outsource per-1,000 values on the band are divided by the
// unsmoothed office rate to recover the per-category fee rate, which we
// then apply to the smoothed office rate. Commission, by contrast, uses
// the scheme-wide ratio (final_scheme_total_commission ÷
// Σ final_total_annual_premium) so it reconciles to the scheme totals
// shown elsewhere in the quote.
type FeeKind = 'binder' | 'outsource' | 'commission' | 'office'

const rateRatio = (band: BandRate, kind: FeeKind): number => {
  if (kind === 'office') return 1
  if (kind === 'commission') return schemeCommissionRate.value
  const office = band.office_rate_per1000
  if (!office || office === 0) return 0
  if (kind === 'binder') return (band.binder_fee_per1000 ?? 0) / office
  return (band.outsource_fee_per1000 ?? 0) / office
}

const finalValue = (
  cat: string,
  bandIdx: number,
  kind: FeeKind,
  g: Gender
): number | null => {
  const rs = findCategory(cat)
  const band = rs?.additional_gla_cover_band_rates?.[bandIdx]
  if (!band) return null
  const smoothed = effectiveSmoothedNumber(cat, bandIdx, g)
  if (smoothed == null) return null
  return smoothed * rateRatio(band, kind)
}

const formatFinal = (
  cat: string,
  bandIdx: number,
  kind: FeeKind,
  g: Gender
): string => {
  const v = finalValue(cat, bandIdx, kind, g)
  if (v == null) return '—'
  return roundUpToTwoDecimalsAccounting(v)
}

// ===== Chart =====
const chartRefs: Record<string, any> = {}
const bindChartRef = (cat: string, el: any) => {
  chartRefs[cat] = el
}

const chartDataFor = (rs: CategorySummary) => {
  const bands = rs.additional_gla_cover_band_rates ?? []
  return bands.map((b, idx) => ({
    band: formatBandLabel(b),
    office_m: getOfficeRate(b, 'male'),
    office_f: getOfficeRate(b, 'female'),
    office_c: getOfficeRate(b, 'combined'),
    smoothed_m: effectiveSmoothedNumber(rs.category, idx, 'male'),
    smoothed_f: effectiveSmoothedNumber(rs.category, idx, 'female'),
    smoothed_c: effectiveSmoothedNumber(rs.category, idx, 'combined')
  }))
}

const chartOptionsFor = (rs: CategorySummary): any => {
  const data = chartDataFor(rs)
  return {
    data,
    background: { fill: 'transparent' },
    height: 320,
    series: [
      {
        type: 'line',
        xKey: 'band',
        yKey: 'office_m',
        yName: 'Office M',
        stroke: '#1976D2',
        strokeWidth: 2,
        marker: { enabled: true, size: 4 }
      },
      {
        type: 'line',
        xKey: 'band',
        yKey: 'office_f',
        yName: 'Office F',
        stroke: '#D81B60',
        strokeWidth: 2,
        marker: { enabled: true, size: 4 }
      },
      {
        type: 'line',
        xKey: 'band',
        yKey: 'office_c',
        yName: 'Office Combined',
        stroke: '#388E3C',
        strokeWidth: 2,
        marker: { enabled: true, size: 4 }
      },
      {
        type: 'line',
        xKey: 'band',
        yKey: 'smoothed_m',
        yName: 'Smoothed M',
        stroke: '#1976D2',
        strokeWidth: 2,
        lineDash: [6, 4],
        marker: { enabled: true, size: 4 }
      },
      {
        type: 'line',
        xKey: 'band',
        yKey: 'smoothed_f',
        yName: 'Smoothed F',
        stroke: '#D81B60',
        strokeWidth: 2,
        lineDash: [6, 4],
        marker: { enabled: true, size: 4 }
      },
      {
        type: 'line',
        xKey: 'band',
        yKey: 'smoothed_c',
        yName: 'Smoothed Combined',
        stroke: '#388E3C',
        strokeWidth: 2,
        lineDash: [6, 4],
        marker: { enabled: true, size: 4 }
      }
    ],
    axes: [
      { type: 'category', position: 'bottom', title: { text: 'Age Band' } },
      {
        type: 'number',
        position: 'left',
        title: { text: 'Rate per 1,000' }
      }
    ],
    legend: { position: 'bottom' }
  }
}

// ===== Excel export of the full summary (mirrors the existing pattern) =====
const exportSummaryToExcel = () => {
  const cats = categories.value.filter(
    (rs) => (rs.additional_gla_cover_band_rates?.length ?? 0) > 0
  )
  if (cats.length === 0) return
  const isFinal = viewTab.value === 'final'
  const wb = XLSX.utils.book_new()
  cats.forEach((rs) => {
    const bands = rs.additional_gla_cover_band_rates ?? []
    const aoa: any[][] = []
    aoa.push([
      `${benefitTitle.value} — ${rs.category} — ${
        isFinal ? 'Final results' : 'Smoothing'
      }`
    ])
    aoa.push([])
    if (isFinal) {
      aoa.push([
        'Age Band',
        'Binder M',
        'Binder F',
        'Binder Combined',
        'Outsource M',
        'Outsource F',
        'Outsource Combined',
        'Commission M',
        'Commission F',
        'Commission Combined',
        'Final Office M',
        'Final Office F',
        'Final Office Combined'
      ])
      bands.forEach((b, idx) => {
        aoa.push([
          formatBandLabel(b),
          finalValue(rs.category, idx, 'binder', 'male') ?? '',
          finalValue(rs.category, idx, 'binder', 'female') ?? '',
          finalValue(rs.category, idx, 'binder', 'combined') ?? '',
          finalValue(rs.category, idx, 'outsource', 'male') ?? '',
          finalValue(rs.category, idx, 'outsource', 'female') ?? '',
          finalValue(rs.category, idx, 'outsource', 'combined') ?? '',
          finalValue(rs.category, idx, 'commission', 'male') ?? '',
          finalValue(rs.category, idx, 'commission', 'female') ?? '',
          finalValue(rs.category, idx, 'commission', 'combined') ?? '',
          finalValue(rs.category, idx, 'office', 'male') ?? '',
          finalValue(rs.category, idx, 'office', 'female') ?? '',
          finalValue(rs.category, idx, 'office', 'combined') ?? ''
        ])
      })
    } else {
      aoa.push([
        'Age Band',
        'OfficeRate M',
        'OfficeRate F',
        'OfficeRate Combined',
        'Weighted M',
        'Weighted F',
        'Weighted Combined',
        'Smoothed M',
        'Smoothed F',
        'Smoothed Combined'
      ])
      bands.forEach((b, idx) => {
        aoa.push([
          formatBandLabel(b),
          Number(getOfficeRate(b, 'male')),
          Number(getOfficeRate(b, 'female')),
          Number(getOfficeRate(b, 'combined')),
          b.weighted_office_rate_per1000_male ?? '',
          b.weighted_office_rate_per1000_female ?? '',
          b.weighted_office_rate_per1000 ?? '',
          effectiveSmoothedNumber(rs.category, idx, 'male') ?? '',
          effectiveSmoothedNumber(rs.category, idx, 'female') ?? '',
          effectiveSmoothedNumber(rs.category, idx, 'combined') ?? ''
        ])
      })
    }
    const ws = XLSX.utils.aoa_to_sheet(aoa)
    const sheetName = String(rs.category ?? 'Category')
      .replace(/[:\\/?*[\]]/g, '_')
      .slice(0, 31)
    XLSX.utils.book_append_sheet(wb, ws, sheetName || 'Category')
  })
  const schemeName = (props.quote as any)?.scheme_name ?? 'quote'
  const fileSuffix = isFinal ? 'final_results' : 'smoothing'
  XLSX.writeFile(
    wb,
    `additional_gla_cover_${fileSuffix}_${schemeName}.xlsx`
  )
}
</script>

<style scoped>
/*
 * Mirrors the styling that lived inside QuoteBenefitSummary.vue for the
 * AGLA summary table — sticky Age Band column, brand-tinted Combined
 * column, lower-weight M/F sub-columns, and strong group dividers.
 *
 * v-table renders headers/cells outside the scoped boundary on Vuetify 3,
 * so :deep() is needed to pierce into the rendered <th>/<td>.
 */

/* Header-level Smoothing / Final tabs sit on the dark card header, so the
 * default Vuetify tab colours are too low-contrast. Swap to white text with
 * a subtle background pill for the active tab. */
:deep(.agla-header-tabs .v-tab) {
  color: rgba(255, 255, 255, 0.78);
  text-transform: none;
  font-weight: 500;
  letter-spacing: 0;
  min-width: auto;
  padding: 0 14px;
}
:deep(.agla-header-tabs .v-tab--selected) {
  color: #ffffff;
  background-color: rgba(255, 255, 255, 0.16);
  border-radius: 4px;
}
:deep(.agla-header-tabs .v-tab:hover) {
  color: #ffffff;
}

/* Sticky Age Band column so horizontal scroll keeps the row label in view. */
:deep(.agla-summary-table th.agla-age-band),
:deep(.agla-summary-table td.agla-age-band) {
  position: sticky;
  left: 0;
  background-color: rgb(var(--v-theme-surface));
  border-right: 2px solid rgba(var(--v-border-color), var(--v-border-opacity));
  z-index: 1;
}
:deep(.agla-summary-table thead th.agla-age-band) {
  z-index: 3;
}

/* M/F sub-columns sit at slightly lower visual weight than Combined. */
:deep(.agla-summary-table .agla-mf-cell) {
  background-color: rgba(0, 63, 88, 0.04);
  color: rgba(var(--v-theme-on-surface), 0.78);
}

/* Combined column is the headline value — bold and brand-tinted. */
:deep(.agla-summary-table .agla-combined-cell) {
  font-weight: 600;
  background-color: rgba(0, 63, 88, 0.08);
}
:deep(.agla-summary-table thead th.agla-combined-cell) {
  color: #003f58;
}

/* Strong divider on the right edge of every group except the last. */
:deep(.agla-summary-table .agla-group-divider) {
  border-right: 2px solid rgba(0, 63, 88, 0.35);
}

/* Factor row sits as a subdued helper strip below the editable smoothed row. */
:deep(.agla-summary-table tr.factor-row td) {
  background-color: rgba(0, 63, 88, 0.02);
  font-size: 0.8rem;
  color: rgba(var(--v-theme-on-surface), 0.6);
}

/* Editable inputs blend into the table and pick up brand colour on focus. */
.smoothed-input,
.factor-input {
  width: 100%;
  padding: 4px 6px;
  border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  border-radius: 4px;
  text-align: right;
  background: rgb(var(--v-theme-surface));
  color: rgb(var(--v-theme-on-surface));
  font-size: 0.85rem;
  outline: none;
}
.smoothed-input:focus,
.factor-input:focus {
  border-color: #003f58;
  box-shadow: 0 0 0 2px rgba(0, 63, 88, 0.15);
}
</style>
