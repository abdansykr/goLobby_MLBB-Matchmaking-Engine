/**
 * Axios instance terpusat untuk seluruh request ke /api/*.
 * Secara otomatis melampirkan JWT dari localStorage ke setiap request.
 *
 * Gunakan `apiClient` ini sebagai pengganti `axios` langsung untuk
 * semua endpoint yang memerlukan autentikasi.
 */
import axios from 'axios'

const apiClient = axios.create({
    baseURL: '/',
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    },
})

// ── Request interceptor: lampirkan JWT jika tersedia ─────────────────────────
apiClient.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('golobby_token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => Promise.reject(error),
)

// ── Response interceptor: handle 401 secara global ───────────────────────────
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Token kedaluwarsa atau invalid — bersihkan storage
            localStorage.removeItem('golobby_token')
            localStorage.removeItem('golobby_user')
        }
        return Promise.reject(error)
    },
)

export default apiClient
