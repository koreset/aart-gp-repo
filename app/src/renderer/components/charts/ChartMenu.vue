<template>
  <div class="chart-menu-wrapper">
    <v-menu location="bottom end" :close-on-content-click="true">
      <template #activator="{ props: menuProps }">
        <v-btn
          v-bind="menuProps"
          icon
          color="primary"
          variant="plain"
          size="small"
          title="Chart options"
        >
          <v-icon>mdi-menu</v-icon>
        </v-btn>
      </template>

      <v-list density="compact" min-width="210" nav>
        <v-list-subheader>Download</v-list-subheader>

        <v-list-item
          prepend-icon="mdi-image-outline"
          title="PNG image"
          @click="downloadPng"
        ></v-list-item>

        <v-list-item
          prepend-icon="mdi-file-pdf-box"
          title="PDF (via print)"
          @click="downloadPdf"
        ></v-list-item>

        <v-list-item
          prepend-icon="mdi-file-delimited-outline"
          title="CSV spreadsheet"
          :disabled="!hasData"
          @click="downloadCsv"
        ></v-list-item>

        <v-divider class="my-1"></v-divider>

        <v-list-item
          prepend-icon="mdi-table-eye"
          title="View underlying data"
          :disabled="!hasData"
          @click="showDialog = true"
        ></v-list-item>
      </v-list>
    </v-menu>

    <!-- ── View Data Dialog ──────────────────────────────────────────────── -->
    <v-dialog v-model="showDialog" max-width="960" scrollable>
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon class="mr-2" color="primary">mdi-table</v-icon>
          {{ title }}
        </v-card-title>
        <v-divider></v-divider>
        <v-card-text style="max-height: 60vh">
          <v-data-table
            :headers="tableHeaders"
            :items="formattedData"
            density="compact"
            :items-per-page="20"
            class="data-preview-table"
          >
            <template
              v-for="h in numericKeys"
              :key="h"
              #[`item.${h}`]="{ item }"
            >
              {{ formatCell(item[h]) }}
            </template>
          </v-data-table>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" size="small" @click="showDialog = false"
            >Close</v-btn
          >
          <v-btn
            color="primary"
            variant="tonal"
            size="small"
            prepend-icon="mdi-file-delimited-outline"
            @click="downloadCsv"
            >Download CSV</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{
  /** Ref to the <ag-charts> Vue component instance */
  chartRef: any
  /** Human-readable title — used as dialog header and download filename */
  title: string
  /** Current data array bound to the chart */
  data: any[]
}>()

const showDialog = ref(false)

const hasData = computed(
  () => Array.isArray(props.data) && props.data.length > 0
)

// ── Table headers auto-derived from data keys ─────────────────────────────────

const tableHeaders = computed(() => {
  if (!hasData.value) return []
  return Object.keys(props.data[0]).map((key) => ({
    title: key.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase()),
    key,
    sortable: true
  }))
})

const numericKeys = computed(() => {
  if (!hasData.value) return []
  return Object.keys(props.data[0]).filter(
    (k) => typeof props.data[0][k] === 'number'
  )
})

const formattedData = computed(() => props.data ?? [])

const formatCell = (val: any) => {
  if (typeof val !== 'number') return val
  return Number.isInteger(val)
    ? val
    : val.toLocaleString('en-ZA', { maximumFractionDigits: 4 })
}

// ── PNG download ──────────────────────────────────────────────────────────────

const downloadPng = () => {
  props.chartRef?.chart?.download(props.title)
}

// ── PDF via print dialog ──────────────────────────────────────────────────────
// Gets the chart canvas, embeds it in a minimal printable HTML page, and opens
// the browser/Electron print dialog so the user can "Save as PDF".

const downloadPdf = () => {
  // Try AG Charts' native base64 first; fall back to canvas querySelector
  let dataUrl: string | null = null

  try {
    // AG Charts Enterprise exposes toBase64Image on the chart instance
    const result = props.chartRef?.chart?.toBase64Image?.('image/png', 1)
    if (typeof result === 'string' && result.startsWith('data:')) {
      dataUrl = result
    }
  } catch (_) {
    // noop — will fall back below
  }

  if (!dataUrl) {
    // Fall back: grab the canvas element rendered inside the chart component
    const canvas: HTMLCanvasElement | null =
      props.chartRef?.$el?.querySelector?.('canvas')
    if (canvas) {
      dataUrl = canvas.toDataURL('image/png')
    }
  }

  if (!dataUrl) {
    console.warn('ChartMenu: could not obtain chart image for PDF export')
    return
  }

  const safe = props.title.replace(/</g, '&lt;').replace(/>/g, '&gt;')
  const scriptClose = '</scr' + 'ipt>'
  const html = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>${safe}</title>
  <style>
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: Arial, sans-serif; padding: 24px; }
    h2 { font-size: 16px; margin-bottom: 16px; color: #333; }
    img { max-width: 100%; display: block; }
    @media print {
      @page { margin: 15mm; }
      body { padding: 0; }
    }
  </style>
</head>
<body>
  <h2>${safe}</h2>
  <img src="${dataUrl}" />
  <script>window.onload = function(){ window.print(); }${scriptClose}
</body>
</html>`

  const blob = new Blob([html], { type: 'text/html;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const w = window.open(url, '_blank')
  // Clean up the object URL after the window has had time to load
  if (w) {
    w.onunload = () => URL.revokeObjectURL(url)
  } else {
    setTimeout(() => URL.revokeObjectURL(url), 30_000)
  }
}

// ── CSV download ──────────────────────────────────────────────────────────────

const downloadCsv = () => {
  if (!hasData.value) return

  const keys = Object.keys(props.data[0])

  const escape = (val: any) => {
    const str = val === null || val === undefined ? '' : String(val)
    return str.includes(',') || str.includes('"') || str.includes('\n')
      ? `"${str.replace(/"/g, '""')}"`
      : str
  }

  const rows = [
    keys.map(escape).join(','),
    ...props.data.map((row) => keys.map((k) => escape(row[k])).join(','))
  ]

  const blob = new Blob([rows.join('\r\n')], {
    type: 'text/csv;charset=utf-8;'
  })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.title.replace(/[/\\?%*:|"<>]/g, '-')}.csv`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}
</script>

<style scoped>
.chart-menu-wrapper {
  display: flex;
  justify-content: flex-end;
}

.data-preview-table :deep(th) {
  font-weight: 600;
  white-space: nowrap;
}
</style>
