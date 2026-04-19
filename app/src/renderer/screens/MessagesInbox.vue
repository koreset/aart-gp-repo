<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  useConversationStore,
  type Conversation,
  type ConversationMessage
} from '@/renderer/store/conversations'
import { useAppStore } from '@/renderer/store/app'
import ConversationService from '@/renderer/api/ConversationService'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'

const route = useRoute()
const store = useConversationStore()
const appStore = useAppStore()

const messageBody = ref('')
const replyTo = ref<ConversationMessage | null>(null)
const messagesContainer = ref<HTMLElement | null>(null)
const showNewConversationDialog = ref(false)
const newTitle = ref('')
const newParticipants = ref<string[]>([])
const newObjectType = ref('')
const newObjectId = ref<number | undefined>()
const orgUsers = ref<{ title: string; value: string; name: string }[]>([])

// Mention autocomplete state
const showMentionMenu = ref(false)
const mentionQuery = ref('')
const mentionStartIndex = ref(-1)
const messageInput = ref<any>(null)
const mentionFilteredUsers = computed(() => {
  const q = mentionQuery.value.toLowerCase()
  if (!q) return orgUsers.value
  return orgUsers.value.filter(
    (u) => u.name.toLowerCase().includes(q) || u.value.toLowerCase().includes(q)
  )
})

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
  return new Date(dateStr).toLocaleString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric'
  })
}

function getLastMessage(conv: Conversation): string {
  if (conv.messages && conv.messages.length > 0) {
    const last = conv.messages[conv.messages.length - 1]
    return last.body.length > 50 ? last.body.slice(0, 50) + '...' : last.body
  }
  return 'No messages yet'
}

async function selectConversation(conv: Conversation) {
  await store.loadConversation(conv.id)
  store.markAsRead(conv.id)
  scrollToBottom()
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

async function fetchOrgUsers() {
  try {
    const org = appStore.getOrganisationName
    if (!org) return
    const res = await GroupPricingService.getOrgUsers({ name: org })
    const users = res.data?.data || res.data || []
    orgUsers.value = users.map((u: any) => ({
      title: `${u.name} (${u.email})`,
      value: u.email,
      name: u.name
    }))
  } catch (e) {
    console.error('Failed to fetch org users:', e)
  }
}

async function createNewConversation() {
  if (newParticipants.value.length === 0) return
  try {
    const res = await ConversationService.createConversation({
      object_type: newObjectType.value || 'general',
      object_id: newObjectId.value || 0,
      title: newTitle.value || 'New conversation',
      participant_emails: newParticipants.value
    })
    const conv = res.data?.data || res.data
    store.activeConversation = conv
    showNewConversationDialog.value = false
    newTitle.value = ''
    newParticipants.value = []
    newObjectType.value = ''
    newObjectId.value = undefined
    store.fetchInbox()
  } catch (e) {
    console.error('Failed to create conversation:', e)
  }
}

function setReply(msg: ConversationMessage) {
  replyTo.value = msg
}

function cancelReply() {
  replyTo.value = null
}

function onMessageInput() {
  const val = messageBody.value
  // Find the native input element from the v-text-field ref
  const inputEl = messageInput.value?.$el?.querySelector(
    'input'
  ) as HTMLInputElement | null
  const cursorPos = inputEl?.selectionStart ?? val.length

  // Look backwards from cursor for an '@' that starts a mention
  const textBeforeCursor = val.slice(0, cursorPos)
  const atIndex = textBeforeCursor.lastIndexOf('@')

  if (atIndex >= 0) {
    // Only trigger if '@' is at start or preceded by a space
    const charBefore = atIndex > 0 ? textBeforeCursor[atIndex - 1] : ' '
    if (charBefore === ' ' || charBefore === '\n' || atIndex === 0) {
      const query = textBeforeCursor.slice(atIndex + 1)
      // Close if user typed a space after the query (finished or cancelled)
      if (!query.includes(' ')) {
        mentionStartIndex.value = atIndex
        mentionQuery.value = query
        showMentionMenu.value = true
        return
      }
    }
  }
  showMentionMenu.value = false
}

function insertMention(user: { name: string; value: string }) {
  const val = messageBody.value
  const before = val.slice(0, mentionStartIndex.value)
  const inputEl = messageInput.value?.$el?.querySelector(
    'input'
  ) as HTMLInputElement | null
  const cursorPos = inputEl?.selectionStart ?? val.length
  const after = val.slice(cursorPos)

  messageBody.value = `${before}@${user.value} ${after}`
  showMentionMenu.value = false
  mentionQuery.value = ''
  mentionStartIndex.value = -1

  // Restore focus to input
  nextTick(() => {
    if (inputEl) {
      inputEl.focus()
      const newPos = before.length + user.value.length + 2 // @+email+space
      inputEl.setSelectionRange(newPos, newPos)
    }
  })
}

// Scroll to bottom when messages change (from global WS handler or local send)
watch(
  () => store.activeConversation?.messages?.length,
  () => scrollToBottom()
)

onMounted(async () => {
  await store.fetchInbox()
  fetchOrgUsers()

  // Auto-open conversation from query param
  const convId = route.query.conversation
  if (convId) {
    await store.loadConversation(Number(convId))
    scrollToBottom()
  }
})

onUnmounted(() => {
  store.clearActive()
})
</script>

<template>
  <v-container fluid>
    <base-card :show-actions="false">
      <template #header>
        <div class="d-flex align-center" style="width: 100%">
          <span class="headline">Messages</span>
          <v-spacer />
          <v-btn
            icon
            size="small"
            variant="text"
            color="white"
            @click="showNewConversationDialog = true"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </div>
      </template>
      <template #default>
        <v-row no-gutters style="height: calc(100vh - 200px)">
          <!-- Left panel: conversation list -->
          <v-col
            cols="4"
            style="
              border-right: 1px solid #e0e0e0;
              height: 100%;
              overflow-y: auto;
              background: #fafafa;
            "
          >
            <v-list
              v-if="store.conversations.length > 0"
              density="compact"
              class="py-0"
            >
              <v-list-item
                v-for="conv in store.conversations"
                :key="conv.id"
                :active="activeConv?.id === conv.id"
                style="cursor: pointer"
                class="py-2"
                @click="selectConversation(conv)"
              >
                <v-list-item-title style="font-size: 13px; font-weight: 500">
                  {{ conv.title || `${conv.object_type} #${conv.object_id}` }}
                </v-list-item-title>
                <v-list-item-subtitle style="font-size: 12px">
                  {{ getLastMessage(conv) }}
                </v-list-item-subtitle>
                <template #append>
                  <span style="font-size: 11px; color: #999">
                    {{ formatDate(conv.updated_at) }}
                  </span>
                </template>
              </v-list-item>
            </v-list>

            <div
              v-else
              class="text-center text-grey py-8"
              style="font-size: 13px"
            >
              No conversations yet
            </div>
          </v-col>

          <!-- Right panel: messages -->
          <v-col cols="8" class="d-flex flex-column" style="height: 100%">
            <template v-if="activeConv">
              <!-- Conversation header -->
              <div
                class="pa-3 d-flex align-center"
                style="border-bottom: 1px solid #e0e0e0"
              >
                <div>
                  <div style="font-size: 14px; font-weight: 600">
                    {{
                      activeConv.title ||
                      `${activeConv.object_type} #${activeConv.object_id}`
                    }}
                  </div>
                  <div style="font-size: 12px; color: #666">
                    {{
                      activeConv.participants
                        ?.map((p) => p.user_name || p.user_email)
                        .join(', ')
                    }}
                  </div>
                </div>
              </div>

              <!-- Messages -->
              <div
                ref="messagesContainer"
                style="flex: 1; overflow-y: auto; padding: 16px"
              >
                <div v-for="msg in messages" :key="msg.id" class="mb-4">
                  <div class="d-flex align-start">
                    <v-avatar
                      size="32"
                      :color="
                        msg.sender_email === currentUserEmail
                          ? 'primary'
                          : 'grey-lighten-1'
                      "
                      class="mr-3 mt-1"
                    >
                      <span style="font-size: 12px; color: white">{{
                        getInitials(msg.sender_name)
                      }}</span>
                    </v-avatar>
                    <div style="flex: 1; min-width: 0">
                      <div class="d-flex align-center mb-1">
                        <span style="font-size: 13px; font-weight: 600">{{
                          msg.sender_name
                        }}</span>
                        <span
                          style="font-size: 11px; color: #999; margin-left: 8px"
                        >
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
                      >
                        {{ msg.body }}
                      </p>
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
                  class="text-center text-grey py-8"
                  style="font-size: 13px"
                >
                  No messages yet
                </div>
              </div>

              <!-- Reply indicator -->
              <div
                v-if="replyTo"
                class="px-4 py-1 bg-grey-lighten-4 d-flex align-center"
                style="border-top: 1px solid #e0e0e0"
              >
                <span style="font-size: 12px; color: #666">
                  Replying to {{ replyTo.sender_name }}
                </span>
                <v-spacer />
                <v-btn icon size="x-small" variant="text" @click="cancelReply">
                  <v-icon size="14">mdi-close</v-icon>
                </v-btn>
              </div>

              <!-- Input -->
              <div
                class="pa-3 d-flex align-center"
                style="border-top: 1px solid #e0e0e0; position: relative"
              >
                <div style="flex: 1; position: relative" class="mr-2">
                  <v-text-field
                    ref="messageInput"
                    v-model="messageBody"
                    placeholder="Type a message... (use @ to mention)"
                    density="compact"
                    variant="outlined"
                    hide-details
                    @input="onMessageInput"
                    @keydown.enter.exact="sendMessage"
                    @keydown.escape="showMentionMenu = false"
                  />
                  <v-card
                    v-if="showMentionMenu && mentionFilteredUsers.length > 0"
                    class="mention-dropdown"
                    elevation="8"
                    style="
                      position: absolute;
                      bottom: 100%;
                      left: 0;
                      right: 0;
                      max-height: 200px;
                      overflow-y: auto;
                      z-index: 100;
                      margin-bottom: 4px;
                    "
                  >
                    <v-list density="compact" class="py-0">
                      <v-list-item
                        v-for="u in mentionFilteredUsers"
                        :key="u.value"
                        style="cursor: pointer"
                        @mousedown.prevent="insertMention(u)"
                      >
                        <template #prepend>
                          <v-icon size="18" class="mr-2" color="primary"
                            >mdi-at</v-icon
                          >
                        </template>
                        <v-list-item-title style="font-size: 13px">
                          {{ u.name }}
                        </v-list-item-title>
                        <v-list-item-subtitle style="font-size: 11px">
                          {{ u.value }}
                        </v-list-item-subtitle>
                      </v-list-item>
                    </v-list>
                  </v-card>
                </div>
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

            <div
              v-else
              class="d-flex align-center justify-center"
              style="height: 100%"
            >
              <div class="text-center text-grey">
                <v-icon size="48" class="mb-2">mdi-message-text-outline</v-icon>
                <p style="font-size: 14px"
                  >Select a conversation to view messages</p
                >
              </div>
            </div>
          </v-col>
        </v-row>
      </template>
    </base-card>

    <!-- New conversation dialog -->
    <v-dialog v-model="showNewConversationDialog" max-width="450">
      <v-card>
        <v-card-title style="font-size: 16px">New Conversation</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newTitle"
            label="Title"
            density="compact"
            variant="outlined"
            class="mb-3"
          />
          <v-select
            v-model="newParticipants"
            :items="orgUsers"
            label="Participants"
            density="compact"
            variant="outlined"
            multiple
            chips
            closable-chips
            hint="Select one or more users"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="showNewConversationDialog = false"
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            variant="flat"
            :disabled="newParticipants.length === 0"
            @click="createNewConversation"
          >
            Create
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>
