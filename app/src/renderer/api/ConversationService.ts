import Api from '@/renderer/api/Api'

export default {
  createConversation(data: {
    object_type: string
    object_id: number
    title?: string
    participant_emails: string[]
  }) {
    return Api.post('/conversations', data)
  },

  getUserConversations(params: { page?: number; page_size?: number }) {
    return Api.get('/conversations', { params })
  },

  getObjectConversations(objectType: string, objectId: number) {
    return Api.get('/conversations/by-object', {
      params: { object_type: objectType, object_id: objectId }
    })
  },

  getConversation(id: number) {
    return Api.get(`/conversations/${id}`)
  },

  sendMessage(
    conversationId: number,
    data: { body: string; message_type?: string; parent_message_id?: number }
  ) {
    return Api.post(`/conversations/${conversationId}/messages`, data)
  },

  editMessage(messageId: number, body: string) {
    return Api.patch(`/conversations/messages/${messageId}`, { body })
  },

  deleteMessage(messageId: number) {
    return Api.delete(`/conversations/messages/${messageId}`)
  },

  addParticipant(
    conversationId: number,
    data: { user_email: string; user_name?: string }
  ) {
    return Api.post(`/conversations/${conversationId}/participants`, data)
  },

  markConversationRead(conversationId: number) {
    return Api.post(`/conversations/${conversationId}/read`)
  }
}
