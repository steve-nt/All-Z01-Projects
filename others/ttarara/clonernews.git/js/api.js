// api.js — network layer for HN Firebase API (με fan-out και ήπιο throttling)
const BASE = 'https://hacker-news.firebaseio.com/v0';

const sleep = (ms) => new Promise(r => setTimeout(r, ms));

export async function throttleBurst(ids, chunk = 25, gapMs = 250) {
  const buckets = [];
  for (let i = 0; i < ids.length; i += chunk) buckets.push(ids.slice(i, i + chunk));
  const out = [];
  for (let bi = 0; bi < buckets.length; bi++) {
    const b = buckets[bi];
    const results = await Promise.all(b.map(getItem));
    out.push(...results.filter(Boolean));
    if (bi < buckets.length - 1) await sleep(gapMs);
  }
  return out;
}

export async function getIds(feed /* 'newstories' | 'jobstories' | 'pollstories' */) {
  const res = await fetch(`${BASE}/${feed}.json`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to load ids');
  return res.json();
}

export async function getItem(id) {
  const res = await fetch(`${BASE}/item/${id}.json`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to load item ' + id);
  return res.json();
}

export async function getUpdates() {
  const res = await fetch(`${BASE}/updates.json`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to load updates');
  return res.json(); // { items: number[], profiles: string[] }
}
