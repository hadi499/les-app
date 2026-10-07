export interface EnglishFlashcardCategory {
  id?: number;
  name: string;
  description: string;
}

export interface EnglishFlashcard {
  id?: number;
  category_id: number;
  category?: EnglishFlashcardCategory;
  question: string;
  options: string[];
  answer: string;
  explanation: string;
}

const API_BASE = `/api/english-flashcards`;
const API_CAT_BASE = `/api/english-flashcard-categories`;

const defaultFetchOpts = {
  credentials: 'include' as RequestCredentials
};

// --- CATEGORIES ---

export async function fetchFlashcardCategories(): Promise<EnglishFlashcardCategory[]> {
  const res = await fetch(API_CAT_BASE, { ...defaultFetchOpts, cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch categories');
  return res.json();
}

export async function createFlashcardCategory(data: EnglishFlashcardCategory): Promise<EnglishFlashcardCategory> {
  const res = await fetch(API_CAT_BASE, {
    ...defaultFetchOpts,
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create category');
  return res.json();
}

export async function updateFlashcardCategory(id: number, data: EnglishFlashcardCategory): Promise<EnglishFlashcardCategory> {
  const res = await fetch(`${API_CAT_BASE}/${id}`, {
    ...defaultFetchOpts,
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update category');
  return res.json();
}

export async function deleteFlashcardCategory(id: number): Promise<void> {
  const res = await fetch(`${API_CAT_BASE}/${id}`, {
    ...defaultFetchOpts,
    method: 'DELETE',
  });
  if (!res.ok) throw new Error('Failed to delete category');
}

// --- FLASHCARDS ---

export async function fetchFlashcards(categoryId?: number): Promise<EnglishFlashcard[]> {
  const url = categoryId ? `${API_BASE}?category_id=${categoryId}` : API_BASE;
  const res = await fetch(url, { ...defaultFetchOpts, cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch flashcards');
  return res.json();
}

export async function createFlashcard(data: EnglishFlashcard): Promise<EnglishFlashcard> {
  const res = await fetch(API_BASE, {
    ...defaultFetchOpts,
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to create flashcard');
  return res.json();
}

export async function updateFlashcard(id: number, data: EnglishFlashcard): Promise<EnglishFlashcard> {
  const res = await fetch(`${API_BASE}/${id}`, {
    ...defaultFetchOpts,
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error('Failed to update flashcard');
  return res.json();
}

export async function deleteFlashcard(id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/${id}`, {
    ...defaultFetchOpts,
    method: 'DELETE',
  });
  if (!res.ok) throw new Error('Failed to delete flashcard');
}
