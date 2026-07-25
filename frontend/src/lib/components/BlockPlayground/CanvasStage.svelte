<script lang="ts">
  import type { LevelDef, Position } from './levels';

  let { 
    level, 
    charPos, 
    charDir, 
    status // 'idle' | 'running' | 'success' | 'fail'
  } = $props<{
    level: LevelDef;
    charPos: Position;
    charDir: 'right' | 'left' | 'up' | 'down';
    status: 'idle' | 'running' | 'success' | 'fail';
  }>();

  // Helper to check if a tile is a wall
  function isWall(x: number, y: number) {
    return level.walls.some((w: Position) => w.x === x && w.y === y);
  }
</script>

<div class="relative w-full aspect-square max-w-md mx-auto bg-green-50 border-4 border-green-200 rounded-xl overflow-hidden shadow-inner p-2">
  <div 
    class="w-full h-full grid gap-1"
    style="grid-template-columns: repeat({level.gridSize.cols}, minmax(0, 1fr)); grid-template-rows: repeat({level.gridSize.rows}, minmax(0, 1fr));"
  >
    {#each Array(level.gridSize.rows) as _, y}
      {#each Array(level.gridSize.cols) as _, x}
        <div class="relative bg-white/60 rounded-lg border border-green-100 flex items-center justify-center overflow-hidden">
          
          <!-- Wall -->
          {#if isWall(x, y)}
            <div class="absolute inset-0 bg-slate-700/80 m-1 rounded bg-[url('data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyMCIgaGVpZ2h0PSIyMCI+PHBhdGggZD0iTTAgMGgyMHYyMEgwem0xMCAwaDEwdjEwSDEwem0wIDEwaDEwdjEwSDEwIiBmaWxsPSJyZ2JhKDI1NSwyNTUsMjU1LDAuMSkiLz48L3N2Zz4=')]"></div>
          {/if}

          <!-- Goal -->
          {#if level.goal.x === x && level.goal.y === y}
            <div class="absolute w-3/4 h-3/4 animate-bounce">
              <svg viewBox="0 0 24 24" fill="#facc15" stroke="#ca8a04" stroke-width="2" class="w-full h-full drop-shadow-md">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
              </svg>
            </div>
          {/if}

          <!-- Character -->
          {#if charPos.x === x && charPos.y === y}
            <div 
              class="absolute w-3/4 h-3/4 transition-all duration-300 z-10 flex items-center justify-center"
              style="
                transform: rotate({
                  charDir === 'right' ? '0deg' : 
                  charDir === 'down' ? '90deg' : 
                  charDir === 'left' ? '180deg' : 
                  '-90deg'
                });
              "
            >
              <div class="w-full h-full bg-blue-500 rounded-full flex items-center justify-end p-1 shadow-[0_4px_10px_rgba(59,130,246,0.5)]">
                <!-- Eyes to show direction -->
                <div class="w-1/4 h-1/4 bg-white rounded-full mr-1 mb-2"></div>
                <div class="w-1/4 h-1/4 bg-white rounded-full mr-1 mt-2 absolute right-1 bottom-1"></div>
                <!-- simple triangle beak/nose -->
                <div class="absolute -right-1 w-2 h-2 bg-yellow-400 rotate-45"></div>
              </div>
            </div>
          {/if}

        </div>
      {/each}
    {/each}
  </div>

  <!-- Status Overlay -->
  {#if status === 'success'}
    <div class="absolute inset-0 bg-white/80 backdrop-blur-sm flex flex-col items-center justify-center z-50 animate-in fade-in duration-300">
      <div class="text-5xl mb-4">🌟</div>
      <h3 class="text-2xl font-bold text-green-600">Hebat!</h3>
      <p class="text-slate-600 font-medium">Kamu berhasil menyelesaikan level ini.</p>
    </div>
  {:else if status === 'fail'}
    <div class="absolute inset-0 bg-white/80 backdrop-blur-sm flex flex-col items-center justify-center z-50 animate-in fade-in duration-300">
      <div class="text-5xl mb-4">😅</div>
      <h3 class="text-2xl font-bold text-orange-500">Ups, hampir!</h3>
      <p class="text-slate-600 font-medium">Ayo periksa kodenya dan coba lagi.</p>
    </div>
  {/if}
</div>
