import { createApp, createRouter } from "mini-framework";
import { store } from "./app/state.js";
import { view } from "./app/view.js";

const root = document.querySelector("#app");
if (!root) {
  throw new Error('TodoMVC app: missing root element "#app".');
}

const router = createRouter({
  routes: [{ path: "/" }, { path: "/active" }, { path: "/completed" }],
});

createApp({
  root,
  store,
  view,
  router,
}).mount();
