<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue'
import ConversationService from '@/renderer/api/ConversationService'
import {
  useConversationStore,
  type Conversation,
  type ConversationMessage
} from '@/renderer/store/conversations'
import { useAppStore } from '@/renderer/store/app'

const props = defineProps<{
  objectType: string
  objectId: number
}>()

const store = useConversationStore()
const appStore = useAppStore()

const messageBody = ref('')
const replyTo = ref<ConversationMessage | null>(null)
const messagesContainer = ref<HTMLElement | null>(null)
const showStartDialog = ref(false)
const newTitle = ref('')
const newParticipants = ref('')
const loading = ref(false)

const currentUserEmail = computed(() => appStore.getUser?.email || '')

const activeConv = computed(() => store.activeConversation)
const messages = computed(() => activeConv.value?.messages || [])

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

function getInitials(name: string): string {
  if (!name) return '?'
  return name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

function formatTime(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function loadConversations() {
  loading.value = true
  try {
    const res = await ConversationService.getObjectConversations(
      props.objectType,
      props.objectId
    )
    const data = res.data?.data || res.data
    const convs: Conversation[] = Array.isArray(data) ? data : []
    if (convs.length > 0) {
      store.activeConversation = convs[0]
      scrollToBottom()
    } else {
      store.activeConversation = null
    }
  } catch (e) {
    console.error('Failed to load conversations:', e)
  } finally {
    loading.value = false
  }
}

async function startConversation() {
  const emails = newParticipants.value
    .split(',')
    .map((e) => e.trim())
    .filter((e) => e)
  try {
    const res = await ConversationService.createConversation({
      object_type: props.objectType,
      object_id: props.objectId,
      title:
        newTitle.value || `${props.objectType} #${props.objectId} discussion`,
      participant_emails: emails
    })
    const conv = res.data?.data || res.data
    store.activeConversation = conv
    showStartDialog.value = false
    newTitle.value = ''
    newParticipants.value = ''
  } catch (e) {
    console.error('Failed to create conversation:', e)
  }
}

async function sendMessage() {
  if (!messageBody.value.trim() || !activeConv.value) return
  try {
    await store.sendMessage(
      activeConv.value.id,
      messageBody.value.trim(),
      replyTo.value?.id
    )
    messageBody.value = ''
    replyTo.value = null
    scrollToBottom()
  } catch (e) {
    console.error('Failed to send message:', e)
  }
}

function setReply(msg: ConversationMessage) {
  replyTo.value = msg
}

function cancelReply() {
  replyTo.value = null
}

// Scroll to bottom when messages change (from global WS handler or local send)
watch(
  () => store.activeConversation?.messages?.length,
  () => scrollToBottom()
)

onMounted(() => {
  loadConversations()
})

onUnmounted(() => {
  store.clearActive()
})
</script>

<template>
  <v-card variant="outlined" class="mt-4">
    <v-card-title class="py-2 px-4 d-flex align-center" style="font-size: 14px">
      <v-icon size="18" class="mr-2">mdi-forum</v-icon>
      Discussion
    </v-card-title>
    <v-divider />

    <!-- Loading -->
    <v-card-text v-if="loading" class="text-center py-6">
      <v-progress-circular indeterminate size="24" />
    </v-card-text>

    <!-- No conversation yet -->
    <v-card-text v-else-if="!activeConv" class="text-center py-6">
      <p class="text-grey mb-3" style="font-size: 13px">No discussion yet</p>
      <v-btn
        size="small"
        color="primary"
        variant="outlined"
        @click="showStartDialog = true"
      >
        Start Discussion
      </v-btn>
    </v-card-text>

    <!-- Conversation messages -->
    <template v-else>
      <div
        ref="messagesContainer"
        style="max-height: 400px; overflow-y: auto; padding: 12px"
      >
        <div v-for="msg in messages" :key="msg.id" class="mb-3">
          <div class="d-flex align-start">
            <v-avatar
              size="28"
              :color="
                msg.sender_email === currentUserEmail
                  ? 'primary'
                  : 'grey-lighten-1'
              "
              class="mr-2 mt-1"
            >
              <span style="font-size: 11px; color: white">{{
                getInitials(msg.sender_name)
              }}</span>
            </v-avatar>
            <div style="flex: 1; min-width: 0">
              <div class="d-flex align-center mb-1">
                <span style="font-size: 12px; font-weight: 600">{{
                  msg.sender_name
                }}</span>
                <span style="font-size: 11px; color: #999; margin-left: 8px">
                  {{ formatTime(msg.created_at) }}
                </span>
                <v-chip
                  v-if="msg.is_edited"
                  size="x-small"
                  variant="text"
                  class="ml-1"
                  >edited</v-chip
                >
              </div>
              <div
                v-if="msg.parent_message_id"
                style="
                  font-size: 11px;
                  color: #888;
                  border-left: 2px solid #ccc;
                  padding-left: 8px;
                  margin-bottom: 4px;
                "
              >
                Replying to a message
              </div>
              <p
                style="
                  font-size: 13px;
                  margin: 0;
                  white-space: pre-wrap;
                  word-break: break-word;
                "
                >{{ msg.body }}</p
              >
              <v-btn
                variant="text"
                size="x-small"
                color="grey"
                class="mt-1 pa-0"
                style="min-width: auto"
                @click="setReply(msg)"
              >
                Reply
              </v-btn>
            </div>
          </div>
        </div>

        <div
          v-if="messages.length === 0"
          class="text-center text-grey py-4"
          style="font-size: 13px"
        >
          No messages yet. Start the conversation!
        </div>
      </div>

      <v-divider />

      <!-- Reply indicator -->
      <div
        v-if="replyTo"
        class="px-4 py-1 bg-grey-lighten-4 d-flex align-center"
      >
        <span style="font-size: 12px; color: #666">
          Replying to {{ replyTo.sender_name }}
        </span>
        <v-spacer />
        <v-btn icon size="x-small" variant="text" @click="cancelReply">
          <v-icon size="14">mdi-close</v-icon>
        </v-btn>
      </div>

      <!-- Message input -->
      <div class="pa-3 d-flex align-center">
        <v-text-field
          v-model="messageBody"
          placeholder="Type a message..."
          density="compact"
          variant="outlined"
          hide-details
          class="mr-2"
          @keyup.enter="sendMessage"
        />
        <v-btn
          icon
          size="small"
          color="primary"
          :disabled="!messageBody.trim()"
          @click="sendMessage"
        >
          <v-icon>mdi-send</v-icon>
        </v-btn>
      </div>
    </template>

    <!-- Start conversation dialog -->
    <v-dialog v-model="showStartDialog" max-width="450">
      <v-card>
        <v-card-title style="font-size: 16px">Start Discussion</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newTitle"
            label="Title (optional)"
            density="compact"
            variant="outlined"
            class="mb-3"
          />
          <v-text-field
            v-model="newParticipants"
            label="Participant emails (comma-separated)"
            density="compact"
            variant="outlined"
            hint="e.g. john@company.com, jane@company.com"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showStartDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :disabled="!newParticipants.trim()"
            @click="startConversation"
          >
            Create
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>
