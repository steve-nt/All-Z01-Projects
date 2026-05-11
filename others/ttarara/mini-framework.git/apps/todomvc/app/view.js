import { h } from "mini-framework";
import {
  addTodoFromInput,
  setNewTodoText,
  toggleTodo,
  deleteTodo,
  toggleAll,
  clearCompleted,
  startEditingTodo,
  setEditingText,
  cancelEditingTodo,
  commitEditingTodo,
} from "./actions.js";

function renderTodoItem(todo, state, store) {
  const isEditing = state.editingId === todo.id;
  const classes = [todo.completed ? "completed" : "", isEditing ? "editing" : ""]
    .filter(Boolean)
    .join(" ");

  return h(
    "li",
    {
      key: todo.id,
      class: classes || undefined,
    },
    h(
      "div",
      // Keys are important so the renderer reuses these DOM nodes
      // and the edit input doesn't lose focus while typing.
      { class: "view", key: "view" },
      h("input", {
        class: "toggle",
        type: "checkbox",
        checked: todo.completed,
        on: {
          change: () => toggleTodo(store, todo.id),
        },
      }),
      h("label", {
        on: {
          dblclick: () => startEditingTodo(store, todo.id),
        },
      }, todo.title),
      h("button", {
        class: "destroy",
        type: "button",
        on: {
          click: () => deleteTodo(store, todo.id),
        },
      }),
    ),
    h("input", {
      key: "edit",
      class: "edit",
      value: isEditing ? state.editingText : todo.title,
      on: {
        input: (event) => setEditingText(store, event.target.value),
        // Use a bubbling focus event so delegated listeners always catch it.
        focusout: (event) => commitEditingTodo(store, todo.id, event.target.value),
        keydown: (event) => {
          const isEnter =
            event.key === "Enter" ||
            event.key === "NumpadEnter" ||
            event.code === "Enter" ||
            event.keyCode === 13 ||
            event.which === 13;
          const isEscape =
            event.key === "Escape" ||
            event.keyCode === 27 ||
            event.which === 27;

          if (isEnter) {
            event.preventDefault();
            // Let the `focusout` handler commit using the live input value.
            event.target.blur();
            return;
          }
          if (isEscape) {
            event.preventDefault();
            cancelEditingTodo(store, todo.id);
          }
        },
      },
    }),
  );
}

function handleNewTodoKeydown(event, store) {
  if (event.key !== "Enter") return;
  event.preventDefault();
  addTodoFromInput(store);
}

export function view({ state, store, route }) {
  const todos = state.todos;

  const routePath = route?.path ?? "/";
  const normalizedRoutePath = routePath === "/active" || routePath === "/completed" ? routePath : "/";

  const filteredTodos =
    normalizedRoutePath === "/active"
      ? todos.filter((t) => !t.completed)
      : normalizedRoutePath === "/completed"
        ? todos.filter((t) => t.completed)
        : todos;

  const remainingCount = todos.reduce((acc, t) => acc + (t.completed ? 0 : 1), 0);
  const completedCount = todos.length - remainingCount;
  const remainingLabel = remainingCount === 1 ? "item" : "items";
  const filter = normalizedRoutePath === "/active" ? "active" : normalizedRoutePath === "/completed" ? "completed" : "all";

  return h(
    "section",
    // Stable key keeps the root node reusable across rerenders.
    { class: "todoapp", key: "todoapp-root" },
    h(
      "header",
      { class: "header", key: "header" },
      h("h1", { key: "title" }, "todos"),
      h("input", {
        // Keep this key stable so input DOM stays consistent on rerenders.
        key: "new-todo-input",
        class: "new-todo",
        placeholder: "What needs to be done?",
        autofocus: true,
        value: state.newTodoText,
        on: {
          input: (event) => {
            setNewTodoText(store, event.target.value);
          },
          keydown: (event) => handleNewTodoKeydown(event, store),
        },
      }),
    ),
    // Main (toggle-all + list) and footer only when there are todos
    ...(todos.length > 0
      ? [
          h(
            "section",
            { class: "main", key: "main" },
            h("input", {
              key: "toggle-all",
              id: "toggle-all",
              class: "toggle-all",
              type: "checkbox",
              checked: todos.every((t) => t.completed),
              on: {
                change: () => toggleAll(store),
              },
            }),
            h("label", {
              key: "toggle-all-label",
              for: "toggle-all",
            }, "Mark all as complete"),
            h(
              "ul",
              { class: "todo-list", key: "todo-list" },
              filteredTodos.map((todo) => renderTodoItem(todo, state, store)),
            ),
          ),
          h(
            "footer",
            { class: "footer", key: "footer" },
            h(
              "span",
              { class: "todo-count", key: "todo-count" },
              h("strong", { key: "todo-count-strong" }, String(remainingCount)),
              " ",
              remainingLabel,
              " left",
            ),
            h(
              "ul",
              { class: "filters", key: "filters" },
              h(
                "li",
                { key: "filter-all" },
                h("a", {
                  key: "filter-all-link",
                  href: "#/",
                  class: filter === "all" ? "selected" : undefined,
                }, "All"),
              ),
              h(
                "li",
                { key: "filter-active" },
                h("a", {
                  key: "filter-active-link",
                  href: "#/active",
                  class: filter === "active" ? "selected" : undefined,
                }, "Active"),
              ),
              h(
                "li",
                { key: "filter-completed" },
                h("a", {
                  key: "filter-completed-link",
                  href: "#/completed",
                  class: filter === "completed" ? "selected" : undefined,
                }, "Completed"),
              ),
            ),
            completedCount > 0
              ? h(
                  "button",
                  {
                    key: "clear-completed",
                    class: "clear-completed",
                    type: "button",
                    on: { click: () => clearCompleted(store) },
                  },
                  "Clear completed",
                )
              : null,
          ),
        ]
      : []),
  );
}
