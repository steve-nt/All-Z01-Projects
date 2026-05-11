import { createStore } from "mini-framework";

// Shared frontend state contract for team integration.
// Keep field names stable so Person 2/3/4 can build on top safely.
export const initialState = {
  todos: [],
  newTodoText: "",
  editingId: null,
  editingText: "",
};

export const store = createStore(initialState);
