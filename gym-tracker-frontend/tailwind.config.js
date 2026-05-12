/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                // Tokyo Night palette
                night: {
                    bg: '#1a1b2e',
                    surface: '#24283b',
                    surfaceAlt: '#2f3346',
                    border: '#3b4261',
                    text: '#c0caf5',
                    muted: '#565f89',
                    blue: '#7aa2f7',
                    teal: '#73daca',
                    cyan: '#b4f9f8',
                    purple: '#bb9af7',
                    green: '#9ece6a',
                    orange: '#e0af68',
                    red: '#f7768e',
                },
            },
            keyframes: {
                'slide-up': {
                    '0%': { transform: 'translateY(100%)', opacity: '0' },
                    '100%': { transform: 'translateY(0)', opacity: '1' },
                },
            },
            animation: {
                'slide-up': 'slide-up 0.3s ease-out',
            },
        },
    },
    plugins: [],
}
