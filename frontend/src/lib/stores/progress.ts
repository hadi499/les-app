import { writable, get } from 'svelte/store';
import type { LessonProgress } from '$lib/types';

export interface GameHistory {
    id: number;
    mode: string;
    score: number;
    created_at: string;
}

export interface LessonHistory {
    id: number;
    lessonId: number;
    wpm: number;
    accuracy: number;
    stars: number;
    created_at: string;
}

export const progressMap = writable<Map<number, LessonProgress>>(new Map());
export const gameScores = writable({ left: 0, right: 0, both: 0, letters: 0, all: 0 });
export const gameHistoryStore = writable<GameHistory[]>([]);
export const lessonHistoryStore = writable<LessonHistory[]>([]);

export async function fetchAllProgress() {
    try {
        const res = await fetch(`/api/typing/progress`, { credentials: 'include' });
        if (res.ok) {
            const data: LessonProgress[] = await res.json();
            const map = new Map<number, LessonProgress>();
            for (const p of data) {
                map.set(p.lessonId, p);
            }
            progressMap.set(map);
        }
    } catch(e) { console.error(e); }
}

export async function fetchAllGameScores() {
    try {
        const res = await fetch(`/api/typing/game-scores`, { credentials: 'include' });
        if (res.ok) {
            const data = await res.json();
            gameScores.set(data);
        }
    } catch(e) { console.error(e); }
}

export async function fetchHistory() {
    try {
        const resGame = await fetch(`/api/typing/history/game`, { credentials: 'include' });
        if (resGame.ok) {
            gameHistoryStore.set(await resGame.json());
        }
        
        const resLesson = await fetch(`/api/typing/history/lesson`, { credentials: 'include' });
        if (resLesson.ok) {
            lessonHistoryStore.set(await resLesson.json());
        }
    } catch(e) { console.error(e); }
}

export function getProgress(lessonId: number): LessonProgress {
    const map = get(progressMap);
    return map.get(lessonId) ?? {
        lessonId, bestWPM: 0, bestAccuracy: 0, completed: false, stars: 0, attempts: 0
    };
}

export function getAllProgress(): LessonProgress[] {
    const map = get(progressMap);
    return Array.from(map.values()).sort((a, b) => a.lessonId - b.lessonId);
}

export async function saveLessonProgress(progress: LessonProgress) {
    try {
        const res = await fetch(`/api/typing/progress`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(progress),
            credentials: 'include'
        });
        if (res.ok) {
            const updated = await res.json();
            progressMap.update(map => {
                map.set(updated.lessonId, updated);
                return map;
            });
        }
    } catch(e) { console.error(e); }
}

export function calculateStars(wpm: number, accuracy: number): number {
    if (accuracy >= 95 && wpm >= 25) return 3;
    if (accuracy >= 85 && wpm >= 15) return 2;
    if (accuracy >= 70 && wpm >= 10) return 1;
    return 0;
}

export function getTotalStars(): number {
    const map = get(progressMap);
    return Array.from(map.values()).reduce((sum, p) => sum + p.stars, 0);
}

export function getMaxStars(): number {
    return 36;
}
