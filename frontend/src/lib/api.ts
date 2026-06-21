import type { Card, PaginatedCards, CardQuery } from '$lib/types';

const API_BASE = `/api`;

const defaultFetchOpts = {
  credentials: 'include' as RequestCredentials
};

export async function fetchCards(query: CardQuery = {}): Promise<PaginatedCards> {
  const params = new URLSearchParams();
  if (query.page) params.set('page', String(query.page));
  if (query.limit) params.set('limit', String(query.limit));
  if (query.size) params.set('size', query.size);
  if (query.search) params.set('search', query.search);
  if (query.cardType) params.set('cardType', query.cardType);
  if (query.all) params.set('all', 'true');

  const url = `${API_BASE}/cards${params.toString() ? '?' + params.toString() : ''}`;
  const res = await fetch(url, defaultFetchOpts);
  if (!res.ok) throw new Error('Failed to fetch cards');
  const json = await res.json();
  return {
    data: json.data.map((c: any) => ({ ...c, id: String(c.ID), size: String(c.size ?? '') })),
    total: json.total,
    page: json.page,
    limit: json.limit,
    totalPages: json.totalPages,
    counts: json.counts ?? { all: 0, '6': 0, '4': 0, '2': 0, 'gambar': 0 },
  };
}

export async function createCard(card: Omit<Card, 'id'>): Promise<Card> {
  const res = await fetch(`${API_BASE}/cards`, {
    ...defaultFetchOpts,
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(card),
  });
  if (!res.ok) throw new Error('Failed to create card');
  const data = await res.json();
  return { ...data, id: String(data.ID), size: String(data.size ?? '') };
}

export async function updateCard(id: string, card: Omit<Card, 'id'>): Promise<Card> {
  const res = await fetch(`${API_BASE}/cards/${id}`, {
    ...defaultFetchOpts,
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(card),
  });
  if (!res.ok) throw new Error('Failed to update card');
  const data = await res.json();
  return { ...data, id: String(data.ID), size: String(data.size ?? '') };
}

export async function deleteCard(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/cards/${id}`, { ...defaultFetchOpts, method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to delete card');
}

export async function fetchTrashCards(): Promise<{ data: Card[]; total: number }> {
  const res = await fetch(`${API_BASE}/cards/trash`, defaultFetchOpts);
  if (!res.ok) throw new Error('Failed to fetch trash');
  const json = await res.json();
  return {
    data: json.data.map((c: any) => ({ ...c, id: String(c.ID), size: String(c.size ?? '') })),
    total: json.total,
  };
}

export async function restoreCard(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/cards/trash/${id}/restore`, { ...defaultFetchOpts, method: 'POST' });
  if (!res.ok) throw new Error('Failed to restore card');
}

export async function forceDeleteCard(id: string): Promise<void> {
  const res = await fetch(`${API_BASE}/cards/trash/${id}/force`, { ...defaultFetchOpts, method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to permanently delete card');
}

export async function fetchTrashCount(): Promise<number> {
  const res = await fetch(`${API_BASE}/cards/trash`, defaultFetchOpts);
  if (!res.ok) throw new Error('Failed to fetch trash count');
  const json = await res.json();
  return json.total ?? 0;
}

export async function emptyTrash(): Promise<void> {
  const res = await fetch(`${API_BASE}/cards/trash/empty`, { ...defaultFetchOpts, method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to empty trash');
}

export async function uploadImage(file: File): Promise<string> {
  const formData = new FormData();
  formData.append('image', file);
  const res = await fetch(`${API_BASE}/upload`, {
    ...defaultFetchOpts,
    method: 'POST',
    body: formData,
  });
  if (!res.ok) throw new Error('Failed to upload image');
  const data = await res.json();
  return data.url;
}

export async function fetchImages(): Promise<{ name: string; url: string }[]> {
  const res = await fetch(`${API_BASE}/images`, defaultFetchOpts);
  if (!res.ok) throw new Error('Failed to fetch images');
  const data = await res.json();
  return data.images ?? [];
}
