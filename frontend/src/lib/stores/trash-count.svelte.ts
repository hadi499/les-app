import { browser } from '$app/environment';
import * as api from '$lib/api';

let count = $state(0);
let loaded = false;

async function load() {
  if (!browser || loaded) return;
  loaded = true;
  try {
    const res = await api.fetchTrashCount();
    count = res;
  } catch {}
}

export const trashCount = {
  get value() { return count; },
  set(n: number) { count = n; },
  inc() { count++; },
  dec() { count = Math.max(0, count - 1); },
  init() { load(); },
};
