<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <span class="headline">Fraud Risk Configuration</span>
              <v-tooltip location="right">
                <template #activator="{ props }">
                  <v-icon v-bind="props" size="small" color="info" class="ml-2">
                    mdi-information-outline
                  </v-icon>
                </template>
                <span>
                  GLM scores each claim's fraud probability from historical
                  patterns. Company rules override the GLM when matched.
                </span>
              </v-tooltip>
            </div>
          </template>
          <template #default>
            <v-tabs v-model="tab" color="primary" grow>
              <v-tab value="model">GLM Coefficients</v-tab>
              <v-tab value="rules">Risk Classification Rules</v-tab>
            </v-tabs>
            <v-window v-model="tab" class="mt-4">
              <v-window-item value="model">
                <FraudRiskCoefficients />
              </v-window-item>
              <v-window-item value="rules">
                <FraudRiskRules />
              </v-window-item>
            </v-window>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import FraudRiskCoefficients from '@/renderer/screens/group_pricing/administration/components/FraudRiskCoefficients.vue'
import FraudRiskRules from '@/renderer/screens/group_pricing/administration/components/FraudRiskRules.vue'

const tab = ref<'model' | 'rules'>('model')
</script>
