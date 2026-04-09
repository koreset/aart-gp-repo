<template>
  <v-container>
    <v-row>
      <v-col cols="6">
        <v-file-input
          v-model:model-value="
            groupStore.group_pricing_quote.uploadData.member_data_file
          "
          variant="outlined"
          density="compact"
          show-size
          accept=".csv"
          placeholder="Click to upload member data"
          label="Member Data"
          prepend-icon="mdi-paperclip"
        >
        </v-file-input>
      </v-col>
      <v-col cols="6">
        <v-file-input
          v-model:model-value="
            groupStore.group_pricing_quote.uploadData.claims_experience_file
          "
          variant="outlined"
          density="compact"
          show-size
          accept=".csv"
          placeholder="Click to upload claims experience data"
          label="Claims Experience"
          prepend-icon="mdi-paperclip"
        >
        </v-file-input>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="4">
        <v-select
          key="user"
          v-model="groupStore.group_pricing_quote.reviewer"
          label="Choose a reviewer in your organisation"
          variant="outlined"
          density="compact"
          :items="orgUsers"
          item-title="user"
          item-value="user"
        ></v-select>
      </v-col>
      <v-col cols="4">
        <v-select
          v-model="groupStore.group_pricing_quote.basis"
          label="Choose a parameter basis"
          variant="outlined"
          density="compact"
          :items="parameterBases"
        ></v-select>
      </v-col>
    </v-row>
  </v-container>
</template>
<script setup lang="ts">
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useGroupPricingStore } from '@/renderer/store/group_pricing'
import { useAppStore } from '@/renderer/store/app'
import { computed, onMounted, ref } from 'vue'

const appStore = useAppStore()
const groupStore = useGroupPricingStore()
const orgUsers: any = ref([])
const organization = computed(
  () => appStore.getLicenseData.data.attributes.metadata.organization
)
const parameterBases = ref([])

const validateForm = () => {
  return true
}

onMounted(() => {
  try {
    GroupPricingService.getOrgUsers({ name: organization.value }).then(
      (res) => {
        const uniqueData = Array.from(
          new Map(res.data.map((entry) => [entry.user, entry])).values()
        )
        orgUsers.value = uniqueData
      }
    )
    GroupPricingService.getParameterBases().then((res) => {
      parameterBases.value = res.data
    })
  } catch (error) {}
})

defineExpose({
  validateForm
})
</script>
<style lang="css" scoped></style>
