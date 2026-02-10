import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import './assets/main.css'
import App from './App.vue'
import DashboardView from './views/DashboardView.vue'

// Router configuration
const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'dashboard',
            component: DashboardView
        },
        // Add more routes as needed
        // {
        //   path: '/lobby/:id',
        //   name: 'lobby',
        //   component: () => import('./views/LobbyView.vue')
        // }
    ]
})

const app = createApp(App)
app.use(router)
app.mount('#app')
