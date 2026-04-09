<template>
  <div>
    <h3>Database Setup</h3>
    <v-form>
      <v-container>
        <v-row
          ><v-col>
            <p>
              Required: Base URL of the API Server AART will communicate with.
            </p>
          </v-col></v-row
        >
        <v-row>
          <v-col cols="8">
            <v-text-field
              v-model="baseApi"
              label="API Url"
              variant="outlined"
              density="compact"
            ></v-text-field>
          </v-col>
          <v-col>
            <v-btn color="primary" rounded @click="saveApiUrl"
              >Apply Setting</v-btn
            >
          </v-col>
        </v-row>
      </v-container>
      <v-divider></v-divider>
    </v-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const baseApi = ref('')

const saveApiUrl = () => {
  window.mainApi?.send('msgSaveBaseUrl', baseApi.value)
}
</script>

<style></style>
