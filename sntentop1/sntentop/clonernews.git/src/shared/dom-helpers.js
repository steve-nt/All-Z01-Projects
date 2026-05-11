// DOM helpers centralize safe element creation so feature views do not repeat XSS-prone patterns.
export const createElement = (tagName, options = {}) => {
  const element = document.createElement(tagName);
  const { className, text, attributes = {}, dataset = {}, children = [] } = options;

  // Class assignment stays direct because it is a controlled sink for styling only.
  if (className) {
    element.className = className;
  }

  // Text content is the preferred sink because it cannot execute markup.
  if (text !== undefined) {
    element.textContent = String(text);
  }

  // Attributes are only copied when they have meaningful values so the output stays clean.
  for (const [name, value] of Object.entries(attributes)) {
    if (value === false || value === null || value === undefined) {
      continue;
    }

    element.setAttribute(name, String(value));
  }

  // Dataset values are serialized explicitly so feature code can keep state in data attributes safely.
  for (const [name, value] of Object.entries(dataset)) {
    if (value === false || value === null || value === undefined) {
      continue;
    }

    element.dataset[name] = String(value);
  }

  appendChildren(element, children);

  return element;
};

// Text nodes are created explicitly so callers never need to reach for innerHTML.
export const createText = (value) => document.createTextNode(String(value));

// This helper keeps the mutation surface small by returning the element after updating text.
export const setText = (element, value) => {
  element.textContent = String(value);

  return element;
};

// Children are appended defensively so falsy placeholders do not leak into rendered output.
export const appendChildren = (parent, children) => {
  for (const child of children) {
    if (child === null || child === undefined || child === false) {
      continue;
    }

    parent.append(child.nodeType ? child : document.createTextNode(String(child)));
  }

  return parent;
};

// Clearing through replaceChildren keeps DOM teardown simple and avoids manual node loops.
export const clearElement = (element) => {
  element.replaceChildren();

  return element;
};

// Lookup stays in one helper so shell modules avoid direct document access.
export const getElementById = (id, target = document) => target.getElementById(id);

// First-child lookup is wrapped so shell modules can avoid reading DOM properties directly.
export const getFirstElementChild = (element) => element?.firstElementChild ?? null;

// Clone behavior is centralized so callers avoid direct Node-type checks.
export const cloneDomNode = (node) => (node instanceof Node ? node.cloneNode(true) : null);

// Element guards are wrapped so non-view modules avoid direct HTMLElement references.
export const isHtmlElement = (value) => value instanceof HTMLElement;

// Child replacement is centralized so shell code can swap views without direct DOM mutation calls.
export const replaceElementChildren = (element, ...children) => {
  const renderableChildren = children.filter(
    (child) => child !== null && child !== undefined && child !== false,
  );
  element.replaceChildren(...renderableChildren);

  return element;
};

// Title updates are wrapped so app-shell code can avoid direct document mutation.
export const setDocumentTitle = (value, target = document) => {
  target.title = String(value);

  return target.title;
};
