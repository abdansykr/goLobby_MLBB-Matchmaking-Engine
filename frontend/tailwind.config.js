/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{vue,js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                // Dark Fantasy Base Colors
                'midnight': {
                    900: '#0a0e27',
                    800: '#131829',
                    700: '#1a1f3a',
                    600: '#252b4a',
                    500: '#2d3659',
                },
                'slate-deep': {
                    900: '#0f172a',
                    800: '#1e293b',
                    700: '#334155',
                },
                // Magic/Cyber Accents
                'electric-violet': {
                    500: '#8b5cf6',
                    400: '#a78bfa',
                    300: '#c4b5fd',
                },
                'cyan-magic': {
                    500: '#06b6d4',
                    400: '#22d3ee',
                    300: '#67e8f9',
                },
                // Prestige/Rank Colors
                'antique-gold': {
                    500: '#d4af37',
                    400: '#e5c158',
                    300: '#f0d896',
                },
                // Status Colors
                'mythic-glow': '#ff6b00',
                'legend-purple': '#9333ea',
                'epic-blue': '#3b82f6',
            },
            boxShadow: {
                'glow-violet': '0 0 20px rgba(139, 92, 246, 0.5), 0 0 40px rgba(139, 92, 246, 0.2)',
                'glow-cyan': '0 0 20px rgba(6, 182, 212, 0.5), 0 0 40px rgba(6, 182, 212, 0.2)',
                'glow-gold': '0 0 20px rgba(212, 175, 55, 0.5), 0 0 40px rgba(212, 175, 55, 0.2)',
                'inner-glow': 'inset 0 0 20px rgba(139, 92, 246, 0.3)',
                'hexagon': '0 4px 20px rgba(0, 0, 0, 0.5), 0 0 40px rgba(139, 92, 246, 0.1)',
            },
            animation: {
                'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
                'glow': 'glow 2s ease-in-out infinite alternate',
                'spin-slow': 'spin 8s linear infinite',
                'radar': 'radar 3s linear infinite',
                'shimmer': 'shimmer 2s linear infinite',
                'float': 'float 3s ease-in-out infinite',
            },
            keyframes: {
                glow: {
                    '0%': {
                        boxShadow: '0 0 20px rgba(139, 92, 246, 0.5), 0 0 40px rgba(139, 92, 246, 0.2)',
                    },
                    '100%': {
                        boxShadow: '0 0 30px rgba(139, 92, 246, 0.8), 0 0 60px rgba(139, 92, 246, 0.4)',
                    },
                },
                radar: {
                    '0%': {
                        transform: 'rotate(0deg) scale(0.8)',
                        opacity: '0.3',
                    },
                    '50%': {
                        opacity: '1',
                    },
                    '100%': {
                        transform: 'rotate(360deg) scale(1.2)',
                        opacity: '0.3',
                    },
                },
                shimmer: {
                    '0%': {
                        backgroundPosition: '-1000px 0',
                    },
                    '100%': {
                        backgroundPosition: '1000px 0',
                    },
                },
                float: {
                    '0%, 100%': {
                        transform: 'translateY(0px)',
                    },
                    '50%': {
                        transform: 'translateY(-10px)',
                    },
                },
            },
            backgroundImage: {
                'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
                'gradient-magic': 'linear-gradient(135deg, #8b5cf6 0%, #06b6d4 100%)',
                'gradient-gold': 'linear-gradient(135deg, #d4af37 0%, #f0d896 100%)',
                'glass-gradient': 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)',
            },
            backdropBlur: {
                'glass': '12px',
            },
        },
    },
    plugins: [],
}
