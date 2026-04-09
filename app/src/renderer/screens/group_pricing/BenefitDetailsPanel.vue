<template>
  <v-expansion-panel elevation="1">
    <v-expansion-panel-title class="font-weight-bold text-primary">
      {{ benefitData.title }}
    </v-expansion-panel-title>
    <v-expansion-panel-text>
      <v-container fluid>
        <template v-if="benefitData.type === 'sum_assured'">
          <div class="section-header">SUM ASSURED</div>
          <v-data-table
            :headers="tableHeaders"
            :items="sumAssuredItems"
            item-key="name"
            class="elevation-1 data-table"
            density="compact"
            :hide-default-footer="true"
          >
          </v-data-table>
        </template>

        <div class="section-header mt-6">RISK PREMIUM</div>
        <v-data-table
          :headers="tableHeaders"
          :items="riskPremiumItems"
          item-key="name"
          class="elevation-1 data-table"
          density="compact"
          :hide-default-footer="true"
        >
        </v-data-table>

        <div class="section-header mt-6">OFFICE PREMIUM</div>
        <v-data-table
          :headers="tableHeaders"
          :items="officePremiumItems"
          item-key="name"
          class="elevation-1 data-table"
          density="compact"
          :hide-default-footer="true"
        >
        </v-data-table>
      </v-container>
    </v-expansion-panel-text>
  </v-expansion-panel>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  benefitData: {
    type: Object,
    required: true
  },
  resultSummary: {
    type: Object,
    required: true
  }
})

const tableHeaders: any = [
  { title: 'Metric', key: 'name', sortable: false, width: '40%' },
  {
    title: 'Theoretical Rate',
    key: 'theoretical',
    align: 'end',
    sortable: false
  },
  {
    title: 'Experience Rated',
    key: 'experience',
    align: 'end',
    sortable: false
  },
  { title: 'Discounted', key: 'discounted', align: 'end', sortable: false }
]

const formatValue = (value, isPercent = false) => {
  if (value === null || value === undefined || isNaN(value)) return '-'
  const num = Math.ceil(Number(value) * 100) / 100
  const formattedNum = num
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ')
  return isPercent ? `${formattedNum}%` : formattedNum
}

// Computed properties to generate table rows from props
const sumAssuredItems = computed(() => {
  const data = props.benefitData.sumAssuredData
  if (!data) return []
  return [
    {
      name: 'Minimum Sum Assured',
      theoretical: formatValue(props.resultSummary[data.min]),
      experience: formatValue(props.resultSummary[data.min]),
      discounted: formatValue(props.resultSummary[data.min])
    },
    {
      name: 'Maximum Sum Assured',
      theoretical: formatValue(props.resultSummary[data.max]),
      experience: formatValue(props.resultSummary[data.max]),
      discounted: formatValue(props.resultSummary[data.max])
    },
    {
      name: 'Maximum FCL Capped Sum Assured',
      theoretical: formatValue(props.resultSummary[data.max_capped]),
      experience: formatValue(props.resultSummary[data.max_capped]),
      discounted: formatValue(props.resultSummary[data.max_capped])
    },
    {
      name: 'Total Capped Sum Assured',
      theoretical: formatValue(props.resultSummary[data.total_capped]),
      experience: formatValue(props.resultSummary[data.total_capped]),
      discounted: formatValue(props.resultSummary[data.total_capped])
    },
    {
      name: 'Average Covered Sum Assured',
      theoretical: formatValue(props.resultSummary[data.average_capped]),
      experience: formatValue(props.resultSummary[data.average_capped]),
      discounted: formatValue(props.resultSummary[data.average_capped])
    }
  ]
})

const riskPremiumItems = computed(() => {
  const data = props.benefitData.riskPremiumData
  if (!data) return []

  const items = [
    {
      name: 'Annual Risk Premium',
      theoretical: formatValue(props.resultSummary[data.annual_premium]),
      experience: formatValue(props.resultSummary[data.exp_annual_premium]),
      discounted: formatValue(props.resultSummary[data.exp_annual_premium])
    },
    {
      name: 'Risk Premium as % of Annual Salary',
      theoretical: formatValue(
        props.resultSummary[data.salary_proportion] * 100,
        true
      ),
      experience: formatValue(
        props.resultSummary[data.exp_salary_proportion] * 100,
        true
      ),
      discounted: formatValue(
        props.resultSummary[data.exp_salary_proportion] * 100,
        true
      )
    }
  ]

  if (props.benefitData.type === 'family_funeral') {
    items.splice(1, 0, {
      name: 'Risk Premium per Member per Month',
      theoretical: formatValue(props.resultSummary[data.annual_premium] / 12),
      experience: formatValue(
        props.resultSummary[data.exp_annual_premium] / 12
      ),
      discounted: formatValue(props.resultSummary[data.exp_annual_premium] / 12)
    })
  } else {
    items.unshift({
      name: 'Expected Number of Claims',
      theoretical: formatValue(props.resultSummary[data.claims]),
      experience: formatValue(props.resultSummary[data.exp_claims]),
      discounted: formatValue(props.resultSummary[data.exp_claims])
    })
    items.splice(2, 0, {
      name: 'Unit Rate per 1000 Sum Assured',
      theoretical: formatValue(props.resultSummary[data.unit_rate]),
      experience: formatValue(props.resultSummary[data.exp_unit_rate]),
      discounted: formatValue(props.resultSummary[data.exp_unit_rate])
    })
  }

  return items
})

const officePremiumItems = computed(() => {
  const data = props.benefitData.officePremiumData
  if (!data) return []

  const items = [
    {
      name: 'Annual Office Premium',
      theoretical: formatValue(props.resultSummary[data.annual_premium]),
      experience: formatValue(props.resultSummary[data.exp_annual_premium]),
      discounted: formatValue(props.resultSummary[data.exp_annual_premium])
    },
    {
      name: 'Office Premium as % of Annual Salary',
      theoretical: formatValue(
        props.resultSummary[data.salary_proportion] * 100,
        true
      ),
      experience: formatValue(
        props.resultSummary[data.exp_salary_proportion] * 100,
        true
      ),
      discounted: formatValue(
        props.resultSummary[data.exp_salary_proportion] * 100,
        true
      )
    }
  ]

  if (props.benefitData.type === 'family_funeral') {
    items.splice(1, 0, {
      name: 'Office Premium per Member per Month',
      theoretical: formatValue(props.resultSummary[data.annual_premium] / 12),
      experience: formatValue(
        props.resultSummary[data.exp_annual_premium] / 12
      ),
      discounted: formatValue(props.resultSummary[data.exp_annual_premium] / 12)
    })
  } else {
    items.splice(1, 0, {
      name: 'Unit Office Premium Rate per 1000 Covered Sum Assured',
      theoretical: formatValue(props.resultSummary[data.unit_rate]),
      experience: formatValue(props.resultSummary[data.exp_unit_rate]),
      discounted: formatValue(props.resultSummary[data.exp_unit_rate])
    })
  }

  return items
})
</script>

<style scoped>
.section-header {
  font-size: 1rem;
  font-weight: 700;
  color: #333;
  padding: 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
  margin-bottom: 16px;
}

.data-table {
  border: 1px solid #e0e0e0;
}

:deep(.v-data-table-header) {
  background-color: #fafafa;
}

:deep(tbody tr:hover) {
  background-color: #f5f5f5 !important;
}
</style>
