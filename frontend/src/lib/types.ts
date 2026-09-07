export type Finger = 'left-pinky' | 'left-ring' | 'left-middle' | 'left-index' | 'right-index' | 'right-middle' | 'right-ring' | 'right-pinky' | 'thumb';

export type Hand = 'left' | 'right';

export interface KeyInfo {
	key: string;
	finger: Finger;
	hand: Hand;
	row: 'number' | 'top' | 'home' | 'bottom';
}

export interface Lesson {
	id: number;
	title: string;
	description: string;
	icon: string;
	keys: string[];
	words: string[];
	sentences: string[];
	unlockScore: number;
}

export interface LessonProgress {
	lessonId: number;
	bestWPM: number;
	bestAccuracy: number;
	completed: boolean;
	stars: number;
	attempts: number;
}

export interface TypingStats {
	wpm: number;
	accuracy: number;
	correctChars: number;
	totalChars: number;
	errors: number;
	time: number;
}

export interface Card {
  id: string;
  card_folder_id?: number | null;
  card_folder?: { id: number; name: string };
  image: string;
  title: string;
  content: string;
  size: string;
  cardType: string;
}

export interface PaginatedCards {
  data: Card[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
  counts: Record<string, number>;
}

export interface CardQuery {
  page?: number;
  limit?: number;
  size?: string;
  search?: string;
  cardType?: string;
  all?: boolean;
}
