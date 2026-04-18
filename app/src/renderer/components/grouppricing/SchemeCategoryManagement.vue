<template>
  <base-card :show-actions="false">
    <template #header>
      <v-icon icon="mdi-tag-multiple" class="mr-3" />
      <span>Scheme Categories</span>
    </template>
    <template #default>
      <p class="mb-6">
        Configure the available scheme categories for your organization. These
        categories will be used to classify group schemes.
      </p>
      <v-row class="mb-4">
        <v-col cols="12" md="4">
          <v-text-field
            v-model="categoryData.name"
            label="Category Name"
            variant="outlined"
            density="compact"
            :error-messages="errors.name"
            placeholder="e.g., Management, Executive"
          />
        </v-col>
        <v-col cols="12" md="5">
          <v-text-field
            v-model="categoryData.description"
            label="Description (Optional)"
            variant="outlined"
            density="compact"
            placeholder="Brief description of this category"
          />
        </v-col>
        <v-col cols="12" md="2">
          <v-switch
            v-model="categoryData.active"
            class="mt-0"
            label="Active"
            color="primary"
          />
        </v-col>
        <v-col cols="12" md="1">
          <v-btn
            color="primary"
            variant="flat"
            size="small"
            icon="mdi-plus"
            :loading="isSaving"
            elevation="2"
            @click="handleAddCategory"
          />
        </v-col>
      </v-row>
      <div v-if="categories.length > 0" class="form-section">
        <h4 class="text-h6 mb-4">Current Categories</h4>
        <div class="categories-grid">
          <v-card
            v-for="category in categories"
            :key="category.id"
            class="category-item mb-3"
            variant="outlined"
          >
            <v-card-text class="py-3">
              <div class="d-flex justify-space-between align-center">
                <div class="flex-grow-1">
                  <div class="d-flex align-center mb-2">
                    <h5 class="text-h6 mr-2">{{ category.name }}</h5>
                    <v-chip
                      :color="category.active ? 'success' : 'grey'"
                      size="small"
                      variant="flat"
                    >
                      {{ category.active ? 'Active' : 'Inactive' }}
                    </v-chip>
                  </div>
                  <p
                    v-if="category.description"
                    class="text-body-2 text-medium-emphasis ma-0"
                  >
                    {{ category.description }}
                  </p>
                </div>
                <div class="d-flex align-center">
                  <v-btn
                    icon="mdi-pencil"
                    variant="text"
                    size="small"
                    color="primary"
                    @click="handleEditCategory(category)"
                  />
                  <v-btn
                    icon="mdi-delete"
                    variant="text"
                    size="small"
                    color="error"
                    :loading="isDeleting"
                    @click="handleDeleteCategory(category)"
                  />
                </div>
              </div>
            </v-card-text>
          </v-card>
        </div>
      </div>
      <div v-else class="text-center py-8">
        <v-icon
          icon="mdi-tag-off"
          size="64"
          color="grey-lighten-1"
          class="mb-4"
        />
        <h4 class="text-h6 mb-2">No Categories Yet</h4>
        <p class="text-body-2 text-medium-emphasis">
          Add your first scheme category to get started
        </p>
      </div>
    </template>
  </base-card>
  <!-- Edit Category Dialog -->
  <v-dialog v-model="editDialog" max-width="500px">
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon icon="mdi-pencil" class="mr-2" />
        Edit Category
      </v-card-title>

      <v-card-text>
        <v-text-field
          v-model="editCategoryData.name"
          label="Category Name"
          variant="outlined"
          density="comfortable"
          class="mb-3"
        />
        <v-text-field
          v-model="editCategoryData.description"
          label="Description (Optional)"
          variant="outlined"
          density="comfortable"
          class="mb-3"
        />
        <v-switch
          v-model="editCategoryData.active"
          label="Active"
          color="primary"
          inset
        />
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="editDialog = false">Cancel</v-btn>
        <v-btn
          color="primary"
          variant="flat"
          :loading="isUpdating"
          @click="handleUpdateCategory"
        >
          Update
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Delete Confirmation Dialog -->
  <confirmation-dialog ref="confirmDeleteAction" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSchemeCategories } from '@/renderer/composables/useSchemeCategories'
import { useNotifications } from '@/renderer/composables/useNotifications'
import { useErrorHandler } from '@/renderer/composables/useErrorHandler'
import ConfirmationDialog from '@/renderer/components/ConfirmDialog.vue'
import type { SchemeCategoryData } from '@/renderer/types/metadata'
import BaseCard from '@/renderer/components/BaseCard.vue'
const {
  categories,
  categoryData,
  editCategoryData,
  isSaving,
  isDeleting,
  isUpdating,
  errors,
  loadSchemeCategories,
  createSchemeCategory,
  updateSchemeCategory,
  deleteSchemeCategory,
  setCategories,
  resetCategoryForm,
  setEditCategoryData
} = useSchemeCategories()

const { showSuccess, showError } = useNotifications()
const { handleValidationError } = useErrorHandler()

const editDialog = ref(false)
const deleteDialog = ref(false)
const categoryToDelete = ref<SchemeCategoryData | null>(null)
const confirmDeleteAction = ref()

const handleAddCategory = async () => {
  if (!categoryData.name.trim()) {
    errors.value.name = 'Category name is required'
    return
  }

  // Check for duplicates
  const exists = categories.value.some(
    (c) =>
      c.name.toLowerCase().trim() === categoryData.name.toLowerCase().trim()
  )
  if (exists) {
    errors.value.name = 'A category with this name already exists'
    return
  }

  try {
    const newCategory = await createSchemeCategory({
      name: categoryData.name.trim(),
      description: categoryData.description?.trim() || '',
      active: categoryData.active
    })

    categories.value.push(newCategory)
    resetCategoryForm()
    showSuccess('Category added successfully')
  } catch (error: any) {
    const message = handleValidationError(error)
    showError(message)
  }
}

const handleEditCategory = (category: SchemeCategoryData) => {
  setEditCategoryData(category)
  editDialog.value = true
}

const handleUpdateCategory = async () => {
  if (!editCategoryData.name.trim()) {
    showError('Category name is required')
    return
  }

  try {
    const updated = await updateSchemeCategory(editCategoryData.id!, {
      name: editCategoryData.name.trim(),
      description: editCategoryData.description?.trim() || '',
      active: editCategoryData.active
    })

    const index = categories.value.findIndex((c) => c.id === updated.id)
    if (index !== -1) {
      categories.value[index] = updated
    }

    editDialog.value = false
    showSuccess('Category updated successfully')
  } catch (error: any) {
    const message = handleValidationError(error)
    showError(message)
  }
}

const handleDeleteCategory = async (category: SchemeCategoryData) => {
  const confirmed = await confirmDeleteAction.value.open(
    `Are you sure you want to delete the category "${categoryToDelete.value?.name}"?`,
    'This action cannot be undone and may affect existing group schemes.'
  )
  if (!confirmed) return

  try {
    categoryToDelete.value = category
    await deleteSchemeCategory(categoryToDelete.value.id!)
    categories.value = categories.value.filter(
      (c) => c.id !== categoryToDelete.value!.id
    )
    showSuccess('Category deleted successfully')
  } catch (error: any) {
    const message = handleValidationError(error)
    showError(message)
  } finally {
    deleteDialog.value = false
    categoryToDelete.value = null
  }
}

const refreshCategories = async () => {
  try {
    const data = await loadSchemeCategories()
    setCategories(data)
  } catch (error) {
    // Error handled by composable notification system
  }
}

onMounted(() => {
  refreshCategories()
})
</script>

<style scoped>
.category-card {
  background: linear-gradient(145deg, #ffffff 0%, #f8f9fa 100%);
}

.form-section {
  margin-bottom: 1.5rem;
}

.categories-grid {
  max-height: 400px;
  overflow-y: auto;
}

.category-item {
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.category-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
