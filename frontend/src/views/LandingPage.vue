<template>
  <div class="flex flex-col justify-center items-center h-screen bg-[#fafafa]">
    <div v-if="isRedirecting" class="text-center">
      <p class="text-xl font-medium mb-4">Redirecting to...</p>
      <p class="text-blue-500">{{ redirectUrl }}</p>
    </div>
    <div v-else class="w-full max-w-3xl px-4">
      <label for="url-input" class="block text-5xl mb-12 font-extrabold text-[#85a7c3] text-center">Shortify Link</label>
      <div class="relative flex w-full mx-auto my-6">
        <input 
          id="url-input"
          type="text" 
          v-model="url"
          placeholder="Enter your URL here" 
          :class="[
            'w-full bg-white rounded-md py-7 px-6 text-base shadow-lg focus:bg-white',
            'border transition duration-200 ease-in-out border-transparent',
            'focus:outline-none'
          ]"
        >
        <button 
          @click="handleShorten" 
          :disabled="!isValid || isLoading"
          class="absolute right-[12px] top-1/2 transform -translate-y-1/2 px-5 py-3 bg-[#85a7c3] text-white rounded-[3px] hover:bg-[#248dc4] text-base font-extralight capitalize disabled:bg-gray-400 disabled:cursor-not-allowed"
        >
          {{ isLoading ? 'Shortening...' : 'Shortify' }}
        </button>
      </div>
      <p class="mt-2 text-sm h-5 text-red-600">{{ errorMessage }}</p>
      <div class="mt-4 min-h-[288px]">
        <div v-if="shortenedUrl">
        <p class="font-medium">Shortened URL:</p>
        <a :href="shortenedUrl" target="_blank" class="text-blue-500 hover:underline">{{ shortenedUrl }}</a>
        <div class="mt-4">
          <p class="font-medium">QR Code:</p>
          <img :src="`data:image/png;base64,${qrCode}`" alt="QR Code" class="mt-2 w-48 h-48" />
        </div>
        </div>
      </div>
      <p v-if="error" class="mt-2 text-sm text-red-600">{{ error }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUrlValidator } from '@/composables/useUrlValidator'
import { useUrlShortener } from '@/composables/useUrlShortener'

const route = useRoute()
const url = ref('')
const { isValid, errorMessage } = useUrlValidator(url)
const { shortenedUrl, qrCode, error, isLoading, shortenUrl, handleRedirect, redirectUrl, isRedirecting } = useUrlShortener()

const handleShorten = async () => {
  if (isValid.value) {
    await shortenUrl(url.value)
  }
}

onMounted(() => {
  const shortId = route.params.shortId
  if (typeof shortId === 'string' && shortId) {
    handleRedirect(shortId)
  }
})
</script>
