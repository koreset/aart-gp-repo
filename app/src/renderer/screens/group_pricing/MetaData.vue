<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12">
        <base-card :show-actions="false">
          <template #header>
            <div class="section-header">
              <v-icon icon="mdi-cog-outline" class="mr-3" />
              Group Pricing Configuration
            </div>
          </template>
          <template #default>
            <p class="text-subtitle-1 text-medium-emphasis mb-6">
              Configure insurer details, broker relationships, and benefit
              customization settings
            </p>
            <v-expansion-panels variant="accordion" multiple>
              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-office-building" class="mr-2" />
                  Insurer Details
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <insurer-data-form />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-account-group" class="mr-2" />
                  Broker Management
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <broker-management />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-bank" class="mr-2" />
                  Reinsurer Management
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <reinsurer-management />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-tag-multiple" class="mr-2" />
                  Scheme Categories
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <scheme-category-management />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-tune" class="mr-2" />
                  Benefits Customization
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <benefits-customization />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-percent-outline" class="mr-2" />
                  Discount Calculation Method
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <discount-method-settings />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-shield-check-outline" class="mr-2" />
                  Free Cover Limit Calculation Method
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <free-cover-limit-method-settings />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-medical-bag" class="mr-2" />
                  Medical Aid Waiver Calculation Method
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <medical-aid-waiver-method-settings />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-calculator-variant" class="mr-2" />
                  PTD Base Rate Calculation Method
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <ptd-base-rate-method-settings />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-alert-octagon-outline" class="mr-2" />
                  Risk Watchlist Thresholds
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <risk-watchlist-threshold-settings />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-file-word-outline" class="mr-2" />
                  Quote Template
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <quote-template-management />
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon icon="mdi-file-document-check-outline" class="mr-2" />
                  On Risk Letter Template
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <on-risk-letter-template-management />
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="hideNotification"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import InsurerDataForm from '@/renderer/components/grouppricing/InsurerDataForm.vue'
import BrokerManagement from '@/renderer/components/grouppricing/BrokerManagement.vue'
import ReinsurerManagement from '@/renderer/components/grouppricing/ReinsurerManagement.vue'
import SchemeCategoryManagement from '@/renderer/components/grouppricing/SchemeCategoryManagement.vue'
import BenefitsCustomization from '@/renderer/components/grouppricing/BenefitsCustomization.vue'
import DiscountMethodSettings from '@/renderer/components/grouppricing/DiscountMethodSettings.vue'
import FreeCoverLimitMethodSettings from '@/renderer/components/grouppricing/FreeCoverLimitMethodSettings.vue'
import MedicalAidWaiverMethodSettings from '@/renderer/components/grouppricing/MedicalAidWaiverMethodSettings.vue'
import PtdBaseRateMethodSettings from '@/renderer/components/grouppricing/PtdBaseRateMethodSettings.vue'
import RiskWatchlistThresholdSettings from '@/renderer/components/grouppricing/RiskWatchlistThresholdSettings.vue'
import QuoteTemplateManagement from '@/renderer/components/grouppricing/QuoteTemplateManagement.vue'
import OnRiskLetterTemplateManagement from '@/renderer/components/grouppricing/OnRiskLetterTemplateManagement.vue'
import { useNotifications } from '@/renderer/composables/useNotifications'
import BaseCard from '@/renderer/components/BaseCard.vue'
// Use the notifications composable for the snackbar
const { snackbar, snackbarMessage, snackbarColor, hideNotification } =
  useNotifications()
</script>

<style scoped>
.metadata-container {
  padding: 2rem 1rem;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.page-header {
  text-align: center;
  padding: 2rem 0;
}

.sections-container {
  max-width: none;
}

/* Enhanced card styling for child components */
:deep(.v-card) {
  border-radius: 5px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.v-card:hover) {
  transform: translateY(-4px);
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15) !important;
}

/* Section headers styling */
:deep(h4) {
  color: #1976d2;
  font-weight: 600;
  font-size: 1.25rem;
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Form field enhancements */
:deep(.v-text-field) {
  margin-bottom: 8px;
}

:deep(.v-text-field .v-field__outline) {
  border-radius: 8px;
}

/* Button enhancements */
:deep(.v-btn) {
  border-radius: 8px;
  text-transform: none;
  font-weight: 500;
  letter-spacing: 0.5px;
}

/* Data grid styling */
:deep(.ag-theme-balham) {
  border-radius: 5px;
  overflow: hidden;
}

/* Animation for section transitions
.sections-container > * {
  animation: fadeInUp 0.6s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
  */

/* Responsive design */
@media (max-width: 768px) {
  .metadata-container {
    padding: 1rem 0.5rem;
  }

  .page-header {
    padding: 1rem 0;
  }

  .page-header h1 {
    font-size: 1.5rem;
  }
}
</style>
