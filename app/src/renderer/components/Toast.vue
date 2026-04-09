<template>
  <v-row>
    <v-col>
      <v-snackbar
        v-model="flash.visible"
        :color="flash.color"
        location="center"
        timeout="4000"
        variant="flat"
      >
        {{ flash.message }}
        <template #actions>
          <v-btn color="white" variant="text" @click="flash.hide()">
            Close
          </v-btn>
        </template>
      </v-snackbar>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { useFlashStore } from '@/renderer/store/flash'
const flash = useFlashStore()
</script>
