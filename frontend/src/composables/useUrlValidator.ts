import { ref, watch, Ref } from 'vue'

export function useUrlValidator(input: Ref<string>) {
  const isValid = ref(false)
  const errorMessage = ref('')

  const validateUrl = (url: string): void => {
    // Regular expression to check for www., http://, https://, and classic domain endings
    const urlPattern = /^(https?:\/\/|www\.)[a-z0-9-]+(\.[a-z0-9-]+)+([/?].*)?$/i

    if (urlPattern.test(url)) {
      isValid.value = true
      errorMessage.value = ''
    } else {
      isValid.value = false
      errorMessage.value = 'Please enter a valid URL with www., http://, or https:// and a valid domain'
    }
  }

  watch(input, (newValue: string) => {
    if (newValue.trim() === '') {
      isValid.value = false
      errorMessage.value = ''
    } else {
      validateUrl(newValue)
    }
  })

  return { isValid, errorMessage }
}
