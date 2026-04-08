/**
 * useAuth — reactive auth state + login/register/logout actions.
 * State bersifat global (module-level refs) agar shared antar komponen.
 */
import { ref, computed } from 'vue'
import apiClient from '../utils/api'

// ── Module-level state (singleton) ───────────────────────────────────────────
const _token = ref(localStorage.getItem('golobby_token') || null)
const _user  = ref(JSON.parse(localStorage.getItem('golobby_user') || 'null'))

export function useAuth() {
    const isLoggedIn = computed(() => !!_token.value)
    const user       = computed(() => _user.value)
    const token      = computed(() => _token.value)

    // ── Helpers ──────────────────────────────────────────────────────────────
    const _persist = (tokenStr, userData) => {
        _token.value = tokenStr
        _user.value  = userData
        localStorage.setItem('golobby_token', tokenStr)
        localStorage.setItem('golobby_user', JSON.stringify(userData))
    }

    const _clear = () => {
        _token.value = null
        _user.value  = null
        localStorage.removeItem('golobby_token')
        localStorage.removeItem('golobby_user')
    }

    // ── Register ─────────────────────────────────────────────────────────────
    /**
     * @returns {{ message: string, username: string }}
     * @throws  Error with .message from server
     */
    const register = async ({ username, email, password, whatsapp }) => {
        try {
            const res = await apiClient.post('/api/auth/register', { username, email, password, whatsapp })
            return res.data
        } catch (err) {
            const msg = err.response?.data?.error || 'Gagal mendaftar'
            throw new Error(msg)
        }
    }

    // ── Login ─────────────────────────────────────────────────────────────────
    /**
     * @returns {{ token, username, email, user_id }}
     * @throws  Error with .message from server
     */
    const login = async ({ email, password }) => {
        try {
            const res = await apiClient.post('/api/auth/login', { email, password })
            const { token: tokenStr, username, email: userEmail, user_id, whatsapp_number, avatar_url } = res.data
            _persist(tokenStr, { username, email: userEmail, user_id, whatsapp_number, avatar_url })
            return res.data
        } catch (err) {
            const msg = err.response?.data?.error || 'Gagal login'
            throw new Error(msg)
        }
    }

    // ── Logout ────────────────────────────────────────────────────────────────
    const logout = () => {
        _clear()
    }

    // ── Profile Methods ──────────────────────────────────────────────────────
    const updateProfile = async (data) => {
        try {
            const res = await apiClient.put('/api/user/profile', data)
            const updated = res.data.user
            if (_user.value) {
                const newUser = { ..._user.value, ...updated }
                _persist(_token.value, newUser)
            }
            return res.data
        } catch (err) {
            throw new Error(err.response?.data?.error || 'Failed to update profile')
        }
    }

    const uploadAvatar = async (file) => {
        try {
            const formData = new FormData()
            formData.append('avatar', file)
            const res = await apiClient.post('/api/user/avatar', formData, {
                headers: { 'Content-Type': 'multipart/form-data' }
            })
            if (_user.value && res.data.avatar_url) {
                const newUser = { ..._user.value, avatar_url: res.data.avatar_url }
                _persist(_token.value, newUser)
            }
            return res.data
        } catch (err) {
            throw new Error(err.response?.data?.error || 'Failed to upload avatar')
        }
    }

    const changePassword = async ({ currentPassword, newPassword }) => {
        try {
            const res = await apiClient.put('/api/user/password', {
                current_password: currentPassword,
                new_password: newPassword
            })
            return res.data
        } catch (err) {
            throw new Error(err.response?.data?.error || 'Gagal mengubah password')
        }
    }

    const resetPassword = async ({ email, username, newPassword }) => {
        try {
            const res = await apiClient.post('/api/auth/reset-password', {
                email,
                username,
                new_password: newPassword
            })
            return res.data
        } catch (err) {
            throw new Error(err.response?.data?.error || 'Gagal mereset password')
        }
    }

    return { isLoggedIn, user, token, register, login, logout, updateProfile, uploadAvatar, changePassword, resetPassword }
}
