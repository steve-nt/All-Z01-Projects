let nextTodoId = 1;

// Simple incremental id generator for this client-side session.
function createTodoId() {
  const id = nextTodoId;
  nextTodoId += 1;
  return id;
}

export function setNewTodoText(store, value) {
  store.setState({ newTodoText: value });
}

export function addTodoFromInput(store) {
  // Trim so whitespace-only entries are treated as empty input.
  const title = store.getState().newTodoText.trim();
  if (!title) {
    return false;
  }

  const newTodo = {
    id: createTodoId(),
    title,
    completed: false,
  };

  // Immutable update: append todo and clear the input field.
  store.setState((prev) => ({
    todos: [...prev.todos, newTodo],
    newTodoText: "",
  }));

  return true;
}

export function toggleTodo(store, id) {
  store.setState((prev) => ({
    todos: prev.todos.map((t) =>
      t.id === id ? { ...t, completed: !t.completed } : t
    ),
  }));
}

export function deleteTodo(store, id) {
  store.setState((prev) => ({
    todos: prev.todos.filter((t) => t.id !== id),
    editingId: prev.editingId === id ? null : prev.editingId,
    editingText: prev.editingId === id ? "" : prev.editingText,
  }));
}

export function toggleAll(store) {
  const { todos } = store.getState();
  if (todos.length === 0) return;
  const allCompleted = todos.every((t) => t.completed);
  store.setState((prev) => ({
    todos: prev.todos.map((t) => ({ ...t, completed: !allCompleted })),
  }));
}

export function clearCompleted(store) {
  const { todos, editingId } = store.getState();
  const completedIds = new Set(todos.filter((t) => t.completed).map((t) => t.id));
  if (completedIds.size === 0) return;

  store.setState((prev) => {
    const editingTodo = prev.editingId != null
      ? prev.todos.find((t) => t.id === prev.editingId)
      : null;

    const shouldClearEditing = editingTodo?.completed === true;

    return {
      todos: prev.todos.filter((t) => !t.completed),
      editingId: shouldClearEditing ? null : prev.editingId,
      editingText: shouldClearEditing ? "" : prev.editingText,
    };
  });
}

// Tracks Escape-cancel for a specific todo so the subsequent blur does not save.
let skipBlurCommitForId = null;

export function startEditingTodo(store, id) {
  const todo = store.getState().todos.find((t) => t.id === id);
  if (!todo) return;

  store.setState({
    editingId: id,
    editingText: todo.title,
  });

  // Focus the newly shown edit input after the render cycle.
  requestAnimationFrame(() => {
    const input = document.querySelector("li.editing .edit");
    if (!input) return;
    input.focus();
    const len = input.value.length;
    input.setSelectionRange(len, len);
  });
}

export function setEditingText(store, value) {
  store.setState({ editingText: value });
}

export function cancelEditingTodo(store, id) {
  skipBlurCommitForId = id;
  store.setState({
    editingId: null,
    editingText: "",
  });
}

export function commitEditingTodo(store, id) {
  if (skipBlurCommitForId === id) {
    skipBlurCommitForId = null;
    return;
  }

  const { editingId, editingText } = store.getState();
  if (editingId !== id) return;

  const title = editingText.trim();
  if (!title) {
    deleteTodo(store, id);
    return;
  }

  store.setState((prev) => ({
    todos: prev.todos.map((t) => (t.id === id ? { ...t, title } : t)),
    editingId: null,
    editingText: "",
  }));
}
