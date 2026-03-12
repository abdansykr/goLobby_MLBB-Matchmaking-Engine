import { ref, onUnmounted } from 'vue'
import axios from 'axios'

const API_BASE = '/api/scrim'

// Auto-deteksi protokol WebSocket berdasarkan halaman yang sedang dibuka.
// - Jika browser buka https:// (production/Cloudflare) → pakai wss://
// - Jika browser buka http:// (lokal) → pakai ws://
// VITE_WS_URL dari .env akan selalu diutamakan jika ada.
const _wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
const _wsHost = window.location.host
const WS_BASE = import.meta.env.VITE_WS_URL || `${_wsProtocol}//${_wsHost}/ws`

export function useScrimAPI() {
    const loading = ref(false)
    const error = ref(null)
    const currentRequest = ref(null)

    // Internal WebSocket state
    let ws = null
    let pollInterval = null
    let matchFoundCallback = null

    // ──────────────────────────────────────────────────────────────
    // WebSocket: Real-time match notification (primary mechanism)
    // ──────────────────────────────────────────────────────────────
    const connectWebSocket = (requestId, onMatchFound) => {
        matchFoundCallback = onMatchFound

        // Close any existing connection
        disconnectWebSocket()

        const wsUrl = `${WS_BASE}?request_id=${requestId}`
        console.log('🔌 Connecting WebSocket:', wsUrl)

        try {
            ws = new WebSocket(wsUrl)

            ws.onopen = () => {
                console.log('✅ WebSocket connected for request:', requestId)
                // Cancel polling once WS is established
                clearPolling()
            }

            ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data)
                    console.log('📨 WS message:', data)

                    if (data.type === 'SCRIM_MATCH_FOUND') {
                        console.log('🎮 Match found via WebSocket!')
                        if (matchFoundCallback) {
                            matchFoundCallback(data)
                        }
                        disconnectWebSocket()
                    }
                } catch (e) {
                    console.error('Failed to parse WS message:', e)
                }
            }

            ws.onerror = (err) => {
                console.warn('⚠️ WebSocket error — falling back to polling:', err)
                // Graceful degradation: start polling as fallback
                startPolling(requestId, onMatchFound)
            }

            ws.onclose = (event) => {
                if (event.code !== 1000) {
                    // Abnormal closure — switch to polling
                    console.warn(`WebSocket closed (code ${event.code}), switching to polling`)
                    startPolling(requestId, onMatchFound)
                }
            }
        } catch (e) {
            console.warn('WebSocket not supported, using polling:', e)
            startPolling(requestId, onMatchFound)
        }
    }

    const disconnectWebSocket = () => {
        if (ws) {
            ws.close(1000, 'Component cleanup')
            ws = null
        }
    }

    // ──────────────────────────────────────────────────────────────
    // Polling: Fallback if WebSocket unavailable
    // ──────────────────────────────────────────────────────────────
    const startPolling = (requestId, onMatchFound) => {
        clearPolling()
        console.log('📡 Starting polling fallback for request:', requestId)

        pollInterval = setInterval(async () => {
            try {
                const data = await checkStatus(requestId)

                if (data.status === 'matched') {
                    clearPolling()
                    if (onMatchFound) {
                        onMatchFound({
                            type: 'SCRIM_MATCH_FOUND',
                            opponent_name: data.match?.opponent_name,
                            opponent_number: data.match?.opponent_number,
                            whatsapp_url: data.match?.whatsapp_url,
                            match_id: data.match?.id,
                            expires_in: data.match?.expires_in
                        })
                    }
                } else if (data.status === 'expired' || data.status === 'cancelled') {
                    clearPolling()
                }
            } catch (err) {
                console.error('Polling error:', err)
            }
        }, 3000) // poll every 3 seconds
    }

    const clearPolling = () => {
        if (pollInterval) {
            clearInterval(pollInterval)
            pollInterval = null
        }
    }

    // ──────────────────────────────────────────────────────────────
    // API: Create Scrim Request
    // ──────────────────────────────────────────────────────────────
    const findMatch = async (formData, onMatchFound) => {
        loading.value = true
        error.value = null

        try {
            const payload = {
                team_name: formData.teamName,
                whatsapp_number: formData.whatsappNumber,
                category: formData.category,
                rank_weight: parseInt(formData.rankWeight)
            }

            console.log('📤 Sending matchmaking request:', payload)
            const response = await axios.post(`${API_BASE}/request`, payload)
            console.log('📥 Request accepted:', response.data)

            currentRequest.value = response.data
            const requestId = response.data.request_id || response.data.id

            if (requestId) {
                // Start WebSocket first, polling as fallback
                connectWebSocket(requestId, onMatchFound)
            }

            return response.data
        } catch (err) {
            console.error('API Error:', {
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

    // ──────────────────────────────────────────────────────────────
    // API: Poll Request Status (also used as WS fallback)
    // ──────────────────────────────────────────────────────────────
    const checkStatus = async (requestId) => {
        try {
            const response = await axios.get(`${API_BASE}/request/${requestId}`)
            return response.data
        } catch (err) {
            console.error('Status check error:', err)
            throw err
        }
    }

    // ──────────────────────────────────────────────────────────────
    // API: Cancel Matchmaking
    // ──────────────────────────────────────────────────────────────
    const cancelMatch = async (requestId) => {
        loading.value = true
        try {
            disconnectWebSocket()
            clearPolling()
            await axios.post(`${API_BASE}/request/${requestId}/cancel`)
            currentRequest.value = null
        } catch (err) {
            console.error('Cancel error:', err)
        } finally {
            loading.value = false
        }
    }

    // Cleanup on component unmount
    onUnmounted(() => {
        disconnectWebSocket()
        clearPolling()
    })

    return {
        loading,
        error,
        currentRequest,
        findMatch,
        checkStatus,
        cancelMatch,
        disconnectWebSocket
    }
}
