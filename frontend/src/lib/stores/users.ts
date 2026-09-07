import { fetchAllGameScores, gameScores } from './progress';
import { get } from 'svelte/store';

// Legacy compatibility: we no longer manage local users. The current user is managed by the Go backend via HTTP cookies.
// We just provide empty implementations for legacy functions to prevent breaking changes in UI components.
// If a component actually needs the real user, it should call /me API or rely on the layout wrapper.

export function getCurrentUserId(): string | null {
	return 'authenticated'; // dummy to satisfy components checking if logged in
}

export function getGameHighScores(userId?: string) {
	return get(gameScores);
}

export async function saveGameHighScore(userId: string, mode: 'left' | 'right' | 'both' | 'letters' | 'all', score: number) {
	try {
        const res = await fetch(`/api/typing/game-scores`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ mode, score }),
            credentials: 'include'
        });
        if (res.ok) {
            await fetchAllGameScores();
        }
    } catch(e) { console.error(e); }
}

export function getUsers() { return []; }
export function createUser() { return { id: 'dummy', name: 'User', avatar: '🧑‍💻', createdAt: Date.now() }; }
export function setCurrentUserId() {}
