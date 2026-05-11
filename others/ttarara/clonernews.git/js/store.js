// store.js — app state + simple in-memory cache
export const state = {
  feed: 'newstories',
  ids: [],            // ids for current feed (newest->oldest)
  cursor: 0,          // how many ids consumed
  pageSize: 20,
  items: new Map(),   // id -> item (cache)
  users: new Map(),   // id -> user (αν θες μελλοντικά)
  filters: { q: '', minScore: 0 },
};

export function setFeed(feed) {
  state.feed = feed;
  state.ids = [];
  state.cursor = 0;
  state.items.clear();
}

export function setFilterQuery(q) {
  state.filters.q = (q || '').trim().toLowerCase();
}

export function setMinScore(min) {
  state.filters.minScore = Number(min) || 0;
}

export function getItemFromCache(id) {
  return state.items.get(id);
}
