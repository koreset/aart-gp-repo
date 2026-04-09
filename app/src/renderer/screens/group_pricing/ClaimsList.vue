<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Claims</span>
          </template>
          <template #default>
            <v-row
              ><v-col
                ><v-btn rounded color="primary" size="small" @click="lodgeClaim"
                  >Lodge Claim</v-btn
                >
              </v-col></v-row
            >
            <v-row>
              <v-col>
                <data-grid
                  :columnDefs="columnDefs"
                  :rowData="rowData"
                  :pagination="true"
                />
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
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import formatValues from '@/renderer/utils/format_values'
const router = useRouter()
const rowData: any = ref([])
const columnDefs: any = ref([])

const lodgeClaim = () => {
  router.push({ name: 'group-pricing-lodge-claim' })
}

onMounted(() => {
  GroupPricingService.getClaims().then((response) => {
    rowData.value = response.data
    createColumnDefs(rowData.value)
  })
})

const createColumnDefs = (data: any) => {
  columnDefs.value = []
  Object.keys(data[0]).forEach((element) => {
    const header: any = {}
    header.headerName = element
    header.field = element
    header.valueFormatter = formatValues
    header.minWidth = 220
    header.filter = true
    header.resizable = true
    header.sortable = true
    columnDefs.value.push(header)
  })
}
</script>
<style lang="css" scoped></style>
