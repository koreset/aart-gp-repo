<template>
  <base-card :show-actions="false">
    <template #header>
      <v-icon icon="mdi-account-group" class="mr-3" />
      <span>Broker Management</span>
    </template>
    <template #default>
      <h5 class="section-title">
        <v-icon icon="mdi-account-plus-outline" size="small" class="mr-2" />
        Add New Broker
      </h5>
      <v-form ref="addForm">
        <v-row>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="brokerData.name"
              placeholder="Enter broker name"
              label="Broker Name"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-account"
              :rules="[(v) => !!v?.trim() || 'Broker name is required']"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-text-field
              v-model="brokerData.contact_email"
              placeholder="Enter contact email"
              label="Contact Email"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-email"
              type="email"
              :rules="[
                (v) => !!v?.trim() || 'Contact email is required',
                (v) => /.+@.+\..+/.test(v) || 'Must be a valid email address'
              ]"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-text-field
              v-model="brokerData.contact_number"
              placeholder="Enter contact number"
              label="Contact Number"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-phone"
              :rules="[(v) => !!v?.trim() || 'Contact number is required']"
            />
          </v-col>
          <v-col cols="12" md="1">
            <v-btn
              color="primary"
              size="small"
              icon="mdi-plus"
              :loading="isSaving"
              elevation="2"
              @click="handleAddBroker"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" md="3">
            <v-text-field
              v-model="brokerData.fsp_number"
              placeholder="FSP license number"
              label="FSP Number"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-certificate"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-select
              v-model="brokerData.fsp_category"
              placeholder="Select FSP category"
              label="FSP Category"
              variant="outlined"
              density="compact"
              :items="fspCategories"
              clearable
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-text-field
              v-model="brokerData.binder_agreement_ref"
              placeholder="Binder agreement reference"
              label="Binder Agreement Ref"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-file-document"
            />
          </v-col>
          <v-col cols="12" md="3">
            <v-text-field
              v-model="brokerData.tied_agent_ref"
              placeholder="Tied agent reference"
              label="Tied Agent Ref"
              variant="outlined"
              density="compact"
              prepend-inner-icon="mdi-link-variant"
            />
          </v-col>
        </v-row>
      </v-form>
      <h5 class="section-title">
        <v-icon icon="mdi-format-list-bulleted" size="small" class="mr-2" />
        Registered Brokers
      </h5>
      <v-row v-if="brokers.length > 0">
        <v-col cols="12">
          <data-grid
            :row-data="rowData"
            :column-defs="columnDefs"
            @row-clicked="handleRowSelection"
          />
        </v-col>
      </v-row>
      <div v-else class="empty-state">
        <v-icon icon="mdi-account-plus" size="64" color="grey-lighten-2" />
        <h3 class="text-h6 text-medium-emphasis mt-4">No Brokers Added</h3>
        <p class="text-body-2 text-medium-emphasis">
          Add your first broker using the form above to get started
        </p>
      </div>
    </template>
  </base-card>
  <v-dialog v-model="editDialog" max-width="700px" persistent>
    <v-card class="edit-dialog">
      <v-card-title class="dialog-header">
        <div class="d-flex justify-space-between align-center w-100">
          <div class="d-flex align-center">
            <v-icon icon="mdi-account-edit" class="mr-3" />
            <span>Edit Broker - {{ selectedBroker?.name }}</span>
          </div>
          <v-btn
            color="error"
            variant="outlined"
            prepend-icon="mdi-delete"
            :loading="isDeleting"
            size="small"
            @click="handleDeleteBroker"
          >
            Delete
          </v-btn>
        </div>
      </v-card-title>

      <v-card-text class="pa-6">
        <v-form ref="editForm">
          <v-row>
            <v-col cols="12">
              <v-text-field
                v-model="editBrokerData.name"
                label="Broker Name"
                variant="outlined"
                prepend-inner-icon="mdi-account"
                :rules="[(v) => !!v?.trim() || 'Broker name is required']"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="editBrokerData.contact_email"
                label="Contact Email"
                variant="outlined"
                prepend-inner-icon="mdi-email"
                type="email"
                :rules="[
                  (v) => !!v?.trim() || 'Contact email is required',
                  (v) => /.+@.+\..+/.test(v) || 'Must be a valid email address'
                ]"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="editBrokerData.contact_number"
                label="Contact Number"
                variant="outlined"
                prepend-inner-icon="mdi-phone"
                :rules="[(v) => !!v?.trim() || 'Contact number is required']"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="editBrokerData.fsp_number"
                label="FSP Number"
                variant="outlined"
                prepend-inner-icon="mdi-certificate"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="editBrokerData.fsp_category"
                label="FSP Category"
                variant="outlined"
                :items="fspCategories"
                clearable
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="editBrokerData.binder_agreement_ref"
                label="Binder Agreement Ref"
                variant="outlined"
                prepend-inner-icon="mdi-file-document"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="editBrokerData.tied_agent_ref"
                label="Tied Agent Ref"
                variant="outlined"
                prepend-inner-icon="mdi-link-variant"
              />
            </v-col>
          </v-row>
        </v-form>
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
          @click="handleUpdateBroker"
        >
          Update Broker
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <confirmation-dialog ref="confirmDialog" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ConfirmationDialog from '@/renderer/components/ConfirmDialog.vue'
import DataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'
import { useBrokers } from '@/renderer/composables/useBrokers'
import { useBrokerGrid } from '@/renderer/composables/useBrokerGrid'
import { useNotifications } from '@/renderer/composables/useNotifications'
import { useErrorHandler } from '@/renderer/composables/useErrorHandler'
import BaseCard from '../BaseCard.vue'

const emit = defineEmits<{
  brokerAdded: []
  brokerUpdated: []
  brokerDeleted: []
}>()

const confirmDialog = ref()
const addForm = ref()
const editForm = ref()
const editDialog = ref(false)
const fspCategories = [
  'Category I',
  'Category II',
  'Category IIA',
  'Category III',
  'Category IV'
]

const {
  brokerData,
  editBrokerData,
  isSaving,
  isDeleting,
  isUpdating,
  loadBrokers,
  createBroker,
  updateBroker,
  deleteBroker,
  resetBrokerForm,
  setEditBrokerData
} = useBrokers()

const {
  brokers,
  rowData,
  columnDefs,
  selectedBroker,
  handleRowSelection: gridHandleRowSelection,
  setBrokers,
  clearSelection
} = useBrokerGrid()

const { showSuccess, showError } = useNotifications()
const { handleApiError } = useErrorHandler()

const handleRowSelection = (event: any) => {
  gridHandleRowSelection(event)
  if (selectedBroker.value) {
    setEditBrokerData(selectedBroker.value)
    editDialog.value = true
  }
}

const handleAddBroker = async () => {
  const { valid } = await addForm.value.validate()
  if (!valid) return

  try {
    await createBroker(brokerData.value)
    await refreshBrokers()
    resetBrokerForm()
    addForm.value.resetValidation()
    showSuccess('Broker added successfully')
    emit('brokerAdded')
  } catch (error: any) {
    const errorMessage = handleApiError(error)
    showError(errorMessage)
  }
}

const handleUpdateBroker = async () => {
  const { valid } = await editForm.value.validate()
  if (!valid) return

  try {
    if (!selectedBroker.value?.id) return

    await updateBroker(selectedBroker.value.id, editBrokerData.value)
    await refreshBrokers()
    closeEditDialog()
    showSuccess('Broker updated successfully')
    emit('brokerUpdated')
  } catch (error: any) {
    const errorMessage = handleApiError(error)
    showError(errorMessage)
  }
}

const handleDeleteBroker = async () => {
  if (!selectedBroker.value?.id) return

  const confirmed = await confirmDialog.value.open(
    `Deleting entry for broker ${selectedBroker.value.id}: ${editBrokerData.value.name}`,
    'This operation is irreversible and can only be performed on brokers with no active quotes.'
  )

  if (!confirmed) return

  try {
    await deleteBroker(selectedBroker.value.id)
    await refreshBrokers()
    closeEditDialog()
    showSuccess('Broker deleted successfully')
    emit('brokerDeleted')
  } catch (error: any) {
    const errorMessage = handleApiError(error)
    showError(errorMessage)
  }
}

const closeEditDialog = () => {
  editDialog.value = false
  clearSelection()
  editForm.value?.resetValidation()
}

const refreshBrokers = async () => {
  try {
    const data = await loadBrokers()
    setBrokers(data)
  } catch (error) {
    // Error handled by composable notification system
  }
}

onMounted(() => {
  refreshBrokers()
})
</script>

<style scoped>
.broker-card {
  border-radius: 16px;
  overflow: hidden;
}

.section-header {
  background: linear-gradient(135deg, #4caf50 0%, #388e3c 100%);
  color: white;
  padding: 1.5rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.form-section {
  margin-bottom: 2rem;
  padding: 1.5rem;
  background: #f8f9fa;
  border-radius: 12px;
  border-left: 4px solid #4caf50;
}

.section-title {
  color: #4caf50;
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
}

.data-grid-container {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
}

.edit-dialog {
  border-radius: 16px;
}

.dialog-header {
  background: linear-gradient(135deg, #2196f3 0%, #1976d2 100%);
  color: white;
  padding: 1.5rem;
}

/* Form field enhancements */
:deep(.v-text-field .v-field__outline) {
  border-radius: 8px;
}

:deep(.v-btn) {
  border-radius: 8px;
  text-transform: none;
  font-weight: 500;
}

/* Hover effects for buttons */
:deep(.v-btn:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* Icon button special styling */
:deep(.v-btn--icon) {
  border-radius: 50%;
}
</style>
