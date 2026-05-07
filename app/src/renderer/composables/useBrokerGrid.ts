import { ref, computed } from 'vue'
import type { BrokerData, ColumnDef } from '@/renderer/types/metadata'
import formatValues from '@/renderer/utils/format_values'

const formatCreationDate = (params: { value?: string | null }): string => {
  const value = params.value
  if (!value) return ''
  const d = new Date(value)
  if (isNaN(d.getTime())) return String(value)
  return d.toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

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
    },
    {
      headerName: 'Created By',
      field: 'created_by',
      valueFormatter: formatValues,
      minWidth: 150,
      sortable: true,
      filter: true,
      resizable: true
    },
    {
      headerName: 'Creation Date',
      field: 'creation_date',
      valueFormatter: formatCreationDate,
      minWidth: 180,
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
