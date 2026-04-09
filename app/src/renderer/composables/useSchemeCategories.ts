import { ref, reactive } from 'vue'
import type {
  SchemeCategoryData,
  CreateSchemeCategoryPayload
} from '@/renderer/types/metadata'
import GroupPricingService from '@/renderer/api/GroupPricingService'

export function useSchemeCategories() {
  const categories = ref<SchemeCategoryData[]>([])
  const isSaving = ref(false)
  const isDeleting = ref(false)
  const isUpdating = ref(false)

  const categoryData = reactive<CreateSchemeCategoryPayload>({
    name: '',
    description: '',
    active: true
  })

  const editCategoryData = reactive<SchemeCategoryData>({
    id: undefined,
    name: '',
    description: '',
    active: true
  })

  const errors = ref<Record<string, string>>({})

  // API functions
  const loadSchemeCategories = async (): Promise<SchemeCategoryData[]> => {
    try {
      const response = await GroupPricingService.getSchemeCategoryMasters()
      return response.data
    } catch (error) {
      console.error('Failed to load scheme categories:', error)
      throw error
    }
  }

  const createSchemeCategory = async (
    payload: CreateSchemeCategoryPayload
  ): Promise<SchemeCategoryData> => {
    isSaving.value = true
    try {
      const response =
        await GroupPricingService.createSchemeCategoryMaster(payload)
      return response.data
    } catch (error) {
      console.error('Failed to create scheme category:', error)
      throw error
    } finally {
      isSaving.value = false
    }
  }

  const updateSchemeCategory = async (
    id: number,
    payload: Partial<SchemeCategoryData>
  ): Promise<SchemeCategoryData> => {
    isUpdating.value = true
    try {
      const response = await GroupPricingService.updateSchemeCategoryMaster(
        id,
        payload
      )
      return response.data
    } catch (error) {
      console.error('Failed to update scheme category:', error)
      throw error
    } finally {
      isUpdating.value = false
    }
  }

  const deleteSchemeCategory = async (id: number): Promise<void> => {
    isDeleting.value = true
    try {
      await GroupPricingService.deleteSchemeCategoryMaster(id)
    } catch (error) {
      console.error('Failed to delete scheme category:', error)
      throw error
    } finally {
      isDeleting.value = false
    }
  }

  const setCategories = (newCategories: SchemeCategoryData[]) => {
    categories.value = newCategories
  }

  const resetCategoryForm = () => {
    categoryData.name = ''
    categoryData.description = ''
    categoryData.active = true
    errors.value = {}
  }

  const setEditCategoryData = (category: SchemeCategoryData) => {
    editCategoryData.id = category.id
    editCategoryData.name = category.name
    editCategoryData.description = category.description || ''
    editCategoryData.active = category.active
  }

  return {
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
  }
}
