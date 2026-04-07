import { ref, onUnmounted } from 'vue'
import apiClient from '../utils/api'

const API_BASE = '/api/scrim'

// Auto-deteksi protokol WebSocket berdasarkan halaman yang sedang dibuka.
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
    const connectWebSocket = (requestId, onMatchFound, onMatchDeclined, onMatchSuccess) => {
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

                    // ── Match Found ──────────────────────────────────
                    if (data.type === 'SCRIM_MATCH_FOUND') {
                        console.log('🎮 Match found via WebSocket!')
                        if (matchFoundCallback) {
                            matchFoundCallback(data)
                        }
                        // DO NOT disconnect WebSocket here! We need it for MATCH_SUCCESS and MATCH_CANCELLED
                    }

                    if (data.type === 'MATCH_SUCCESS') {
                        console.log('🎉 Match confirmed by ALL users!')
                        if (onMatchSuccess) onMatchSuccess(data)
                        disconnectWebSocket()
                    }

                    // ── Opponent declined the match (legacy single-notify) ──
                    if (data.type === 'MATCH_DECLINED') {
                        console.warn('❌ Opponent declined (MATCH_DECLINED):', data.reason)
                        disconnectWebSocket()
                        clearPolling()
                        if (onMatchDeclined) onMatchDeclined(data)
                    }

                    // ── Canonical: match cancelled, sent to BOTH players ───
                    if (data.type === 'MATCH_CANCELLED') {
                        console.warn('🚫 Match cancelled (MATCH_CANCELLED):', data.reason)
                        disconnectWebSocket()
                        clearPolling()
                        if (onMatchDeclined) onMatchDeclined(data)
                    }
                } catch (e) {
                    console.error('Failed to parse WS message:', e)
                }
            }

            ws.onerror = (err) => {
                console.warn('⚠️ WebSocket error — falling back to polling:', err)
                startPolling(requestId, onMatchFound)
            }

            ws.onclose = (event) => {
                if (event.code !== 1000) {
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
        }, 3000)
    }

    const clearPolling = () => {
        if (pollInterval) {
            clearInterval(pollInterval)
            pollInterval = null
        }
    }

    // ──────────────────────────────────────────────────────────────
    // API: Create Scrim Request  (uses apiClient → JWT auto-attached)
    // ──────────────────────────────────────────────────────────────
    const findMatch = async (formData, onMatchFound, onMatchDeclined, onMatchSuccess) => {
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
            const response = await apiClient.post(`${API_BASE}/request`, payload)
            console.log('📥 Request accepted:', response.data)

            currentRequest.value = response.data
            const requestId = response.data.request_id || response.data.id

            if (requestId) {
                connectWebSocket(requestId, onMatchFound, onMatchDeclined, onMatchSuccess)
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
    // API: Poll Request Status  (uses apiClient → JWT auto-attached)
    // ──────────────────────────────────────────────────────────────
    const checkStatus = async (requestId) => {
        try {
            const response = await apiClient.get(`${API_BASE}/request/${requestId}`)
            return response.data
        } catch (err) {
            console.error('Status check error:', err)
            throw err
        }
    }

    // ──────────────────────────────────────────────────────────────
    // API: Cancel Matchmaking (while searching)
    // ──────────────────────────────────────────────────────────────
    const cancelMatch = async (requestId) => {
        loading.value = true
        try {
            disconnectWebSocket()
            clearPolling()
            await apiClient.post(`${API_BASE}/request/${requestId}/cancel`)
            currentRequest.value = null
        } catch (err) {
            console.error('Cancel error:', err)
        } finally {
            loading.value = false
        }
    }

    // ──────────────────────────────────────────────────────────────
    // API: Confirm Match  (uses apiClient → JWT auto-attached)
    // ──────────────────────────────────────────────────────────────
    const confirmMatch = async (matchId, requestId) => {
        try {
            await apiClient.post(`${API_BASE}/match/${matchId}/confirm?request_id=${requestId}`)
        } catch (err) {
            console.warn('Confirm match warning (may already be resolved):', err)
        }
    }

    // ──────────────────────────────────────────────────────────────
    // API: Reject Match — atomic, broadcasts MATCH_CANCELLED to both players.
    // ──────────────────────────────────────────────────────────────
    const rejectMatch = async (matchId, requestId) => {
        try {
            disconnectWebSocket()
            clearPolling()
            await apiClient.post(`${API_BASE}/match/${matchId}/reject?request_id=${requestId}`)
            currentRequest.value = null
        } catch (err) {
            // Match may already be expired/cancelled — idempotent, that's OK
            console.warn('Reject match warning (may already be resolved):', err)
        }
    }

    // Legacy — kept for backward compatibility
    const declineMatch = rejectMatch

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
        confirmMatch,
        declineMatch,  // alias → rejectMatch
        rejectMatch,
    }
}
