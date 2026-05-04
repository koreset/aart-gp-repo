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
                    Tick a module to grant access to that menu tab. Function-level
                    permissions inside a module become available once module access
                    is granted.
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
                    class="d-flex justify-space-between text-caption text-medium-emphasis my-2"
                  >
                    <span
                      >{{ baselineCount }} module{{ baselineCount === 1 ? '' : 's' }}
                      + {{ specialCount }} function{{ specialCount === 1 ? '' : 's' }}
                      selected</span
                    >
                    <span>
                      {{ filteredCategories.length }} of
                      {{ categories.length }} sections shown
                    </span>
                  </div>

                  <div v-if="!isSuperuser">
                    <template v-for="cat in filteredCategories" :key="cat.category">
                      <div class="text-overline mt-4 mb-1 text-medium-emphasis">
                        {{ cat.category }}
                      </div>
                      <v-expansion-panels multiple variant="accordion">
                        <v-expansion-panel
                          v-for="group in cat.groups"
                          :key="group.baseline.slug"
                        >
                          <v-expansion-panel-title>
                            <div class="d-flex align-center w-100">
                              <v-checkbox-btn
                                :model-value="selected.has(group.baseline.slug)"
                                :indeterminate="isPartial(group)"
                                color="primary"
                                @click.stop
                                @update:model-value="toggleBaseline(group)"
                              />
                              <span class="font-weight-medium mr-2">{{
                                group.baseline.name
                              }}</span>
                              <v-chip
                                v-if="group.specials.length > 0"
                                size="x-small"
                                variant="tonal"
                                :color="
                                  selected.has(group.baseline.slug)
                                    ? 'primary'
                                    : 'default'
                                "
                              >
                                {{ countSelectedSpecials(group) }} /
                                {{ group.specials.length }}
                              </v-chip>
                              <v-spacer />
                            </div>
                          </v-expansion-panel-title>
                          <v-expansion-panel-text>
                            <p class="text-caption text-medium-emphasis mb-2">
                              {{ group.baseline.description }}
                            </p>
                            <div
                              v-if="group.specials.length > 0"
                              class="d-flex justify-end mb-2"
                            >
                              <v-btn
                                size="x-small"
                                variant="text"
                                :disabled="!selected.has(group.baseline.slug)"
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
                              v-if="group.specials.length === 0"
                              class="text-caption text-medium-emphasis"
                            >
                              No additional function-level permissions in this
                              section.
                            </div>
                            <v-list density="compact" class="pl-2">
                              <v-list-item
                                v-for="s in filteredSpecials(group)"
                                :key="s.slug"
                                class="px-2"
                              >
                                <template #prepend>
                                  <v-checkbox-btn
                                    :model-value="selected.has(s.slug)"
                                    :disabled="!selected.has(group.baseline.slug)"
                                    @update:model-value="
                                      toggleSpecial(group, s, $event)
                                    "
                                  />
                                </template>
                                <v-list-item-title class="text-body-2">{{
                                  s.name
                                }}</v-list-item-title>
                                <v-list-item-subtitle class="text-caption">{{
                                  s.description
                                }}</v-list-item-subtitle>
                              </v-list-item>
                            </v-list>
                          </v-expansion-panel-text>
                        </v-expansion-panel>
                      </v-expansion-panels>
                    </template>
                  </div>

                  <v-row class="mt-4">
                    <v-col class="d-flex">
                      <v-btn
                        color="primary"
                        size="small"
                        rounded
                        :disabled="!roleName"
                        @click="saveRole"
                        >Save</v-btn
                      >
                      <v-btn
                        size="small"
                        variant="text"
                        class="ml-2"
                        @click="cancelEdit"
                        >Cancel</v-btn
                      >
                    </v-col>
                  </v-row>
                </v-card>
              </v-col>
            </v-row>

            <v-row>
              <v-col>
                <v-table class="trans-tables">
                  <thead>
                    <tr class="table-row">
                      <th class="table-col text-left">Role Name</th>
                      <th class="table-col text-left">Description</th>
                      <th class="table-col text-left">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="role in userRoles" :key="role.role_name">
                      <td>{{ role.role_name }}</td>
                      <td>{{ role.description }}</td>
                      <td class="text-right">
                        <v-btn
                          size="small"
                          variant="text"
                          @click="viewPermissions(role)"
                        >
                          View Permissions
                        </v-btn>
                        <v-btn
                          v-if="hasPermission('system:manage_roles')"
                          size="small"
                          variant="text"
                          @click="assignPermissions(role)"
                        >
                          Edit Permissions
                        </v-btn>
                        <v-btn
                          v-if="hasPermission('system:manage_roles')"
                          size="small"
                          variant="text"
                          icon
                          @click="deleteRole(role)"
                        >
                          <v-icon color="red">mdi-delete</v-icon>
                        </v-btn>
                      </td>
                    </tr>
                  </tbody>
                </v-table>
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
import { computed, onMounted, ref } from 'vue'
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
const search = ref('')
const rolePermissions = ref<GPPermission[]>([])

const userRoles = ref<any[]>([])

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
    if (p.tier === 'special' && p.parent_slug && byBaseline.has(p.parent_slug)) {
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
      g.specials.sort(
        (a, b) => (a.display_order ?? 0) - (b.display_order ?? 0)
      )
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

const openNewRole = () => {
  roleId.value = 0
  roleName.value = ''
  roleDescription.value = ''
  selected.value = new Set()
  search.value = ''
  addRole.value = true
}

const cancelEdit = () => {
  addRole.value = false
  roleId.value = 0
  roleName.value = ''
  roleDescription.value = ''
  selected.value = new Set()
  search.value = ''
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
  addRole.value = true
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

<style lang="css" scoped></style>
