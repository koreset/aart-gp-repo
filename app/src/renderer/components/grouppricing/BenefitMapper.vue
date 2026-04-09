<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Customize Benefit Names</span>
          </template>
          <template #default>
            <v-row
              v-for="(benefit, index) in benefits"
              :key="index"
              class="d-flex my-2"
            >
              <v-col cols="3">
                <strong>{{ benefit.benefit_name }}</strong>
              </v-col>
              <v-col cols="4">
                <v-text-field
                  v-model="benefit.benefit_alias"
                  :label="`Alias for ${benefit.benefit_name}`"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.alias]"
                  hide-details="auto"
                />
              </v-col>
              <v-col v-if="benefit.benefit_alias" cols="4">
                <v-text-field
                  v-model="benefit.benefit_alias_code"
                  :label="`Code for ${benefit.benefit_alias}`"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.alias_code]"
                  hide-details="auto"
                />
              </v-col>
            </v-row>
            <v-btn
              :disabled="!isValid || loading"
              :loading="loading"
              rounded
              size="small"
              color="primary"
              class="mt-4"
              @click="submitBenefits"
            >
              Save Mappings
            </v-btn>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>
<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { onMounted, ref } from 'vue'
import GroupPricingSerivce from '@/renderer/api/GroupPricingService'

defineProps<{ loading?: boolean }>()

interface Benefit {
  id: number
  benefit_name: string
  benefit_alias: string
  benefit_code: string
  benefit_alias_code: string
  is_mapped: boolean
}

const emit = defineEmits<{
  (e: 'submit', value: Benefit[]): void
}>()

const benefits = ref<Benefit[]>([])

onMounted(() => {
  GroupPricingSerivce.getBenefitMaps().then((res) => {
    if (res.data.length > 0) {
      benefits.value = res.data.map((benefit: any) => ({
        id: benefit.id,
        benefit_name: benefit.benefit_name,
        benefit_alias: benefit.benefit_alias,
        benefit_code: benefit.benefit_code,
        benefit_alias_code: benefit.benefit_alias_code,
        is_mapped: benefit.is_mapped
      }))
    }
  })
})

function submitBenefits() {
  emit('submit', benefits.value)
}

const isValid = ref(true)

const rules = {
  alias: (val: string) => {
    if (!val) return true
    const exists = benefits.value.filter(
      (benefit) =>
        benefit.benefit_alias?.trim().toLowerCase() ===
        val?.trim().toLowerCase()
    ).length
    if (exists === 1) {
      isValid.value = true
      return true
    } else {
      isValid.value = false
      return 'This alias is already in used'
    }
  },
  alias_code: (val: string) => {
    if (!val) return true
    const exists = benefits.value.filter(
      (benefit) =>
        benefit.benefit_alias_code?.trim().toLowerCase() ===
        val?.trim().toLowerCase()
    ).length
    if (exists === 1) {
      isValid.value = true
      return true
    } else {
      isValid.value = false
      return 'This alias code is already in used'
    }
  }
}
</script>
<style lang="css" scoped></style>
