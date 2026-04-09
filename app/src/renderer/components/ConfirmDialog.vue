<template>
  <v-dialog v-model="dialog" persistent max-width="500" @keydown.esc="cancel">
    <v-card>
      <v-card-title class="headline"
        ><v-icon class="mr-3" color="red" size="25">mdi-alert-circle</v-icon
        >{{ dialogTitle }}</v-card-title
      >
      <v-card-text>{{ dialogMessage }}</v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          variant="outlined"
          class="mb-3"
          size="small"
          rounded
          color="red"
          @click="cancel"
          >No</v-btn
        >
        <v-btn
          variant="outlined"
          class="mb-3"
          size="small"
          rounded
          color="red"
          @click="agree"
          >Yes</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const dialog = ref(false)
// Use different names for props to avoid conflict with the setup context
const dialogTitle = ref('')
const dialogMessage = ref('')
const resolve = ref<null | ((value: unknown) => void)>(null)
const reject = ref<null | ((reason?: any) => void)>(null)

const open = (title: string, message: string) => {
  dialog.value = true
  dialogTitle.value = title
  dialogMessage.value = message
  return new Promise((_resolve, _reject) => {
    resolve.value = _resolve
    reject.value = _reject
  })
}

const cancel = () => {
  if (reject.value) {
    reject.value(false)
  }
  dialog.value = false
}

const agree = () => {
  if (resolve.value) {
    resolve.value(true)
  }
  dialog.value = false
}

defineExpose({
  open
})
</script>
