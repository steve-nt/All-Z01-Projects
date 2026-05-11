// ui.js — DOM builders & renderers (append + prepend)
import { state, getItemFromCache } from './store.js';
import { getItem, throttleBurst } from './api.js';

const feedEl = document.getElementById('feed');
const loaderEl = document.getElementById('loader');
const detailsPanel = document.getElementById('details');
const detailsArticle = document.getElementById('details-article');
const commentsEl = document.getElementById('comments');

export function setLoading(v){ loaderEl.hidden = !v; }
export function clearFeed(){ feedEl.innerHTML = ''; }

export function appendItems(items){ renderItems(items, false); }
export function prependItems(items){ renderItems(items, true); }

function renderItems(items, prepend){
  const frag = document.createDocumentFragment();
  let count = 0;
  for (const it of items){
    if(!it) continue;
    // cache inline
    if (it && it.id != null) state.items.set(it.id, it);
    if (!passesFilters(it)) continue;
    frag.appendChild(renderListItem(it));
    count++;
  }
  if (count === 0) return;
  if (prepend && feedEl.firstChild) {
    feedEl.insertBefore(frag, feedEl.firstChild);
  } else {
    feedEl.appendChild(frag);
  }
}

function passesFilters(it){
  const q = state.filters.q;
  const min = state.filters.minScore;
  const title = (it.title || '').toLowerCase();
  const okTitle = q ? title.includes(q) : true;
  const score = typeof it.score === 'number' ? it.score : 0;
  const okScore = score >= min;
  const okType = it.type !== 'comment'; // comments στη λίστα τα αγνοούμε
  return okTitle && okScore && okType;
}

function renderListItem(it){
  const li = document.createElement('li');
  li.tabIndex = 0;
  li.setAttribute('role','article');

  const left = document.createElement('div');
  const right = document.createElement('div');

  const a = document.createElement('a');
  a.className = 'item-title link';
  a.textContent = it.title || '(no title)';
  a.href = it.url ? it.url : `https://news.ycombinator.com/item?id=${it.id}`;
  a.target = '_blank';
  a.rel = 'noreferrer';

  const meta = document.createElement('div');
  meta.className = 'item-meta';
  const time = new Date((it.time || 0) * 1000);
  meta.textContent = [
    it.by ? `by ${it.by}` : '',
    time ? time.toLocaleString() : '',
    typeof it.score === 'number' ? `• ${it.score} points` : '',
    typeof it.descendants === 'number' ? `• ${it.descendants} comments` : '',
  ].filter(Boolean).join('  •  ');

  left.appendChild(a);
  left.appendChild(meta);

  const openBtn = document.createElement('button');
  openBtn.className = 'badge';
  openBtn.textContent = 'Open';
  openBtn.addEventListener('click', () => openDetails(it.id));

  right.appendChild(openBtn);

  li.appendChild(left);
  li.appendChild(right);
  li.addEventListener('keypress', (e)=>{ if(e.key==='Enter') openDetails(it.id); });
  return li;
}

export async function openDetails(id){
  detailsPanel.hidden = false;
  detailsArticle.innerHTML = '<h2>Loading…</h2>';
  commentsEl.innerHTML = '';

  let item = getItemFromCache(id);
  if (!item) { item = await getItem(id); if(item && item.id) state.items.set(item.id, item); }

  detailsArticle.innerHTML = `
    <h2>${escapeHtml(item.title || '(no title)')}</h2>
    <div class="kv">
      <span>by ${escapeHtml(item.by || 'unknown')}</span>
      <span>•</span>
      <span>${new Date((item.time||0)*1000).toLocaleString()}</span>
      ${typeof item.score === 'number' ? `<span>•</span><span>score ${item.score}</span>` : ''}
    </div>
    ${item.text ? `<p>${item.text}</p>` : ''}
    ${item.url ? `<p><a class="link" href="${item.url}" target="_blank" rel="noreferrer">Open original ↗</a></p>` : ''}
  `;

  const kids = Array.isArray(item.kids) ? item.kids.slice(0, 200) : [];
  if (kids.length){
    const tree = await loadCommentTree(kids);
    commentsEl.appendChild(tree);
  } else {
    commentsEl.innerHTML = '<p class="item-meta">No comments.</p>';
  }
}

async function loadCommentTree(ids){
  const items = await throttleBurst(ids, 25, 200);
  const byId = new Map(items.filter(Boolean).map(x=>[x.id, x]));
  const frag = document.createDocumentFragment();
  for (const id of ids){
    const node = buildCommentNode(byId.get(id), byId);
    if (node) frag.appendChild(node);
  }
  return frag;
}

function buildCommentNode(c, byId){
  if (!c || c.deleted || c.dead) return null;
  const el = document.createElement('div');
  el.className = 'comment';
  el.innerHTML = `
    <div class="meta">${escapeHtml(c.by || 'unknown')} • ${new Date((c.time||0)*1000).toLocaleString()}</div>
    <div class="text">${c.text || ''}</div>
  `;
  if (Array.isArray(c.kids) && c.kids.length){
    const childFrag = document.createDocumentFragment();
    for (const kid of c.kids.slice(0, 50)){
      childFrag.appendChild(buildCommentNode(byId.get(kid), byId) || document.createTextNode(''));
    }
    el.appendChild(childFrag);
  }
  return el;
}

export function closeDetails(){
  detailsPanel.hidden = true;
  detailsArticle.innerHTML = '';
  commentsEl.innerHTML = '';
}

function escapeHtml(s){
  return String(s).replace(/[&<>"']/g, m => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[m]));
}
