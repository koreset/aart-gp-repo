<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Roles Management</span>
          </template>
          <template #default>
            <v-row>
              <v-col>
                <v-btn
                  v-if="hasPermission('system:manage_roles')"
                  color="primary"
                  size="small"
                  rounded
                  @click="openNewRole"
                  >Add a Role</v-btn
                >
              </v-col>
            </v-row>

            <v-row v-if="addRole">
              <v-col>
                <v-card variant="outlined" class="pa-4">
                  <v-row>
                    <v-col cols="12" md="4">
                      <v-text-field
                        v-model="roleName"
                        label="Role Name"
                        variant="outlined"
                        density="compact"
                      />
                    </v-col>
                    <v-col cols="12" md="8">
                      <v-text-field
                        v-model="roleDescription"
                        label="Description"
                        variant="outlined"
                        density="compact"
                      />
                    </v-col>
                  </v-row>

                  <v-divider class="my-2" />

                  <div class="d-flex align-center">
                    <v-icon class="mr-2">mdi-shield-key-outline</v-icon>
                    <span class="text-subtitle-1 font-weight-medium"
                      >Permissions</span
                    >
                  </div>
                  <p class="text-caption text-medium-emphasis mb-3">
                    Tick a module to grant access to that menu tab.
                    Function-level permissions inside a module become available
                    once module access is granted.
                  </p>

                  <v-alert
                    v-if="isSuperuser"
                    type="warning"
                    variant="tonal"
                    density="compact"
                    class="mb-3"
                  >
                    Superuser is enabled — overrides every checkbox below.
                  </v-alert>

                  <v-checkbox
                    v-model="isSuperuser"
                    label="Superuser (Administrator)"
                    color="error"
                    hint="Grants unrestricted access to every module. Use sparingly."
                    persistent-hint
                    density="compact"
                  />

                  <v-text-field
                    v-model="search"
                    prepend-inner-icon="mdi-magnify"
                    label="Filter permissions"
                    density="compact"
                    variant="outlined"
                    clearable
                    class="mt-3"
                  />

                  <div
                    class="d-flex justify-space-between align-center flex-wrap text-caption text-medium-emphasis my-2"
                  >
                    <span>
                      <strong>{{ baselineCount }}</strong> module{{
                        baselineCount === 1 ? '' : 's'
                      }}
                      + <strong>{{ specialCount }}</strong> function{{
                        specialCount === 1 ? '' : 's'
                      }}
                      selected
                      <span v-if="isDirty" class="ml-2 text-warning">
                        · {{ dirtyCount }} unsaved change{{
                          dirtyCount === 1 ? '' : 's'
                        }}
                      </span>
                    </span>
                    <span v-if="search">
                      {{ filteredCategories.length }} of
                      {{ categories.length }} sections shown
                    </span>
                  </div>

                  <div v-if="!isSuperuser && filteredCategories.length > 0">
                    <v-row no-gutters class="role-editor-grid">
                      <!-- Vertical module list (left pane) -->
                      <v-col cols="12" md="3" class="role-editor-sidebar">
                        <v-list
                          density="compact"
                          nav
                          class="pa-1 module-list"
                          color="primary"
                          :selected="activeTab ? [activeTab] : []"
                        >
                          <v-list-item
                            v-for="cat in filteredCategories"
                            :key="cat.category"
                            :value="cat.category"
                            :active="activeTab === cat.category"
                            @click="activeTab = cat.category"
                          >
                            <v-list-item-title class="text-body-2">{{
                              cat.category
                            }}</v-list-item-title>
                            <template #append>
                              <v-chip
                                size="x-small"
                                variant="tonal"
                                :color="
                                  categorySelectionCount(cat) > 0
                                    ? 'primary'
                                    : undefined
                                "
                              >
                                {{ categorySelectionCount(cat) }}/{{
                                  categoryTotalCount(cat)
                                }}
                              </v-chip>
                            </template>
                          </v-list-item>
                        </v-list>
                      </v-col>

                      <!-- Active category content (right pane) -->
                      <v-col
                        cols="12"
                        md="9"
                        class="role-editor-content pl-md-3"
                      >
                        <template
                          v-for="cat in filteredCategories"
                          :key="cat.category"
                        >
                          <div v-if="cat.category === activeTab">
                            <div
                              v-for="group in cat.groups"
                              :key="group.baseline.slug"
                              class="module-group mb-3"
                            >
                              <div class="module-header">
                                <v-checkbox-btn
                                  :model-value="
                                    selected.has(group.baseline.slug)
                                  "
                                  :indeterminate="isPartial(group)"
                                  color="primary"
                                  hide-details
                                  density="compact"
                                  @update:model-value="toggleBaseline(group)"
                                />
                                <div class="module-header-text">
                                  <div class="font-weight-medium">
                                    {{ group.baseline.name }}
                                  </div>
                                  <div
                                    v-if="group.baseline.description"
                                    class="text-caption text-medium-emphasis"
                                  >
                                    {{ group.baseline.description }}
                                  </div>
                                </div>
                                <v-chip
                                  v-if="group.specials.length > 0"
                                  size="x-small"
                                  variant="tonal"
                                  :color="
                                    selected.has(group.baseline.slug)
                                      ? 'primary'
                                      : undefined
                                  "
                                >
                                  {{ countSelectedSpecials(group) }} /
                                  {{ group.specials.length }}
                                </v-chip>
                              </div>

                              <div
                                v-if="group.specials.length === 0"
                                class="text-caption text-medium-emphasis px-3 pb-2"
                              >
                                No additional function-level permissions in this
                                section.
                              </div>

                              <div v-else class="px-3 pb-2">
                                <div class="d-flex justify-end ga-1 mb-1">
                                  <v-btn
                                    size="x-small"
                                    variant="text"
                                    :disabled="
                                      !selected.has(group.baseline.slug)
                                    "
                                    @click="selectAllInGroup(group)"
                                    >Select all</v-btn
                                  >
                                  <v-btn
                                    size="x-small"
                                    variant="text"
                                    :disabled="
                                      !group.specials.some((s) =>
                                        selected.has(s.slug)
                                      )
                                    "
                                    @click="clearAllInGroup(group)"
                                    >Clear all</v-btn
                                  >
                                </div>
                                <div
                                  v-for="s in filteredSpecials(group)"
                                  :key="s.slug"
                                  class="permission-row"
                                  :title="s.description"
                                >
                                  <v-checkbox-btn
                                    :model-value="selected.has(s.slug)"
                                    :disabled="
                                      !selected.has(group.baseline.slug)
                                    "
                                    hide-details
                                    density="compact"
                                    @update:model-value="
                                      toggleSpecial(group, s, $event)
                                    "
                                  />
                                  <span class="permission-name">{{
                                    s.name
                                  }}</span>
                                  <span
                                    v-if="s.description"
                                    class="permission-desc text-medium-emphasis"
                                  >
                                    {{ s.description }}
                                  </span>
                                </div>
                              </div>
                            </div>
                          </div>
                        </template>
                      </v-col>
                    </v-row>
                  </div>
                  <div
                    v-else-if="!isSuperuser"
                    class="text-center text-medium-emphasis py-6"
                  >
                    No permissions match your search.
                  </div>

                  <v-row class="mt-4">
                    <v-col class="d-flex align-center">
                      <v-btn
                        color="primary"
                        size="small"
                        rounded
                        :disabled="!roleName || !isDirty"
                        @click="saveRole"
                      >
                        Save<span v-if="isDirty">
                          ({{ dirtyCount }} change{{
                            dirtyCount === 1 ? '' : 's'
                          }})</span
                        >
                      </v-btn>
                      <v-btn
                        size="small"
                        variant="text"
                        class="ml-2"
                        @click="discardChanges"
                      >
                        {{ isDirty ? 'Discard changes' : 'Close' }}
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col cols="12" md="5">
                <v-text-field
                  v-model="rolesSearch"
                  prepend-inner-icon="mdi-magnify"
                  label="Search roles"
                  density="compact"
                  variant="outlined"
                  clearable
                  hide-details
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-data-table
                  :headers="roleHeaders"
                  :items="userRoles"
                  :search="rolesSearch"
                  :items-per-page="10"
                  :items-per-page-options="[5, 10, 25, 50]"
                  density="comfortable"
                  class="trans-tables"
                >
                  <template v-slot:[`item.actions`]="{ item }">
                    <div
                      class="d-flex justify-end align-center ga-1 flex-nowrap"
                    >
                      <v-btn
                        size="small"
                        variant="text"
                        class="text-none"
                        @click="viewPermissions(item)"
                      >
                        View Permissions
                      </v-btn>
                      <v-btn
                        v-if="hasPermission('system:manage_roles')"
                        size="small"
                        variant="text"
                        class="text-none"
                        @click="assignPermissions(item)"
                      >
                        Edit Permissions
                      </v-btn>
                      <v-btn
                        v-if="hasPermission('system:manage_roles')"
                        size="small"
                        variant="text"
                        icon
                        @click="deleteRole(item)"
                      >
                        <v-icon color="red">mdi-delete</v-icon>
                      </v-btn>
                    </div>
                  </template>
                  <template #no-data>
                    <div class="text-center text-medium-emphasis py-4">
                      No roles found.
                    </div>
                  </template>
                </v-data-table>
              </v-col>
            </v-row>

            <v-row v-if="rolePermissions.length > 0">
              <v-col>
                <v-table>
                  <thead>
                    <tr class="table-row">
                      <th class="table-col text-left">Permission Name</th>
                      <th class="table-col text-left">Description</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="permission in rolePermissions"
                      :key="permission.slug"
                    >
                      <td>{{ permission.name }}</td>
                      <td>{{ permission.description }}</td>
                    </tr>
                  </tbody>
                  <tfoot>
                    <tr class="table-row">
                      <td colspan="2" class="text-center">
                        <v-btn size="small" variant="text" @click="closeTable"
                          >Close</v-btn
                        >
                      </td>
                    </tr>
                  </tfoot>
                </v-table>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" :timeout="timeout">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false"
          >Close</v-btn
        >
      </template>
    </v-snackbar>
    <confirm-dialog ref="confirmDeleteAction" />
    <confirm-dialog ref="confirmDisableBaseline" />
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { computed, onMounted, ref, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

interface GPPermission {
  id?: number
  slug: string
  name: string
  description: string
  category?: string
  tier?: 'baseline' | 'special'
  parent_slug?: string
  display_order?: number
  disabled?: boolean
}

interface PermissionGroup {
  baseline: GPPermission
  specials: GPPermission[]
}

interface CategorySection {
  category: string
  groups: PermissionGroup[]
}

const { hasPermission } = usePermissionCheck()
const snackbar = ref(false)
const snackbarMessage = ref('')
const timeout = ref(3000)
const confirmDeleteAction = ref()
const confirmDisableBaseline = ref()

const addRole = ref(false)
const roleName = ref('')
const roleDescription = ref('')
const roleId = ref(0)
const permissions = ref<GPPermission[]>([])
const selected = ref<Set<string>>(new Set())
const initialSelected = ref<Set<string>>(new Set())
const initialRoleName = ref('')
const initialRoleDescription = ref('')
const search = ref('')
const activeTab = ref<string>('')
const rolePermissions = ref<GPPermission[]>([])

const userRoles = ref<any[]>([])
const rolesSearch = ref('')
const roleHeaders = [
  { title: 'Role Name', key: 'role_name', sortable: true, width: '20%' },
  { title: 'Description', key: 'description', sortable: true },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false,
    align: 'end' as const,
    width: 320,
    cellProps: { style: 'white-space: nowrap;' }
  }
]

onMounted(() => {
  fetchRoles()
  fetchPermissions()
})

const isSuperuser = computed({
  get: () => selected.value.has('system:admin'),
  set: (v: boolean) => {
    if (v) selected.value.add('system:admin')
    else selected.value.delete('system:admin')
  }
})

const categories = computed<CategorySection[]>(() => {
  const active = permissions.value.filter((p) => !p.disabled)
  const byBaseline = new Map<string, PermissionGroup>()
  for (const p of active) {
    if (p.tier === 'baseline') {
      byBaseline.set(p.slug, { baseline: p, specials: [] })
    }
  }
  for (const p of active) {
    if (
      p.tier === 'special' &&
      p.parent_slug &&
      byBaseline.has(p.parent_slug)
    ) {
      byBaseline.get(p.parent_slug)!.specials.push(p)
    }
  }
  const byCategory = new Map<string, CategorySection>()
  for (const g of byBaseline.values()) {
    const cat = g.baseline.category || 'Other'
    if (!byCategory.has(cat)) byCategory.set(cat, { category: cat, groups: [] })
    byCategory.get(cat)!.groups.push(g)
  }
  for (const sec of byCategory.values()) {
    sec.groups.sort(
      (a, b) =>
        (a.baseline.display_order ?? 0) - (b.baseline.display_order ?? 0)
    )
    for (const g of sec.groups) {
      g.specials.sort((a, b) => (a.display_order ?? 0) - (b.display_order ?? 0))
    }
  }
  // Hide System category — system:admin is rendered as the Superuser toggle
  // and admin:* specials are gated on it.
  return [...byCategory.values()]
    .filter((sec) => sec.category !== 'System')
    .sort(
      (a, b) =>
        (a.groups[0]?.baseline.display_order ?? 0) -
        (b.groups[0]?.baseline.display_order ?? 0)
    )
})

const matchesSearch = (text: string) => {
  if (!search.value) return true
  return text.toLowerCase().includes(search.value.toLowerCase())
}

const filteredCategories = computed<CategorySection[]>(() => {
  if (!search.value) return categories.value
  const out: CategorySection[] = []
  for (const sec of categories.value) {
    const groups = sec.groups.filter((g) => {
      if (matchesSearch(g.baseline.name) || matchesSearch(g.baseline.slug)) {
        return true
      }
      return g.specials.some(
        (s) => matchesSearch(s.name) || matchesSearch(s.slug)
      )
    })
    if (groups.length > 0) out.push({ category: sec.category, groups })
  }
  return out
})

const filteredSpecials = (group: PermissionGroup) => {
  if (!search.value) return group.specials
  if (matchesSearch(group.baseline.name)) return group.specials
  return group.specials.filter(
    (s) => matchesSearch(s.name) || matchesSearch(s.slug)
  )
}

const baselineCount = computed(() => {
  let count = 0
  for (const p of permissions.value) {
    if (p.tier === 'baseline' && selected.value.has(p.slug)) count++
  }
  return count
})

const specialCount = computed(() => {
  let count = 0
  for (const p of permissions.value) {
    if (p.tier === 'special' && selected.value.has(p.slug)) count++
  }
  return count
})

const countSelectedSpecials = (group: PermissionGroup) =>
  group.specials.filter((s) => selected.value.has(s.slug)).length

const categorySelectionCount = (cat: CategorySection) => {
  let n = 0
  for (const g of cat.groups) {
    if (selected.value.has(g.baseline.slug)) n++
    for (const s of g.specials) if (selected.value.has(s.slug)) n++
  }
  return n
}

const categoryTotalCount = (cat: CategorySection) => {
  let n = 0
  for (const g of cat.groups) {
    n += 1 + g.specials.length
  }
  return n
}

const dirtyCount = computed(() => {
  let n = 0
  if (roleName.value !== initialRoleName.value) n++
  if (roleDescription.value !== initialRoleDescription.value) n++
  for (const slug of selected.value) {
    if (!initialSelected.value.has(slug)) n++
  }
  for (const slug of initialSelected.value) {
    if (!selected.value.has(slug)) n++
  }
  return n
})

const isDirty = computed(() => dirtyCount.value > 0)

// Snap activeTab to a sensible default whenever the visible category list
// changes — e.g. after permissions load, when search filters the list, or
// when entering edit mode for an existing role. Prefer the first category
// that already has any selected permission so editors see selections on
// arrival; otherwise fall back to the first visible category.
watch(
  filteredCategories,
  (cats) => {
    if (cats.length === 0) {
      activeTab.value = ''
      return
    }
    if (cats.some((c) => c.category === activeTab.value)) return
    const withSelection = cats.find((c) => categorySelectionCount(c) > 0)
    activeTab.value = (withSelection ?? cats[0]).category
  },
  { immediate: true }
)

const isPartial = (group: PermissionGroup) => {
  const has = countSelectedSpecials(group)
  return (
    selected.value.has(group.baseline.slug) &&
    has > 0 &&
    has < group.specials.length
  )
}

const toggleSpecial = (
  group: PermissionGroup,
  perm: GPPermission,
  on: boolean | null
) => {
  if (on) {
    selected.value.add(perm.slug)
    if (!selected.value.has(group.baseline.slug)) {
      selected.value.add(group.baseline.slug)
      snackbarMessage.value = `Module access enabled for ${group.baseline.name}`
      snackbar.value = true
    }
  } else {
    selected.value.delete(perm.slug)
  }
  // Force reactivity (Set mutation isn't reactive on its own)
  selected.value = new Set(selected.value)
}

const toggleBaseline = async (group: PermissionGroup) => {
  const turnOff = selected.value.has(group.baseline.slug)
  if (turnOff) {
    const childCount = countSelectedSpecials(group)
    if (childCount > 0) {
      const ok = await confirmDisableBaseline.value.open(
        'Disable module access?',
        `This also removes ${childCount} function permission(s) under ${group.baseline.name}.`
      )
      if (!ok) return
      for (const s of group.specials) selected.value.delete(s.slug)
    }
    selected.value.delete(group.baseline.slug)
  } else {
    selected.value.add(group.baseline.slug)
  }
  selected.value = new Set(selected.value)
}

const selectAllInGroup = (group: PermissionGroup) => {
  if (!selected.value.has(group.baseline.slug)) {
    selected.value.add(group.baseline.slug)
  }
  for (const s of group.specials) selected.value.add(s.slug)
  selected.value = new Set(selected.value)
}

const clearAllInGroup = (group: PermissionGroup) => {
  for (const s of group.specials) selected.value.delete(s.slug)
  selected.value = new Set(selected.value)
}

const snapshotInitial = () => {
  initialSelected.value = new Set(selected.value)
  initialRoleName.value = roleName.value
  initialRoleDescription.value = roleDescription.value
}

const openNewRole = () => {
  roleId.value = 0
  roleName.value = ''
  roleDescription.value = ''
  selected.value = new Set()
  search.value = ''
  activeTab.value = ''
  addRole.value = true
  snapshotInitial()
}

const cancelEdit = () => {
  addRole.value = false
  roleId.value = 0
  roleName.value = ''
  roleDescription.value = ''
  selected.value = new Set()
  initialSelected.value = new Set()
  initialRoleName.value = ''
  initialRoleDescription.value = ''
  search.value = ''
  activeTab.value = ''
}

const discardChanges = async () => {
  if (isDirty.value) {
    const ok = await confirmDeleteAction.value.open(
      'Discard changes?',
      `You have ${dirtyCount.value} unsaved change${dirtyCount.value === 1 ? '' : 's'}. Discard them and close the editor?`
    )
    if (!ok) return
  }
  cancelEdit()
}

const closeTable = () => {
  rolePermissions.value = []
}

const assignPermissions = (role: any) => {
  selected.value = new Set(
    (role.permissions || []).map((p: GPPermission) => p.slug)
  )
  roleName.value = role.role_name
  roleDescription.value = role.description
  roleId.value = role.id
  search.value = ''
  activeTab.value = ''
  addRole.value = true
  snapshotInitial()
}

const fetchRoles = async () => {
  try {
    const response = await GroupPricingService.getUserRoles()
    if (response.status !== 200) throw new Error('Network response was not ok')
    userRoles.value = response.data
  } catch (error) {
    console.error('Error fetching roles:', error)
  }
}

const fetchPermissions = async () => {
  try {
    const response = await GroupPricingService.getPermissions()
    if (response.status !== 200) throw new Error('Network response was not ok')
    permissions.value = response.data
  } catch (error) {
    console.error('Error fetching permissions:', error)
  }
}

const saveRole = async () => {
  try {
    const slugMap = new Map(permissions.value.map((p) => [p.slug, p]))
    const payload = [...selected.value]
      .map((slug) => slugMap.get(slug))
      .filter((p): p is GPPermission => !!p)

    await GroupPricingService.createUserRole({
      id: roleId.value,
      role_name: roleName.value,
      description: roleDescription.value,
      permissions: payload
    })

    const resp = await GroupPricingService.getUserRoles()
    if (resp.status !== 200) throw new Error('Network response was not ok')
    userRoles.value = resp.data

    cancelEdit()
  } catch (error: any) {
    console.error('Error saving role:', error.data)
    snackbarMessage.value = 'Error: ' + error.data + '. Please try again.'
    snackbar.value = true
  }
}

const deleteRole = async (role: any) => {
  try {
    const confirm = await confirmDeleteAction.value.open(
      'Delete Role',
      `Are you sure you want to delete the role "${role.role_name}"?`
    )
    if (!confirm) return
    const response = await GroupPricingService.deleteUserRole(role.id)
    if (response.status !== 200) throw new Error('Network response was not ok')
    userRoles.value = userRoles.value.filter((r: any) => r.id !== role.id)
  } catch (error) {
    console.error('Error deleting role:', error)
    snackbarMessage.value = 'cannot delete a role that is in use'
    snackbar.value = true
  }
}

const viewPermissions = (role: any) => {
  GroupPricingService.getRolePermissions(role.id)
    .then((response) => {
      rolePermissions.value = response.data
    })
    .catch((error) => {
      console.error('Error fetching permissions:', error)
    })
}
</script>

<style lang="css" scoped>
.role-editor-grid {
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  overflow: hidden;
}

.role-editor-sidebar {
  border-right: 1px solid rgba(0, 0, 0, 0.08);
  background: rgba(0, 0, 0, 0.015);
  max-height: 540px;
  overflow-y: auto;
}

.role-editor-content {
  max-height: 540px;
  overflow-y: auto;
  padding: 8px;
}

.module-list :deep(.v-list-item) {
  min-height: 36px;
  border-radius: 4px;
  margin-bottom: 2px;
}

.module-group {
  border: 1px solid rgba(0, 0, 0, 0.06);
  border-radius: 4px;
}

.module-header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(0, 0, 0, 0.02);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  text-align: left;
}

.module-header-text {
  flex: 1 1 auto;
  min-width: 0;
  text-align: left;
}

.permission-row {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
  padding: 2px 4px;
  border-radius: 3px;
  min-height: 28px;
  text-align: left;
}

.permission-row:hover {
  background: rgba(0, 0, 0, 0.025);
}

/* Vuetify's v-checkbox-btn renders a v-selection-control which defaults to
   flex-grow: 1, eating all available row width. Pin it to its natural size so
   sibling labels sit immediately to the right of the checkbox. */
.permission-row :deep(.v-selection-control),
.module-header :deep(.v-selection-control) {
  flex: 0 0 auto;
  min-width: 0;
  width: auto;
}

.permission-name {
  font-size: 0.875rem;
  font-weight: 500;
  flex: 0 0 auto;
}

.permission-desc {
  font-size: 0.75rem;
  flex: 0 1 auto;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 960px) {
  .role-editor-sidebar {
    max-height: 200px;
    border-right: none;
    border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  }
  .role-editor-content {
    max-height: none;
  }
  .permission-desc {
    white-space: normal;
  }
}
</style>
