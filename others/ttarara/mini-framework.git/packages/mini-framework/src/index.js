// Mini-Framework public API (v0).
// Person 2: VNode contract + renderer (mount/patch) + keyed lists.

// --- A. Events contract ---
//
// Declarative format: props.on = { click, input, submit, keydown, change, ... }
// - Keys are DOM event type names (e.g. "click", "input", "submit", "keydown", "change").
// - Values are handler functions. Handler signature: (event: Event) => void.
// - The framework passes the native Event object through to your handler.
// - Events are handled via delegation (one listener per event type at app root);
//   do not call addEventListener in app code.

// --- B. VNode creation contract ---
//
// Final VNode shape returned by h(tag, props, ...children):
//   { tag: string, props: object, key?: string | number, children: VNode[] }
//
// - tag: element tag name ("div", "input", ...) or "#text" for text nodes.
// - props: attributes, properties, and framework features (e.g. on, class, style).
//   The "key" prop is lifted to VNode.key and not applied to the DOM.
// - key: optional; used for list reconciliation (see keyed lists below).
// - children: normalized array of child VNodes (elements or #text).
//
// Children normalization rules:
// - null, undefined, false are ignored.
// - Arrays are flattened recursively ([a, [b, c]] -> a, b, c).
// - Strings and numbers become text VNodes ({ tag: "#text", props: { nodeValue }, children: [] }).
//
// Key behavior for lists:
// - Pass key in props: h("li", { key: item.id }, ...). The renderer uses key to match
//   old and new list items across updates so DOM nodes are reused and reordered instead
//   of recreated, preserving focus and form state (e.g. Todo edit mode).
// - Keys should be unique among siblings. Missing keys fall back to index-based match.

/**
 * @typedef {Object} VNode
 * @property {string} tag - Tag name or "#text"
 * @property {Object} props - Attributes, properties, style, on (events)
 * @property {string|number} [key] - Optional key for list reconciliation
 * @property {VNode[]} children - Normalized child VNodes
 */

/**
 * Create a VNode. Key from props.key is lifted to VNode.key and not set on the DOM.
 * @param {string} tag
 * @param {Object | null} props
 * @param {...any} children
 * @returns {VNode}
 */
export function h(tag, props, ...children) {
  const p = props ?? {};
  const vnode = {
    tag,
    props: p,
    children: normalizeChildren(children),
  };
  if (p.key !== undefined && p.key !== null) {
    vnode.key = p.key;
  }
  return vnode;
}

/**
 * Normalize children: ignore null/undefined/false, flatten arrays, coerce primitives to text VNodes.
 * @param {any[]} children
 * @returns {VNode[]}
 */
function normalizeChildren(children) {
  const out = [];
  const push = (v) => {
    if (v === null || v === undefined || v === false) return;
    if (Array.isArray(v)) {
      v.forEach(push);
      return;
    }
    if (typeof v === "string" || typeof v === "number") {
      out.push({ tag: "#text", props: { nodeValue: String(v) }, children: [] });
      return;
    }
    out.push(v);
  };
  children.forEach(push);
  return out;
}

export { createStore } from "./store.js";

/**
 * @typedef {Object} Route
 * @property {string} path
 * @property {string} [name]
 */

/**
 * @typedef {Object} RouteMatch
 * @property {string} path
 * @property {Record<string,string>} params
 * @property {URLSearchParams} query
 * @property {string} hash
 * @property {boolean} [isFallback]
 */

/**
 * Small hash-based router.
 *
 * Contract:
 * - mode: "hash" only (for now)
 * - routes: array of { path, name? }
 * - getLocation(): RouteMatch for current URL
 * - match(path): RouteMatch or null
 * - navigate(pathOrDescriptor): update hash (e.g. "/active" -> "#/active")
 * - subscribe(fn): listen to location changes, returns unsubscribe()
 *
 * @param {{ routes?: Route[]; mode?: "hash" }} options
 */
export function createRouter({ routes = [], mode = "hash" } = {}) {
  if (mode !== "hash") {
    throw new Error(
      `[mini-framework] createRouter({ mode }): only "hash" mode is supported right now.`,
    );
  }

  /** @type {Route[]} */
  const routeTable = Array.isArray(routes) ? routes.slice() : [];

  /** @type {Set<(match: RouteMatch) => void>} */
  const listeners = new Set();

  function normalizeHash(rawHash) {
    const h = rawHash || "";
    const trimmed = h.startsWith("#") ? h.slice(1) : h;
    if (!trimmed || trimmed === "/") return "/";
    return trimmed.startsWith("/") ? trimmed : "/" + trimmed;
  }

  function getDefaultRoute() {
    if (routeTable.length === 0) return null;
    const wildcard = routeTable.find((r) => r.path === "*" || r.path === "/*");
    return wildcard ?? routeTable[0];
  }

  /**
   * @returns {RouteMatch}
   */
  function getLocation() {
    const path = normalizeHash(window.location.hash);
    const query = new URLSearchParams(window.location.search ?? "");
    const hash = window.location.hash || "";
    // Avoid shadowing the `match()` function (TDZ error).
    const matched = match(path);
    if (matched) return matched;

    const fallbackRoute = getDefaultRoute();
    if (fallbackRoute) {
      return {
        path: fallbackRoute.path,
        params: {},
        query,
        hash,
        isFallback: true,
      };
    }

    return {
      path,
      params: {},
      query,
      hash,
      isFallback: true,
    };
  }

  /**
   * @param {string} path
   * @returns {RouteMatch | null}
   */
  function match(path) {
    const normalized = normalizeHash(path);
    const query = new URLSearchParams(window.location.search ?? "");
    const hash = window.location.hash || "";

    const route = routeTable.find((r) => r.path === normalized);
    if (!route) return null;

    return {
      path: route.path,
      params: {},
      query,
      hash,
    };
  }

  function toHashPath(path) {
    const normalized = normalizeHash(path);
    return normalized === "/" ? "#/" : "#" + normalized;
  }

  function navigate(to) {
    if (typeof to === "string") {
      window.location.hash = toHashPath(to);
      return;
    }
    if (to && typeof to === "object" && to.name) {
      const route = routeTable.find((r) => r.name === to.name);
      if (!route) {
        console.warn(
          `[mini-framework] router.navigate({ name: "${to.name}" }): no route found with that name.`,
        );
        return;
      }
      window.location.hash = toHashPath(route.path);
      return;
    }
    console.warn(
      "[mini-framework] router.navigate(to): expected string path or { name } descriptor.",
    );
  }

  function emit() {
    const location = getLocation();
    // snapshot listeners for safe mutation during emit
    [...listeners].forEach((fn) => {
      try {
        fn(location);
      } catch (err) {
        console.error("[mini-framework] router subscriber threw:", err);
      }
    });
  }

  // Single hashchange listener per router instance.
  const hashListener = () => emit();
  window.addEventListener("hashchange", hashListener);

  function subscribe(fn) {
    if (typeof fn !== "function") {
      throw new TypeError(
        "router.subscribe(listener): listener must be a function.",
      );
    }
    listeners.add(fn);
    // Immediately notify with current location so app can render initial route.
    fn(getLocation());
    return () => {
      listeners.delete(fn);
    };
  }

  return {
    navigate,
    getLocation,
    match,
    subscribe,
    destroy() {
      window.removeEventListener("hashchange", hashListener);
      listeners.clear();
    },
  };
}

export function createApp({ root, view, store, router }) {
  if (!root) {
    throw new Error(
      "[mini-framework] createApp({ root }): root is required (HTMLElement or selector string).",
    );
  }
  if (typeof view !== "function") {
    throw new Error(
      "[mini-framework] createApp({ view }): view must be a function that returns a VNode.",
    );
  }

  /** @type {Element | null} */
  let rootEl = null;
  /** @type {RenderContext | null} */
  let renderContext = null;

  let unsubscribeStore = null;
  let unsubscribeRouter = null;
  let isMounted = false;
  /** @type {VNode | null} */
  let lastRootVNode = null;
  /** @type {Element | Text | null} */
  let lastRootDOM = null;

  let isRendering = false;
  let isRenderScheduled = false;
  let frameId = null;

  /** @type {RenderContext} */
  function ensureRenderContext() {
    if (!rootEl) {
      throw new Error(
        "[mini-framework] internal error: root element missing during render.",
      );
    }
    if (!renderContext) {
      renderContext = {
        root: rootEl,
        handlerMap: new WeakMap(),
        eventTypes: new Set(),
        delegatedListeners: new Map(),
      };
    }
    return renderContext;
  }

  const ctx = {
    get state() {
      return store?.getState?.();
    },
    get route() {
      return router?.getLocation?.();
    },
    store,
    router,
  };

  function renderOnce() {
    if (isRendering) {
      console.warn(
        "[mini-framework] render was re-entered; ignoring nested render call.",
      );
      return;
    }
    isRendering = true;
    try {
      const vnode = view(ctx);
      if (!vnode || typeof vnode !== "object" || !vnode.tag) {
        throw new Error(
          "[mini-framework] view(ctx) must return a valid VNode (use h()).",
        );
      }
      const context = ensureRenderContext();
      if (lastRootVNode === null) {
        lastRootVNode = vnode;
        lastRootDOM = mount(vnode, context);
        rootEl.replaceChildren(lastRootDOM);
      } else {
        patch(lastRootDOM, lastRootVNode, vnode, context);
        lastRootVNode = vnode;
      }
    } finally {
      isRendering = false;
    }
  }

  function scheduleRender() {
    if (!isMounted) return;
    if (isRenderScheduled) return;
    isRenderScheduled = true;
    frameId = requestAnimationFrame(() => {
      isRenderScheduled = false;
      frameId = null;
      renderOnce();
    });
  }

  function removeDelegatedListeners() {
    if (!renderContext || !rootEl) return;
    for (const [eventType, fn] of renderContext.delegatedListeners) {
      rootEl.removeEventListener(eventType, fn);
    }
    renderContext.delegatedListeners.clear();
    renderContext.eventTypes.clear();
  }

  function resolveRootElement() {
    if (rootEl) return rootEl;
    if (typeof root === "string") {
      const el = document.querySelector(root);
      if (!el) {
        throw new Error(
          `[mini-framework] createApp: root selector "${root}" did not match any element.`,
        );
      }
      rootEl = el;
    } else if (root && root.nodeType === 1) {
      rootEl = root;
    } else {
      throw new Error(
        "[mini-framework] createApp({ root }): root must be an Element or selector string.",
      );
    }
    return rootEl;
  }

  return {
    mount() {
      if (isMounted) return;
      resolveRootElement();
      isMounted = true;

      renderOnce();

      if (store?.subscribe) {
        unsubscribeStore = store.subscribe(scheduleRender);
      }
      if (router?.subscribe) {
        let sawInitial = false;
        unsubscribeRouter = router.subscribe(() => {
          if (!sawInitial) {
            // Skip the immediate first emission; we've already rendered with current route.
            sawInitial = true;
            return;
          }
          scheduleRender();
        });
      }
    },
    unmount() {
      if (!isMounted) return;
      isMounted = false;

      if (frameId != null) {
        cancelAnimationFrame(frameId);
        frameId = null;
      }
      isRenderScheduled = false;

      if (unsubscribeStore) unsubscribeStore();
      if (unsubscribeRouter) unsubscribeRouter();
      unsubscribeStore = null;
      unsubscribeRouter = null;

      removeDelegatedListeners();
      lastRootVNode = null;
      lastRootDOM = null;

      if (rootEl) {
        rootEl.replaceChildren(); // Remove all renderer-managed nodes; handlerMap keys are GC'd.
      }
    },
  };
}

// --- C. DOM abstraction + renderer ---
//
// Mount: VNode -> real DOM (elements + text).
// Props: attributes (id, class, data-*, aria-*), form control properties (value, checked, etc.),
//       style as object { color: "red" } (documented in README).
// Events: delegated at app root; handler map (DOM element -> { eventType: handler }) drives dispatch.
// Update strategy: diff/patch with keyed list reconciliation for stable Todo items.
// Unmount: root.replaceChildren() + remove delegated listeners + clear handler map refs (no leaks).

/**
 * @typedef {Object} RenderContext
 * @property {Element} root - App root element (delegation target)
 * @property {WeakMap<Element, Object<string, function(Event): void>>} handlerMap - element -> { eventType: handler }
 * @property {Set<string>} eventTypes - event types seen so far
 * @property {Map<string, function(Event): void>} delegatedListeners - eventType -> root listener (for cleanup)
 */

const FORM_CONTROL_PROPS = new Set([
  "value",
  "checked",
  "disabled",
  "selected",
  "indeterminate",
  "readOnly",
  "multiple",
]);

/**
 * Create a delegated listener for one event type. Walks from event.target up to context.root,
 * finds the first element with a handler in handlerMap, and invokes it with the native Event.
 * @param {string} eventType
 * @param {RenderContext} context
 * @returns {function(Event): void}
 */
function createDelegate(eventType, context) {
  return function delegatedListener(e) {
    for (let node = e.target; node && node !== context.root; node = node.parentElement) {
      const handlers = context.handlerMap.get(node);
      if (handlers && typeof handlers[eventType] === "function") {
        handlers[eventType](e);
        break;
      }
    }
  };
}

/**
 * Ensure a delegated listener for this event type is attached to context.root (lazy, once per type).
 * @param {string} eventType
 * @param {RenderContext} context
 */
function ensureDelegatedListener(eventType, context) {
  if (context.delegatedListeners.has(eventType)) return;
  const fn = createDelegate(eventType, context);
  context.root.addEventListener(eventType, fn);
  context.delegatedListeners.set(eventType, fn);
}

/**
 * Apply prop changes to a DOM element. If context is provided, props.on is stored in
 * context.handlerMap (delegation) instead of adding per-element listeners, so rerenders
 * do not double-fire. Old props are used to remove stale classes/attributes/styles.
 * @param {Element} el
 * @param {Object} oldProps
 * @param {Object} newProps
 * @param {RenderContext} [context]
 */
function applyProps(el, oldProps, newProps, context) {
  const prev = oldProps ?? {};
  const next = newProps ?? {};

  if (context) {
    context.handlerMap.set(el, {}); // clear handlers; "on" branch below will overwrite if present
  }

  for (const [k, oldValue] of Object.entries(prev)) {
    if (k === "key") continue;
    if (k in next && next[k] !== null && next[k] !== undefined && next[k] !== false) {
      continue;
    }

    if (k === "on") {
      if (context) {
        context.handlerMap.set(el, {});
      }
      continue;
    }

    if (k === "class" || k === "className") {
      el.className = "";
      el.removeAttribute("class");
      continue;
    }

    if (k === "style") {
      if (oldValue && typeof oldValue === "object" && !Array.isArray(oldValue)) {
        for (const cssProp of Object.keys(oldValue)) {
          el.style[cssProp] = "";
        }
      } else {
        el.removeAttribute("style");
      }
      continue;
    }

    if (FORM_CONTROL_PROPS.has(k) && k in el) {
      if (k === "value") {
        el[k] = "";
      } else {
        el[k] = false;
      }
      el.removeAttribute(k);
      continue;
    }

    if (typeof oldValue === "boolean" && k in el) {
      el[k] = false;
      el.removeAttribute(k);
      continue;
    }

    if (k in el) {
      try {
        el[k] = "";
      } catch (_) {}
    }
    el.removeAttribute(k);
  }

  for (const [k, v] of Object.entries(next)) {
    if (k === "key") continue;
    if (v === null || v === undefined || v === false) continue;
    if (k === "on") {
      const handlers = {};
      if (v && typeof v === "object") {
        for (const [eventName, handler] of Object.entries(v)) {
          if (typeof handler === "function") {
            handlers[eventName] = handler;
            if (context) {
              context.eventTypes.add(eventName);
              ensureDelegatedListener(eventName, context);
            }
          }
        }
      }
      if (context) {
        context.handlerMap.set(el, handlers); // always update (or clear) so patch reflects current VNode
      } else if (Object.keys(handlers).length) {
        for (const [eventName, handler] of Object.entries(handlers)) {
          el.addEventListener(eventName, handler);
        }
      }
      continue;
    }
    if (k === "class" || k === "className") {
      el.className = String(v);
      continue;
    }
    if (k === "style" && v && typeof v === "object" && !Array.isArray(v)) {
      const prevStyle = prev.style;
      if (prevStyle && typeof prevStyle === "object" && !Array.isArray(prevStyle)) {
        for (const cssProp of Object.keys(prevStyle)) {
          if (!(cssProp in v)) {
            el.style[cssProp] = "";
          }
        }
      } else if (typeof prevStyle === "string") {
        el.removeAttribute("style");
      }
      for (const [cssProp, cssVal] of Object.entries(v)) {
        if (cssVal != null && cssVal !== "") {
          el.style[cssProp] = cssVal;
        }
      }
      continue;
    }
    if (k === "style" && typeof v === "string") {
      el.setAttribute("style", v);
      continue;
    }
    if (FORM_CONTROL_PROPS.has(k) && k in el) {
      try {
        el[k] = v;
        continue;
      } catch (_) {}
    }
    if (typeof v === "boolean" && k in el) {
      el[k] = v;
      continue;
    }
    if (k in el) {
      try {
        el[k] = v;
        continue;
      } catch (_) {}
    }
    el.setAttribute(k, String(v));
  }
}

/**
 * Mount a single VNode to a new DOM node (element or text). No reuse.
 * @param {VNode} vnode
 * @param {RenderContext} [context]
 * @returns {Element | Text}
 */
function mount(vnode, context) {
  if (!vnode || typeof vnode !== "object") {
    return document.createTextNode("");
  }
  if (vnode.tag === "#text") {
    return document.createTextNode(vnode.props?.nodeValue ?? "");
  }
  const el = document.createElement(vnode.tag);
  applyProps(el, {}, vnode.props, context);
  for (const child of vnode.children ?? []) {
    el.appendChild(mount(child, context));
  }
  return el;
}

/**
 * Patch existing DOM to match new VNode. Reuses DOM when tag and key match.
 * @param {Element | Text} el - Existing DOM node (from previous mount/patch)
 * @param {VNode} oldVNode
 * @param {VNode} newVNode
 * @param {RenderContext} [context]
 */
function patch(el, oldVNode, newVNode, context) {
  if (!newVNode || typeof newVNode !== "object") {
    const parent = el.parentNode;
    if (parent) parent.replaceChild(document.createTextNode(""), el);
    return;
  }
  if (newVNode.tag === "#text") {
    const text = document.createTextNode(newVNode.props?.nodeValue ?? "");
    if (el.parentNode) el.parentNode.replaceChild(text, el);
    return;
  }
  if (oldVNode?.tag === "#text" || el.nodeType !== 1) {
    const newEl = mount(newVNode, context);
    if (el.parentNode) el.parentNode.replaceChild(newEl, el);
    return;
  }
  const domEl = /** @type {Element} */ (el);
  if (oldVNode?.tag !== newVNode.tag) {
    const newEl = mount(newVNode, context);
    if (domEl.parentNode) domEl.parentNode.replaceChild(newEl, domEl);
    return;
  }
  applyProps(domEl, oldVNode?.props, newVNode.props, context);
  patchChildren(domEl, oldVNode?.children ?? [], newVNode.children ?? [], context);
}

/**
 * Keyed list reconciliation: match old and new children by key (or index), patch in place, reorder.
 * Preserves DOM nodes for matching keys so focus and form state (e.g. Todo edit) are stable.
 * @param {Element} parentEl
 * @param {VNode[]} oldCh
 * @param {VNode[]} newCh
 * @param {RenderContext} [context]
 */
function patchChildren(parentEl, oldCh, newCh, context) {
  const oldLen = oldCh.length;
  const newLen = newCh.length;
  const currentChildren = Array.from(parentEl.childNodes);

  if (newLen === 0) {
    parentEl.replaceChildren();
    return;
  }

  const oldKeyToIdx = new Map();
  for (let i = 0; i < oldLen; i++) {
    const v = oldCh[i];
    if (v?.key != null) oldKeyToIdx.set(v.key, i);
  }

  const matched = [];
  for (let i = 0; i < newLen; i++) {
    const n = newCh[i];
    if (n?.key != null && oldKeyToIdx.has(n.key)) {
      matched.push({ newIdx: i, oldIdx: oldKeyToIdx.get(n.key) });
      oldKeyToIdx.delete(n.key);
    }
  }

  const newDomList = [];
  for (let i = 0; i < newLen; i++) {
    const m = matched.find((x) => x.newIdx === i);
    if (m !== undefined) {
      const oldIdx = m.oldIdx;
      const dom = currentChildren[oldIdx];
      patch(dom, oldCh[oldIdx], newCh[i], context);
      newDomList.push(dom);
    } else {
      newDomList.push(mount(newCh[i], context));
    }
  }

  for (let i = 0; i < newDomList.length; i++) {
    const desired = newDomList[i];
    const currentAtI = parentEl.childNodes[i];
    if (currentAtI !== desired) {
      parentEl.insertBefore(desired, currentAtI ?? null);
    }
  }
  while (parentEl.childNodes.length > newDomList.length) {
    parentEl.removeChild(parentEl.lastChild);
  }
}
