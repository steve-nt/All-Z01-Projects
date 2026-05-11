<template>
  <div class="welcome welcome--painting">
    <div class="welcome-card">
      <!-- Painting-style header strip -->
      <div class="welcome-stripe" aria-hidden="true">
        <span class="welcome-stripe-block welcome-stripe-block--1"></span>
        <span class="welcome-stripe-block welcome-stripe-block--2"></span>
        <span class="welcome-stripe-block welcome-stripe-block--3"></span>
        <span class="welcome-stripe-block welcome-stripe-block--4"></span>
        <span class="welcome-stripe-block welcome-stripe-block--5"></span>
      </div>

      <!-- Scrollable body so panel never exceeds viewport -->
      <div class="welcome-card-body">
        <!-- App name on card -->
        <p class="welcome-app-name">
        <span class="welcome-app-name-main">AthensZone</span>
        <span class="welcome-app-name-tag">Social</span>
      </p>

      <!-- Hero illustration (mono + accent) -->
      <div class="welcome-hero" aria-hidden="true">
        <svg class="welcome-hero-svg" viewBox="0 0 200 120" fill="none" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <linearGradient id="welcome-blob" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" style="stop-color:#e2e8f0"/>
              <stop offset="100%" style="stop-color:#cbd5e1"/>
            </linearGradient>
          </defs>
          <circle cx="100" cy="75" r="55" fill="url(#welcome-blob)" opacity="0.6"/>
          <g stroke="#0f172a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="100" cy="58" r="24" fill="#f8fafc"/>
            <path d="M100 84v14 M100 84l-10 12 M100 84l10 12"/>
            <path d="M82 112h36"/>
            <path d="M88 58l-8-8 M112 58l8-8"/>
            <path d="M88 52 Q100 46 112 52" fill="none" stroke="#0f172a"/>
            <path d="M124 50l4-6 4 2-2 6-6-2z" fill="#f8fafc" stroke="#0f172a"/>
            <rect x="68" y="42" width="18" height="28" rx="4" fill="#fff" stroke="#0f172a"/>
            <circle cx="77" cy="50" r="2" fill="#0f172a"/>
            <rect x="72" y="58" width="10" height="6" rx="1" fill="#0f172a" opacity="0.25"/>
          </g>
          <path d="M40 35l2 5 5 2-5 2-2 5-2-5-5-2 5-2z" fill="#0f172a" opacity="0.15"/>
          <path d="M162 28l1.5 4 4 1.5-4 1.5-1.5 4-1.5-4-4-1.5 4-1.5z" fill="#0f172a" opacity="0.15"/>
        </svg>
      </div>

      <h1 class="welcome-title">Connect with people who matter</h1>
      <p class="welcome-tagline">
        Share posts, join groups, and stay in touch—all in one place.
      </p>

      <ul class="welcome-features" aria-label="What you can do">
        <li>
          <span class="welcome-feature-icon" aria-hidden="true">
            <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg"><rect x="8" y="8" width="32" height="28" rx="5" stroke="currentColor" stroke-width="2" fill="var(--surface-2)"/><path d="M8 18h32" stroke="currentColor" stroke-width="2" stroke-linecap="round"/><path d="M14 26h12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" opacity="0.7"/></svg>
          </span>
          <span>Create posts and share updates</span>
        </li>
        <li>
          <span class="welcome-feature-icon" aria-hidden="true">
            <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg"><circle cx="24" cy="18" r="8" stroke="currentColor" stroke-width="2" fill="var(--surface-2)"/><circle cx="14" cy="34" r="6" stroke="currentColor" stroke-width="2" fill="var(--surface-2)"/><circle cx="34" cy="34" r="6" stroke="currentColor" stroke-width="2" fill="var(--surface-2)"/><path d="M18 28c0-2 2.5-4 6-4s6 2 6 4" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>
          </span>
          <span>Join groups and events</span>
        </li>
        <li>
          <span class="welcome-feature-icon" aria-hidden="true">
            <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M8 12v24l6-6h6l10-10V12H8z" stroke="currentColor" stroke-width="2" fill="var(--surface-2)" stroke-linejoin="round"/><circle cx="18" cy="20" r="2" fill="currentColor" opacity="0.6"/><circle cx="24" cy="20" r="2" fill="currentColor" opacity="0.6"/></svg>
          </span>
          <span>Message friends privately or in groups</span>
        </li>
      </ul>

      <!-- CTAs inside the white panel -->
      <div class="welcome-actions">
        <RouterLink to="/login" class="welcome-btn welcome-btn--primary">Log in</RouterLink>
        <RouterLink to="/register" class="welcome-btn welcome-btn--secondary">Create account</RouterLink>
      </div>

      <!-- Painting-style footer strip -->
      <div class="welcome-drawing-2" aria-hidden="true">
        <svg class="welcome-drawing-2-svg" viewBox="0 0 200 80" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="20" y="20" width="50" height="36" rx="6" fill="#f1f5f9" stroke="#0f172a" stroke-width="1.5"/>
          <rect x="26" y="26" width="38" height="24" rx="2" fill="#0f172a"/>
          <rect x="45" y="62" width="22" height="4" rx="1" fill="#0f172a"/>
          <path d="M100 62v-28c0-4 4-8 10-8s10 4 10 8v28" stroke="#475569" stroke-width="2" stroke-linecap="round" fill="none"/>
          <ellipse cx="110" cy="22" rx="18" ry="12" fill="#94a3b8"/>
          <path d="M168 42h-14v18c0 2 3 4 7 4s7-2 7-4V42z" fill="#64748b" stroke="#0f172a" stroke-width="1.5"/>
          <ellipse cx="161" cy="40" rx="8" ry="2.5" fill="#94a3b8"/>
          <circle cx="38" cy="14" r="5" fill="none" stroke="#0f172a" stroke-width="1.2" stroke-dasharray="2 2"/>
          <path d="M152 10l2.5 5 5 1.2-3.5 3.5 1 4.5-5-2.5-5 2.5 1-4.5-3.5-3.5 5-1.2z" fill="#0f172a" opacity="0.2"/>
        </svg>
      </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink } from "vue-router";
</script>

<style scoped>
.welcome {
  height: 100%;
  min-height: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px 12px;
  overflow: hidden;
  box-sizing: border-box;
}

/* Painting-style page background: soft gradients + subtle shapes, stays B&W-connected */
.welcome--painting {
  background:
    radial-gradient(ellipse 120% 80% at 20% 20%, rgba(15, 23, 42, 0.06) 0%, transparent 50%),
    radial-gradient(ellipse 80% 100% at 80% 80%, rgba(15, 23, 42, 0.05) 0%, transparent 50%),
    linear-gradient(160deg, #fafafa 0%, #f1f5f9 40%, #e2e8f0 100%);
}

@media (prefers-color-scheme: dark) {
  .welcome--painting {
    background:
      radial-gradient(ellipse 120% 80% at 20% 20%, rgba(255, 255, 255, 0.04) 0%, transparent 50%),
      radial-gradient(ellipse 80% 100% at 80% 80%, rgba(255, 255, 255, 0.03) 0%, transparent 50%),
      linear-gradient(160deg, #0f172a 0%, #1e293b 40%, #0b0b0b 100%);
  }
}

.welcome-card {
  width: 100%;
  max-width: 420px;
  max-height: 100%;
  background: var(--surface);
  color: var(--surface-text);
  border: 1px solid var(--border);
  border-radius: 24px;
  padding: 0;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.12),
    0 0 0 1px rgba(0, 0, 0, 0.04);
  text-align: center;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Top “painting” stripe: blocks of color (B&W palette + one accent) */
.welcome-stripe {
  display: flex;
  height: 6px;
  margin: 0;
}

.welcome-stripe-block {
  flex: 1;
  display: block;
}

/* Top “painting” stripe: B&W + one accent for a modern painting feel */
.welcome-stripe-block--1 { background: #0f172a; }
.welcome-stripe-block--2 { background: #475569; }
.welcome-stripe-block--3 { background: #0d9488; }
.welcome-stripe-block--4 { background: #e2e8f0; }
.welcome-stripe-block--5 { background: #0f172a; }

.welcome-card-body {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.welcome-app-name {
  margin: 0;
  padding: 8px 18px 0;
  font-size: 1.1rem;
  font-weight: 700;
  letter-spacing: -0.025em;
  color: var(--surface-text);
}

.welcome-app-name-main {
  font-weight: 700;
}

.welcome-app-name-tag {
  font-weight: 500;
  margin-left: 0.25em;
  opacity: 0.85;
}

.welcome-hero {
  margin: 0;
  padding: 12px 18px 8px;
  background: linear-gradient(180deg, rgba(241, 245, 249, 0.8) 0%, transparent 100%);
}

.welcome-hero-svg {
  width: 100%;
  max-width: 180px;
  max-height: 90px;
  height: auto;
  margin: 0 auto;
  display: block;
}

.welcome-title {
  margin: 0 0 6px;
  padding: 0 18px;
  font-size: clamp(1.25rem, 3.5vw, 1.5rem);
  font-weight: 800;
  line-height: 1.2;
  letter-spacing: -0.03em;
  color: var(--surface-text);
}

.welcome-tagline {
  margin: 0 0 14px;
  padding: 0 18px;
  font-size: 0.9rem;
  line-height: 1.4;
  color: var(--surface-text);
  opacity: 0.85;
}

.welcome-features {
  list-style: none;
  padding: 0 18px;
  margin: 0 0 14px;
  width: 100%;
  max-width: 320px;
  margin-left: auto;
  margin-right: auto;
  text-align: left;
}

.welcome-features li {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 0;
  font-size: 0.85rem;
  color: var(--surface-text);
}

.welcome-feature-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--surface-text);
  background: #f1f5f9;
  border-radius: 10px;
}

.welcome-feature-icon svg {
  width: 22px;
  height: 22px;
}

/* Login / Register inside the white panel */
.welcome-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 18px 16px;
}

.welcome-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 12px 20px;
  font-size: 0.95rem;
  font-weight: 600;
  border-radius: 12px;
  text-decoration: none;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.welcome-btn:focus-visible {
  outline: 2px solid var(--surface-text);
  outline-offset: 2px;
}

.welcome-btn--primary {
  background: #0f172a;
  color: #fff;
  border: 2px solid #0f172a;
}

.welcome-btn--primary:hover {
  background: #1e293b;
  border-color: #1e293b;
  box-shadow: 0 4px 14px rgba(15, 23, 42, 0.35);
}

.welcome-btn--secondary {
  background: transparent;
  color: var(--surface-text);
  border: 2px solid var(--border-strong);
}

.welcome-btn--secondary:hover {
  background: #f1f5f9;
  border-color: var(--surface-text);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);
}

.welcome-drawing-2 {
  margin: 0;
  padding: 10px 18px;
  background: #f8fafc;
  border-top: 1px dashed var(--border);
}

.welcome-drawing-2-svg {
  width: 100%;
  max-width: 160px;
  max-height: 56px;
  height: auto;
  margin: 0 auto;
  display: block;
}
</style>
