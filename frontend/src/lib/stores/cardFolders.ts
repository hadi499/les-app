import { writable } from 'svelte/store';
import { env } from '$env/dynamic/public';

const API_URL = env.PUBLIC_API_URL || 'http://localhost:8080';

export interface CardFolder {
    id: number;
    name: string;
}

export const cardFolders = writable<CardFolder[]>([]);

export async function fetchCardFolders() {
    try {
        const res = await fetch(`${API_URL}/api/card-folders`, {
            credentials: 'include'
        });
        if (res.ok) {
            const data = await res.json();
            cardFolders.set(data.data || []);
        }
    } catch (err) {
        console.error('Failed to fetch card folders:', err);
    }
}

export async function createCardFolder(name: string) {
    try {
        const res = await fetch(`${API_URL}/api/card-folders`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name }),
            credentials: 'include'
        });
        if (res.ok) {
            await fetchCardFolders();
            return true;
        }
    } catch (err) {
        console.error('Failed to create card folder:', err);
    }
    return false;
}

export async function updateCardFolder(id: number, name: string) {
    try {
        const res = await fetch(`${API_URL}/api/card-folders/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name }),
            credentials: 'include'
        });
        if (res.ok) {
            await fetchCardFolders();
            return true;
        }
    } catch (err) {
        console.error('Failed to update card folder:', err);
    }
    return false;
}

export async function deleteCardFolder(id: number) {
    try {
        const res = await fetch(`${API_URL}/api/card-folders/${id}`, {
            method: 'DELETE',
            credentials: 'include'
        });
        if (res.ok) {
            await fetchCardFolders();
            return true;
        }
    } catch (err) {
        console.error('Failed to delete card folder:', err);
    }
    return false;
}
