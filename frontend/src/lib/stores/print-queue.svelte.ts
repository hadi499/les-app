import { browser } from '$app/environment';
import type { Card } from '$lib/types';

let queue = $state<Card[]>([]);

if (browser) {
  const saved = localStorage.getItem('print-queue');
  if (saved) {
    try {
      queue = JSON.parse(saved);
    } catch {
      queue = [];
    }
  }
}

function save() {
  if (!browser) return;
  localStorage.setItem('print-queue', JSON.stringify(queue));
}

export const printQueue = {
  get items() { return queue; },
  get count() { return queue.length; },

  has(id: string) { return queue.some(c => c.id === id); },

  /** Toggle: add if not present, remove first occurrence if present */
  toggle(card: Card) {
    const idx = queue.findIndex(c => c.id === card.id);
    if (idx !== -1) {
      queue = [...queue.slice(0, idx), ...queue.slice(idx + 1)];
    } else {
      queue = [...queue, card];
    }
    save();
  },

  /** Remove first occurrence only (support duplicate copies) */
  removeOne(id: string) {
    const idx = queue.findIndex(c => c.id === id);
    if (idx !== -1) {
      queue = [...queue.slice(0, idx), ...queue.slice(idx + 1)];
      save();
    }
  },

  /** Remove all occurrences of a card */
  removeAll(id: string) {
    queue = queue.filter(c => c.id !== id);
    save();
  },

  /** Add a copy even if already in queue (for printing duplicates) */
  pushCopy(card: Card) {
    queue = [...queue, card];
    save();
  },

  clear() {
    queue = [];
    save();
  }
};
