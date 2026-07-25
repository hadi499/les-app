export interface Position {
  x: number;
  y: number;
}

export interface LevelDef {
  id: number;
  title: string;
  gridSize: { cols: number; rows: number };
  start: Position;
  startDir: 'right' | 'left' | 'up' | 'down';
  goal: Position;
  walls: Position[];
  allowedBlocks: string[]; // e.g. ['move_forward', 'turn_left', 'turn_right']
}

export const levels: Record<number, LevelDef> = {
  1: {
    id: 1,
    title: 'Level 1: Maju ke Bintang',
    gridSize: { cols: 5, rows: 1 },
    start: { x: 0, y: 0 },
    startDir: 'right',
    goal: { x: 4, y: 0 },
    walls: [],
    allowedBlocks: ['move_forward']
  },
  2: {
    id: 2,
    title: 'Level 2: Belok Kiri',
    gridSize: { cols: 5, rows: 5 },
    start: { x: 0, y: 4 },
    startDir: 'right',
    goal: { x: 4, y: 0 },
    walls: [
      { x: 1, y: 4 }, { x: 2, y: 4 }, { x: 3, y: 4 }, { x: 4, y: 4 }, // block bottom row
      { x: 1, y: 3 }, { x: 2, y: 3 }, { x: 3, y: 3 }, { x: 4, y: 3 }, // block next row
      { x: 1, y: 2 }, { x: 2, y: 2 }, { x: 3, y: 2 }, { x: 4, y: 2 },
      { x: 1, y: 1 }, { x: 2, y: 1 }, { x: 3, y: 1 }, { x: 4, y: 1 }
    ], // Essentially just an L-shape path: move up, then right. Actually start is (0,4) facing right.
    // Wait, (0,4) facing right, but wall is at (1,4). So they must turn left immediately.
    allowedBlocks: ['move_forward', 'turn_left', 'turn_right']
  },
  3: {
    id: 3,
    title: 'Level 3: Labirin Kecil',
    gridSize: { cols: 5, rows: 5 },
    start: { x: 0, y: 0 },
    startDir: 'right',
    goal: { x: 4, y: 4 },
    walls: [
      { x: 1, y: 0 },
      { x: 1, y: 1 },
      { x: 3, y: 4 },
      { x: 3, y: 3 },
      { x: 3, y: 2 },
    ],
    allowedBlocks: ['move_forward', 'turn_left', 'turn_right']
  },
  4: {
    id: 4,
    title: 'Level 4: Pengulangan (Loop)',
    gridSize: { cols: 6, rows: 1 },
    start: { x: 0, y: 0 },
    startDir: 'right',
    goal: { x: 5, y: 0 },
    walls: [],
    allowedBlocks: ['move_forward', 'repeat_n'] // Introduce loops
  }
};
