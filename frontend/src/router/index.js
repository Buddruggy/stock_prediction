import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: Home
        },
        {
            path: '/history',
            name: 'history',
            component: () => import('../views/History.vue')
        },
        {
            path: '/about',
            name: 'about',
            component: () => import('../views/About.vue')
        }
    ]
})

export default router
