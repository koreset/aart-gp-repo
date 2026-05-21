<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center flex-wrap">
              <div>
                <span class="analytics-title">Claims Analytics</span>
                <div class="analytics-subtitle">
                  Claims performance metrics, trends, and analysis across all
                  schemes.
                </div>
              </div>
            </div>
          </template>
          <template #default>
            <claims-analytics-dashboard :schemes="schemes" />
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ClaimsAnalyticsDashboard from './components/ClaimsAnalyticsDashboard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const schemes = ref<any[]>([])

onMounted(async () => {
  try {
    // Coverage-history universe — currently or previously in-force schemes
    // only, so the analytics scheme filter doesn't expose pre-quote schemes.
    const res = await GroupPricingService.getSchemesWithCoverageHistory()
    schemes.value = res.data?.data ?? res.data ?? []
  } catch {
    schemes.value = []
  }
})
</script>

<style scoped>
.analytics-title {
  font-size: 1.15rem;
  font-weight: 700;
  letter-spacing: 0.2px;
  color: #ffffff;
}

.analytics-subtitle {
  font-size: 0.78rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.88);
  margin-top: 2px;
}
</style>
