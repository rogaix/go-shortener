import { ref } from 'vue'
import axios from 'axios'

export function useUrlShortener() {
  const shortenedUrl = ref('')
  const qrCode = ref('')
  const error = ref('')
  const isLoading = ref(false)
  const redirectUrl = ref('')
  const isRedirecting = ref(false)

  const shortenUrl = async (url: string) => {
    isLoading.value = true
    error.value = ''
    try {
      const response = await axios.post('http://localhost:8080/api/shorten', { url })
      shortenedUrl.value = response.data.short_url
      qrCode.value = response.data.qr_code
      console.log("response", response)
    } catch (err) {
      error.value = 'An error occurred while shortening the URL'
      console.error('Error shortening URL:', err)
    } finally {
      isLoading.value = false
    }
  }

  const handleRedirect = async (shortId: string) => {
    redirectUrl.value = ''
    error.value = ''
    try {
      const response = await axios.get(`http://localhost:8080/${shortId}`)
      if (response.data && response.data.long_url) {
        isRedirecting.value = true
        redirectUrl.value = response.data.long_url
        qrCode.value = response.data.qr_code
        setTimeout(() => {
          window.location.href = redirectUrl.value
        }, 1500)
      } else {
        console.error('Invalid response from server:', response.data)
        error.value = 'Invalid response from server'
      }
    } catch (err) {
      console.error('Error fetching redirect URL:', err)
      error.value = 'Error fetching redirect URL'
    }
  }

  return {
    shortenedUrl,
    qrCode,
    error,
    isLoading,
    shortenUrl,
    handleRedirect,
    redirectUrl,
    isRedirecting
  }
}
