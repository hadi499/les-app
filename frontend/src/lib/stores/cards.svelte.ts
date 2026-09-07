import { browser } from '$app/environment';
import type { Card, CardQuery } from '$lib/types';
import * as api from '$lib/api';

let items = $state<Card[]>([]);
let total = $state(0);
let totalPages = $state(1);
let currentPage = $state(1);
let currentLimit = $state(10);
let counts = $state<Record<string, number>>({ '6': 0, '4': 0, '2': 0, 'gambar': 0 });
let loading = $state(true);
let error = $state<string | null>(null);
let lastQuery = $state<CardQuery>({});

async function load(query: CardQuery = {}) {
  if (!browser) return;
  const params = { ...query };
  if (!params.limit) params.limit = 1000;
  lastQuery = params;
  loading = true;
  error = null;
  try {
    const res = await api.fetchCards(params);
    items = res.data;
    total = res.total;
    totalPages = res.totalPages;
    currentPage = res.page;
    currentLimit = res.limit;
    counts = res.counts ?? { '6': 0, '4': 0, '2': 0 };
  } catch (e: any) {
    error = e.message || 'Gagal memuat data';
  } finally {
    loading = false;
  }
}

export const cardsStore = {
  get items() { return items; },
  get total() { return total; },
  get totalPages() { return totalPages; },
  get page() { return currentPage; },
  get limit() { return currentLimit; },
  get loading() { return loading; },
  get error() { return error; },
  get counts() { return counts; },

  async fetch(query?: CardQuery) {
    await load(query ?? lastQuery);
  },

  async add(card: Omit<Card, 'id'>) {
    await api.createCard(card);
    await load(lastQuery);
  },

  async update(id: string, data: Partial<Omit<Card, 'id'>>) {
    const existing = items.find(c => c.id === id);
    if (!existing) return;
    await api.updateCard(id, { ...existing, ...data });
    await load(lastQuery);
  },

  async remove(id: string) {
    await api.deleteCard(id);
    await load(lastQuery);
  },

  async uploadImage(file: File): Promise<string> {
    return api.uploadImage(file);
  }
};
