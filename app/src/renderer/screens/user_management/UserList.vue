<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">User Management</span>
          </template>
          <template #default>
            <v-row>
              <v-col>
                <v-table class="trans-tables">
                  <thead>
                    <tr class="table-row">
                      <th class="table-col text-left">Name</th>
                      <th class="table-col text-left">Email</th>
                      <th class="table-col text-left">Role</th>
                      <th class="table-col text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="user in orgUsers" :key="user">
                      <td>{{ user.name }}</td>
                      <td>
                        {{ user.email }}
                      </td>
                      <td>
                        {{ user.gp_role }}
                      </td>
                      <td class="text-right">
                        <v-btn
                          v-if="hasPermission('system:manage_users')"
                          size="small"
                          variant="text"
                          @click="assignRole(user)"
                        >
                          Assign Role
                        </v-btn>
                        <v-btn
                          v-if="hasPermission('system:manage_users')"
                          size="small"
                          variant="text"
                          @click="removeRole(user)"
                        >
                          Remove Role
                        </v-btn>
                      </td>
                    </tr>
                  </tbody>
                </v-table>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-dialog v-model="dialog" max-width="400px">
                  <v-card>
                    <v-card-title>Select Role</v-card-title>
                    <v-card-text>
                      <v-select
                        v-model="selectedRole"
                        :items="availableRoles"
                        item-title="role_name"
                        item-value="role_name"
                        label="Choose a role"
                        variant="outlined"
                        density="compact"
                        return-object
                      ></v-select>
                    </v-card-text>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn variant="text" @click="dialog = false"
                        >Cancel</v-btn
                      >
                      <v-btn variant="text" @click="applyRoleChange"
                        >Apply</v-btn
                      >
                    </v-card-actions>
                  </v-card>
                </v-dialog>
              </v-col>
            </v-row>
            <v-dialog></v-dialog>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" :timeout="3000">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>
<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import { onMounted, ref } from 'vue'
import formatValues from '@/renderer/utils/helpers.js'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const { hasPermission } = usePermissionCheck()
const snackbar = ref(false)
const snackbarMessage = ref('')
const columnDefs: any = ref([])
const rowData: any = ref([])
const selectedRole: any = ref(null)
const selectedUser: any = ref(null)
const availableRoles: any = ref([
  { role_name: 'Admin', slug: 'admin' },
  { role_name: 'User', slug: 'user' },
  { role_name: 'Viewer', slug: 'viewer' }
])
const orgUsers: any = ref([])

const dialog = ref(false)

const applyRoleChange = async () => {
  console.log('Selected Role:', selectedRole.value)

  // call api to update the user role
  // updateRole(user)
  // update the user role in the orgUsers array
  const user = orgUsers.value.find((user) => user.id === selectedUser.value.id)

  if (user) {
    user.gp_role = selectedRole.value.role_name
    user.gp_role_id = selectedRole.value.id
    GroupPricingService.updateUserRole(user)
      .then((response) => {
        if (response.status !== 201) {
          throw new Error('Network response was not ok')
        }
        snackbarMessage.value = 'Role assigned successfully'
        snackbar.value = true
      })
      .catch((error) => {
        console.error('Error updating role:', error)
        snackbarMessage.value = 'Failed to assign role'
        snackbar.value = true
      })
  }

  dialog.value = false
}

const assignRole = (user) => {
  selectedRole.value = user.gp_role
  selectedUser.value = user
  dialog.value = true
}

const removeRole = (user) => {
  GroupPricingService.removeUserRole(user)
    .then((response) => {
      if (response.status !== 200) {
        throw new Error('Network response was not ok')
      }
      user.gp_role = ''
      user.gp_role_id = null
      snackbarMessage.value = 'Role removed successfully'
      snackbar.value = true
    })
    .catch((error) => {
      console.error('Error removing role:', error)
      snackbarMessage.value = 'Failed to remove role'
      snackbar.value = true
    })
}

onMounted(async () => {
  const res = await GroupPricingService.getUserRoles()

  availableRoles.value = res.data

  try {
    const result = await GroupPricingService.getOrgUsers({ name: 'AART' })
    orgUsers.value = result.data

    rowData.value = orgUsers.value
    createColumnDefs(rowData.value)
  } catch (error) {}
})

const createColumnDefs = (data) => {
  columnDefs.value = []
  Object.keys(data[0]).forEach((element) => {
    const header: any = {}
    header.headerName = element
    header.field = element
    header.valueFormatter = formatValues
    header.minWidth = 220
    header.filter = true
    header.resizable = true
    header.sortable = true
    columnDefs.value.push(header)
  })
}


</script>
<style lang="css" scoped></style>
