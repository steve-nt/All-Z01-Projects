// @vitest-environment jsdom

// DOM helper tests stay in jsdom because they exercise safe browser sink behavior directly.
// Public API under test: createElement, createText, appendChildren, setText, and clearElement.
// Constraints: tests verify safe DOM sinks and do not use unsafe HTML insertion APIs.
import { describe, expect, it } from 'vitest';

import {
  appendChildren,
  clearElement,
  createElement,
  createText,
  setText,
} from '../../../src/shared/dom-helpers.js';

describe('dom helpers', () => {
  it('creates elements using safe sinks', () => {
    // This assertion guards the safe-sink contract that later feature views rely on.
    const element = createElement('a', {
      className: 'link',
      text: 'Read more',
      attributes: {
        href: 'https://example.com',
        target: '_blank',
      },
      dataset: {
        testid: 'story-link',
      },
    });

    expect(element.tagName).toBe('A');
    expect(element.className).toBe('link');
    expect(element.textContent).toBe('Read more');
    expect(element.getAttribute('href')).toBe('https://example.com');
    expect(element.dataset.testid).toBe('story-link');
  });

  it('appends mixed children and clears content', () => {
    // Mixed child types ensure the helper handles text, nodes, and cleanup in one pass.
    const element = document.createElement('div');
    const span = document.createElement('span');
    span.textContent = 'World';

    appendChildren(element, ['Hello ', createText('there '), span]);

    expect(element.textContent).toBe('Hello there World');

    setText(element, 'Reset');
    expect(element.textContent).toBe('Reset');

    clearElement(element);
    expect(element.childNodes).toHaveLength(0);
  });
});
