<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <h3 class="mb-0">Invoice Detail</h3>
      </template>
      <template #default>
        <page-header
          title="Invoice Detail"
          :subtitle="detail?.invoice_number"
          icon="mdi-receipt-text-outline"
          :breadcrumbs="[
            { text: 'Invoices', to: { name: 'group-pricing-invoices' } },
            { text: detail?.invoice_number ?? `Invoice ${invoiceId}` }
          ]"
        />

        <!-- Invoice Header Card -->
        <v-row class="mb-4">
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-text>
                <v-row>
                  <v-col cols="12" md="6">
                    <div class="text-caption text-medium-emphasis">Scheme</div>
                    <div class="text-body-1 font-weight-medium mb-2">{{
                      detail?.scheme_name ?? '—'
                    }}</div>
                    <div class="text-caption text-medium-emphasis">Contact</div>
                    <div class="text-body-1"
                      >{{ detail?.contact_name ?? '—' }}
                      {{
                        detail?.contact_email ? `(${detail.contact_email})` : ''
                      }}</div
                    >
                  </v-col>
                  <v-col cols="12" md="3">
                    <div class="text-caption text-medium-emphasis"
                      >Invoice #</div
                    >
                    <div class="text-body-1 font-weight-medium mb-2">{{
                      detail?.invoice_number ?? '—'
                    }}</div>
                    <div class="text-caption text-medium-emphasis"
                      >Issue Date</div
                    >
                    <div class="text-body-1">{{
                      detail?.issue_date ?? '—'
                    }}</div>
                  </v-col>
                  <v-col cols="12" md="3">
                    <div class="text-caption text-medium-emphasis"
                      >Due Date</div
                    >
                    <div
                      class="text-body-1 font-weight-medium mb-2"
                      :class="isOverdue ? 'text-error' : ''"
                    >
                      {{ detail?.due_date ?? '—' }}
                    </div>
                    <div class="text-caption text-medium-emphasis">Balance</div>
                    <div class="text-h6 font-weight-bold text-primary">{{
                      fmtCurrency(detail?.balance)
                    }}</div>
                  </v-col>
                </v-row>
              </v-card-text>
              <v-divider />
              <v-card-actions>
                <v-btn
                  variant="outlined"
                  size="small"
                  prepend-icon="mdi-send-outline"
                  :loading="markingSent"
                  @click="handleMarkSent"
                >
                  Mark as Sent
                </v-btn>
                <v-btn
                  v-if="hasPermission('premiums:adjustment_note')"
                  variant="outlined"
                  size="small"
                  prepend-icon="mdi-note-plus-outline"
                  @click="adjDialog = true"
                >
                  Adjustment Note
                </v-btn>
                <v-btn
                  variant="outlined"
                  size="small"
                  prepend-icon="mdi-printer"
                  @click="printPage"
                >
                  Print
                </v-btn>
                <v-spacer />
                <v-btn
                  v-if="hasPermission('premiums:record_payment')"
                  color="primary"
                  size="small"
                  prepend-icon="mdi-cash-plus"
                  @click="paymentDialog = true"
                >
                  Record Payment
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-col>
        </v-row>

        <!-- Line Items -->
        <v-row class="mb-4">
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Invoice Line Items</v-card-title
              >
              <v-card-text class="pa-0">
                <v-data-table
                  :headers="lineItemHeaders"
                  :items="detail?.line_items ?? []"
                  density="compact"
                  hide-default-footer
                >
                  <template #[`item.base_premium`]="{ item }: { item: any }">{{
                    fmtCurrency(item.base_premium)
                  }}</template>
                  <template #[`item.adjustment`]="{ item }: { item: any }">{{
                    fmtCurrency(item.adjustment)
                  }}</template>
                  <template #[`item.total`]="{ item }: { item: any }">{{
                    fmtCurrency(item.total)
                  }}</template>
                  <template #bottom>
                    <tr class="font-weight-bold">
                      <td colspan="2" class="text-right pa-3">Total</td>
                      <td class="pa-3"></td>
                      <td class="pa-3">{{
                        fmtCurrency(detail?.net_payable)
                      }}</td>
                    </tr>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Adjustments -->
        <v-row v-if="(detail?.adjustments ?? []).length" class="mb-4">
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Adjustment Notes</v-card-title
              >
              <v-card-text class="pa-0">
                <v-data-table
                  :headers="adjHeaders"
                  :items="detail?.adjustments ?? []"
                  density="compact"
                  hide-default-footer
                >
                  <template #[`item.type`]="{ item }: { item: any }">
                    <v-chip
                      :color="item.type === 'credit' ? 'success' : 'error'"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ item.type }}
                    </v-chip>
                  </template>
                  <template #[`item.amount`]="{ item }: { item: any }">{{
                    fmtCurrency(item.amount)
                  }}</template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

        <!-- Payment History -->
        <v-row>
          <v-col cols="12">
            <v-card variant="outlined">
              <v-card-title class="text-subtitle-1 pa-3"
                >Payment History</v-card-title
              >
              <v-card-text class="pa-0">
                <v-data-table
                  :headers="paymentHeaders"
                  :items="detail?.payment_history ?? []"
                  density="compact"
                  hide-default-footer
                >
                  <template #[`item.amount`]="{ item }: { item: any }">{{
                    fmtCurrency(item.amount)
                  }}</template>
                  <template #[`item.status`]="{ item }: { item: any }">
                    <v-chip
                      :color="item.status === 'matched' ? 'success' : 'warning'"
                      size="x-small"
                      variant="tonal"
                    >
                      {{ item.status }}
                    </v-chip>
                  </template>
                </v-data-table>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- Record Payment Dialog -->
    <v-dialog v-model="paymentDialog" max-width="520" persistent>
      <v-card>
        <v-card-title>Record Payment</v-card-title>
        <v-card-text>
          <v-form ref="payFormRef">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="payForm.payment_date"
                  label="Payment Date *"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="[(v) => !!v || 'Payment date is required']"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="payForm.method"
                  label="Method *"
                  :items="['eft', 'cheque', 'cash', 'debit_order']"
                  variant="outlined"
                  density="compact"
                  :rules="[(v) => !!v || 'Method is required']"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="payForm.amount"
                  label="Amount *"
                  type="number"
                  variant="outlined"
                  density="compact"
                  :rules="[
                    (v) => !!v || 'Amount is required',
                    (v) => v > 0 || 'Amount must be greater than 0'
                  ]"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="payForm.bank_reference"
                  label="Bank Reference"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12">
                <v-textarea
                  v-model="payForm.notes"
                  label="Notes"
                  variant="outlined"
                  density="compact"
                  rows="2"
                />
              </v-col>
            </v-row>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="paymentDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="savingPayment"
            @click="handleRecordPayment"
            >Save Payment</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Adjustment Note Dialog -->
    <v-dialog v-model="adjDialog" max-width="480" persistent>
      <v-card>
        <v-card-title>Create Adjustment Note</v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" md="6">
              <v-select
                v-model="adjForm.type"
                label="Type *"
                :items="['credit', 'debit']"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="adjForm.amount"
                label="Amount *"
                type="number"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model="adjForm.reason"
                label="Reason *"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12">
              <v-text-field
                v-model="adjForm.reference"
                label="Reference"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="plain" @click="adjDialog = false">Cancel</v-btn>
          <v-btn color="primary" :loading="savingAdj" @click="handleCreateAdj"
            >Create Note</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar
      v-model="snackbar"
      :color="snackbarColor"
      :timeout="3500"
      centered
    >
      {{ snackbarText }}
      <template #actions>
        <v-btn variant="text" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import PremiumManagementService from '@/renderer/api/PremiumManagementService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import PageHeader from '@/renderer/components/PageHeader.vue'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()

const props = defineProps<{ invoiceId: string }>()

const loading = ref(false)
const markingSent = ref(false)
const savingPayment = ref(false)
const payFormRef = ref<any>(null)
const savingAdj = ref(false)
const paymentDialog = ref(false)
const adjDialog = ref(false)
const snackbar = ref(false)
const snackbarText = ref('')
const snackbarColor = ref('success')

const detail = ref<any>(null)

const today = new Date().toISOString().slice(0, 10)

const isOverdue = computed(
  () => detail.value?.balance > 0 && detail.value?.due_date < today
)

const payForm = ref({
  payment_date: today,
  method: 'eft',
  amount: 0,
  bank_reference: '',
  notes: ''
})

const adjForm = ref({ type: 'credit', amount: 0, reason: '', reference: '' })

const lineItemHeaders = [
  { title: 'Benefit', key: 'benefit' },
  { title: 'Members', key: 'member_count' },
  { title: 'Base Premium', key: 'base_premium' },
  { title: 'Adjustment', key: 'adjustment' },
  { title: 'Total', key: 'total' }
]

const adjHeaders = [
  { title: 'Type', key: 'type' },
  { title: 'Description', key: 'description' },
  { title: 'Amount', key: 'amount' },
  { title: 'By', key: 'created_by' },
  { title: 'Date', key: 'created_at' }
]

const paymentHeaders = [
  { title: 'Date', key: 'payment_date' },
  { title: 'Method', key: 'method' },
  { title: 'Amount', key: 'amount' },
  { title: 'Bank Ref', key: 'bank_reference' },
  { title: 'Status', key: 'status' },
  { title: 'By', key: 'recorded_by' }
]

function fmtCurrency(val: number | undefined) {
  if (val == null) return '—'
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    maximumFractionDigits: 0
  }).format(val)
}

async function loadDetail() {
  loading.value = true
  try {
    const res = await PremiumManagementService.getInvoiceDetail(
      Number(props.invoiceId)
    )
    detail.value = res.data.data
    payForm.value.amount = detail.value?.balance ?? 0
  } catch (e) {
    console.error('Failed to load invoice detail', e)
  } finally {
    loading.value = false
  }
}

async function handleMarkSent() {
  markingSent.value = true
  try {
    await PremiumManagementService.markInvoiceSent(Number(props.invoiceId))
    await loadDetail()
    showSnack('Invoice marked as sent')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    markingSent.value = false
  }
}

async function handleRecordPayment() {
  const { valid } = (await payFormRef.value?.validate()) ?? { valid: true }
  if (!valid) return
  savingPayment.value = true
  try {
    await PremiumManagementService.recordPayment({
      scheme_id: detail.value.scheme_id,
      invoice_id: detail.value.id,
      ...payForm.value
    })
    paymentDialog.value = false
    await loadDetail()
    showSnack('Payment recorded')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    savingPayment.value = false
  }
}

async function handleCreateAdj() {
  savingAdj.value = true
  try {
    const { default: Api } = await import('@/renderer/api/Api')
    await Api().post(
      `/group-pricing/premiums/invoices/${props.invoiceId}/adjustments`,
      {
        invoice_id: Number(props.invoiceId),
        scheme_id: detail.value.scheme_id,
        ...adjForm.value
      }
    )
    adjDialog.value = false
    await loadDetail()
    showSnack('Adjustment note created')
  } catch (e: any) {
    showSnack(e?.response?.data?.message ?? 'Failed', 'error')
  } finally {
    savingAdj.value = false
  }
}

function showSnack(msg: string, color = 'success') {
  snackbarText.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

function printPage() {
  window.print()
}

onMounted(loadDetail)
</script>

<style scoped>
@media print {
  nav,
  header,
  .v-navigation-drawer,
  .v-app-bar {
    display: none !important;
  }
}
</style>
