<!-- eslint-disable no-use-before-define -->
<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Group Pricing Dashboard</span>
          </template>
          <template #default>
            <v-row class="align-center">
              <v-col cols="4">
                <v-select
                  v-model="selectedYear"
                  variant="outlined"
                  density="compact"
                  :items="availableYears"
                  label="Select Year"
                  placeholder="Select Year"
                  @update:model-value="refreshDashboard"
                ></v-select>
              </v-col>
              <v-col v-if="financialYearInfo" cols="auto">
                <v-chip
                  size="small"
                  color="primary"
                  variant="tonal"
                  prepend-icon="mdi-calendar-range"
                >
                  {{ financialYearInfo.financial_year_label }}
                </v-chip>
              </v-col>
            </v-row>

            <v-divider class="mb-5 mt-n6"></v-divider>

            <!-- ── Section 1: Summary KPIs ───────────────────────────────── -->
            <v-row class="d-flex justify-center">
              <v-col v-for="card in cards" :key="card.title" :cols="card.flex">
                <v-card
                  variant="tonal"
                  color="primary"
                  class="dash-card"
                  :style="card.route ? 'cursor:pointer' : ''"
                  @click="
                    card.route ? router.push({ name: card.route }) : undefined
                  "
                >
                  <v-card-subtitle
                    ><h5>{{ card.title }}</h5></v-card-subtitle
                  >
                  <v-card-text>
                    <v-row>
                      <v-col class="d-flex justify-center">
                        <h2>{{ addFormat(card) }}</h2>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-divider class="my-7"></v-divider>

            <!-- ── Quotes Data banner ─────────────────────────────────────── -->
            <v-alert
              v-model="showQuotesBanner"
              type="info"
              variant="tonal"
              color="blue"
              closable
              density="compact"
              class="mb-4"
              icon="mdi-file-document-outline"
            >
              <strong>Quotes Data</strong> — Sections below show activity across
              all quotes (all statuses) for the selected year.
            </v-alert>

            <!-- ── Section 3: Conversion & Quotes Overview ───────────────── -->
            <v-row class="align-center mb-2">
              <v-col>
                <span class="text-subtitle-1 font-weight-bold"
                  >Conversion &amp; Quotes Overview</span
                >
              </v-col>
              <v-col cols="auto">
                <v-select
                  v-model="conversionViewBy"
                  :items="['Count', 'Annual Premium']"
                  density="compact"
                  variant="outlined"
                  hide-details
                  style="max-width: 180px"
                  label="View by"
                  @update:model-value="changeConversionDataSource"
                />
              </v-col>
            </v-row>
            <v-row class="pane-bg pa-5 ma-3">
              <v-col cols="3">
                <ChartMenu
                  :chart-ref="newQuoteCharts"
                  title="New Quotes"
                  :data="newQuoteOptions?.data ?? []"
                />
                <ag-charts
                  v-if="newQuoteOptions"
                  ref="newQuoteCharts"
                  :options="newQuoteOptions"
                ></ag-charts>
              </v-col>
              <v-col cols="3">
                <ChartMenu
                  :chart-ref="renewalsCharts"
                  title="Renewals"
                  :data="renewalsOptions?.data ?? []"
                />
                <ag-charts
                  v-if="renewalsOptions"
                  ref="renewalsCharts"
                  :options="renewalsOptions"
                ></ag-charts>
              </v-col>
              <v-col cols="3">
                <ChartMenu
                  :chart-ref="conversionCharts"
                  title="Quote Conversion"
                  :data="conversionOptions?.data ?? []"
                />
                <ag-charts
                  v-if="conversionOptions"
                  ref="conversionCharts"
                  :options="conversionOptions"
                ></ag-charts>
              </v-col>
              <v-col cols="3">
                <ChartMenu
                  :chart-ref="inForceCharts"
                  title="Accepted Quotes"
                  :data="inForceSchemesOptions?.data ?? []"
                />
                <ag-charts
                  v-if="inForceSchemesOptions"
                  ref="inForceCharts"
                  :options="inForceSchemesOptions"
                ></ag-charts>
              </v-col>
            </v-row>

            <!-- Broker Performance Table -->
            <v-row class="align-center mt-5 mb-2">
              <v-col>
                <span class="text-subtitle-1 font-weight-bold"
                  >Broker Performance (Top 10)</span
                >
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-table density="compact" class="card-bg broker-table">
                  <thead>
                    <tr>
                      <th class="text-left">Broker</th>
                      <th class="text-center">Total Quotes</th>
                      <th class="text-center">Accepted</th>
                      <th class="text-center">Conversion Rate</th>
                      <th class="text-right">Total Premium</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="b in brokerMetrics" :key="b.broker_id">
                      <td>{{ b.broker_name || 'Direct' }}</td>
                      <td class="text-center">{{ b.total_quotes }}</td>
                      <td class="text-center">{{ b.accepted_quotes }}</td>
                      <td class="text-center">
                        <v-chip
                          size="small"
                          :color="
                            b.conversion_rate >= 50
                              ? 'success'
                              : b.conversion_rate >= 25
                                ? 'warning'
                                : 'error'
                          "
                          variant="tonal"
                          >{{ b.conversion_rate }}%</v-chip
                        >
                      </td>
                      <td class="text-right"
                        >R{{ formatNumber(b.total_premium) }}</td
                      >
                    </tr>
                    <tr v-if="!brokerMetrics || brokerMetrics.length === 0">
                      <td colspan="5" class="text-center text-grey py-4">
                        No broker data for selected year
                      </td>
                    </tr>
                  </tbody>
                </v-table>
              </v-col>
            </v-row>

            <v-divider class="my-7"></v-divider>

            <!-- ── Section 2: Sales Pipeline ─────────────────────────────── -->
            <v-row class="align-center mb-2">
              <v-col>
                <span class="text-subtitle-1 font-weight-bold"
                  >Sales Pipeline</span
                >
              </v-col>
              <v-col cols="3">
                <v-select
                  v-model="selectedDataView"
                  placeholder="Select a data view"
                  variant="outlined"
                  density="compact"
                  :items="['Annual Premium', 'Count']"
                  @update:model-value="changeChartDataSource"
                ></v-select>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="6" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="trendCharts"
                  title="Monthly Quote Volume"
                  :data="monthlyTrendOptions?.data ?? []"
                />
                <ag-charts
                  v-if="monthlyTrendOptions"
                  ref="trendCharts"
                  :options="monthlyTrendOptions"
                ></ag-charts>
              </v-col>
              <v-col class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="funnelCharts"
                  title="Pipeline by Stage"
                  :data="pipelineFunnelOptions?.data ?? []"
                />
                <ag-charts
                  v-if="pipelineFunnelOptions"
                  ref="funnelCharts"
                  :options="pipelineFunnelOptions"
                ></ag-charts>
              </v-col>
            </v-row>
            <v-row class="mt-4">
              <v-col cols="6" class="card-bg mx-1">
                <!-- header row: title + train button -->
                <div
                  class="d-flex align-center justify-space-between px-2 pt-2"
                >
                  <span class="text-body-2 font-weight-bold"
                    >Win Probability Distribution</span
                  >
                  <v-btn
                    size="x-small"
                    color="primary"
                    variant="tonal"
                    prepend-icon="mdi-brain"
                    :loading="trainLoading"
                    @click="openTrainDialog"
                    >Train Model</v-btn
                  >
                </div>

                <!-- model metadata -->
                <div
                  v-if="modelInfo"
                  class="px-2 pb-1 text-caption text-grey-darken-1 d-flex flex-wrap align-center gap-3"
                >
                  <span
                    >Last trained:
                    {{
                      modelInfo.trained_at
                        ? new Date(modelInfo.trained_at).toLocaleDateString(
                            'en-ZA'
                          )
                        : '—'
                    }}</span
                  >

                  <!-- Accuracy with tooltip -->
                  <v-tooltip location="bottom" max-width="300">
                    <template #activator="{ props: ttp }">
                      <span
                        v-bind="ttp"
                        class="model-meta-chip"
                        :class="accuracyClass"
                      >
                        Accuracy:
                        {{
                          modelInfo.accuracy != null
                            ? modelInfo.accuracy + '%'
                            : '—'
                        }}
                        <v-icon size="10" class="ml-1"
                          >mdi-information-outline</v-icon
                        >
                      </span>
                    </template>
                    <span>
                      <strong>Model Accuracy</strong><br />
                      Percentage of quotes correctly classified as won or lost
                      on the held-out test set.<br /><br />
                      <strong>Interpretation:</strong><br />
                      &lt; 60% — Low confidence, insufficient data<br />
                      60–75% — Moderate, usable as a guide<br />
                      75–85% — Good, reliable for prioritisation<br />
                      &gt; 85% — High confidence
                    </span>
                  </v-tooltip>

                  <!-- AUC with click-through dialog -->
                  <span
                    class="model-meta-chip model-meta-clickable"
                    @click="infoDialog = 'auc'"
                  >
                    AUC: {{ modelInfo.auc ?? '—' }}
                    <v-icon size="10" class="ml-1"
                      >mdi-information-outline</v-icon
                    >
                  </span>

                  <!-- Training size with click-through dialog -->
                  <span
                    class="model-meta-chip model-meta-clickable"
                    @click="infoDialog = 'training_size'"
                  >
                    Training size: {{ modelInfo.training_size ?? '—' }}
                    <v-icon size="10" class="ml-1"
                      >mdi-information-outline</v-icon
                    >
                  </span>
                </div>

                <!-- AUC info dialog -->
                <v-dialog v-model="aucDialogOpen" max-width="520" scrollable>
                  <v-card>
                    <v-card-title class="d-flex align-center gap-2">
                      <v-icon color="primary">mdi-chart-bell-curve</v-icon>
                      AUC — Area Under the ROC Curve
                    </v-card-title>
                    <v-divider></v-divider>
                    <v-card-text>
                      <p class="mb-3">
                        AUC measures how well the model
                        <strong>separates won quotes from lost quotes</strong>
                        across all possible score thresholds — independent of
                        any fixed cut-off.
                      </p>
                      <p class="mb-1 font-weight-medium">How to read it:</p>
                      <v-table density="compact" class="mb-3">
                        <thead>
                          <tr
                            ><th>AUC range</th><th>Meaning</th
                            ><th>Confidence</th></tr
                          >
                        </thead>
                        <tbody>
                          <tr
                            ><td>0.50</td><td>No better than random guessing</td
                            ><td class="text-error">None</td></tr
                          >
                          <tr
                            ><td>0.50–0.65</td><td>Weak signal, limited data</td
                            ><td class="text-warning">Low</td></tr
                          >
                          <tr
                            ><td>0.65–0.75</td
                            ><td>Moderate — useful directional guide</td
                            ><td class="text-info">Moderate</td></tr
                          >
                          <tr
                            ><td>0.75–0.85</td><td>Good discrimination</td
                            ><td class="text-success">Good</td></tr
                          >
                          <tr
                            ><td>0.85–1.00</td><td>Excellent separation</td
                            ><td class="text-success">High</td></tr
                          >
                        </tbody>
                      </v-table>
                      <v-alert
                        type="info"
                        variant="tonal"
                        density="compact"
                        class="mb-3"
                      >
                        <strong
                          >Current AUC: {{ modelInfo?.auc ?? '—' }}</strong
                        >
                        <span v-if="modelInfo?.auc">
                          &nbsp;— {{ aucInterpretation }}
                        </span>
                      </v-alert>
                      <p class="text-caption text-grey-darken-1">
                        AUC is more reliable than accuracy when wins and losses
                        are imbalanced (e.g. 80% of quotes are lost). A high AUC
                        means the model correctly ranks a randomly chosen won
                        quote above a randomly chosen lost quote.
                      </p>
                    </v-card-text>
                    <v-divider></v-divider>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn
                        variant="text"
                        size="small"
                        @click="infoDialog = null"
                        >Close</v-btn
                      >
                    </v-card-actions>
                  </v-card>
                </v-dialog>

                <!-- Training size info dialog -->
                <v-dialog
                  v-model="trainingSizeDialogOpen"
                  max-width="520"
                  scrollable
                >
                  <v-card>
                    <v-card-title class="d-flex align-center gap-2">
                      <v-icon color="primary">mdi-database-outline</v-icon>
                      Training Size — Data Credibility
                    </v-card-title>
                    <v-divider></v-divider>
                    <v-card-text>
                      <p class="mb-3">
                        Training size is the number of historical quotes (with
                        known outcomes) used to fit the model. More data means
                        the model's weights are
                        <strong>more statistically stable</strong> and
                        generalisable.
                      </p>
                      <p class="mb-1 font-weight-medium"
                        >Credibility thresholds:</p
                      >
                      <v-table density="compact" class="mb-3">
                        <thead>
                          <tr
                            ><th>Sample size</th><th>Credibility</th
                            ><th>Recommendation</th></tr
                          >
                        </thead>
                        <tbody>
                          <tr
                            ><td>&lt; 20</td
                            ><td class="text-error">Insufficient</td
                            ><td>Training blocked</td></tr
                          >
                          <tr
                            ><td>20–49</td><td class="text-warning">Low</td
                            ><td
                              >Scores are indicative only — treat with
                              caution</td
                            ></tr
                          >
                          <tr
                            ><td>50–99</td><td class="text-info">Moderate</td
                            ><td>Useful for pipeline prioritisation</td></tr
                          >
                          <tr
                            ><td>100–499</td><td class="text-success">Good</td
                            ><td
                              >Reliable for underwriting and sales decisions</td
                            ></tr
                          >
                          <tr
                            ><td>500+</td><td class="text-success">High</td
                            ><td
                              >Production-grade predictive credibility</td
                            ></tr
                          >
                        </tbody>
                      </v-table>
                      <v-alert
                        :type="trainingSizeAlertType"
                        variant="tonal"
                        density="compact"
                        class="mb-3"
                      >
                        <strong
                          >Current training size:
                          {{ modelInfo?.training_size ?? '—' }}</strong
                        >
                        <span v-if="modelInfo?.training_size">
                          &nbsp;— {{ trainingSizeInterpretation }}
                        </span>
                      </v-alert>
                      <p class="text-caption text-grey-darken-1">
                        The model uses an 80/20 chronological train/test split,
                        so the displayed size reflects the training portion
                        only. As more quotes reach terminal status (accepted,
                        rejected, expired etc.), retrain the model to improve
                        credibility.
                      </p>
                    </v-card-text>
                    <v-divider></v-divider>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn
                        variant="text"
                        size="small"
                        @click="infoDialog = null"
                        >Close</v-btn
                      >
                    </v-card-actions>
                  </v-card>
                </v-dialog>

                <ChartMenu
                  :chart-ref="winProbCharts"
                  title="Win Probability Distribution"
                  :data="winProbDistOptions?.data ?? []"
                />
                <ag-charts
                  v-if="winProbDistOptions"
                  ref="winProbCharts"
                  :options="winProbDistOptions"
                ></ag-charts>
                <div v-else class="text-center text-grey py-6 text-caption">
                  Train the model to see win probability scores for active
                  quotes
                </div>
              </v-col>

              <!-- Monthly Conversion Rate Trend -->
              <v-col class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="convRateCharts"
                  title="Monthly Conversion Rate Trend"
                  :data="convRateTrendOptions?.data ?? []"
                />
                <ag-charts
                  v-if="convRateTrendOptions"
                  ref="convRateCharts"
                  :options="convRateTrendOptions"
                ></ag-charts>
                <div v-else class="text-center text-grey py-6 text-caption">
                  No conversion data for selected year
                </div>
              </v-col>
            </v-row>

            <!-- Industry Mix & Scheme Size row -->
            <v-row class="mt-4">
              <v-col cols="6" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="industryPipelineCharts"
                  title="Industry Mix of Active Pipeline"
                  :data="industryPipelineOptions?.data ?? []"
                />
                <ag-charts
                  v-if="industryPipelineOptions"
                  ref="industryPipelineCharts"
                  :options="industryPipelineOptions"
                ></ag-charts>
                <div v-else class="text-center text-grey py-6 text-caption">
                  No active pipeline quotes for selected year
                </div>
              </v-col>
              <v-col class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="schemeSizeCharts"
                  title="Scheme Size Distribution (Active Pipeline)"
                  :data="schemeSizeOptions?.data ?? []"
                />
                <ag-charts
                  v-if="schemeSizeOptions"
                  ref="schemeSizeCharts"
                  :options="schemeSizeOptions"
                ></ag-charts>
                <div v-else class="text-center text-grey py-6 text-caption">
                  No active pipeline quotes for selected year
                </div>
              </v-col>
            </v-row>

            <!-- Train Model confirmation dialog -->
            <v-dialog v-model="trainDialog" max-width="460" persistent>
              <v-card>
                <v-card-title class="d-flex align-center gap-2">
                  <v-icon color="warning">mdi-alert-circle-outline</v-icon>
                  Train Win Probability Model
                </v-card-title>
                <v-divider></v-divider>
                <v-card-text class="pt-4">
                  <p
                    >This will retrain the logistic regression model on all
                    historical quote outcomes.</p
                  >
                  <v-alert
                    type="warning"
                    variant="tonal"
                    class="mt-3"
                    density="compact"
                  >
                    At least <strong>20 terminal quotes</strong> (accepted,
                    in-force, rejected, not-taken-up, or expired) are required.
                    Training will be skipped silently if this threshold is not
                    met.
                  </v-alert>
                  <p class="mt-3 text-caption text-grey-darken-1">
                    Training runs in the background and may take a few seconds.
                    All active quotes will be automatically rescored once
                    training completes.
                  </p>
                </v-card-text>
                <v-divider></v-divider>
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn
                    variant="text"
                    size="small"
                    @click="trainDialog = false"
                    >Cancel</v-btn
                  >
                  <v-btn
                    color="primary"
                    variant="tonal"
                    size="small"
                    prepend-icon="mdi-brain"
                    :loading="trainLoading"
                    @click="triggerTraining"
                    >Start Training</v-btn
                  >
                </v-card-actions>
              </v-card>
            </v-dialog>

            <!-- Snackbar feedback -->
            <v-snackbar
              v-model="snackbar.show"
              :color="snackbar.color"
              timeout="4000"
              location="bottom end"
            >
              {{ snackbar.message }}
              <template #actions>
                <v-btn
                  variant="text"
                  size="small"
                  @click="snackbar.show = false"
                  >Dismiss</v-btn
                >
              </template>
            </v-snackbar>

            <v-divider class="my-7"></v-divider>

            <!-- ── In-Force Portfolio banner ──────────────────────────────── -->
            <v-alert
              v-model="showInforceBanner"
              type="success"
              variant="tonal"
              color="green"
              closable
              density="compact"
              class="mb-4"
              icon="mdi-shield-check-outline"
            >
              <strong>In-Force Portfolio</strong> — Sections below reflect the
              accepted and in-force book. Use the Data Source toggle on exposure
              charts to switch between In-Force, Quotes-only, or All data.
            </v-alert>

            <!-- ── Section 4: Portfolio Analysis ─────────────────────────── -->
            <v-row class="align-center mb-2">
              <v-col>
                <span class="text-subtitle-1 font-weight-bold"
                  >Portfolio Analysis</span
                >
              </v-col>
              <v-col cols="auto" class="d-flex align-center gap-2">
                <v-tooltip :text="exposureLockTooltip" location="bottom">
                  <template #activator="{ props: ttp }">
                    <span v-bind="ttp">
                      <v-btn
                        variant="tonal"
                        :color="isSelectedFYOpen ? 'primary' : 'grey'"
                        size="small"
                        :prepend-icon="
                          isSelectedFYOpen ? 'mdi-refresh' : 'mdi-lock-outline'
                        "
                        :loading="rebuildLoading"
                        :disabled="!isSelectedFYOpen"
                        @click="triggerExposureRebuild"
                      >
                        {{
                          isSelectedFYOpen ? 'Refresh Exposure' : 'Data Locked'
                        }}
                      </v-btn>
                    </span>
                  </template>
                </v-tooltip>
                <v-chip
                  v-if="!isSelectedFYOpen && financialYearInfo"
                  size="x-small"
                  color="grey"
                  variant="tonal"
                  prepend-icon="mdi-lock-outline"
                >
                  {{ financialYearInfo.financial_year_label }} closed
                </v-chip>
              </v-col>
            </v-row>
            <v-row class="d-flex justify-center">
              <v-col cols="5.1" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="gicCharts"
                  title="Income Statement Components"
                  :data="gicOptions?.data ?? []"
                />
                <ag-charts
                  v-if="conversionOptions"
                  ref="gicCharts"
                  :options="gicOptions"
                ></ag-charts>
              </v-col>
              <v-col cols="5.1" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="revenueCharts"
                  title="Expected Revenue and Claims by Benefit"
                  :data="revenueOptions?.data ?? []"
                />
                <ag-charts
                  v-if="conversionOptions"
                  ref="revenueCharts"
                  :options="revenueOptions"
                ></ag-charts>
              </v-col>
            </v-row>

            <v-row class="mt-4 align-center">
              <v-col cols="3">
                <v-select
                  v-model="selectedBenefit"
                  :items="benefitList"
                  variant="outlined"
                  density="compact"
                  label="Benefit"
                  placeholder="Select Benefit"
                  @update:model-value="getExposureData"
                ></v-select>
              </v-col>
              <v-col cols="4">
                <v-btn-toggle
                  v-model="exposureDataSource"
                  mandatory
                  density="compact"
                  variant="outlined"
                  color="primary"
                  @update:model-value="refreshDashboard"
                >
                  <v-btn value="inforce" size="small">In-Force</v-btn>
                  <v-btn value="all" size="small">All Quotes</v-btn>
                  <v-btn value="quotes" size="small">Quotes Only</v-btn>
                </v-btn-toggle>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="5.1" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="expCharts"
                  :title="`${selectedBenefit} Exposure by Age Band — ${exposureDataSourceLabel}`"
                  :data="exposureOptions?.data ?? []"
                />
                <ag-charts
                  v-if="conversionOptions"
                  ref="expCharts"
                  :options="exposureOptions"
                ></ag-charts>
              </v-col>
              <v-col cols="5.1" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="expGenderCharts"
                  :title="`${selectedBenefit} Exposure by Age Band & Gender — ${exposureDataSourceLabel}`"
                  :data="exposureGenderOptions?.data ?? []"
                />
                <ag-charts
                  v-if="conversionOptions"
                  ref="expGenderCharts"
                  :options="exposureGenderOptions"
                ></ag-charts>
              </v-col>
            </v-row>

            <!-- ── Exposure Trend: Year-over-Year ─────────────────────────── -->
            <v-row class="mt-4">
              <v-col cols="11.2" class="card-bg mx-1">
                <div
                  class="d-flex align-center justify-space-between pr-2 pt-2"
                >
                  <ChartMenu
                    :chart-ref="exposureTrendCharts"
                    :title="`${selectedBenefit} Exposure Trend by Year & Age Band — ${exposureDataSourceLabel}`"
                    :data="exposureTrendOptions?.data ?? []"
                  />
                </div>
                <ag-charts
                  v-if="exposureTrendOptions"
                  ref="exposureTrendCharts"
                  :options="exposureTrendOptions"
                ></ag-charts>
                <div v-else class="text-grey text-body-2 text-center py-6">
                  No trend data available
                </div>
              </v-col>
            </v-row>

            <v-row class="mt-4">
              <v-col cols="5.1" class="card-bg mx-1">
                <div
                  class="d-flex align-center justify-space-between pr-2 pt-2"
                >
                  <ChartMenu
                    :chart-ref="provinceCharts"
                    title="Exposure by Region"
                    :data="provinceBarOptions?.data ?? []"
                  />
                  <v-select
                    v-model="provinceViewBy"
                    :items="viewByOptions"
                    density="compact"
                    variant="outlined"
                    hide-details
                    style="max-width: 160px"
                    label="View by"
                  />
                </div>
                <ag-charts
                  v-if="rawProvinceData.length"
                  ref="provinceCharts"
                  :options="provinceBarOptions"
                ></ag-charts>
                <div v-else class="text-grey text-body-2 text-center py-6">
                  No region data for selected year
                </div>
              </v-col>
              <v-col cols="5.1" class="card-bg mx-1">
                <div
                  class="d-flex align-center justify-space-between pr-2 pt-2"
                >
                  <ChartMenu
                    :chart-ref="industryAgeCharts"
                    title="Exposure by Industry & Age Band"
                    :data="industryAgeOptions?.data ?? []"
                  />
                  <v-select
                    v-model="industryAgeViewBy"
                    :items="viewByOptions"
                    density="compact"
                    variant="outlined"
                    hide-details
                    style="max-width: 160px"
                    label="View by"
                  />
                </div>
                <ag-charts
                  v-if="rawIndustryAgeData.length"
                  ref="industryAgeCharts"
                  :options="industryAgeOptions"
                ></ag-charts>
                <div v-else class="text-grey text-body-2 text-center py-6">
                  No industry/age data for selected year
                </div>
              </v-col>
            </v-row>

            <v-divider class="my-7"></v-divider>

            <!-- ── Section 5: Pricing Intelligence ───────────────────────── -->
            <v-row class="align-center mb-2">
              <v-col>
                <span class="text-subtitle-1 font-weight-bold"
                  >Pricing Intelligence</span
                >
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="5" class="card-bg mx-1">
                <ChartMenu
                  :chart-ref="ratesCharts"
                  title="Avg Rate per R1000 SA by Benefit"
                  :data="ratesBarOptions?.data ?? []"
                />
                <ag-charts
                  v-if="ratesBarOptions"
                  ref="ratesCharts"
                  :options="ratesBarOptions"
                ></ag-charts>
              </v-col>
              <v-col class="mx-1">
                <v-row>
                  <v-col
                    v-for="pm in pricingCards"
                    :key="pm.title"
                    cols="6"
                    class="pb-3"
                  >
                    <v-card variant="tonal" :color="pm.color" class="dash-card">
                      <v-card-subtitle
                        ><h5>{{ pm.title }}</h5></v-card-subtitle
                      >
                      <v-card-text>
                        <v-row>
                          <v-col class="d-flex justify-center">
                            <h2>{{ pm.value }}</h2>
                          </v-col>
                        </v-row>
                      </v-card-text>
                    </v-card>
                  </v-col>
                </v-row>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import ChartMenu from '@/renderer/components/charts/ChartMenu.vue'
import { AgChartOptions } from 'ag-charts-types'
import { AgCharts } from 'ag-charts-vue3'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useStatusBarStore } from '@/renderer/store/statusBar'

const statusBarStore = useStatusBarStore()
const router = useRouter()

const paperTheme = {
  palette: {
    fills: ['#006f9b', '#ff7faa', '#00994d', '#000000', '#00a0dd'],
    strokes: ['#003f58', '#934962', '#004a25', '#914d1d', '#006288']
  }
}

const selectedYear = ref<string | null>(null)
const selectedDataView = ref<string | null>(null)
const conversionViewBy = ref<string>('Count')
const selectedBenefit = ref<string | null>(null)
const exposureDataSource = ref<string>('inforce')

// Section banner visibility
const showQuotesBanner = ref(true)
const showInforceBanner = ref(true)

// Financial year info from backend
const financialYearInfo = ref<any>(null)

// Exposure rebuild state
const rebuildLoading = ref(false)

const gicCharts: any = ref(null)
const revenueCharts: any = ref(null)
const expCharts: any = ref(null)
const expGenderCharts: any = ref(null)
const exposureTrendCharts: any = ref(null)
const trendCharts: any = ref(null)
const funnelCharts: any = ref(null)
const winProbCharts: any = ref(null)
const convRateCharts: any = ref(null)
const industryPipelineCharts: any = ref(null)
const schemeSizeCharts: any = ref(null)
const newQuoteCharts: any = ref(null)
const renewalsCharts: any = ref(null)
const conversionCharts: any = ref(null)
const inForceCharts: any = ref(null)
const ratesCharts: any = ref(null)
const newQuotesTotalCount: any = ref(0)

// Exposure trend chart options
const exposureTrendOptions: any = ref<AgChartOptions | null>(null)

// Sales Pipeline — new chart options
const convRateTrendOptions: any = ref<AgChartOptions | null>(null)
const industryPipelineOptions: any = ref<AgChartOptions | null>(null)
const schemeSizeOptions: any = ref<AgChartOptions | null>(null)

// New section refs
const brokerMetrics: any = ref([])
const pricingCards: any = ref([])

// Win probability model admin
const modelInfo: any = ref(null)
const trainDialog = ref(false)
const trainLoading = ref(false)
const snackbar = ref({ show: false, message: '', color: 'success' })

// Info dialogs — 'auc' | 'training_size' | null
const infoDialog = ref<string | null>(null)
const aucDialogOpen = computed({
  get: () => infoDialog.value === 'auc',
  set: (v) => {
    if (!v) infoDialog.value = null
  }
})
const trainingSizeDialogOpen = computed({
  get: () => infoDialog.value === 'training_size',
  set: (v) => {
    if (!v) infoDialog.value = null
  }
})

const aucInterpretation = computed(() => {
  const auc = modelInfo.value?.auc ?? 0
  if (auc >= 0.85)
    return 'Excellent discrimination — high predictive confidence'
  if (auc >= 0.75) return 'Good discrimination — reliable for prioritisation'
  if (auc >= 0.65) return 'Moderate — useful directional guide'
  if (auc >= 0.5) return 'Weak signal — limited data or high noise'
  return 'No better than random — retrain with more data'
})

const trainingSizeInterpretation = computed(() => {
  const n = modelInfo.value?.training_size ?? 0
  if (n >= 500) return 'High credibility — production-grade confidence'
  if (n >= 100) return 'Good credibility — reliable for business decisions'
  if (n >= 50) return 'Moderate credibility — suitable for pipeline guidance'
  if (n >= 20) return 'Low credibility — indicative only, treat with caution'
  return 'Insufficient data'
})

const trainingSizeAlertType = computed(() => {
  const n = modelInfo.value?.training_size ?? 0
  if (n >= 100) return 'success'
  if (n >= 50) return 'info'
  if (n >= 20) return 'warning'
  return 'error'
})

const accuracyClass = computed(() => {
  const acc = modelInfo.value?.accuracy ?? 0
  if (acc >= 75) return 'accuracy-good'
  if (acc >= 60) return 'accuracy-moderate'
  return 'accuracy-low'
})

const openTrainDialog = () => {
  trainDialog.value = true
}

const triggerTraining = async () => {
  trainLoading.value = true
  trainDialog.value = false
  try {
    await GroupPricingService.trainWinProbabilityModel()
    snackbar.value = {
      show: true,
      message:
        'Training job started. Scores will refresh automatically once complete.',
      color: 'success'
    }
    // Poll model info after a short delay to reflect new model
    setTimeout(async () => {
      try {
        const res = await GroupPricingService.getWinProbabilityModelInfo()
        if (res.data?.data) modelInfo.value = res.data.data
      } catch {}
    }, 5000)
  } catch {
    snackbar.value = {
      show: true,
      message: 'Failed to start training job. Please try again.',
      color: 'error'
    }
  } finally {
    trainLoading.value = false
  }
}

let data: any = null

const availableYears = computed(() => {
  const years: any = []
  const currentYear = new Date().getFullYear()
  for (let i = currentYear; i >= currentYear - 5; i--) {
    years.push(i.toString())
  }
  return years
})

const benefitList = ['All', 'GLA', 'SGLA', 'PTD', 'CI', 'PHI', 'TTD']
const inForceInnerLabel: any = ref<string>('')

onMounted(() => {
  selectedYear.value = new Date().getFullYear().toString()
  selectedBenefit.value = 'All'
  selectedDataView.value = 'Count'
  statusBarStore.set([
    { icon: 'mdi-calendar', text: `Year: ${selectedYear.value}` },
    {
      icon: 'mdi-shield-check-outline',
      text: `Benefit: ${selectedBenefit.value}`
    }
  ])
  refreshDashboard()
  getExposureData()
  getExposureTrendData()
  // Load financial year label for the year selector
  GroupPricingService.getFinancialYearInfo(selectedYear.value)
    .then((res) => {
      if (res.data?.data) financialYearInfo.value = res.data.data
    })
    .catch(() => {})
  // Load model metadata for the win probability card
  GroupPricingService.getWinProbabilityModelInfo()
    .then((res) => {
      if (res.data?.data) modelInfo.value = res.data.data
    })
    .catch(() => {})
})
onUnmounted(() => statusBarStore.clear())

// ── Formatters ────────────────────────────────────────────────────────────────

const formatWithCommas = (value: number) =>
  value.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')

const formatNumber = (num) => {
  if (num >= 1_000_000) {
    return (num / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M'
  } else if (num >= 1_000) {
    return (num / 1_000).toFixed(1).replace(/\.0$/, '') + 'K'
  }
  return num?.toString() ?? '0'
}

const roundUpToTwoDecimals = (num) => {
  const roundedNum = Math.ceil(num * 100) / 100
  return roundedNum
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ')
}

const addFormat = (card: any) => {
  if (card.data_type === 'currency') {
    return `R${formatWithCommas(card.value)}`
  }
  return card.value
}

// ── Chart options (initial / defaults) ───────────────────────────────────────

const conversionOptions: any = ref<AgChartOptions>({
  data: [
    { asset: 'Converted', amount: 39 },
    { asset: 'Unconverted', amount: 61 }
  ],
  title: { text: 'Quote Conversion', fontSize: 14, fontWeight: 'bold' },
  series: [
    {
      type: 'donut',
      calloutLabelKey: 'asset',
      calloutLabel: { enabled: false },
      angleKey: 'amount',
      innerRadiusRatio: 0.6,
      innerLabels: [
        { text: 'Total Quotes', fontWeight: 'bold' },
        { text: '100', spacing: 4, fontSize: 14, color: 'black' }
      ],
      innerCircle: { fill: '#c9fdc9' },
      showInLegend: true
    }
  ]
})

const inForceSchemesOptions: any = ref<AgChartOptions>({
  data: [
    { asset: 'Converted', amount: 39 },
    { asset: 'Unconverted', amount: 61 }
  ],
  title: { text: 'Accepted Quotes', fontSize: 14, fontWeight: 'bold' },
  series: [
    {
      type: 'donut',
      calloutLabelKey: 'asset',
      calloutLabel: { enabled: false },
      angleKey: 'amount',
      innerRadiusRatio: 0.6,
      innerLabels: [
        { text: 'Total Quotes', fontWeight: 'bold' },
        {
          text: `${inForceInnerLabel.value}`,
          spacing: 4,
          fontSize: 14,
          color: 'black'
        }
      ],
      innerCircle: { fill: '#c9fdc9' },
      showInLegend: true
    }
  ]
})

const newQuoteOptions: any = ref<AgChartOptions>({
  data: [
    { asset: 'InProgress', amount: 39 },
    { asset: 'Approved', amount: 61 },
    { asset: 'Accepted', amount: 20 }
  ],
  title: { text: 'New Quotes', fontSize: 14, fontWeight: 'bold' },
  series: [
    {
      type: 'donut',
      calloutLabelKey: 'asset',
      calloutLabel: { enabled: false },
      angleKey: 'amount',
      innerRadiusRatio: 0.6,
      innerLabels: [
        { text: 'Total Quotes', fontWeight: 'bold' },
        {
          text: `${newQuotesTotalCount.value}`,
          spacing: 4,
          fontSize: 14,
          color: 'black'
        }
      ],
      innerCircle: { fill: '#c9fdc9' },
      showInLegend: true
    }
  ]
})

const renewalsOptions: any = ref<AgChartOptions>({
  data: [
    { asset: 'InProgress', amount: 39 },
    { asset: 'Approved', amount: 61 },
    { asset: 'Accepted', amount: 20 }
  ],
  title: { text: 'Renewals', fontSize: 14, fontWeight: 'bold' },
  series: [
    {
      type: 'donut',
      calloutLabelKey: 'asset',
      calloutLabel: { enabled: false },
      angleKey: 'amount',
      innerRadiusRatio: 0.6,
      innerLabels: [
        { text: 'Total Quotes', fontWeight: 'bold' },
        {
          text: `${newQuotesTotalCount.value}`,
          spacing: 4,
          fontSize: 14,
          color: 'black'
        }
      ],
      innerCircle: { fill: '#c9fdc9' },
      showInLegend: true
    }
  ]
})

const revenueOptions: any = ref<AgChartOptions>({
  data: [],
  title: {
    text: 'Expected Revenue & Claims by Benefit',
    fontSize: 14,
    fontWeight: 'bold'
  },
  background: { fill: 'aliceblue' },
  series: [
    { type: 'bar', xKey: 'type', yKey: 'revenue', yName: 'Revenue' },
    { type: 'bar', xKey: 'type', yKey: 'claims', yName: 'Claims' }
  ],
  axes: [
    {
      type: 'number',
      position: 'left',
      label: { formatter: (params) => `R${formatNumber(params.value)}` },
      title: { text: 'Amount', fontSize: 10, fontWeight: 'bold' }
    },
    { type: 'category', position: 'bottom' }
  ]
})

const gicOptions: any = ref<AgChartOptions>({
  data: [],
  title: {
    text: 'Income Statement Components',
    fontSize: 14,
    fontWeight: 'bold'
  },
  background: { fill: 'aliceblue' },
  series: [
    { type: 'bar', xKey: 'type', yKey: 'expected', yName: 'Expected' },
    { type: 'bar', xKey: 'type', yKey: 'actual', yName: 'Actual' }
  ],
  axes: [
    {
      type: 'number',
      position: 'left',
      label: { formatter: (params) => `R${formatNumber(params.value)}` }
    },
    { type: 'category', position: 'bottom' }
  ]
})

const exposureOptions: any = ref<AgChartOptions>({
  theme: paperTheme,
  data: [],
  title: {
    text: `${selectedBenefit.value} Exposure`,
    fontSize: 14,
    fontWeight: 'bold'
  },
  legend: { enabled: true },
  background: { fill: 'aliceblue' },
  series: [
    {
      type: 'bar',
      xKey: 'age_band',
      xName: 'Age Band',
      yKey: 'total_sum_assured',
      yName: 'Total Sum Assured',
      stroke: '#000000',
      strokeWidth: 1
    }
  ],
  axes: [
    {
      type: 'number',
      position: 'left',
      label: { formatter: (params) => `R${formatNumber(params.value)}` },
      title: { text: 'Total Sum Assured', fontSize: 10, fontWeight: 'bold' }
    },
    { type: 'category', position: 'bottom' }
  ]
})

const exposureGenderOptions: any = ref<AgChartOptions>({
  data: [],
  title: {
    text: `${selectedBenefit.value} Exposure`,
    fontSize: 14,
    fontWeight: 'bold'
  },
  background: { fill: 'aliceblue' },
  series: [
    {
      type: 'bar',
      xKey: 'age_band',
      xName: 'Age Band',
      yKey: 'male_sum_assured',
      yName: 'Male Sum Assured'
    },
    {
      type: 'bar',
      xKey: 'age_band',
      xName: 'Age Band',
      yKey: 'female_sum_assured',
      yName: 'Female Sum Assured'
    }
  ],
  axes: [
    {
      type: 'number',
      position: 'left',
      label: { formatter: (params) => `R${formatNumber(params.value)}` },
      title: { text: 'Total Sum Assured', fontSize: 10, fontWeight: 'bold' }
    },
    { type: 'category', position: 'bottom' }
  ]
})

// ── New chart options ─────────────────────────────────────────────────────────

const monthlyTrendOptions: any = ref<AgChartOptions>({
  data: [],
  title: { text: 'Monthly Quote Volume', fontSize: 14, fontWeight: 'bold' },
  background: { fill: 'aliceblue' },
  series: [
    {
      type: 'line',
      xKey: 'month_name',
      yKey: 'new_business_count',
      yName: 'New Business'
    },
    {
      type: 'line',
      xKey: 'month_name',
      yKey: 'renewal_count',
      yName: 'Renewal'
    }
  ],
  axes: [
    {
      type: 'number',
      position: 'left',
      title: { text: 'Quotes', fontSize: 10, fontWeight: 'bold' }
    },
    { type: 'category', position: 'bottom' }
  ]
})

const pipelineFunnelOptions: any = ref<AgChartOptions>({
  data: [],
  title: { text: 'Pipeline by Stage', fontSize: 14, fontWeight: 'bold' },
  background: { fill: 'aliceblue' },
  series: [
    {
      type: 'bar',
      direction: 'horizontal',
      xKey: 'stage_label',
      yKey: 'count',
      yName: 'Quotes'
    }
  ],
  axes: [
    {
      type: 'number',
      position: 'bottom',
      title: { text: 'Quotes', fontSize: 10, fontWeight: 'bold' }
    },
    { type: 'category', position: 'left' }
  ]
})

const winProbDistOptions: any = ref<AgChartOptions | null>(null)

// ── Exposure by Region + Industry by Age ─────────────────────────────────────
const provinceCharts: any = ref(null)
const industryAgeCharts: any = ref(null)

// Raw data stored on dashboard load; chart options are derived via computed
const rawProvinceData = ref<any[]>([])
const rawIndustryAgeData = ref<any[]>([])

// Dropdown selections — default values match what actuaries most commonly need
const provinceViewBy = ref<'Count' | 'Sum Assured'>('Count')
const industryAgeViewBy = ref<'Count' | 'Sum Assured'>('Sum Assured')
const viewByOptions = ['Count', 'Sum Assured']

const exposureDataSourceLabel = computed(() => {
  if (exposureDataSource.value === 'inforce') return 'In-Force'
  if (exposureDataSource.value === 'quotes') return 'Quotes Only'
  return 'All Quotes'
})

// True when today's date falls within the financial year that corresponds to
// the currently selected year.  When false the Refresh Exposure button is
// disabled and locked — past financial years must not be overwritten.
const isSelectedFYOpen = computed(() => {
  const info = financialYearInfo.value
  if (!info) return false // not yet loaded — disable until we know

  const yem: number = info.year_end_month // e.g. 3 = March
  const fyStartYear: number = info.financial_year_start
  const fyEndYear: number = info.financial_year_end

  // First day of the financial year: month after year-end in fyStartYear
  // e.g. year_end_month=3 → FY starts on 1 Apr fyStartYear
  const startMonth = (yem % 12) + 1 // 3 → 4 (April)
  const fyStart = new Date(fyStartYear, startMonth - 1, 1)

  // Last day of the financial year: last day of year_end_month in fyEndYear
  // new Date(year, month, 0) gives the last day of the previous month
  const fyEnd = new Date(fyEndYear, yem, 0)
  fyEnd.setHours(23, 59, 59, 999)

  const today = new Date()
  return today >= fyStart && today <= fyEnd
})

const exposureLockTooltip = computed(() =>
  isSelectedFYOpen.value
    ? 'Rebuild exposure data for this financial year'
    : `${financialYearInfo.value?.financial_year_label ?? 'This financial year'} is closed — historical data is locked`
)

const provinceBarOptions = computed<any>(() => {
  if (!rawProvinceData.value.length) return undefined
  const isCount = provinceViewBy.value === 'Count'
  const yKey = isCount ? 'member_count' : 'total_salary'
  const yLabel = isCount ? 'Member Count' : 'Total Annual Salary (R)'
  return {
    data: rawProvinceData.value,
    title: { text: 'Exposure by Region', fontSize: 14, fontWeight: 'bold' },
    background: { fill: 'aliceblue' },
    series: [
      {
        type: 'bar',
        direction: 'horizontal',
        xKey: 'region',
        yKey,
        yName: isCount ? 'Members' : 'Total Salary',
        fill: '#006f9b',
        stroke: '#000000',
        strokeWidth: 1
      }
    ],
    axes: [
      { type: 'category', position: 'left' },
      {
        type: 'number',
        position: 'bottom',
        title: { text: yLabel, fontSize: 10, fontWeight: 'bold' },
        label: isCount
          ? {}
          : { formatter: (p: any) => `R${formatNumber(p.value)}` }
      }
    ]
  }
})

const industryAgeOptions = computed<any>(() => {
  if (!rawIndustryAgeData.value.length) return undefined
  const isCount = industryAgeViewBy.value === 'Count'
  const raw = rawIndustryAgeData.value
  const valueKey = isCount ? 'record_count' : 'total_sum_assured'

  // Unique age bands in ascending order (backend already orders by min_age)
  const bandSet: string[] = []
  for (const r of raw) {
    if (!bandSet.includes(r.age_band)) bandSet.push(r.age_band)
  }

  // Unique industries — become the grouped bar series (one colour per industry)
  const industrySet: string[] = []
  for (const r of raw) {
    if (r.industry && !industrySet.includes(r.industry))
      industrySet.push(r.industry)
  }

  // Pivot: one row per age band, one column per industry
  const pivotData = bandSet.map((band) => {
    const row: any = { age_band: band }
    for (const ind of industrySet) {
      const match = raw.find(
        (r: any) => r.age_band === band && r.industry === ind
      )
      row[ind] = match?.[valueKey] ?? 0
    }
    return row
  })

  const industryColors = [
    '#006f9b',
    '#ff7f00',
    '#00994d',
    '#d62728',
    '#9467bd',
    '#8c564b',
    '#e377c2',
    '#7f7f7f',
    '#bcbd22',
    '#17becf'
  ]
  const yTitle = isCount ? 'Number of Lives Exposed' : 'Total Sum Assured (R)'
  return {
    data: pivotData,
    title: {
      text: 'Exposure by Industry & Age Band',
      fontSize: 14,
      fontWeight: 'bold'
    },
    background: { fill: 'aliceblue' },
    legend: {
      enabled: true,
      position: 'right',
      item: {
        label: { fontSize: 12, fontWeight: 'bold', color: '#333' },
        marker: { size: 14 },
        paddingY: 6
      }
    },
    series: industrySet.map((ind, i) => ({
      type: 'bar',
      xKey: 'age_band',
      yKey: ind,
      yName: ind,
      fill: industryColors[i % industryColors.length],
      stroke: '#000000',
      strokeWidth: 1
    })),
    axes: [
      {
        type: 'category',
        position: 'bottom',
        title: { text: 'Age Band', fontSize: 10, fontWeight: 'bold' }
      },
      {
        type: 'number',
        position: 'left',
        title: { text: yTitle, fontSize: 10, fontWeight: 'bold' },
        label: isCount
          ? {}
          : { formatter: (p: any) => `R${formatNumber(p.value)}` }
      }
    ]
  }
})

const ratesBarOptions: any = ref<AgChartOptions>({
  data: [],
  title: {
    text: 'Avg Rate per R1,000 SA by Benefit',
    fontSize: 14,
    fontWeight: 'bold'
  },
  background: { fill: 'aliceblue' },
  series: [
    { type: 'bar', xKey: 'benefit', yKey: 'rate', yName: 'Rate per R1,000' }
  ],
  axes: [
    {
      type: 'number',
      position: 'left',
      title: { text: 'Rate (R)', fontSize: 10, fontWeight: 'bold' }
    },
    { type: 'category', position: 'bottom' }
  ]
})

// ── KPI cards default ─────────────────────────────────────────────────────────

const cards = ref<any[]>([
  {
    title: 'Annual Premium',
    value: 0,
    flex: 2,
    route: 'group-pricing-schemes',
    data_type: 'currency'
  },
  {
    title: 'Scheme Count',
    value: 0,
    flex: 2,
    route: 'group-pricing-schemes',
    data_type: 'number'
  },
  {
    title: 'Expected Claims',
    value: 0,
    flex: 2,
    route: 'group-pricing-claims-management',
    data_type: 'number'
  },
  {
    title: 'Expense Recovery',
    value: 0,
    flex: 2,
    route: null,
    data_type: 'number'
  },
  {
    title: 'Conversion Rate',
    value: '0%',
    flex: 2,
    route: null,
    data_type: 'plain'
  },
  {
    title: 'Expected Loss Ratio',
    value: '0%',
    flex: 2,
    route: null,
    data_type: 'plain'
  }
])

// ── changeConversionDataSource — updates Conversion & Quotes Overview donuts ──

const changeConversionDataSource = () => {
  if (!data) return

  const isPremium = conversionViewBy.value === 'Annual Premium'

  let unconvertedQuotes = 0
  let totalQuotes = 0
  let convertedQuotesAccepted = 0
  let inForceSchemesRenewal = 0
  let inForceSchemesNew = 0
  let inForceSchemesTotal = ''
  let newQuotesInProgress = 0
  let newQuotesPendingReview = 0
  let newQuotesApproved = 0
  let newQuotesAccepted = 0
  let totalnewQuotes = ''
  let renewalQuotesInProgress = 0
  let renewalQuotesPendingReview = 0
  let renewalQuotesApproved = 0
  let renewalQuotesAccepted = 0
  let totalrenewalQuotes = ''

  if (isPremium) {
    convertedQuotesAccepted =
      data.total_quotes_converted_premium_accepted / 1000000
    unconvertedQuotes = data.total_quotes_unconverted_premium / 1000000
    totalQuotes = data.conversion_total_premium / 1000000
    inForceSchemesRenewal = data.renewals_in_force_premium / 1000000
    inForceSchemesNew = data.new_business_in_force_premium / 1000000
    inForceSchemesTotal = `${roundUpToTwoDecimals(data.total_in_force_premium / 1000000)}m`
    newQuotesInProgress = data.new_quotes_in_progress_premium / 1000000
    newQuotesPendingReview =
      (data.new_quotes_pending_review_premium ?? 0) / 1000000
    newQuotesApproved = data.new_quotes_approved_premium / 1000000
    newQuotesAccepted =
      (data.new_quotes_in_force_premium + data.new_quotes_accepted_premium) /
      1000000
    totalnewQuotes = `${roundUpToTwoDecimals(newQuotesInProgress + newQuotesPendingReview + newQuotesApproved + newQuotesAccepted)}m`
    renewalQuotesInProgress = data.renewals_quotes_in_progress_premium / 1000000
    renewalQuotesPendingReview =
      (data.renewals_quotes_pending_review_premium ?? 0) / 1000000
    renewalQuotesApproved = data.renewals_quotes_approved_premium / 1000000
    renewalQuotesAccepted =
      (data.renewals_quotes_in_force_premium +
        data.renewals_quotes_accepted_premium) /
      1000000
    totalrenewalQuotes = `${roundUpToTwoDecimals(renewalQuotesInProgress + renewalQuotesPendingReview + renewalQuotesApproved + renewalQuotesAccepted)}m`
  } else {
    convertedQuotesAccepted = data.total_quotes_converted_count_accepted
    unconvertedQuotes = data.total_quotes_unconverted_count
    totalQuotes = data.conversion_total_count
    inForceSchemesRenewal = data.renewals_in_force_count
    inForceSchemesNew = data.new_business_in_force_count
    inForceSchemesTotal = data.total_in_force_count
    newQuotesInProgress = data.new_quotes_in_progress_count
    newQuotesPendingReview = data.new_quotes_pending_review_count ?? 0
    newQuotesApproved = data.new_quotes_approved_count
    newQuotesAccepted =
      data.new_quotes_in_force_count + data.new_quotes_accepted_count
    totalnewQuotes = roundUpToTwoDecimals(
      newQuotesInProgress +
        newQuotesPendingReview +
        newQuotesApproved +
        newQuotesAccepted
    )
    renewalQuotesInProgress = data.renewals_quotes_in_progress_count
    renewalQuotesPendingReview = data.renewals_quotes_pending_review_count ?? 0
    renewalQuotesApproved = data.renewals_quotes_approved_count
    renewalQuotesAccepted =
      data.renewals_quotes_in_force_count + data.renewals_quotes_accepted_count
    totalrenewalQuotes = roundUpToTwoDecimals(
      renewalQuotesInProgress +
        renewalQuotesPendingReview +
        renewalQuotesApproved +
        renewalQuotesAccepted
    )
  }

  conversionOptions.value = {
    ...conversionOptions.value,
    data: [
      { asset: 'accepted', amount: convertedQuotesAccepted },
      { asset: 'approved', amount: unconvertedQuotes }
    ],
    series: [
      {
        ...conversionOptions.value.series[0],
        innerLabels: [
          {
            text: `${roundUpToTwoDecimals((convertedQuotesAccepted / (totalQuotes || 1)) * 100)}%`,
            spacing: 4,
            fontSize: 14,
            color: 'black'
          }
        ]
      }
    ]
  }

  inForceSchemesOptions.value = {
    ...inForceSchemesOptions.value,
    data: [
      { asset: 'Renewal', amount: inForceSchemesRenewal },
      { asset: 'New Business', amount: inForceSchemesNew }
    ],
    series: [
      {
        ...inForceSchemesOptions.value.series[0],
        innerLabels: [
          {
            text: `${inForceSchemesTotal}`,
            spacing: 4,
            fontSize: 14,
            color: 'black'
          }
        ]
      }
    ]
  }

  const statusLegendConfig = {
    position: 'bottom' as const,
    orientation: 'horizontal' as const,
    item: {
      label: { fontSize: 10 },
      marker: { size: 8 },
      paddingX: 10,
      paddingY: 4
    }
  }

  const statusFills = ['#5B8DB8', '#F0A500', '#4CAF82', '#2E7D5E']

  newQuoteOptions.value = {
    ...newQuoteOptions.value,
    data: [
      { asset: 'In Progress', amount: newQuotesInProgress },
      { asset: 'Under Review', amount: newQuotesPendingReview },
      { asset: 'Approved', amount: newQuotesApproved },
      { asset: 'Accepted', amount: newQuotesAccepted }
    ],
    legend: statusLegendConfig,
    series: [
      {
        type: 'donut',
        calloutLabelKey: 'asset',
        calloutLabel: { enabled: false },
        angleKey: 'amount',
        fills: statusFills,
        innerRadiusRatio: 0.6,
        innerLabels: [
          { text: totalnewQuotes, spacing: 4, fontSize: 13, color: 'black' }
        ],
        innerCircle: { fill: '#f4f9f4' },
        showInLegend: true
      }
    ]
  }

  renewalsOptions.value = {
    ...renewalsOptions.value,
    data: [
      { asset: 'In Progress', amount: renewalQuotesInProgress },
      { asset: 'Under Review', amount: renewalQuotesPendingReview },
      { asset: 'Approved', amount: renewalQuotesApproved },
      { asset: 'Accepted', amount: renewalQuotesAccepted }
    ],
    legend: statusLegendConfig,
    series: [
      {
        type: 'donut',
        calloutLabelKey: 'asset',
        calloutLabel: { enabled: false },
        angleKey: 'amount',
        fills: statusFills,
        innerRadiusRatio: 0.6,
        innerLabels: [
          { text: totalrenewalQuotes, spacing: 4, fontSize: 13, color: 'black' }
        ],
        innerCircle: { fill: '#f4f9f4' },
        showInLegend: true
      }
    ]
  }
}

// ── changeChartDataSource — updates Sales Pipeline charts ─────────────────────

const changeChartDataSource = () => {
  if (!data) return

  const isPremium = selectedDataView.value === 'Annual Premium'

  // Monthly trend toggle
  if (data.monthly_quote_trend) {
    monthlyTrendOptions.value = {
      ...monthlyTrendOptions.value,
      data: data.monthly_quote_trend,
      series: isPremium
        ? [
            {
              type: 'line',
              xKey: 'month_name',
              yKey: 'new_business_premium',
              yName: 'New Business (R)'
            },
            {
              type: 'line',
              xKey: 'month_name',
              yKey: 'renewal_premium',
              yName: 'Renewal (R)'
            }
          ]
        : [
            {
              type: 'line',
              xKey: 'month_name',
              yKey: 'new_business_count',
              yName: 'New Business'
            },
            {
              type: 'line',
              xKey: 'month_name',
              yKey: 'renewal_count',
              yName: 'Renewal'
            }
          ],
      axes: [
        {
          type: 'number',
          position: 'left',
          title: {
            text: isPremium ? 'Premium (R)' : 'Quotes',
            fontSize: 10,
            fontWeight: 'bold'
          },
          ...(isPremium
            ? { label: { formatter: (p) => `R${formatNumber(p.value)}` } }
            : {})
        },
        { type: 'category', position: 'bottom' }
      ]
    }
  }

  // Funnel toggle
  if (data.quote_funnel) {
    pipelineFunnelOptions.value = {
      ...pipelineFunnelOptions.value,
      data: data.quote_funnel.map((s) => ({
        ...s,
        stage_label: s.stage.replace(/_/g, ' ')
      })),
      series: [
        {
          type: 'bar',
          direction: 'horizontal',
          xKey: 'stage_label',
          yKey: isPremium ? 'premium' : 'count',
          yName: isPremium ? 'Premium (R)' : 'Quotes'
        }
      ],
      axes: [
        {
          type: 'number',
          position: 'bottom',
          title: {
            text: isPremium ? 'Premium (R)' : 'Quotes',
            fontSize: 10,
            fontWeight: 'bold'
          },
          ...(isPremium
            ? { label: { formatter: (p) => `R${formatNumber(p.value)}` } }
            : {})
        },
        { type: 'category', position: 'left' }
      ]
    }
  }
}

// ── getExposureData ───────────────────────────────────────────────────────────

const getExposureData = async () => {
  try {
    if (selectedBenefit.value && selectedYear.value) {
      const res = await GroupPricingService.getExposureData(
        selectedYear.value,
        selectedBenefit.value,
        exposureDataSource.value
      )
      const sourceLabel = exposureDataSourceLabel.value
      const titleText = `${selectedBenefit.value} Exposure (${selectedYear.value}) — ${sourceLabel}`
      exposureOptions.value = {
        ...exposureOptions.value,
        title: { text: titleText, fontSize: 14, fontWeight: 'bold' },
        data: res.data || []
      }
      exposureGenderOptions.value = {
        ...exposureGenderOptions.value,
        title: { text: titleText, fontSize: 14, fontWeight: 'bold' },
        data: res.data || []
      }
    }
  } catch (error) {
    console.error(error)
  }
}

// ── getExposureTrendData ──────────────────────────────────────────────────────

const getExposureTrendData = async () => {
  try {
    const res = await GroupPricingService.getExposureTrend(
      selectedBenefit.value || 'All',
      exposureDataSource.value
    )
    const rows: any[] = res.data?.data || []
    if (!rows.length) {
      exposureTrendOptions.value = null
      return
    }

    // Collect unique age bands (preserving min_age order)
    const ageBandOrder: Map<string, number> = new Map()
    rows.forEach((r) => {
      if (!ageBandOrder.has(r.age_band)) ageBandOrder.set(r.age_band, r.min_age)
    })
    const ageBands = [...ageBandOrder.entries()]
      .sort((a, b) => a[1] - b[1])
      .map((e) => e[0])

    // Collect unique financial years
    const years = [...new Set(rows.map((r) => r.financial_year))].sort()

    // Pivot: for each year build an object keyed by age_band
    const pivoted: Record<number, any> = {}
    years.forEach((y) => {
      pivoted[y] = { year: String(y) }
      ageBands.forEach((ab) => {
        pivoted[y][ab] = 0
      })
    })
    rows.forEach((r) => {
      if (pivoted[r.financial_year]) {
        pivoted[r.financial_year][r.age_band] = r.total_sum_assured
      }
    })

    const chartData = years.map((y) => pivoted[y])
    const sourceLabel = exposureDataSourceLabel.value

    exposureTrendOptions.value = {
      data: chartData,
      title: {
        text: `${selectedBenefit.value || 'All'} Exposure Trend by Year & Age Band — ${sourceLabel}`,
        fontSize: 14,
        fontWeight: 'bold'
      },
      series: ageBands.map((ab) => ({
        type: 'bar',
        xKey: 'year',
        yKey: ab,
        yName: ab,
        stacked: false
      })),
      axes: [
        {
          type: 'category',
          position: 'bottom',
          title: { text: 'Financial Year' }
        },
        {
          type: 'number',
          position: 'left',
          title: { text: 'Sum Assured (R)' },
          label: {
            formatter: (params: any) => `R${formatNumber(params.value)}`
          }
        }
      ]
    }
  } catch (error) {
    console.error(error)
  }
}

// ── triggerExposureRebuild ────────────────────────────────────────────────────

const triggerExposureRebuild = async () => {
  if (!selectedYear.value) return
  rebuildLoading.value = true
  try {
    await GroupPricingService.rebuildExposureData(selectedYear.value)
    await getExposureData()
    await getExposureTrendData()
  } catch (error) {
    console.error(error)
  } finally {
    rebuildLoading.value = false
  }
}

// ── refreshDashboard ──────────────────────────────────────────────────────────

const refreshDashboard = async () => {
  if (!selectedYear.value) return

  // Refresh financial year label alongside dashboard data
  GroupPricingService.getFinancialYearInfo(selectedYear.value)
    .then((res) => {
      if (res.data?.data) financialYearInfo.value = res.data.data
    })
    .catch(() => {})

  const res = await GroupPricingService.getDashboardData(
    selectedYear.value,
    exposureDataSource.value
  )
  data = res.data

  // Existing charts
  changeConversionDataSource()
  changeChartDataSource()
  gicOptions.value = {
    ...gicOptions.value,
    data: res.data.income_statement_components
  }
  revenueOptions.value = {
    ...revenueOptions.value,
    data: res.data.revenue_benefits
  }

  // KPI cards from backend — normalise flex to 2 so all 6 fit in one row
  const backendCards = (res.data.card_data ?? []).map((c: any) => ({
    ...c,
    flex: 2
  }))

  // Computed KPI: Conversion Rate
  const convTotal = data.conversion_total_count || 0
  const convAccepted = data.total_quotes_converted_count_accepted || 0
  const convPct =
    convTotal > 0 ? roundUpToTwoDecimals((convAccepted / convTotal) * 100) : 0

  // Computed KPI: Expected Loss Ratio
  const elr = data.pricing_metrics?.expected_loss_ratio ?? 0

  cards.value = [
    ...backendCards,
    {
      title: 'Conversion Rate',
      value: `${convPct}%`,
      flex: 2,
      route: null,
      data_type: 'plain'
    },
    {
      title: 'Expected Loss Ratio',
      value: `${elr}%`,
      flex: 2,
      route: null,
      data_type: 'plain'
    }
  ]

  // Broker metrics table
  brokerMetrics.value = data.broker_metrics || []

  // Pricing intelligence section
  if (data.pricing_metrics) {
    const pm = data.pricing_metrics

    // Rates bar
    ratesBarOptions.value = {
      ...ratesBarOptions.value,
      data: [
        { benefit: 'GLA', rate: pm.avg_gla_rate_per_1000 },
        { benefit: 'PTD', rate: pm.avg_ptd_rate_per_1000 },
        { benefit: 'CI', rate: pm.avg_ci_rate_per_1000 },
        { benefit: 'SGLA', rate: pm.avg_sgla_rate_per_1000 }
      ]
    }

    // Pricing mini-cards
    const alr = pm.actual_loss_ratio ?? 0
    pricingCards.value = [
      { title: 'Avg Discount', value: `${pm.avg_discount}%`, color: 'primary' },
      {
        title: 'Avg Commission',
        value: `${pm.avg_commission_pct}%`,
        color: 'primary'
      },
      {
        title: 'Expected Loss Ratio',
        value: `${elr}%`,
        color: elr > 100 ? 'error' : elr > 80 ? 'warning' : 'success'
      },
      {
        title: 'Actual Loss Ratio',
        value: `${alr}%`,
        color: alr > 100 ? 'error' : alr > 80 ? 'warning' : 'success'
      }
    ]
  }

  // Win probability distribution — one series per band so each bar gets its own colour
  if (data.win_probability_bands) {
    const bands = data.win_probability_bands
    const bandConfig = [
      { key: '0-20', band: '0–20%', color: '#d32f2f' },
      { key: '20-40', band: '20–40%', color: '#f57c00' },
      { key: '40-60', band: '40–60%', color: '#0288d1' },
      { key: '60-80', band: '60–80%', color: '#388e3c' },
      { key: '80-100', band: '80–100%', color: '#1b5e20' }
    ]
    const bandData = bandConfig.map((b) => ({
      band: b.band,
      count: bands[b.key] ?? 0
    }))
    winProbDistOptions.value = {
      data: bandData,
      title: {
        text: 'Win Probability Distribution (Active Quotes)',
        fontSize: 14,
        fontWeight: 'bold'
      },
      background: { fill: 'aliceblue' },
      series: bandConfig.map((b) => ({
        type: 'bar',
        xKey: 'band',
        yKey: 'count',
        yName: b.band,
        fill: b.color,
        stroke: b.color,
        showInLegend: false
      })),
      axes: [
        {
          type: 'number',
          position: 'left',
          title: { text: 'Quotes', fontSize: 10 },
          min: 0
        },
        {
          type: 'category',
          position: 'bottom',
          title: { text: 'Probability Band', fontSize: 10 }
        }
      ]
    }
  }

  // Store raw data; computed options rebuild automatically when view-by changes
  rawProvinceData.value = data.exposure_by_province ?? []
  rawIndustryAgeData.value = data.industry_by_age ?? []

  // ── Monthly Conversion Rate Trend ────────────────────────────────────────
  if (data.monthly_conversion_trend?.length) {
    const convRows = data.monthly_conversion_trend
    // Only render months that have at least one closed quote in either type
    const hasData = convRows.some((r: any) => r.nb_total > 0 || r.ren_total > 0)
    convRateTrendOptions.value = hasData
      ? {
          data: convRows,
          title: {
            text: 'Monthly Conversion Rate — New Business vs Renewal',
            fontSize: 14,
            fontWeight: 'bold'
          },
          background: { fill: 'aliceblue' },
          series: [
            {
              type: 'line',
              xKey: 'month_name',
              yKey: 'nb_conv_rate',
              yName: 'New Business',
              stroke: '#006f9b',
              marker: { enabled: true, fill: '#006f9b' },
              tooltip: {
                renderer: (params: any) =>
                  `<b>${params.datum.month_name}</b><br/>NB: ${params.datum.nb_conv_rate}% (${params.datum.nb_accepted}/${params.datum.nb_total})`
              }
            },
            {
              type: 'line',
              xKey: 'month_name',
              yKey: 'ren_conv_rate',
              yName: 'Renewal',
              stroke: '#ff7faa',
              marker: { enabled: true, fill: '#ff7faa' },
              tooltip: {
                renderer: (params: any) =>
                  `<b>${params.datum.month_name}</b><br/>Renewal: ${params.datum.ren_conv_rate}% (${params.datum.ren_accepted}/${params.datum.ren_total})`
              }
            }
          ],
          axes: [
            {
              type: 'category',
              position: 'bottom',
              title: { text: 'Month' }
            },
            {
              type: 'number',
              position: 'left',
              title: { text: 'Conversion Rate (%)' },
              min: 0,
              max: 100,
              label: { formatter: (p: any) => `${p.value}%` }
            }
          ],
          legend: { enabled: true }
        }
      : null
  } else {
    convRateTrendOptions.value = null
  }

  // ── Industry Mix of Active Pipeline ──────────────────────────────────────
  if (data.industry_pipeline?.length) {
    const ipRows: any[] = data.industry_pipeline
    industryPipelineOptions.value = {
      data: ipRows,
      title: {
        text: 'Industry Mix — Active Pipeline (Quote Count & Annual Premium)',
        fontSize: 14,
        fontWeight: 'bold'
      },
      background: { fill: 'aliceblue' },
      series: [
        {
          type: 'bar',
          direction: 'horizontal',
          xKey: 'industry',
          yKey: 'quote_count',
          yName: 'Quote Count',
          fill: '#006f9b'
        },
        {
          type: 'bar',
          direction: 'horizontal',
          xKey: 'industry',
          yKey: 'total_premium',
          yName: 'Annual Premium (R)',
          fill: '#00994d'
        }
      ],
      axes: [
        { type: 'category', position: 'left' },
        {
          type: 'number',
          position: 'bottom',
          title: { text: 'Count / Premium' },
          label: { formatter: (p: any) => formatNumber(p.value) }
        }
      ],
      legend: { enabled: true }
    }
  } else {
    industryPipelineOptions.value = null
  }

  // ── Scheme Size Distribution ──────────────────────────────────────────────
  if (data.scheme_size_distribution?.length) {
    const ssRows: any[] = data.scheme_size_distribution
    schemeSizeOptions.value = {
      data: ssRows,
      title: {
        text: 'Scheme Size Distribution — Active Pipeline (Lives)',
        fontSize: 14,
        fontWeight: 'bold'
      },
      background: { fill: 'aliceblue' },
      series: [
        {
          type: 'bar',
          xKey: 'size_band',
          yKey: 'quote_count',
          yName: 'Quotes',
          fill: '#006f9b',
          tooltip: {
            renderer: (params: any) =>
              `<b>${params.datum.size_band} lives</b><br/>Quotes: ${params.datum.quote_count}<br/>Premium: R${formatNumber(params.datum.total_premium)}`
          }
        }
      ],
      axes: [
        {
          type: 'category',
          position: 'bottom',
          title: { text: 'Lives (Members)' }
        },
        {
          type: 'number',
          position: 'left',
          title: { text: 'Number of Quotes' },
          min: 0
        }
      ],
      legend: { enabled: false }
    }
  } else {
    schemeSizeOptions.value = null
  }

  getExposureData()
  getExposureTrendData()
}
</script>

<style lang="css" scoped>
.model-meta-chip {
  display: inline-flex;
  align-items: center;
  padding: 2px 6px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.05);
}

.model-meta-clickable {
  cursor: pointer;
  text-decoration: underline dotted;
}
.model-meta-clickable:hover {
  background: rgba(0, 111, 155, 0.12);
  color: #006f9b;
}

.accuracy-good {
  color: #2e7d32;
}
.accuracy-moderate {
  color: #e65100;
}
.accuracy-low {
  color: #c62828;
}

.card-bg {
  background-color: #f5f5f5;
  border: 1px solid #dadada;
  border-radius: 3px;
  margin: 0px;
  margin-bottom: 16px;
  padding: 0px;
}

.pane-bg {
  background-color: #f5f5f5;
  border: 1px solid #f5f5f5;
  border-radius: 3px;
  margin: 0px;
  padding: 0px;
}

.dash-card {
  border-radius: 10px;
  border: 1px solid #b2cbe1;
}

.broker-table thead th {
  font-weight: 600;
  background-color: #eaf3fb;
}
</style>
