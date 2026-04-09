<template>
  <div v-if="editor" class="editor-container">
    <div class="toolbar">
      <button
        :class="{ 'is-active': editor.isActive('bold') }"
        @click="editor.chain().focus().toggleBold().run()"
      >
        Bold
      </button>
      <button
        :class="{ 'is-active': editor.isActive('italic') }"
        @click="editor.chain().focus().toggleItalic().run()"
      >
        Italic
      </button>
      <button
        :class="{ 'is-active': editor.isActive('strike') }"
        @click="editor.chain().focus().toggleStrike().run()"
      >
        Strike
      </button>
      <button
        :class="{ 'is-active': editor.isActive('paragraph') }"
        @click="editor.chain().focus().setParagraph().run()"
      >
        Paragraph
      </button>
      <button
        :class="{ 'is-active': editor.isActive('heading', { level: 1 }) }"
        @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
      >
        H1
      </button>
      <button
        :class="{ 'is-active': editor.isActive('heading', { level: 2 }) }"
        @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
      >
        H2
      </button>
      <button
        :class="{ 'is-active': editor.isActive('bulletList') }"
        @click="editor.chain().focus().toggleBulletList().run()"
      >
        Bullet List
      </button>
    </div>

    <editor-content :editor="editor" />
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const editor = useEditor({
  content: props.modelValue,
  extensions: [StarterKit],
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})
</script>

<style>
.editor-container {
  border: 1px solid #ccc;
  border-radius: 5px;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  padding: 8px;
  border-bottom: 1px solid #ccc;
  background-color: #f5f5f5;
}

.toolbar button {
  padding: 4px 8px;
  border: 1px solid #ddd;
  background-color: #fff;
  border-radius: 3px;
  cursor: pointer;
}

.toolbar button.is-active {
  background-color: #38a873;
  color: white;
  font-weight: bold;
}

.editor-container .ProseMirror {
  padding: 10px;
  min-height: 150px;
  outline: none;
}

.ProseMirror ul,
.ProseMirror ol {
  padding-left: 25px;
}
</style>
