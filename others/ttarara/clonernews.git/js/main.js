// main.js — glue code
import { getIds, throttleBurst } from './api.js';
import { state, setFeed, setFilterQuery, setMinScore } from './store.js';
import { setLoading, appendItems, clearFeed, closeDetails } from './ui.js';
import { startLive } from './live.js';

const tabs = document.querySelectorAll('.tab');
const search = document.getElementById('search');
const onlyTop = document.getElementById('only-top');
const loadMoreBtn = document.getElementById('load-more');
const loaderEl = document.getElementById('loader');
const closeDetailsBtn = document.getElementById('close-details');

bootstrap();

function bootstrap(){
  tabs.forEach(btn => btn.addEventListener('click', onTab));
  closeDetailsBtn.addEventListener('click', closeDetails);

  const debounced = debounce((v)=>{ setFilterQuery(v); reRenderVisible(); }, 300);
  search.addEventListener('input', (e)=> debounced(e.target.value));
  onlyTop.addEventListener('change', (e)=> { setMinScore(e.target.checked ? 100 : 0); reRenderVisible(); });

  loadMoreBtn.addEventListener('click', loadNextPage);
  window.addEventListener('scroll', onScrollNearBottom);

  // αρχικό feed
  onTabClick('newstories');

  // Live (auto updates κάθε 5″)
  startLive();
}

function onTab(e){
  const feed = e.currentTarget.dataset.feed;
  onTabClick(feed);
}

async function onTabClick(feed){
  tabs.forEach(b => b.classList.toggle('is-active', b.dataset.feed === feed));
  setFeed(feed);
  search.value = '';
  onlyTop.checked = false;
  setFilterQuery('');
  setMinScore(0);

  clearFeed();
  loaderEl.hidden = false;
  const ids = await getIds(feed);
  state.ids = ids;
  state.cursor = 0;
  await loadNextPage();
  loaderEl.hidden = true;
}

async function loadNextPage(){
  if (state.cursor >= state.ids.length) return;
  setLoading(true);
  const slice = state.ids.slice(state.cursor, state.cursor + state.pageSize);
  const items = await throttleBurst(slice, 25, 200);
  items.sort((a,b)=> (b?.time||0) - (a?.time||0));
  appendItems(items);
  state.cursor += slice.length;
  setLoading(false);
}

function onScrollNearBottom(){
  const near = window.innerHeight + window.scrollY >= (document.body.offsetHeight - 1200);
  if (near) loadNextPage();
}

function reRenderVisible(){
  const cached = Array.from(state.items.values());
  cached.sort((a,b)=> (b?.time||0) - (a?.time||0));
  clearFeed();
  appendItems(cached);
}

function debounce(fn, ms){
  let t = 0;
  return (...args) => { clearTimeout(t); t = setTimeout(()=> fn(...args), ms); };
}
