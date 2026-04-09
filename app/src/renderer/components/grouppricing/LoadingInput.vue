<template>
  <v-form>
    <v-container>
      <v-row>
        <v-col cols="4"
          ><v-text-field
            v-model="acceleratedBenefitDiscount"
            v-bind="acceleratedBenefitDiscountAttrs"
            :error-messages="errors.accelerated_benefit_discount"
            variant="outlined"
            density="compact"
            type="number"
            label="Accelerated Benefit Discount"
            placeholder="Enter a value"
          ></v-text-field>
        </v-col>
        <v-col cols="4"
          ><v-text-field
            v-model="commissionLoading"
            v-bind="commissionLoadingAttrs"
            :error-messages="errors.commission_loading"
            variant="outlined"
            density="compact"
            type="number"
            label="Commission Rate (%)"
            placeholder="Enter a value"
          ></v-text-field>
        </v-col>
        <v-col cols="4"
          ><v-text-field
            v-model="profitLoading"
            v-bind="profitLoadingAttrs"
            :error-messages="errors.profit_loading"
            variant="outlined"
            density="compact"
            type="number"
            label="Profit Loading (%)"
            placeholder="Enter a value"
          ></v-text-field>
        </v-col>
        <v-col cols="4"
          ><v-text-field
            v-model="expenseLoading"
            v-bind="expenseLoadingAttrs"
            :error-messages="errors.expense_loading"
            variant="outlined"
            density="compact"
            type="number"
            label="Expense Loading (%)"
            placeholder="Enter a value"
          ></v-text-field>
        </v-col>
        <v-col cols="4"
          ><v-text-field
            v-model="loadingsDiscount"
            v-bind="loadingsDiscountAttrs"
            :error-messages="errors.loadings_discount"
            variant="outlined"
            density="compact"
            type="number"
            label="Overall Premium Discount"
            placeholder="Enter a value"
          ></v-text-field>
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>
<script setup lang="ts">
import { useForm } from 'vee-validate'
import * as yup from 'yup'

const validationSchema = yup.object({
  accelerated_benefit_discount: yup
    .number()
    .typeError('Value must be a number')
    .required('Accelerated Benefit Discount is required')
    .min(0, 'Value must be greater than or equal to 0'),
  commission_loading: yup
    .number()
    .typeError('Value must be a number')
    .required('Commission Rate is required')
    .min(0, 'Value must be greater than or equal to 0'),
  profit_loading: yup
    .number()
    .typeError('Value must be a number')
    .required('Profit Loading is required')
    .min(0, 'Value must be greater than or equal to 0'),
  expense_loading: yup
    .number()
    .typeError('Value must be a number')
    .required('Expense Loading is required')
    .min(0, 'Value must be greater than or equal to 0'),
  loadings_discount: yup // Corrected from 'discount'
    .number()
    .typeError('Value must be a number')
    .required('Overall Premium Discount is required')
    .min(0, 'Value must be greater than or equal to 0')
})

const { handleSubmit, errors, defineField } = useForm({
  validationSchema,
  initialValues: {
    accelerated_benefit_discount: 0, // Added optional chaining for safety
    commission_loading: 0,
    profit_loading: 0,
    expense_loading: 0,
    loadings_discount: 0 // form field 'loadings_discount' maps to store 'discount'
  }
})

const [acceleratedBenefitDiscount, acceleratedBenefitDiscountAttrs] =
  defineField('accelerated_benefit_discount')
const [commissionLoading, commissionLoadingAttrs] =
  defineField('commission_loading')
const [profitLoading, profitLoadingAttrs] = defineField('profit_loading')
const [expenseLoading, expenseLoadingAttrs] = defineField('expense_loading')
const [loadingsDiscount, loadingsDiscountAttrs] =
  defineField('loadings_discount') // This is the vee-validate field name

const validateForm = handleSubmit((values) => {
  // Ensure the 'loadings' object exists in the store

  // console.log('Updated groupStore.group_pricing_quote.loadings:', loadingsStore);
  return true // Indicate successful handling
})

defineExpose({
  validateForm
})
</script>
<!-- <script setup lang="ts">
import { useGroupPricingStore } from '@/renderer/store/group_pricing'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
const groupStore = useGroupPricingStore()

const validationSchema = yup.object({
  accelerated_benefit_discount: yup
    .number()
    .required('Accelerated Benefit Discount is required')
    .min(0, 'Value must be greater than or equal to 0'),
  commission_loading: yup
    .number()
    .required('Commission Rate is required')
    .min(0, 'Value must be greater than or equal to 0'),
  profit_loading: yup
    .number()
    .required('Profit Loading is required')
    .min(0, 'Value must be greater than or equal to 0'),
  expense_loading: yup
    .number()
    .required('Expense Loading is required')
    .min(0, 'Value must be greater than or equal to 0'),
  discount: yup
    .number()
    .required('Overall Premium Discount is required')
    .min(0, 'Value must be greater than or equal to 0')
})

const { handleSubmit, errors, defineField } = useForm({
  validationSchema,
  initialValues: {
    accelerated_benefit_discount:
      groupStore.group_pricing_quote.loadings.accelerated_benefit_discount,
    commission_loading: groupStore.group_pricing_quote.loadings.commission_loading,
    profit_loading: groupStore.group_pricing_quote.loadings.profit_loading,
    expense_loading: groupStore.group_pricing_quote.loadings.expense_loading,
    loadings_discount: groupStore.group_pricing_quote.loadings.discount
  }
})

const [acceleratedBenefitDiscount, acceleratedBenefitDiscountAttrs] = defineField(
  'accelerated_benefit_discount'
)
const [commissionLoading, commissionLoadingAttrs] = defineField('commission_loading')
const [profitLoading, profitLoadingAttrs] = defineField('profit_loading')
const [expenseLoading, expenseLoadingAttrs] = defineField('expense_loading')
const [loadingsDiscount, loadingsDiscountAttrs] = defineField('loadings_discount')

const validateForm = handleSubmit((values) => {
  return true
})

defineExpose({
  validateForm
})
</script> -->
<style lang="css" scoped></style>
