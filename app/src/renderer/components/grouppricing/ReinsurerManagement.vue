<template>
  <base-card :show-actions="false">
    <template #header>
      <v-icon icon="mdi-office-building" class="mr-3" />
      <span>Reinsurer Management</span>
    </template>
    <template #default>
      <h5 class="section-title">
        <v-icon icon="mdi-plus-circle-outline" size="small" class="mr-2" />
        Add New Reinsurer
      </h5>
      <v-row>
        <v-col cols="12" md="3">
          <v-text-field
            v-model="reinsurerData.name"
            placeholder="Enter reinsurer name"
            label="Reinsurer Name"
            variant="outlined"
            density="compact"
            prepend-inner-icon="mdi-office-building"
          />
        </v-col>
        <v-col cols="12" md="2">
          <v-text-field
            v-model="reinsurerData.code"
            placeholder="Enter code"
            label="Reinsurer Code"
            variant="outlined"
            density="compact"
            prepend-inner-icon="mdi-identifier"
          />
        </v-col>
        <v-col cols="12" md="3">
          <v-text-field
            v-model="reinsurerData.contact_email"
            placeholder="Enter contact email"
            label="Contact Email"
            variant="outlined"
            density="compact"
            prepend-inner-icon="mdi-email"
            type="email"
          />
        </v-col>
        <v-col cols="12" md="3">
          <v-text-field
            v-model="reinsurerData.contact_person"
            placeholder="Enter contact person"
            label="Contact Person"
            variant="outlined"
            density="compact"
            prepend-inner-icon="mdi-account"
          />
        </v-col>
        <v-col cols="12" md="1">
          <v-btn
            color="primary"
            size="small"
            icon="mdi-plus"
            :loading="isSaving"
            elevation="2"
            @click="handleAddReinsurer"
          />
        </v-col>
      </v-row>
      <h5 class="section-title">
        <v-icon icon="mdi-format-list-bulleted" size="small" class="mr-2" />
        Registered Reinsurers
      </h5>
      <v-row v-if="reinsurers.length > 0">
        <v-col cols="12">
          <data-grid
            :row-data="rowData"
            :column-defs="columnDefs"
            @row-clicked="handleRowSelection"
          />
        </v-col>
      </v-row>
      <div v-else class="empty-state">
        <v-icon
          icon="mdi-office-building-plus"
          size="64"
          color="grey-lighten-2"
        />
        <h3 class="text-h6 text-medium-emphasis mt-4">No Reinsurers Added</h3>
        <p class="text-body-2 text-medium-emphasis">
          Add your first reinsurer using the form above to get started
        </p>
      </div>
    </template>
  </base-card>
  <v-dialog v-model="editDialog" max-width="800px" persistent>
    <v-card class="edit-dialog">
      <v-card-title class="dialog-header">
        <div class="d-flex justify-space-between align-center w-100">
          <div class="d-flex align-center">
            <v-icon icon="mdi-office-building-edit" class="mr-3" />
            <span>Edit Reinsurer - {{ selectedReinsurer?.name }}</span>
          </div>
          <div class="d-flex gap-2">
            <v-btn
              color="warning"
              variant="outlined"
              prepend-icon="mdi-cancel"
              :loading="isDeactivating"
              size="small"
              @click="openDeactivateDialog"
            >
              Deactivate
            </v-btn>
            <v-btn
              color="error"
              variant="outlined"
              prepend-icon="mdi-delete"
              :loading="isDeleting"
              size="small"
              @click="handleDeleteReinsurer"
            >
              Delete
            </v-btn>
          </div>
        </div>
      </v-card-title>

      <v-card-text class="pa-6">
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="editReinsurerData.name"
              label="Reinsurer Name"
              variant="outlined"
              prepend-inner-icon="mdi-office-building"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="editReinsurerData.code"
              label="Reinsurer Code"
              variant="outlined"
              prepend-inner-icon="mdi-identifier"
              disabled
              hint="Code cannot be changed"
              persistent-hint
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="editReinsurerData.contact_email"
              label="Contact Email"
              variant="outlined"
              prepend-inner-icon="mdi-email"
              type="email"
            />
          </v-col>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="editReinsurerData.contact_person"
              label="Contact Person"
              variant="outlined"
              prepend-inner-icon="mdi-account"
            />
          </v-col>
        </v-row>
      </v-card-text>

      <v-card-actions class="pa-6">
        <v-spacer />
        <v-btn
          color="grey"
          variant="outlined"
          prepend-icon="mdi-close"
          @click="closeEditDialog"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          prepend-icon="mdi-content-save"
          :loading="isUpdating"
          @click="handleUpdateReinsurer"
        >
          Update Reinsurer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Deactivate Dialog -->
  <v-dialog v-model="deactivateDialog" max-width="500px" persistent>
    <v-card>
      <v-card-title class="dialog-header">
        <v-icon icon="mdi-cancel" class="mr-3" />
        <span>Deactivate Reinsurer - {{ selectedReinsurer?.name }}</span>
      </v-card-title>

      <v-card-text class="pa-6">
        <v-alert type="warning" variant="tonal" class="mb-4">
          Deactivating this reinsurer will remove it from all selection lists.
          This action can only be performed if the reinsurer is not involved in
          any active treaties.
        </v-alert>
        <v-textarea
          v-model="deactivationReason"
          label="Reason for Deactivation *"
          variant="outlined"
          rows="3"
          placeholder="Please provide a reason for deactivating this reinsurer"
          :rules="[(v) => !!v || 'Reason is required']"
        />
      </v-card-text>

      <v-card-actions class="pa-6">
        <v-spacer />
        <v-btn
          color="grey"
          variant="outlined"
          prepend-icon="mdi-close"
          @click="closeDeactivateDialog"
        >
          Cancel
        </v-btn>
        <v-btn
          color="warning"
          prepend-icon="mdi-cancel"
          :loading="isDeactivating"
          :disabled="!deactivationReason"
          @click="handleDeactivateReinsurer"
        >
          Deactivate Reinsurer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <confirmation-dialog ref="confirmDialog" />
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import ConfirmationDialog from '@/renderer/components/ConfirmDialog.vue'
import DataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useNotifications } from '@/renderer/composables/useNotifications'
import { useErrorHandler } from '@/renderer/composables/useErrorHandler'
import BaseCard from '../BaseCard.vue'

const emit = defineEmits<{
  reinsurerAdded: []
  reinsurerUpdated: []
  reinsurerDeleted: []
}>()

const confirmDialog = ref()
const editDialog = ref(false)
const deactivateDialog = ref(false)

// State
const reinsurers = ref<any[]>([])
const selectedReinsurer = ref<any>(null)
const isSaving = ref(false)
const isDeleting = ref(false)
const isUpdating = ref(false)
const isDeactivating = ref(false)
const deactivationReason = ref('')

const reinsurerData = ref({
  name: '',
  code: '',
  contact_email: '',
  contact_person: ''
})

const editReinsurerData = ref({
  id: 0,
  name: '',
  code: '',
  contact_email: '',
  contact_person: ''
})

const { showSuccess, showError } = useNotifications()
const { handleApiError } = useErrorHandler()

// Grid configuration
const columnDefs = ref([
  {
    headerName: 'Reinsurer Name',
    field: 'name',
    filter: true,
    sortable: true,
    flex: 2
  },
  {
    headerName: 'Code',
    field: 'code',
    filter: true,
    sortable: true,
    flex: 1
  },
  {
    headerName: 'Contact Email',
    field: 'contact_email',
    filter: true,
    sortable: true,
    flex: 2
  },
  {
    headerName: 'Contact Person',
    field: 'contact_person',
    filter: true,
    sortable: true,
    flex: 2
  },
  {
    headerName: 'Created By',
    field: 'created_by',
    filter: true,
    sortable: true,
    flex: 1
  }
])

const rowData = computed(() => reinsurers.value)

// Methods
const loadReinsurers = async () => {
  try {
    const response = await GroupPricingService.getReinsurers()
    reinsurers.value = response.data || []
  } catch (error) {
    handleApiError(error, 'Failed to load reinsurers')
  }
}

const handleAddReinsurer = async () => {
  if (!reinsurerData.value.name || !reinsurerData.value.code) {
    showError('Please fill in reinsurer name and code')
    return
  }

  isSaving.value = true
  try {
    await GroupPricingService.createReinsurer(reinsurerData.value)
    showSuccess('Reinsurer added successfully')
    resetForm()
    await loadReinsurers()
    emit('reinsurerAdded')
  } catch (error: any) {
    if (error.response?.status === 409) {
      showError('Reinsurer with this code already exists')
    } else {
      handleApiError(error, 'Failed to add reinsurer')
    }
  } finally {
    isSaving.value = false
  }
}

const handleRowSelection = (event: any) => {
  if (event.data) {
    selectedReinsurer.value = event.data
    editReinsurerData.value = { ...event.data }
    editDialog.value = true
  }
}

const handleUpdateReinsurer = async () => {
  if (!editReinsurerData.value.name) {
    showError('Please fill in reinsurer name')
    return
  }

  isUpdating.value = true
  try {
    await GroupPricingService.updateReinsurer(
      editReinsurerData.value.id,
      editReinsurerData.value
    )
    showSuccess('Reinsurer updated successfully')
    closeEditDialog()
    await loadReinsurers()
    emit('reinsurerUpdated')
  } catch (error) {
    handleApiError(error, 'Failed to update reinsurer')
  } finally {
    isUpdating.value = false
  }
}

const handleDeleteReinsurer = async () => {
  const confirmed = await confirmDialog.value?.open(
    'Delete Reinsurer',
    `Are you sure you want to delete ${selectedReinsurer.value?.name}? This action cannot be undone.`,
    'error'
  )

  if (!confirmed) return

  isDeleting.value = true
  try {
    await GroupPricingService.deleteReinsurer(selectedReinsurer.value.id)
    showSuccess('Reinsurer deleted successfully')
    closeEditDialog()
    await loadReinsurers()
    emit('reinsurerDeleted')
  } catch (error) {
    handleApiError(error, 'Failed to delete reinsurer')
  } finally {
    isDeleting.value = false
  }
}

const resetForm = () => {
  reinsurerData.value = {
    name: '',
    code: '',
    contact_email: '',
    contact_person: ''
  }
}

const closeEditDialog = () => {
  editDialog.value = false
  selectedReinsurer.value = null
  editReinsurerData.value = {
    id: 0,
    name: '',
    code: '',
    contact_email: '',
    contact_person: ''
  }
}

const openDeactivateDialog = () => {
  deactivationReason.value = ''
  deactivateDialog.value = true
}

const closeDeactivateDialog = () => {
  deactivateDialog.value = false
  deactivationReason.value = ''
}

const handleDeactivateReinsurer = async () => {
  if (!deactivationReason.value.trim()) {
    showError('Please provide a reason for deactivation')
    return
  }

  isDeactivating.value = true
  try {
    await GroupPricingService.deactivateReinsurer(
      selectedReinsurer.value.id,
      deactivationReason.value
    )
    showSuccess('Reinsurer deactivated successfully')
    closeDeactivateDialog()
    closeEditDialog()
    await loadReinsurers()
    emit('reinsurerUpdated')
  } catch (error: any) {
    if (error.response?.status === 409) {
      // Conflict - reinsurer has active treaties
      showError(
        error.response?.data?.error ||
          'Cannot deactivate: reinsurer is involved in active treaties'
      )
    } else {
      handleApiError(error, 'Failed to deactivate reinsurer')
    }
  } finally {
    isDeactivating.value = false
  }
}

onMounted(() => {
  loadReinsurers()
})
</script>

<style scoped>
.section-title {
  font-size: 1rem;
  font-weight: 600;
  color: #424242;
  margin: 1.5rem 0 1rem;
  display: flex;
  align-items: center;
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.edit-dialog {
  border-radius: 12px;
}

.dialog-header {
  background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
  color: white;
  padding: 1.5rem;
}
</style>
