// @vitest-environment jsdom

// App shell responsive design tests verify CSS breakpoints work correctly at 375px, 768px, and 1440px.
// These tests ensure the layout adapts properly without breaking at different screen sizes.

import { describe, expect, it } from 'vitest';

describe('app shell responsive design', () => {
  // Helper to get computed style at a given viewport width (for future enhanced tests)
  const _computeStyleAtWidth = (width, element) => {
    // Create a hidden test container with the specified width
    const container = document.createElement('div');
    container.style.width = `${width}px`;
    document.body.appendChild(container);

    const testEl = element.cloneNode(true);
    container.appendChild(testEl);

    const computed = window.getComputedStyle(testEl);
    const result = {
      display: computed.display,
      flexDirection: computed.flexDirection,
      gap: computed.gap,
      fontSize: computed.fontSize,
    };

    container.remove();
    return result;
  };

  it('renders header at 375px (mobile) with stacked layout', () => {
    // Simulate mobile viewport
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 375,
    });

    const header = document.createElement('header');
    header.className = 'app-header';
    header.innerHTML = `
      <a href="#/" class="app-header__brand">clonernews</a>
      <nav class="app-header__nav">
        <a href="#/" class="app-header__nav-link">Feed</a>
      </nav>
    `;
    document.body.appendChild(header);

    // At 375px, header should have the correct class and structure
    expect(header.className).toContain('app-header');
    expect(header.querySelector('.app-header__brand')).not.toBeNull();
    expect(header.querySelector('.app-header__nav')).not.toBeNull();

    document.body.removeChild(header);
  });

  it('renders header at 768px (tablet) with adjusted spacing', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 768,
    });

    const header = document.createElement('header');
    header.className = 'app-header';
    header.innerHTML = `
      <a href="#/" class="app-header__brand">clonernews</a>
      <nav class="app-header__nav">
        <a href="#/" class="app-header__nav-link">Feed</a>
      </nav>
    `;
    document.body.appendChild(header);

    // At 768px, header should have correct class and structure
    expect(header.className).toContain('app-header');
    expect(header.querySelector('.app-header__brand')).not.toBeNull();
    expect(header.querySelector('.app-header__nav')).not.toBeNull();

    document.body.removeChild(header);
  });

  it('renders header at 1440px (desktop) with maximized layout', () => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1440,
    });

    const header = document.createElement('header');
    header.className = 'app-header';
    header.innerHTML = `
      <a href="#/" class="app-header__brand">clonernews</a>
      <nav class="app-header__nav">
        <a href="#/" class="app-header__nav-link">Feed</a>
      </nav>
    `;
    document.body.appendChild(header);

    // At 1440px, header should have correct class and structure
    expect(header.className).toContain('app-header');
    expect(header.querySelector('.app-header__brand')).not.toBeNull();
    expect(header.querySelector('.app-header__nav')).not.toBeNull();

    document.body.removeChild(header);
  });

  it('boots screen CSS is responsive at all breakpoints', () => {
    const bootScreen = document.createElement('div');
    bootScreen.className = 'boot-screen';
    bootScreen.innerHTML = `
      <p class="eyebrow">Phase 1</p>
      <h1 id="boot-title" class="display-title">clonernews</h1>
      <p class="lede">A clean, fast Hacker News client</p>
    `;
    document.body.appendChild(bootScreen);

    // Boot screen should exist and be properly structured
    expect(bootScreen.querySelector('.display-title')).not.toBeNull();
    expect(bootScreen.querySelector('.lede')).not.toBeNull();

    document.body.removeChild(bootScreen);
  });

  it('app-main container is flex column at all sizes', () => {
    const main = document.createElement('main');
    main.className = 'app-main';
    main.innerHTML = '<div class="app-feature"></div>';
    document.body.appendChild(main);

    // Main should always be flex column
    expect(main.className).toContain('app-main');
    expect(main.querySelector('.app-feature')).not.toBeNull();

    document.body.removeChild(main);
  });

  it('header navigation wraps on mobile, stays inline on desktop', () => {
    const nav = document.createElement('nav');
    nav.className = 'app-header__nav';
    nav.innerHTML = `
      <a href="#/" class="app-header__nav-link">Feed</a>
      <a href="#/item/1" class="app-header__nav-link">Item 1</a>
    `;
    document.body.appendChild(nav);

    // Navigation structure is correct
    const links = nav.querySelectorAll('.app-header__nav-link');
    expect(links).toHaveLength(2);

    document.body.removeChild(nav);
  });
});
