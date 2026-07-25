<script lang="ts">
  import { onMount } from 'svelte';
  import CanvasStage from './CanvasStage.svelte';
  import type { LevelDef, Position } from './levels';

  let { level, blockSequence, onNextLevel } = $props<{
    level: LevelDef;
    blockSequence: any[];
    onNextLevel: () => void;
  }>();

  let charPos = $state<Position>({ x: 0, y: 0 });
  let charDir = $state<'right' | 'left' | 'up' | 'down'>('right');
  let status = $state<'idle' | 'running' | 'success' | 'fail'>('idle');

  // Reset to initial state when level changes
  $effect(() => {
    level; // dependency
    reset();
  });

  function reset() {
    charPos = { ...level.start };
    charDir = level.startDir;
    status = 'idle';
  }

  // Flatten the block sequence (unroll loops) to make stepping simple
  function flattenSequence(seq: any[]): any[] {
    const result: any[] = [];
    for (const block of seq) {
      if (block.type === 'repeat_n') {
        const times = block.times || 1;
        const body = flattenSequence(block.do || []);
        for (let i = 0; i < times; i++) {
          result.push(...body);
        }
      } else {
        result.push(block);
      }
    }
    return result;
  }

  // Delay helper
  const delay = (ms: number) => new Promise(res => setTimeout(res, ms));

  function isWall(x: number, y: number) {
    return level.walls.some((w: Position) => w.x === x && w.y === y);
  }

  function isValidMove(x: number, y: number) {
    if (x < 0 || x >= level.gridSize.cols) return false;
    if (y < 0 || y >= level.gridSize.rows) return false;
    if (isWall(x, y)) return false;
    return true;
  }

  async function run() {
    if (status === 'running') return;
    
    reset();
    status = 'running';

    const instructions = flattenSequence(blockSequence);
    
    // Add small initial delay
    await delay(300);

    for (const inst of instructions) {
      if (status !== 'running') break; // stopped

      if (inst.type === 'move_forward') {
        let newX = charPos.x;
        let newY = charPos.y;
        if (charDir === 'right') newX += 1;
        else if (charDir === 'left') newX -= 1;
        else if (charDir === 'up') newY -= 1;
        else if (charDir === 'down') newY += 1;

        if (isValidMove(newX, newY)) {
          charPos = { x: newX, y: newY };
        } else {
          // Hit wall or boundary
          status = 'fail';
          return;
        }
      } else if (inst.type === 'turn_left') {
        if (charDir === 'right') charDir = 'up';
        else if (charDir === 'up') charDir = 'left';
        else if (charDir === 'left') charDir = 'down';
        else if (charDir === 'down') charDir = 'right';
      } else if (inst.type === 'turn_right') {
        if (charDir === 'right') charDir = 'down';
        else if (charDir === 'down') charDir = 'left';
        else if (charDir === 'left') charDir = 'up';
        else if (charDir === 'up') charDir = 'right';
      }

      await delay(500); // Wait 500ms between moves to show animation
    }

    if (status === 'running') {
      // Finished all instructions. Did we reach the goal?
      if (charPos.x === level.goal.x && charPos.y === level.goal.y) {
        status = 'success';
      } else {
        status = 'fail';
      }
    }
  }
</script>

<div class="flex flex-col h-full bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
  <!-- Controls -->
  <div class="p-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
    <h2 class="font-bold text-slate-800">{level.title}</h2>
    <div class="flex gap-2">
      <button 
        class="px-4 py-2 bg-slate-200 hover:bg-slate-300 text-slate-700 font-semibold rounded-xl transition flex items-center gap-2"
        onclick={reset}
        disabled={status === 'running'}
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
        Reset
      </button>
      <button 
        class="px-6 py-2 bg-green-500 hover:bg-green-600 text-white font-bold rounded-xl shadow-sm shadow-green-500/30 transition flex items-center gap-2"
        onclick={run}
        disabled={status === 'running'}
      >
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd"></path></svg>
        Jalankan
      </button>
    </div>
  </div>

  <!-- Stage -->
  <div class="flex-1 p-6 flex flex-col items-center justify-center bg-slate-50 relative">
    <CanvasStage 
      level={level}
      charPos={charPos}
      charDir={charDir}
      status={status}
    />

    {#if status === 'success'}
      <div class="mt-8 animate-in slide-in-from-bottom-4">
        <button 
          class="px-8 py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold text-lg rounded-xl shadow-lg shadow-blue-600/30 transition flex items-center gap-2"
          onclick={onNextLevel}
        >
          Lanjut ke Level Berikutnya
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7"></path></svg>
        </button>
      </div>
    {/if}
  </div>
</div>
