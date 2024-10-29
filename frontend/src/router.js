import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from './views/LandingPage.vue'

const routes = [
  {
    path: '/:shortId?',
    name: 'LandingPage',
    component: LandingPage
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
