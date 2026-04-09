import { defineStore } from 'pinia'
import { ref } from 'vue'
import ConversationService from '@/renderer/api/ConversationService'

export interface ConversationParticipant {
  id: number
  conversation_id: number
  user_email: string
  user_name: string
  joined_at: string
  last_read_at: string | null
}

export interface ConversationMessage {
  id: number
  conversation_id: number
  sender_email: string
  sender_name: string
  body: string
  message_type: string
  parent_message_id: number | null
  is_edited: boolean
  created_at: string
  updated_at: string
}

export interface Conversation {
  id: number
  object_type: string
  object_id: number
  title: string
  created_by: string
  is_archived: boolean
  created_at: string
  updated_at: string
  participants: ConversationParticipant[]
  messages: ConversationMessage[]
}

export const useConversationStore = defineStore('conversations', () => {
  const conversations = ref<Conversation[]>([])
  const activeConversation = ref<Conversation | null>(null)
  const loading = ref(false)
  const total = ref(0)

  async function fetchInbox(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const res = await ConversationService.getUserConversations({
        page,
        page_size: pageSize
      })
      const data = res.data?.data || res.data
      conversations.value = data.conversations || []
      total.value = data.total || 0
    } catch (e) {
      console.error('Failed to fetch conversations:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchObjectConversations(
    objectType: string,
    objectId: number
  ) {
    loading.value = true
    try {
      const res = await ConversationService.getObjectConversations(
        objectType,
        objectId
      )
      const data = res.data?.data || res.data
      conversations.value = Array.isArray(data) ? data : []
    } catch (e) {
      console.error('Failed to fetch object conversations:', e)
    } finally {
      loading.value = false
    }
  }

  async function loadConversation(id: number) {
    loading.value = true
    try {
      const res = await ConversationService.getConversation(id)
      const conv = res.data?.data || res.data
      // Ensure messages array is always initialized
      if (!conv.messages) conv.messages = []
      if (!conv.participants) conv.participants = []
      activeConversation.value = conv
    } catch (e) {
      console.error('Failed to load conversation:', e)
    } finally {
      loading.value = false
    }
  }

  async function sendMessage(
    conversationId: number,
    body: string,
    parentMessageId?: number
  ) {
    const payload: any = { body }
    if (parentMessageId) {
      payload.parent_message_id = parentMessageId
    }
    const res = await ConversationService.sendMessage(conversationId, payload)
    const msg: ConversationMessage = res.data?.data || res.data
    if (
      activeConversation.value &&
      activeConversation.value.id === conversationId
    ) {
      if (!activeConversation.value.messages)
        activeConversation.value.messages = []
      activeConversation.value.messages.push(msg)
    }
    return msg
  }

  function addIncomingMessage(msg: ConversationMessage) {
    if (
      activeConversation.value &&
      activeConversation.value.id === msg.conversation_id
    ) {
      if (!activeConversation.value.messages)
        activeConversation.value.messages = []
      // Avoid duplicates
      if (!activeConversation.value.messages.find((m) => m.id === msg.id)) {
        activeConversation.value.messages.push(msg)
      }
    }
  }

  async function markAsRead(conversationId: number) {
    try {
      await ConversationService.markConversationRead(conversationId)
    } catch (e) {
      console.error('Failed to mark conversation as read:', e)
    }
  }

  function clearActive() {
    activeConversation.value = null
  }

  return {
    conversations,
    activeConversation,
    loading,
    total,
    fetchInbox,
    fetchObjectConversations,
    loadConversation,
    sendMessage,
    addIncomingMessage,
    markAsRead,
    clearActive
  }
})
