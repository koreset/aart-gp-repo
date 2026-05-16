<template>
  <v-dialog
    :model-value="open"
    persistent
    max-width="800px"
    @update:model-value="$emit('update:open', $event)"
  >
    <v-card>
      <v-card-title>
        {{ isEdit ? 'Edit Rule' : 'New Rule' }}
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" md="6">
            <v-text-field
              v-model="form.name"
              label="Name *"
              variant="outlined"
              density="compact"
              :error-messages="errors.name"
            />
          </v-col>
          <v-col cols="6" md="3">
            <v-select
              v-model="form.risk_level"
              :items="riskLevelOptions"
              label="Risk Level *"
              variant="outlined"
              density="compact"
            />
          </v-col>
          <v-col cols="6" md="3">
            <v-text-field
              v-model.number="form.priority"
              type="number"
              label="Priority"
              variant="outlined"
              density="compact"
              hint="Higher number wins first"
              persistent-hint
            />
          </v-col>
          <v-col cols="12">
            <v-textarea
              v-model="form.description"
              label="Description"
              variant="outlined"
              density="compact"
              rows="2"
            />
          </v-col>
          <v-col cols="12">
            <v-switch
              v-model="form.enabled"
              color="success"
              label="Enabled"
              density="compact"
              hide-details
            />
          </v-col>
        </v-row>

        <v-divider class="my-3" />

        <div class="text-subtitle-2 mb-2">Conditions</div>
        <v-alert
          v-if="errors.conditions"
          type="error"
          variant="tonal"
          density="compact"
          class="mb-3"
        >
          {{ errors.conditions }}
        </v-alert>
        <ConditionGroupEditor
          v-model:node="rootGroup"
          :features="features"
          :depth="0"
        />

        <v-alert
          v-if="saveError"
          type="error"
          variant="tonal"
          density="compact"
          class="mt-3"
        >
          {{ saveError }}
        </v-alert>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="$emit('update:open', false)">
          Cancel
        </v-btn>
        <v-btn color="primary" :loading="saving" @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, h, defineComponent, type PropType } from 'vue'
import FraudRiskService, {
  type FraudRiskRule,
  type FraudFeatureSpec,
  type FraudRuleNode,
  type FraudRuleGroup,
  type FraudRuleLeaf,
  type FraudRuleOp,
  type FraudRiskLevel
} from '@/renderer/api/FraudRiskService'

const props = defineProps<{
  open: boolean
  rule: FraudRiskRule | null
  features: FraudFeatureSpec[]
}>()

const emit = defineEmits<{
  (e: 'update:open', v: boolean): void
  (e: 'saved'): void
}>()

const riskLevelOptions: { title: string; value: FraudRiskLevel }[] = [
  { title: 'Low Risk', value: 'low' },
  { title: 'Medium Risk', value: 'medium' },
  { title: 'High Risk', value: 'high' },
  { title: 'Critical Risk', value: 'critical' }
]

const isEdit = computed(() => props.rule != null)

const form = ref({
  name: '',
  description: '',
  risk_level: 'medium' as FraudRiskLevel,
  priority: 50,
  enabled: true
})

// Root condition group; rules are evaluated as one "all" group at the top.
const rootGroup = ref<FraudRuleGroup>({ all: [] })

const errors = ref<{ name?: string; conditions?: string }>({})
const saving = ref(false)
const saveError = ref<string | null>(null)

watch(
  () => props.open,
  (v) => {
    if (!v) return
    errors.value = {}
    saveError.value = null
    if (props.rule) {
      form.value = {
        name: props.rule.name,
        description: props.rule.description || '',
        risk_level: props.rule.risk_level,
        priority: props.rule.priority,
        enabled: props.rule.enabled
      }
      rootGroup.value = normaliseRoot(props.rule.conditions)
    } else {
      form.value = {
        name: '',
        description: '',
        risk_level: 'medium',
        priority: 50,
        enabled: true
      }
      rootGroup.value = { all: [] }
    }
  },
  { immediate: true }
)

function normaliseRoot(node: FraudRuleNode | undefined): FraudRuleGroup {
  if (!node) return { all: [] }
  if ('all' in node || 'any' in node) return { ...node } as FraudRuleGroup
  return { all: [node] }
}

function validate(): boolean {
  errors.value = {}
  if (!form.value.name.trim()) {
    errors.value.name = 'Name is required'
  }
  const children =
    (rootGroup.value.all as FraudRuleNode[] | undefined) ||
    (rootGroup.value.any as FraudRuleNode[] | undefined) ||
    []
  if (children.length === 0) {
    errors.value.conditions = 'Add at least one condition'
  }
  return !errors.value.name && !errors.value.conditions
}

async function save() {
  if (!validate()) return
  saving.value = true
  saveError.value = null
  try {
    const payload: Partial<FraudRiskRule> = {
      name: form.value.name.trim(),
      description: form.value.description,
      risk_level: form.value.risk_level,
      priority: form.value.priority || 50,
      enabled: form.value.enabled,
      conditions: rootGroup.value
    }
    if (isEdit.value && props.rule) {
      await FraudRiskService.updateRule(props.rule.id, payload)
    } else {
      await FraudRiskService.createRule(payload)
    }
    emit('saved')
  } catch (err: any) {
    saveError.value =
      err?.response?.data || err?.message || 'Failed to save rule'
  } finally {
    saving.value = false
  }
}

// Recursive in-template editor for an AND/OR group. Implemented inline (no
// <template> recursion gotchas) using h() so the component can refer to
// itself.
const ConditionGroupEditor = defineComponent({
  name: 'ConditionGroupEditor',
  props: {
    node: { type: Object as PropType<FraudRuleGroup>, required: true },
    features: { type: Array as PropType<FraudFeatureSpec[]>, required: true },
    depth: { type: Number, default: 0 }
  },
  emits: ['update:node'],
  setup(props, { emit }) {
    const isAll = computed(() => 'all' in props.node)
    const children = computed<FraudRuleNode[]>({
      get() {
        return (
          (props.node.all as FraudRuleNode[]) ||
          (props.node.any as FraudRuleNode[]) ||
          []
        )
      },
      set(v) {
        emit('update:node', isAll.value ? { all: v } : { any: v })
      }
    })

    function setMode(mode: 'all' | 'any') {
      const arr = children.value
      emit('update:node', mode === 'all' ? { all: arr } : { any: arr })
    }

    function addLeaf() {
      const ruleFeatures = props.features.filter((f) => f.used_by_rules)
      const first = ruleFeatures[0]
      const leaf: FraudRuleLeaf = {
        field: first?.name || 'claim_amount',
        op: first?.kind === 'string' ? 'eq' : 'gt',
        value: first?.kind === 'numeric' || first?.kind === 'bool' ? 0 : ''
      }
      children.value = [...children.value, leaf]
    }

    function addGroup() {
      const group: FraudRuleGroup = { all: [] }
      children.value = [...children.value, group]
    }

    function updateChild(i: number, next: FraudRuleNode) {
      const copy = [...children.value]
      copy[i] = next
      children.value = copy
    }

    function removeChild(i: number) {
      const copy = [...children.value]
      copy.splice(i, 1)
      children.value = copy
    }

    return () => {
      return h(
        'div',
        {
          class: 'condition-group',
          style: {
            border: '1px solid rgba(0,0,0,0.12)',
            borderRadius: '4px',
            padding: '8px',
            marginBottom: '8px',
            background:
              props.depth % 2 === 0 ? 'rgba(0,0,0,0.02)' : 'transparent'
          }
        },
        [
          h(
            'div',
            { class: 'd-flex align-center mb-2', style: { gap: '8px' } },
            [
              h(
                'v-btn-toggle',
                {
                  modelValue: isAll.value ? 'all' : 'any',
                  'onUpdate:modelValue': (v: 'all' | 'any') => setMode(v),
                  density: 'compact',
                  mandatory: true,
                  variant: 'outlined'
                },
                [
                  h('v-btn', { value: 'all', size: 'small' }, () => 'ALL of'),
                  h('v-btn', { value: 'any', size: 'small' }, () => 'ANY of')
                ]
              ),
              h('v-spacer'),
              h(
                'v-btn',
                {
                  size: 'small',
                  variant: 'text',
                  prependIcon: 'mdi-plus',
                  onClick: addLeaf
                },
                () => 'Condition'
              ),
              h(
                'v-btn',
                {
                  size: 'small',
                  variant: 'text',
                  prependIcon: 'mdi-plus-box-multiple',
                  onClick: addGroup
                },
                () => 'Group'
              )
            ]
          ),
          ...children.value.map((child, i) => {
            if ('all' in child || 'any' in child) {
              return h('div', { key: i, class: 'd-flex' }, [
                h(ConditionGroupEditor, {
                  style: { flex: 1 },
                  node: child as FraudRuleGroup,
                  features: props.features,
                  depth: props.depth + 1,
                  'onUpdate:node': (n: FraudRuleNode) => updateChild(i, n)
                }),
                h('v-btn', {
                  icon: 'mdi-close',
                  size: 'small',
                  variant: 'text',
                  color: 'error',
                  onClick: () => removeChild(i)
                })
              ])
            }
            return h(LeafEditor, {
              key: i,
              leaf: child as FraudRuleLeaf,
              features: props.features,
              'onUpdate:leaf': (n: FraudRuleLeaf) => updateChild(i, n),
              onRemove: () => removeChild(i)
            })
          })
        ]
      )
    }
  }
})

const LeafEditor = defineComponent({
  name: 'LeafEditor',
  props: {
    leaf: { type: Object as PropType<FraudRuleLeaf>, required: true },
    features: { type: Array as PropType<FraudFeatureSpec[]>, required: true }
  },
  emits: ['update:leaf', 'remove'],
  setup(props, { emit }) {
    const ruleFeatures = computed(() =>
      props.features.filter((f) => f.used_by_rules)
    )

    const selectedFeature = computed(() =>
      ruleFeatures.value.find((f) => f.name === props.leaf.field)
    )

    const opOptions = computed<{ title: string; value: FraudRuleOp }[]>(() => {
      const kind = selectedFeature.value?.kind
      if (kind === 'string') {
        return [
          { title: '=', value: 'eq' },
          { title: '≠', value: 'ne' },
          { title: 'in', value: 'in' }
        ]
      }
      return [
        { title: '=', value: 'eq' },
        { title: '≠', value: 'ne' },
        { title: '>', value: 'gt' },
        { title: '≥', value: 'gte' },
        { title: '<', value: 'lt' },
        { title: '≤', value: 'lte' }
      ]
    })

    function update<K extends keyof FraudRuleLeaf>(
      key: K,
      v: FraudRuleLeaf[K]
    ) {
      emit('update:leaf', { ...props.leaf, [key]: v })
    }

    function onFieldChange(name: string) {
      const f = ruleFeatures.value.find((x) => x.name === name)
      const defaultValue = f?.kind === 'numeric' || f?.kind === 'bool' ? 0 : ''
      emit('update:leaf', {
        field: name,
        op: f?.kind === 'string' ? 'eq' : 'gt',
        value: defaultValue
      })
    }

    return () =>
      h(
        'div',
        {
          class: 'd-flex align-center mb-2',
          style: { gap: '8px' }
        },
        [
          h('v-select', {
            modelValue: props.leaf.field,
            'onUpdate:modelValue': onFieldChange,
            items: ruleFeatures.value.map((f) => ({
              title: f.description || f.name,
              value: f.name
            })),
            label: 'Feature',
            density: 'compact',
            variant: 'outlined',
            hideDetails: true,
            style: { flex: 2 }
          }),
          h('v-select', {
            modelValue: props.leaf.op,
            'onUpdate:modelValue': (v: FraudRuleOp) => update('op', v),
            items: opOptions.value,
            label: 'Op',
            density: 'compact',
            variant: 'outlined',
            hideDetails: true,
            style: { width: '110px' }
          }),
          selectedFeature.value?.kind === 'string'
            ? selectedFeature.value?.choices?.length
              ? h('v-combobox', {
                  modelValue: props.leaf.value,
                  'onUpdate:modelValue': (v: any) => update('value', v),
                  items: selectedFeature.value.choices,
                  label: 'Value',
                  density: 'compact',
                  variant: 'outlined',
                  hideDetails: true,
                  multiple: props.leaf.op === 'in',
                  style: { flex: 2 }
                })
              : h('v-text-field', {
                  modelValue: props.leaf.value,
                  'onUpdate:modelValue': (v: any) => update('value', v),
                  label: 'Value',
                  density: 'compact',
                  variant: 'outlined',
                  hideDetails: true,
                  style: { flex: 2 }
                })
            : h('v-text-field', {
                modelValue: props.leaf.value,
                'onUpdate:modelValue': (v: any) =>
                  update('value', v === '' ? 0 : Number(v)),
                type: 'number',
                label: 'Value',
                density: 'compact',
                variant: 'outlined',
                hideDetails: true,
                style: { flex: 2 }
              }),
          h('v-btn', {
            icon: 'mdi-close',
            size: 'small',
            variant: 'text',
            color: 'error',
            onClick: () => emit('remove')
          })
        ]
      )
  }
})
</script>
