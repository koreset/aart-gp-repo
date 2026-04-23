<template>
  <v-container>
    <v-row>
      <v-col>
        <v-stepper
          v-model:model-value="position"
          class="smaller-font"
          dark
          alt-labels
        >
          <v-stepper-header>
            <template v-for="(step, index) in steps" :key="step.value">
              <v-stepper-item
                :title="step.title"
                :value="step.value"
              ></v-stepper-item>
              <v-divider
                v-if="(index as number) < steps.length - 1"
              ></v-divider>
            </template>
          </v-stepper-header>

          <v-stepper-window>
            <template v-for="step in steps" :key="step.value">
              <v-stepper-window-item :value="step.value">
                <component
                  :is="step.component"
                  ref="currentStep"
                  @fetch-quote-by-scheme="handleFetchQuoteByScheme"
                  @all_schemes_saved="areAllSchemesSaved = true"
                />
              </v-stepper-window-item>
            </template>

            <v-divider class="my-5"></v-divider>
            <v-row>
              <v-col class="d-flex">
                <v-btn
                  class="me-auto"
                  rounded
                  size="small"
                  color="primary"
                  @click="goToQuotes"
                  >Back to Quotes</v-btn
                >
                <v-btn
                  class="ml-2 mb-3"
                  size="small"
                  rounded
                  color="primary"
                  :disabled="isPrevDisabled"
                  @click="movePrev"
                  >Prev</v-btn
                >

                <v-btn
                  class="ml-4"
                  size="small"
                  rounded
                  color="primary"
                  :disabled="isNextDisabled"
                  @click="moveNext"
                  >Next</v-btn
                >

                <v-btn
                  v-if="position === steps.length"
                  class="ml-9 mb-3"
                  size="small"
                  rounded
                  color="primary"
                  :disabled="!areAllSchemesSaved || isGenerating"
                  :loading="isGenerating"
                  @click="generateQuote"
                  >Generate Quote</v-btn
                >
              </v-col>
            </v-row>
          </v-stepper-window>
        </v-stepper>
      </v-col>
    </v-row>
    <confirmation-dialog ref="confirmDialog" />
  </v-container>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeMount, shallowRef } from 'vue'
import { useRouter, useRoute, onBeforeRouteLeave } from 'vue-router'
import { useGroupPricingStore } from '@/renderer/store/group_pricing'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import Generalnput from '@/renderer/components/grouppricing/Generalnput.vue'
import AdditionalBenefits from '@/renderer/components/grouppricing/AdditionalBenefits.vue'
import ConfirmationDialog from '@/renderer/components/ConfirmDialog.vue'

const groupStore = useGroupPricingStore()
const router = useRouter()
const route = useRoute()
const quoteId = ref(route.params.id)
const benefitMaps: any = ref([])

const steps: any = shallowRef([])

onBeforeMount(async () => {
  const res = await GroupPricingService.getBenefitMaps()
  benefitMaps.value = res.data
  steps.value = [
    { title: 'General', value: 1, component: Generalnput },
    { title: 'Benefits Configuration', value: 2, component: AdditionalBenefits }
  ]
})

const position = ref(1)

const isPrevDisabled = computed(() => position.value <= 1)
const isNextDisabled = computed(() => position.value >= steps.value.length)
const areAllSchemesSaved = ref(false)
const isGenerating = ref(false)
const currentStep: any = ref(null)
const isDirty = ref(true) // starts dirty; set false after successful generation
const confirmDialog = ref()

onBeforeRouteLeave(async (_to, _from, next) => {
  if (!isDirty.value) return next()
  try {
    await confirmDialog.value.open(
      'Unsaved changes',
      'Your quote data will be lost if you leave. Continue?'
    )
    next()
  } catch {
    next(false)
  }
})

const moveNext = async () => {
  try {
    // Always treat currentStep.value as an array of component refs
    const stepRefs = Array.isArray(currentStep.value)
      ? currentStep.value
      : [currentStep.value]
    const stepInstance = stepRefs[position.value - 1]
    const isValid =
      stepInstance && stepInstance.validateForm
        ? await stepInstance.validateForm()
        : false
    if (!isValid) {
      return
    }
    position.value++
  } catch (e) {
    console.log(e)
  }
}

const movePrev = () => {
  position.value--
}

const goToQuotes = () => {
  groupStore.resetGroupPricingQuote()
  router.push({ name: 'group-pricing-quotes' })
}

const handleFetchQuoteByScheme = async (schemeName: string) => {
  try {
    const response = await GroupPricingService.getQuoteBySchemeName(schemeName)
    if (response && response.data) {
      // Update the store with the fetched quote data to prepopulate the form
      groupStore.group_pricing_quote = {
        ...groupStore.group_pricing_quote,
        ...response.data,
        // Reset commencement_date as it should be set for the new quote
        commencement_date: null,
        quote_type: 'Renewal',
        quote_id: 0,
        id: 0,
        edit_mode: false
      }
      console.log('Prepopulated quote with data:', response.data)
    }
  } catch (error) {
    console.error('Error fetching quote by scheme name:', error)
  }
}

const generateQuote = async () => {
  // Always treat currentStep.value as an array of component refs
  const stepRefs = Array.isArray(currentStep.value)
    ? currentStep.value
    : [currentStep.value]
  const lastStepInstance = stepRefs[steps.value.length - 1]
  const isValid =
    lastStepInstance && lastStepInstance.validateForm
      ? await lastStepInstance.validateForm()
      : false
  if (!isValid) return

  groupStore.group_pricing_quote.occupation_class = 0

  isGenerating.value = true
  try {
    await GroupPricingService.generateQuote(
      JSON.stringify(groupStore.group_pricing_quote)
    )
    groupStore.resetGroupPricingQuote()
    isDirty.value = false
    router.push({ name: 'group-pricing-quotes' })
  } catch (error: any) {
    console.error('Quote generation failed:', error)
  } finally {
    isGenerating.value = false
  }
}

onMounted(() => {
  GroupPricingService.getIndustries().then((res) => {
    groupStore.industries = res.data
  })
  GroupPricingService.getSchemesInforce().then((res) => {
    console.log('schemes', res.data)
    // filter out schemes whose status is not 'InForce'
    groupStore.groupSchemes = res.data.filter(
      (scheme: any) => scheme.status === 'in_force'
    )
  })
  if (quoteId.value) {
    GroupPricingService.getQuote(quoteId.value).then((res) => {
      groupStore.group_pricing_quote = res.data
      groupStore.group_pricing_quote.commencement_date = null
      groupStore.group_pricing_quote.edit_mode = true
      console.log('group_pricing_quote', groupStore.group_pricing_quote)
    })
  }
})
</script>
<style lang="css" scoped>
.smaller-font {
  font-size: 14px;
}
</style>
