import { ref } from 'vue'
import axios from 'axios'

const API_BASE = '/api/scrim'

export function useScrimAPI() {
    const loading = ref(false)
    const error = ref(null)
    const currentRequest = ref(null)

    // 1. Create Request
    const findMatch = async (formData) => {
        loading.value = true
        error.value = null
        try {
            const payload = {
                team_name: formData.teamName,
                whatsapp_number: formData.whatsappNumber,
                category: formData.category,
                rank_weight: parseInt(formData.rankWeight)
            }

            console.log('Sending payload:', payload)

            const response = await axios.post(`${API_BASE}/request`, payload)

            console.log('Response received:', response.data)
            currentRequest.value = response.data
            return response.data
        } catch (err) {
            console.error('API Error Details:', {
                status: err.response?.status,
                data: err.response?.data,
                message: err.message
            })
            error.value = err.response?.data?.error || err.message || 'Failed to start matchmaking'
            throw err
        } finally {
            loading.value = false
        }
    }

    // 2. Poll Status
    const checkStatus = async (requestId) => {
        try {
            const response = await axios.get(`${API_BASE}/request/${requestId}`)
            return response.data
        } catch (err) {
            console.error('Polling error:', err)
            throw err
        }
    }

    // 3. Cancel Request
    const cancelMatch = async (requestId) => {
        loading.value = true
        try {
            await axios.post(`${API_BASE}/request/${requestId}/cancel`)
            currentRequest.value = null
        } catch (err) {
            console.error('Cancel error:', err)
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        currentRequest,
        findMatch,
        checkStatus,
        cancelMatch
    }
}
