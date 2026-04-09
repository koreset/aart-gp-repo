<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center flex-wrap">
              <div>
                <span class="headline">Claims Analytics</span>
                <div class="text-caption text-medium-emphasis">
                  Claims performance metrics, trends, and analysis across all
                  schemes.
                </div>
              </div>
              <v-btn
                size="small"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="
                  router.push({ name: 'group-pricing-claims-management' })
                "
              >
                Back to Claims
              </v-btn>
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
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ClaimsAnalyticsDashboard from './components/ClaimsAnalyticsDashboard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const router = useRouter()
const schemes = ref<any[]>([])

onMounted(async () => {
  try {
    const res = await GroupPricingService.getSchemesInforce()
    schemes.value = res.data?.data ?? res.data ?? []
  } catch {
    schemes.value = []
  }
})
</script>
