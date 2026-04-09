import { ref } from 'vue'
import type { BrokerData, CreateBrokerPayload } from '@/renderer/types/metadata'
import GroupPricingService from '@/renderer/api/GroupPricingService'

export function useBrokers() {
  const isLoading = ref(false)
  const isSaving = ref(false)
  const isDeleting = ref(false)
  const isUpdating = ref(false)

  const brokerData = ref<CreateBrokerPayload>({
    name: '',
    contact_email: '',
    contact_number: '',
    fsp_number: '',
    fsp_category: '',
    binder_agreement_ref: '',
    tied_agent_ref: ''
  })

  const editBrokerData = ref<BrokerData>({
    name: '',
    contact_email: '',
    contact_number: '',
    fsp_number: '',
    fsp_category: '',
    binder_agreement_ref: '',
    tied_agent_ref: ''
  })

  const loadBrokers = async (): Promise<BrokerData[]> => {
    try {
      isLoading.value = true
      const res = await GroupPricingService.getBrokers()
      return res.data || []
    } finally {
      isLoading.value = false
    }
  }

  const createBroker = async (payload: CreateBrokerPayload): Promise<any> => {
    try {
      isSaving.value = true
      const res = await GroupPricingService.createBroker(payload)
      return res
    } finally {
      isSaving.value = false
    }
  }

  const updateBroker = async (id: number, payload: Partial<BrokerData>) => {
    try {
      isUpdating.value = true
      const res = await GroupPricingService.updateBroker(id, payload)
      return res
    } finally {
      isUpdating.value = false
    }
  }

  const deleteBroker = async (id: number) => {
    try {
      isDeleting.value = true
      const res = await GroupPricingService.deleteBroker(id)
      return res
    } finally {
      isDeleting.value = false
    }
  }

  const resetBrokerForm = () => {
    brokerData.value = {
      name: '',
      contact_email: '',
      contact_number: '',
      fsp_number: '',
      fsp_category: '',
      binder_agreement_ref: '',
      tied_agent_ref: ''
    }
  }

  const setEditBrokerData = (broker: BrokerData) => {
    editBrokerData.value = { ...broker }
  }

  return {
    brokerData,
    editBrokerData,
    isLoading,
    isSaving,
    isDeleting,
    isUpdating,
    loadBrokers,
    createBroker,
    updateBroker,
    deleteBroker,
    resetBrokerForm,
    setEditBrokerData
  }
}
