import { ref, computed } from 'vue'
import type { BrokerData, ColumnDef } from '@/renderer/types/metadata'
import formatValues from '@/renderer/utils/format_values'

export function useBrokerGrid() {
  const brokers = ref<BrokerData[]>([])
  const selectedBroker = ref<BrokerData | null>(null)

  const rowData = computed(() => brokers.value)

  const columnDefs = computed<ColumnDef[]>(() => [
    {
      headerName: 'ID',
      field: 'id',
      valueFormatter: formatValues,
      minWidth: 100,
      sortable: true,
      filter: true,
      resizable: true
    },
    {
      headerName: 'Name',
      field: 'name',
      valueFormatter: formatValues,
      minWidth: 200,
      sortable: true,
      filter: true,
      resizable: true
    },
    {
      headerName: 'Contact Email',
      field: 'contact_email',
      valueFormatter: formatValues,
      minWidth: 250,
      sortable: true,
      filter: true,
      resizable: true
    },
    {
      headerName: 'Contact Number',
      field: 'contact_number',
      valueFormatter: formatValues,
      minWidth: 200,
      sortable: true,
      filter: true,
      resizable: true
    },
    {
      headerName: 'FSP Number',
      field: 'fsp_number',
      valueFormatter: formatValues,
      minWidth: 150,
      sortable: true,
      filter: true,
      resizable: true
    }
  ])

  const handleRowSelection = (event: any) => {
    selectedBroker.value = event.data
  }

  const setBrokers = (newBrokers: BrokerData[]) => {
    brokers.value = newBrokers
  }

  const clearSelection = () => {
    selectedBroker.value = null
  }

  return {
    brokers,
    rowData,
    columnDefs,
    selectedBroker,
    handleRowSelection,
    setBrokers,
    clearSelection
  }
}
