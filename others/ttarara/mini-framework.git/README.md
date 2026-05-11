# Mini Framework Documentation

## Overview

`mini-framework` is a small JavaScript framework used to build UI with:

- virtual nodes instead of manual DOM code
- a global store for state
- a small hash router
- declarative events

You write a `view(ctx)` function that returns UI with `h(...)`.  
The framework renders that UI, listens for state and route changes, and updates the DOM.

This is the public API used by the TodoMVC app in this repository.

The TodoMVC application in this repository is used to demonstrate and validate the framework features during the audit.

## How to Run the Project

This project runs with **Vite**.

From the repository root, install dependencies:

```bash
npm install
```

Start the development server:

```bash
npm run dev
```

What you should expect:

- Vite starts a local development server
- a local URL appears in the terminal, usually something like `http://localhost:5173/`
- open that URL in your browser
- you should see the TodoMVC app built with this framework

## Features

### DOM abstraction
You create UI with JavaScript objects called **VNodes** instead of writing `document.createElement(...)` everywhere.

### Routing system
The framework includes a small **hash router**. It is useful for URLs like:

- `#/`
- `#/active`
- `#/completed`

### State management
The framework includes a small **store** with:

- `getState()`
- `setState(...)`
- `subscribe(...)`

### Event handling
Events are declared inside VNode props with `on: { ... }`.  
The framework handles them through **event delegation** at the app root.

## Getting Started

Basic app example:

```js
import { createApp, createStore, h } from "mini-framework";

const store = createStore({ count: 0 });

function view({ state, store }) {
  return h(
    "div",
    { class: "card" },
    h("h1", null, "Counter"),
    h("p", null, "Count: ", String(state.count)),
    h(
      "button",
      {
        on: {
          click: () => store.setState((prev) => ({ count: prev.count + 1 })),
        },
      },
      "Increment",
    ),
  );
}

createApp({
  root: document.querySelector("#app"),
  store,
  view,
}).mount();
```

What happens here:

- `createStore(...)` creates the app state
- `view(...)` returns VNodes
- `createApp(...)` connects the root, store, and view
- `mount()` renders the app

## Creating Elements

Use `h(tag, props, ...children)` to create a VNode.

```js
import { h } from "mini-framework";

const title = h("h1", null, "Hello");
```

Example with a `div`:

```js
const box = h("div", { class: "box" }, "This is a box");
```

### `h(...)` parameters

- `tag`: HTML tag name like `"div"`, `"button"`, `"input"`
- `props`: attributes, properties, events, styles
- `children`: nested elements, text, numbers, or arrays

### Child rules

- `null`, `undefined`, and `false` are ignored
- arrays are flattened
- strings and numbers become text nodes

Example:

```js
const node = h(
  "div",
  null,
  "Hello ",
  123,
  [h("span", null, "world")],
);
```

## Attributes and Props

Pass attributes and properties in the second argument of `h(...)`.

```js
const input = h("input", {
  id: "new-todo",
  class: "new-todo",
  type: "text",
  placeholder: "What needs to be done?",
});
```

### Common attributes

```js
const link = h("a", {
  href: "#/active",
  class: "selected",
  title: "Show active todos",
}, "Active");
```

### Form properties

Some props are applied as DOM properties so inputs behave correctly:

- `value`
- `checked`
- `disabled`
- `selected`
- `readOnly`
- `multiple`

Example:

```js
const checkbox = h("input", {
  type: "checkbox",
  checked: true,
});
```

### Style

Style can be an object:

```js
const card = h("div", {
  style: {
    color: "red",
    marginTop: "8px",
  },
}, "Styled text");
```

Or a string:

```js
const card = h("div", {
  style: "color: red; margin-top: 8px;",
}, "Styled text");
```

## Events

Events are declared with `props.on`.

```js
const button = h(
  "button",
  {
    on: {
      click: (event) => {
        console.log("clicked", event.target);
      },
    },
  },
  "Click me",
);
```

You can use normal DOM event names, for example:

- `click`
- `input`
- `change`
- `submit`
- `keydown`
- `dblclick`
- `focusout`

Example with input:

```js
const field = h("input", {
  value: "hello",
  on: {
    input: (event) => {
      console.log(event.target.value);
    },
  },
});
```

Important:

- app code should use `on: { ... }`
- app code should not manually attach listeners everywhere with `addEventListener(...)`
- the framework passes the native browser event object to your handler

## Nesting Elements

You can nest elements by passing children inside `h(...)`.

```js
const page = h(
  "div",
  { class: "container" },
  h("h1", null, "Todo App"),
  h("p", null, "This is a small example."),
  h(
    "button",
    {
      on: {
        click: () => console.log("save"),
      },
    },
    "Save",
  ),
);
```

You can also render lists:

```js
const items = ["A", "B", "C"];

const list = h(
  "ul",
  null,
  items.map((item, index) =>
    h("li", { key: index }, item),
  ),
);
```

Keys are important in lists because they help the renderer reuse the correct DOM nodes.

## State Management

Use `createStore(initialState)` to create shared state.

```js
import { createStore } from "mini-framework";

const store = createStore({
  count: 0,
  name: "Andy",
});
```

### `getState()`

Returns the current state:

```js
const state = store.getState();
console.log(state.count);
```

### `setState(...)`

You can update state with an object:

```js
store.setState({ count: 1 });
```

Or with an updater function:

```js
store.setState((prev) => ({
  count: prev.count + 1,
}));
```

### `subscribe(...)`

You can listen for changes:

```js
const unsubscribe = store.subscribe(() => {
  console.log("new state:", store.getState());
});
```

### Store behavior

- if previous state and next state are plain objects, the store does a **shallow merge**
- otherwise, the state is replaced
- state should be treated as immutable

In a framework app, store changes trigger a re-render through `createApp(...)`.

## Routing

Use `createRouter(...)` for hash-based routing.

```js
import { createRouter } from "mini-framework";

const router = createRouter({
  routes: [
    { path: "/" },
    { path: "/active" },
    { path: "/completed" },
  ],
});
```

### Route methods

#### `getLocation()`

Returns the current route information:

```js
const location = router.getLocation();
console.log(location.path);
```

#### `navigate(path)`

Changes the hash URL:

```js
router.navigate("/active");
```

This updates the URL to:

```txt
#/active
```

#### `subscribe(fn)`

Runs a callback when the route changes:

```js
router.subscribe((location) => {
  console.log("route changed:", location.path);
});
```

### Typical routing use

In a view, you can read the current route and render different content:

```js
function view({ route }) {
  const currentPath = route?.path ?? "/";

  return h(
    "div",
    null,
    h("p", null, "Current route: ", currentPath),
  );
}
```

## How It Works Internally

This section explains the main ideas simply.

### View -> VNodes

`h(...)` does not create a real DOM element immediately.  
It creates a plain JavaScript object called a **VNode**.

Example shape:

```js
{
  tag: "div",
  props: { class: "card" },
  children: [
    { tag: "#text", props: { nodeValue: "hello" }, children: [] }
  ]
}
```

The framework reads this object tree and uses it to build or update the real DOM.

### VNodes -> DOM

On the first render:

- the framework calls your `view(ctx)`
- it gets a VNode tree
- it creates real DOM nodes from that tree
- it inserts them into the root element

This first render is called **mounting**.

### State change -> re-render

When you use `createApp(...)`, the app subscribes to the store.

That means:

- an action calls `store.setState(...)`
- the store notifies subscribers
- the framework schedules a render
- `view(ctx)` runs again with the new state
- the new VNode tree is compared with the old one
- only the changed DOM is patched

This is why the UI stays in sync with the data.

After the first render, the framework does not rebuild everything from scratch.  
It patches the existing DOM.

For lists, keys help the framework keep the correct DOM nodes when items move or update.

### Router -> view update

When you pass a router into `createApp(...)`, the app also listens for route changes.

That means:

- the URL hash changes
- the router notifies subscribers
- the framework re-renders
- your `view({ route, state, ... })` can show the correct screen or filter

This is how the TodoMVC filters are connected to the URL.

### Event delegation

Events use **delegation**.

That means the framework:

- stores handlers from `props.on`
- attaches one root listener per event type
- catches bubbling events at the app root
- finds the correct element handler
- runs your function with the native event

This avoids attaching lots of separate listeners during every render.

## API Summary

### `h(tag, props, ...children)`
Creates a VNode.

### `createStore(initialState)`
Creates a store with:

- `getState()`
- `setState(next)`
- `subscribe(listener)`

### `createRouter({ routes, mode })`
Creates a hash router with:

- `navigate(path)`
- `getLocation()`
- `match(path)`
- `subscribe(listener)`

### `createApp({ root, view, store, router })`
Creates an app instance with:

- `mount()`
- `unmount()`

## Where To Look In This Repo

If you want to study the implementation:

- framework source: `packages/mini-framework/src/index.js`
- store source: `packages/mini-framework/src/store.js`
- TodoMVC app entry: `apps/todomvc/main.js`
- TodoMVC view: `apps/todomvc/app/view.js`
- TodoMVC actions: `apps/todomvc/app/actions.js`

## Quick Manual Test (TodoMVC)

Use these steps for a fast audit test:

1. Add 2 todos.
   Expected: both appear in the list and the footer becomes visible.

2. Check one todo, then uncheck it.
   Expected: the item gets the completed style, then returns to normal.

3. Click `Active`, `Completed`, and `All`.
   Expected: the URL hash changes, the correct items are shown, and only one filter is selected at a time.

4. Double-click a todo label and edit it.
   Expected: the item enters edit mode, saves correctly, and leaves edit mode after saving.

5. Click `Clear completed`.
   Expected: only completed todos are removed.

6. Delete the remaining todos.
   Expected: the list becomes empty and the main section/footer disappear.


## Authors

📝✅ Todo crew:

- 👩‍💻 [Charoula Tarara](https://discordapp.com/users/1242540766879023160)
- 📋 [Georgia Marouli](https://discordapp.com/users/1277216244910522371)
- ✍️ [Andriana Stas](https://discordapp.com/users/780150798927134740)
- 🗂️ [Iana Kopylova](https://discordapp.com/users/1279339146833297509)