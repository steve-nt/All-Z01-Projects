// live.js — Αυτόματη ένθεση νέων items κάθε 5s + ενημερωτικό banner
import { getUpdates, getItem } from './api.js';
import { state } from './store.js';
import { prependItems } from './ui.js';

const banner = document.getElementById('live-banner');
const text = document.getElementById('live-text');

let timer = null;

export function startLive(){
  stopLive();
  banner.hidden = true;               // θα προβάλλεται μόνο όταν μπαίνουν νέα
  timer = setInterval(checkAndInsert, 5000);
}

export function stopLive(){
  if (timer){ clearInterval(timer); timer = null; }
}

async function checkAndInsert(){
  try{
    const { items = [] } = await getUpdates();
    if (!Array.isArray(items) || items.length === 0) return;

    // ids που δεν έχουμε ήδη στη μνήμη
    const unseen = items.filter(id => !state.items.has(id));
    if (!unseen.length) return;

    // Κατέβασε μέχρι 20 για να είμαστε ευγενικοί με το API
    const fetched = (await Promise.allSettled(unseen.slice(0, 20).map(getItem)))
      .filter(r => r.status === 'fulfilled' && r.value)
      .map(r => r.value)
      .sort((a,b)=> (b?.time||0) - (a?.time||0));

    if (!fetched.length) return;

    // cache inline και ένθεση στην κορυφή (UI filters ισχύουν)
    for (const it of fetched) if (it && it.id != null) state.items.set(it.id, it);
    prependItems(fetched);

    // Εμφάνισε banner για λίγο με το πλήθος
    text.textContent = `+${fetched.length} new item(s)`;
    banner.hidden = false;
    setTimeout(()=>{ banner.hidden = true; }, 3500);
  }catch(e){
    // αθόρυβο: best-effort live
  }
}
