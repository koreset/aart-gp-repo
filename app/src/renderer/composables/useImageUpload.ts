import { ref } from 'vue'

export function useImageUpload() {
  const imagePreview = ref<string | null>(null)
  const logoFile = ref<File | null>(null)

  const handleFileChange = (file: File) => {
    if (!file) {
      return
    }

    if (!file.type.startsWith('image/')) {
      throw new Error('Please select a valid image file')
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      if (e.target?.result) {
        imagePreview.value = e.target.result as string
      }
    }
    reader.readAsDataURL(file)
    logoFile.value = file
  }

  const clearImage = () => {
    imagePreview.value = null
    logoFile.value = null
  }

  return {
    imagePreview,
    logoFile,
    handleFileChange,
    clearImage
  }
}
