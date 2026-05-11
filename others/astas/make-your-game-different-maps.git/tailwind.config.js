/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{html,js}"],
  theme: {
    extend: {
      fontFamily: {
        silkscreen: ["Silkscreen", "monospace"],
        pressstart: ["'Press Start 2P'", "monospace"],
      },
      colors: {
        hudBg: "rgba(0,0,0,0.5)",
        hudText: "#e0f7fa",
        accent: "#ffd54f",
        lowTime: "#ff4444",
        btnPrimary: "#9cff57",
        btnSecondary: "#a855f7",
        btnTertiary: "#ffdd33",
        btnPause: "#00e0c6",
      },
      boxShadow: {
        hud: "0 0.3vh 1.5vh rgba(0,0,0,0.7), inset 0 0 1.5vh rgba(2,136,209,0.4)",
        btn: "0 4px 0 #01579b, 0 6px 8px rgba(0,0,0,0.3)",
        btnHover: "0 6px 0 #01579b, 0 8px 12px rgba(0,0,0,0.35)",
      },
      animation: {
        bubble: "ocean-wave 8s linear infinite",
        notif: "notification-pulse 0.8s ease-in-out",
        flood: "bubbleFlood var(--duration) ease-out forwards",
        flash: "flashFade var(--duration) ease-out forwards",
        miss: "miss-float 1.5s ease-out forwards",
        story: "story-fade-in 0.8s ease-in-out",
      },
    },
  },
  plugins: [],
};
